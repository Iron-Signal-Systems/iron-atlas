#!/usr/bin/env python3
from __future__ import annotations

import stat
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]

REQUIRED = {
    "docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md": [
        "unrestricted execution context",
        "Production mode shall ignore and reject",
        "Provider groups, roles, administrative status",
        "transaction-local PostgreSQL actor context",
        "CSRF",
        "Trusted proxy",
    ],
    "docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md": [
        "IA-AUTH-004",
        "IA-AUTH-016",
        "Unique provider-subject mapping",
        "continued PostgreSQL pooled-identity isolation",
    ],
    "docs/testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md": [
        "Confused-deputy assertions",
        "session fixation",
        "identity remapping",
        "Redaction assertions",
    ],
    "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD-TEMPLATE.md": [
        "Trusted Proxy Profile",
        "Production-boundary exception prohibited",
        "Production readiness",
    ],
    "tools/validation/phase-gates/validate_phase1_step3_contract.sh": [
        "phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1",
        "1a750f7de791f567184c6f48e18eaec2933b8a14",
        "contract only",
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
        normalized_text = " ".join(text.split())
        for token in tokens:
            normalized_token = " ".join(token.split())
            check(
                f"{relative} contains {token}",
                normalized_token in normalized_text,
            )

    system_requirements = (
        ROOT / "docs/requirements/SYSTEM-REQUIREMENTS.md"
    ).read_text(encoding="utf-8")
    for number in range(1, 17):
        requirement = f"IA-AUTH-{number:03d}"
        check(
            f"system requirements contain {requirement}",
            requirement in system_requirements,
        )

    predecessor = (
        ROOT / "docs/acceptance/PHASE-1-STEP-2-ACCEPTANCE-RECORD.md"
    ).read_text(encoding="utf-8")
    for token in [
        "Trusted Authentication and Governed Actor Resolution Boundary",
        "ordinary headers must not select the actor",
        "session, cookie, CSRF, replay, logout, expiry, and trusted-proxy controls",
    ]:
        check(f"accepted predecessor contains {token}", token in predecessor)

    identity_sql = (
        ROOT / "sql/schema/migrations/002_governed_identity.sql"
    ).read_text(encoding="utf-8")
    for token in [
        "CREATE TABLE atlas.actor",
        "CREATE TABLE atlas.identity_provider",
        "CREATE TABLE atlas.external_identity",
        "UNIQUE (provider_id, provider_subject)",
    ]:
        check(f"accepted identity schema contains {token}", token in identity_sql)

    role_sql = (
        ROOT / "sql/schema/migrations/003_roles_and_authority.sql"
    ).read_text(encoding="utf-8")
    for token in [
        "CREATE TABLE atlas.role_binding",
        "valid_until timestamptz",
        "granted_by_actor_id",
        "grant_reason text NOT NULL",
    ]:
        check(f"accepted role schema contains {token}", token in role_sql)

    accepted_record = (
        ROOT / "docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md"
    )
    check(
        "Step 3 accepted record absent before formal acceptance",
        not accepted_record.exists(),
    )

    gate = (
        ROOT
        / "tools/validation/phase-gates/"
        "validate_phase1_step3_contract.sh"
    )
    check(
        "Step 3 contract gate executable",
        gate.is_file() and bool(gate.stat().st_mode & stat.S_IXUSR),
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
        print("Phase 1 Step 3 contract validation FAILED.")
        return 1

    print("Phase 1 Step 3 contract validation PASSED.")
    print(
        "This proves contract, requirements, traceability, test planning, "
        "acceptance-template, predecessor, and repository registration only."
    )
    print(
        "It does not prove executable production authentication, provider "
        "integration, sessions, CSRF, trusted-proxy enforcement, or "
        "production readiness."
    )
    return 0

if __name__ == "__main__":
    raise SystemExit(main())
