# Iron Atlas

> An Iron Signal Systems project
>
> Built on purpose. Backed by discipline. Engineered to endure.
>
> Development status: Phase 1 Step 2 is accepted; the Phase 1 Step 3 authentication foundation, governed actor resolver, and bounded OIDC discovery/JWKS/ID-token verification checkpoints are merged; authorization-code exchange with PKCE S256 and bounded one-time in-memory preauthentication transactions is the active implementation candidate; no HTTP login/callback route, durable session, cookie, CSRF, logout, or trusted-proxy implementation is accepted; not ready for production use

Iron Atlas is an authoritative, version-controlled system for infrastructure documentation, diagrams, inventory, automated discovery, project tracking, change management, validation, preventive health analysis, and formal acceptance.

The project is designed for network technicians, network administrators, network-security staff, reviewers, auditors, and infrastructure teams. It combines a lightweight Go service and HTML5 interface with modular collectors, vendor parsers, normalized infrastructure records, governed changes, independent approvals, topology generation, and replaceable external-system adapters.

## Product Position

Iron Atlas is an infrastructure-intelligence, documentation, and integration platform. It complements established operational systems instead of recreating capabilities that those systems already perform well.

- Zabbix remains responsible for continuous monitoring, alerting, graphing, escalation, maintenance, and availability history. Atlas adds authoritative infrastructure identity, topology, evidence provenance, reconciliation, and governed assistance for maps, dashboards, templates, low-level discovery, and reports.
- Graylog remains responsible for centralized log and SNMP-trap collection, search, retention, and investigation. Atlas adds device, site, interface, VLAN, role, and topology context that can improve lookup tables, pipelines, streams, queries, dashboards, and reports.
- Security Onion remains responsible for network-security monitoring, packet analysis, detection, and investigation. Atlas may provide governed asset and topology context without becoming its detection engine.
- Cisco, Fortinet, and other infrastructure platforms remain responsible for operation and enforcement. Atlas observes, documents, compares, explains, and validates their state; it is not initially a controller or automated-remediation system.

Atlas may consume approved evidence from existing systems and may generate recommendations, exports, definitions, maps, dashboards, queries, reports, and other reviewable integration artifacts for them. Any future write or provisioning integration must be separately accepted, previewable, attributable, bounded, idempotent where practical, reversible where practical, and subject to the applicable approval policy.

## Initial Scope

- Go-first service architecture with an embedded HTML5 interface.
- Role-aware workspaces for network operations and security teams.
- Independent two-person control for governed changes.
- Firewall configuration ingestion and semantic analysis for FortiGate, OPNsense, and pfSense.
- Cisco IOS and IOS XE evidence collection, with the first-value slice centered on Catalyst 9300L/9300 access switching, Catalyst 9500 core and distribution switching, and Catalyst 9800 wireless controllers.
- Compatibility profiles for Catalyst 2960 through 2960-X and Catalyst 9200 after the first Cisco slice is established.
- Thirty-day full technical-support evidence collection plus lighter recurring health collection.
- Port, trunk, CDP/LLDP, MAC, spanning-tree, ACL, VLAN-pruning, and port-channel analysis.
- Draw.io source governance and generated/curated diagram separation.
- Canonical telemetry and infrastructure context with replaceable Zabbix, Graylog/syslog, webhook, SIEM, and OpenMetrics adapters.
- Reviewable generation and reconciliation support for external-system maps, dashboards, queries, lookup data, templates, reports, and documentation.
- Minimal Arch Linux deployment using systemd and only required runtime packages.

## Current Executable Boundary

The accepted Phase 0 baseline contains the embedded HTML5 interface, module registry, memory-backed change workflow, initial vendor parsers, native Go Zabbix sender adapter, and repository validation framework.

Accepted Phase 1 Step 1 adds manifest-driven PostgreSQL migrations, database ownership and role contracts, governed identity and authority persistence, database-enforced independent approval, append-only history, and disposable PostgreSQL tests.

Accepted Phase 1 Step 2 adds:

