#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"
source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"

implementation_base="6c912428a90b125f1b826729593e11ed914c12e9"
results_root="$repo_root/test-framework/test-results/phase-gates"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/phase1-step3-authenticated-session-$run_id"
latest_report="$results_root/phase1-step3-authenticated-session-latest-report.txt"

validation_report_init \
  "Iron Atlas Phase 1 Step 3 authenticated server-side session boundary" \
  "$report_dir"

revalidate_http_checkpoint() {
    isolated_gate_revalidate \
        "$repo_root" \
        "$implementation_base" \
        "tools/validation/phase-gates/validate_phase1_step3_http_login_callback.sh"
}

blocked=""

if ! validation_run \
  "SSH-signed post-PR-13 implementation base is an ancestor" \
  git merge-base --is-ancestor "$implementation_base" HEAD
then
  blocked="SSH-signed implementation-base ancestry check"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "HTTP login and callback checkpoint remains valid" \
    revalidate_http_checkpoint
  then
    blocked="HTTP login and callback checkpoint"
  fi
else
  validation_skip \
    "HTTP login and callback checkpoint remains valid" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "authenticated-session static contract" \
    python3 tools/validation/validate_phase1_step3_authenticated_session.py
  then
    blocked="authenticated-session static contract"
  fi
else
  validation_skip \
    "authenticated-session static contract" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "authenticated-session regression" \
    ./test-framework/authentication/test_phase1_step3_authenticated_session.sh
  then
    blocked="authenticated-session regression"
  fi
else
  validation_skip \
    "authenticated-session regression" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "complete test framework" \
    env IRON_ATLAS_HTTP_PREDECESSOR_ALREADY_VALIDATED=1 \
    ./test-framework/run_all.sh
  then
    blocked="complete test framework"
  fi
else
  validation_skip \
    "complete test framework" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  validation_run \
    "repository validation" \
    ./tools/validation/validate_repository.sh || true
else
  validation_skip \
    "repository validation" \
    "blocked by $blocked"
fi

validation_note \
  "Scope: implementation candidate only; this does not establish session rotation, sliding activity refresh, bounded session-count or cleanup policy, logout, administrative revocation workflow, CSRF, trusted-proxy enforcement, production application wiring, authentication audit persistence, MFA enforcement, local TOTP enrollment or recovery, representative-provider compatibility, formal Step 3 acceptance, or production readiness."

if validation_report_finish \
  "$report_dir/final-report.txt" \
  "$latest_report"
then
  exit 0
fi

exit 1
