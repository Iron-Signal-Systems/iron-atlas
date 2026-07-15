package httpui

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/change"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/modules"
)

type Dependencies struct {
	Logger              *slog.Logger
	Policy              *authz.Policy
	Changes             change.Service
	Modules             modules.Registry
	DevelopmentIdentity bool
}

type Server struct {
	deps      Dependencies
	templates *template.Template
	mux       *http.ServeMux
}

func New(deps Dependencies) (*Server, error) {
	if deps.Logger == nil || deps.Policy == nil || deps.Changes == nil {
		return nil, errors.New("logger, policy, and change service are required")
	}
	templates, err := template.ParseFS(webFiles, "templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}
	s := &Server{deps: deps, templates: templates, mux: http.NewServeMux()}
	s.routes()
	return s, nil
}

func (s *Server) routes() {
	staticFS, _ := fs.Sub(webFiles, "static")
	s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	s.mux.HandleFunc("GET /healthz", s.health)
	s.mux.HandleFunc("GET /readyz", s.ready)
	s.mux.HandleFunc("GET /api/v1/status", s.apiStatus)
	s.mux.HandleFunc("GET /api/v1/changes", s.apiChanges)
	s.mux.HandleFunc("POST /api/v1/changes/{id}/approve", s.apiApprove)
	s.mux.HandleFunc("GET /changes", s.changesPage)
	s.mux.HandleFunc("GET /modules", s.modulesPage)
	s.mux.HandleFunc("GET /", s.dashboard)
}

func (s *Server) Handler() http.Handler {
	return requestLog(s.deps.Logger, securityHeaders(s.mux))
}

func (s *Server) actor(r *http.Request) authz.Actor {
	if !s.deps.DevelopmentIdentity {
		return authz.Actor{}
	}
	id := strings.TrimSpace(r.Header.Get("X-Iron-Atlas-Actor"))
	if id == "" {
		id = "network-tech-01"
	}
	rawRoles := strings.TrimSpace(r.Header.Get("X-Iron-Atlas-Roles"))
	if rawRoles == "" {
		rawRoles = string(authz.RoleNetworkTech)
	}
	roles := make([]authz.Role, 0)
	for _, value := range strings.Split(rawRoles, ",") {
		value = strings.TrimSpace(value)
		if value != "" {
			roles = append(roles, authz.Role(value))
		}
	}
	return authz.Actor{ID: id, Roles: roles}
}

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	actor := s.actor(r)
	data := pageData{
		Title:               "Dashboard",
		Actor:               actor,
		Changes:             s.deps.Changes.List(),
		Modules:             s.deps.Modules.List(),
		Now:                 time.Now().UTC(),
		DevelopmentIdentity: s.deps.DevelopmentIdentity,
	}
	s.render(w, "dashboard.html", data)
}

func (s *Server) changesPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "changes.html", pageData{Title: "Change management", Actor: s.actor(r), Changes: s.deps.Changes.List(), DevelopmentIdentity: s.deps.DevelopmentIdentity})
}

func (s *Server) modulesPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, "modules.html", pageData{Title: "Modules", Actor: s.actor(r), Modules: s.deps.Modules.List(), DevelopmentIdentity: s.deps.DevelopmentIdentity})
}

func (s *Server) render(w http.ResponseWriter, name string, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.templates.ExecuteTemplate(w, name, data); err != nil {
		s.deps.Logger.Error("template rendering failed", "template", name, "error", err)
	}
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready", "storage": "memory-candidate"})
}

func (s *Server) apiStatus(w http.ResponseWriter, r *http.Request) {
	actor := s.actor(r)
	if err := s.deps.Policy.Require(actor, authz.PermissionViewDashboard); err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"status": "candidate", "actor": actor.ID, "roles": actor.RoleNames(), "modules": len(s.deps.Modules.List())})
}

func (s *Server) apiChanges(w http.ResponseWriter, r *http.Request) {
	actor := s.actor(r)
	if err := s.deps.Policy.Require(actor, authz.PermissionViewDashboard); err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, s.deps.Changes.List())
}

func (s *Server) apiApprove(w http.ResponseWriter, r *http.Request) {
	actor := s.actor(r)
	var body struct {
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 16<<10)).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "valid JSON body is required"})
		return
	}
	request, err := s.deps.Changes.Approve(r.PathValue("id"), actor, strings.TrimSpace(body.Reason))
	if err != nil {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, request)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

type pageData struct {
	Title               string
	Actor               authz.Actor
	Changes             []change.Request
	Modules             []modules.Module
	Now                 time.Time
	DevelopmentIdentity bool
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self'; script-src 'self'; img-src 'self' data:; frame-ancestors 'none'; base-uri 'none'; form-action 'self'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

func requestLog(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("http request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(started).Milliseconds())
	})
}
