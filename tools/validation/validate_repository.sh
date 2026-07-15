#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

skip_go=false
[[ "${1:-}" == "--skip-go" ]] && skip_go=true

pass=0
fail=0
check() {
  local name="$1"; shift
  if "$@"; then
    printf 'PASS: %s\n' "$name"
    pass=$((pass+1))
  else
    printf 'FAIL: %s\n' "$name"
    fail=$((fail+1))
  fi
}

required=(
  README.md
  SECURITY.md
  docs/README.md
  docs/architecture/TARGET-ARCHITECTURE.md
  docs/architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md
  docs/architecture/FIREWALL-CONFIGURATION-SEMANTIC-ANALYSIS.md
  docs/architecture/CISCO-EVIDENCE-COLLECTION-AND-PREVENTIVE-HEALTH.md
  docs/architecture/CISCO-TRUNK-AND-ENDPOINT-ATTRIBUTION.md
  docs/architecture/EXTERNAL-SYSTEM-INDEPENDENT-TELEMETRY.md
  docs/testing/TESTING-AND-ADVERSARIAL-VALIDATION-MODEL.md
  cmd/atlasd/main.go
  internal/change/change.go
  integrations/zabbix/sender.go
  sql/schema/manifests/atlas.manifest
)
for file in "${required[@]}"; do check "required file $file" test -f "$file"; done

check "no raw evidence directories tracked" bash -c 'for d in ./raw-evidence ./evidence; do [[ ! -d "$d" ]] || ! find "$d" -type f -print -quit | grep -q .; done'
check "Markdown links" python3 tools/validation/validate_docs.py
check "migration contract" python3 tools/validation/validate_migrations.py
check "Draw.io XML" python3 -c 'import xml.etree.ElementTree as ET; ET.parse("diagrams/source/curated/architecture/ARCH-001-iron-atlas-context.drawio")'

if ! $skip_go; then
  check "Go format" bash -c 'test -z "$(gofmt -l cmd internal modules integrations)"'
  check "Go vet" go vet ./...
  check "Go tests" go test ./...
fi

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
