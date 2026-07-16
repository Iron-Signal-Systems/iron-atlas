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
  "governed actor resolution static validator" \
  python3 tools/validation/validate_phase1_step3_governed_actor_resolution.py

check \
  "governed actor resolution go test -race" \
  go test -race ./internal/authentication/postgresql

check \
  "governed actor resolution migration validation" \
  python3 tools/validation/validate_migrations.py

check \
  "governed actor resolution gate shell syntax" \
  bash -n \
    tools/validation/phase-gates/validate_phase1_step3_governed_actor_resolution.sh

check \
  "governed actor resolution regression shell syntax" \
  bash -n \
    test-framework/authentication/test_phase1_step3_governed_actor_resolution.sh

check \
  "Step 3 accepted record is absent" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md

printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
