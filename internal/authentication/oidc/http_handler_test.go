package oidc

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Iron-Signal-Systems/iron-atlas/internal/authentication"
)

const testIssuer = "https://identity.example.test/oidc"

type fakeBrowserAuthorizationFlow struct {
	beginRequest AuthorizationRequest
	beginErr     error
	issuer       string
	complete     func(context.Context, string, string) (authentication.Principal, error)
	cancel       func(context.Context, string) error
}

func (f *fakeBrowserAuthorizationFlow) Begin(context.Context) (AuthorizationRequest, error) {
	return f.beginRequest, f.beginErr
}

func (f *fakeBrowserAuthorizationFlow) Complete(
	ctx context.Context,
	state string,
	code string,
) (authentication.Principal, error) {
	if f.complete == nil {
		return authentication.Principal{}, authentication.ErrAuthenticationInvalid
	}
	return f.complete(ctx, state, code)
}

func (f *fakeBrowserAuthorizationFlow) Cancel(ctx context.Context, state string) error {
	if f.cancel == nil {
		return nil
	}
	return f.cancel(ctx, state)
}

func (f *fakeBrowserAuthorizationFlow) IssuerURL() string {
	if f.issuer == "" {
		return testIssuer
	}
	return f.issuer
}

type recordingPrincipalHandler struct {
	mu         sync.Mutex
	principals []authentication.Principal
	status     int
}

func (h *recordingPrincipalHandler) ServeVerifiedPrincipal(
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

func (h *recordingPrincipalHandler) count() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.principals)
}

func TestHTTPLoginCreatesBoundSecureStateCookie(t *testing.T) {
	now := time.Date(2026, 7, 19, 19, 0, 0, 0, time.UTC)
	state := testState(1)
	location := "https://identity.example.test/authorize?state=" + url.QueryEscape(state)
	flow := &fakeBrowserAuthorizationFlow{beginRequest: AuthorizationRequest{
		AuthorizationURL: location,
		State:            state,
		ExpiresAt:        now.Add(5 * time.Minute),
	}}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, now)

	request := httptest.NewRequest(http.MethodGet, LoginPath, nil)
	response := httptest.NewRecorder()
	handler.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusFound)
	}
	if got := response.Header().Get("Location"); got != location {
		t.Fatalf("location = %q, want %q", got, location)
	}
	if strings.Contains(response.Body.String(), state) {
		t.Fatal("response body exposed OIDC state")
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie count = %d, want 1", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != stateCookieName || cookie.Value != state {
		t.Fatalf("unexpected state cookie: %#v", cookie)
	}
	if cookie.Path != "/" || !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("unsafe cookie attributes: %#v", cookie)
	}
	if cookie.Domain != "" {
		t.Fatalf("cookie domain = %q, want host-only", cookie.Domain)
	}
	if cookie.MaxAge != 300 {
		t.Fatalf("cookie max age = %d, want 300", cookie.MaxAge)
	}
	assertBrowserNoStore(t, response.Header())
}

func TestHTTPCallbackProducesVerifiedPrincipal(t *testing.T) {
	now := time.Date(2026, 7, 19, 19, 0, 0, 0, time.UTC)
	state := testState(2)
	principal := authentication.Principal{
		ProviderID:      "provider-1",
		Subject:         "subject-1",
		AuthenticatedAt: now,
	}
	flow := &fakeBrowserAuthorizationFlow{
		complete: func(_ context.Context, gotState string, code string) (authentication.Principal, error) {
			if gotState != state || code != "authorization-code" {
				t.Fatalf("complete got state %q and code %q", gotState, code)
			}
			return principal, nil
		},
	}
	receiver := &recordingPrincipalHandler{}
	handler := newHTTPHandlerForTest(t, flow, receiver, now)

	request := callbackRequest(state, "authorization-code")
	response := httptest.NewRecorder()
	handler.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if receiver.count() != 1 {
		t.Fatalf("verified-principal count = %d, want 1", receiver.count())
	}
	assertClearedStateCookie(t, response)
	assertBrowserNoStore(t, response.Header())
}

