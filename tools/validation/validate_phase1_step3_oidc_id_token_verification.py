#!/usr/bin/env python3
from pathlib import Path
import re

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
    text = candidate.read_text()
    normalized = compact(text)
    for token in tokens:
        check(
            f"{path} contains {token}",
            compact(token) in normalized,
        )
    return text


verifier = require(
    "internal/authentication/oidc/verifier.go",
    [
        "coreoidc.NewProvider",
        "SupportedSigningAlgs",
        "authentication.ErrAuthenticationUnavailable",
        "authentication.ErrAuthenticationInvalid",
        "subtle.ConstantTimeCompare",
        "token.VerifyAccessToken",
        "rejectDuplicateFields",
        "authorization-code response support",
        "maximum token age",
    ],
)
check(
    "OIDC verifier never enables insecure signature verification",
    "InsecureSkipSignatureCheck" not in verifier,
)
check(
    "OIDC verifier never skips issuer verification",
    "SkipIssuerCheck" not in verifier,
)
check(
    "OIDC verifier never skips audience verification",
    "SkipClientIDCheck" not in verifier,
)
check(
    "OIDC verifier prohibits HTTP issuer shortcuts",
    "InsecureIssuerURLContext" not in verifier,
)
check(
    "OIDC verifier does not implement request-header identity",
    "X-Iron-Atlas-Actor" not in verifier
    and "X-Iron-Atlas-Roles" not in verifier,
)

tests = require(
    "internal/authentication/oidc/verifier_test.go",
    [
        "TestVerifierAcceptsValidToken",
        "TestVerifierEnforcesAccessTokenHash",
        "TestVerifierUsesBoundedAuthenticationTime",
        "TestVerifierFailsClosedForInvalidProtocolState",
        "TestVerifierRejectsDuplicateSensitiveClaim",
        "TestVerifierRejectsDuplicateSensitiveHeader",
        "TestVerifierRejectsOversizedOrMalformedToken",
        "TestVerifierRefreshesKeysOnRotation",
        "TestVerifierClassifiesKeyProviderOutage",
        "TestVerifierClassifiesUnknownKeyAsInvalidWhenProviderResponds",
        "TestVerifierClassifiesJWKSServiceFailureAsUnavailable",
        "TestVerifierSupportsConcurrentReadOnlyVerification",
        "TestNewRejectsInsecureOrUnboundedConfiguration",
        "httptest.NewTLSServer",
    ],
)
check(
    "OIDC tests do not depend on a live external provider",
    "accounts.google.com" not in tests
    and "login.microsoftonline.com" not in tests,
)

go_mod = require(
    "go.mod",
    [
        "github.com/coreos/go-oidc/v3 v3.19.0",
        "go 1.25.0",
        "toolchain go1.26.5",
    ],
)
check(
    "OIDC library is a direct pinned requirement",
    bool(re.search(
        r"(?m)^\s*github\.com/coreos/go-oidc/v3 v3\.19\.0\s*$",
        go_mod,
    )),
)

require(
    "docs/architecture/OIDC-ID-TOKEN-VERIFICATION-IMPLEMENTATION.md",
    [
        "exact HTTPS issuer",
        "no duplicate security-sensitive",
        "access-token hash verification",
        "does not implement a browser login route",
    ],
)
require(
    "docs/decisions/ADR-0006-OIDC-ID-TOKEN-VERIFICATION-LIBRARIES.md",
    [
        "v3.19.0",
        "Custom JOSE and OIDC implementation",
        "nonce verification",
        "vulnerability scanning",
    ],
)
require(
    "test-framework/authentication/test_phase1_step3_oidc_id_token_verification.sh",
    [
        "OIDC ID-token static validator",
        "go test -race",
        "go mod verify",
        "govulncheck",
    ],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_oidc_id_token_verification.sh",
    [
        "3ad3220c51179d3772d90da7f1025c4d41382922",
        "implementation candidate only",
    ],
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text().splitlines()
for path in (
    "internal/authentication/oidc/verifier.go",
    "internal/authentication/oidc/verifier_test.go",
    "docs/architecture/OIDC-ID-TOKEN-VERIFICATION-IMPLEMENTATION.md",
    "docs/decisions/ADR-0006-OIDC-ID-TOKEN-VERIFICATION-LIBRARIES.md",
    "tools/validation/validate_phase1_step3_oidc_id_token_verification.py",
    "test-framework/authentication/test_phase1_step3_oidc_id_token_verification.sh",
    "tools/validation/phase-gates/validate_phase1_step3_oidc_id_token_verification.sh",
):
    check(f"file manifest contains {path}", path in manifest)

print()
print(f"PASS checks: {passes}")
print(f"FAIL checks: {len(failures)}")
if failures:
    print("Phase 1 Step 3 OIDC ID-token verification validation FAILED.")
    raise SystemExit(1)

print("Phase 1 Step 3 OIDC ID-token verification validation PASSED.")
print(
    "This proves the bounded discovery, JWKS, signature, issuer, audience, "
    "authorized-party, nonce, time, stable-subject, access-token-hash, "
    "duplicate-field, rotation, outage, race, and concurrency contracts."
)
print(
    "It does not prove authorization-code exchange, PKCE transaction storage, "
    "browser sessions, cookies, CSRF, logout, trusted proxies, representative "
    "provider compatibility, formal Step 3 acceptance, or production readiness."
)
