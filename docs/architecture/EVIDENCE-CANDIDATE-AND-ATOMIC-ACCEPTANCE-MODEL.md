# Evidence Candidate and Atomic-Acceptance Model

## Status

Normative architecture alignment candidate.

## Purpose

Atlas must never expose a partially parsed, partially normalized, internally inconsistent, or concurrently changing evidence set as accepted truth.

## State model

```text
RECEIVED
→ QUARANTINED or STAGED
→ CANDIDATE
→ VALIDATING
→ ACCEPTED or REJECTED
→ SUPERSEDED
```

A failed candidate never mutates the prior accepted state.

## Candidate identity

A candidate records source identity, logical scope, collection profile, bundle schema, observation and collection times, content digests, parser and normalizer versions, policy version, candidate ID, and parent snapshot when applicable. Identity is immutable after finalization.

## Validation

Before acceptance Atlas verifies structure, schema, digest and signature where required, source and scope, resource limits, record status, replay and duplication, parser and normalization success, referential consistency, configured-versus-observed distinction, protected-field rules, exact tool and policy versions, and explicit unknown, incomplete, ambiguous, and conflicting results.

Validation may classify findings but may not silently repair evidence in a way that destroys lineage.

## Atomic publication

Acceptance is one transaction or equivalent atomic operation that records the candidate and result, creates the immutable accepted snapshot, publishes all normalized records together, advances the current pointer for exactly one source and scope, preserves the predecessor, and records responsible actor, service, policy, and time.

Readers see either the complete previous accepted state or the complete new accepted state, never a mixture.

## Concurrency

- One lineage position has one winning acceptance.
- Stale candidates cannot overwrite newer accepted state.
- Equivalent candidates may deduplicate by digest.
- Different concurrent candidates remain distinct and use deterministic ordering or governed resolution.
- Parallel validation keeps all outputs isolated.
- Acceptance retries are idempotent.

## Rejection and quarantine

Rejected or quarantined evidence retains bounded metadata identifying classification, source, scope, candidate, validation version, failure stage, sanitized reason, retry eligibility, and retention policy. Rejected content is never queried as accepted evidence.

## Supersession and rollback

Rollback does not rewrite history. It creates a governed current-state decision that selects a previously accepted snapshot or a newly reconstructed candidate, with reason and approval evidence. Superseded snapshots remain available for historical query, comparison, incident reconstruction, audit, reproduction, and rollback planning.

## Cross-source acceptance

Cisco, FortiGate, BloodHound-derived, IFI, documentation, and monitoring snapshots are accepted independently. A correlation build records the exact accepted inputs it consumed. One source update does not silently reinterpret another source's history.

## Required validation

Test transaction rollback, interruption, concurrent acceptance, stale candidate rejection, duplicate digest handling, malformed lineage, partial-publication prevention, reader consistency, supersession, rollback, corruption, and deterministic rebuild.
