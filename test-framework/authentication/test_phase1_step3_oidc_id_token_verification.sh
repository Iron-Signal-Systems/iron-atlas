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
  "OIDC ID-token static validator" \
  python3 tools/validation/validate_phase1_step3_oidc_id_token_verification.py

check \
  "OIDC ID-token go test -race" \
  go test -race ./internal/authentication/oidc

check \
  "OIDC ID-token go vet" \
  go vet ./internal/authentication/oidc

check \
  "OIDC dependency module verification" \
  go mod verify

check \
  "OIDC dependency vulnerability analysis" \
  go tool govulncheck ./internal/authentication/oidc

check \
  "OIDC ID-token phase-gate shell syntax" \
  bash -n \
    tools/validation/phase-gates/validate_phase1_step3_oidc_id_token_verification.sh

check \
  "OIDC ID-token regression shell syntax" \
  bash -n \
    test-framework/authentication/test_phase1_step3_oidc_id_token_verification.sh

check \
  "Step 3 accepted record is absent" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md

printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
