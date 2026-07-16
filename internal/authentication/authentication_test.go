package authentication

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

func TestParseMode(t *testing.T) {
	for _, test := range []struct {
		raw  string
		want Mode
	}{
		{raw: "development", want: ModeDevelopment},
		{raw: " DEVELOPMENT ", want: ModeDevelopment},
		{raw: "production", want: ModeProduction},
	} {
		got, err := ParseMode(test.raw)
		if err != nil {
			t.Fatalf("ParseMode(%q): %v", test.raw, err)
		}
		if got != test.want {
			t.Fatalf("ParseMode(%q) = %q, want %q", test.raw, got, test.want)
		}
	}
	if _, err := ParseMode("automatic"); err == nil {
		t.Fatal("unsupported authentication mode must fail closed")
	}
}

func TestDevelopmentModeInjectsImmutableResolvedIdentity(t *testing.T) {
	now := time.Date(2026, 7, 16, 10, 0, 0, 0, time.UTC)
	middleware, err := New(Options{
		Mode: ModeDevelopment,
		Now:  func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}

	var first authz.Actor
	var second authz.Actor
	handler := middleware.Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var ok bool
			first, ok = ActorFromContext(r.Context())
			if !ok {
				t.Fatal("resolved actor is missing")
			}
			first.Roles[0] = authz.RolePlatformAdmin
			second, ok = ActorFromContext(r.Context())
			if !ok {
				t.Fatal("resolved actor disappeared")
			}
			principal, ok := PrincipalFromContext(r.Context())
			if !ok {
				t.Fatal("principal is missing")
			}
			if principal.ProviderID != "development" ||
				principal.Subject != "auditor-01" ||
				!principal.AuthenticatedAt.Equal(now) {
				t.Fatalf("unexpected principal: %#v", principal)
			}
			w.WriteHeader(http.StatusNoContent)
		},
	))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	request.Header.Set(DevelopmentActorHeader, "auditor-01")
	request.Header.Set(DevelopmentRolesHeader, "auditor")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.Code)
	}
	if first.ID != "auditor-01" {
		t.Fatalf("unexpected actor ID %q", first.ID)
	}
	if len(second.Roles) != 1 || second.Roles[0] != authz.RoleAuditor {
		t.Fatalf("context identity was mutable: %#v", second.Roles)
	}
}

func TestDevelopmentModeRejectsDuplicateOrUnknownRoleInput(t *testing.T) {
	middleware, err := New(Options{Mode: ModeDevelopment})
	if err != nil {
		t.Fatal(err)
	}
	nextCalled := false
	handler := middleware.Handler(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		},
	))

	for name, configure := range map[string]func(*http.Request){
		"duplicate actor header": func(request *http.Request) {
			request.Header.Add(DevelopmentActorHeader, "one")
			request.Header.Add(DevelopmentActorHeader, "two")
		},
		"unknown role": func(request *http.Request) {
			request.Header.Set(DevelopmentRolesHeader, "god_mode")
		},
		"duplicate role": func(request *http.Request) {
			request.Header.Set(DevelopmentRolesHeader, "viewer,viewer")
		},
	} {
		t.Run(name, func(t *testing.T) {
			nextCalled = false
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			configure(request)
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d", response.Code)
			}
			if nextCalled {
				t.Fatal("invalid identity reached the protected handler")
			}
		})
	}
}

func TestProductionModeRejectsDevelopmentHeadersBeforeAdapter(t *testing.T) {
	authenticatorCalled := false
	resolverCalled := false
	middleware, err := New(Options{
		Mode: ModeProduction,
		Authenticator: AuthenticatorFunc(
			func(context.Context, *http.Request) (Principal, error) {
				authenticatorCalled = true
				return Principal{}, nil
			},
		),
		ActorResolver: ActorResolverFunc(
			func(context.Context, Principal) (authz.Actor, error) {
				resolverCalled = true
				return authz.Actor{}, nil
			},
		),
	})
	if err != nil {
		t.Fatal(err)
	}
	handler := middleware.Handler(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			t.Fatal("protected handler must not run")
		},
	))
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(DevelopmentActorHeader, "attacker")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", response.Code)
	}
	if authenticatorCalled || resolverCalled {
		t.Fatal("development headers reached production authentication components")
	}
}

func TestProductionModeWithoutAdapterFailsClosed(t *testing.T) {
	middleware, err := New(Options{Mode: ModeProduction})
	if err != nil {
		t.Fatal(err)
	}
	nextCalled := false
	handler := middleware.Handler(http.HandlerFunc(
		func(http.ResponseWriter, *http.Request) {
			nextCalled = true
		},
	))
	response := httptest.NewRecorder()
	handler.ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/api/v1/status", nil),
	)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", response.Code)
	}
	if nextCalled {
		t.Fatal("unauthenticated production request reached handler")
	}
}

func TestProductionAdapterAndResolverSetServerSideIdentity(t *testing.T) {
	now := time.Date(2026, 7, 16, 10, 0, 0, 0, time.UTC)
	middleware, err := New(Options{
		Mode: ModeProduction,
		Authenticator: AuthenticatorFunc(
			func(context.Context, *http.Request) (Principal, error) {
				return Principal{
					ProviderID:      "oidc-test",
					Subject:         "subject-123",
					AuthenticatedAt: now,
				}, nil
			},
		),
		ActorResolver: ActorResolverFunc(
			func(_ context.Context, principal Principal) (authz.Actor, error) {
				if principal.Subject != "subject-123" {
					t.Fatalf("unexpected subject %q", principal.Subject)
				}
				return authz.Actor{
					ID:    "actor-123",
					Roles: []authz.Role{authz.RoleViewer},
				}, nil
			},
		),
	})
	if err != nil {
		t.Fatal(err)
	}

	handler := middleware.Handler(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			actor, ok := ActorFromContext(r.Context())
			if !ok || actor.ID != "actor-123" {
				t.Fatalf("unexpected resolved actor: %#v", actor)
			}
			w.WriteHeader(http.StatusNoContent)
		},
	))
	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/status?actor=attacker",
		nil,
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", response.Code)
	}
}

func TestHealthAndReadinessPathsRemainPublic(t *testing.T) {
	middleware, err := New(Options{Mode: ModeProduction})
	if err != nil {
		t.Fatal(err)
	}
	handler := middleware.Handler(http.HandlerFunc(
		func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNoContent)
		},
	))
	for _, path := range []string{"/healthz", "/readyz", "/static/site.css"} {
		response := httptest.NewRecorder()
		handler.ServeHTTP(
			response,
			httptest.NewRequest(http.MethodGet, path, nil),
		)
		if response.Code != http.StatusNoContent {
			t.Fatalf("%s: expected 204, got %d", path, response.Code)
		}
	}
}

func TestNestedIdentityMiddlewareFailsClosed(t *testing.T) {
	middleware, err := New(Options{Mode: ModeDevelopment})
	if err != nil {
		t.Fatal(err)
	}
	handler := middleware.Handler(
		middleware.Handler(http.HandlerFunc(
			func(http.ResponseWriter, *http.Request) {
				t.Fatal("nested identity middleware reached handler")
			},
		)),
	)
	response := httptest.NewRecorder()
	handler.ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/", nil),
	)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", response.Code)
	}
}
