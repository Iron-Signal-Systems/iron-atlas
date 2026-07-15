# Changelog

## Unreleased

### Accepted

- Accepted the Phase 1 Step 1 PostgreSQL governance foundation as a non-production development boundary under annotated tag `phase-1-step-1-postgresql-governance-foundation-complete-v1`.
- Recorded the exact candidate commit, deterministic Git archive hash, validation evidence, limitations, security assumptions, temporary single-maintainer development exception, and exact Phase 1 Step 2 work.

### Fixed

- Removed ambiguous PL/pgSQL actor variable and column resolution from the governed change and approval functions.
- Added disposable-database regression assertions proving the internal actor helper remains outside the application API while the governed change API resolves the intended active actor.

### Added

- Phase 1 Step 1 manifest-driven PostgreSQL migration framework.
- Database owner, schema owner, migrator, application, read-only, auditor, and test-runner role contracts.
- Governed actor, external identity, role, authority, change, approval, decision, and audit persistence.
- Database-enforced requester and approver independence.
- Append-only approval, status-history, decision, audit, and migration records.
- Disposable PostgreSQL correctness, security, idempotency, and concurrency tests.
- Phase 1 Step 1 validation gate with isolated revalidation of the accepted Phase 0 predecessor.
- Phase 0 acceptance errata preserving the historical tag while correcting record-generation defects.

### Changed

- Archived the Phase 0 monolithic SQL design candidate outside the executable migration manifest.
- Replaced the initial migration manifest with ordered Phase 1 migrations.
- Updated CI, repository validation, documentation indexes, testing guidance, and the roadmap for Phase 1 Step 1.
- Kept generated test results outside Git; formal acceptance evidence remains in immutable acceptance records.

## Phase 0

- Initial Iron Atlas repository architecture.
- Go HTML5 service candidate.
- Independent change-approval implementation and tests.
- Initial firewall and Cisco parser boundaries.
- Native Go Zabbix sender adapter.
- Arch Linux deployment baseline.
- Documentation, testing, validation, and phase-gate structure.
- Phase 0 accepted as a non-production development baseline under annotated tag
  `phase-0-repository-and-executable-baseline-complete-v1`.
