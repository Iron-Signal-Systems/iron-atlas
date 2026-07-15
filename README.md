# Iron Atlas

> An Iron Signal Systems project
>
> Built on purpose. Backed by discipline. Engineered to endure.
>
> Development status: Phase 0 accepted as a non-production development baseline; Phase 1 Step 1 PostgreSQL foundation implementation candidate; not ready for production use

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

The accepted Phase 0 baseline contains a compilable standard-library-only Go implementation that demonstrates:

- A responsive HTML5 dashboard and role-oriented navigation.
- Health and readiness endpoints.
- A module registry.
- A memory-backed change workflow with requester/approver independence.
- A native Go Zabbix sender-protocol adapter.
- A FortiGate hierarchical configuration parser foundation.
- OPNsense and pfSense XML import probes.
- Cisco command-bundle parsing and trunk endpoint-attribution rules.
- Repository validation and unit tests.

Phase 1 Step 1 adds the manifest-driven PostgreSQL migration framework, ownership and role contracts, governed identity and authority persistence, independent approval enforcement, append-only history, and disposable database tests. It does **not** connect the HTML5 service to PostgreSQL or claim production authentication, live collection, complete vendor semantic analysis, or production readiness.

## Quick Start

```bash
go test ./...
go run ./cmd/atlasd
```

Open `http://127.0.0.1:8080`.

The initial server runs in development identity mode unless configured otherwise. Development identity headers are never an acceptable production authentication boundary.

## Repository Layout

```text
.
├── cmd/                     Go executables
├── configs/                 Configuration examples
├── deployment/              Arch Linux and systemd material
├── diagrams/                Draw.io sources and publication boundary
├── docs/                    Normative architecture and project documentation
├── integrations/            Replaceable external-system adapters
├── internal/                Shared application implementation
├── modules/                 Vendor and capability modules
├── projects/                Project-governance templates and records
├── changes/                 Change-management templates and records
├── sql/                     Planned PostgreSQL schema and migration manifest
├── test-framework/          Test orchestration and retained results boundary
└── tools/validation/        Repository checks and phase gates
```

## Validation

```bash
./tools/validation/validate_repository.sh
./tools/validation/phase-gates/validate_phase0_step1.sh
```

## Documentation

Start with:

- [Documentation index](docs/README.md)
- [Target architecture](docs/architecture/TARGET-ARCHITECTURE.md)
- [Change management and two-person control](docs/architecture/CHANGE-MANAGEMENT-AND-TWO-PERSON-CONTROL.md)
- [Firewall semantic analysis](docs/architecture/FIREWALL-CONFIGURATION-SEMANTIC-ANALYSIS.md)
- [Cisco evidence collection](docs/architecture/CISCO-EVIDENCE-COLLECTION-AND-PREVENTIVE-HEALTH.md)
- [Cisco trunk analysis](docs/architecture/CISCO-TRUNK-AND-ENDPOINT-ATTRIBUTION.md)
- [External-system-independent telemetry](docs/architecture/EXTERNAL-SYSTEM-INDEPENDENT-TELEMETRY.md)
- [Phased implementation roadmap](docs/roadmap/IMPLEMENTATION-ROADMAP.md)
- [Phase 0 accepted baseline](docs/acceptance/PHASE-0-STEP-1-ACCEPTANCE-RECORD.md)
- [Phase 0 acceptance errata](docs/acceptance/PHASE-0-STEP-1-ACCEPTANCE-ERRATA.md)
- [PostgreSQL migration and ownership model](docs/architecture/POSTGRESQL-MIGRATION-AND-OWNERSHIP-MODEL.md)
- [PostgreSQL database security boundary](docs/architecture/POSTGRESQL-DATABASE-SECURITY-BOUNDARY.md)

## Security Boundary

Raw firewall backups, technical-support output, credentials, SSH private keys, shared secrets, certificates, and unredacted evidence are prohibited from Git. The repository contains contracts, code, schemas, fixtures, redacted examples, and accepted documentation only.

## License

BSD 3-Clause. See [LICENSE](LICENSE).
