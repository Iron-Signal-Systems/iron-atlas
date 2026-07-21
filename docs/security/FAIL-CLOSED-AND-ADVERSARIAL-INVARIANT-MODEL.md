# Fail-Closed and Adversarial Invariant Model

## Status

Normative security alignment candidate.

## Core invariants

1. Untrusted input never selects an actor, authority, credential, source scope, accepted snapshot, or write target.
2. Missing, malformed, duplicate, ambiguous, stale, conflicting, or excessive security state fails closed.
3. A UI, adapter, provider claim, parser, collector, or integration cannot grant Atlas authority.
4. An unaccepted candidate cannot become current through partial progress, restart, retry, or concurrency.
5. Failure never erases the prior accepted state.
6. Secrets and protected evidence do not enter logs, manifests, retained test output, URLs, browser storage, or generic errors.
7. Resource exhaustion is bounded and attributable.
8. Development mechanisms cannot activate in production mode.
9. Optional modules cannot expand core privilege.
10. Unknown is represented as unknown, never as absent, safe, or permitted.

## Adversarial classes

Validation includes spoofing, identity confusion, replay, duplicate consumption, cardinality ambiguity, malformed and excessive structured data, injection, downgrade, stale policy or actor state, concurrent disablement and revocation, dependency outage, clock manipulation, queue saturation, disk/memory/CPU/connection pressure, crash and restart, signer/digest/audience/provenance substitution, cross-scope confusion, secret reflection, and confused-deputy behavior across Atlas, IFI, OIDC, PostgreSQL, and source systems.

## Enforcement points

Hostile classes are tested at every relevant enforcement point: admission, protocol verification, actor resolution, session creation, authorization, database routine, collection, bundle validation, parser, normalization, candidate acceptance, correlation, query, export, and administration.

## Error behavior

External errors are stable and non-sensitive. Internal retained evidence may contain bounded identity, stage, policy version, and classification but not secrets or unrestricted input.

Error handling must not retry unauthorized behavior, convert ambiguity into a default, leak protected existence, return partial accepted results, mask cleanup failure, convert validator failure into success, or continue after integrity becomes unknown.

## Resource observations

Correctness and resource observations remain separate. Atlas records time, CPU, memory, I/O, database, WAL, storage growth, and host/tool fingerprints. Security limits remain enforceable before performance budgets are statistically justified.

## Required gate evidence

Every gate records exact predecessor/candidate, threat classes, enforcement points, positive and negative cases, concurrency and race results, outage behavior, resource observations, redaction, toolchain, limitations, canonical clean-clone result, and explicit nonclaims.

## Explicit nonclaims

Passing a bounded adversarial gate does not prove absence of defects, independent review, complete vendor coverage, or production readiness.
