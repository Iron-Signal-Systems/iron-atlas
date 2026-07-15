# Phase Gates

Phase gates are checkpoint-specific.

Do not weaken a historical gate merely so it passes against a later working
tree. Revalidate historical boundaries at their exact accepted commit when
absence of later artifacts is part of the gate.

## Current Gates

- `validate_phase0_step1.sh` validates the Phase 0 implementation candidate.
- `validate_phase0_acceptance.sh` validates the synchronized Phase 0
  non-production acceptance boundary.
