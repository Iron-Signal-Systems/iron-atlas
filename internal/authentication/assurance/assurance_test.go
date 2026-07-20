package assurance

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

const testPolicyVersion = "phase-1-step-3-assurance-v1"

func testPolicy(t *testing.T) *Policy {
	t.Helper()
	policy, err := NewPolicy(PolicyConfig{
		Version:                  testPolicyVersion,
		RequireMFA:               true,
		MaximumAuthenticationAge: 10 * time.Minute,
		AcceptedMFAContexts:      []string{"urn:example:acr:mfa"},
		AcceptedMFAMethodSets: []MethodSet{
			{"pwd", "otp"},
			{"mfa"},
		},
		PhishingResistantRoles: []authz.Role{
			authz.RolePlatformAdmin,
			authz.RoleChangeApprover,
		},
		AcceptedPhishingResistantContexts: []string{"urn:example:acr:phishing-resistant"},
		AcceptedPhishingResistantMethodSets: []MethodSet{
			{"hwk", "user"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	return policy
}

func testPrincipal(now time.Time, contextValue string, methods ...string) authentication.Principal {
	return authentication.Principal{
		ProviderID:      "oidc-test",
		Subject:         "subject-123",
		AuthenticatedAt: now.Add(-time.Minute),
		Assurance: authentication.Assurance{
			Context: contextValue,
			Methods: append([]string(nil), methods...),
		},
	}
}

func TestPolicyAcceptsExplicitProviderMFA(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	actor := authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RoleNetworkTech}}
	for name, principal := range map[string]authentication.Principal{
		"accepted context":    testPrincipal(now, "urn:example:acr:mfa", "pwd"),
		"accepted method set": testPrincipal(now, "", "otp", "pwd", "provider-extra"),
	} {
		t.Run(name, func(t *testing.T) {
			decision := testPolicy(t).Evaluate(principal, actor, now)
			if !decision.Satisfied() {
				t.Fatalf("outcome = %q reason = %q", decision.Outcome, decision.ReasonCode)
			}
			if !decision.Assurance.MFAAuthenticated {
				t.Fatal("MFA assurance was not established")
			}
			if !decision.Assurance.MFAAuthenticatedAt.Equal(principal.AuthenticatedAt) {
				t.Fatalf("MFA time = %s", decision.Assurance.MFAAuthenticatedAt)
			}
			if decision.Assurance.SecurityPolicyVersion != testPolicyVersion {
				t.Fatalf("policy version = %q", decision.Assurance.SecurityPolicyVersion)
			}
		})
	}
}

func TestPolicyDoesNotInferMFAFromUnknownClaims(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	actor := authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RoleViewer}}
	for name, principal := range map[string]authentication.Principal{
		"missing assurance":  testPrincipal(now, ""),
		"unknown context":    testPrincipal(now, "urn:unknown:acr", "pwd"),
		"partial method set": testPrincipal(now, "", "otp"),
		"provider asserted boolean ignored": func() authentication.Principal {
			principal := testPrincipal(now, "", "pwd")
			principal.Assurance.MFAAuthenticated = true
			principal.Assurance.MFAAuthenticatedAt = principal.AuthenticatedAt
			return principal
		}(),
	} {
		t.Run(name, func(t *testing.T) {
			decision := testPolicy(t).Evaluate(principal, actor, now)
			if decision.Outcome != OutcomeAdditionalAuthenticationRequired || decision.ReasonCode != ReasonMFARequired {
				t.Fatalf("decision = %#v", decision)
			}
		})
	}
}

func TestPolicyRequiresPhishingResistanceForHighImpactRoles(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	actor := authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RolePlatformAdmin}}
	ordinaryMFA := testPrincipal(now, "urn:example:acr:mfa", "pwd", "otp")
	decision := testPolicy(t).Evaluate(ordinaryMFA, actor, now)
	if decision.Outcome != OutcomePhishingResistantRequired {
		t.Fatalf("ordinary MFA outcome = %q", decision.Outcome)
	}

	phishingResistant := testPrincipal(now, "urn:example:acr:phishing-resistant", "hwk", "user")
	decision = testPolicy(t).Evaluate(phishingResistant, actor, now)
	if !decision.Satisfied() {
		t.Fatalf("phishing-resistant outcome = %q", decision.Outcome)
	}
}

func TestPolicyRequiresFreshAuthentication(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	principal := testPrincipal(now, "urn:example:acr:mfa", "pwd", "otp")
	principal.AuthenticatedAt = now.Add(-20 * time.Minute)
	decision := testPolicy(t).Evaluate(
		principal,
		authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RoleViewer}},
		now,
	)
	if decision.Outcome != OutcomeStepUpRequired || decision.ReasonCode != ReasonPrimaryAuthenticationStale {
		t.Fatalf("decision = %#v", decision)
	}
}

