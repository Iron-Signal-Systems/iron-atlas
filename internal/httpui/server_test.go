package httpui

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/change"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/health"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/modules"
)

func testServer(t *testing.T, ready health.Checker) *Server {
	t.Helper()
	policy := authz.DefaultPolicy()
	server, err := New(Dependencies{
		Logger:              slog.New(slog.NewTextHandler(io.Discard, nil)),
		Policy:              policy,
		Changes:             change.NewMemoryService(policy),
		Modules:             modules.DefaultRegistry(),
		Readiness:           ready,
		DevelopmentIdentity: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return server
}

func TestDashboardAndSecurityHeaders(t *testing.T) {
	server := testServer(t, health.Static{DependencyName: "memory"})
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

func (failedReadiness) Name() string                { return "postgresql" }
func (failedReadiness) Check(context.Context) error { return errors.New("offline") }

func TestReadinessFailsClosed(t *testing.T) {
	server := testServer(t, failedReadiness{})
	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", response.Code)
	}
}