- A replaceable PostgreSQL implementation of the Go change-service interface.
- A bounded least-privileged `pgxpool` connection pool.
- Transaction-local authenticated actor context.
- Persistent change creation and approval through accepted PostgreSQL functions.
- Rollback and pooled-connection identity-isolation tests.
- PostgreSQL-aware readiness behavior.

Step 2 does **not** establish production authentication, credential delivery, TLS provisioning, backup recovery, high availability, live collection, or production readiness.


The merged Phase 1 Step 3 authentication-foundation checkpoint adds:

- Typed `development` and `production` authentication modes.
- A dedicated authentication middleware and private immutable request-context identity.
- Future production authenticator and governed actor-resolver interfaces.
- Explicit rejection of development identity headers in production mode.
- Fail-closed protected routes when no production adapter is configured.
- Public health, readiness, and static-asset routes that do not manufacture an actor.

This checkpoint does not by itself implement a production authentication adapter, sessions, CSRF protection, trusted-proxy enforcement, or production authentication.

## Quick Start

Memory-backed development mode remains the default:

```bash
go test ./...
go run ./cmd/atlasd
```

PostgreSQL development mode requires the accepted migrations, runtime grants, and a protected runtime connection string:

```bash
export IRON_ATLAS_CHANGE_STORE=postgresql
export IRON_ATLAS_AUTHENTICATION_MODE=development # controlled local testing only
export IRON_ATLAS_DATABASE_URL='postgres://atlas_application:REDACTED@localhost/iron_atlas?sslmode=verify-full'
go run ./cmd/atlasd
```

Do not commit a real database URL or credentials. Open `http://127.0.0.1:8080` after startup.

Memory mode defaults to `development` authentication for the Phase 0 demonstration. PostgreSQL mode defaults to `production`, where protected routes fail closed until a trusted production adapter and governed actor resolver are configured. Controlled PostgreSQL testing may explicitly set `IRON_ATLAS_AUTHENTICATION_MODE=development`. The legacy `IRON_ATLAS_DEV_IDENTITY` boolean is rejected, and development identity headers are never an acceptable production authentication boundary.

## Repository Layout

```text
.
├── cmd/                     Go executables
├── configs/                 Non-secret configuration examples
├── deployment/              Arch Linux and systemd material
├── diagrams/                Draw.io sources and publication boundary
├── docs/                    Normative architecture and project documentation
├── integrations/            Replaceable external-system adapters
├── internal/                Shared application and PostgreSQL runtime implementation
├── modules/                 Vendor and capability modules
├── projects/                Project-governance templates and records
├── changes/                 Change-management templates and records
├── sql/                     Governed PostgreSQL schema and migration manifest
├── test-framework/          Test orchestration and transient local results
├── validation/              Toolchain contract and committed sanitized evidence
└── tools/validation/        Repository checks, evidence tools, and phase gates
```

## Validation

```bash
python3 tools/validation/validate_toolchain.py
./test-framework/run_all.sh
./tools/validation/phase-gates/validate_phase1_step3_oidc_authorization_code_pkce.sh
```

Formal acceptance additionally requires the exact pushed commit to pass `tools/validation/verify_canonical_clone.sh` from a clean clone of the canonical GitHub repository. Retained validation evidence is sanitized and committed under `validation/evidence/`.

## Documentation

Start with:

