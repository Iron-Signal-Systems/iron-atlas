# Module Runtime and Failure-Containment Model

## Status

Normative architecture alignment candidate. This contract does not claim that every boundary is implemented.

## Purpose

A malformed source, slow collector, exhausted parser, failed integration, or defective optional module must not become loss of the entire Atlas service, corruption of accepted evidence, or fabricated certainty.

## Runtime ownership

Atlas distinguishes process, service, module, work-item, and goroutine boundaries. A goroutine is an implementation mechanism, never an authority or lifecycle boundary by itself. Every goroutine must belong to an owned service or work item with cancellation, completion, and cleanup rules.

```text
external source
→ bounded adapter
→ immutable evidence candidate
→ validation and normalization
→ canonical Atlas contracts
→ correlation and answer services
```

Core services do not depend on vendor-specific packages or optional integrations. A module may be absent, disabled, degraded, stale, or failed without changing the meaning of other evidence.

## Containment requirements

Each module must:

- own its context and cancel all descendants during shutdown;
- set operation, network, parse, persistence, and publication deadlines;
- use bounded queues, bytes, workers, retries, records, and findings;
- reject or defer work at admission limits instead of growing without bound;
- avoid mutable package-global operational state;
- recover a panic only where the entire work item can be invalidated;
- discard all unaccepted partial output after failure;
- preserve the last accepted state when a successor candidate fails;
- expose degraded state without treating missing evidence as proof of safety;
- separate health, readiness, and evidence freshness; and
- keep credentials, private evidence, and unrestricted source fragments out of logs.

## Failure classes

Failures are classified as input rejection, source unavailable, dependency unavailable, timeout, cancellation, resource limit, adapter defect, validation failure, conflicting evidence, persistence failure, publication failure, or internal invariant violation.

The classification must be stable, bounded, sanitized, attributable to one module and work item, and suitable for retry policy.

## Panic and invariant behavior

Panic recovery is allowed only at a top-level work-item or module supervisor that can mark the work failed, discard output, release leases, retain sanitized diagnostics, and avoid reusing questionable mutable state. A core integrity violation may terminate the process so a service manager restarts a clean instance. Atlas must not continue in an unknown integrity state merely to appear available.

## Backpressure and retry

Every producer-consumer edge declares maximum queued items, bytes, workers, admission wait, overflow behavior, and cancellation behavior. Silent dropping and unbounded buffering are prohibited.

Retry is allowed only for classified transient failure. It is bounded by attempts, elapsed time, backoff, jitter, and a circuit-open state. Repeated deterministic failure stops automatic retry and requires operator action.

## Degraded-state model

Atlas uses at least `HEALTHY`, `DEGRADED`, `UNAVAILABLE`, `STALE`, `BLOCKED`, `REJECTED`, and `UNKNOWN`. Overall process readiness must not hide a failed evidence source. Queries carry source-specific state and freshness.

## Security and privilege

A module receives only the credentials, network reachability, storage, database routines, and authorities required for its role. Optional modules do not inherit Atlas administrative authority, unrestricted evidence, another module's credentials, source-system write access, or authority to accept their own output.

## Required validation

Implementation gates test panic containment, cancellation, saturation, resource exhaustion, retry exhaustion, dependency outage, partial-output discard, concurrent shutdown and completion, stale-state reporting, module disablement, one-adapter failure while other capabilities remain available, and race-free cleanup.

## Explicit nonclaims

This contract does not require microservices or containers. A modular monolith remains acceptable when ownership, privilege, resource, and failure boundaries are explicit and tested.
