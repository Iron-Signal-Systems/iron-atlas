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


authentication = require(
    "internal/authentication/authentication.go",
    [
        "type Assurance struct",
        "BoundActorID string",
        "func (a Assurance) Validate() error",
        "session-bound actor does not match current governed actor",
        "clonePrincipal",
    ],
)
check(
    "request context clones authentication methods",
    "cloned.Assurance.Methods = append(" in authentication
    and "principal.Assurance.Methods..." in authentication,
)

session = require(
    "internal/authentication/session/session.go",
    [
        '__Host-iron_atlas_session',
        "identifierBytes = 32",
        "sha256.Sum256",
        "base64.RawURLEncoding",
        "request.CookiesNamed",
        "http.SameSiteLaxMode",
        "ServeVerifiedPrincipal",
        "func (s *Service) Authenticate",
        "principal.BoundActorID = record.ActorID",
        "SecurityPolicyVersion",
        "Cache-Control",
        "no-store",
    ],
)
check(
    "session service does not log identifier material",
    "slog." not in session and "log." not in session,
)
check(
    "session identifier is not accepted from URL data",
    "Query().Get" not in session and "FormValue" not in session,
)
check(
    "session success redirect uses a narrow local-path allowlist",
    'strings.HasPrefix(raw, "//")' in session
    and "safeLocalPathCharacter" in session
    and "character == '/'" in session
    and "character == '~'" in session,
)
check(
    "session success redirect does not use general URL parsing",
    "url.Parse" not in session,
)
store = require(
    "internal/authentication/session/postgresql/store.go",
    [
        "atlas.create_authenticated_session",
        "atlas.authenticate_session",
        "ErrSessionNotFound",
        "ErrSessionUnavailable",
        "IdentifierDigest[:]",
    ],
)
check(
    "PostgreSQL store uses controlled routines rather than direct table SQL",
    "FROM atlas.authenticated_session" not in store
    and "INSERT INTO atlas.authenticated_session" not in store,
)

migration = require(
    "sql/schema/migrations/008_authenticated_session.sql",
    [
        "CREATE TABLE atlas.authenticated_session",
        "identifier_digest bytea NOT NULL UNIQUE",
        "octet_length(identifier_digest) = 32",
        "authenticated_session_provider_fk",
        "authenticated_session_actor_fk",
        "CREATE OR REPLACE FUNCTION atlas.create_authenticated_session",
        "CREATE OR REPLACE FUNCTION atlas.authenticate_session",
        "SECURITY DEFINER",
        "SET search_path = pg_catalog, atlas",
        "transaction_timestamp() < s.idle_expires_at",
        "transaction_timestamp() < s.absolute_expires_at",
        "REVOKE ALL ON TABLE atlas.authenticated_session FROM PUBLIC",
    ],
)
check(
    "migration does not grant session table access to application role",
    "GRANT SELECT ON atlas.authenticated_session" not in migration,
)

database_tests = require(
    "test-framework/database/test_database.py",
    [
        "session created before identity remap",
        "session authenticates before identity remap",
        "identity remap invalidates existing session",
    ],
)
check(
    "database remap test uses controlled session routines",
    "atlas.create_authenticated_session" in database_tests
    and "atlas.authenticate_session" in database_tests,
)

require(
    "sql/schema/manifests/atlas.manifest",
    ["migrations/008_authenticated_session.sql"],
)
require(
    "sql/bootstrap/runtime-grants.sql",
    [
        "GRANT EXECUTE ON FUNCTION atlas.create_authenticated_session",
        "GRANT EXECUTE ON FUNCTION atlas.authenticate_session(bytea)",
    ],
)

unit_tests = require(
    "internal/authentication/session/session_test.go",
    [
        "TestVerifiedPrincipalCreatesDigestOnlySecureSession",
        "TestVerifiedPrincipalRejectsAlreadySessionBoundPrincipal",
        "TestVerifiedPrincipalFailsClosedForRandomnessAndStoreFailure",
        "TestAuthenticateReturnsServerBoundPrincipal",
        "TestAuthenticateRejectsDuplicateMalformedUnknownAndExpiredCookies",
        "TestAuthenticateClassifiesStoreOutage",
        "TestServiceConcurrentAuthentication",
        "TestNewRejectsUnsafeConfiguration",
    ],
)
check(
    "session tests do not contain a live provider dependency",
    "accounts.google.com" not in unit_tests
    and "login.microsoftonline.com" not in unit_tests,
)
check(
    "session tests do not print raw cookie structures",
    'cookies = %#v' not in unit_tests and 'cookie = %#v' not in unit_tests,
)
require(
    "internal/authentication/assurance_test.go",
    [
        "TestAssuranceValidation",
        "TestResolvedIdentityRejectsSessionActorRemapping",
        "TestResolvedIdentityClonesAssuranceMethods",
    ],
)
store_unit_tests = require(
    "internal/authentication/session/postgresql/store_test.go",
    [
        "TestStoreCreateUsesControlledFunction",
        "TestStoreFindMapsControlledResult",
    ],
)
store_integration_tests = require(
    "internal/authentication/session/postgresql/store_integration_test.go",
    [
        "TestIntegrationStoreCreatesAndAuthenticatesSession",
        "TestIntegrationStoreRejectsActorMismatchAndUnknownLookup",
        "TestIntegrationApplicationCannotReadSessionTable",
        "TestIntegrationConcurrentSessionLookup",
    ],
)
check(
    "session tests avoid complete record or cookie structure output",
    "%#v" not in unit_tests
    and "%#v" not in store_unit_tests
    and "%#v" not in store_integration_tests,
)