- [Documentation index](docs/README.md)
- [Target architecture](docs/architecture/TARGET-ARCHITECTURE.md)
- [Operational-system complement and integration model](docs/architecture/OPERATIONAL-SYSTEM-COMPLEMENT-AND-INTEGRATION-MODEL.md)
- [Change management and two-person control](docs/architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md)
- [PostgreSQL migration and ownership model](docs/architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md)
- [PostgreSQL database security boundary](docs/architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md)
- [Go PostgreSQL runtime and identity context](docs/architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md)
- [Trusted authentication and governed actor resolution](docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md)
- [OIDC authorization-code and PKCE transaction implementation](docs/architecture/OIDC-AUTHORIZATION-CODE-AND-PKCE-TRANSACTION-IMPLEMENTATION.md)
- [ADR-0004 — pgx PostgreSQL runtime driver](docs/decisions/ADR-0004-PGX-POSTGRESQL-RUNTIME-DRIVER.md)
- [Go PostgreSQL runtime integration testing](docs/testing/GO-POSTGRESQL-RUNTIME-INTEGRATION-TESTING.md)
- [Portable validation and canonical repository acceptance](docs/architecture/PORTABLE-VALIDATION-AND-CANONICAL-REPOSITORY-ACCEPTANCE.md)
- [Canonical clean-clone validation workflow](docs/operations/CANONICAL-CLEAN-CLONE-VALIDATION.md)
- [Firewall semantic analysis](docs/architecture/FIREWALL-CONFIGURATION-SEMANTIC-ANALYSIS.md)
- [Cisco evidence collection](docs/architecture/CISCO-EVIDENCE-COLLECTION-AND-PREVENTIVE-HEALTH.md)
- [Cisco trunk analysis](docs/architecture/CISCO-TRUNK-AND-ENDPOINT-ATTRIBUTION.md)
- [External-system-independent telemetry](docs/architecture/EXTERNAL-SYSTEM-INDEPENDENT-TELEMETRY.md)
- [Atlas primary-focus execution plan](docs/roadmap/ATLAS-PRIMARY-FOCUS-EXECUTION-PLAN.md)
- [Phased implementation roadmap](docs/roadmap/IMPLEMENTATION-ROADMAP.md)
- [Phase 0 accepted baseline](docs/acceptance/PHASE-0-STEP-1-ACCEPTANCE-RECORD.md)
- [Phase 0 acceptance errata](docs/acceptance/PHASE-0-STEP-1-ACCEPTANCE-ERRATA.md)
- [Phase 1 Step 1 accepted PostgreSQL governance foundation](docs/acceptance/PHASE-1-STEP-1-ACCEPTANCE-RECORD.md)
- [Phase 1 Step 2 accepted Go PostgreSQL runtime, identity-context, and portable-validation boundary](docs/acceptance/PHASE-1-STEP-2-ACCEPTANCE-RECORD.md)

## Security Boundary

Raw firewall backups, technical-support output, credentials, database URLs, SSH private keys, shared secrets, certificates, and unredacted evidence are prohibited from Git. The repository contains contracts, code, schemas, fixtures, redacted examples, and accepted documentation only.

## License

BSD 3-Clause. See [LICENSE](LICENSE).

## Governed Actor Resolution Candidate

The merged bounded Phase 1 Step 3 checkpoint adds a least-privileged PostgreSQL
`ActorResolver` behind `atlas.resolve_governed_actor(text, text)`. The
application role receives function execution only, not broad access to governed
identity or role tables. Missing, inactive, disabled, retired, unmapped, or
unsupported governed state fails closed.

This checkpoint still does not implement an external identity-provider adapter,
sessions, CSRF protection, trusted-proxy enforcement, or production
authentication.

## OIDC ID-Token Verification Candidate

The merged bounded Phase 1 Step 3 predecessor verifies exact HTTPS provider
discovery, remote JWKS signatures, issuer, audience, authorized party, permitted
asymmetric algorithms, expiry, issued-at, not-before, nonce, stable subject,
access-token hash when present, duplicate sensitive claims, key rotation, and
provider outage behavior.

That verifier checkpoint does not itself provide browser login routes,
authorization-code exchange, PKCE transaction handling, cookies, sessions,
CSRF, logout, trusted-proxy enforcement, or production authentication.

## OIDC Authorization-Code and PKCE Candidate

The active bounded Phase 1 Step 3 candidate creates 256-bit state, nonce, and
PKCE verifier values, stores only a SHA-256 state digest, requires discovered
PKCE S256 support, atomically consumes each transaction once, exchanges the code
through the exact discovered HTTPS token endpoint and redirect URI, bounds the
provider response, and returns only a verified provider-neutral principal.

The candidate uses a bounded in-memory transaction store. Process restart
invalidates outstanding login attempts. It does not add HTTP login or callback
routes, browser cookies, durable sessions, governed actor wiring, CSRF, logout,
trusted-proxy enforcement, production credential delivery, formal Step 3
acceptance, or production readiness.
