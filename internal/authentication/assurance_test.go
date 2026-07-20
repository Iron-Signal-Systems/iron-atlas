package authentication

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authz"
)

func TestAssuranceValidation(t *testing.T) {
	now := time.Date(2026, 7, 19, 20, 0, 0, 0, time.UTC)
	valid := Assurance{
		Context:               "urn:iron-atlas:assurance:provider-mfa",
		Methods:               []string{"pwd", "otp"},
		MFAAuthenticated:      true,
		MFAAuthenticatedAt:    now,
		SecurityPolicyVersion: "phase-1-step-3-session-v1",
	}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}

	for name, assurance := range map[string]Assurance{
		"duplicate methods": {
			Methods: []string{"otp", "otp"},
		},
		"unnormalized method": {
			Methods: []string{" otp"},
		},
		"MFA missing time": {
			MFAAuthenticated: true,
		},
		"non-MFA with time": {
			MFAAuthenticatedAt: now,
		},
		"oversized context": {
			Context: strings.Repeat("x", 257),
		},
	} {
		t.Run(name, func(t *testing.T) {
			if err := assurance.Validate(); err == nil {
				t.Fatal("expected assurance rejection")
			}
		})
	}
}

func TestResolvedIdentityRejectsSessionActorRemapping(t *testing.T) {
	identity := ResolvedIdentity{
		Principal: Principal{
			ProviderID:      "oidc-test",
			Subject:         "subject-123",
			AuthenticatedAt: time.Now().UTC(),
			BoundActorID:    "actor-original",
		},
		Actor: authz.Actor{
			ID:    "actor-remapped",
			Roles: []authz.Role{authz.RoleViewer},
		},
	}
	if err := identity.Validate(); err == nil {
		t.Fatal("expected session actor remapping rejection")
	}
}

func TestResolvedIdentityClonesAssuranceMethods(t *testing.T) {
	identity := ResolvedIdentity{
		Principal: Principal{
			ProviderID:      "oidc-test",
			Subject:         "subject-123",
			AuthenticatedAt: time.Now().UTC(),
			Assurance: Assurance{
				Methods: []string{"pwd", "otp"},
			},
		},
		Actor: authz.Actor{ID: "actor-123"},
	}
	ctx := withIdentity(context.Background(), identity)
	first, ok := ResolvedIdentityFromContext(ctx)
	if !ok {
		t.Fatal("identity was not stored")
	}
	first.Principal.Assurance.Methods[0] = "tampered"
	second, ok := ResolvedIdentityFromContext(ctx)
	if !ok {
		t.Fatal("identity was not stored")
	}
	if second.Principal.Assurance.Methods[0] != "pwd" {
		t.Fatalf("stored assurance mutated: %#v", second.Principal.Assurance.Methods)
	}
}
