#!/usr/bin/env bash
set -Eeuo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$repo_root"
source "$repo_root/tools/validation/lib/reporting.sh"
implementation_base="a0ab1ad19cf48ba11d97b3a9e87acd7b68e1eb60"
results_root="$repo_root/test-framework/test-results/phase-gates"
run_id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
report_dir="$results_root/architecture-roadmap-alignment-$run_id"
latest_report="$results_root/architecture-roadmap-alignment-latest-report.txt"
validation_report_init "Iron Atlas architecture and roadmap alignment boundary" "$report_dir"
blocked=""
if ! validation_run "SSH-signed BUSL boundary is an ancestor" git merge-base --is-ancestor "$implementation_base" HEAD; then blocked="signed base ancestry"; fi
if [[ -z "$blocked" ]]; then
  if ! validation_run "alignment static contract" python3 tools/validation/validate_architecture_roadmap_alignment.py; then blocked="alignment static contract"; fi
else validation_skip "alignment static contract" "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then
  if ! validation_run "alignment regression" ./test-framework/governance/test_architecture_roadmap_alignment.sh; then blocked="alignment regression"; fi
else validation_skip "alignment regression" "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then
  if ! validation_run "complete test framework" ./test-framework/run_all.sh; then blocked="complete test framework"; fi
else validation_skip "complete test framework" "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then validation_run "repository validation" ./tools/validation/validate_repository.sh || true; else validation_skip "repository validation" "blocked by $blocked"; fi
validation_note "Scope: documentation, security, architecture, roadmap, gate, acceptance, and governance alignment only. No runtime or schema change; no Phase 1 Step 3 or Phase 2 acceptance; no production-readiness claim."
if validation_report_finish "$report_dir/final-report.txt" "$latest_report"; then exit 0; fi
exit 1
