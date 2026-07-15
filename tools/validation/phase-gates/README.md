# Phase Gates

Phase gates are checkpoint-specific. Historical gates are preserved and re-run at their exact accepted commit when later work may change repository assumptions.

## Gates

- `validate_phase0_step1.sh` — Phase 0 implementation candidate.
- `validate_phase0_acceptance.sh` — accepted Phase 0 non-production boundary.
- `validate_phase1_step1.sh` — Phase 1 Step 1 PostgreSQL migration and governed-identity candidate; revalidates the accepted Phase 0 tag in an isolated local clone.
- `validate_phase1_step1_acceptance.sh` — accepted Phase 1 Step 1 non-production PostgreSQL governance boundary; re-runs the implementation gate.
