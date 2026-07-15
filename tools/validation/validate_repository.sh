#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

skip_go=false
skip_database=false
for arg in "$@"; do
  [[ "$arg" == "--skip-go" ]] && skip_go=true
  [[ "$arg" == "--skip-database" ]] && skip_database=true
done

pass=0
fail=0
check() {
  local name="$1"; shift
  if "$@"; then printf 'PASS: %s\n' "$name"; pass=$((pass+1));
  else printf 'FAIL: %s\n' "$name"; fail=$((fail+1)); fi
}

required=(
  README.md SECURITY.md docs/README.md
  docs/architecture/TARGET-ARCHITECTURE.md
  docs/architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md
  docs/architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md
  docs/architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md
  docs/testing/POSTGRESQL-DISPOSABLE-DATABASE-TESTING.md
  cmd/atlasd/main.go internal/change/change.go integrations/zabbix/sender.go
  sql/schema/manifests/atlas.manifest
  sql/bootstrap/production-role-contract.sql
  tools/database/apply_migrations.sh
)
for file in "${required[@]}"; do check "required file $file" test -f "$file"; done

check "no raw evidence directories tracked" bash -c 'for d in ./raw-evidence ./evidence; do [[ ! -d "$d" ]] || ! find "$d" -type f -print -quit | grep -q .; done'
check "Markdown links" python3 tools/validation/validate_docs.py
check "migration contract" python3 tools/validation/validate_migrations.py
check "database security static contract" python3 tools/validation/validate_sql_static.py
check "Draw.io XML" python3 -c 'import xml.etree.ElementTree as ET; ET.parse("diagrams/source/curated/architecture/ARCH-001-iron-atlas-context.drawio")'
check "source SHA-256 records" python3 tools/validation/validate_source_checksums.py
check "file manifest synchronized" python3 - <<'PY'
from pathlib import Path
import subprocess
root=Path.cwd()
actual=sorted({x for x in subprocess.check_output(["git","ls-files","--cached","--others","--exclude-standard"],text=True).splitlines() if (root/x).is_file()})
recorded=(root/'FILE-MANIFEST.txt').read_text().splitlines()
raise SystemExit(0 if actual == recorded else 1)
PY

if ! $skip_go; then
  check "Go format" bash -c 'test -z "$(gofmt -l cmd internal modules integrations)"'
  check "Go vet" go vet ./...
  check "Go tests" go test ./...
fi
if ! $skip_database; then
  check "disposable PostgreSQL tests" ./test-framework/database/run_disposable_postgres.sh
fi

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
