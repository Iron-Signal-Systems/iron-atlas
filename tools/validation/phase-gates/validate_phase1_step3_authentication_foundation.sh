#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

contract_merge="ce57772440c17035f808048609de8596b0f18944"

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

echo "== Iron Atlas Phase 1 Step 3 authentication foundation =="

check \
  "accepted Step 3 contract merge is an ancestor" \
  git merge-base --is-ancestor "$contract_merge" HEAD

check \
  "accepted Step 3 contract remains valid" \
  ./tools/validation/phase-gates/validate_phase1_step3_contract.sh

check \
  "authentication foundation static contract" \
  python3 tools/validation/validate_phase1_step3_authentication_foundation.py

check \
  "targeted authentication foundation race tests" \
  go test -race ./internal/authentication ./internal/app ./internal/httpui

check \
  "repository validation without disposable database" \
  ./tools/validation/validate_repository.sh --skip-database

check \
  "complete test framework" \
  ./test-framework/run_all.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"

if (( fail != 0 )); then
  echo "Phase 1 Step 3 authentication foundation validation FAILED."
  exit 1
fi

cat <<'MSG'

Phase 1 Step 3 authentication foundation validation PASSED.

This is an implementation candidate only. It proves typed authentication
modes, development-header isolation, production fail-closed behavior,
immutable request-context identity, future adapter/resolver seams, tests,
documentation synchronization, and preserved predecessor controls.

It does not prove production identity-provider integration, provider-backed
actor resolution, sessions, CSRF, trusted-proxy enforcement, or production
readiness.
MSG
