#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"
source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"

implementation_base="e4ae9de5a5757d1a53c04f0b17163919bc688b04"
results_root="$repo_root/test-framework/test-results/phase-gates"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/phase1-step3-authentication-assurance-$run_id"
latest_report="$results_root/phase1-step3-authentication-assurance-latest-report.txt"

validation_report_init \
  "Iron Atlas Phase 1 Step 3 authentication assurance boundary" \
  "$report_dir"

revalidate_session_checkpoint() {
    isolated_gate_revalidate \
        "$repo_root" \
        "$implementation_base" \
        "tools/validation/phase-gates/validate_phase1_step3_authenticated_session.sh"
}

blocked=""

if ! validation_run \
  "SSH-signed post-PR-14 implementation base is an ancestor" \
  git merge-base --is-ancestor "$implementation_base" HEAD
then
  blocked="SSH-signed implementation-base ancestry check"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "authenticated-session checkpoint remains valid" \
    revalidate_session_checkpoint
  then
    blocked="authenticated-session checkpoint"
  fi
else
  validation_skip \
    "authenticated-session checkpoint remains valid" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "authentication-assurance static contract" \
    python3 tools/validation/validate_phase1_step3_authentication_assurance.py
  then
    blocked="authentication-assurance static contract"
  fi
else
  validation_skip \
    "authentication-assurance static contract" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "authentication-assurance regression" \
    ./test-framework/authentication/test_phase1_step3_authentication_assurance.sh
  then
    blocked="authentication-assurance regression"
  fi
else
  validation_skip \
    "authentication-assurance regression" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "complete test framework" \
    env IRON_ATLAS_SESSION_PREDECESSOR_ALREADY_VALIDATED=1 \
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
  "Scope: implementation candidate only; this establishes provider-neutral assurance normalization and policy enforcement before session creation. It does not establish local TOTP enrollment, QR-code generation, TOTP verification, recovery codes, WebAuthn, session rotation, logout, administrative revocation workflow, CSRF, trusted-proxy enforcement, production application wiring, representative-provider compatibility, formal Step 3 acceptance, or production readiness."

if validation_report_finish \
  "$report_dir/final-report.txt" \
  "$latest_report"
then
  exit 0
fi

exit 1
