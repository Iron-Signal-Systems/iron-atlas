package oidc

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	coreoidc "github.com/coreos/go-oidc/v3/oidc"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
)

const (
	defaultHTTPTimeout  = 5 * time.Second
	defaultClockSkew    = time.Minute
	defaultMaxTokenAge  = 15 * time.Minute
	defaultMaxTokenSize = 64 << 10
	maxHeaderSize       = 8 << 10
	maxPayloadSize      = 56 << 10
)

var permittedAlgorithms = map[string]struct{}{
	coreoidc.RS256: {},
	coreoidc.PS256: {},
	coreoidc.ES256: {},
	coreoidc.EdDSA: {},
}

// Config defines one trusted OIDC provider-verification profile.
//
// This component verifies provider discovery metadata and ID tokens. It does
// not implement browser redirects, authorization-code exchange, cookies,
// sessions, CSRF, logout, or trusted-proxy handling.
type Config struct {
	ProviderID               string
	IssuerURL                string
	ClientID                 string
	AllowedSigningAlgorithms []string
	HTTPClient               *http.Client
	Now                      func() time.Time
	MaxClockSkew             time.Duration
	MaxTokenAge              time.Duration
	MaxTokenBytes            int
}

type keyVerificationState struct {
	dependencyUnavailable bool
}

type keyVerificationStateKey struct{}

type classifyingKeySet struct {
	delegate coreoidc.KeySet
}

func (k classifyingKeySet) VerifySignature(
	ctx context.Context,
	rawJWT string,
) ([]byte, error) {
	payload, err := k.delegate.VerifySignature(ctx, rawJWT)
	if err == nil {
		return payload, nil
	}

	// coreos/go-oidc v3.19.0 wraps the KeySet error with %v after this
	// method returns, which intentionally does not preserve the error chain.
	// Record dependency classification in request-local state before that
	// wrapping occurs. A completed remote refresh that simply cannot verify
	// the signature remains an invalid assertion, not an outage.
	if keySetDependencyUnavailable(err) {
		if state, ok := ctx.Value(
			keyVerificationStateKey{},
		).(*keyVerificationState); ok {
			state.dependencyUnavailable = true
		}
	}
	return nil, err
}

func keySetDependencyUnavailable(err error) bool {
	if unavailable(err) {
		return true
	}
	return err.Error() != "failed to verify id token signature"
}

// Verifier validates one OIDC provider's signed ID tokens and normalizes the
// verified result into the provider-neutral authentication.Principal contract.
type Verifier struct {
	providerID               string
	clientID                 string
	issuerURL                string
	authorizationEndpoint    string
	tokenEndpoint            string
	tokenEndpointAuthMethods []string
	codeChallengeMethods     []string
	allowedAlgs              map[string]struct{}
	httpClient               *http.Client
	tokenVerifier            *coreoidc.IDTokenVerifier
	now                      func() time.Time
	maxClockSkew             time.Duration
	maxTokenAge              time.Duration
	maxTokenBytes            int
}

