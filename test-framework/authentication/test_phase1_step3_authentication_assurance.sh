#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"

results_root="$repo_root/test-framework/test-results/authentication"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/phase1-step3-authentication-assurance-$run_id"
latest_report="$results_root/phase1-step3-authentication-assurance-latest-report.txt"

validation_report_init \
  "Phase 1 Step 3 authentication assurance regression" \
  "$report_dir"

validation_run \
  "repository-managed govulncheck preflight" \
  go tool govulncheck -version || true

validation_run \
  "authentication-assurance static validator" \
  python3 tools/validation/validate_phase1_step3_authentication_assurance.py || true

validation_run \
  "authentication-assurance format" \
  bash -c 'test -z "$(gofmt -l \
    internal/authentication/assurance \
    internal/authentication/oidc/verifier.go \
    internal/authentication/oidc/verifier_test.go \
    internal/authentication/oidc/http_handler.go \
    internal/authentication/oidc/http_handler_test.go \
    internal/authentication/session/session.go \
    internal/authentication/session/session_test.go \
    internal/authentication/session/postgresql/store.go \
    internal/authentication/session/postgresql/store_test.go)"' || true

validation_run \
  "authentication-assurance go test -race" \
  go test -race \
    ./internal/authentication/assurance \
    ./internal/authentication/oidc \
    ./internal/authentication/session \
    ./internal/authentication/session/postgresql || true

validation_run \
  "authentication-assurance go vet" \
  go vet \
    ./internal/authentication/assurance \
    ./internal/authentication/oidc \
    ./internal/authentication/session \
    ./internal/authentication/session/postgresql || true

validation_run \
  "authentication-assurance dependency module verification" \
  go mod verify || true

validation_run \
  "authentication-assurance dependency vulnerability analysis" \
  go tool govulncheck \
    ./internal/authentication/assurance \
    ./internal/authentication/oidc \
    ./internal/authentication/session/... || true

validation_run \
  "authentication-assurance phase-gate shell syntax" \
  bash -n tools/validation/phase-gates/validate_phase1_step3_authentication_assurance.sh || true

validation_run \
  "authentication-assurance regression shell syntax" \
  bash -n test-framework/authentication/test_phase1_step3_authentication_assurance.sh || true

validation_run \
  "Step 3 accepted record is absent" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md || true

validation_report_finish \
  "$report_dir/final-report.txt" \
  "$latest_report"
