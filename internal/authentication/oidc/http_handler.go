package oidc

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Iron-Signal-Systems/atlas/internal/authentication"
)

const (
	LoginPath                    = "/auth/login"
	CallbackPath                 = "/auth/callback"
	stateCookieName              = "__Host-iron_atlas_oidc_state"
	maximumAuthorizationLocation = 8 << 10
	maximumCallbackQuery         = 16 << 10
	maximumCallbackValue         = 8 << 10
	maximumProviderMessage       = 4 << 10
)

// BrowserAuthorizationFlow is the bounded protocol seam required by the HTTP
// login and callback boundary. AuthorizationCodeFlow implements this contract.
type BrowserAuthorizationFlow interface {
	Begin(context.Context) (AuthorizationRequest, error)
	Complete(context.Context, string, string) (authentication.Principal, error)
	Cancel(context.Context, string) error
	IssuerURL() string
}

// VerifiedPrincipalHandler receives only a principal that completed the
// authorization-code, PKCE, nonce, signature, issuer, and audience checks.
// The authenticated-session candidate uses this seam to create bounded
// server-side session state without exposing provider tokens to the browser.
type VerifiedPrincipalHandler interface {
	ServeVerifiedPrincipal(
		http.ResponseWriter,
		*http.Request,
		authentication.Principal,
	)
}

type VerifiedPrincipalHandlerFunc func(
	http.ResponseWriter,
	*http.Request,
	authentication.Principal,
)

func (f VerifiedPrincipalHandlerFunc) ServeVerifiedPrincipal(
	writer http.ResponseWriter,
	request *http.Request,
	principal authentication.Principal,
) {
	f(writer, request, principal)
}

// HTTPHandlerConfig defines the bounded browser-facing OIDC checkpoint.
type HTTPHandlerConfig struct {
	Flow            BrowserAuthorizationFlow
	VerifiedHandler VerifiedPrincipalHandler
	Now             func() time.Time
}

// HTTPHandler implements only the browser login and callback boundary. It does
// not create an Atlas session, authenticate later requests, enforce CSRF, or
// trust forwarded proxy metadata.
type HTTPHandler struct {
	flow            BrowserAuthorizationFlow
	verifiedHandler VerifiedPrincipalHandler
	issuer          string
	now             func() time.Time
}

func NewHTTPHandler(config HTTPHandlerConfig) (*HTTPHandler, error) {
	if config.Flow == nil {
		return nil, errors.New("browser authorization flow is required")
	}
	if config.VerifiedHandler == nil {
		return nil, errors.New("verified principal handler is required")
	}
	issuer, err := trustedCallbackIssuer(config.Flow.IssuerURL())
	if err != nil {
		return nil, fmt.Errorf("OIDC issuer: %w", err)
	}
	if issuer == "" {
		return nil, errors.New("OIDC issuer is required")
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	return &HTTPHandler{
		flow:            config.Flow,
		verifiedHandler: config.VerifiedHandler,
		issuer:          issuer,
		now:             now,
	}, nil
}

func (h *HTTPHandler) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(LoginPath, h.Login)
	mux.HandleFunc(CallbackPath, h.Callback)
	return mux
}

