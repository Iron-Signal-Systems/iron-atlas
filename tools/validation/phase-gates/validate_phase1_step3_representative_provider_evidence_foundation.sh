#!/usr/bin/env bash
set -Eeuo pipefail
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"; cd "$root"
source "$root/tools/validation/lib/reporting.sh"
source "$root/tools/validation/lib/isolated_gate_revalidation.sh"
base=e7824049852855f15d26686600fc42802b8a38ff
results="$root/test-framework/test-results/phase-gates"; id="$(date -u +%Y%m%dT%H%M%SZ)-$$"
dir="$results/phase1-step3-representative-provider-evidence-foundation-$id"; latest="$results/phase1-step3-representative-provider-evidence-foundation-latest-report.txt"
validation_report_init 'Iron Atlas Phase 1 Step 3 representative-provider evidence foundation' "$dir"
revalidate(){ isolated_gate_revalidate "$root" "$base" tools/validation/phase-gates/validate_phase1_step3_provider_neutral_assurance_evidence.sh; }
blocked=''
validation_run 'SSH-signed provider-neutral assurance-evidence boundary is an ancestor' git merge-base --is-ancestor "$base" HEAD || blocked='signed base ancestry'
if [[ -z "$blocked" ]]; then validation_run 'provider-neutral assurance-evidence boundary remains valid' revalidate || blocked='provider-neutral boundary'; else validation_skip 'provider-neutral assurance-evidence boundary remains valid' "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then validation_run 'representative-provider evidence-foundation static contract' python3 tools/validation/validate_phase1_step3_representative_provider_evidence_foundation.py || blocked='static contract'; else validation_skip 'representative-provider evidence-foundation static contract' "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then validation_run 'representative-provider evidence-foundation regression' ./test-framework/authentication/test_phase1_step3_representative_provider_evidence_foundation.sh || blocked='regression'; else validation_skip 'representative-provider evidence-foundation regression' "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then validation_run 'complete test framework' ./test-framework/run_all.sh || blocked='complete framework'; else validation_skip 'complete test framework' "blocked by $blocked"; fi
if [[ -z "$blocked" ]]; then validation_run 'repository validation' ./tools/validation/validate_repository.sh || true; else validation_skip 'repository validation' "blocked by $blocked"; fi
validation_note 'Scope: bounded evidence contract and sanitized observation boundary only. This establishes strict observation-only bundles, literal assurance-claim preservation, digest and path binding, deterministic validation, and secret rejection. It does not establish compatibility with a named provider, live provider behavior, provider-specific semantic mapping, session lifecycle completion, CSRF, trusted proxies, emergency access, production wiring, formal Phase 1 Step 3 acceptance, or production readiness.'
validation_report_finish "$dir/final-report.txt" "$latest"