func TestHTTPCallbackRejectsStateMismatchBeforeExchange(t *testing.T) {
	var completeCalls atomic.Int32
	flow := &fakeBrowserAuthorizationFlow{
		complete: func(context.Context, string, string) (authentication.Principal, error) {
			completeCalls.Add(1)
			return authentication.Principal{}, nil
		},
	}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	request := callbackRequest(testState(3), "authorization-code")
	request.Header.Set("Cookie", stateCookieName+"="+testState(4))
	response := httptest.NewRecorder()

	handler.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	if completeCalls.Load() != 0 {
		t.Fatalf("complete calls = %d, want 0", completeCalls.Load())
	}
	assertClearedStateCookie(t, response)
}

func TestHTTPCallbackRejectsDuplicateAndUnsupportedParameters(t *testing.T) {
	state := testState(5)
	tests := []string{
		CallbackPath + "?state=" + state + "&state=" + state + "&code=code",
		CallbackPath + "?state=" + state + "&code=one&code=two",
		CallbackPath + "?state=" + state + "&code=code&return_to=https://attacker.test/",
		CallbackPath + "?state=" + state + "&code=code&error=access_denied",
	}
	flow := &fakeBrowserAuthorizationFlow{}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	for _, target := range tests {
		t.Run(target, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, target, nil)
			request.AddCookie(&http.Cookie{Name: stateCookieName, Value: state})
			response := httptest.NewRecorder()
			handler.Handler().ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
			}
		})
	}
}

func TestHTTPCallbackConsumesProviderErrorWithoutReflectingDetails(t *testing.T) {
	state := testState(6)
	marker := "private-provider-description"
	var canceled atomic.Int32
	flow := &fakeBrowserAuthorizationFlow{
		cancel: func(_ context.Context, gotState string) error {
			if gotState != state {
				t.Fatalf("cancel state = %q, want %q", gotState, state)
			}
			canceled.Add(1)
			return nil
		},
	}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	target := CallbackPath + "?state=" + url.QueryEscape(state) +
		"&error=access_denied&error_description=" + url.QueryEscape(marker)
	request := httptest.NewRequest(http.MethodGet, target, nil)
	request.AddCookie(&http.Cookie{Name: stateCookieName, Value: state})
	response := httptest.NewRecorder()

	handler.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	if canceled.Load() != 1 {
		t.Fatalf("cancel calls = %d, want 1", canceled.Load())
	}
	if strings.Contains(response.Body.String(), marker) {
		t.Fatal("provider error description was reflected")
	}
}

func TestHTTPCallbackAllowsOnlyOneConcurrentConsumer(t *testing.T) {
	state := testState(7)
	principal := authentication.Principal{
		ProviderID:      "provider-1",
		Subject:         "subject-1",
		AuthenticatedAt: time.Now().UTC(),
	}
	var consumed atomic.Bool
	flow := &fakeBrowserAuthorizationFlow{
		complete: func(context.Context, string, string) (authentication.Principal, error) {
			if !consumed.CompareAndSwap(false, true) {
				return authentication.Principal{}, authentication.ErrAuthenticationInvalid
			}
			return principal, nil
		},
	}
	receiver := &recordingPrincipalHandler{}
	handler := newHTTPHandlerForTest(t, flow, receiver, time.Now().UTC())

	const attempts = 32
	statuses := make(chan int, attempts)
	var group sync.WaitGroup
	group.Add(attempts)
	for range attempts {
		go func() {
			defer group.Done()
			response := httptest.NewRecorder()
			handler.Handler().ServeHTTP(response, callbackRequest(state, "authorization-code"))
			statuses <- response.Code
		}()
	}
	group.Wait()
	close(statuses)

	successes := 0
	for status := range statuses {
		if status == http.StatusNoContent {
			successes++
			continue
		}
		if status != http.StatusUnauthorized {
			t.Fatalf("unexpected callback status %d", status)
		}
	}
	if successes != 1 || receiver.count() != 1 {
		t.Fatalf("successes = %d, verified principals = %d, want 1 and 1", successes, receiver.count())
	}
}

func TestHTTPCallbackBoundsAndRedactsInput(t *testing.T) {
	state := testState(8)
	marker := strings.Repeat("sensitive-code-", 2000)
	flow := &fakeBrowserAuthorizationFlow{}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	request := httptest.NewRequest(
		http.MethodGet,
		CallbackPath+"?state="+url.QueryEscape(state)+"&code="+url.QueryEscape(marker),
		nil,
	)
	request.AddCookie(&http.Cookie{Name: stateCookieName, Value: state})
	response := httptest.NewRecorder()

	handler.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
	if strings.Contains(response.Body.String(), "sensitive-code") {
		t.Fatal("oversized authorization code was reflected")
	}
}

