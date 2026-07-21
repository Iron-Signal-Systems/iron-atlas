#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
results_dir="$repo_root/test-framework/test-results"
mkdir -p "$results_dir"
log="$results_dir/latest.log"
summary="$results_dir/latest-summary.txt"

exec > >(tee "$log") 2>&1
cd "$repo_root"

source "$repo_root/tools/validation/lib/reporting.sh"
source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"

run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_dir/validation-reporting/run-all-$run_id"
validation_report_init "Iron Atlas complete test framework" "$report_dir"

run() {
    validation_run "$@" || true
}

revalidate_authentication_assurance_checkpoint() {
    isolated_gate_revalidate \
        "$repo_root" \
        "cc93fdd2311ca188ad03b0bd94293156ff243973" \
        "tools/validation/phase-gates/validate_phase1_step3_authentication_assurance.sh"
}

run "Business Source License 1.1 static validation" python3 tools/validation/validate_licensing.py
run "Business Source License 1.1 regression" ./test-framework/governance/test_business_source_license_transition.sh
run "architecture and roadmap alignment static validation" python3 tools/validation/validate_architecture_roadmap_alignment.py
run "architecture and roadmap alignment regression" ./test-framework/governance/test_architecture_roadmap_alignment.sh
run "provider-neutral assurance-evidence static validation" python3 tools/validation/validate_phase1_step3_provider_neutral_assurance_evidence.py
run "provider-neutral assurance-evidence regression" ./test-framework/authentication/test_phase1_step3_provider_neutral_assurance_evidence.sh
run "go format check" bash -c 'test -z "$(gofmt -l cmd internal modules integrations)"'
run "go module verification" go mod verify
run "go vet" go vet ./...
run "go test" go test -race ./...
run "migration static validation" python3 tools/validation/validate_migrations.py
run "database security static validation" python3 tools/validation/validate_sql_static.py
run "Go PostgreSQL runtime static validation" python3 tools/validation/validate_go_postgresql_runtime.py
run "Phase 1 Step 3 authentication-assurance checkpoint revalidation" revalidate_authentication_assurance_checkpoint
run "validation reporting static validation" python3 tools/validation/validate_validation_reporting.py
run "validation reporting regression" ./test-framework/validation/test_validation_reporting.sh
run "phase-gate exit propagation" ./test-framework/phase-gates/test_isolated_gate_revalidation.sh
run "external toolchain validation" python3 tools/validation/validate_toolchain.py
run "portable acceptance static validation" python3 tools/validation/validate_portable_acceptance.py
run "portable validation regression" ./test-framework/portability/test_portable_validation.sh
run "committed validation evidence" python3 tools/validation/validate_committed_evidence.py
run "disposable PostgreSQL tests" ./test-framework/database/run_disposable_postgres.sh
run "repository validation" ./tools/validation/validate_repository.sh --skip-go --skip-database

validation_note "Resource observation: RECORDED_BY_DATABASE_TEST"
validation_note "Performance thresholds: NOT_EVALUATED"

report_status=0
validation_report_finish "$report_dir/final-report.txt" "$summary" || report_status=$?
exit "$report_status"
