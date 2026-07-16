# Iron Atlas

> An Iron Signal Systems project
>
> Built on purpose. Backed by discipline. Engineered to endure.
>
> Development status: Phase 1 Step 2 Go PostgreSQL runtime, identity-context, and portable-validation boundary accepted for non-production development; Phase 1 Step 3 trusted-authentication and governed-actor-resolution contract is the active candidate; no executable production authentication is accepted; not ready for production use

Iron Atlas is an authoritative, version-controlled system for infrastructure documentation, diagrams, inventory, automated discovery, project tracking, change management, validation, preventive health analysis, and formal acceptance.

The project is designed for network technicians, network administrators, network-security staff, reviewers, auditors, and infrastructure teams. It combines a lightweight Go service and HTML5 interface with modular collectors, vendor parsers, normalized infrastructure records, governed changes, independent approvals, topology generation, and replaceable external-system adapters.

## Initial Scope

- Go-first service architecture with an embedded HTML5 interface.
- Role-aware workspaces for network operations and security teams.
- Independent two-person control for governed changes.
- Firewall configuration ingestion and semantic analysis for FortiGate, OPNsense, and pfSense.
- Cisco IOS and IOS XE evidence collection for Catalyst 2960 through 2960-X, 9200, 9300, and 9500.
- Catalyst 9800 wireless-controller collection.
- Thirty-day full technical-support evidence collection plus lighter recurring health collection.
- Port, trunk, CDP/LLDP, MAC, spanning-tree, ACL, VLAN-pruning, and port-channel analysis.
- Draw.io source governance and generated/curated diagram separation.
- Canonical telemetry with replaceable Zabbix, webhook, syslog, and OpenMetrics adapters.
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

## Quick Start

Memory-backed development mode remains the default:

```bash
go test ./...
go run ./cmd/atlasd
```

PostgreSQL development mode requires the accepted migrations, runtime grants, and a protected runtime connection string:

```bash
export IRON_ATLAS_CHANGE_STORE=postgresql
export IRON_ATLAS_DEV_IDENTITY=true # controlled local testing only
export IRON_ATLAS_DATABASE_URL='postgres://atlas_application:REDACTED@localhost/iron_atlas?sslmode=verify-full'
go run ./cmd/atlasd
```

Do not commit a real database URL or credentials. Open `http://127.0.0.1:8080` after startup.

Memory mode defaults to development identity headers for the Phase 0 demonstration. PostgreSQL mode defaults them off and requires an explicit opt-in for controlled local testing. Development identity headers are never an acceptable production authentication boundary.

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
./tools/validation/phase-gates/validate_phase1_step2.sh
```

Formal acceptance additionally requires the exact pushed commit to pass `tools/validation/verify_canonical_clone.sh` from a clean clone of the canonical GitHub repository. Retained validation evidence is sanitized and committed under `validation/evidence/`.

## Documentation

Start with:

- [Documentation index](docs/README.md)
- [Target architecture](docs/architecture/TARGET-ARCHITECTURE.md)
- [Change management and two-person control](docs/architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md)
- [PostgreSQL migration and ownership model](docs/architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md)
- [PostgreSQL database security boundary](docs/architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md)
- [Go PostgreSQL runtime and identity context](docs/architecture/GO-POSTGRESQL-RUNTIME-AND-IDENTITY-CONTEXT.md)
- [Trusted authentication and governed actor resolution](docs/architecture/TRUSTED-AUTHENTICATION-AND-GOVERNED-ACTOR-RESOLUTION.md)
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
