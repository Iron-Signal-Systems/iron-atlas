# ADR-0004 — pgx PostgreSQL Runtime Driver

## Status

Phase 1 Step 2 implementation candidate.

## Context

Iron Atlas requires a production-capable Go PostgreSQL driver and bounded
connection pool for runtime persistence, transaction control, context-aware
cancellation, PostgreSQL-specific types, and disposable integration testing.
The Go standard library defines `database/sql`, but it does not include a
PostgreSQL wire-protocol implementation.

The project preference remains to minimize dependencies and use the standard
library where it provides the required capability. A database driver is a
justified exception because implementing and maintaining the PostgreSQL wire
protocol, authentication methods, cancellation behavior, type system, and pool
semantics inside Iron Atlas would create greater security and reliability risk.

## Decision

Use `github.com/jackc/pgx/v5` and its `pgxpool` package for the isolated
PostgreSQL runtime adapter.

The dependency shall be:

- Pinned in `go.mod` and verified by `go.sum`.
- Isolated under `internal/database/postgresql` and
  `internal/change/postgresql`.
- Hidden behind Iron Atlas domain interfaces.
- Updated only through a reviewed change with tests and documentation.
- Excluded from vendor-neutral domain packages and HTML handlers.
- Used with explicit pool limits, timeouts, transactions, and context
  cancellation.

Phase 1 Step 2 pins `pgx/v5` version `v5.10.0` and sets the module language
baseline to Go 1.25.

## Security Consequences

- Iron Atlas does not implement a custom database wire protocol.
- The dependency becomes part of the software supply chain and must be covered
  by future SBOM, vulnerability, provenance, and update controls.
- Database URLs and credentials remain runtime secrets and are not embedded in
  the binary or repository.
- The driver does not receive migration or ownership authority.
- Acting identity is set only through transaction-local context and is tested
  against pooled-connection leakage.

## Alternatives Considered

### Implement a PostgreSQL driver in Iron Atlas

Rejected because it would duplicate a complex security-sensitive protocol and
substantially increase maintenance and validation burden.

### Use `database/sql` with another PostgreSQL driver

Not selected for this phase. `pgx` provides direct PostgreSQL behavior and a
native bounded pool while remaining isolated behind the same domain boundary.
The adapter may still be replaced later without changing the domain model.

### Shell out to `psql`

Rejected for runtime persistence because it complicates transaction ownership,
cancellation, pooling, error handling, and secret exposure. `psql` remains
appropriate for migration and disposable-test orchestration.
