package oidc

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
)

const (
	testAuthorizationCode = "authorization-code-value"
	testClientID          = "iron-atlas-flow-test"
	testClientSecret      = "test-client-secret"
	testRedirectURL       = "https://atlas.example/auth/oidc/callback"
)

type authorizationProviderEmulator struct {
	server *httptest.Server
	client *http.Client
	now    time.Time

	mu                sync.Mutex
	key               *rsa.PrivateKey
	keyID             string
	expectedNonce     string
	expectedChallenge string
	tokenAvailable    bool
	oversizedResponse bool
	tokenCalls        int
}

func newAuthorizationProviderEmulator(
	t *testing.T,
) *authorizationProviderEmulator {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	provider := &authorizationProviderEmulator{
		now:            time.Date(2026, 7, 16, 17, 30, 0, 0, time.UTC),
		key:            key,
		keyID:          "flow-key-1",
		tokenAvailable: true,
	}
	provider.server = httptest.NewTLSServer(http.HandlerFunc(provider.handle))
	provider.client = provider.server.Client()
	provider.client.Timeout = 2 * time.Second
	t.Cleanup(provider.server.Close)
	return provider
}

func (p *authorizationProviderEmulator) handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	switch r.URL.Path {
	case "/.well-known/openid-configuration":
		writeTestJSON(w, map[string]any{
			"issuer":                                p.server.URL,
			"authorization_endpoint":                p.server.URL + "/authorize",
			"token_endpoint":                        p.server.URL + "/token",
			"jwks_uri":                              p.server.URL + "/jwks",
			"response_types_supported":              []string{"code"},
			"subject_types_supported":               []string{"public"},
			"id_token_signing_alg_values_supported": []string{"RS256"},
			"token_endpoint_auth_methods_supported": []string{"client_secret_basic"},
			"code_challenge_methods_supported":      []string{"S256"},
		})
	case "/jwks":
		p.mu.Lock()
		key := p.key
		keyID := p.keyID
		p.mu.Unlock()
		writeTestJSON(w, jose.JSONWebKeySet{
			Keys: []jose.JSONWebKey{{
				Key:       &key.PublicKey,
				KeyID:     keyID,
				Algorithm: string(jose.RS256),
				Use:       "sig",
			}},
		})
	case "/token":
		p.handleToken(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *authorizationProviderEmulator) handleToken(
	w http.ResponseWriter,
	r *http.Request,
) {
	p.mu.Lock()
	p.tokenCalls++
	available := p.tokenAvailable
	oversized := p.oversizedResponse
	expectedNonce := p.expectedNonce
	expectedChallenge := p.expectedChallenge
	key := p.key
	keyID := p.keyID
	p.mu.Unlock()

	if !available {
		http.Error(w, "temporarily unavailable", http.StatusServiceUnavailable)
		return
	}
	if oversized {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(
			w,
			strings.Repeat("x", defaultMaximumTokenResponse+1024),
		)
		return
	}
	if err := r.ParseForm(); err != nil {
		writeTestJSON(w, map[string]any{"error": "invalid_request"})
		return
	}
	clientID, clientSecret, basic := r.BasicAuth()
	if !basic ||
		clientID != testClientID ||
		clientSecret != testClientSecret ||
		r.Form.Get("grant_type") != "authorization_code" ||
		r.Form.Get("code") != testAuthorizationCode ||
		r.Form.Get("redirect_uri") != testRedirectURL ||
		pkceChallenge(r.Form.Get("code_verifier")) != expectedChallenge ||
		expectedNonce == "" {
		w.WriteHeader(http.StatusBadRequest)
		writeTestJSON(w, map[string]any{"error": "invalid_grant"})
		return
	}

	accessToken := "opaque-access-token"
	rawIDToken, err := signAuthorizationToken(
		key,
		keyID,
		map[string]any{
			"iss":     p.server.URL,
			"sub":     "subject-flow-123",
			"aud":     testClientID,
			"exp":     p.now.Add(5 * time.Minute).Unix(),
			"iat":     p.now.Add(-time.Minute).Unix(),
			"nonce":   expectedNonce,
			"at_hash": accessTokenHash(accessToken),
		},
	)
	if err != nil {
		http.Error(w, "token signing failed", http.StatusInternalServerError)
		return
	}
	writeTestJSON(w, map[string]any{
		"access_token": accessToken,
		"token_type":   "Bearer",
		"expires_in":   300,
		"id_token":     rawIDToken,
	})
}

func signAuthorizationToken(
	key *rsa.PrivateKey,
	keyID string,
	claims map[string]any,
) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).
			WithType("JWT").
			WithHeader("kid", keyID),
	)
	if err != nil {
		return "", err
	}
	object, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}
	return object.CompactSerialize()
}

