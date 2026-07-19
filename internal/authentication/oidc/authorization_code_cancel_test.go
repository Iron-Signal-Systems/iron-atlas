package oidc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
)

func TestAuthorizationCodeFlowCancelConsumesStateExactlyOnce(t *testing.T) {
	now := time.Date(2026, 7, 19, 19, 0, 0, 0, time.UTC)
	state := base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{0x11}, 32))
	nonce := base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{0x22}, 32))
	verifier := base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{0x33}, 32))
	redirectURL := "https://atlas.example.test/auth/callback"

	store, err := NewMemoryPreauthenticationStore(4)
	if err != nil {
		t.Fatal(err)
	}
	transaction := PreauthenticationTransaction{
		StateDigest:  sha256.Sum256([]byte(state)),
		Nonce:        nonce,
		PKCEVerifier: verifier,
		RedirectURL:  redirectURL,
		CreatedAt:    now,
		ExpiresAt:    now.Add(5 * time.Minute),
	}
	if err := store.Create(context.Background(), transaction); err != nil {
		t.Fatal(err)
	}

	flow := &AuthorizationCodeFlow{
		verifier: &Verifier{issuerURL: "https://identity.example.test/oidc"},
		store:    store,
		oauthConfig: oauth2.Config{
			RedirectURL: redirectURL,
		},
		now: func() time.Time { return now },
	}

	if got := flow.IssuerURL(); got != "https://identity.example.test/oidc" {
		t.Fatalf("issuer URL = %q", got)
	}
	if err := flow.Cancel(context.Background(), state); err != nil {
		t.Fatalf("first cancel failed: %v", err)
	}
	if err := flow.Cancel(context.Background(), state); !errors.Is(
		err,
		authentication.ErrAuthenticationInvalid,
	) {
		t.Fatalf("second cancel error = %v, want authentication invalid", err)
	}
}