func (h *HTTPHandler) Login(writer http.ResponseWriter, request *http.Request) {
	browserNoStore(writer)
	if request.Method != http.MethodGet {
		writer.Header().Set("Allow", http.MethodGet)
		writeBrowserFailure(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if request.URL == nil || request.URL.RawQuery != "" {
		clearStateCookie(writer)
		writeBrowserFailure(writer, http.StatusBadRequest, "invalid login request")
		return
	}
	result, err := h.flow.Begin(request.Context())
	if err != nil {
		clearStateCookie(writer)
		writeBrowserFailure(
			writer,
			http.StatusServiceUnavailable,
			"authentication service unavailable",
		)
		return
	}
	if err := validateAuthorizationRequest(result, h.now().UTC()); err != nil {
		clearStateCookie(writer)
		writeBrowserFailure(
			writer,
			http.StatusServiceUnavailable,
			"authentication service unavailable",
		)
		return
	}

	http.SetCookie(writer, stateCookie(result, h.now().UTC()))
	writer.Header().Set("Location", result.AuthorizationURL)
	writer.WriteHeader(http.StatusFound)
}

func (h *HTTPHandler) Callback(writer http.ResponseWriter, request *http.Request) {
	browserNoStore(writer)
	clearStateCookie(writer)
	if request.Method != http.MethodGet {
		writer.Header().Set("Allow", http.MethodGet)
		writeBrowserFailure(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	callback, err := parseCallback(request)
	if err != nil {
		writeBrowserFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}
	cookieState, err := callbackCookieState(request)
	if err != nil || !constantTimeEqual(callback.State, cookieState) {
		writeBrowserFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}
	if callback.Issuer != "" && callback.Issuer != h.issuer {
		if err := h.flow.Cancel(request.Context(), callback.State); err != nil {
			writeAuthenticationError(writer, err)
			return
		}
		writeBrowserFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}

	if callback.ProviderError != "" {
		if err := h.flow.Cancel(request.Context(), callback.State); err != nil {
			writeAuthenticationError(writer, err)
			return
		}
		writeBrowserFailure(writer, http.StatusUnauthorized, "authentication failed")
		return
	}

	principal, err := h.flow.Complete(
		request.Context(),
		callback.State,
		callback.Code,
	)
	if err != nil {
		writeAuthenticationError(writer, err)
		return
	}
	h.verifiedHandler.ServeVerifiedPrincipal(writer, request, principal)
}

type callbackValues struct {
	State         string
	Code          string
	Issuer        string
	ProviderError string
}

func parseCallback(request *http.Request) (callbackValues, error) {
	if request == nil || request.URL == nil {
		return callbackValues{}, errors.New("callback request is required")
	}
	if len(request.URL.RawQuery) == 0 || len(request.URL.RawQuery) > maximumCallbackQuery {
		return callbackValues{}, errors.New("callback query is missing or too large")
	}
	values, err := url.ParseQuery(request.URL.RawQuery)
	if err != nil {
		return callbackValues{}, errors.New("callback query is malformed")
	}
	allowed := map[string]struct{}{
		"state":             {},
		"code":              {},
		"iss":               {},
		"session_state":     {},
		"error":             {},
		"error_description": {},
		"error_uri":         {},
	}
	for key := range values {
		if _, ok := allowed[key]; !ok {
			return callbackValues{}, errors.New("callback query contains an unsupported parameter")
		}
	}

	state, err := exactlyOne(values, "state", maximumCallbackValue, tokenValue)
	if err != nil {
		return callbackValues{}, err
	}
	code, codePresent, err := optionalOne(values, "code", maximumCallbackValue, tokenValue)
	if err != nil {
		return callbackValues{}, err
	}
	issuer, _, err := optionalOne(values, "iss", maximumCallbackValue, textValue)
	if err != nil {
		return callbackValues{}, err
	}
	_, sessionStatePresent, err := optionalOne(
		values,
		"session_state",
		maximumCallbackValue,
		tokenValue,
	)
	if err != nil {
		return callbackValues{}, err
	}
	providerError, errorPresent, err := optionalOne(
		values,
		"error",
		maximumCallbackValue,
		tokenValue,
	)
	if err != nil {
		return callbackValues{}, err
	}
	_, descriptionPresent, err := optionalOne(
		values,
		"error_description",
		maximumProviderMessage,
		textValue,
	)
	if err != nil {
		return callbackValues{}, err
	}
	_, errorURIPresent, err := optionalOne(
		values,
		"error_uri",
		maximumProviderMessage,
		textValue,
	)
	if err != nil {
		return callbackValues{}, err
	}
	if codePresent == errorPresent {
		return callbackValues{}, errors.New("callback must contain exactly one result")
	}
	if codePresent && (descriptionPresent || errorURIPresent) {
		return callbackValues{}, errors.New("provider error metadata requires an error result")
	}
	if errorPresent && sessionStatePresent {
		return callbackValues{}, errors.New("session state is permitted only on successful callbacks")
	}
	return callbackValues{
		State:         state,
		Code:          code,
		Issuer:        issuer,
		ProviderError: providerError,
	}, nil
}

type valueValidator func(string) error

func exactlyOne(
	values url.Values,
	name string,
	maximum int,
	validate valueValidator,
) (string, error) {
	value, present, err := optionalOne(values, name, maximum, validate)
	if err != nil {
		return "", err
	}
	if !present {
		return "", fmt.Errorf("callback parameter %s is required", name)
	}
	return value, nil
}

func optionalOne(
	values url.Values,
	name string,
	maximum int,
	validate valueValidator,
) (string, bool, error) {
	items, present := values[name]
	if !present {
		return "", false, nil
	}
	if len(items) != 1 {
		return "", false, fmt.Errorf("callback parameter %s must not repeat", name)
	}
	value := items[0]
	if len(value) > maximum || validate(value) != nil {
		return "", false, fmt.Errorf("callback parameter %s is invalid", name)
	}
	return value, true, nil
}

func tokenValue(value string) error {
	if value == "" || value != strings.TrimSpace(value) || !utf8.ValidString(value) {
		return errors.New("token value is missing or unnormalized")
	}
	for _, item := range value {
		if unicode.IsSpace(item) || unicode.IsControl(item) {
			return errors.New("token value contains whitespace or control data")
		}
	}
	return nil
}

func textValue(value string) error {
	if value == "" || value != strings.TrimSpace(value) || !utf8.ValidString(value) {
		return errors.New("text value is missing or unnormalized")
	}
	for _, item := range value {
		if unicode.IsControl(item) {
			return errors.New("text value contains control data")
		}
	}
	return nil
}

func callbackCookieState(request *http.Request) (string, error) {
	cookies := request.CookiesNamed(stateCookieName)
	if len(cookies) != 1 {
		return "", errors.New("exactly one state cookie is required")
	}
	value := cookies[0].Value
	if err := validateRandomToken("state cookie", value); err != nil {
		return "", err
	}
	return value, nil
}

func constantTimeEqual(left string, right string) bool {
	if len(left) != len(right) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(left), []byte(right)) == 1
}

func validateAuthorizationRequest(result AuthorizationRequest, now time.Time) error {
	if err := validateRandomToken("state", result.State); err != nil {
		return err
	}
	if result.ExpiresAt.IsZero() || !result.ExpiresAt.After(now) {
		return errors.New("authorization request is already expired")
	}
	if len(result.AuthorizationURL) == 0 || len(result.AuthorizationURL) > maximumAuthorizationLocation {
		return errors.New("authorization location is missing or too large")
	}
	location, err := url.Parse(result.AuthorizationURL)
	if err != nil || location.Scheme != "https" || location.Host == "" || location.User != nil || location.Fragment != "" {
		return errors.New("authorization location is not a trusted HTTPS URL")
	}
	return nil
}

func trustedCallbackIssuer(raw string) (string, error) {
	location, err := url.Parse(raw)
	if err != nil || location.Scheme != "https" || location.Host == "" || location.User != nil || location.Fragment != "" || location.RawQuery != "" {
		return "", errors.New("issuer must be an exact HTTPS origin or path")
	}
	return location.String(), nil
}

func stateCookie(result AuthorizationRequest, now time.Time) *http.Cookie {
	remaining := result.ExpiresAt.Sub(now)
	if remaining <= 0 {
		remaining = time.Second
	}
	maxAge := int(remaining / time.Second)
	if remaining%time.Second != 0 {
		maxAge++
	}
	if maxAge < 1 {
		maxAge = 1
	}
	return &http.Cookie{
		Name:     stateCookieName,
		Value:    result.State,
		Path:     "/",
		Expires:  result.ExpiresAt.UTC(),
		MaxAge:   maxAge,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func clearStateCookie(writer http.ResponseWriter) {
	http.SetCookie(writer, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(1, 0).UTC(),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func browserNoStore(writer http.ResponseWriter) {
	writer.Header().Set("Cache-Control", "no-store")
	writer.Header().Set("Pragma", "no-cache")
	writer.Header().Set("Referrer-Policy", "no-referrer")
	writer.Header().Set("X-Content-Type-Options", "nosniff")
}

func writeAuthenticationError(writer http.ResponseWriter, err error) {
	if errors.Is(err, authentication.ErrAuthenticationUnavailable) {
		writeBrowserFailure(
			writer,
			http.StatusServiceUnavailable,
			"authentication service unavailable",
		)
		return
	}
	writeBrowserFailure(writer, http.StatusUnauthorized, "authentication failed")
}

func writeBrowserFailure(writer http.ResponseWriter, status int, message string) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	_, _ = writer.Write([]byte(message + "\n"))
}
