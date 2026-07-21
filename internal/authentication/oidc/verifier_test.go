package oidc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
)

type providerEmulator struct {
	t      *testing.T
	server *httptest.Server
	client *http.Client
	now    time.Time

	mu        sync.RWMutex
	key       *rsa.PrivateKey
	keyID     string
	available bool
}

func newProviderEmulator(t *testing.T) *providerEmulator {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	emulator := &providerEmulator{
		t:         t,
		now:       time.Date(2026, 7, 16, 13, 0, 0, 0, time.UTC),
		key:       key,
		keyID:     "key-1",
		available: true,
	}
	emulator.server = httptest.NewTLSServer(http.HandlerFunc(emulator.handle))
	emulator.client = emulator.server.Client()
	emulator.client.Timeout = 2 * time.Second
	t.Cleanup(emulator.server.Close)
	return emulator
}

func (e *providerEmulator) handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	e.mu.RLock()
	available := e.available
	key := e.key
	keyID := e.keyID
	e.mu.RUnlock()

	if !available {
		http.Error(w, "provider unavailable", http.StatusServiceUnavailable)
		return
	}
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		writeTestJSON(w, map[string]any{
			"issuer":                                e.server.URL,
			"authorization_endpoint":                e.server.URL + "/authorize",
			"token_endpoint":                        e.server.URL + "/token",
			"jwks_uri":                              e.server.URL + "/jwks",
			"response_types_supported":              []string{"code"},
			"subject_types_supported":               []string{"public"},
			"id_token_signing_alg_values_supported": []string{"RS256"},
			"token_endpoint_auth_methods_supported": []string{"client_secret_basic"},
		})
	case "/jwks":
		writeTestJSON(w, jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{{
				Key:       &key.PublicKey,
				KeyID:     keyID,
				Algorithm: string(jose.RS256),
				Use:       "sig",
			}},
		})
	default:
		http.NotFound(w, r)
	}
}

func writeTestJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(value)
}

func accessTokenHash(accessToken string) string {
	digest := sha256.Sum256([]byte(accessToken))
	return base64.RawURLEncoding.EncodeToString(
		digest[:len(digest)/2],
	)
}

