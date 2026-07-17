#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

implementation_base="36394c917a7c60350f229fc80df2066a0c132681"
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

printf '== Iron Atlas Phase 1 Step 3 OIDC authorization-code and PKCE transaction ==\n'

check \
  "merged PR #6 implementation base is an ancestor" \
  git merge-base --is-ancestor "$implementation_base" HEAD

check \
  "OIDC ID-token verification checkpoint remains valid" \
  ./tools/validation/phase-gates/validate_phase1_step3_oidc_id_token_verification.sh

check \
  "OIDC authorization-code and PKCE static contract" \
  python3 tools/validation/validate_phase1_step3_oidc_authorization_code_pkce.py

check \
  "OIDC authorization-code and PKCE regression" \
  ./test-framework/authentication/test_phase1_step3_oidc_authorization_code_pkce.sh

check \
  "complete test framework" \
  ./test-framework/run_all.sh

check \
  "repository validation" \
  ./tools/validation/validate_repository.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
if (( fail != 0 )); then
  printf '\nPhase 1 Step 3 OIDC authorization-code and PKCE validation FAILED.\n'
  exit 1
fi

printf '\nPhase 1 Step 3 OIDC authorization-code and PKCE validation PASSED.\n'
printf '\nThis is an implementation candidate only. It does not establish HTTP login\n'
printf 'or callback routes, browser cookies, durable sessions, CSRF, logout,\n'
printf 'trusted-proxy enforcement, production credential delivery, formal Step 3\n'
printf 'acceptance, representative-provider compatibility, or production readiness.\n'
