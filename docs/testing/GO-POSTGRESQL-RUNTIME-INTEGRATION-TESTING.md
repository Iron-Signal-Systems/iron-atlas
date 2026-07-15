# Go PostgreSQL Runtime Integration Testing

## Status

Accepted as the Phase 1 Step 2 disposable runtime-integration test boundary under tag `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`.

## Disposable Boundary

The existing disposable PostgreSQL runner now applies the accepted migrations and runtime grants, executes the SQL governance suite, and then runs the Go PostgreSQL integration tests against the same temporary Unix-socket-only database.

The test cluster:

- Uses `initdb --auth=trust` only inside a private temporary directory.
- Listens on no TCP address.
- Uses a temporary Unix socket.
- Is removed after each run.
- Does not initialize or modify `/var/lib/postgres/data`.
- Does not require the host PostgreSQL service to be enabled.

## Go Integration Coverage

The Step 2 suite proves:

- The application pool connects using the least-privileged runtime role.
- PostgreSQL mode defaults development identity headers off.
- Encoded connection options cannot preconfigure the acting identity.
- The pool reports readiness while the database is available.
- Actor context is visible inside the governed transaction.
- Actor context is absent after commit.
- Actor context is absent after rollback.
- A governed write is rolled back when the callback fails.
- 500 sequential transactions do not leak actor identity through one pooled connection.
- Eight concurrent workers complete 75 actor-bound transactions each without cross-actor contamination.
- The Go change adapter creates a change through `atlas.create_change_request`.
- Requester self-approval remains rejected by PostgreSQL.
- An independent approver can approve through `atlas.record_approval`.
- A failed unauthorized operation does not contaminate the next pooled transaction.
- A failing or missing isolated predecessor validator causes the Step 2 phase gate to fail rather than being masked by temporary-clone cleanup.

## Test Commands

```bash
./test-framework/run_all.sh
./tools/validation/phase-gates/validate_phase1_step2.sh
```

The integration packages are selected with the `integration` build tag by the disposable database runner. Ordinary `go test ./...` does not require a running database.

## Interpretation

Passing tests prove the tested Go pool, transaction, identity-context, service-function, rollback, readiness, and disposable-database boundary. They do not prove production credentials, authentication, network encryption, recovery, high availability, or production readiness.

## Canonical Reproduction

The complete integration boundary must run from a clean canonical GitHub clone with the requirements in `validation/toolchain-requirements.json`. The disposable PostgreSQL runner creates and destroys its own cluster; a workstation PostgreSQL service is not part of the proof. Results deliberately retained for Step 2 acceptance are recorded, sanitized, checksummed, validated, and committed under `validation/evidence/`.
