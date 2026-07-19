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

check   "HTTP login and callback static validator"   python3 tools/validation/validate_phase1_step3_http_login_callback.py

check   "HTTP login and callback format"   bash -c 'test -z "$(gofmt -l     internal/authentication/oidc/authorization_code.go     internal/authentication/oidc/authorization_code_cancel_test.go     internal/authentication/oidc/http_handler.go     internal/authentication/oidc/http_handler_test.go)"'

check   "HTTP login and callback go test -race"   go test -race ./internal/authentication/oidc

check   "HTTP login and callback go vet"   go vet ./internal/authentication/oidc

check   "HTTP login and callback dependency module verification"   go mod verify

check   "HTTP login and callback dependency vulnerability analysis"   govulncheck ./internal/authentication/oidc

check   "HTTP login and callback phase-gate shell syntax"   bash -n tools/validation/phase-gates/validate_phase1_step3_http_login_callback.sh

check   "HTTP login and callback regression shell syntax"   bash -n test-framework/authentication/test_phase1_step3_http_login_callback.sh

check   "Step 3 accepted record is absent"   test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md

printf 'PASS checks: %d\nFAIL checks: %d\n' "$pass" "$fail"
(( fail == 0 ))
