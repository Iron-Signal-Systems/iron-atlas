package httpui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/change"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/health"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/modules"
)

type Dependencies struct {
	Logger         *slog.Logger
	Policy         *authz.Policy
	Changes        change.Service
	Modules        modules.Registry
	Readiness      health.Checker
	Authentication *authentication.Middleware
}

type Server struct {
	deps      Dependencies
	templates *template.Template
	mux       *http.ServeMux
}

func New(deps Dependencies) (*Server, error) {
	if deps.Logger == nil ||
		deps.Policy == nil ||
		deps.Changes == nil ||
		deps.Readiness == nil ||
		deps.Authentication == nil {
		return nil, errors.New(
			"logger, policy, change service, readiness checker, and " +
				"authentication middleware are required",
		)
	}
	templates, err := template.ParseFS(webFiles, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	server := &Server{
		deps:      deps,
		templates: templates,
		mux:       http.NewServeMux(),
	}
	server.routes()
	return server, nil
}

func (s *Server) routes() {
	staticFS, _ := fs.Sub(webFiles, "static")
	s.mux.Handle(
		"GET /static/",
		http.StripPrefix(
			"/static/",
			http.FileServer(http.FS(staticFS)),
		),
	)
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /readyz", s.ready)
	s.mux.HandleFunc("GET /api/v1/status", s.apiStatus)
	s.mux.HandleFunc("GET /api/v1/changes", s.apiChanges)
	s.mux.HandleFunc(
		"POST /api/v1/changes/{id}/approve",
		s.apiApprove,
	)
	s.mux.HandleFunc("GET /changes", s.changesPage)
	s.mux.HandleFunc("GET /modules", s.modulesPage)
	s.mux.HandleFunc("GET /", s.dashboard)
}

func (s *Server) Handler() http.Handler {
	authenticated := s.deps.Authentication.Handler(s.mux)
	return requestLog(
		s.deps.Logger,
		securityHeaders(authenticated),
	)
}

func (s *Server) actor(
	w http.ResponseWriter,
	request *http.Request,
) (authz.Actor, bool) {
	actor, ok := authentication.ActorFromContext(
		request.Context(),
	)
	if !ok {
		writeJSON(
			w,
			http.StatusUnauthorized,
			map[string]string{"error": "authentication required"},
		)
		return authz.Actor{}, false
	}
	return actor, true
}

func (s *Server) dashboard(
	w http.ResponseWriter,
	request *http.Request,
) {
	if request.URL.Path != "/" {
		http.NotFound(w, request)
		return
	}
	actor, ok := s.actor(w, request)
	if !ok {
		return
	}
	changes, err := s.deps.Changes.List(request.Context())
	if err != nil {
		s.dependencyFailure(w, "list dashboard changes", err)
		return
	}
	data := pageData{
		Title:   "Dashboard",
		Actor:   actor,
		Changes: changes,
		Modules: s.deps.Modules.List(),
		Now:     time.Now().UTC(),
		DevelopmentMode: s.deps.Authentication.Mode() ==
			authentication.ModeDevelopment,
	}
	s.render(w, "dashboard.html", data)
}

func (s *Server) changesPage(
	w http.ResponseWriter,
	request *http.Request,
) {
	actor, ok := s.actor(w, request)
	if !ok {
		return
	}
	changes, err := s.deps.Changes.List(request.Context())
	if err != nil {
		s.dependencyFailure(w, "list change page", err)
		return
	}
	s.render(w, "changes.html", pageData{
		Title:   "Change management",
		Actor:   actor,
		Changes: changes,
		DevelopmentMode: s.deps.Authentication.Mode() ==
			authentication.ModeDevelopment,
	})
}

func (s *Server) modulesPage(
	w http.ResponseWriter,
	request *http.Request,
) {
	actor, ok := s.actor(w, request)
	if !ok {
		return
	}
	s.render(w, "modules.html", pageData{
		Title:   "Modules",
		Actor:   actor,
		Modules: s.deps.Modules.List(),
		DevelopmentMode: s.deps.Authentication.Mode() ==
			authentication.ModeDevelopment,
	})
}

func (s *Server) render(
	w http.ResponseWriter,
	name string,
	data pageData,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		s.deps.Logger.Error(
			"template rendering failed",
			"template",
			name,
			"error",
			err,
		)
	}
}

func (s *Server) health(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{"status": "ok"},
	)
}

