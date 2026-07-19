package oidc

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	coreoidc "github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
)

const (
	defaultPreauthenticationTTL       = 5 * time.Minute
	maximumPreauthenticationTTL       = 10 * time.Minute
	defaultMaximumAuthorizationCode   = 8 << 10
	defaultMaximumTokenResponse       = 1 << 20
	defaultMaximumPreauthentication   = 4096
	authorizationRandomTokenBytes     = 32
	maximumAuthorizationURLBytes      = 8 << 10
	maximumClientSecretBytes          = 4 << 10
	maximumScopeCount                 = 16
	maximumScopeBytes                 = 128
	maximumAuthorizationCodeBytesHard = 64 << 10
	maximumTokenResponseBytesHard     = 8 << 20
)

var (
	ErrPreauthenticationInvalid = errors.New(
		"preauthentication transaction is invalid or unavailable",
	)
	ErrPreauthenticationCapacity = errors.New(
		"preauthentication transaction capacity is exhausted",
	)
	ErrPreauthenticationUnavailable = errors.New(
		"preauthentication transaction store is unavailable",
	)
	errTokenResponseTooLarge = errors.New("OIDC token response is too large")
)

// PreauthenticationTransaction is the server-side state required to bind one
// authorization-code callback to one initiation attempt.
//
// StateDigest is a SHA-256 digest. The raw state value is never retained by the
// store. Nonce and PKCEVerifier are short-lived secrets and must not be logged.
type PreauthenticationTransaction struct {
	StateDigest  [sha256.Size]byte
	Nonce        string
	PKCEVerifier string
	RedirectURL  string
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func (t PreauthenticationTransaction) validate() error {
	if t.StateDigest == ([sha256.Size]byte{}) {
		return errors.New("state digest is required")
	}
	if err := validateRandomToken("nonce", t.Nonce); err != nil {
		return err
	}
	if err := validateRandomToken("PKCE verifier", t.PKCEVerifier); err != nil {
		return err
	}
	if _, err := trustedHTTPSURL("redirect URL", t.RedirectURL); err != nil {
		return err
	}
	if t.CreatedAt.IsZero() || t.ExpiresAt.IsZero() {
		return errors.New("preauthentication timestamps are required")
	}
	if !t.ExpiresAt.After(t.CreatedAt) {
		return errors.New("preauthentication expiry must follow creation")
	}
	if t.ExpiresAt.Sub(t.CreatedAt) > maximumPreauthenticationTTL {
		return errors.New("preauthentication lifetime exceeds the maximum")
	}
	return nil
}

// PreauthenticationStore provides atomic one-time transaction creation and
// consumption. A successful Consume must make every later consume of the same
// digest fail closed.
type PreauthenticationStore interface {
	Create(context.Context, PreauthenticationTransaction) error
	Consume(
		context.Context,
		[sha256.Size]byte,
		time.Time,
	) (PreauthenticationTransaction, error)
}

// MemoryPreauthenticationStore is the bounded in-memory candidate store.
//
// Restart invalidates all outstanding transactions. Durable restart-surviving
// storage is deliberately not claimed by this checkpoint.
type MemoryPreauthenticationStore struct {
	mu         sync.Mutex
	maxEntries int
	records    map[[sha256.Size]byte]PreauthenticationTransaction
}

func NewMemoryPreauthenticationStore(
	maxEntries int,
) (*MemoryPreauthenticationStore, error) {
	if maxEntries == 0 {
		maxEntries = defaultMaximumPreauthentication
	}
	if maxEntries < 1 || maxEntries > 1_000_000 {
		return nil, errors.New(
			"preauthentication capacity must be between one and one million",
		)
	}
	return &MemoryPreauthenticationStore{
		maxEntries: maxEntries,
		records: make(
			map[[sha256.Size]byte]PreauthenticationTransaction,
			maxEntries,
		),
	}, nil
}

func (s *MemoryPreauthenticationStore) Create(
	ctx context.Context,
	transaction PreauthenticationTransaction,
) error {
	if s == nil {
		return ErrPreauthenticationUnavailable
	}
	if ctx == nil {
		return ErrPreauthenticationUnavailable
	}
	if err := ctx.Err(); err != nil {
		return ErrPreauthenticationUnavailable
	}
	if err := transaction.validate(); err != nil {
		return ErrPreauthenticationInvalid
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.purgeExpired(transaction.CreatedAt.UTC())
	if _, exists := s.records[transaction.StateDigest]; exists {
		return ErrPreauthenticationInvalid
	}
	if len(s.records) >= s.maxEntries {
		return ErrPreauthenticationCapacity
	}
	s.records[transaction.StateDigest] = transaction
	return nil
}

func (s *MemoryPreauthenticationStore) Consume(
	ctx context.Context,
	stateDigest [sha256.Size]byte,
	now time.Time,
) (PreauthenticationTransaction, error) {
	if s == nil {
		return PreauthenticationTransaction{}, ErrPreauthenticationUnavailable
	}
	if ctx == nil {
		return PreauthenticationTransaction{}, ErrPreauthenticationUnavailable
	}
	if err := ctx.Err(); err != nil {
		return PreauthenticationTransaction{}, ErrPreauthenticationUnavailable
	}
	if stateDigest == ([sha256.Size]byte{}) || now.IsZero() {
		return PreauthenticationTransaction{}, ErrPreauthenticationInvalid
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.purgeExpired(now.UTC())
	transaction, exists := s.records[stateDigest]
	if !exists {
		return PreauthenticationTransaction{}, ErrPreauthenticationInvalid
	}
	delete(s.records, stateDigest)

	if err := transaction.validate(); err != nil {
		return PreauthenticationTransaction{}, ErrPreauthenticationInvalid
	}
	current := now.UTC()
	if current.Before(transaction.CreatedAt.UTC()) ||
		!transaction.ExpiresAt.After(current) {
		return PreauthenticationTransaction{}, ErrPreauthenticationInvalid
	}
	return transaction, nil
}

func (s *MemoryPreauthenticationStore) purgeExpired(now time.Time) {
	for digest, transaction := range s.records {
		if !transaction.ExpiresAt.After(now) {
			delete(s.records, digest)
		}
	}
}

// AuthorizationCodeFlowConfig defines one bounded OIDC authorization-code and
// PKCE transaction profile.
type AuthorizationCodeFlowConfig struct {
	Verifier                  *Verifier
	RedirectURL               string
	ClientSecret              string
	Scopes                    []string
	Store                     PreauthenticationStore
	Random                    io.Reader
	Now                       func() time.Time
	PreauthenticationTTL      time.Duration
	MaxAuthorizationCodeBytes int
	MaxTokenResponseBytes     int64
}

// AuthorizationRequest is the result of creating one short-lived,
// single-purpose browser authorization request.
//
// State is returned because the later HTTP boundary must bind it to the same
// browser initiation. State must never be logged or placed in a persistent URL.
type AuthorizationRequest struct {
	AuthorizationURL string
	State            string
	ExpiresAt        time.Time
}

// AuthorizationCodeFlow creates one-time state, nonce, and PKCE transactions,
// exchanges one authorization code, and verifies the returned ID token.
//
// This component does not implement HTTP login or callback routes, cookies,
// durable sessions, CSRF, logout, trusted-proxy enforcement, or actor
// resolution.
type AuthorizationCodeFlow struct {
	verifier                  *Verifier
	store                     PreauthenticationStore
	oauthConfig               oauth2.Config
	httpClient                *http.Client
	random                    io.Reader
	now                       func() time.Time
	preauthenticationTTL      time.Duration
	maxAuthorizationCodeBytes int
}

func NewAuthorizationCodeFlow(
	config AuthorizationCodeFlowConfig,
) (*AuthorizationCodeFlow, error) {
	if config.Verifier == nil ||
		config.Verifier.tokenVerifier == nil ||
		config.Verifier.httpClient == nil {
		return nil, errors.New("OIDC verifier is required")
	}
	if config.Store == nil {
		return nil, errors.New("preauthentication store is required")
	}
	if !contains(config.Verifier.codeChallengeMethods, "S256") {
		return nil, errors.New(
			"provider does not advertise PKCE S256 support",
		)
	}
	redirectURL, err := trustedHTTPSURL("redirect URL", config.RedirectURL)
	if err != nil {
		return nil, err
	}
	if len(config.ClientSecret) > maximumClientSecretBytes ||
		!utf8.ValidString(config.ClientSecret) {
		return nil, errors.New("OIDC client secret is invalid or too large")
	}
	for _, r := range config.ClientSecret {
		if unicode.IsControl(r) {
			return nil, errors.New("OIDC client secret contains a control character")
		}
	}

	scopes, err := normalizedScopes(config.Scopes)
	if err != nil {
		return nil, err
	}
	authStyle, err := tokenEndpointAuthStyle(
		config.Verifier.tokenEndpointAuthMethods,
		config.ClientSecret != "",
	)
	if err != nil {
		return nil, err
	}

	now := config.Now
	if now == nil {
		now = time.Now
	}
	randomSource := config.Random
	if randomSource == nil {
		randomSource = rand.Reader
	}
	ttl := config.PreauthenticationTTL
	if ttl == 0 {
		ttl = defaultPreauthenticationTTL
	}
	if ttl < time.Minute || ttl > maximumPreauthenticationTTL {
		return nil, errors.New(
			"preauthentication lifetime must be between one and ten minutes",
		)
	}
	maximumCodeBytes := config.MaxAuthorizationCodeBytes
	if maximumCodeBytes == 0 {
		maximumCodeBytes = defaultMaximumAuthorizationCode
	}
	if maximumCodeBytes < 128 ||
		maximumCodeBytes > maximumAuthorizationCodeBytesHard {
		return nil, errors.New(
			"maximum authorization-code size must be between 128 bytes and 64 KiB",
		)
	}
	maximumTokenBytes := config.MaxTokenResponseBytes
	if maximumTokenBytes == 0 {
		maximumTokenBytes = defaultMaximumTokenResponse
	}
	if maximumTokenBytes < 1024 ||
		maximumTokenBytes > maximumTokenResponseBytesHard {
		return nil, errors.New(
			"maximum token-response size must be between 1 KiB and 8 MiB",
		)
	}

	httpClient := boundedHTTPClient(
		config.Verifier.httpClient,
		maximumTokenBytes,
	)

	return &AuthorizationCodeFlow{
		verifier: config.Verifier,
		store:    config.Store,
		oauthConfig: oauth2.Config{
			ClientID:     config.Verifier.clientID,
			ClientSecret: config.ClientSecret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   config.Verifier.authorizationEndpoint,
				TokenURL:  config.Verifier.tokenEndpoint,
				AuthStyle: authStyle,
			},
			RedirectURL: redirectURL,
			Scopes:      scopes,
		},
		httpClient:                httpClient,
		random:                    randomSource,
		now:                       now,
		preauthenticationTTL:      ttl,
		maxAuthorizationCodeBytes: maximumCodeBytes,
	}, nil
}

// Begin creates a one-time preauthentication transaction and returns the exact
// provider authorization URL.
func (f *AuthorizationCodeFlow) Begin(
	ctx context.Context,
) (AuthorizationRequest, error) {
	if f == nil || f.store == nil || f.random == nil || f.now == nil {
		return AuthorizationRequest{}, authentication.ErrAuthenticationUnavailable
	}
	if ctx == nil {
		return AuthorizationRequest{}, authentication.ErrAuthenticationInvalid
	}
	if err := ctx.Err(); err != nil {
		return AuthorizationRequest{}, authentication.ErrAuthenticationUnavailable
	}

	state, err := randomToken(f.random)
	if err != nil {
		return AuthorizationRequest{}, authentication.ErrAuthenticationUnavailable
	}
	nonce, err := randomToken(f.random)
	if err != nil {
		return AuthorizationRequest{}, authentication.ErrAuthenticationUnavailable
	}
	pkceVerifier, err := randomToken(f.random)
	if err != nil {
		return AuthorizationRequest{}, authentication.ErrAuthenticationUnavailable
	}

	authorizationURL := f.oauthConfig.AuthCodeURL(
		state,
		coreoidc.Nonce(nonce),
		oauth2.S256ChallengeOption(pkceVerifier),
	)
	if len(authorizationURL) > maximumAuthorizationURLBytes {
		return AuthorizationRequest{}, authentication.ErrAuthenticationInvalid
	}

	now := f.now().UTC()
	expiresAt := now.Add(f.preauthenticationTTL)
	transaction := PreauthenticationTransaction{
		StateDigest:  sha256.Sum256([]byte(state)),
		Nonce:        nonce,
		PKCEVerifier: pkceVerifier,
		RedirectURL:  f.oauthConfig.RedirectURL,
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
	}
	if err := f.store.Create(ctx, transaction); err != nil {
		if errors.Is(err, ErrPreauthenticationInvalid) {
			return AuthorizationRequest{}, authentication.ErrAuthenticationInvalid
		}
		return AuthorizationRequest{}, authentication.ErrAuthenticationUnavailable
	}

	return AuthorizationRequest{
		AuthorizationURL: authorizationURL,
		State:            state,
		ExpiresAt:        expiresAt,
	}, nil
}

// Complete atomically consumes one preauthentication transaction, exchanges one
// authorization code with its exact PKCE verifier and redirect URI, and verifies
// the returned ID token.
func (f *AuthorizationCodeFlow) Complete(
	ctx context.Context,
	state string,
	authorizationCode string,
) (authentication.Principal, error) {
	if f == nil ||
		f.store == nil ||
		f.verifier == nil ||
		f.httpClient == nil ||
		f.now == nil {
		return authentication.Principal{}, authentication.ErrAuthenticationUnavailable
	}
	if ctx == nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	if err := ctx.Err(); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationUnavailable
	}
	if err := validateRandomToken("state", state); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	if err := validateAuthorizationCode(
		authorizationCode,
		f.maxAuthorizationCodeBytes,
	); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	now := f.now().UTC()
	transaction, err := f.store.Consume(
		ctx,
		sha256.Sum256([]byte(state)),
		now,
	)
	if err != nil {
		if errors.Is(err, ErrPreauthenticationUnavailable) ||
			errors.Is(err, ErrPreauthenticationCapacity) {
			return authentication.Principal{}, authentication.ErrAuthenticationUnavailable
		}
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	if transaction.RedirectURL != f.oauthConfig.RedirectURL {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	exchangeContext := context.WithValue(
		ctx,
		oauth2.HTTPClient,
		f.httpClient,
	)
	token, err := f.oauthConfig.Exchange(
		exchangeContext,
		authorizationCode,
		oauth2.VerifierOption(transaction.PKCEVerifier),
	)
	if err != nil {
		return authentication.Principal{}, classifyExchangeError(err)
	}
	if token == nil ||
		token.AccessToken == "" ||
		len(token.AccessToken) > f.verifier.maxTokenBytes ||
		!utf8.ValidString(token.AccessToken) ||
		!strings.EqualFold(token.TokenType, "Bearer") {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok || rawIDToken == "" {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	principal, err := f.verifier.Verify(
		ctx,
		rawIDToken,
		transaction.Nonce,
		token.AccessToken,
	)
	if err != nil {
		if errors.Is(err, authentication.ErrAuthenticationUnavailable) {
			return authentication.Principal{}, authentication.ErrAuthenticationUnavailable
		}
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	return principal, nil
}

// Cancel atomically consumes one valid preauthentication transaction without
// exchanging an authorization code. It is used when the provider returns a
// bounded error callback or when callback issuer binding fails.
func (f *AuthorizationCodeFlow) Cancel(
	ctx context.Context,
	state string,
) error {
	if f == nil || f.store == nil || f.now == nil {
		return authentication.ErrAuthenticationUnavailable
	}
	if ctx == nil {
		return authentication.ErrAuthenticationInvalid
	}
	if err := ctx.Err(); err != nil {
		return authentication.ErrAuthenticationUnavailable
	}
	if err := validateRandomToken("state", state); err != nil {
		return authentication.ErrAuthenticationInvalid
	}

	transaction, err := f.store.Consume(
		ctx,
		sha256.Sum256([]byte(state)),
		f.now().UTC(),
	)
	if err != nil {
		if errors.Is(err, ErrPreauthenticationUnavailable) ||
			errors.Is(err, ErrPreauthenticationCapacity) {
			return authentication.ErrAuthenticationUnavailable
		}
		return authentication.ErrAuthenticationInvalid
	}
	if transaction.RedirectURL != f.oauthConfig.RedirectURL {
		return authentication.ErrAuthenticationInvalid
	}
	return nil
}

// IssuerURL returns the exact trusted issuer bound to this flow.
func (f *AuthorizationCodeFlow) IssuerURL() string {
	if f == nil || f.verifier == nil {
		return ""
	}
	return f.verifier.issuerURL
}

func randomToken(source io.Reader) (string, error) {
	raw := make([]byte, authorizationRandomTokenBytes)
	if _, err := io.ReadFull(source, raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}

func validateRandomToken(name string, value string) error {
	if value == "" ||
		value != strings.TrimSpace(value) ||
		!utf8.ValidString(value) {
		return fmt.Errorf("%s is missing or unnormalized", name)
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil ||
		len(decoded) != authorizationRandomTokenBytes ||
		base64.RawURLEncoding.EncodeToString(decoded) != value {
		return fmt.Errorf("%s is not a canonical 256-bit base64url value", name)
	}
	return nil
}

func validateAuthorizationCode(code string, maximum int) error {
	if code == "" ||
		code != strings.TrimSpace(code) ||
		len(code) > maximum ||
		!utf8.ValidString(code) {
		return errors.New("authorization code is missing, unnormalized, or too large")
	}
	for _, r := range code {
		if unicode.IsSpace(r) || unicode.IsControl(r) {
			return errors.New("authorization code contains whitespace or control data")
		}
	}
	return nil
}

func normalizedScopes(values []string) ([]string, error) {
	if len(values) > maximumScopeCount {
		return nil, errors.New("too many OIDC scopes")
	}
	scopes := make([]string, 0, len(values)+1)
	seen := make(map[string]struct{}, len(values)+1)
	scopes = append(scopes, coreoidc.ScopeOpenID)
	seen[coreoidc.ScopeOpenID] = struct{}{}

	for _, value := range values {
		if value == coreoidc.ScopeOpenID {
			continue
		}
		if value == "" ||
			value != strings.TrimSpace(value) ||
			len(value) > maximumScopeBytes ||
			!utf8.ValidString(value) {
			return nil, errors.New("OIDC scope is missing, unnormalized, or too large")
		}
		for _, r := range value {
			if unicode.IsSpace(r) || unicode.IsControl(r) {
				return nil, errors.New("OIDC scope contains whitespace or control data")
			}
		}
		if _, duplicate := seen[value]; duplicate {
			return nil, fmt.Errorf("duplicate OIDC scope %q", value)
		}
		seen[value] = struct{}{}
		scopes = append(scopes, value)
		if len(scopes) > maximumScopeCount {
			return nil, errors.New("too many OIDC scopes")
		}
	}
	return scopes, nil
}

func tokenEndpointAuthStyle(
	methods []string,
	hasClientSecret bool,
) (oauth2.AuthStyle, error) {
	if hasClientSecret {
		if contains(methods, "client_secret_basic") {
			return oauth2.AuthStyleInHeader, nil
		}
		if contains(methods, "client_secret_post") {
			return oauth2.AuthStyleInParams, nil
		}
		return 0, errors.New(
			"provider advertises no supported client-secret authentication method",
		)
	}
	if contains(methods, "none") {
		return oauth2.AuthStyleInParams, nil
	}
	return 0, errors.New(
		"provider requires client authentication but no client secret is configured",
	)
}

func classifyExchangeError(err error) error {
	if errors.Is(err, errTokenResponseTooLarge) ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return authentication.ErrAuthenticationUnavailable
	}

	var networkError net.Error
	if errors.As(err, &networkError) {
		return authentication.ErrAuthenticationUnavailable
	}

	var retrieveError *oauth2.RetrieveError
	if errors.As(err, &retrieveError) {
		if retrieveError.Response != nil &&
			retrieveError.Response.StatusCode >= http.StatusInternalServerError {
			return authentication.ErrAuthenticationUnavailable
		}
		switch retrieveError.ErrorCode {
		case "server_error", "temporarily_unavailable", "invalid_client":
			return authentication.ErrAuthenticationUnavailable
		default:
			return authentication.ErrAuthenticationInvalid
		}
	}
	return authentication.ErrAuthenticationUnavailable
}

func boundedHTTPClient(source *http.Client, maximum int64) *http.Client {
	clone := *source
	base := clone.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	clone.Transport = responseLimitTransport{
		base:    base,
		maximum: maximum,
	}
	clone.CheckRedirect = func(
		_ *http.Request,
		_ []*http.Request,
	) error {
		return http.ErrUseLastResponse
	}
	return &clone
}

type responseLimitTransport struct {
	base    http.RoundTripper
	maximum int64
}

func (t responseLimitTransport) RoundTrip(
	request *http.Request,
) (*http.Response, error) {
	response, err := t.base.RoundTrip(request)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		response.Body = &boundedReadCloser{
			source:    response.Body,
			remaining: t.maximum,
		}
	}
	return response, nil
}

type boundedReadCloser struct {
	source    io.ReadCloser
	remaining int64
	exceeded  bool
}

func (r *boundedReadCloser) Read(buffer []byte) (int, error) {
	if r.exceeded {
		return 0, errTokenResponseTooLarge
	}
	limit := r.remaining + 1
	if int64(len(buffer)) > limit {
		buffer = buffer[:int(limit)]
	}
	n, err := r.source.Read(buffer)
	if int64(n) > r.remaining {
		allowed := int(r.remaining)
		r.remaining = 0
		r.exceeded = true
		if allowed > 0 {
			return allowed, errTokenResponseTooLarge
		}
		return 0, errTokenResponseTooLarge
	}
	r.remaining -= int64(n)
	return n, err
}

func (r *boundedReadCloser) Close() error {
	return r.source.Close()
}
