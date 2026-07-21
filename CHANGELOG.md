# Changelog

## Unreleased

### Added

- Added the provider-neutral OIDC assurance-evidence candidate using Atlas-controlled synthetic claims, explicit `auth_time` correlation for `acr` and `amr`, exact governed method sets, hostile tests, synchronized documentation, and a bounded implementation gate.
- Added the architecture and roadmap alignment candidate defining module runtime failure containment, scheduled evidence ingestion, monitoring and freshness, immutable candidates and atomic acceptance, signed purpose-limited Atlas–IFI snapshots, fail-closed adversarial invariants, external-provider MFA requirements, and signed candidate/post-merge repository trust.
- Added an alignment record preserving historical Phase 0, Phase 1 Steps 1 and 2, and existing Step 3 checkpoints without relabeling them as formal acceptance.

- Added the Phase 1 Step 3 authentication-assurance candidate with bounded OIDC `acr`, `amr`, and `auth_time` normalization, explicit provider-MFA and phishing-resistant policy, stale-authentication step-up, callback hardening, mandatory MFA and policy-version binding before session persistence, synchronized documentation, hostile tests, and an implementation phase gate.
- Made the actionable `FINAL RESULT` the terminal output line for migrated validation commands, including the primary failing check and cause on failure.
- Added meaningful validation failure reporting with primary-cause extraction, per-check logs, deduplicated cascades, skipped-dependent checks, and final report paths.
- Added repository-managed, version-pinned `govulncheck` invoked through `go tool`, removing workstation PATH dependence.
- Revalidated the frozen HTTP login checkpoint from its exact signed predecessor commit instead of running historical assertions against successor documentation.
- Added the Phase 1 Step 3 authenticated server-side session candidate with digest-only PostgreSQL persistence, secure opaque cookies, current governed actor re-resolution, assurance metadata for future MFA, hostile and race tests, synchronized documentation, and an implementation phase gate.

- Added the Phase 1 Step 3 HTTP login and callback candidate with secure host-only browser state binding, exact callback parsing, issuer checks, provider-error cancellation, verified-principal handoff, hostile and race tests, synchronized documentation, and an implementation phase gate.
- Replaced the mandatory Iron Atlas roadmap with a preserved-history Phase 0–12 core sequence covering evidence protection, operating-system-based Cisco normalization, FortiGate, BloodHound identity context, cross-source intelligence, change decisions, interface completion, read-only collection, production recovery, and controlled-pilot acceptance.
- Added a planned phase-and-gate map that preserves existing validators, defines contract, implementation, integration, and acceptance gates, and keeps optional integration modules outside core operational acceptance.
- Added the canonical 1254 × 1254 Iron Atlas crest with governed digest, accessibility text, README rendering, and brand-asset boundaries.
- Added an exact published ISRAS 0.1.4 adoption candidate using the official release initializer, pinned reusable workflow, project pin, release-artifact verification, and preserved Atlas-specific validation.

- A pinned maintained YAML v4 node decoder, Atlas node-conversion boundary, multiline quoted-scalar fixture, hostile admission tests, explicit resource limits, and cancellation checks for the experimental FortiGate YAML adapter.
- Upload-safe `fortigate-inspect -redact` summaries that omit local paths, device identity, FortiOS version, and finding details while retaining counts and bounded decoder positions.
- Upload-safe structure diagnostics that report only public FortiOS section labels and aggregate layout counts while omitting private scalar values, unknown keys, and VDOM names.
- Upload-safe semantic-quality aggregates for normalized record kinds, resolved and unresolved reference classes, built-ins, findings, and coverage warnings.
- Stable vendor-independent normalized-record metrics and explicit built-in-reference records.
- A bounded compatibility repair for confirmed Fortinet-generated invalid YAML forms, with semantic, false-positive, and resource-limit tests.
- ADR-0007 and synchronized FortiGate YAML snapshot architecture documentation.

- Phase 1 Step 3 OIDC authorization-code exchange with PKCE S256, 256-bit state and nonce, SHA-256 state-digest storage, bounded one-time in-memory preauthentication transactions, exact redirect binding, token-response bounds, replay resistance, concurrency proofs, and secret-redaction tests.
- Phase 1 Step 3 authorization-code and PKCE architecture, testing, traceability, static validation, regression, and phase-gate synchronization.

- Phase 1 Step 3 typed authentication-mode middleware, immutable request identity, future authenticator and actor-resolver interfaces, and targeted hostile tests.
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

- Accepted the architecture and roadmap alignment as a documentation and governance boundary at signed commit `12569192da89a1a34f4ebfe107c4d02c60cbdb09`, with PR #17, merge commit `5de9e1f5f9770f12b56a046dc735b769cc842a02`, and successful validate, Portable validation, and ISRAS hosted runs recorded in the alignment evidence.

- Removed Atlas-local password, TOTP-secret, QR-enrollment, and ordinary recovery-code ownership from the required authentication roadmap; successor work now begins with provider-neutral assurance evidence, followed by evidence-backed representative-provider compatibility, session lifecycle, CSRF, trusted proxy, production wiring, governed emergency access, integration, and formal Step 3 acceptance.
- Synchronized README, architecture, requirements, testing, roadmap, gates, acceptance, governance, and validation around the signed BUSL boundary.

