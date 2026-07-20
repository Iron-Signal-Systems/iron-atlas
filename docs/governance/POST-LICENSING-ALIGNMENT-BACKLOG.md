# Post-Licensing Architecture and Roadmap Alignment Backlog

> Status: Governed next-work register
>
> This document records deferred discussion decisions. It does not implement
> them and does not expand the Business Source License transition PR into a
> runtime or architecture rewrite.

## Sequence

After the BSL transition is merged and closed with a signed post-merge
boundary, create a dedicated branch and PR for architecture, roadmap,
requirements, testing, and phase-gate alignment.

## Authentication direction correction

Iron Atlas is not intended to become a local password or TOTP identity
provider.

The alignment change must:

- remove local TOTP enrollment, secret generation, QR enrollment, and local
  recovery-code implementation from the required roadmap;
- state that user primary authentication and MFA occur at an approved external
  OIDC identity provider;
- retain Atlas validation of `acr`, `amr`, `auth_time`, authentication age,
  step-up requirements, privileged-role requirements, and phishing-resistant
  assurance;
- add representative-provider compatibility testing;
- specify session rotation, logout, administrative revocation, emergency
  access, CSRF, trusted-proxy, and production-wiring boundaries;
- prohibit Atlas from storing user passwords or provider TOTP seeds; and
- preserve the already merged authentication-assurance implementation and its
  historical validation claims.

## Required architecture artifacts

The alignment PR must create and synchronize:

```text
docs/architecture/module-runtime-and-failure-containment-model.md
docs/architecture/scheduled-evidence-ingestion-model.md
docs/architecture/monitoring-alerting-and-evidence-freshness-model.md
docs/architecture/evidence-candidate-and-atomic-acceptance-model.md
docs/architecture/atlas-ifi-snapshot-integration-contract.md
docs/security/fail-closed-and-adversarial-invariant-model.md
docs/security/mfa-and-authentication-assurance-requirements.md
```

File naming may be normalized to the repository's uppercase convention, but the
meaning and coverage must remain.

## Module runtime and failure containment

Define:

- module-process and goroutine failure boundaries;
- cancellation and timeout propagation;
- restart behavior and degraded-state reporting;
- bounded queues and backpressure;
- evidence-source isolation;
- resource ceilings and hostile-input containment; and
- the conditions under which one adapter may fail without taking down Atlas.

## Scheduled ingestion and evidence freshness

Define:

- scheduled and operator-triggered collection;
- source-specific freshness and staleness;
- retry, jitter, backoff, and outage behavior;
- duplicate and overlapping collection control;
- snapshot identity and lineage;
- reconciliation between configuration and runtime evidence; and
- explicit unknown, stale, incomplete, and conflicting states.

## Candidate and atomic acceptance model

Define:

- immutable candidate evidence sets;
- validation before acceptance;
- all-or-nothing publication of accepted state;
- rollback and supersession;
- exact evidence and toolchain lineage;
- separation of current candidate, accepted state, and historical state; and
- concurrency behavior when collection or validation overlaps.

## Atlas–Iron File Intelligence boundary

Preserve this permanent direction:

```text
IFI authoritative state
→ immutable signed purpose-limited evidence snapshot
→ Atlas IFI import adapter
→ validated Atlas correlation records
```

Atlas must not:

- query the IFI PostgreSQL database directly;
- share IFI schemas, service credentials, or unrestricted storage;
- assume IFI availability;
- copy unrestricted file content or sensitive metadata;
- write back to IFI;
- make IFI evidence-lifecycle decisions; or
- inherit IFI's collection or forensic authority.

IFI remains authoritative for file intelligence. Atlas remains authoritative
for Atlas-side validation, infrastructure reachability, correlation,
explanation, change impact, and leadership-facing conclusions.

## Phase and acceptance reconciliation

The alignment PR must:

- preserve all historical accepted tags, gates, records, and claims;
- distinguish accepted, merged implementation, active candidate, and planned
  work;
- correct future phase numbering and sequencing without rewriting history;
- identify historical gates that remain valid;
- identify any gate requiring exact-commit isolated revalidation;
- keep PR #15 authentication assurance as an implementation checkpoint rather
  than formal Step 3 acceptance; and
- synchronize README, architecture, requirements, testing, roadmap, phase-gate
  plan, acceptance template, and changelog in the same change set.

## Work that follows alignment

Only after the alignment boundary is accepted should implementation resume with
bounded work for:

1. external-provider compatibility;
2. completed server-side session lifecycle;
3. logout and administrative revocation;
4. CSRF enforcement;
5. trusted-proxy enforcement;
6. production application wiring;
7. emergency and recovery access controls; and
8. formal Phase 1 Step 3 acceptance preparation.
