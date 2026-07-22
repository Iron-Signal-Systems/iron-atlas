package session

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/atlas/internal/authz"
)

type fakeStore struct {
	mu          sync.Mutex
	created     []CreateRequest
	findDigest  [sha256.Size]byte
	findCalls   int
	createError error
	findError   error
	record      Record
}

func (f *fakeStore) Create(
	_ context.Context,
	request CreateRequest,
) (Record, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.createError != nil {
		return Record{}, f.createError
	}
	f.created = append(f.created, request)
	return Record{
		Principal:             request.Principal,
		ActorID:               request.ActorID,
		CreatedAt:             request.CreatedAt,
		LastActivityAt:        request.CreatedAt,
		IdleExpiresAt:         request.IdleExpiresAt,
		AbsoluteExpiresAt:     request.AbsoluteExpiresAt,
		SecurityPolicyVersion: request.SecurityPolicyVersion,
	}, nil
}

func (f *fakeStore) Find(
	_ context.Context,
	digest [sha256.Size]byte,
) (Record, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.findCalls++
	f.findDigest = digest
	if f.findError != nil {
		return Record{}, f.findError
	}
	return f.record, nil
}

func testPrincipal() authentication.Principal {
	return authentication.Principal{
		ProviderID:      "oidc-test",
		Subject:         "subject-123",
		AuthenticatedAt: time.Date(2026, 7, 19, 20, 0, 0, 0, time.UTC),
		Assurance: authentication.Assurance{
			Context:          "urn:iron-atlas:assurance:provider-mfa",
			Methods:          []string{"pwd", "otp"},
			MFAAuthenticated: true,
			MFAAuthenticatedAt: time.Date(
				2026,
				7,
				19,
				20,
				0,
				0,
				0,
				time.UTC,
			),
			SecurityPolicyVersion: "phase-1-step-3-session-v1",
		},
	}
}

func newTestService(t *testing.T, store *fakeStore, random []byte) *Service {
	t.Helper()
	service, err := New(Config{
		Store: store,
		Resolver: authentication.ActorResolverFunc(func(
			context.Context,
			authentication.Principal,
		) (authz.Actor, error) {
			return authz.Actor{
				ID:    "actor-123",
				Roles: []authz.Role{authz.RoleNetworkTech},
			}, nil
		}),
		Random:                bytes.NewReader(random),
		Now:                   func() time.Time { return time.Date(2026, 7, 19, 20, 1, 0, 0, time.UTC) },
		IdleLifetime:          20 * time.Minute,
		AbsoluteLifetime:      4 * time.Hour,
		SuccessLocation:       "/",
		SecurityPolicyVersion: "phase-1-step-3-session-v1",
		RequireMFA:            true,
	})
	if err != nil {
		t.Fatal(err)
	}
	return service
}

func TestVerifiedPrincipalCreatesDigestOnlySecureSession(t *testing.T) {
	random := bytes.Repeat([]byte{0x42}, identifierBytes)
	store := &fakeStore{}
	service := newTestService(t, store, random)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)
	service.ServeVerifiedPrincipal(response, request, testPrincipal())

	if response.Code != http.StatusSeeOther {
		t.Fatalf("status = %d body = %q", response.Code, response.Body.String())
	}
	if response.Header().Get("Location") != "/" {
		t.Fatalf("Location = %q", response.Header().Get("Location"))
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("session cookie count = %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != CookieName || cookie.Path != "/" ||
		!cookie.Secure || !cookie.HttpOnly ||
		cookie.SameSite != http.SameSiteLaxMode || cookie.Domain != "" {
		t.Fatal("session cookie attributes are invalid")
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	if len(store.created) != 1 {
		t.Fatalf("created sessions = %d", len(store.created))
	}
	created := store.created[0]
	if created.ActorID != "actor-123" {
		t.Fatalf("actor = %q", created.ActorID)
	}
	if created.SecurityPolicyVersion != "phase-1-step-3-session-v1" {
		t.Fatalf("policy = %q", created.SecurityPolicyVersion)
	}
	decoded, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		t.Fatal(err)
	}
	expectedDigest := sha256.Sum256(decoded)
	if created.IdentifierDigest != expectedDigest {
		t.Fatal("stored digest does not match cookie identifier")
	}
	if bytes.Contains(created.IdentifierDigest[:], []byte(cookie.Value)) {
		t.Fatal("raw session identifier reached persistent request")
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q", response.Header().Get("Cache-Control"))
	}
}

func TestVerifiedPrincipalRejectsAlreadySessionBoundPrincipal(t *testing.T) {
	store := &fakeStore{}
	service := newTestService(t, store, bytes.Repeat([]byte{0x42}, identifierBytes))
	principal := testPrincipal()
	principal.BoundActorID = "actor-123"
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)

	service.ServeVerifiedPrincipal(response, request, principal)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", response.Code)
	}
	store.mu.Lock()
	defer store.mu.Unlock()
	if len(store.created) != 0 {
		t.Fatalf("created sessions = %d", len(store.created))
	}
}