func New(ctx context.Context, config Config) (*Verifier, error) {
	if ctx == nil {
		return nil, errors.New("context is required")
	}
	providerID, err := boundedIdentifier("provider ID", config.ProviderID, 256)
	if err != nil {
		return nil, err
	}
	issuerURL, err := trustedHTTPSURL("issuer URL", config.IssuerURL)
	if err != nil {
		return nil, err
	}
	clientID, err := boundedIdentifier("client ID", config.ClientID, 512)
	if err != nil {
		return nil, err
	}

	allowed, ordered, err := signingAlgorithms(config.AllowedSigningAlgorithms)
	if err != nil {
		return nil, err
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	clockSkew := config.MaxClockSkew
	if clockSkew == 0 {
		clockSkew = defaultClockSkew
	}
	if clockSkew < 0 || clockSkew > 5*time.Minute {
		return nil, errors.New("maximum clock skew must be between zero and five minutes")
	}
	maxTokenAge := config.MaxTokenAge
	if maxTokenAge == 0 {
		maxTokenAge = defaultMaxTokenAge
	}
	if maxTokenAge < time.Minute || maxTokenAge > time.Hour {
		return nil, errors.New("maximum token age must be between one minute and one hour")
	}
	maxTokenBytes := config.MaxTokenBytes
	if maxTokenBytes == 0 {
		maxTokenBytes = defaultMaxTokenSize
	}
	if maxTokenBytes < 1024 || maxTokenBytes > 1<<20 {
		return nil, errors.New("maximum token size must be between 1024 bytes and 1 MiB")
	}

	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultHTTPTimeout}
	} else {
		clone := *httpClient
		httpClient = &clone
		if httpClient.Timeout <= 0 || httpClient.Timeout > 30*time.Second {
			return nil, errors.New(
				"OIDC HTTP client timeout must be greater than zero and no more than 30 seconds",
			)
		}
	}

	providerContext := coreoidc.ClientContext(ctx, httpClient)
	provider, err := coreoidc.NewProvider(providerContext, issuerURL)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: OIDC provider discovery failed",
			authentication.ErrAuthenticationUnavailable,
		)
	}

	endpoint := provider.Endpoint()
	authorizationEndpoint, err := trustedHTTPSURL(
		"authorization endpoint",
		endpoint.AuthURL,
	)
	if err != nil {
		return nil, err
	}
	tokenEndpoint, err := trustedHTTPSURL("token endpoint", endpoint.TokenURL)
	if err != nil {
		return nil, err
	}
	var metadata struct {
		JWKSURL                  string   `json:"jwks_uri"`
		SigningAlgorithms        []string `json:"id_token_signing_alg_values_supported"`
		ResponseTypesSupported   []string `json:"response_types_supported"`
		SubjectTypesSupported    []string `json:"subject_types_supported"`
		TokenEndpointAuthMethods []string `json:"token_endpoint_auth_methods_supported"`
		CodeChallengeMethods     []string `json:"code_challenge_methods_supported"`
	}
	if err := provider.Claims(&metadata); err != nil {
		return nil, fmt.Errorf(
			"%w: OIDC provider metadata is invalid",
			authentication.ErrAuthenticationInvalid,
		)
	}
	if _, err := trustedHTTPSURL("JWKS endpoint", metadata.JWKSURL); err != nil {
		return nil, err
	}
	if !intersects(metadata.SigningAlgorithms, allowed) {
		return nil, fmt.Errorf(
			"%w: provider advertises no permitted ID-token signing algorithm",
			authentication.ErrAuthenticationInvalid,
		)
	}
	if !contains(metadata.ResponseTypesSupported, "code") {
		return nil, fmt.Errorf(
			"%w: provider does not advertise authorization-code response support",
			authentication.ErrAuthenticationInvalid,
		)
	}
	if len(metadata.SubjectTypesSupported) == 0 {
		return nil, fmt.Errorf(
			"%w: provider does not advertise a subject type",
			authentication.ErrAuthenticationInvalid,
		)
	}

	remoteKeySet := coreoidc.NewRemoteKeySet(
		providerContext,
		metadata.JWKSURL,
	)
	tokenVerifier := coreoidc.NewVerifier(
		issuerURL,
		classifyingKeySet{delegate: remoteKeySet},
		&coreoidc.Config{
			ClientID:             clientID,
			SupportedSigningAlgs: ordered,
			Now: func() time.Time {
				return now().UTC().Add(-clockSkew)
			},
		},
	)

	return &Verifier{
		providerID:               providerID,
		clientID:                 clientID,
		issuerURL:                issuerURL,
		authorizationEndpoint:    authorizationEndpoint,
		tokenEndpoint:            tokenEndpoint,
		tokenEndpointAuthMethods: append([]string(nil), metadata.TokenEndpointAuthMethods...),
		codeChallengeMethods:     append([]string(nil), metadata.CodeChallengeMethods...),
		allowedAlgs:              allowed,
		httpClient:               httpClient,
		tokenVerifier:            tokenVerifier,
		now:                      now,
		maxClockSkew:             clockSkew,
		maxTokenAge:              maxTokenAge,
		maxTokenBytes:            maxTokenBytes,
	}, nil
}

