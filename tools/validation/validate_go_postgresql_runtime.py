#!/usr/bin/env python3
from pathlib import Path
import re
import sys

root = Path(__file__).resolve().parents[2]
errors: list[str] = []

required = [
    "go.mod",
    "internal/database/postgresql/pool.go",
    "internal/database/postgresql/pool_integration_test.go",
    "internal/change/postgresql/service.go",
    "internal/change/postgresql/service_integration_test.go",
    "internal/health/health.go",
    "docs/architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md",
    "docs/decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md",
    "docs/testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md",
    "tools/validation/lib/isolated_gate_revalidation.sh",
    "test-framework/phase-gates/test_isolated_gate_revalidation.sh",
]
for name in required:
    if not (root / name).is_file():
        errors.append(f"missing runtime file: {name}")

if errors:
    print("\n".join(f"FAIL: {error}" for error in errors), file=sys.stderr)
    raise SystemExit(1)

mod = (root / "go.mod").read_text()
if not re.search(r"^go 1\.(25|26)(?:\.\d+)?$", mod, re.MULTILINE):
    errors.append("Go module baseline must be Go 1.25 or 1.26")
if "github.com/jackc/pgx/v5 v5.10.0" not in mod:
    errors.append("pgx v5.10.0 must be pinned")
if "replace github.com/jackc/pgx" in mod:
    errors.append("local pgx replacement must not be committed")

pool = (root / "internal/database/postgresql/pool.go").read_text()
for token in [
    "pgxpool.ParseConfig",
    "rejectActorRuntimeParameters",
    "BeginTx",
    "set_config",
    "actorSetting",
    "true",
    "tx.Commit",
    "tx.Rollback",
]:
    if token not in pool:
        errors.append(f"pool boundary missing token: {token}")
if re.search(r"(?im)^\s*(?:SET|set)\s+atlas\.actor_id", pool):
    errors.append("actor context must not be set at session scope")
if "set_config($1, $2, false)" in pool:
    errors.append("actor context must not use session-scoped set_config")

service = (root / "internal/change/postgresql/service.go").read_text()
for token in [
    "atlas.create_change_request",
    "atlas.record_approval",
    "WithActor",
    "atlas.change_approval_summary",
]:
    if token not in service:
        errors.append(f"PostgreSQL change adapter missing token: {token}")
for prohibited in [
    "INSERT INTO atlas.change_request",
    "UPDATE atlas.change_request",
    "DELETE FROM atlas.change_request",
    "INSERT INTO atlas.approval_action",
]:
    if prohibited in service:
        errors.append(f"adapter bypasses accepted function boundary: {prohibited}")

integration = (root / "internal/database/postgresql/pool_integration_test.go").read_text()
for token in ["i < 500", "workers = 8", "iterations = 75", "RollbackClearsActorAndData"]:
    if token not in integration:
        errors.append(f"identity-isolation test missing token: {token}")

app = (root / "internal/app/app.go").read_text()
for token in [
    "IRON_ATLAS_CHANGE_STORE",
    "IRON_ATLAS_DATABASE_URL",
    "ChangeStorePostgreSQL",
    "developmentIdentityFromEnvironment",
    "return store == ChangeStoreMemory",
    "StartupTimeout",
]:
    if token not in app:
        errors.append(f"application runtime configuration missing token: {token}")


gate = (root / "tools/validation/phase-gates/validate_phase1_step2.sh").read_text()
for token in [
    "isolated_gate_revalidation.sh",
    "isolated_gate_revalidate",
    "validate_phase1_step1_acceptance.sh",
]:
    if token not in gate:
        errors.append(f"Step 2 gate missing failure-propagating predecessor token: {token}")

helper = (root / "tools/validation/lib/isolated_gate_revalidation.sh").read_text()
for token in ["rc=$?", 'return "$rc"', "validator is missing or not executable"]:
    if token not in helper:
        errors.append(f"isolated gate helper missing token: {token}")

gate_test = (root / "test-framework/phase-gates/test_isolated_gate_revalidation.sh").read_text()
for token in [
    "failing isolated predecessor validator returns failure",
    "missing isolated predecessor validator returns failure",
]:
    if token not in gate_test:
        errors.append(f"phase-gate failure propagation test missing token: {token}")

adr = (root / "docs/decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md").read_text()
for token in ["v5.10.0", "Go 1.25", "software supply chain", "transaction-local"]:
    if token not in adr:
        errors.append(f"pgx ADR missing token: {token}")

for path in root.rglob("*"):
    if not path.is_file() or ".git" in path.parts:
        continue
    if path.name in {"SOURCE-SHA256SUMS.txt", "FILE-MANIFEST.txt"}:
        continue
    try:
        text = path.read_text()
    except UnicodeDecodeError:
        continue
    if re.search(r"postgres(?:ql)?://[^\s:'\"]+:[^@\s]+@", text, re.IGNORECASE) and "REDACTED" not in text:
        errors.append(f"possible committed database credential in {path.relative_to(root)}")

if errors:
    print("\n".join(f"FAIL: {error}" for error in errors), file=sys.stderr)
    raise SystemExit(1)

print("validated Go PostgreSQL runtime, transaction-local identity, and integration-test contracts")
