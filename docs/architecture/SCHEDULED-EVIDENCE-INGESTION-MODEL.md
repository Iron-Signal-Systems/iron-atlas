# Scheduled Evidence Ingestion Model

## Status

Normative architecture alignment candidate. This document does not claim live collection is implemented or accepted.

## Purpose

Atlas must acquire evidence predictably without overwhelming infrastructure, creating overlapping collection storms, or treating collection success as proof that evidence is complete or current.

## Trigger classes

A collection has exactly one trigger: governed schedule, operator request, accepted change-validation workflow, accepted source event, or controlled recovery/backfill. Every request records trigger, requester or scheduler identity, source, profile, scope, creation time, and correlation identity.

## Governed schedules

A schedule declares source, collection profile, enabled state, timezone, recurrence, permitted start window, jitter, maximum duration, concurrency key, retry policy, freshness objective, blackout behavior, and evidence classification. Atlas does not invent a default production source, credential, scope, or profile.

## Admission and overlap

Atlas derives a concurrency key from source, profile, logical scope, and credential boundary.

- At most one active collection exists for a key unless safe parallelism is explicitly accepted.
- Equivalent queued requests may coalesce only with equivalent scope and security context.
- Operator work may supersede queued work but does not silently cancel active work.
- Missed schedules are classified rather than replayed without bound.
- Queue size, age, bytes, and workers are bounded.

## Collection lifecycle

```text
requested
→ admitted or rejected
→ started
→ source interaction
→ immutable bundle finalized
→ candidate validation
→ accepted, quarantined, rejected, or failed
```

A successful collection attempt is not accepted evidence. Only the atomic-acceptance boundary may publish it.

## Retry

Transient retry uses maximum attempts, bounded elapsed time, backoff, jitter, cancellation, source-specific rate limits, and circuit-open behavior. Authentication, authorization, malformed response, invariant, and deterministic parser failures do not retry indefinitely.

## Source protection

Collectors use least privilege, fixed command or query scope, response limits, read-only behavior unless a separate write boundary is accepted, and explicit device, tenant, VDOM, VRF, domain, or logical-scope identity. Credentials and source secrets never enter logs or retained evidence.

## Time semantics

Each bundle records schedule time, request time, collection start/end, source observation time, receipt time, validation time, and acceptance time separately. Missing source observation time is not silently replaced by receipt time.

## Configuration and runtime evidence

Configured state and observed runtime state are separate evidence classes with independent schedules and freshness. Disagreement remains explicit. Neither overwrites the other.

## Outage behavior

The prior accepted snapshot remains immutable and ages normally. The source becomes `STALE`, `UNAVAILABLE`, or `UNKNOWN`; answers carry that limitation; no empty snapshot replaces accepted evidence; recovery follows bounded retry and overlap rules.

## Required validation

Test schedule drift, jitter, duplicate triggers, overlap, cancellation, outage, authentication failure, oversized and partial responses, late completion, restart recovery, retry exhaustion, backlog saturation, and concurrent manual and scheduled requests.