func (p *authorizationProviderEmulator) acceptAuthorization(
	t *testing.T,
	request AuthorizationRequest,
) url.Values {
	t.Helper()
	parsed, err := url.Parse(request.AuthorizationURL)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Scheme != "https" ||
		parsed.Host != strings.TrimPrefix(p.server.URL, "https://") ||
		parsed.Path != "/authorize" {
		t.Fatalf("unexpected authorization endpoint: %s", parsed.String())
	}
	query := parsed.Query()
	if query.Get("state") != request.State {
		t.Fatalf("authorization state does not match returned state")
	}
	if query.Get("code_challenge_method") != "S256" {
		t.Fatalf(
			"code_challenge_method = %q",
			query.Get("code_challenge_method"),
		)
	}
	if err := validateRandomToken("nonce", query.Get("nonce")); err != nil {
		t.Fatalf("invalid nonce: %v", err)
	}
	if err := validateRandomToken(
		"code challenge",
		query.Get("code_challenge"),
	); err != nil {
		t.Fatalf("invalid code challenge: %v", err)
	}

	p.mu.Lock()
	p.expectedNonce = query.Get("nonce")
	p.expectedChallenge = query.Get("code_challenge")
	p.mu.Unlock()
	return query
}

func (p *authorizationProviderEmulator) calls() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.tokenCalls
}

func (p *authorizationProviderEmulator) setTokenAvailable(value bool) {
	p.mu.Lock()
	p.tokenAvailable = value
	p.mu.Unlock()
}

func (p *authorizationProviderEmulator) setOversizedResponse(value bool) {
	p.mu.Lock()
	p.oversizedResponse = value
	p.mu.Unlock()
}

func newTestAuthorizationFlow(
	t *testing.T,
	provider *authorizationProviderEmulator,
	now func() time.Time,
	ttl time.Duration,
	capacity int,
) (*AuthorizationCodeFlow, *MemoryPreauthenticationStore) {
	t.Helper()
	verifier, err := New(context.Background(), Config{
		ProviderID:               "oidc-flow-test",
		IssuerURL:                provider.server.URL,
		ClientID:                 testClientID,
		AllowedSigningAlgorithms: []string{"RS256"},
		HTTPClient:               provider.client,
		Now:                      func() time.Time { return provider.now },
		MaxClockSkew:             time.Minute,
		MaxTokenAge:              15 * time.Minute,
	})
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewMemoryPreauthenticationStore(capacity)
	if err != nil {
		t.Fatal(err)
	}
	if now == nil {
		now = func() time.Time { return provider.now }
	}
	flow, err := NewAuthorizationCodeFlow(AuthorizationCodeFlowConfig{
		Verifier:             verifier,
		RedirectURL:          testRedirectURL,
		ClientSecret:         testClientSecret,
		Scopes:               []string{"profile", "email"},
		Store:                store,
		Now:                  now,
		PreauthenticationTTL: ttl,
	})
	if err != nil {
		t.Fatal(err)
	}
	return flow, store
}

func pkceChallenge(verifier string) string {
	digest := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(digest[:])
}

func TestAuthorizationCodeFlowBeginsWithBoundStateNonceAndPKCE(
	t *testing.T,
) {
	provider := newAuthorizationProviderEmulator(t)
	flow, store := newTestAuthorizationFlow(
		t,
		provider,
		nil,
		5*time.Minute,
		64,
	)

	request, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	query := provider.acceptAuthorization(t, request)

	if query.Get("response_type") != "code" ||
		query.Get("client_id") != testClientID ||
		query.Get("redirect_uri") != testRedirectURL {
		t.Fatalf("unexpected authorization query: %v", query)
	}
	if query.Get("client_secret") != "" {
		t.Fatal("client secret appeared in authorization URL")
	}
	if query.Get("scope") != "openid profile email" {
		t.Fatalf("scope = %q", query.Get("scope"))
	}
	if err := validateRandomToken("state", request.State); err != nil {
		t.Fatalf("invalid state: %v", err)
	}
	if !request.ExpiresAt.Equal(provider.now.Add(5 * time.Minute)) {
		t.Fatalf("expiry = %s", request.ExpiresAt)
	}

	digest := sha256.Sum256([]byte(request.State))
	store.mu.Lock()
	transaction, exists := store.records[digest]
	store.mu.Unlock()
	if !exists {
		t.Fatal("preauthentication transaction was not stored")
	}
	if transaction.Nonce == request.State ||
		transaction.PKCEVerifier == request.State {
		t.Fatal("raw state was retained as another transaction secret")
	}
	if transaction.StateDigest != digest {
		t.Fatal("stored state digest does not match")
	}
}

