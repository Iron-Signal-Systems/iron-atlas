#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
results_dir="$repo_root/test-framework/test-results"
mkdir -p "$results_dir"
log="$results_dir/latest.log"
summary="$results_dir/latest-summary.txt"

exec > >(tee "$log") 2>&1
cd "$repo_root"

pass=0
fail=0
run() {
  local name="$1"; shift
  printf '\n== %s ==\n' "$name"
  if "$@"; then printf 'PASS: %s\n' "$name"; pass=$((pass+1));
  else printf 'FAIL: %s\n' "$name"; fail=$((fail+1)); fi
}

run "go format check" bash -c 'test -z "$(gofmt -l cmd internal modules integrations)"'
run "go module verification" go mod verify
run "go vet" go vet ./...
run "go test" go test -race ./...
run "migration static validation" python3 tools/validation/validate_migrations.py
run "database security static validation" python3 tools/validation/validate_sql_static.py
run "Go PostgreSQL runtime static validation" python3 tools/validation/validate_go_postgresql_runtime.py
run "Phase 1 Step 3 contract static validation" python3 tools/validation/validate_phase1_step3_contract.py
run "Phase 1 Step 3 contract regression" ./test-framework/authentication/test_phase1_step3_contract.sh
run "phase-gate exit propagation" ./test-framework/phase-gates/test_isolated_gate_revalidation.sh
run "external toolchain validation" python3 tools/validation/validate_toolchain.py
run "portable acceptance static validation" python3 tools/validation/validate_portable_acceptance.py
run "portable validation regression" ./test-framework/portability/test_portable_validation.sh
run "committed validation evidence" python3 tools/validation/validate_committed_evidence.py
run "disposable PostgreSQL tests" ./test-framework/database/run_disposable_postgres.sh
run "repository validation" ./tools/validation/validate_repository.sh --skip-go --skip-database

{
  echo "PASS checks: $pass"
  echo "FAIL checks: $fail"
  if (( fail == 0 )); then echo "Correctness result: PASS"; else echo "Correctness result: FAIL"; fi
  echo "Resource observation: RECORDED_BY_DATABASE_TEST"
  echo "Performance thresholds: NOT_EVALUATED"
} | tee "$summary"

(( fail == 0 ))
