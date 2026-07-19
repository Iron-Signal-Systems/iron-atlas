#!/usr/bin/env python3
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
passes = 0
failures: list[str] = []


def compact(value: str) -> str:
    return " ".join(value.split())


def check(name: str, condition: bool) -> None:
    global passes
    if condition:
        passes += 1
        print(f"PASS: {name}")
    else:
        failures.append(name)
        print(f"FAIL: {name}")


def require(path: str, tokens: list[str]) -> str:
    candidate = ROOT / path
    check(f"required file {path}", candidate.is_file())
    if not candidate.is_file():
        return ""
    text = candidate.read_text(encoding="utf-8")
    normalized = compact(text)
    for token in tokens:
        check(f"{path} contains {token}", compact(token) in normalized)
    return text


flow = require(
    "internal/authentication/oidc/authorization_code.go",
    [
        "func (f *AuthorizationCodeFlow) Cancel",
        "f.store.Consume",
        "func (f *AuthorizationCodeFlow) IssuerURL",
        "return f.verifier.issuerURL",
    ],
)
check(
    "authorization flow cancellation does not log protocol material",
    "slog." not in flow and "log." not in flow,
)

handler = require(
    "internal/authentication/oidc/http_handler.go",
    [
        'LoginPath = "/auth/login"',
        'CallbackPath = "/auth/callback"',
        "__Host-iron_atlas_oidc_state",
        "type BrowserAuthorizationFlow interface",
        "type VerifiedPrincipalHandler interface",
        "http.SameSiteLaxMode",
        "request.CookiesNamed",
        "subtle.ConstantTimeCompare",
        "url.ParseQuery",
        "callback query contains an unsupported parameter",
        "h.flow.Cancel",
        "h.flow.Complete",
        "ServeVerifiedPrincipal",
        "Cache-Control",
        "no-store",
        "Referrer-Policy",
        "no-referrer",
    ],
)
check(
    "login handler uses an empty-body redirect rather than http.Redirect",
    "http.Redirect" not in handler,
)
check(
    "handler does not accept request-selected return locations",
    "return_to" not in handler and "redirect_uri" not in handler,
)
check(
    "handler does not log callback or cookie material",
    "slog." not in handler and "log." not in handler,
)
check(
    "handler does not create an Atlas session",
    "SessionID" not in handler
    and "session.Store" not in handler
    and "CreateSession" not in handler,
)

tests = require(
    "internal/authentication/oidc/http_handler_test.go",
    [
        "TestHTTPLoginCreatesBoundSecureStateCookie",
        "TestHTTPCallbackProducesVerifiedPrincipal",
        "TestHTTPCallbackRejectsStateMismatchBeforeExchange",
        "TestHTTPCallbackRejectsDuplicateAndUnsupportedParameters",
        "TestHTTPCallbackConsumesProviderErrorWithoutReflectingDetails",
        "TestHTTPCallbackAllowsOnlyOneConcurrentConsumer",
        "TestHTTPCallbackBoundsAndRedactsInput",
        "TestHTTPRoutesRejectUnsafeMethods",
        "TestNewHTTPHandlerRejectsUnsafeConfiguration",
        "TestHTTPCallbackClassifiesProviderOutage",
    ],
)
check(
    "HTTP tests do not depend on a live provider",
    "accounts.google.com" not in tests
    and "login.microsoftonline.com" not in tests,
)

require(
    "internal/authentication/oidc/authorization_code_cancel_test.go",
    [
        "TestAuthorizationCodeFlowCancelConsumesStateExactlyOnce",
        "authentication.ErrAuthenticationInvalid",
    ],
)
require(
    "docs/architecture/OIDC-HTTP-LOGIN-AND-CALLBACK-IMPLEMENTATION.md",
    [
        "bounded implementation candidate",
        "__Host-iron_atlas_oidc_state",
        "constant time",
        "atomically consumes error callbacks",
        "VerifiedPrincipalHandler",
        "deliberately not wired",
        "does not establish",
    ],
)
require(
    "docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md",
    ["HTTP login and callback implementation checkpoint", "verified-principal handoff"],
)
require(
    "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md",
    ["HTTP login and callback implementation status", "browser state-binding cookie"],
)
require(
    "docs/testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md",
    ["HTTP login and callback implementation campaign", "provider-error cancellation"],
)
require(
    "test-framework/authentication/test_phase1_step3_http_login_callback.sh",
    ["HTTP login and callback static validator", "go test -race", "go mod verify", "govulncheck"],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_http_login_callback.sh",
    [
        "28ec1eab5b5c4e69731e9b0a79fe6105beab316d",
        "implementation candidate only",
        "does not establish durable sessions",
    ],
)

app = require("internal/app/app.go", [])
server = require("internal/httpui/server.go", [])
check(
    "HTTP checkpoint is deliberately not wired into atlasd",
    "NewHTTPHandler" not in app
    and '"/auth/login"' not in server
    and '"/auth/callback"' not in server,
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text(encoding="utf-8").splitlines()
for path in (
    "internal/authentication/oidc/http_handler.go",
    "internal/authentication/oidc/http_handler_test.go",
    "internal/authentication/oidc/authorization_code_cancel_test.go",
    "docs/architecture/OIDC-HTTP-LOGIN-AND-CALLBACK-IMPLEMENTATION.md",
    "tools/validation/validate_phase1_step3_http_login_callback.py",
    "test-framework/authentication/test_phase1_step3_http_login_callback.sh",
    "tools/validation/phase-gates/validate_phase1_step3_http_login_callback.sh",
):
    check(f"file manifest contains {path}", path in manifest)

check(
    "formal Step 3 acceptance record remains absent",
    not (ROOT / "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md").exists(),
)

print()
print(f"PASS checks: {passes}")
print(f"FAIL checks: {len(failures)}")
if failures:
    print("Phase 1 Step 3 HTTP login and callback validation FAILED.")
    raise SystemExit(1)

print("Phase 1 Step 3 HTTP login and callback validation PASSED.")
print(
    "This proves bounded browser login and callback routes, secure state-cookie "
    "binding, exact callback parsing, issuer checks, provider-error cancellation, "
    "generic failure handling, one-consumer behavior, and verified-principal handoff."
)
print(
    "It does not prove durable sessions, protected-route authentication, CSRF, "
    "logout, trusted proxies, production wiring, production credential delivery, "
    "representative-provider compatibility, formal Step 3 acceptance, or "
    "production readiness."
)