require(
    "docs/architecture/AUTHENTICATED-SERVER-SIDE-SESSION-IMPLEMENTATION.md",
    [
        "digest-only persistence",
        "current governed actor and role re-resolution",
        "narrow ASCII allowlist",
        "Authentication assurance and future MFA",
        "RFC 6238 TOTP",
        "WebAuthn",
        "does not establish",
    ],
)
require(
    "docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md",
    [
        "Authenticated server-side session implementation checkpoint",
        "Authentication assurance and MFA roadmap",
        "Google Authenticator",
        "1Password",
    ],
)
require(
    "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md",
    [
        "Authenticated-session implementation status",
        "authentication assurance and MFA policy",
        "TOTP enrollment, verification, and recovery",
    ],
)
require(
    "docs/requirements/SYSTEM-REQUIREMENTS.md",
    [
        "IA-AUTH-017",
        "IA-AUTH-018",
        "IA-AUTH-019",
        "WebAuthn",
        "RFC 6238 TOTP",
        "silent administrator bypass",
    ],
)
require(
    "docs/testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md",
    [
        "Authenticated-session implementation campaign",
        "Authentication assurance and MFA successor campaigns",
    ],
)
require(
    "docs/roadmap/IMPLEMENTATION-ROADMAP.md",
    [
        "authenticated server-side session is the active bounded",
        "Provider-neutral authentication assurance",
        "phishing-resistant MFA",
        "RFC 6238 TOTP fallback",
    ],
)
require(
    "docs/roadmap/PHASE-GATE-PLAN.md",
    [
        "validate_phase1_step3_authenticated_session.sh",
        "validate_phase1_step3_authentication_assurance.sh",
        "validate_phase1_step3_totp_enrollment_verification_recovery.sh",
    ],
)

require(
    "test-framework/authentication/test_phase1_step3_authenticated_session.sh",
    ["go test -race", "go mod verify", "govulncheck"],
)
require(
    "test-framework/run_all.sh",
    [
        "Phase 1 Step 3 HTTP login and callback predecessor revalidation",
        "isolated_gate_revalidate",
        "Phase 1 Step 3 authenticated-session static validation",
        "Phase 1 Step 3 authenticated-session regression",
    ],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_authenticated_session.sh",
    [
        "6c912428a90b125f1b826729593e11ed914c12e9",
        "implementation candidate only",
        "does not establish session rotation",
    ],
)

app = require("internal/app/app.go", [])
check(
    "authenticated-session checkpoint remains deliberately unwired from atlasd",
    "session.New(" not in app
    and "sessionpostgresql.New(" not in app
    and "CookieName" not in app,
)

manifest = (ROOT / "FILE-MANIFEST.txt").read_text(encoding="utf-8").splitlines()
for path in (
    "docs/architecture/AUTHENTICATED-SERVER-SIDE-SESSION-IMPLEMENTATION.md",
    "internal/authentication/assurance_test.go",
    "internal/authentication/session/session.go",
    "internal/authentication/session/session_test.go",
    "internal/authentication/session/postgresql/store.go",
    "internal/authentication/session/postgresql/store_test.go",
    "internal/authentication/session/postgresql/store_integration_test.go",
    "sql/schema/migrations/008_authenticated_session.sql",
    "test-framework/authentication/test_phase1_step3_authenticated_session.sh",
    "tools/validation/phase-gates/validate_phase1_step3_authenticated_session.sh",
    "tools/validation/validate_phase1_step3_authenticated_session.py",
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
    print("Phase 1 Step 3 authenticated-session validation FAILED.")
    raise SystemExit(1)

print("Phase 1 Step 3 authenticated-session validation PASSED.")
print(
    "This proves digest-only PostgreSQL session persistence, secure opaque browser "
    "cookies, bounded creation-time validity, protected-request authentication, "
    "current governed actor re-resolution, actor-remapping rejection, assurance "
    "metadata retention, least-privileged controlled routines, and hostile and "
    "concurrent behavior."
)
print(
    "It does not prove session rotation, sliding activity refresh, bounded session-"
    "count or cleanup policy, logout, administrative revocation workflow, CSRF, "
    "trusted proxies, production wiring, "
    "MFA enforcement, local TOTP enrollment or recovery, representative-provider "
    "compatibility, formal Step 3 acceptance, or production readiness."
)
