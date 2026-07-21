#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"

results_root="$repo_root/test-framework/test-results/authentication"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/phase1-step3-provider-neutral-assurance-evidence-$run_id"
latest_report="$results_root/phase1-step3-provider-neutral-assurance-evidence-latest-report.txt"

validation_report_init \
  "Phase 1 Step 3 provider-neutral assurance evidence regression" \
  "$report_dir"

validation_run \
  "provider-neutral assurance-evidence static validator" \
  python3 tools/validation/validate_phase1_step3_provider_neutral_assurance_evidence.py || true

validation_run \
  "provider-neutral assurance evidence format" \
  bash -c 'test -z "$(gofmt -l internal/authentication/assurance internal/authentication/oidc)"' || true

validation_run \
  "provider-neutral assurance evidence go test -race" \
  go test -race ./internal/authentication/assurance ./internal/authentication/oidc || true

validation_run \
  "provider-neutral assurance evidence go vet" \
  go vet ./internal/authentication/assurance ./internal/authentication/oidc || true

validation_run \
  "provider-neutral assurance evidence module verification" \
  go mod verify || true

validation_run \
  "repository-managed govulncheck preflight" \
  go tool govulncheck -version || true

validation_run \
  "provider-neutral assurance evidence vulnerability analysis" \
  go tool govulncheck ./internal/authentication/assurance ./internal/authentication/oidc || true

validation_run \
  "provider-neutral phase-gate shell syntax" \
  bash -n tools/validation/phase-gates/validate_phase1_step3_provider_neutral_assurance_evidence.sh || true

validation_run \
  "provider-neutral regression shell syntax" \
  bash -n test-framework/authentication/test_phase1_step3_provider_neutral_assurance_evidence.sh || true

validation_run \
  "formal Step 3 accepted record is absent" \
  test ! -e docs/acceptance/PHASE-1-STEP-3-ACCEPTANCE-RECORD.md || true

validation_report_finish \
  "$report_dir/final-report.txt" \
  "$latest_report"
