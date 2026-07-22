package httpui

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/atlas/internal/authz"
	"github.com/Iron-Signal-Systems/atlas/internal/change"
	"github.com/Iron-Signal-Systems/atlas/internal/health"
	"github.com/Iron-Signal-Systems/atlas/internal/modules"
)

func testServer(
	t *testing.T,
	ready health.Checker,
	mode authentication.Mode,
) *Server {
	t.Helper()
	policy := authz.DefaultPolicy()
	authenticationMiddleware, err := authentication.New(
		authentication.Options{Mode: mode},
	)
	if err != nil {
		t.Fatal(err)
	}
	server, err := New(Dependencies{
		Logger: slog.New(
			slog.NewTextHandler(io.Discard, nil),
		),
		Policy:         policy,
		Changes:        change.NewMemoryService(policy),
		Modules:        modules.DefaultRegistry(),
		Readiness:      ready,
		Authentication: authenticationMiddleware,
	})
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestDashboardAndSecurityHeaders(t *testing.T) {
	server := testServer(
		t,
		health.Static{DependencyName: "memory"},
		authentication.ModeDevelopment,
	)
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	if response.Header().Get("Content-Security-Policy") == "" {
		t.Fatal("expected content security policy")
	}
}

type failedReadiness struct{}

func (failedReadiness) Name() string { return "postgresql" }

func (failedReadiness) Check(context.Context) error {
	return errors.New("offline")
}

func TestReadinessFailsClosedWithoutAuthentication(t *testing.T) {
	server := testServer(
		t,
		failedReadiness{},
		authentication.ModeProduction,
	)
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
}

func TestProductionModeRejectsDevelopmentHeaders(t *testing.T) {
	server := testServer(
		t,
		health.Static{DependencyName: "memory"},
		authentication.ModeProduction,
	)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/status",
		nil,
	)
	request.Header.Set(
		authentication.DevelopmentActorHeader,
		"attacker",
	)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
}

func TestProductionModeRequiresAuthentication(t *testing.T) {
	server := testServer(
		t,
		health.Static{DependencyName: "memory"},
		authentication.ModeProduction,
	)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/status",
		nil,
	)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
}

func TestQueryCannotSelectDevelopmentActor(t *testing.T) {
	server := testServer(
		t,
		health.Static{DependencyName: "memory"},
		authentication.ModeDevelopment,
	)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/status?actor=platform-admin",
		nil,
	)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	var body struct {
		Actor string `json:"actor"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Actor != "network-tech-01" {
		t.Fatalf(
			"query selected actor %q instead of server-side default",
			body.Actor,
		)
	}
}

func TestDevelopmentHeaderIdentityIsInjectedBeforeHandler(t *testing.T) {
	server := testServer(
		t,
		health.Static{DependencyName: "memory"},
		authentication.ModeDevelopment,
	)
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/status",
		nil,
	)
	request.Header.Set(
		authentication.DevelopmentActorHeader,
		"auditor-01",
	)
	request.Header.Set(
		authentication.DevelopmentRolesHeader,
		string(authz.RoleAuditor),
	)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
	var body struct {
		Actor string   `json:"actor"`
		Roles []string `json:"roles"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Actor != "auditor-01" {
		t.Fatalf("unexpected actor %q", body.Actor)
	}
	if len(body.Roles) != 1 || body.Roles[0] != string(authz.RoleAuditor) {
		t.Fatalf("unexpected roles %#v", body.Roles)
	}
}