- Prospectively transitioned Iron Atlas from BSD 3-Clause to Business Source License 1.1 (`BUSL-1.1`) from the signed `cc93fdd` predecessor, with no Additional Use Grant, a 2030-07-18 Change Date, AGPLv3-only Change License, preserved historical BSD text, explicit trademark separation, machine-readable validation, and a governed post-licensing alignment backlog.
- Corrected future authentication direction so Atlas consumes and enforces approved external-provider MFA assurance rather than implementing local password or TOTP credential ownership.

- Replaced the FortiGate adapter's handwritten physical-line YAML grammar with maintained node decoding while preserving Atlas normalization, source-location, evidence, and snapshot contracts.
- Added bounded layout detection for canonical fixtures, native direct CMDB sections, and native VDOM containers while retaining fail-closed rejection of unrelated YAML.
- Extended the firewall reference graph with object-kind resolution, explicit built-ins, unresolved-versus-ambiguous classification, and normalized-reference limits.

- Clarified Atlas as an active complement to Zabbix, Graylog, Security Onion, Draw.io, and vendor platforms rather than a replacement for their mature operational functions.
- Reordered the unaccepted first-product sequence to prioritize an offline Catalyst 9300L/9300, 9500, and 9800 infrastructure-value slice before Fortinet policy analysis.
- Added governed direction for Zabbix maps, dashboards, templates, discovery, and reports; Graylog lookup, enrichment, queries, pipelines, streams, dashboards, and reports; and future separately accepted external-system provisioning.
- Synchronized current documentation status with the merged Phase 1 Step 3 OIDC discovery, JWKS, and ID-token verification checkpoint.

- Replaced the legacy development-identity boolean with explicit `development` and `production` authentication modes.
- Moved HTTP actor construction out of handlers and into the authentication boundary.
- Raised the Go module baseline to Go 1.25 for the accepted `pgx` v5 runtime dependency.
- Made the change-service interface context-aware and persistence-neutral.
- Kept memory mode as the default development store while adding explicit PostgreSQL mode.
- Updated the HTML5 and API handlers to fail closed on persistence errors.
- Updated CI and the disposable PostgreSQL runner to execute Go database integration tests.

### Security

- Upgraded the indirect `golang.org/x/text` dependency to `v0.39.0` to remediate `GO-2026-5970` reported by hosted vulnerability validation.
- FortiGate YAML fails closed on oversized or excessive input, aliases, anchors, custom tags, duplicate keys, multiple documents, unsupported scalar forms, and normalized-record or finding limit violations.
- Upload-safe semantic output uses fixed allowlists and fallback classifications so source-derived labels, names, values, paths, and finding details cannot be reflected into retained logs.

- Production mode now rejects development identity headers and protected requests fail closed when no trusted adapter is configured.
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

- Corrected nested validation cause extraction so a subordinate runner's terminal `FINAL RESULT: FAIL` takes precedence over earlier intentional error fixtures, and synchronized reporting validation with exact authentication-assurance checkpoint revalidation.

- Repaired confirmed Fortinet-generated invalid YAML adjacent multi-value fragments and restricted literal object-name keys while leaving ordinary valid YAML unchanged and broader malformed YAML fail-closed.
- Repaired admission of native FortiGate YAML backups whose root does not use the prototype's synthetic `global` or `vdom` wrapper.

- Synchronized the Go PostgreSQL runtime validator with the typed authentication-mode boundary while retaining explicit PostgreSQL production-default and controlled-development assertions.
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

### Phase 1 Step 3 PostgreSQL governed actor-resolution candidate

- Added migration 007 with a fixed-search-path, least-privileged
  `atlas.resolve_governed_actor(text, text)` function.
- Added a PostgreSQL `authentication.ActorResolver` with explicit Atlas role
  mapping and fail-closed unknown-state handling.
- Added disposable-database, integration, hostile-input, concurrency, race,
  static-validation, and phase-gate coverage.
- Preserved the prohibition on broad application-role reads of governed
  identity and role tables.

### Phase 1 Step 3 OIDC ID-token verification candidate

- Added pinned `coreos/go-oidc` verification behind a stricter Atlas boundary.
- Added exact HTTPS discovery and endpoint checks, asymmetric algorithm
  allowlisting, bounded input and time policy, nonce and authorized-party
  validation, duplicate sensitive-claim rejection, and access-token-hash
  verification.
- Added a disposable TLS provider emulator covering signature, issuer,
  audience, nonce, time, malformed input, key rotation, outage, race, and
  concurrency behavior.
- Preserved explicit nonclaims for authorization-code exchange, PKCE state,
  browser sessions, CSRF, logout, trusted proxies, and production readiness.
- Synchronized the standard hosted `validate` workflow with the
  repository-owned pinned tool bootstrap so `govulncheck` and future declared
  Go validation tools are available before the complete test framework runs.
