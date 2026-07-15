# Go PostgreSQL Runtime and Identity Context Boundary

## Status

Accepted as the non-production Phase 1 Step 2 runtime and identity-context boundary under tag `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`.

## Purpose

Phase 1 Step 2 connects the Go service to the accepted Phase 1 Step 1 PostgreSQL governance foundation without moving identity, approval, or history enforcement into the HTML5 client or ordinary handler code.

The Go runtime uses a replaceable domain service interface and a PostgreSQL adapter implemented with `pgx` and `pgxpool` under [ADR-0004](../decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md). PostgreSQL remains the authoritative enforcement point for durable actor status, authority, requester independence, approval concurrency, and append-only history.

## Runtime Modes

Iron Atlas supports two explicit change-store modes:

- `memory` — the default non-production demonstration store.
- `postgresql` — the persistent runtime adapter introduced in Step 2.

PostgreSQL mode requires `IRON_ATLAS_DATABASE_URL`. Real connection strings, passwords, certificates, and tokens are runtime secrets and are prohibited from Git.

Memory mode defaults the temporary development-header identity boundary on for the Phase 0 demonstration. PostgreSQL mode defaults that boundary off. It can be enabled only through an explicit `IRON_ATLAS_DEV_IDENTITY=true` setting for controlled local testing; it is never a production authentication mechanism.

## Dependency Direction

```text
HTTP and HTML5 handlers
        ↓
change.Service domain interface
        ↓
PostgreSQL change adapter
        ↓
least-privileged PostgreSQL pool
        ↓
accepted atlas service functions and projections
```

The domain and HTTP layers do not depend directly on `pgx`. The PostgreSQL-specific implementation remains isolated under `internal/database/postgresql` and `internal/change/postgresql`.

## Pool Boundary

The application pool:

- Uses only the `atlas_application` login identity.
- Has bounded minimum and maximum connections.
- Applies connection, statement, lock, and idle-in-transaction timeouts.
- Uses an explicit application name.
- Does not apply migrations.
- Does not own the database, schema, functions, or tables.
- Does not set an actor at session scope.
- Rejects connection options that attempt to preconfigure `atlas.actor_id`, including encoded options.
- Fails startup when PostgreSQL mode is selected and the dependency cannot be reached.

Production credential delivery and TLS certificate provisioning remain outside Step 2. The adapter accepts the externally delivered connection string but never logs it.

## Transaction-Scoped Actor Context

Every governed mutation follows this sequence:

1. Receive the authenticated actor from the server-side identity boundary.
2. Reject an empty actor before opening a transaction.
3. Begin a database transaction.
4. Call `set_config('atlas.actor_id', actor_id, true)` inside that transaction.
5. Execute an accepted security-definer service function.
6. Commit only when the complete operation succeeds.
7. Roll back on every failure path.

The third argument to `set_config` is `true`, making the value transaction-local. No handler, JSON field, form value, query parameter, or ordinary SQL statement can select a different acting identity.

The runtime must never execute either of the following outside a transaction:

```text
SET atlas.actor_id = ...
set_config('atlas.actor_id', ..., false)
```

## Pooled-Connection Identity Isolation

A pooled connection may serve many actors over time. Identity safety therefore depends on transaction-local context and complete transaction closure, not on the assumption that one connection belongs to one actor.

Step 2 tests:

- 500 sequential actor changes through one physical pool connection.
- 600 concurrent actor-bound transactions across a bounded pool.
- Actor-context removal after commit.
- Actor-context removal after rollback.
- A failed unauthorized transaction followed by a successful authorized transaction.
- Data rollback together with identity-context rollback.

These tests prove the implemented boundary under a disposable PostgreSQL cluster. They do not prove production identity-provider correctness.

## Change Persistence

The PostgreSQL adapter uses only the accepted runtime surface:

- `atlas.create_change_request(text, text, integer)`
- `atlas.record_approval(text, text, text)`
- `atlas.change_request`
- `atlas.change_approval_summary`

It does not insert, update, or delete governed tables directly.

The database derives the actor from transaction context. Roles supplied in development HTTP headers are used only by the in-process demonstration authorization layer and are not trusted as durable PostgreSQL authority records.

## Health and Readiness

- `/healthz` is a liveness check and remains independent of PostgreSQL state.
- `/readyz` checks the selected store dependency.
- PostgreSQL failure returns HTTP 503 from readiness.
- Dependency failures in change reads or writes fail closed and do not expose connection details to the client.

## Failure Boundary

The runtime fails closed for:

- Missing database URL in PostgreSQL mode
- Invalid pool limits
- Startup connection failure
- Missing authenticated actor
- Transaction begin failure
- Actor-context binding failure
- Governed function failure
- Commit failure
- Database read failure
- Readiness dependency failure

## Explicit Exclusions

Phase 1 Step 2 does not establish:

- Production authentication
- NPS, LDAP, SAML, OIDC, or other identity-provider integration
- Production database credentials or secret rotation
- Production TLS certificate issuance
- Database backup or restoration
- High availability or failover
- Protected evidence persistence
- Live device collection
- Production deployment acceptance

## Portable Validation Boundary

The Step 2 runtime is not eligible for acceptance from a workstation-only validation. The exact pushed commit must pass the Step 2 phase gate from a clean clone of the canonical GitHub repository after validating the declared external toolchain and verifying pinned Go modules. Sanitized retained evidence is committed under `validation/evidence/`.
