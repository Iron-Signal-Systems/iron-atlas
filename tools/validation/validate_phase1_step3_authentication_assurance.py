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


assurance = require(
    "internal/authentication/assurance/assurance.go",
    [
        "type PolicyConfig struct",
        "OutcomeAdditionalAuthenticationRequired",
        "OutcomePhishingResistantRequired",
        "MaximumAuthenticationAge",
        "AcceptedMFAMethodSets",
        "PhishingResistantRoles",
        "SecurityPolicyVersion",
        "ServeVerifiedPrincipal",
        "additional authentication required",
    ],
)
check(
    "assurance policy does not consume provider roles or groups",
    "groups" not in assurance.lower() and "provider roles" not in assurance.lower(),
)
check(
    "assurance service is stateless",
    "sync.Mutex" not in assurance and "map[string]Decision" not in assurance,
)

verifier = require(
    "internal/authentication/oidc/verifier.go",
    [
        'json:"acr"',
        'json:"amr"',
        '"acr": {}',
        '"amr": {}',
        "normalizedAssuranceClaims",
        "Assurance: assurance",
    ],
)
check(
    "OIDC verifier does not infer MFA success",
    "MFAAuthenticated: true" not in verifier,
)

handler = require(
    "internal/authentication/oidc/http_handler.go",
    [
        "issuer string",
        "issuer: issuer",
        "callback.Issuer != h.issuer",
        "provider error metadata requires an error result",
        "session state is permitted only on successful callbacks",
    ],
)
check(
    "callback issuer is not re-read during comparison",
    "callback.Issuer != h.flow.IssuerURL()" not in handler,
)

session = require(
    "internal/authentication/session/session.go",
    [
        "RequireMFA bool",
        "authenticated sessions require MFA assurance enforcement",
        "principal.Assurance.SecurityPolicyVersion != s.securityPolicyVersion",
        "!principal.Assurance.MFAAuthenticated",
        "session MFA assurance is required",
    ],
)
check(
    "session service does not manufacture policy assurance",
    "principal.Assurance.SecurityPolicyVersion = s.securityPolicyVersion" not in session,
)

store = require(
    "internal/authentication/session/postgresql/store.go",
    [
        "session MFA assurance is required",
        "principal and session security policy versions differ",
    ],
)

require(
    "sql/schema/migrations/009_authentication_assurance.sql",
    [
        "authenticated_session_requires_mfa",
        "MFA assurance policy introduced",
        "mfa_authenticated",
        "mfa_authenticated_at",
    ],
)
require(
    "sql/schema/manifests/atlas.manifest",
    ["migrations/009_authentication_assurance.sql"],
)

unit_tests = require(
    "internal/authentication/assurance/assurance_test.go",
    [
        "TestPolicyAcceptsExplicitProviderMFA",
        "TestPolicyDoesNotInferMFAFromUnknownClaims",
        "TestPolicyRequiresPhishingResistanceForHighImpactRoles",
        "TestPolicyRequiresFreshAuthentication",
        "TestServiceForwardsOnlySatisfiedAssurance",
        "TestServiceRejectsUnsatisfiedAssuranceWithoutSessionHandoff",
        "TestServiceConcurrentEvaluation",
    ],
)
check(
    "assurance tests contain no live identity provider dependency",
    "accounts.google.com" not in unit_tests
    and "login.microsoftonline.com" not in unit_tests,
)

require(
    "internal/authentication/oidc/verifier_test.go",
    [
        "TestVerifierNormalizesAssuranceClaims",
        "TestVerifierRejectsMalformedAssuranceClaims",
        "TestVerifierRejectsDuplicateAssuranceClaims",
    ],
)
require(
    "internal/authentication/oidc/http_handler_test.go",
    [
        "TestHTTPHandlerRetainsValidatedIssuer",
        "TestHTTPCallbackRejectsConflictingProviderMetadata",
    ],
)
require(
    "internal/authentication/session/session_test.go",
    ["TestVerifiedPrincipalRequiresSatisfiedMFAPolicy"],
)

require(
    "docs/architecture/AUTHENTICATION-ASSURANCE-IMPLEMENTATION.md",
    [
        "provider-neutral authentication-assurance normalization",
        "Only `satisfied` may reach the authenticated-session service",
        "does not infer MFA",
        "governed RFC 6238 TOTP lifecycle",
    ],
)
require(
    "docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md",
    ["Authentication assurance implementation checkpoint"],
)
require(
    "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md",
    ["Authentication-assurance implementation status"],
)
require(
    "docs/testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md",
    ["Authentication-assurance implementation campaign"],
)
require(
    "docs/roadmap/IMPLEMENTATION-ROADMAP.md",
    ["authentication assurance is the active bounded trusted-authentication candidate"],
)
require(
    "docs/roadmap/PHASE-GATE-PLAN.md",
    ["The active implementation candidate is `validate_phase1_step3_authentication_assurance.sh`"],
)
require(
    "docs/README.md",
    ["Authentication assurance implementation"],
)

require(
    "test-framework/authentication/test_phase1_step3_authentication_assurance.sh",
    ["go test -race", "go tool govulncheck", "authentication assurance"],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_authentication_assurance.sh",
    [
        "e4ae9de5a5757d1a53c04f0b17163919bc688b04",
        "authenticated-session checkpoint remains valid",
        "implementation candidate only",
    ],
)
require(
    "test-framework/run_all.sh",
    [
        "authenticated-session predecessor revalidation",
        "authentication-assurance static validation",
        "authentication-assurance regression",
    ],
)

app = require("internal/app/app.go", [])
check(
    "assurance checkpoint remains deliberately unwired from atlasd",
    "assurance.NewService(" not in app,
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text(encoding="utf-8").splitlines()
for path in (
    "docs/architecture/AUTHENTICATION-ASSURANCE-IMPLEMENTATION.md",
    "internal/authentication/assurance/assurance.go",
    "internal/authentication/assurance/assurance_test.go",
    "sql/schema/migrations/009_authentication_assurance.sql",
    "test-framework/authentication/test_phase1_step3_authentication_assurance.sh",
    "tools/validation/phase-gates/validate_phase1_step3_authentication_assurance.sh",
    "tools/validation/validate_phase1_step3_authentication_assurance.py",
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
    print("Phase 1 Step 3 authentication-assurance validation FAILED.")
    raise SystemExit(1)

print("Phase 1 Step 3 authentication-assurance validation PASSED.")
print(
    "This proves provider-neutral acr, amr, and auth_time normalization, exact "
    "versioned assurance policy, required MFA enforcement before session creation, "
    "role-sensitive phishing-resistant outcomes, stale-authentication step-up, "
    "callback hardening, and hostile and concurrent behavior."
)
print(
    "It does not prove local TOTP enrollment, QR-code generation, TOTP verification, "
    "recovery codes, WebAuthn, session lifecycle completion, CSRF, trusted proxies, "
    "production wiring, representative-provider compatibility, formal Step 3 "
    "acceptance, or production readiness."
)