func TestAuthorizationCodeFlowCompletesExactlyOnce(t *testing.T) {
	provider := newAuthorizationProviderEmulator(t)
	flow, _ := newTestAuthorizationFlow(
		t,
		provider,
		nil,
		5*time.Minute,
		64,
	)

	request, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, request)

	principal, err := flow.Complete(
		context.Background(),
		request.State,
		testAuthorizationCode,
	)
	if err != nil {
		t.Fatal(err)
	}
	if principal.ProviderID != "oidc-flow-test" ||
		principal.Subject != "subject-flow-123" ||
		!principal.AuthenticatedAt.Equal(provider.now.Add(-time.Minute)) {
		t.Fatalf("unexpected principal: %#v", principal)
	}

	_, err = flow.Complete(
		context.Background(),
		request.State,
		testAuthorizationCode,
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("replay error = %v", err)
	}
	if provider.calls() != 1 {
		t.Fatalf("token calls = %d, want 1", provider.calls())
	}
}

func TestAuthorizationCodeFlowRejectsUnknownAndExpiredState(t *testing.T) {
	provider := newAuthorizationProviderEmulator(t)
	current := provider.now
	flow, _ := newTestAuthorizationFlow(
		t,
		provider,
		func() time.Time { return current },
		time.Minute,
		64,
	)

	unknownState, err := randomToken(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, err = flow.Complete(
		context.Background(),
		unknownState,
		testAuthorizationCode,
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("unknown-state error = %v", err)
	}

	request, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, request)
	current = current.Add(2 * time.Minute)

	_, err = flow.Complete(
		context.Background(),
		request.State,
		testAuthorizationCode,
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("expired-state error = %v", err)
	}
	if provider.calls() != 0 {
		t.Fatalf("token endpoint called %d times", provider.calls())
	}
}

func TestAuthorizationCodeFlowAllowsOnlyOneConcurrentConsumer(
	t *testing.T,
) {
	provider := newAuthorizationProviderEmulator(t)
	flow, _ := newTestAuthorizationFlow(
		t,
		provider,
		nil,
		5*time.Minute,
		64,
	)
	request, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, request)

	const attempts = 50
	var wait sync.WaitGroup
	results := make(chan error, attempts)
	for index := 0; index < attempts; index++ {
		wait.Add(1)
		go func() {
			defer wait.Done()
			_, err := flow.Complete(
				context.Background(),
				request.State,
				testAuthorizationCode,
			)
			results <- err
		}()
	}
	wait.Wait()
	close(results)

	successes := 0
	invalid := 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, authentication.ErrAuthenticationInvalid):
			invalid++
		default:
			t.Fatalf("unexpected concurrent error: %v", err)
		}
	}
	if successes != 1 || invalid != attempts-1 {
		t.Fatalf("successes=%d invalid=%d", successes, invalid)
	}
	if provider.calls() != 1 {
		t.Fatalf("token calls = %d, want 1", provider.calls())
	}
}

func TestAuthorizationCodeFlowClassifiesInvalidCodeAndProviderOutage(
	t *testing.T,
) {
	provider := newAuthorizationProviderEmulator(t)
	flow, _ := newTestAuthorizationFlow(
		t,
		provider,
		nil,
		5*time.Minute,
		64,
	)

	invalidRequest, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, invalidRequest)
	_, err = flow.Complete(
		context.Background(),
		invalidRequest.State,
		"wrong-authorization-code",
	)
	if !errors.Is(err, authentication.ErrAuthenticationInvalid) {
		t.Fatalf("invalid-code error = %v", err)
	}

	outageRequest, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, outageRequest)
	provider.setTokenAvailable(false)
	_, err = flow.Complete(
		context.Background(),
		outageRequest.State,
		testAuthorizationCode,
	)
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("outage error = %v", err)
	}
}

func TestAuthorizationCodeFlowBoundsTokenResponse(t *testing.T) {
	provider := newAuthorizationProviderEmulator(t)
	flow, _ := newTestAuthorizationFlow(
		t,
		provider,
		nil,
		5*time.Minute,
		64,
	)
	request, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, request)
	provider.setOversizedResponse(true)

	_, err = flow.Complete(
		context.Background(),
		request.State,
		testAuthorizationCode,
	)
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("oversized-response error = %v", err)
	}
}

