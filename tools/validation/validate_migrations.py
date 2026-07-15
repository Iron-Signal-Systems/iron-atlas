#!/usr/bin/env python3
from pathlib import Path
import re
import sys

root = Path(__file__).resolve().parents[2]
manifest = root / "sql/schema/manifests/atlas.manifest"
errors: list[str] = []

required_order = [
    "BEGIN;",
    "SET LOCAL lock_timeout = '5s';",
    "SET LOCAL statement_timeout = '1min';",
    "SET LOCAL idle_in_transaction_session_timeout = '1min';",
    "SET LOCAL ROLE atlas_schema_owner;",
    "COMMIT;",
]

entries: list[str] = []
for raw in manifest.read_text().splitlines():
    line = raw.split("#", 1)[0].strip()
    if line:
        entries.append(line)

if not entries:
    errors.append("migration manifest is empty")

expected_version = 1
for entry in entries:
    match = re.fullmatch(r"migrations/(\d{3})_[a-z0-9_]+\.sql", entry)
    if not match:
        errors.append(f"invalid manifest entry: {entry}")
        continue
    version = int(match.group(1))
    if version != expected_version:
        errors.append(f"expected migration {expected_version:03d}, found {version:03d}")
    expected_version += 1

    path = manifest.parent.parent / entry
    if not path.is_file():
        errors.append(f"missing migration: {entry}")
        continue
    text = path.read_text()
    positions = []
    for token in required_order:
        if token not in text:
            errors.append(f"{entry}: missing {token}")
        else:
            positions.append(text.index(token))
    if len(positions) == len(required_order) and positions != sorted(positions):
        errors.append(f"{entry}: execution-safety tokens are out of order")
    if re.search(r"\bCREATE\s+ROLE\b", text, re.IGNORECASE):
        errors.append(f"{entry}: role creation belongs in bootstrap, not migrations")
    if re.search(r"\bGRANT\b[^;]*\bTO\s+PUBLIC\b", text, re.IGNORECASE | re.DOTALL):
        errors.append(f"{entry}: grants to PUBLIC are prohibited")
    for function in re.finditer(
        r"CREATE\s+(?:OR\s+REPLACE\s+)?FUNCTION.*?\$\$;",
        text,
        re.IGNORECASE | re.DOTALL,
    ):
        block = function.group(0)
        if "SECURITY DEFINER" in block.upper() and "SET search_path" not in block:
            errors.append(f"{entry}: SECURITY DEFINER function lacks fixed search_path")

candidate = root / "sql/schema/candidates/000_initial_atlas_schema.sql"
if not candidate.is_file():
    errors.append("Phase 0 SQL candidate is not archived")
if any("000_initial_atlas_schema.sql" in entry for entry in entries):
    errors.append("Phase 0 SQL candidate must not remain in executable manifest")

role_contract = (root / "sql/bootstrap/production-role-contract.sql").read_text()
for role in (
    "atlas_database_owner", "atlas_schema_owner", "atlas_migrator",
    "atlas_application", "atlas_readonly", "atlas_auditor",
):
    if role not in role_contract:
        errors.append(f"production role contract missing {role}")
if "SUPERUSER" in role_contract and "NOSUPERUSER" not in role_contract:
    errors.append("production role contract may not create superusers")

if errors:
    print("\n".join(f"FAIL: {error}" for error in errors), file=sys.stderr)
    raise SystemExit(1)

print(f"validated {len(entries)} ordered migrations and production role contract")
