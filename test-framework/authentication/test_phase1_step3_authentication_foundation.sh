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

check \
  "authentication foundation static validator" \
  python3 tools/validation/validate_phase1_step3_authentication_foundation.py

check \
  "authentication foundation go test -race" \
  go test -race ./internal/authentication ./internal/app ./internal/httpui

check \
  "authentication foundation gate shell syntax" \
  bash -n \
    tools/validation/phase-gates/validate_phase1_step3_authentication_foundation.sh

check \
  "authentication foundation regression shell syntax" \
  bash -n \
    test-framework/authentication/test_phase1_step3_authentication_foundation.sh

check \
  "Step 3 accepted record is absent" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md

printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
