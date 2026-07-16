#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

accepted_predecessor="c6ad0d8d5c6268e5bd850eae646bd2e21ed7f3f5"
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

printf '== Iron Atlas Phase 1 Step 3 governed actor resolution ==\n'

check \
  "accepted authentication foundation merge is an ancestor" \
  git merge-base --is-ancestor "$accepted_predecessor" HEAD

check \
  "accepted authentication foundation remains valid" \
  ./tools/validation/phase-gates/validate_phase1_step3_authentication_foundation.sh

check \
  "governed actor resolution static contract" \
  python3 tools/validation/validate_phase1_step3_governed_actor_resolution.py

check \
  "governed actor resolution regression" \
  ./test-framework/authentication/test_phase1_step3_governed_actor_resolution.sh

check \
  "complete test framework" \
  ./test-framework/run_all.sh

check \
  "repository validation" \
  ./tools/validation/validate_repository.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
if (( fail != 0 )); then
  printf '\nPhase 1 Step 3 governed actor resolution validation FAILED.\n'
  exit 1
fi

printf '\nPhase 1 Step 3 governed actor resolution validation PASSED.\n'
printf '\nThis is an implementation candidate only. It does not establish an\n'
printf 'external authentication provider, sessions, CSRF, trusted-proxy\n'
printf 'enforcement, formal Step 3 acceptance, or production readiness.\n'