func (s *Server) ready(
	w http.ResponseWriter,
	request *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		request.Context(),
		2*time.Second,
	)
	defer cancel()
	if err := s.deps.Readiness.Check(ctx); err != nil {
		s.deps.Logger.Warn(
			"readiness dependency failed",
			"dependency",
			s.deps.Readiness.Name(),
			"error",
			err,
		)
		writeJSON(
			w,
			http.StatusServiceUnavailable,
			map[string]string{
				"status":     "not_ready",
				"dependency": s.deps.Readiness.Name(),
			},
		)
		return
	}
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status":     "ready",
			"dependency": s.deps.Readiness.Name(),
		},
	)
}

func (s *Server) apiStatus(
	w http.ResponseWriter,
	request *http.Request,
) {
	actor, ok := s.actor(w, request)
	if !ok {
		return
	}
	if err := s.deps.Policy.Require(
		actor,
		authz.PermissionViewDashboard,
	); err != nil {
		writeJSON(
			w,
			http.StatusForbidden,
			map[string]string{"error": err.Error()},
		)
		return
	}
	writeJSON(
		w,
		http.StatusOK,
		map[string]any{
			"status":  "candidate",
			"actor":   actor.ID,
			"roles":   actor.RoleNames(),
			"modules": len(s.deps.Modules.List()),
			"storage": s.deps.Readiness.Name(),
		},
	)
}

func (s *Server) apiChanges(
	w http.ResponseWriter,
	request *http.Request,
) {
	actor, ok := s.actor(w, request)
	if !ok {
		return
	}
	if err := s.deps.Policy.Require(
		actor,
		authz.PermissionViewDashboard,
	); err != nil {
		writeJSON(
			w,
			http.StatusForbidden,
			map[string]string{"error": err.Error()},
		)
		return
	}
	changes, err := s.deps.Changes.List(request.Context())
	if err != nil {
		s.dependencyFailure(w, "list API changes", err)
		return
	}
	writeJSON(w, http.StatusOK, changes)
}

func (s *Server) apiApprove(
	w http.ResponseWriter,
	request *http.Request,
) {
	actor, ok := s.actor(w, request)
	if !ok {
		return
	}
	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(
		http.MaxBytesReader(w, request.Body, 16<<10),
	).Decode(&body); err != nil {
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "valid JSON body is required",
			},
		)
		return
	}
	changeRequest, err := s.deps.Changes.Approve(
		request.Context(),
		request.PathValue("id"),
		actor,
		strings.TrimSpace(body.Reason),
	)
	if err != nil {
		s.writeChangeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, changeRequest)
}

func (s *Server) writeChangeError(
	w http.ResponseWriter,
	err error,
) {
	switch {
	case errors.Is(err, change.ErrInvalid):
		writeJSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
	case errors.Is(err, change.ErrForbidden):
		writeJSON(
			w,
			http.StatusForbidden,
			map[string]string{"error": err.Error()},
		)
	case errors.Is(err, change.ErrNotFound):
		writeJSON(
			w,
			http.StatusNotFound,
			map[string]string{"error": err.Error()},
		)
	case errors.Is(err, change.ErrConflict):
		writeJSON(
			w,
			http.StatusConflict,
			map[string]string{"error": err.Error()},
		)
	default:
		s.deps.Logger.Error("change operation failed", "error", err)
		writeJSON(
			w,
			http.StatusServiceUnavailable,
			map[string]string{
				"error": "change service unavailable",
			},
		)
	}
}

func (s *Server) dependencyFailure(
	w http.ResponseWriter,
	operation string,
	err error,
) {
	s.deps.Logger.Error(
		"dependency operation failed",
		"operation",
		operation,
		"error",
		err,
	)
	http.Error(
		w,
		"service dependency unavailable",
		http.StatusServiceUnavailable,
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

type pageData struct {
	Title           string
	Actor           authz.Actor
	Changes         []change.Request
	Modules         []modules.Module
	Now             time.Time
	DevelopmentMode bool
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(
				"Content-Security-Policy",
				"default-src 'self'; style-src 'self'; "+
					"script-src 'self'; img-src 'self' data:; "+
					"frame-ancestors 'none'; base-uri 'none'; "+
					"form-action 'self'",
			)
			w.Header().Set(
				"X-Content-Type-Options",
				"nosniff",
			)
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set(
				"Permissions-Policy",
				"camera=(), microphone=(), geolocation=()",
			)
			next.ServeHTTP(w, r)
		},
	)
}

func requestLog(
	logger *slog.Logger,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			next.ServeHTTP(w, r)
			logger.Info(
				"http request",
				"method",
				r.Method,
				"path",
				r.URL.Path,
				"duration_ms",
				time.Since(started).Milliseconds(),
			)
		},
	)
}