// Verify validates a signed ID token for the configured provider.
//
// expectedNonce is mandatory because this verifier is intended for an
// authorization-code flow. accessToken may be empty when the ID token does not
// contain at_hash. If at_hash is present, the access token is required and its
// hash is verified.
func (v *Verifier) Verify(
	ctx context.Context,
	rawIDToken string,
	expectedNonce string,
	accessToken string,
) (authentication.Principal, error) {
	if v == nil || v.tokenVerifier == nil || v.httpClient == nil {
		return authentication.Principal{}, fmt.Errorf(
			"%w: OIDC verifier is unavailable",
			authentication.ErrAuthenticationUnavailable,
		)
	}
	if ctx == nil {
		return authentication.Principal{}, fmt.Errorf(
			"%w: context is required",
			authentication.ErrAuthenticationInvalid,
		)
	}
	if rawIDToken == "" ||
		rawIDToken != strings.TrimSpace(rawIDToken) ||
		len(rawIDToken) > v.maxTokenBytes ||
		!utf8.ValidString(rawIDToken) {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	nonce, err := boundedIdentifier("expected nonce", expectedNonce, 512)
	if err != nil || len(nonce) < 16 {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	if err := inspectJWT(rawIDToken, v.allowedAlgs); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	verificationState := &keyVerificationState{}
	verifyContext := context.WithValue(
		coreoidc.ClientContext(ctx, v.httpClient),
		keyVerificationStateKey{},
		verificationState,
	)
	token, err := v.tokenVerifier.Verify(verifyContext, rawIDToken)
	if err != nil {
		if verificationState.dependencyUnavailable || unavailable(err) {
			return authentication.Principal{}, fmt.Errorf(
				"%w: OIDC key verification unavailable",
				authentication.ErrAuthenticationUnavailable,
			)
		}
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	var claims struct {
		AuthorizedParty    string          `json:"azp"`
		NotBefore          json.RawMessage `json:"nbf"`
		AuthenticationTime json.RawMessage `json:"auth_time"`
	}
	if err := token.Claims(&claims); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	now := v.now().UTC()
	if token.Issuer != v.issuerURL ||
		token.Subject == "" ||
		token.Subject != strings.TrimSpace(token.Subject) ||
		token.IssuedAt.IsZero() ||
		token.IssuedAt.After(now.Add(v.maxClockSkew)) ||
		now.Sub(token.IssuedAt) > v.maxTokenAge+v.maxClockSkew {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	if subtle.ConstantTimeCompare(
		[]byte(token.Nonce),
		[]byte(nonce),
	) != 1 {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	if len(token.Audience) > 1 && claims.AuthorizedParty == "" {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	if claims.AuthorizedParty != "" &&
		claims.AuthorizedParty != v.clientID {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	notBefore, present, err := numericDate(claims.NotBefore)
	if err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	if present && notBefore.After(now.Add(v.maxClockSkew)) {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}

	authenticatedAt := token.IssuedAt
	authTime, present, err := numericDate(claims.AuthenticationTime)
	if err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	if present {
		if authTime.After(now.Add(v.maxClockSkew)) ||
			authTime.After(token.IssuedAt.Add(v.maxClockSkew)) {
			return authentication.Principal{}, authentication.ErrAuthenticationInvalid
		}
		authenticatedAt = authTime
	}

	if token.AccessTokenHash != "" {
		if accessToken == "" ||
			len(accessToken) > v.maxTokenBytes ||
			token.VerifyAccessToken(accessToken) != nil {
			return authentication.Principal{}, authentication.ErrAuthenticationInvalid
		}
	}

	principal := authentication.Principal{
		ProviderID:      v.providerID,
		Subject:         token.Subject,
		AuthenticatedAt: authenticatedAt.UTC(),
	}
	if err := principal.Validate(); err != nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	return principal, nil
}

func signingAlgorithms(values []string) (
	map[string]struct{},
	[]string,
	error,
) {
	if len(values) == 0 {
		values = []string{coreoidc.RS256}
	}
	allowed := make(map[string]struct{}, len(values))
	ordered := make([]string, 0, len(values))
	for _, raw := range values {
		value := strings.TrimSpace(raw)
		if value == "" || value != raw {
			return nil, nil, errors.New("signing algorithms must be nonempty and normalized")
		}
		if _, ok := permittedAlgorithms[value]; !ok {
			return nil, nil, fmt.Errorf("unsupported or symmetric signing algorithm %q", value)
		}
		if _, duplicate := allowed[value]; duplicate {
			return nil, nil, fmt.Errorf("duplicate signing algorithm %q", value)
		}
		allowed[value] = struct{}{}
		ordered = append(ordered, value)
	}
	return allowed, ordered, nil
}

func boundedIdentifier(name string, value string, maximum int) (string, error) {
	if value == "" ||
		value != strings.TrimSpace(value) ||
		len(value) > maximum ||
		!utf8.ValidString(value) {
		return "", fmt.Errorf("%s is missing, unnormalized, invalid, or too large", name)
	}
	for _, r := range value {
		if r < 0x20 || r == 0x7f {
			return "", fmt.Errorf("%s contains a control character", name)
		}
	}
	return value, nil
}

func trustedHTTPSURL(name string, raw string) (string, error) {
	if raw == "" || raw != strings.TrimSpace(raw) || len(raw) > 2048 {
		return "", fmt.Errorf("%s is missing, unnormalized, or too large", name)
	}
	parsed, err := url.Parse(raw)
	if err != nil ||
		parsed.Scheme != "https" ||
		parsed.Host == "" ||
		parsed.User != nil ||
		parsed.RawQuery != "" ||
		parsed.Fragment != "" {
		return "", fmt.Errorf("%s must be an exact HTTPS URL without userinfo, query, or fragment", name)
	}
	return raw, nil
}

func intersects(values []string, allowed map[string]struct{}) bool {
	for _, value := range values {
		if _, ok := allowed[value]; ok {
			return true
		}
	}
	return false
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func unavailable(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var networkError net.Error
	return errors.As(err, &networkError)
}

func numericDate(raw json.RawMessage) (time.Time, bool, error) {
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return time.Time{}, false, nil
	}
	value, err := strconv.ParseFloat(string(raw), 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return time.Time{}, false, errors.New("invalid NumericDate")
	}
	seconds, fraction := math.Modf(value)
	nanoseconds := int64(fraction * float64(time.Second))
	return time.Unix(int64(seconds), nanoseconds).UTC(), true, nil
}

func inspectJWT(raw string, allowed map[string]struct{}) error {
	parts := strings.Split(raw, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return errors.New("JWT must contain exactly three nonempty segments")
	}
	header, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil || len(header) > maxHeaderSize {
		return errors.New("JWT header is invalid or oversized")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || len(payload) > maxPayloadSize {
		return errors.New("JWT payload is invalid or oversized")
	}

	if err := rejectDuplicateFields(
		header,
		map[string]struct{}{"alg": {}, "kid": {}, "typ": {}},
	); err != nil {
		return err
	}
	if err := rejectDuplicateFields(
		payload,
		map[string]struct{}{
			"iss": {}, "sub": {}, "aud": {}, "exp": {}, "iat": {},
			"nbf": {}, "nonce": {}, "azp": {}, "auth_time": {},
			"at_hash": {},
		},
	); err != nil {
		return err
	}

	var protected struct {
		Algorithm string `json:"alg"`
	}
	if err := json.Unmarshal(header, &protected); err != nil {
		return err
	}
	if protected.Algorithm == "" ||
		strings.EqualFold(protected.Algorithm, "none") {
		return errors.New("JWT algorithm is missing or prohibited")
	}
	if _, ok := allowed[protected.Algorithm]; !ok {
		return errors.New("JWT algorithm is not permitted")
	}
	return nil
}

func rejectDuplicateFields(
	data []byte,
	sensitive map[string]struct{},
) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	first, err := decoder.Token()
	if err != nil || first != json.Delim('{') {
		return errors.New("security object must be a JSON object")
	}
	seen := make(map[string]struct{})
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		key, ok := token.(string)
		if !ok {
			return errors.New("security object key is not a string")
		}
		if _, relevant := sensitive[key]; relevant {
			if _, duplicate := seen[key]; duplicate {
				return fmt.Errorf("duplicate security-sensitive field %q", key)
			}
			seen[key] = struct{}{}
		}
		var discard json.RawMessage
		if err := decoder.Decode(&discard); err != nil {
			return err
		}
	}
	last, err := decoder.Token()
	if err != nil || last != json.Delim('}') {
		return errors.New("security object is malformed")
	}
	if token, err := decoder.Token(); err != io.EOF || token != nil {
		return errors.New("security object contains trailing JSON")
	}
	return nil
}
