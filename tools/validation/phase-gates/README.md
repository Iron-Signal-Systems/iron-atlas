# Phase Gates

Phase gates are checkpoint-specific. Historical gates are preserved and re-run at their exact accepted commit when later work may change repository assumptions.

## Gates

- `validate_phase0_step1.sh` — Phase 0 implementation candidate.
- `validate_phase0_acceptance.sh` — accepted Phase 0 non-production boundary.
- `validate_phase1_step1.sh` — Phase 1 Step 1 PostgreSQL migration and governed-identity candidate; revalidates the accepted Phase 0 tag in an isolated local clone.
- `validate_phase1_step1_acceptance.sh` — accepted Phase 1 Step 1 PostgreSQL governance boundary.
- `validate_phase1_step2.sh` — Phase 1 Step 2 Go PostgreSQL runtime and transaction-local identity-context candidate; revalidates accepted Step 1 in an isolated local clone.

Step 2 uses the tested `isolated_gate_revalidate` helper. The helper preserves the predecessor validator exit status after temporary-clone cleanup, so a failed historical gate cannot be reported as a passing revalidation.

## Canonical Repository Requirement

A local phase gate is necessary but not sufficient for acceptance. The exact pushed commit must also pass the applicable gate through `tools/validation/verify_canonical_clone.sh`. Retained output is recorded with `tools/validation/record_validation_evidence.sh` and committed below `validation/evidence/`.
