# Architecture and Roadmap Alignment Record

## Status

Implementation candidate based on signed BUSL boundary:

```text
a0ab1ad19cf48ba11d97b3a9e87acd7b68e1eb60
establish SSH-signed post-merge validation boundary for Business Source License 1.1
```

## Decision

This alignment converts the post-licensing discussion register into normative architecture, security, authentication, roadmap, gate, testing, acceptance, and repository-governance direction without changing runtime code or database schema.

## Historical preservation

- Phase 0 remains accepted under its existing tag and records.
- Phase 1 Steps 1 and 2 remain accepted under their existing tags and records.
- Existing Phase 1 Step 3 validators remain historical implementation checkpoints.
- PR #15 authentication assurance remains a merged implementation checkpoint, not formal Step 3 acceptance.
- No historical gate, tag, record, result, or nonclaim is renamed or weakened.

## Authentication correction

Atlas relies on approved external OIDC providers for primary authentication and MFA. Atlas validates provider assurance but does not own user passwords, provider TOTP seeds, QR enrollment, or ordinary recovery codes. The planned Atlas-local TOTP checkpoint is removed from the required sequence.

## Architecture artifacts

This candidate adds module runtime/failure containment, scheduled ingestion, evidence freshness, atomic acceptance, Atlas–IFI signed snapshot integration, fail-closed adversarial invariants, provider-owned MFA requirements, and signed candidate/post-merge trust governance.

## Phase 1 Step 3 successor order

1. representative external-provider compatibility;
2. session rotation, expiry, cleanup, logout, and administrative revocation;
3. CSRF enforcement;
4. trusted-proxy and transport enforcement;
5. production authentication wiring;
6. emergency and recovery access controls;
7. Step 3 integration; and
8. formal Step 3 acceptance preparation.

## Evidence-platform direction

Phase 2 remains the vendor-independent evidence intake, protection, candidate, atomic acceptance, parser isolation, recovery, and hostile-evidence foundation. Cisco, FortiGate, BloodHound-derived, and IFI evidence sources build on that boundary.

## Nonclaims

This alignment does not implement the contracts, accept Phase 1 Step 3, accept Phase 2, establish live collection, create IFI integration, or establish production readiness.

## Historical checkpoint revalidation

After this alignment changes successor documentation, the frozen Phase 1 Step 3
authentication-assurance checkpoint is revalidated in an isolated local clone
at the exact signed boundary:

```text
cc93fdd2311ca188ad03b0bd94293156ff243973
```

The historical validator, regression, phase gate, implementation, and original
documentation are executed unchanged at that commit. The aligned successor tree
is validated by the architecture-and-roadmap alignment gate.

This preserves historical evidence without weakening frozen validators or
forcing successor documentation to repeat superseded Atlas-local TOTP roadmap
commitments.
