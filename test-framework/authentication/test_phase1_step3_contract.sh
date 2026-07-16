#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

pass=0
fail=0

check() {
  local name="$1"
  shift
  if "$@"; then
    printf 'PASS: %s\n' "$name"
    pass=$((pass + 1))
  else
    printf 'FAIL: %s\n' "$name"
    fail=$((fail + 1))
  fi
}

check "Step 3 static contract" \
  python3 tools/validation/validate_phase1_step3_contract.py

check "Step 3 contract gate shell syntax" \
  bash -n tools/validation/phase-gates/validate_phase1_step3_contract.sh

check "Step 3 regression shell syntax" \
  bash -n test-framework/authentication/test_phase1_step3_contract.sh

check "Step 3 acceptance record is not present" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md

printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
