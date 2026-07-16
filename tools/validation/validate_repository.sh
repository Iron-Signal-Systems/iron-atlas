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
  docs/architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md
  docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md
  docs/architecture/PORTABLE-VALIDATION-AND-CANONICAL-REPOSITORY-ACCEPTANCE.md
  docs/decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md
  docs/decisions/ADR-0005-CANONICAL-REPOSITORY-REPRODUCIBILITY.md
  docs/decisions/ADR-0006-OIDC-ID-TOKEN-VERIFICATION-LIBRARIES.md
  docs/architecture/OIDC-ID-TOKEN-VERIFICATION-IMPLEMENTATION.md
  docs/testing/POSTGRESQL-DISPOSABLE-DATABASE-TESTING.md
  docs/testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md
  docs/testing/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION-TESTING.md
  docs/requirements/PHASE-1-STEP-3-REQUIREMENTS-TRACEABILITY.md
  docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD-TEMPLATE.md
  docs/operations/CANONICAL-CLEAN-CLONE-VALIDATION.md
  validation/toolchain-requirements.json
  validation/evidence/README.md
  cmd/atlasd/main.go internal/change/change.go
  internal/authentication/authentication.go
  internal/authentication/authentication_test.go
  internal/authentication/postgresql/resolver.go
  internal/authentication/postgresql/resolver_test.go
  internal/authentication/postgresql/resolver_integration_test.go
  internal/authentication/oidc/verifier.go
  internal/authentication/oidc/verifier_test.go
  internal/database/postgresql/pool.go
  internal/change/postgresql/service.go
  internal/health/health.go
  integrations/zabbix/sender.go
  sql/schema/manifests/atlas.manifest
  sql/bootstrap/production-role-contract.sql
  tools/database/apply_migrations.sh
  tools/validation/validate_go_postgresql_runtime.py
  tools/validation/validate_phase1_step3_contract.py
  tools/validation/validate_phase1_step3_authentication_foundation.py
  tools/validation/validate_phase1_step3_governed_actor_resolution.py
  tools/validation/validate_phase1_step3_oidc_id_token_verification.py
  tools/validation/validate_portable_acceptance.py
  tools/validation/validate_committed_evidence.py
  tools/validation/validate_toolchain.py
  tools/validation/record_validation_evidence.sh
  tools/validation/verify_canonical_clone.sh
  tools/validation/lib/isolated_gate_revalidation.sh
  tools/validation/phase-gates/validate_phase1_step2.sh
  tools/validation/phase-gates/validate_phase1_step3_contract.sh
  tools/validation/phase-gates/validate_phase1_step3_authentication_foundation.sh
  tools/validation/phase-gates/validate_phase1_step3_governed_actor_resolution.sh
  tools/validation/phase-gates/validate_phase1_step3_oidc_id_token_verification.sh
  test-framework/authentication/test_phase1_step3_contract.sh
  test-framework/authentication/test_phase1_step3_authentication_foundation.sh
  test-framework/authentication/test_phase1_step3_governed_actor_resolution.sh
  test-framework/authentication/test_phase1_step3_oidc_id_token_verification.sh
  test-framework/phase-gates/test_isolated_gate_revalidation.sh
  test-framework/portability/test_portable_validation.sh
)
for file in "${required[@]}"; do check "required file $file" test -f "$file"; done

check "no raw evidence directories tracked" bash -c 'for d in ./raw-evidence ./evidence; do [[ ! -d "$d" ]] || ! find "$d" -type f -print -quit | grep -q .; done'
check "no committed database secrets" bash -c '! git grep -nEI "postgres(ql)?://[^[:space:]:\"'"'"']+:[^@[:space:]]+@" -- ":(exclude)SOURCE-SHA256SUMS.txt" ":(exclude)FILE-MANIFEST.txt" | grep -v REDACTED'
check "Markdown links" python3 tools/validation/validate_docs.py
check "migration contract" python3 tools/validation/validate_migrations.py
check "database security static contract" python3 tools/validation/validate_sql_static.py
check "Go PostgreSQL runtime static contract" python3 tools/validation/validate_go_postgresql_runtime.py
check "Phase 1 Step 3 contract" python3 tools/validation/validate_phase1_step3_contract.py
check "Phase 1 Step 3 authentication foundation" python3 tools/validation/validate_phase1_step3_authentication_foundation.py
check "Phase 1 Step 3 governed actor resolution" python3 tools/validation/validate_phase1_step3_governed_actor_resolution.py
check "Phase 1 Step 3 OIDC ID-token verification" python3 tools/validation/validate_phase1_step3_oidc_id_token_verification.py
check "portable acceptance static contract" python3 tools/validation/validate_portable_acceptance.py
check "committed validation evidence" python3 tools/validation/validate_committed_evidence.py
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
  check "Go module verification" go mod verify
  check "Go vet" go vet ./...
  check "Go tests" go test ./...
fi
if ! $skip_database; then
  check "disposable PostgreSQL tests" ./test-framework/database/run_disposable_postgres.sh
fi

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
