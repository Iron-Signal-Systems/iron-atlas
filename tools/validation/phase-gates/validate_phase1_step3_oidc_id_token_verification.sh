#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

accepted_predecessor="3ad3220c51179d3772d90da7f1025c4d41382922"
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

printf '== Iron Atlas Phase 1 Step 3 OIDC ID-token verification ==\n'

check \
  "accepted governed actor-resolution merge is an ancestor" \
  git merge-base --is-ancestor "$accepted_predecessor" HEAD

check \
  "accepted governed actor-resolution checkpoint remains valid" \
  ./tools/validation/phase-gates/validate_phase1_step3_governed_actor_resolution.sh

check \
  "OIDC ID-token verification static contract" \
  python3 tools/validation/validate_phase1_step3_oidc_id_token_verification.py

check \
  "OIDC ID-token verification regression" \
  ./test-framework/authentication/test_phase1_step3_oidc_id_token_verification.sh

check \
  "complete test framework" \
  ./test-framework/run_all.sh

check \
  "repository validation" \
  ./tools/validation/validate_repository.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
if (( fail != 0 )); then
  printf '\nPhase 1 Step 3 OIDC ID-token verification validation FAILED.\n'
  exit 1
fi

printf '\nPhase 1 Step 3 OIDC ID-token verification validation PASSED.\n'
printf '\nThis is an implementation candidate only. It does not establish\n'
printf 'authorization-code exchange, PKCE transaction storage, browser sessions,\n'
printf 'cookies, CSRF, logout, trusted-proxy enforcement, formal Step 3\n'
printf 'acceptance, representative-provider compatibility, or production readiness.\n'
