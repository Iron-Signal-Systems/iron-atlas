# Phase Gates

Phase gates are checkpoint-specific. Historical gates are preserved and re-run at their exact accepted commit when later work may change repository assumptions.

## Gates

- `validate_phase0_step1.sh` — Phase 0 implementation candidate.
- `validate_phase0_acceptance.sh` — accepted Phase 0 non-production boundary.
- `validate_phase1_step1.sh` — Phase 1 Step 1 PostgreSQL migration and governed-identity candidate; revalidates the accepted Phase 0 tag in an isolated local clone.
- `validate_phase1_step1_acceptance.sh` — accepted Phase 1 Step 1 PostgreSQL governance boundary.
- `validate_phase1_step2.sh` — Phase 1 Step 2 Go PostgreSQL runtime and transaction-local identity-context candidate; revalidates accepted Step 1 in an isolated local clone.
- `validate_phase1_step2_acceptance.sh` — formal Phase 1 Step 2 acceptance boundary; verifies the exact implementation/evidence chain, retained evidence integrity, synchronized documentation, and the still-passing implementation gate.
- `validate_phase1_step3_contract.sh` — Phase 1 Step 3 phase-entry contract; verifies requirements, architecture, traceability, testing, acceptance-template, accepted predecessor, and repository synchronization. It is not the final executable Step 3 gate.

- `validate_phase1_step3_authentication_foundation.sh` — first Phase 1 Step 3 implementation gate; verifies typed authentication modes, development-header isolation, immutable request identity, production fail-closed behavior, future adapter/resolver seams, targeted race tests, the accepted contract predecessor, and complete repository validation. It is not Step 3 acceptance.


Step 2 uses the tested `isolated_gate_revalidate` helper. The helper preserves the predecessor validator exit status after temporary-clone cleanup, so a failed historical gate cannot be reported as a passing revalidation.

## Canonical Repository Requirement

A local phase gate is necessary but not sufficient for acceptance. The exact pushed commit must also pass the applicable gate through `tools/validation/verify_canonical_clone.sh`. Retained output is recorded with `tools/validation/record_validation_evidence.sh` and committed below `validation/evidence/`.

- `validate_phase1_step3_governed_actor_resolution.sh` — second Phase 1
  Step 3 implementation gate; verifies the least-privileged PostgreSQL resolver
  function, explicit role mapping, fail-closed governed-state behavior,
  disposable database and race tests, accepted authentication-foundation
  predecessor, and complete repository validation. It is not Step 3
  acceptance.

- `validate_phase1_step3_oidc_id_token_verification.sh` — third Phase 1
  Step 3 implementation gate; verifies pinned OIDC dependencies, exact HTTPS
  discovery, JWKS signature verification, issuer, audience, authorized party,
  nonce, time, stable-subject, duplicate sensitive-field, key-rotation, outage,
  race, concurrency, documentation, manifests, the accepted governed
  actor-resolution predecessor, and complete repository validation. It is not
  Step 3 acceptance.

- `validate_phase1_step3_oidc_authorization_code_pkce.sh` — fourth Phase 1
  Step 3 implementation gate; verifies state, nonce, PKCE S256, SHA-256
  state-digest storage, atomic one-time consumption, exact redirect and token
  endpoint binding, bounded token exchange, replay, outage, race, concurrency,
  redaction, predecessor revalidation, and complete repository validation. It
  is not Step 3 acceptance.
