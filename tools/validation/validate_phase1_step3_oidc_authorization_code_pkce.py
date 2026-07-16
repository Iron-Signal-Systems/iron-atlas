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
        check(
            f"{path} contains {token}",
            compact(token) in normalized,
        )
    return text


flow = require(
    "internal/authentication/oidc/authorization_code.go",
    [
        "type PreauthenticationStore interface",
        "type MemoryPreauthenticationStore struct",
        "sha256.Sum256([]byte(state))",
        "oauth2.S256ChallengeOption(pkceVerifier)",
        "oauth2.VerifierOption(transaction.PKCEVerifier)",
        "coreoidc.Nonce(nonce)",
        "provider does not advertise PKCE S256 support",
        "responseLimitTransport",
        "http.ErrUseLastResponse",
        "tokenEndpointAuthStyle",
        "f.store.Consume",
        "authentication.ErrAuthenticationUnavailable",
        "authentication.ErrAuthenticationInvalid",
    ],
)
check(
    "authorization-code flow does not create cookies",
    "http.SetCookie" not in flow and "SameSite" not in flow,
)
check(
    "authorization-code flow does not add an HTTP callback route",
    '"/auth/' not in flow and '"/login' not in flow,
)
check(
    "authorization-code flow never logs protocol secrets",
    "slog." not in flow and "log." not in flow,
)
check(
    "authorization-code flow does not enable insecure OIDC shortcuts",
    "InsecureSkipSignatureCheck" not in flow
    and "InsecureIssuerURLContext" not in flow
    and "SkipIssuerCheck" not in flow
    and "SkipClientIDCheck" not in flow,
)
check(
    "authorization-code flow never returns provider tokens",
    "AccessToken string" not in flow
    and "RefreshToken string" not in flow
    and "IDToken string" not in flow,
)

verifier = require(
    "internal/authentication/oidc/verifier.go",
    [
        "authorizationEndpoint",
        "tokenEndpoint",
        "tokenEndpointAuthMethods",
        "codeChallengeMethods",
        'json:"code_challenge_methods_supported"',
    ],
)
check(
    "OIDC verifier still prohibits insecure discovery overrides",
    "InsecureIssuerURLContext" not in verifier,
)

tests = require(
    "internal/authentication/oidc/authorization_code_test.go",
    [
        "TestAuthorizationCodeFlowBeginsWithBoundStateNonceAndPKCE",
        "TestAuthorizationCodeFlowCompletesExactlyOnce",
        "TestAuthorizationCodeFlowRejectsUnknownAndExpiredState",
        "TestAuthorizationCodeFlowAllowsOnlyOneConcurrentConsumer",
        "TestAuthorizationCodeFlowClassifiesInvalidCodeAndProviderOutage",
        "TestAuthorizationCodeFlowBoundsTokenResponse",
        "TestAuthorizationCodeFlowDoesNotExposeSecretsInErrors",
        "TestMemoryPreauthenticationStoreIsBoundedAndCleansExpiredState",
        "TestNewAuthorizationCodeFlowRejectsUnsafeConfiguration",
        "TestAuthorizationCodeFlowFailsClosedWhenRandomnessFails",
        "httptest.NewTLSServer",
        '"code_challenge_methods_supported"',
    ],
)
check(
    "authorization-code tests do not depend on a live provider",
    "accounts.google.com" not in tests
    and "login.microsoftonline.com" not in tests,
)

require(
    "docs/architecture/OIDC-AUTHORIZATION-CODE-AND-PKCE-TRANSACTION-IMPLEMENTATION.md",
    [
        "one atomic one-time consume operation",
        "raw state is returned",
        "SHA-256 state digest",
        "advertised `S256` support",
        "Consumption before exchange",
        "does not add HTTP login or callback routes",
        "does not establish",
    ],
)
require(
    "docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md",
    [
        "authorization-code and PKCE transaction implementation checkpoint",
        "durable restart-surviving preauthentication storage",
    ],
)
require(
    "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md",
    [
        "authorization-code and PKCE implementation status",
        "one-time state consumption",
    ],
)
require(
    "docs/testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md",
    [
        "authorization-code and PKCE transaction campaign",
        "exactly one concurrent consumer",
    ],
)
require(
    "test-framework/authentication/test_phase1_step3_oidc_authorization_code_pkce.sh",
    [
        "authorization-code and PKCE static validator",
        "go test -race",
        "go mod verify",
        "govulncheck",
    ],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_oidc_authorization_code_pkce.sh",
    [
        "36394c917a7c60350f229fc80df2066a0c132681",
        "implementation candidate only",
    ],
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text(encoding="utf-8").splitlines()
for path in (
    "internal/authentication/oidc/authorization_code.go",
    "internal/authentication/oidc/authorization_code_test.go",
    "docs/architecture/OIDC-AUTHORIZATION-CODE-AND-PKCE-TRANSACTION-IMPLEMENTATION.md",
    "tools/validation/validate_phase1_step3_oidc_authorization_code_pkce.py",
    "test-framework/authentication/test_phase1_step3_oidc_authorization_code_pkce.sh",
    "tools/validation/phase-gates/validate_phase1_step3_oidc_authorization_code_pkce.sh",
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
    print(
        "Phase 1 Step 3 OIDC authorization-code and PKCE transaction "
        "validation FAILED."
    )
    raise SystemExit(1)

print(
    "Phase 1 Step 3 OIDC authorization-code and PKCE transaction "
    "validation PASSED."
)
print(
    "This proves bounded state, nonce, PKCE S256, exact redirect, atomic "
    "one-time consume, code exchange, verified-principal, response-size, "
    "outage, replay, race, and redaction contracts."
)
print(
    "It does not prove HTTP login or callback routes, browser cookies, "
    "durable sessions, CSRF, logout, trusted proxies, production credential "
    "delivery, formal Step 3 acceptance, representative-provider "
    "compatibility, or production readiness."
)
