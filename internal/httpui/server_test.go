package httpui

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/change"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/modules"
)

func TestDashboardAndSecurityHeaders(t *testing.T) {
	policy := authz.DefaultPolicy()
	server, err := New(Dependencies{
		Logger:              slog.New(slog.NewTextHandler(io.Discard, nil)),
		Policy:              policy,
		Changes:             change.NewMemoryService(policy),
		Modules:             modules.DefaultRegistry(),
		DevelopmentIdentity: true,
	})
	if err != nil {
		t.Fatal(err)
	}
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
