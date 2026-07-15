# Phase 1 Step 2 Acceptance Record

## Record Status

- Decision: Accepted as a non-production Go PostgreSQL runtime, transaction-local identity-context, and portable-validation boundary
- Product: Iron Atlas
- Repository: `Iron-Signal-Systems/iron-atlas`
- Canonical repository: `https://github.com/Iron-Signal-Systems/iron-atlas.git`
- Branch: `dev`
- Acceptance date: `2026-07-15T17:05:07Z`
- Accepted tag: `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`

## Predecessor

- Accepted Phase 1 Step 1 tag: `phase-1-step-1-postgresql-governance-foundation-complete-v1`
- Accepted Phase 1 Step 1 commit: `f41d932beff01e3faf1aeb73d6386a22c95cdda8`
- Isolated predecessor revalidation: PASS

## Candidate

- Candidate implementation commit: `a56de3f9a9859f8daffea29b772091359dd3c3c9`
- Repository-complete evidence boundary commit: `8ea3e715aed989c742e6c9231614417122ead24d`
- Candidate Git archive SHA-256: `9985322850cd3ebd7be387783b91d1908f5140f787f367e6a13a3d522e25da83`
- Toolchain requirements SHA-256: `1a81bd0377a3e59a8002b4a6e0ae5b99e54280f9fcfd5c39f35d4761ad1149a1`
- Go version: `go version go1.26.5-X:nodwarf5 linux/amd64`
- pgx version: `v5.10.0`
- PostgreSQL server and client versions: `postgres (PostgreSQL) 18.4; psql (PostgreSQL) 18.4`
- Non-secret host-class fingerprint: `Arch Linux; Linux 6.6.87.2-microsoft-standard-WSL2 x86_64 GNU/Linux; CPUs=28; MemTotal=16223360 KiB; git version 2.55.0; Python 3.14.6`

## Accepted Scope

- Replaceable Go PostgreSQL adapter behind the domain change-service interface.
- Bounded least-privileged `pgxpool` connection pool.
- Transaction-local actor context with no session-scoped acting identity.
- Pooled-connection identity isolation across commit, rollback, failure, sequential reuse, and concurrency.
- Governed change creation and approval through accepted PostgreSQL service functions.
- Requester self-approval denial and database-authoritative durable authorization.
- PostgreSQL-aware readiness while liveness remains process-oriented.
- Disposable PostgreSQL integration tests using repository-provided environments.
- Portable validation, declared external toolchain requirements, evidence redaction, integrity records, and canonical clean-clone verification.
- Failure-preserving isolated historical-gate revalidation.

## Validation Evidence

- Local implementation gate: PASS — `validation/evidence/phase-1-step-2/local-phase-gate/20260715T154526Z-a56de3f9a985/`
- Canonical clean-clone validation: PASS — `validation/evidence/phase-1-step-2/canonical-clean-clone/20260715T154723Z-777019f9a55b/`
- Canonical clone commit: `777019f9a55b5cb2988a13496005a40c9789a47a`
- Applicable validator: `tools/validation/phase-gates/validate_phase1_step2.sh`
- Committed evidence paths: `validation/evidence/phase-1-step-2/local-phase-gate/20260715T154526Z-a56de3f9a985/` and `validation/evidence/phase-1-step-2/canonical-clean-clone/20260715T154723Z-777019f9a55b/`
- Evidence SHA-256 records: `validation/evidence/phase-1-step-2/local-phase-gate/20260715T154526Z-a56de3f9a985/sha256sums.txt` and `validation/evidence/phase-1-step-2/canonical-clean-clone/20260715T154723Z-777019f9a55b/sha256sums.txt`
- Historical predecessor validator failure propagation: PASS
- Repository validation: 46 PASS, 0 FAIL
- Test framework: 14 PASS, 0 FAIL
- Phase gate: 9 PASS, 0 FAIL
- Sequential identity-isolation iterations: 500
- Concurrent identity-isolation operations: 600
- Correctness result: PASS
- Resource observation: RECORDED_BY_DATABASE_TEST
- Performance thresholds: NOT_EVALUATED

## Reproducibility Statement

No implementation step may be accepted unless a clean clone from the canonical GitHub repository can execute its applicable validation using only version-controlled project artifacts, declared and verifiable external toolchain requirements, disposable test environments, and explicitly supplied non-repository secrets.

The exact pushed repository boundary `777019f9a55b5cb2988a13496005a40c9789a47a` was validated through a fresh clone of the canonical `dev` branch. Its sanitized transcript, environment fingerprint, metadata, summary, and checksums are committed in the repository-complete evidence boundary `8ea3e715aed989c742e6c9231614417122ead24d`.

The acceptance decision becomes effective only after the exact acceptance commit containing this record passes `validate_phase1_step2_acceptance.sh` through `verify_canonical_clone.sh`, that sanitized transcript is committed under `validation/evidence/phase-1-step-2/acceptance-canonical-clean-clone/`, and the annotated acceptance tag is published.

## Review and Approval

- Requester and implementer: John Wood (kb2vhn@gmail.com)
- Independent reviewer: Not assigned for this non-production single-maintainer development boundary
- Conflicts checked: Temporary single-maintainer development exception recorded below
- Operational two-person approval: Not exercised by this development acceptance

### Temporary Single-Maintainer Development Exception

The repository owner accepts this non-production development boundary under a temporary single-maintainer exception because it does not connect to production authentication, production credentials, live collectors, protected infrastructure evidence, or operational infrastructure.

This exception:

- Does not authorize production deployment or production use.
- Does not authorize production credentials, certificates, authentication, or live collection.
- Does not permit a requester to approve their own operational change.
- Does not weaken the documented two-person operational change-control contract.
- Expires before any production-boundary acceptance or operational use.

## Explicit Exclusions

- Production authentication and identity-provider integration
- Production credential delivery and rotation
- Production TLS and certificate provisioning
- Backup and restoration validation
- High availability
- Production connection and resource budgets
- Live infrastructure collection
- Protected evidence intake and storage
- Production monitoring and alert delivery
- Production performance budgets
- Production readiness

## Decision

The Phase 1 Step 2 boundary represented by implementation commit `a56de3f9a9859f8daffea29b772091359dd3c3c9` and repository-complete evidence boundary `8ea3e715aed989c742e6c9231614417122ead24d` is accepted for continued non-production development, subject to the canonical acceptance-commit verification and tag-publication conditions stated above.

The annotated acceptance tag shall point to the acceptance commit containing this record, not to either evidence-preparation commit.

## Exact Next Work

Begin Phase 1 Step 3 — Trusted Authentication and Governed Actor Resolution Boundary:

- Replace development identity headers in PostgreSQL mode with a pluggable trusted authentication adapter.
- Resolve authenticated external identities through governed external-identity and actor records.
- Fail closed for missing, unmapped, inactive, duplicated, or ambiguous identities.
- Bind the resolved actor to immutable server-side request context; request bodies, forms, query parameters, and ordinary headers must not select the actor.
- Define session, cookie, CSRF, replay, logout, expiry, and trusted-proxy controls appropriate to the selected authentication adapter.
- Test spoofing, confused-deputy, actor-status change, identity-remapping, session-expiry, and concurrent-request behavior.
- Keep production credential delivery, database TLS deployment, backup recovery, and high availability outside Step 3 unless separately implemented and accepted.
