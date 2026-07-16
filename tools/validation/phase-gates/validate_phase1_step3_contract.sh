#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"

predecessor_tag="phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1"
accepted_dev_merge="1a750f7de791f567184c6f48e18eaec2933b8a14"

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

echo "== Iron Atlas Phase 1 Step 3 contract validation =="

check "accepted Step 2 tag exists" \
  git rev-parse --verify -q "${predecessor_tag}^{commit}"

check "accepted Step 2 tag is an ancestor" \
  git merge-base --is-ancestor "$predecessor_tag" HEAD

check "accepted ISRAS-enabled dev merge is an ancestor" \
  git merge-base --is-ancestor "$accepted_dev_merge" HEAD

check "Step 3 static contract" \
  python3 tools/validation/validate_phase1_step3_contract.py

check "repository validation without disposable database" \
  ./tools/validation/validate_repository.sh --skip-database

check "complete test framework" \
  ./test-framework/run_all.sh

printf '\nPASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"

if (( fail != 0 )); then
  echo "Phase 1 Step 3 contract validation FAILED."
  exit 1
fi

cat <<'MSG'

Phase 1 Step 3 contract validation PASSED.

This is a phase-entry contract only. It proves that requirements,
architecture, traceability, testing, acceptance-template, predecessor, and
repository registration are synchronized.

It does not prove executable production authentication, provider integration,
sessions, CSRF, trusted-proxy enforcement, or production readiness.
MSG