func TestHTTPRoutesRejectUnsafeMethods(t *testing.T) {
	flow := &fakeBrowserAuthorizationFlow{}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	for _, path := range []string{LoginPath, CallbackPath} {
		request := httptest.NewRequest(http.MethodPost, path, strings.NewReader("state=bad"))
		response := httptest.NewRecorder()
		handler.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s status = %d, want %d", path, response.Code, http.StatusMethodNotAllowed)
		}
	}
}

func TestNewHTTPHandlerRejectsUnsafeConfiguration(t *testing.T) {
	handler := &recordingPrincipalHandler{}
	if _, err := NewHTTPHandler(HTTPHandlerConfig{VerifiedHandler: handler}); err == nil {
		t.Fatal("expected missing-flow error")
	}
	if _, err := NewHTTPHandler(HTTPHandlerConfig{Flow: &fakeBrowserAuthorizationFlow{}}); err == nil {
		t.Fatal("expected missing verified-handler error")
	}
	if _, err := NewHTTPHandler(HTTPHandlerConfig{
		Flow:            &fakeBrowserAuthorizationFlow{issuer: "http://identity.example.test"},
		VerifiedHandler: handler,
	}); err == nil {
		t.Fatal("expected insecure issuer error")
	}
}

func TestHTTPLoginClassifiesUnavailableWithoutLeakingError(t *testing.T) {
	marker := "private-randomness-failure"
	flow := &fakeBrowserAuthorizationFlow{beginErr: fmt.Errorf("%s: %w", marker, authentication.ErrAuthenticationUnavailable)}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	response := httptest.NewRecorder()
	handler.Handler().ServeHTTP(response, httptest.NewRequest(http.MethodGet, LoginPath, nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	if strings.Contains(response.Body.String(), marker) {
		t.Fatal("internal login error was reflected")
	}
}

func TestHTTPCallbackClassifiesProviderOutage(t *testing.T) {
	state := testState(9)
	flow := &fakeBrowserAuthorizationFlow{
		complete: func(context.Context, string, string) (authentication.Principal, error) {
			return authentication.Principal{}, errors.Join(
				authentication.ErrAuthenticationUnavailable,
				errors.New("private provider outage"),
			)
		},
	}
	handler := newHTTPHandlerForTest(t, flow, &recordingPrincipalHandler{}, time.Now().UTC())
	response := httptest.NewRecorder()
	handler.Handler().ServeHTTP(response, callbackRequest(state, "authorization-code"))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
	if strings.Contains(response.Body.String(), "private provider outage") {
		t.Fatal("provider outage detail was reflected")
	}
}

func newHTTPHandlerForTest(
	t *testing.T,
	flow BrowserAuthorizationFlow,
	receiver VerifiedPrincipalHandler,
	now time.Time,
) *HTTPHandler {
	t.Helper()
	handler, err := NewHTTPHandler(HTTPHandlerConfig{
		Flow:            flow,
		VerifiedHandler: receiver,
		Now:             func() time.Time { return now },
	})
	if err != nil {
		t.Fatal(err)
	}
	return handler
}

func callbackRequest(state string, code string) *http.Request {
	target := CallbackPath + "?state=" + url.QueryEscape(state) + "&code=" + url.QueryEscape(code)
	request := httptest.NewRequest(http.MethodGet, target, nil)
	request.AddCookie(&http.Cookie{Name: stateCookieName, Value: state})
	return request
}

func testState(fill byte) string {
	return base64.RawURLEncoding.EncodeToString(bytes.Repeat([]byte{fill}, 32))
}

func assertBrowserNoStore(t *testing.T, header http.Header) {
	t.Helper()
	if header.Get("Cache-Control") != "no-store" || header.Get("Pragma") != "no-cache" {
		t.Fatalf("missing no-store headers: %#v", header)
	}
	if header.Get("Referrer-Policy") != "no-referrer" {
		t.Fatalf("referrer policy = %q", header.Get("Referrer-Policy"))
	}
}

func assertClearedStateCookie(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	cookies := response.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cleared cookie count = %d, want 1", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != stateCookieName || cookie.MaxAge >= 0 || cookie.Value != "" {
		t.Fatalf("state cookie was not cleared: %#v", cookie)
	}
	if cookie.Path != "/" || !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cleared cookie lost security attributes: %#v", cookie)
	}
}