func TestNewPolicyRejectsAmbiguousOrUnsafeConfiguration(t *testing.T) {
	tests := map[string]PolicyConfig{
		"missing version": {},
		"MFA disabled":    {Version: testPolicyVersion},
		"duplicate context": {
			Version:             testPolicyVersion,
			AcceptedMFAContexts: []string{"a", "a"},
		},
		"duplicate method": {
			Version:               testPolicyVersion,
			AcceptedMFAMethodSets: []MethodSet{{"pwd", "pwd"}},
		},
		"duplicate role": {
			Version:                testPolicyVersion,
			PhishingResistantRoles: []authz.Role{authz.RolePlatformAdmin, authz.RolePlatformAdmin},
		},
		"unknown role": {
			Version:                testPolicyVersion,
			PhishingResistantRoles: []authz.Role{"unknown"},
		},
		"excessive age": {
			Version:                  testPolicyVersion,
			MaximumAuthenticationAge: 25 * time.Hour,
		},
	}
	for name, config := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := NewPolicy(config); err == nil {
				t.Fatal("unsafe policy configuration was accepted")
			}
		})
	}
}

type recordingHandler struct {
	mu         sync.Mutex
	principals []authentication.Principal
	status     int
}

func (h *recordingHandler) ServeVerifiedPrincipal(
	writer http.ResponseWriter,
	_ *http.Request,
	principal authentication.Principal,
) {
	h.mu.Lock()
	h.principals = append(h.principals, principal)
	h.mu.Unlock()
	status := h.status
	if status == 0 {
		status = http.StatusNoContent
	}
	writer.WriteHeader(status)
}

func (h *recordingHandler) count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.principals)
}

func TestServiceForwardsOnlySatisfiedAssurance(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	downstream := &recordingHandler{}
	service, err := NewService(ServiceConfig{
		Resolver: authentication.ActorResolverFunc(func(
			context.Context,
			authentication.Principal,
		) (authz.Actor, error) {
			return authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RoleNetworkTech}}, nil
		}),
		Policy: testPolicy(t),
		Next:   downstream,
		Now:    func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)
	service.ServeVerifiedPrincipal(
		response,
		request,
		testPrincipal(now, "urn:example:acr:mfa", "pwd", "otp"),
	)
	if response.Code != http.StatusNoContent || downstream.count() != 1 {
		t.Fatalf("status = %d downstream = %d", response.Code, downstream.count())
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("Cache-Control = %q", response.Header().Get("Cache-Control"))
	}
}

func TestServiceRejectsUnsatisfiedAssuranceWithoutSessionHandoff(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	downstream := &recordingHandler{}
	service, err := NewService(ServiceConfig{
		Resolver: authentication.ActorResolverFunc(func(
			context.Context,
			authentication.Principal,
		) (authz.Actor, error) {
			return authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RoleNetworkTech}}, nil
		}),
		Policy: testPolicy(t),
		Next:   downstream,
		Now:    func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)
	service.ServeVerifiedPrincipal(response, request, testPrincipal(now, "", "pwd"))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", response.Code)
	}
	if downstream.count() != 0 {
		t.Fatalf("downstream calls = %d", downstream.count())
	}
	if len(response.Result().Cookies()) != 0 {
		t.Fatal("unsatisfied assurance emitted a cookie")
	}
}

func TestServiceClassifiesResolverOutage(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	service, err := NewService(ServiceConfig{
		Resolver: authentication.ActorResolverFunc(func(
			context.Context,
			authentication.Principal,
		) (authz.Actor, error) {
			return authz.Actor{}, authentication.ErrAuthenticationUnavailable
		}),
		Policy: testPolicy(t),
		Next:   &recordingHandler{},
		Now:    func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)
	service.ServeVerifiedPrincipal(
		response,
		request,
		testPrincipal(now, "urn:example:acr:mfa", "pwd", "otp"),
	)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestServiceConcurrentEvaluation(t *testing.T) {
	now := time.Date(2026, 7, 20, 20, 0, 0, 0, time.UTC)
	downstream := &recordingHandler{}
	service, err := NewService(ServiceConfig{
		Resolver: authentication.ActorResolverFunc(func(
			context.Context,
			authentication.Principal,
		) (authz.Actor, error) {
			return authz.Actor{ID: "actor-123", Roles: []authz.Role{authz.RoleNetworkTech}}, nil
		}),
		Policy: testPolicy(t),
		Next:   downstream,
		Now:    func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}

	const attempts = 64
	var wait sync.WaitGroup
	errorsFound := make(chan error, attempts)
	for range attempts {
		wait.Add(1)
		go func() {
			defer wait.Done()
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "https://atlas.example/auth/callback", nil)
			service.ServeVerifiedPrincipal(
				response,
				request,
				testPrincipal(now, "urn:example:acr:mfa", "pwd", "otp"),
			)
			if response.Code != http.StatusNoContent {
				errorsFound <- errors.New("unexpected response status")
			}
		}()
	}
	wait.Wait()
	close(errorsFound)
	for err := range errorsFound {
		t.Fatal(err)
	}
	if downstream.count() != attempts {
		t.Fatalf("downstream calls = %d", downstream.count())
	}
}