func TestVerifiedPrincipalFailsClosedForRandomnessAndStoreFailure(t *testing.T) {
	for name, testCase := range map[string]struct {
		random []byte
		store  *fakeStore
	}{
		"randomness unavailable": {
			random: nil,
			store:  &fakeStore{},
		},
		"identifier conflict": {
			random: bytes.Repeat([]byte{0x42}, identifierBytes),
			store:  &fakeStore{createError: ErrSessionConflict},
		},
	} {
		t.Run(name, func(t *testing.T) {
			service := newTestService(t, testCase.store, testCase.random)
			response := httptest.NewRecorder()
			request := httptest.NewRequest(
				http.MethodGet,
				"https://atlas.example/auth/callback",
				nil,
			)
			service.ServeVerifiedPrincipal(response, request, testPrincipal())
			if response.Code != http.StatusServiceUnavailable {
				t.Fatalf("status = %d", response.Code)
			}
			if len(response.Result().Cookies()) != 0 {
				t.Fatal("failure response must not issue a session cookie")
			}
		})
	}
}

func TestAuthenticateReturnsServerBoundPrincipal(t *testing.T) {
	raw := bytes.Repeat([]byte{0x24}, identifierBytes)
	identifier := base64.RawURLEncoding.EncodeToString(raw)
	now := time.Date(2026, 7, 19, 20, 1, 0, 0, time.UTC)
	store := &fakeStore{record: Record{
		Principal:             testPrincipal(),
		ActorID:               "actor-123",
		CreatedAt:             now.Add(-time.Minute),
		LastActivityAt:        now.Add(-time.Minute),
		IdleExpiresAt:         now.Add(20 * time.Minute),
		AbsoluteExpiresAt:     now.Add(4 * time.Hour),
		SecurityPolicyVersion: "phase-1-step-3-session-v1",
	}}
	service := newTestService(t, store, bytes.Repeat([]byte{0x42}, identifierBytes))
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/changes", nil)
	request.AddCookie(&http.Cookie{Name: CookieName, Value: identifier})

	principal, err := service.Authenticate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if principal.BoundActorID != "actor-123" {
		t.Fatalf("bound actor = %q", principal.BoundActorID)
	}
	if principal.Assurance.SecurityPolicyVersion != "phase-1-step-3-session-v1" {
		t.Fatalf("policy = %q", principal.Assurance.SecurityPolicyVersion)
	}
	expectedDigest := sha256.Sum256(raw)
	if store.findDigest != expectedDigest {
		t.Fatal("lookup digest does not match cookie identifier")
	}
}

func TestAuthenticateRejectsDuplicateMalformedUnknownAndExpiredCookies(t *testing.T) {
	validRaw := bytes.Repeat([]byte{0x11}, identifierBytes)
	valid := base64.RawURLEncoding.EncodeToString(validRaw)
	for name, configure := range map[string]func(*http.Request, *fakeStore){
		"missing": func(_ *http.Request, _ *fakeStore) {},
		"duplicate": func(request *http.Request, _ *fakeStore) {
			request.AddCookie(&http.Cookie{Name: CookieName, Value: valid})
			request.AddCookie(&http.Cookie{Name: CookieName, Value: valid})
		},
		"malformed": func(request *http.Request, _ *fakeStore) {
			request.AddCookie(&http.Cookie{Name: CookieName, Value: "not-a-session"})
		},
		"expired": func(request *http.Request, store *fakeStore) {
			request.AddCookie(&http.Cookie{Name: CookieName, Value: valid})
			now := time.Date(2026, 7, 19, 20, 1, 0, 0, time.UTC)
			store.record = Record{
				Principal:             testPrincipal(),
				ActorID:               "actor-123",
				CreatedAt:             now.Add(-2 * time.Hour),
				LastActivityAt:        now.Add(-time.Hour),
				IdleExpiresAt:         now.Add(-time.Minute),
				AbsoluteExpiresAt:     now.Add(time.Hour),
				SecurityPolicyVersion: "phase-1-step-3-session-v1",
			}
		},
		"unknown": func(request *http.Request, store *fakeStore) {
			request.AddCookie(&http.Cookie{Name: CookieName, Value: valid})
			store.findError = ErrSessionNotFound
		},
	} {
		t.Run(name, func(t *testing.T) {
			store := &fakeStore{}
			service := newTestService(t, store, bytes.Repeat([]byte{0x42}, identifierBytes))
			request := httptest.NewRequest(http.MethodGet, "https://atlas.example/changes", nil)
			configure(request, store)
			_, err := service.Authenticate(context.Background(), request)
			if !errors.Is(err, authentication.ErrAuthenticationRequired) {
				t.Fatalf("error = %v", err)
			}
			if (name == "missing" || name == "duplicate" || name == "malformed") && store.findCalls != 0 {
				t.Fatalf("store lookups = %d", store.findCalls)
			}
		})
	}
}