func TestAuthorizationCodeFlowDoesNotExposeSecretsInErrors(t *testing.T) {
	provider := newAuthorizationProviderEmulator(t)
	flow, _ := newTestAuthorizationFlow(
		t,
		provider,
		nil,
		5*time.Minute,
		64,
	)
	request, err := flow.Begin(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	provider.acceptAuthorization(t, request)

	sensitiveCode := "sensitive-authorization-code"
	_, err = flow.Complete(
		context.Background(),
		request.State,
		sensitiveCode,
	)
	if err == nil {
		t.Fatal("expected failure")
	}
	message := err.Error()
	for _, secret := range []string{
		request.State,
		sensitiveCode,
		testClientSecret,
	} {
		if strings.Contains(message, secret) {
			t.Fatalf("error exposed a secret: %q", message)
		}
	}
}

func TestMemoryPreauthenticationStoreIsBoundedAndCleansExpiredState(
	t *testing.T,
) {
	store, err := NewMemoryPreauthenticationStore(1)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 7, 16, 18, 0, 0, 0, time.UTC)
	first := testPreauthenticationTransaction(t, now, time.Minute)
	second := testPreauthenticationTransaction(
		t,
		now.Add(time.Second),
		time.Minute,
	)

	if err := store.Create(context.Background(), first); err != nil {
		t.Fatal(err)
	}
	if err := store.Create(context.Background(), second); !errors.Is(
		err,
		ErrPreauthenticationCapacity,
	) {
		t.Fatalf("capacity error = %v", err)
	}
	if _, err := store.Consume(
		context.Background(),
		first.StateDigest,
		now,
	); err != nil {
		t.Fatal(err)
	}
	if err := store.Create(context.Background(), second); err != nil {
		t.Fatal(err)
	}

	third := testPreauthenticationTransaction(
		t,
		now.Add(2*time.Minute),
		time.Minute,
	)
	if err := store.Create(context.Background(), third); err != nil {
		t.Fatalf("expired cleanup did not free capacity: %v", err)
	}
}

func TestNewAuthorizationCodeFlowRejectsUnsafeConfiguration(t *testing.T) {
	provider := newAuthorizationProviderEmulator(t)
	verifier, err := New(context.Background(), Config{
		ProviderID:               "oidc-flow-test",
		IssuerURL:                provider.server.URL,
		ClientID:                 testClientID,
		AllowedSigningAlgorithms: []string{"RS256"},
		HTTPClient:               provider.client,
		Now:                      func() time.Time { return provider.now },
	})
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewMemoryPreauthenticationStore(4)
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]AuthorizationCodeFlowConfig{
		"missing store": {
			Verifier:     verifier,
			RedirectURL:  testRedirectURL,
			ClientSecret: testClientSecret,
		},
		"HTTP redirect": {
			Verifier:     verifier,
			RedirectURL:  "http://atlas.example/callback",
			ClientSecret: testClientSecret,
			Store:        store,
		},
		"unbounded lifetime": {
			Verifier:             verifier,
			RedirectURL:          testRedirectURL,
			ClientSecret:         testClientSecret,
			Store:                store,
			PreauthenticationTTL: 30 * time.Minute,
		},
		"duplicate scope": {
			Verifier:     verifier,
			RedirectURL:  testRedirectURL,
			ClientSecret: testClientSecret,
			Store:        store,
			Scopes:       []string{"profile", "profile"},
		},
	}

	for name, config := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := NewAuthorizationCodeFlow(config); err == nil {
				t.Fatal("expected configuration error")
			}
		})
	}
}

func TestAuthorizationCodeFlowFailsClosedWhenRandomnessFails(
	t *testing.T,
) {
	provider := newAuthorizationProviderEmulator(t)
	verifier, err := New(context.Background(), Config{
		ProviderID:               "oidc-flow-test",
		IssuerURL:                provider.server.URL,
		ClientID:                 testClientID,
		AllowedSigningAlgorithms: []string{"RS256"},
		HTTPClient:               provider.client,
		Now:                      func() time.Time { return provider.now },
	})
	if err != nil {
		t.Fatal(err)
	}
	store, err := NewMemoryPreauthenticationStore(4)
	if err != nil {
		t.Fatal(err)
	}
	flow, err := NewAuthorizationCodeFlow(AuthorizationCodeFlowConfig{
		Verifier:     verifier,
		RedirectURL:  testRedirectURL,
		ClientSecret: testClientSecret,
		Store:        store,
		Random:       failingReader{},
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = flow.Begin(context.Background())
	if !errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		t.Fatalf("randomness error = %v", err)
	}
}

func testPreauthenticationTransaction(
	t *testing.T,
	created time.Time,
	lifetime time.Duration,
) PreauthenticationTransaction {
	t.Helper()
	state, err := randomToken(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	nonce, err := randomToken(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	verifier, err := randomToken(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return PreauthenticationTransaction{
		StateDigest:  sha256.Sum256([]byte(state)),
		Nonce:        nonce,
		PKCEVerifier: verifier,
		RedirectURL:  testRedirectURL,
		CreatedAt:    created,
		ExpiresAt:    created.Add(lifetime),
	}
}

type failingReader struct{}

func (failingReader) Read(_ []byte) (int, error) {
	return 0, errors.New("random source unavailable")
}
