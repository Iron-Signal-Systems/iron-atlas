#!/usr/bin/env python3
from __future__ import annotations

import stat
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]

REQUIRED = {
    "internal/authentication/authentication.go": [
        'ModeDevelopment Mode = "development"',
        'ModeProduction  Mode = "production"',
        "type Authenticator interface",
        "type ActorResolver interface",
        "ResolvedIdentityFromContext",
        "development identity headers are prohibited",
        "authentication service unavailable",
        "publicPath",
    ],
    "internal/authentication/authentication_test.go": [
        "TestDevelopmentModeInjectsImmutableResolvedIdentity",
        "TestProductionModeRejectsDevelopmentHeadersBeforeAdapter",
        "TestProductionModeWithoutAdapterFailsClosed",
        "TestProductionAdapterAndResolverSetServerSideIdentity",
        "TestNestedIdentityMiddlewareFailsClosed",
    ],
    "internal/app/app.go": [
        "IRON_ATLAS_AUTHENTICATION_MODE",
        "IRON_ATLAS_DEV_IDENTITY is no longer supported",
        "AuthenticationMode authentication.Mode",
        "authentication.ModeProduction",
    ],
    "internal/app/app_test.go": [
        "TestPostgreSQLModeDefaultsToProductionAuthentication",
        "TestLegacyDevelopmentIdentitySettingIsRejected",
        "TestInvalidAuthenticationModeIsRejected",
    ],
    "internal/httpui/server.go": [
        "Authentication *authentication.Middleware",
        "authentication.ActorFromContext",
        "DevelopmentMode",
    ],
    "internal/httpui/server_test.go": [
        "TestProductionModeRejectsDevelopmentHeaders",
        "TestProductionModeRequiresAuthentication",
        "TestQueryCannotSelectDevelopmentActor",
    ],
    "tools/validation/phase-gates/validate_phase1_step3_authentication_foundation.sh": [
        "ce57772440c17035f808048609de8596b0f18944",
        "authentication foundation",
        "implementation candidate only",
    ],
    "test-framework/authentication/test_phase1_step3_authentication_foundation.sh": [
        "authentication foundation static validator",
        "go test -race",
    ],
}

def main() -> int:
    pass_count = 0
    failures: list[str] = []

    def check(name: str, condition: bool) -> None:
        nonlocal pass_count
        if condition:
            pass_count += 1
            print(f"PASS: {name}")
        else:
            failures.append(name)
            print(f"FAIL: {name}")

    for relative, tokens in REQUIRED.items():
        path = ROOT / relative
        check(f"required file {relative}", path.is_file())
        if not path.is_file():
            continue
        text = path.read_text(encoding="utf-8")
        normalized = " ".join(text.split())
        for token in tokens:
            check(
                f"{relative} contains {token}",
                " ".join(token.split()) in normalized,
            )

    app_text = (ROOT / "internal/app/app.go").read_text(encoding="utf-8")
    server_text = (
        ROOT / "internal/httpui/server.go"
    ).read_text(encoding="utf-8")
    main_text = (ROOT / "cmd/atlasd/main.go").read_text(encoding="utf-8")

    check(
        "application no longer exposes DevelopmentIdentity boolean",
        "DevelopmentIdentity" not in app_text,
    )
    check(
        "HTTP handlers no longer read development identity headers",
        "X-Iron-Atlas-Actor" not in server_text
        and "X-Iron-Atlas-Roles" not in server_text,
    )
    check(
        "daemon logs typed authentication mode",
        '"authentication_mode"' in main_text
        and "development_identity" not in main_text,
    )

    readme = (ROOT / "README.md").read_text(encoding="utf-8")
    for token in [
        "IRON_ATLAS_AUTHENTICATION_MODE=development",
        "PostgreSQL mode defaults to `production`",
        "protected routes fail closed",
        "IRON_ATLAS_DEV_IDENTITY",
    ]:
        check(f"README contains {token}", token in readme)

    architecture = (
        ROOT
        / "docs/architecture/"
        "TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md"
    ).read_text(encoding="utf-8")
    for token in [
        "Authentication foundation implementation checkpoint",
        "No external provider adapter is accepted",
        "private immutable request context",
    ]:
        check(f"architecture contains {token}", token in architecture)

    traceability = (
        ROOT
        / "docs/requirements/"
        "PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md"
    ).read_text(encoding="utf-8")
    for token in [
        "Authentication foundation implementation status",
        "IA-AUTH-008",
        "IA-AUTH-009",
        "IA-AUTH-014",
        "IA-AUTH-016",
    ]:
        check(f"traceability contains {token}", token in traceability)

    accepted_record = (
        ROOT / "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md"
    )
    check(
        "Step 3 accepted record remains absent",
        not accepted_record.exists(),
    )

    gate = (
        ROOT
        / "tools/validation/phase-gates/"
        "validate_phase1_step3_authentication_foundation.sh"
    )
    regression = (
        ROOT
        / "test-framework/authentication/"
        "test_phase1_step3_authentication_foundation.sh"
    )
    check(
        "authentication foundation gate is executable",
        gate.is_file() and bool(gate.stat().st_mode & stat.S_IXUSR),
    )
    check(
        "authentication foundation regression is executable",
        regression.is_file()
        and bool(regression.stat().st_mode & stat.S_IXUSR),
    )

    manifest = (ROOT / "FILE-MANIFEST.txt").read_text(
        encoding="utf-8"
    ).splitlines()
    for relative in REQUIRED:
        check(f"file manifest contains {relative}", relative in manifest)

    print()
    print(f"PASS checks: {pass_count}")
    print(f"FAIL checks: {len(failures)}")
    if failures:
        print(
            "Phase 1 Step 3 authentication foundation "
            "validation FAILED."
        )
        return 1

    print(
        "Phase 1 Step 3 authentication foundation "
        "validation PASSED."
    )
    print(
        "This proves typed authentication-mode configuration, "
        "development-header isolation, production fail-closed behavior, "
        "immutable request-context identity, adapter and resolver seams, "
        "tests, documentation, and repository registration."
    )
    print(
        "It does not prove a production identity-provider adapter, "
        "provider-backed actor resolution, sessions, CSRF, trusted-proxy "
        "enforcement, or production readiness."
    )
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