func (e *providerEmulator) verifier(t *testing.T) *Verifier {
	t.Helper()
	verifier, err := New(context.Background(), Config{
		ProviderID:               "oidc-test",
		IssuerURL:                e.server.URL,
		ClientID:                 "iron-atlas-test",
		AllowedSigningAlgorithms: []string{"RS256"},
		HTTPClient:               e.client,
		Now:                      func() time.Time { return e.now },
		MaxClockSkew:             time.Minute,
		MaxTokenAge:              15 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	return verifier
}

func (e *providerEmulator) token(
	t *testing.T,
	overrides map[string]any,
) string {
	t.Helper()
	e.mu.RLock()
	key := e.key
	keyID := e.keyID
	e.mu.RUnlock()
	return signedToken(t, key, keyID, e.defaultClaims(overrides), jose.RS256)
}

func (e *providerEmulator) defaultClaims(
	overrides map[string]any,
) map[string]any {
	claims := map[string]any{
		"iss":   e.server.URL,
		"sub":   "subject-123",
		"aud":   "iron-atlas-test",
		"exp":   e.now.Add(5 * time.Minute).Unix(),
		"iat":   e.now.Add(-time.Minute).Unix(),
		"nonce": "0123456789abcdef0123456789abcdef",
	}
	for key, value := range overrides {
		if value == nil {
			delete(claims, key)
			continue
		}
		claims[key] = value
	}
	return claims
}

func signedToken(
	t *testing.T,
	key *rsa.PrivateKey,
	keyID string,
	claims map[string]any,
	algorithm jose.SignatureAlgorithm,
) string {
	t.Helper()
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	return signedPayload(t, key, keyID, payload, algorithm)
}

func signedPayload(
	t *testing.T,
	key *rsa.PrivateKey,
	keyID string,
	payload []byte,
	algorithm jose.SignatureAlgorithm,
) string {
	t.Helper()
	options := (&jose.SignerOptions{}).
		WithType("JWT").
		WithHeader("kid", keyID)
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: algorithm, Key: key},
		options,
	)
	if err != nil {
		t.Fatal(err)
	}
	object, err := signer.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := object.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestVerifierAcceptsValidToken(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	principal, err := verifier.Verify(
		context.Background(),
		provider.token(t, nil),
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if principal.ProviderID != "oidc-test" ||
		principal.Subject != "subject-123" ||
		!principal.AuthenticatedAt.Equal(provider.now.Add(-time.Minute)) {
		t.Fatalf("unexpected principal: %#v", principal)
	}
}

func TestVerifierEnforcesAccessTokenHash(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	accessToken := "opaque-access-token-value"
	raw := provider.token(t, map[string]any{
		"at_hash": accessTokenHash(accessToken),
	})

	principal, err := verifier.Verify(
		context.Background(),
		raw,
		"0123456789abcdef0123456789abcdef",
		accessToken,
	)
	if err != nil {
		t.Fatal(err)
	}
	if principal.Subject != "subject-123" {
		t.Fatalf("unexpected principal: %#v", principal)
	}

	for name, supplied := range map[string]string{
		"missing access token": "",
		"wrong access token":   "wrong-access-token",
	} {
		t.Run(name, func(t *testing.T) {
			_, err := verifier.Verify(
				context.Background(),
				raw,
				"0123456789abcdef0123456789abcdef",
				supplied,
			)
			if !errors.Is(
				err,
				authentication.ErrAuthenticationInvalid,
			) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestVerifierUsesBoundedAuthenticationTime(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	authenticatedAt := provider.now.Add(-5 * time.Minute)

	principal, err := verifier.Verify(
		context.Background(),
		provider.token(t, map[string]any{
			"auth_time": authenticatedAt.Unix(),
		}),
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if !principal.AuthenticatedAt.Equal(authenticatedAt) {
		t.Fatalf(
			"authenticated time = %s, want %s",
			principal.AuthenticatedAt,
			authenticatedAt,
		)
	}

	_, err = verifier.Verify(
		context.Background(),
		provider.token(t, map[string]any{
			"auth_time": provider.now.Add(2 * time.Minute).Unix(),
		}),
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("future authentication time error = %v", err)
	}
}

func TestVerifierFailsClosedForInvalidProtocolState(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	rogueKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]string{
		"wrong issuer": provider.token(t, map[string]any{
			"iss": "https://issuer.invalid",
		}),
		"wrong audience": provider.token(t, map[string]any{
			"aud": "other-client",
		}),
		"missing authorized party for multiple audiences": provider.token(
			t,
			map[string]any{"aud": []string{"iron-atlas-test", "other-client"}},
		),
		"wrong authorized party": provider.token(t, map[string]any{
			"aud": []string{"iron-atlas-test", "other-client"},
			"azp": "other-client",
		}),
		"expired": provider.token(t, map[string]any{
			"exp": provider.now.Add(-2 * time.Minute).Unix(),
		}),
		"future issued at": provider.token(t, map[string]any{
			"iat": provider.now.Add(2 * time.Minute).Unix(),
		}),
		"stale issued at": provider.token(t, map[string]any{
			"iat": provider.now.Add(-20 * time.Minute).Unix(),
		}),
		"future not before": provider.token(t, map[string]any{
			"nbf": provider.now.Add(2 * time.Minute).Unix(),
		}),
		"missing subject": provider.token(t, map[string]any{"sub": nil}),
		"unnormalized subject": provider.token(t, map[string]any{
			"sub": " subject-123",
		}),
		"invalid signature": signedToken(
			t,
			rogueKey,
			"key-1",
			provider.defaultClaims(nil),
			jose.RS256,
		),
		"prohibited algorithm": signedToken(
			t,
			provider.key,
			provider.keyID,
			provider.defaultClaims(nil),
			jose.PS256,
		),
	}

	for name, raw := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := verifier.Verify(
				context.Background(),
				raw,
				"0123456789abcdef0123456789abcdef",
				"",
			)
			if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}

	t.Run("wrong nonce", func(t *testing.T) {
		_, err := verifier.Verify(
			context.Background(),
			provider.token(t, nil),
			"fedcba9876543210fedcba9876543210",
			"",
		)
		if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
			t.Fatalf("error = %v", err)
		}
	})
}

func TestVerifierRejectsDuplicateSensitiveClaim(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	payload := []byte(
		`{"iss":"` + provider.server.URL + `",` +
			`"sub":"subject-123","aud":"iron-atlas-test",` +
			`"exp":` + fmt.Sprintf("%d", provider.now.Add(5*time.Minute).Unix()) + `,` +
			`"iat":` + fmt.Sprintf("%d", provider.now.Add(-time.Minute).Unix()) + `,` +
			`"nonce":"0123456789abcdef0123456789abcdef",` +
			`"nonce":"attacker"}`,
	)
	raw := signedPayload(t, provider.key, provider.keyID, payload, jose.RS256)
	_, err := verifier.Verify(
		context.Background(),
		raw,
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestVerifierRejectsDuplicateSensitiveHeader(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)

	header := base64.RawURLEncoding.EncodeToString(
		[]byte(`{"alg":"RS256","alg":"none","kid":"key-1"}`),
	)
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{}`))
	raw := header + "." + payload + ".signature"

	_, err := verifier.Verify(
		context.Background(),
		raw,
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestVerifierRejectsOversizedOrMalformedToken(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	for _, raw := range []string{
		"not-a-jwt",
		strings.Repeat("a", defaultMaxTokenSize+1),
		" a.b.c",
	} {
		_, err := verifier.Verify(
			context.Background(),
			raw,
			"0123456789abcdef0123456789abcdef",
			"",
		)
		if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
			t.Fatalf("error = %v", err)
		}
	}
}

func TestVerifierRefreshesKeysOnRotation(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)

	if _, err := verifier.Verify(
		context.Background(),
		provider.token(t, nil),
		"0123456789abcdef0123456789abcdef",
		"",
	); err != nil {
		t.Fatal(err)
	}

	newKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	provider.mu.Lock()
	provider.key = newKey
	provider.keyID = "key-2"
	provider.mu.Unlock()

	if _, err := verifier.Verify(
		context.Background(),
		provider.token(t, nil),
		"0123456789abcdef0123456789abcdef",
		"",
	); err != nil {
		t.Fatal(err)
	}
}

func TestVerifierClassifiesKeyProviderOutage(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)

	newKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	provider.mu.Lock()
	provider.key = newKey
	provider.keyID = "unknown-while-offline"
	provider.mu.Unlock()
	raw := provider.token(t, nil)
	provider.server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err = verifier.Verify(
		ctx,
		raw,
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("error = %v", err)
	}
}

func TestVerifierClassifiesUnknownKeyAsInvalidWhenProviderResponds(
	t *testing.T,
) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	unknownKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}

	raw := signedToken(
		t,
		unknownKey,
		"unknown-but-provider-responsive",
		provider.defaultClaims(nil),
		jose.RS256,
	)
	_, err = verifier.Verify(
		context.Background(),
		raw,
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("error = %v", err)
	}
}

func TestVerifierClassifiesJWKSServiceFailureAsUnavailable(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	unknownKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	raw := signedToken(
		t,
		unknownKey,
		"unknown-while-jwks-fails",
		provider.defaultClaims(nil),
		jose.RS256,
	)

	provider.mu.Lock()
	provider.available = false
	provider.mu.Unlock()

	_, err = verifier.Verify(
		context.Background(),
		raw,
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("error = %v", err)
	}
}

func TestVerifierSupportsConcurrentReadOnlyVerification(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	raw := provider.token(t, nil)

	const operations = 100
	var wait sync.WaitGroup
	errorsFound := make(chan error, operations)
	for index := 0; index < operations; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			principal, err := verifier.Verify(
				context.Background(),
				raw,
				"0123456789abcdef0123456789abcdef",
				"",
			)
			if err != nil {
				errorsFound <- err
				return
			}
			if principal.Subject != "subject-123" {
				errorsFound <- errors.New("wrong subject")
			}
		}()
	}
	wait.Wait()
	close(errorsFound)
	for err := range errorsFound {
		t.Fatal(err)
	}
}

func TestNewRejectsInsecureOrUnboundedConfiguration(t *testing.T) {
	provider := newProviderEmulator(t)
	for name, config := range map[string]Config{
		"HTTP issuer": {
			ProviderID: "oidc-test",
			IssuerURL:  "http://issuer.example",
			ClientID:   "client",
		},
		"symmetric algorithm": {
			ProviderID:               "oidc-test",
			IssuerURL:                provider.server.URL,
			ClientID:                 "client",
			AllowedSigningAlgorithms: []string{"HS256"},
			HTTPClient:               provider.client,
		},
		"zero client timeout": {
			ProviderID: "oidc-test",
			IssuerURL:  provider.server.URL,
			ClientID:   "client",
			HTTPClient: &http.Client{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := New(context.Background(), config); err == nil {
				t.Fatal("insecure configuration was accepted")
			}
		})
	}
}

func TestVerifierNormalizesAssuranceClaims(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	authenticatedAt := provider.now.Add(-4 * time.Minute)
	principal, err := verifier.Verify(
		context.Background(),
		provider.token(t, map[string]any{
			"auth_time": authenticatedAt.Unix(),
			"acr":       "urn:example:acr:mfa",
			"amr":       []string{"pwd", "otp"},
		}),
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if principal.Assurance.Context != "urn:example:acr:mfa" {
		t.Fatalf("assurance context = %q", principal.Assurance.Context)
	}
	if len(principal.Assurance.Methods) != 2 ||
		principal.Assurance.Methods[0] != "pwd" ||
		principal.Assurance.Methods[1] != "otp" {
		t.Fatalf("assurance methods = %v", principal.Assurance.Methods)
	}
	if principal.Assurance.MFAAuthenticated ||
		!principal.Assurance.MFAAuthenticatedAt.IsZero() {
		t.Fatal("OIDC verifier must not infer Atlas MFA acceptance")
	}
}

func TestVerifierRejectsMalformedAssuranceClaims(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	for name, overrides := range map[string]map[string]any{
		"unnormalized context": {"acr": " urn:example:acr:mfa"},
		"wrong amr shape":      {"amr": "pwd"},
		"empty method":         {"amr": []string{"pwd", ""}},
		"duplicate method":     {"amr": []string{"pwd", "pwd"}},
		"too many methods": {
			"amr": []string{
				"m01", "m02", "m03", "m04", "m05", "m06", "m07", "m08", "m09",
				"m10", "m11", "m12", "m13", "m14", "m15", "m16", "m17",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := verifier.Verify(
				context.Background(),
				provider.token(t, overrides),
				"0123456789abcdef0123456789abcdef",
				"",
			)
			if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestVerifierRejectsDuplicateAssuranceClaims(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)
	for name, duplicate := range map[string]string{
		"acr": `"acr":"urn:one","acr":"urn:two"`,
		"amr": `"amr":["pwd"],"amr":["otp"]`,
	} {
		t.Run(name, func(t *testing.T) {
			payload := []byte(
				`{"iss":"` + provider.server.URL + `",` +
					`"sub":"subject-123","aud":"iron-atlas-test",` +
					`"exp":` + fmt.Sprintf("%d", provider.now.Add(5*time.Minute).Unix()) + `,` +
					`"iat":` + fmt.Sprintf("%d", provider.now.Add(-time.Minute).Unix()) + `,` +
					`"nonce":"0123456789abcdef0123456789abcdef",` + duplicate + `}`,
			)
			raw := signedPayload(t, provider.key, provider.keyID, payload, jose.RS256)
			_, err := verifier.Verify(
				context.Background(),
				raw,
				"0123456789abcdef0123456789abcdef",
				"",
			)
			if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}
}

func TestVerifierDoesNotInferMFAFromSuccessfulLogin(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)

	principal, err := verifier.Verify(
		context.Background(),
		provider.token(t, nil),
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if principal.Assurance.Context != "" ||
		len(principal.Assurance.Methods) != 0 ||
		principal.Assurance.MFAAuthenticated {
		t.Fatalf("successful login manufactured assurance: %#v", principal.Assurance)
	}
}

func TestVerifierRequiresAuthenticationTimeForAssuranceEvidence(t *testing.T) {
	provider := newProviderEmulator(t)
	verifier := provider.verifier(t)

	for name, claims := range map[string]map[string]any{
		"acr only": {
			"acr": "urn:iron-atlas:test:mfa",
		},
		"amr only": {
			"amr": []string{"pwd", "otp"},
		},
		"acr and amr": {
			"acr": "urn:iron-atlas:test:mfa",
			"amr": []string{"pwd", "otp"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := verifier.Verify(
				context.Background(),
				provider.token(t, claims),
				"0123456789abcdef0123456789abcdef",
				"",
			)
			if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
				t.Fatalf("error = %v", err)
			}
		})
	}

	authenticatedAt := provider.now.Add(-2 * time.Minute)
	principal, err := verifier.Verify(
		context.Background(),
		provider.token(t, map[string]any{
			"acr":       "urn:iron-atlas:test:mfa",
			"amr":       []string{"pwd", "otp"},
			"auth_time": authenticatedAt.Unix(),
		}),
		"0123456789abcdef0123456789abcdef",
		"",
	)
	if err != nil {
		t.Fatal(err)
	}
	if principal.Assurance.Context != "urn:iron-atlas:test:mfa" ||
		len(principal.Assurance.Methods) != 2 ||
		!principal.AuthenticatedAt.Equal(authenticatedAt) {
		t.Fatalf("unexpected assurance evidence: %#v", principal)
	}
}
