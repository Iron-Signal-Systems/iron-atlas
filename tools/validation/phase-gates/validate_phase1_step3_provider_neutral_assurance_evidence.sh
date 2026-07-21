#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"
source "$repo_root/tools/validation/lib/isolated_gate_revalidation.sh"

implementation_base="2347d21f779768f40496a93cb1d9140cc3b6e0ce"
results_root="$repo_root/test-framework/test-results/phase-gates"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/phase1-step3-provider-neutral-assurance-evidence-$run_id"
latest_report="$results_root/phase1-step3-provider-neutral-assurance-evidence-latest-report.txt"

validation_report_init \
  "Iron Atlas Phase 1 Step 3 provider-neutral assurance evidence boundary" \
  "$report_dir"

revalidate_alignment_boundary() {
    isolated_gate_revalidate \
        "$repo_root" \
        "$implementation_base" \
        "tools/validation/phase-gates/validate_architecture_roadmap_alignment.sh"
}

blocked=""

if ! validation_run \
  "SSH-signed architecture-alignment evidence-closure base is an ancestor" \
  git merge-base --is-ancestor "$implementation_base" HEAD
then
  blocked="SSH-signed implementation-base ancestry check"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "architecture-alignment evidence boundary remains valid" \
    revalidate_alignment_boundary
  then
    blocked="architecture-alignment evidence boundary"
  fi
else
  validation_skip \
    "architecture-alignment evidence boundary remains valid" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "provider-neutral assurance-evidence static contract" \
    python3 tools/validation/validate_phase1_step3_provider_neutral_assurance_evidence.py
  then
    blocked="provider-neutral assurance-evidence static contract"
  fi
else
  validation_skip \
    "provider-neutral assurance-evidence static contract" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "provider-neutral assurance-evidence regression" \
    ./test-framework/authentication/test_phase1_step3_provider_neutral_assurance_evidence.sh
  then
    blocked="provider-neutral assurance-evidence regression"
  fi
else
  validation_skip \
    "provider-neutral assurance-evidence regression" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "complete test framework" \
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
  "Scope: bounded implementation candidate only. This establishes Atlas-controlled provider-neutral assurance evidence, explicit auth_time correlation, exact governed method sets, and fail-closed synthetic policy cases. It does not establish compatibility with a named provider, live hosted MFA, session lifecycle completion, CSRF, trusted proxies, emergency access, production wiring, formal Phase 1 Step 3 acceptance, or production readiness."

if validation_report_finish \
  "$report_dir/final-report.txt" \
  "$latest_report"
then
  exit 0
fi

exit 1
