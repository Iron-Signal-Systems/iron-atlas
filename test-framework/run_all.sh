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

revalidate_http_checkpoint() {
    isolated_gate_revalidate         "$repo_root"         "6c912428a90b125f1b826729593e11ed914c12e9"         "tools/validation/phase-gates/validate_phase1_step3_http_login_callback.sh"
}

run "go format check" bash -c 'test -z "$(gofmt -l cmd internal modules integrations)"'
run "go module verification" go mod verify
run "go vet" go vet ./...
run "go test" go test -race ./...
run "migration static validation" python3 tools/validation/validate_migrations.py
run "database security static validation" python3 tools/validation/validate_sql_static.py
run "Go PostgreSQL runtime static validation" python3 tools/validation/validate_go_postgresql_runtime.py
run "Phase 1 Step 3 contract static validation" python3 tools/validation/validate_phase1_step3_contract.py
run "Phase 1 Step 3 contract regression" ./test-framework/authentication/test_phase1_step3_contract.sh
run "Phase 1 Step 3 authentication foundation static validation" python3 tools/validation/validate_phase1_step3_authentication_foundation.py
run "Phase 1 Step 3 authentication foundation regression" ./test-framework/authentication/test_phase1_step3_authentication_foundation.sh
run "Phase 1 Step 3 governed actor resolution static validation" python3 tools/validation/validate_phase1_step3_governed_actor_resolution.py
run "Phase 1 Step 3 governed actor resolution regression" ./test-framework/authentication/test_phase1_step3_governed_actor_resolution.sh
run "Phase 1 Step 3 OIDC ID-token verification static validation" python3 tools/validation/validate_phase1_step3_oidc_id_token_verification.py
run "Phase 1 Step 3 OIDC ID-token verification regression" ./test-framework/authentication/test_phase1_step3_oidc_id_token_verification.sh
run "Phase 1 Step 3 OIDC authorization-code and PKCE static validation" python3 tools/validation/validate_phase1_step3_oidc_authorization_code_pkce.py
run "Phase 1 Step 3 OIDC authorization-code and PKCE regression" ./test-framework/authentication/test_phase1_step3_oidc_authorization_code_pkce.sh
if [[ "${IRON_ATLAS_HTTP_PREDECESSOR_ALREADY_VALIDATED:-0}" == "1" ]]; then
    validation_skip         "Phase 1 Step 3 HTTP login and callback predecessor revalidation"         "already validated by the calling phase gate"
else
    run         "Phase 1 Step 3 HTTP login and callback predecessor revalidation"         revalidate_http_checkpoint
fi
run "Phase 1 Step 3 authenticated-session static validation" python3 tools/validation/validate_phase1_step3_authenticated_session.py
run "Phase 1 Step 3 authenticated-session regression" ./test-framework/authentication/test_phase1_step3_authenticated_session.sh
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
