#!/usr/bin/env python3
from pathlib import Path

root = Path(__file__).resolve().parents[2]
checks = {
    "requester independence": ("sql/schema/migrations/005_change_approvals.sql", "requester cannot approve own change"),
    "authority enforcement": ("sql/schema/migrations/005_change_approvals.sql", "change.approve"),
    "actor context": ("sql/schema/migrations/004_governed_changes.sql", "atlas.actor_id"),
    "append-only trigger": ("sql/schema/migrations/006_append_only_history.sql", "reject_append_only_mutation"),
    "migration advisory lock": ("tools/database/apply_migrations.sh", "pg_advisory_lock"),
    "application has no table DML grant": ("sql/bootstrap/runtime-grants.sql", "GRANT EXECUTE ON FUNCTION"),
}
errors=[]
for name,(filename,token) in checks.items():
    if token not in (root/filename).read_text():
        errors.append(f"{name}: missing {token} in {filename}")
text=(root/'sql/bootstrap/runtime-grants.sql').read_text().upper()
for prohibited in ('GRANT INSERT', 'GRANT UPDATE', 'GRANT DELETE', 'GRANT TRUNCATE', 'GRANT CREATE'):
    if prohibited in text:
        errors.append(f"runtime grants contain prohibited application grant: {prohibited}")
if errors:
    for error in errors: print(f"FAIL: {error}")
    raise SystemExit(1)
print(f"validated {len(checks)} database security static contracts")
