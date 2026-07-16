# Changelog

## Unreleased

### Added

- Phase 1 Step 3 trusted-authentication and governed-actor-resolution architecture contract.
- Phase 1 Step 3 requirements traceability, adversarial testing model, acceptance template, static validator, phase-entry gate, and regression test.
- Canonical GitHub clean-clone validation as a mandatory acceptance invariant.
- Machine-readable external toolchain requirements and verification.
- Sanitized, checksummed, committed validation-evidence recording and validation.
- Repository-provided canonical clone verifier and portability regression tests.

- Phase 1 Step 2 replaceable Go PostgreSQL change-service adapter.
- Bounded least-privileged `pgxpool` runtime configuration.
- Transaction-local authenticated actor context for governed mutations.
- Persistent change creation and approval through accepted PostgreSQL functions.
- Database-aware health and readiness dependency behavior.
- Sequential, concurrent, commit, rollback, and failed-transaction actor-isolation tests.
- Go PostgreSQL runtime architecture, testing, acceptance-template, and phase-gate documentation.

### Changed

- Raised the Go module baseline to Go 1.25 for the accepted `pgx` v5 runtime dependency.
- Made the change-service interface context-aware and persistence-neutral.
- Kept memory mode as the default development store while adding explicit PostgreSQL mode.
- Updated the HTML5 and API handlers to fail closed on persistence errors.
- Updated CI and the disposable PostgreSQL runner to execute Go database integration tests.

### Security

- Defined fail-closed external-identity resolution, Atlas-owned role authority, immutable request identity, bounded server-side sessions, CSRF, replay, trusted-proxy, and authentication-secret redaction requirements.
- Corrected the Step 2 isolated-predecessor gate so cleanup cannot mask a failed or missing historical validator.
- Database connection strings remain runtime-only secrets and are prohibited from committed configuration.
- Acting identity is set only with transaction-local `set_config(..., true)` and never at pooled-session scope.
- Governed PostgreSQL writes continue to use only accepted security-definer service functions.

### Accepted

- Accepted the Phase 1 Step 2 Go PostgreSQL runtime, transaction-local identity-context, and portable-validation boundary under annotated tag `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`.
- Recorded the exact implementation and evidence commit chain, deterministic archive and toolchain hashes, committed local and canonical clean-clone evidence, limitations, temporary development exception, and exact Step 3 work.
- Accepted the Phase 1 Step 1 PostgreSQL governance foundation as a non-production development boundary under annotated tag `phase-1-step-1-postgresql-governance-foundation-complete-v1`.
- Recorded the exact candidate commit, deterministic Git archive hash, validation evidence, limitations, security assumptions, temporary single-maintainer development exception, and exact Phase 1 Step 2 work.

### Fixed

- Removed ambiguous PL/pgSQL actor variable and column resolution from the governed change and approval functions.
- Added a disposable-database regression assertion proving actor context resolves to the intended active actor.

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