func TestAuthenticateClassifiesStoreOutage(t *testing.T) {
	raw := bytes.Repeat([]byte{0x33}, identifierBytes)
	store := &fakeStore{findError: ErrSessionUnavailable}
	service := newTestService(t, store, bytes.Repeat([]byte{0x42}, identifierBytes))
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/changes", nil)
	request.AddCookie(&http.Cookie{
		Name:  CookieName,
		Value: base64.RawURLEncoding.EncodeToString(raw),
	})

	_, err := service.Authenticate(context.Background(), request)
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("error = %v", err)
	}
}

func TestServiceConcurrentAuthentication(t *testing.T) {
	raw := bytes.Repeat([]byte{0x55}, identifierBytes)
	identifier := base64.RawURLEncoding.EncodeToString(raw)
	now := time.Date(2026, 7, 19, 20, 1, 0, 0, time.UTC)
	store := &fakeStore{record: Record{
		Principal:             testPrincipal(),
		ActorID:               "actor-123",
		CreatedAt:             now.Add(-time.Minute),
		LastActivityAt:        now.Add(-time.Minute),
		IdleExpiresAt:         now.Add(20 * time.Minute),
		AbsoluteExpiresAt:     now.Add(4 * time.Hour),
		SecurityPolicyVersion: "phase-1-step-3-session-v1",
	}}
	service := newTestService(t, store, bytes.Repeat([]byte{0x42}, identifierBytes))

	const operations = 100
	var wait sync.WaitGroup
	errorsFound := make(chan error, operations)
	for range operations {
		wait.Add(1)
		go func() {
			defer wait.Done()
			request := httptest.NewRequest(http.MethodGet, "https://atlas.example/changes", nil)
			request.AddCookie(&http.Cookie{Name: CookieName, Value: identifier})
			principal, err := service.Authenticate(context.Background(), request)
			if err != nil {
				errorsFound <- err
				return
			}
			if principal.BoundActorID != "actor-123" {
				errorsFound <- errors.New("session actor binding changed")
			}
		}()
	}
	wait.Wait()
	close(errorsFound)
	for err := range errorsFound {
		t.Fatal(err)
	}
}

func TestNewRejectsUnsafeConfiguration(t *testing.T) {
	store := &fakeStore{}
	resolver := authentication.ActorResolverFunc(func(
		context.Context,
		authentication.Principal,
	) (authz.Actor, error) {
		return authz.Actor{ID: "actor-123"}, nil
	})
	for name, config := range map[string]Config{
		"missing store":    {Resolver: resolver, SecurityPolicyVersion: "v1"},
		"missing resolver": {Store: store, SecurityPolicyVersion: "v1"},
		"external redirect": {
			Store: store, Resolver: resolver, SecurityPolicyVersion: "v1",
			SuccessLocation: "https://attacker.example/",
		},
		"protocol-relative redirect": {
			Store: store, Resolver: resolver, SecurityPolicyVersion: "v1",
			SuccessLocation: "//attacker.example/",
		},
		"backslash redirect": {
			Store: store, Resolver: resolver, SecurityPolicyVersion: "v1",
			SuccessLocation: "/\\attacker.example/",
		},
		"encoded separator redirect": {
			Store: store, Resolver: resolver, SecurityPolicyVersion: "v1",
			SuccessLocation: "/%2f%2fattacker.example/",
		},
		"query-bearing redirect": {
			Store: store, Resolver: resolver, SecurityPolicyVersion: "v1",
			SuccessLocation: "/?next=attacker",
		},
		"idle exceeds absolute": {
			Store: store, Resolver: resolver, SecurityPolicyVersion: "v1",
			IdleLifetime: 2 * time.Hour, AbsoluteLifetime: time.Hour,
		},
		"missing policy": {Store: store, Resolver: resolver},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := New(config); err == nil {
				t.Fatal("expected configuration rejection")
			}
		})
	}
}

func TestVerifiedPrincipalRequiresSatisfiedMFAPolicy(t *testing.T) {
	for name, mutate := range map[string]func(*authentication.Principal){
		"missing MFA": func(principal *authentication.Principal) {
			principal.Assurance.MFAAuthenticated = false
			principal.Assurance.MFAAuthenticatedAt = time.Time{}
		},
		"wrong policy version": func(principal *authentication.Principal) {
			principal.Assurance.SecurityPolicyVersion = "other-policy"
		},
		"future MFA time": func(principal *authentication.Principal) {
			principal.Assurance.MFAAuthenticatedAt = time.Date(
				2026, 7, 19, 20, 5, 0, 0, time.UTC,
			)
		},
	} {
		t.Run(name, func(t *testing.T) {
			store := &fakeStore{}
			service := newTestService(t, store, bytes.Repeat([]byte{0x42}, identifierBytes))
			principal := testPrincipal()
			mutate(&principal)
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)
			service.ServeVerifiedPrincipal(response, request, principal)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d", response.Code)
			}
			if len(response.Result().Cookies()) != 0 {
				t.Fatal("unsatisfied assurance emitted a session cookie")
			}
			store.mu.Lock()
			defer store.mu.Unlock()
			if len(store.created) != 0 {
				t.Fatalf("created sessions = %d", len(store.created))
			}
		})
	}
}
