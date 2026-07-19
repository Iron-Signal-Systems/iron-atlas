#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

implementation_base="28ec1eab5b5c4e69731e9b0a79fe6105beab316d"
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

printf '== Iron Atlas Phase 1 Step 3 HTTP login and callback boundary ==\n'

check   "merged PR #12 implementation base is an ancestor"   git merge-base --is-ancestor "$implementation_base" HEAD

check   "OIDC authorization-code and PKCE checkpoint remains valid"   ./tools/validation/phase-gates/validate_phase1_step3_oidc_authorization_code_pkce.sh

check   "HTTP login and callback static contract"   python3 tools/validation/validate_phase1_step3_http_login_callback.py

check   "HTTP login and callback regression"   ./test-framework/authentication/test_phase1_step3_http_login_callback.sh

check   "complete test framework"   ./test-framework/run_all.sh

check   "repository validation"   ./tools/validation/validate_repository.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
if (( fail != 0 )); then
  printf '\nPhase 1 Step 3 HTTP login and callback validation FAILED.\n'
  exit 1
fi

printf '\nPhase 1 Step 3 HTTP login and callback validation PASSED.\n'
printf '\nThis is an implementation candidate only. It establishes bounded GET login\n'
printf 'and callback handlers, secure short-lived browser state binding, exact\n'
printf 'callback cardinality, issuer binding, provider-error cancellation, replay\n'
printf 'resistance, generic failure classification, and a verified-principal handoff.\n'
printf '\nIt does not establish durable sessions, session cookies, protected-route\n'
printf 'authentication, session rotation or expiry, logout, revocation, CSRF,\n'
printf 'trusted-proxy enforcement, production application wiring, authentication\n'
printf 'audit persistence, production credential delivery, representative-provider\n'
printf 'compatibility, formal Step 3 acceptance, or production readiness.\n'
