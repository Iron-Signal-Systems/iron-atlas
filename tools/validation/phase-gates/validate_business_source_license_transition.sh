#!/usr/bin/env bash
set -Eeuo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"

implementation_base="cc93fdd2311ca188ad03b0bd94293156ff243973"
results_root="$repo_root/test-framework/test-results/phase-gates"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/business-source-license-transition-$run_id"
latest_report="$results_root/business-source-license-transition-latest-report.txt"

validation_report_init \
  "Atlas Business Source License 1.1 transition" \
  "$report_dir"

blocked=""

if ! validation_run \
  "SSH-signed authentication-assurance boundary is an ancestor" \
  git merge-base --is-ancestor "$implementation_base" HEAD
then
  blocked="signed implementation-base ancestry check"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "Business Source License 1.1 static contract" \
    python3 tools/validation/validate_licensing.py
  then
    blocked="licensing static contract"
  fi
else
  validation_skip \
    "Business Source License 1.1 static contract" \
    "blocked by $blocked"
fi

if [[ -z "$blocked" ]]; then
  if ! validation_run \
    "Business Source License 1.1 regression" \
    ./test-framework/governance/test_business_source_license_transition.sh
  then
    blocked="licensing regression"
  fi
else
  validation_skip \
    "Business Source License 1.1 regression" \
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
  "Scope: prospective repository licensing and governance candidate only. This preserves the historical BSD license, declares BUSL-1.1 parameters, separates trademarks, and records deferred alignment work. It does not provide legal advice, alter runtime behavior, grant production-use rights, complete architecture alignment, complete Phase 1 Step 3, or establish production readiness."

if validation_report_finish \
  "$report_dir/final-report.txt" \
  "$latest_report"
then
  exit 0
fi

exit 1
