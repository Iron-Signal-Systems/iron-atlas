# Disposable PostgreSQL Testing

## Purpose

Phase 1 database tests create a temporary PostgreSQL cluster for each run. They never initialize or modify `/var/lib/postgres/data` and never require a persistent Iron Atlas database.

## Prerequisites

On Arch Linux:

```bash
sudo pacman -S --needed postgresql
```

The test discovers PostgreSQL binaries through `PATH` or `pg_config --bindir`.

## Test Lifecycle

`test-framework/database/run_disposable_postgres.sh`:

1. Creates a temporary data directory and Unix-socket directory.
2. Initializes a trust-authenticated local-only test cluster.
3. Starts PostgreSQL without TCP listening.
4. Creates the development role topology and disposable database.
5. Applies the ordered migrations.
6. Applies runtime grants.
7. Executes correctness, security, idempotency, and concurrency tests through independent `psql` connections.
8. Stops PostgreSQL and removes the temporary cluster through a shell trap.

## Required Proofs

- All manifest migrations apply in order.
- A second migration run is idempotent.
- Migration history contains the expected hashes.
- The application role cannot modify governed tables directly.
- Requesters cannot approve their own changes.
- Unauthorized actors cannot approve changes.
- Independent authorized actors can satisfy a two-person requirement.
- Duplicate concurrent approval from one actor produces only one current approval.
- Append-only records reject update and delete.
- Runtime grants do not grant DDL or ownership authority.

## Correctness and Resource Separation

Correctness results are pass or fail. PostgreSQL version, host fingerprint, elapsed time, and temporary database size are observations and do not become performance thresholds during this step.
