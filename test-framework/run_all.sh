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
  if "$@"; then
    printf 'PASS: %s\n' "$name"
    pass=$((pass+1))
  else
    printf 'FAIL: %s\n' "$name"
    fail=$((fail+1))
  fi
}

run "go format check" bash -c 'test -z "$(gofmt -l cmd internal modules integrations)"'
run "go vet" go vet ./...
run "go test" go test -race ./...
run "repository validation" ./tools/validation/validate_repository.sh --skip-go

{
  echo "PASS checks: $pass"
  echo "FAIL checks: $fail"
  if (( fail == 0 )); then echo "Correctness result: PASS"; else echo "Correctness result: FAIL"; fi
  echo "Resource observation: NOT_RECORDED"
  echo "Performance thresholds: NOT_EVALUATED"
} | tee "$summary"

(( fail == 0 ))
