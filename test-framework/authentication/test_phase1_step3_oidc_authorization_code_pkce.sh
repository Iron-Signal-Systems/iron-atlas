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
  "OIDC authorization-code and PKCE static validator" \
  python3 tools/validation/validate_phase1_step3_oidc_authorization_code_pkce.py

check \
  "OIDC authorization-code and PKCE format" \
  bash -c 'test -z "$(gofmt -l internal/authentication/oidc/authorization_code.go internal/authentication/oidc/authorization_code_test.go)"'

check \
  "OIDC authorization-code and PKCE go test -race" \
  go test -race ./internal/authentication/oidc

check \
  "OIDC authorization-code and PKCE go vet" \
  go vet ./internal/authentication/oidc

check \
  "OIDC authorization-code dependency module verification" \
  go mod verify

check \
  "OIDC authorization-code dependency vulnerability analysis" \
  go tool govulncheck ./internal/authentication/oidc

check \
  "OIDC authorization-code phase-gate shell syntax" \
  bash -n \
    tools/validation/phase-gates/validate_phase1_step3_oidc_authorization_code_pkce.sh

check \
  "OIDC authorization-code regression shell syntax" \
  bash -n \
    test-framework/authentication/test_phase1_step3_oidc_authorization_code_pkce.sh

check \
  "Step 3 accepted record is absent" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md

printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
