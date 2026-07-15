# Phase 1 Step 1 Acceptance Record

## Record Status

- Decision: Accepted as a non-production PostgreSQL governance foundation
- Product: Iron Atlas
- Repository: `Iron-Signal-Systems/iron-atlas`
- Branch: `dev`
- Acceptance date: `2026-07-15T09:05:05Z`
- Accepted tag: `phase-1-step-1-postgresql-governance-foundation-complete-v1`

## Predecessor

- Accepted Phase 0 tag: `phase-0-repository-and-executable-baseline-complete-v1`
- Accepted Phase 0 commit: `6b24494e1e443eb8175af204e5a2e8ff66b2a2c6`
- Isolated predecessor revalidation: PASS

## Candidate

- Candidate implementation commit: `f76d5d0be180a2faa9bfc8a229ba63b686117b86`
- Candidate short commit: `f76d5d0be180`
- Candidate Git archive SHA-256: `2a2d0a5e9749aa6c152393069aad08cf6ca6efb1f3a061b3fec8eadc8c766df8`
- Archive command: `git archive --format=tar --prefix=iron-atlas-phase1-step1/ f76d5d0be180a2faa9bfc8a229ba63b686117b86`
- Go version: `go version go1.26.5-X:nodwarf5 linux/amd64`
- Python version: `Python 3.14.6`
- PostgreSQL server: `postgres (PostgreSQL) 18.4`
- PostgreSQL client: `psql (PostgreSQL) 18.4`
- Host fingerprint: `Arch Linux; Linux 7.1.3-arch1-3 x86_64 GNU/Linux; CPUs=2; MemTotal=3866668 KiB`

## Accepted Scope

Phase 1 Step 1 accepts the following non-production development boundary:

- Manifest-driven and checksum-governed PostgreSQL migrations.
- Ordered migration application and idempotency protection.
- Separate database owner, schema owner, migrator, application, read-only,
  auditor, and development test-runner roles.
- Governed actor and external-identity persistence.
- Role binding and authority persistence.
- Governed change-request, approval, decision, status-history, and audit records.
- Database-enforced requester and approver independence.
- Independent-session approval concurrency protection.
- Append-only migration, approval, decision, status-history, and audit records.
- Least-privilege runtime function boundary.
- Disposable Unix-socket-only PostgreSQL correctness and security tests.
- Isolated revalidation of the frozen Phase 0 predecessor.

## Validation Evidence

- Phase 1 Step 1 implementation gate: PASS
- Phase-gate checks: 8 PASS, 0 FAIL
- Repository validation: PASS
- Repository checks: 21 PASS, 0 FAIL
- Complete test framework: PASS
- Test-framework checks: 7 PASS, 0 FAIL
- Go formatting: PASS
- Go vet: PASS
- Go tests: PASS
- Race-enabled Go tests: PASS
- Migration static validation: PASS
- Database security static validation: PASS
- Six migrations applied and recorded: PASS
- Migration checksum enforcement: PASS
- Migration idempotency: PASS
- Requester self-approval denial: PASS
- Independent-session approval concurrency: PASS
- Append-only protections: PASS
- Internal actor helper remains outside application API: PASS
- Governed change API resolves actor context: PASS
- Correctness result: PASS
- Resource observation: RECORDED_BY_DATABASE_TEST
- Performance thresholds: NOT_EVALUATED
- Pre-acceptance transcript: `/tmp/iron-atlas-phase1-step1-pre-acceptance-validation.log` on the acceptance host

## Review and Approval

- Requester and implementer: John Wood (kb2vhn@gmail.com)
- Approval authority: John Wood, repository owner
- Independent human reviewer: Not assigned for this non-production single-maintainer development baseline
- Conflicts checked: Temporary single-maintainer development exception recorded below
- Operational two-person approval: Not exercised by this development acceptance

### Temporary Single-Maintainer Development Exception

Phase 1 Step 1 establishes and tests the database controls required to enforce
independent operational approval, but it does not connect to production identity,
production credentials, live collectors, or operational infrastructure.

For this reason, the repository owner accepts this development boundary under a
temporary single-maintainer exception.

This exception:

- Does not authorize a production deployment.
- Does not authorize production credentials or authentication.
- Does not authorize live infrastructure collection.
- Does not authorize an operational infrastructure change.
- Does not permit a requester to approve their own operational change.
- Does not weaken the documented two-person operational change-control contract.
- Expires before production-boundary acceptance or operational use.
- Must not be used as precedent for bypassing independent operational approval.

## Explicit Exclusions

The following capabilities remain outside the accepted Phase 1 Step 1 boundary:

- Production authentication
- Production credentials and secret delivery
- Go service PostgreSQL integration
- Production connection pooling
- Transaction-scoped runtime identity context from the Go service
- Live infrastructure collection
- Protected evidence intake and storage
- Production database installation and service configuration
- Production database backup and restoration
- Production database high availability
- Production monitoring and alert delivery
- Complete firewall or Cisco semantic analysis
- Production performance budgets
- Production readiness

## Security Assumptions

- The disposable database proves only the tested SQL and privilege boundary.
- Development database roles and credentials are not production credentials.
- The application role does not own the database, schema, or governed tables.
- The application role cannot apply migrations or directly execute internal helper functions.
- Requester and approver separation is enforced by database state and tested across independent sessions.
- A passing migration or concurrency test does not prove production identity integration, backup recovery, availability, or operational correctness.

## Decision

The Phase 1 Step 1 PostgreSQL governance foundation represented by candidate
commit `f76d5d0be180a2faa9bfc8a229ba63b686117b86` is accepted for continued non-production development.

The annotated acceptance tag shall point to the acceptance commit containing
this record.

## Exact Next Work

Begin Phase 1 Step 2 — Go PostgreSQL Runtime and Identity Context Boundary:

- Add a replaceable Go PostgreSQL adapter.
- Establish least-privileged connection-pool behavior.
- Bind authenticated runtime identity to transaction-scoped database context.
- Prove that identity context cannot leak between pooled connections.
- Add health and readiness behavior for database dependency state.
- Add disposable database integration, concurrency, rollback, and failure tests.
- Keep production authentication and credential delivery outside the Step 2 boundary unless separately accepted.
