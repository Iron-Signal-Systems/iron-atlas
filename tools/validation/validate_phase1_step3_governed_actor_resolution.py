#!/usr/bin/env python3
from pathlib import Path
import re

ROOT = Path(__file__).resolve().parents[2]
failures: list[str] = []
passes = 0


def normalized(text: str) -> str:
    return " ".join(text.split())


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
    compact = normalized(text)
    for token in tokens:
        check(f"{path} contains {token}", normalized(token) in compact)
    return text


migration = require(
    "sql/schema/migrations/007_governed_actor_resolution.sql",
    [
        "CREATE OR REPLACE FUNCTION atlas.resolve_governed_actor",
        "RETURNS TABLE",
        "SECURITY DEFINER",
        "SET search_path = pg_catalog, atlas",
        "ip.active",
        "a.actor_status = 'ACTIVE'",
        "rd.active",
        "rb.valid_from <= transaction_timestamp()",
        "REVOKE ALL ON FUNCTION atlas.resolve_governed_actor",
    ],
)
check(
    "resolver migration does not grant table SELECT",
    "GRANT SELECT" not in migration.upper(),
)

resolver = require(
    "internal/authentication/postgresql/resolver.go",
    [
        "var _ authentication.ActorResolver = (*Resolver)(nil)",
        "atlas.resolve_governed_actor($1, $2)",
        "authentication.ErrIdentityResolutionFailed",
        "authentication.ErrAuthenticationUnavailable",
        "PLATFORM_ADMINISTRATOR",
    ],
)
check(
    "resolver does not query governed tables directly",
    not re.search(
        r"FROM\s+atlas\.(actor|identity_provider|external_identity|role_binding)",
        resolver,
        re.IGNORECASE,
    ),
)

require(
    "internal/authentication/postgresql/resolver_test.go",
    [
        "TestResolverMapsGovernedRoles",
        "TestResolverFailsClosedForMissingMapping",
        "TestResolverClassifiesDatabaseFailureAsUnavailable",
        "TestResolverRejectsUnknownOrDuplicateRoles",
        "TestResolverRejectsUnnormalizedPrincipal",
    ],
)
require(
    "internal/authentication/postgresql/resolver_integration_test.go",
    [
        "TestIntegrationResolverLoadsActiveActorAndRoles",
        "TestIntegrationResolverRejectsInactiveOrUnknownState",
        "TestIntegrationResolverExcludesExpiredAndInactiveRoles",
        "TestIntegrationResolverConcurrentIsolation",
    ],
)

grants = require(
    "sql/bootstrap/runtime-grants.sql",
    [
        "GRANT EXECUTE ON FUNCTION atlas.resolve_governed_actor(text, text) TO atlas_application",
    ],
)
for table in (
    "atlas.actor",
    "atlas.identity_provider",
    "atlas.external_identity",
    "atlas.role_definition",
    "atlas.role_binding",
):
    check(
        f"application has no broad SELECT grant on {table}",
        not re.search(
            rf"GRANT\s+SELECT\s+ON[^;]*\b{re.escape(table)}\b[^;]*"
            r"\bTO\s+atlas_application\b",
            grants,
            re.IGNORECASE | re.DOTALL,
        ),
    )

require(
    "sql/schema/manifests/atlas.manifest",
    ["migrations/007_governed_actor_resolution.sql"],
)
require(
    "docs/architecture/POSTGRESQL-GOVERNED-ACTOR-RESOLUTION-IMPLEMENTATION.md",
    [
        "Least-privilege database interface",
        "returns no row",
        "empty role set",
        "does not implement an external authentication provider",
    ],
)
require(
    "test-framework/authentication/test_phase1_step3_governed_actor_resolution.sh",
    [
        "governed actor resolution static validator",
        "go test -race",
    ],
)
require(
    "tools/validation/phase-gates/validate_phase1_step3_governed_actor_resolution.sh",
    [
        "c6ad0d8d5c6268e5bd850eae646bd2e21ed7f3f5",
        "implementation candidate only",
    ],
)

file_manifest = (ROOT / "FILE-MANIFEST.txt").read_text().splitlines()
for path in (
    "sql/schema/migrations/007_governed_actor_resolution.sql",
    "internal/authentication/postgresql/resolver.go",
    "internal/authentication/postgresql/resolver_test.go",
    "internal/authentication/postgresql/resolver_integration_test.go",
    "docs/architecture/POSTGRESQL-GOVERNED-ACTOR-RESOLUTION-IMPLEMENTATION.md",
    "tools/validation/validate_phase1_step3_governed_actor_resolution.py",
    "test-framework/authentication/test_phase1_step3_governed_actor_resolution.sh",
    "tools/validation/phase-gates/validate_phase1_step3_governed_actor_resolution.sh",
):
    check(f"file manifest contains {path}", path in file_manifest)

print()
print(f"PASS checks: {passes}")
print(f"FAIL checks: {len(failures)}")
if failures:
    print("Phase 1 Step 3 governed actor resolution validation FAILED.")
    raise SystemExit(1)

print("Phase 1 Step 3 governed actor resolution validation PASSED.")
print(
    "This proves the least-privileged PostgreSQL resolver function, explicit "
    "role mapping, fail-closed governed-state handling, targeted tests, "
    "documentation, and repository registration."
)
print(
    "It does not prove an external provider adapter, sessions, CSRF, "
    "trusted-proxy enforcement, or production readiness."
)
