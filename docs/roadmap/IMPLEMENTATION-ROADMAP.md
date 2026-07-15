# Implementation Roadmap

## Phase 0 — Repository and Executable Baseline

**Status:** Accepted as a non-production development baseline under tag `phase-0-repository-and-executable-baseline-complete-v1`.

- Mission, architecture, requirements, and boundaries
- Go module and embedded HTML5 interface
- RBAC and requester-approval independence candidate
- Module registry and parser contracts
- Initial FortiGate, OPNsense, pfSense, Cisco, and Zabbix code candidates
- Unit tests and repository phase gate
- Minimal Arch deployment material

Acceptance does not claim production authentication, persistence, collection, or semantic completeness.

## Phase 1 — PostgreSQL Foundation and Governed Identity

**Status:** Step 1 accepted under tag `phase-1-step-1-postgresql-governance-foundation-complete-v1`; later Phase 1 work is not accepted.

### Step 1 — Migration and Database Governance Foundation

- Manifest-driven migrations
- Production role and ownership topology
- Governed actors and external identities
- Role binding and authority records
- Change, approval, decision, and audit records
- Database-enforced requester and approver independence
- Append-only history
- Disposable database tests and concurrency proofs

Step 1 does not connect the Go service to PostgreSQL or establish production authentication.

### Later Phase 1 Work

- Go PostgreSQL runtime adapter
- Transaction and connection-pool identity-context handling
- Identity-provider integration boundary
- Production credential delivery and rotation
- Database backup and restoration test boundary

## Phase 2 — Evidence Intake and Storage

- Mutually authenticated ingestion API
- Durable staging queue
- Signed evidence bundles
- Encrypted content-addressed storage
- Redaction and classification
- Parser isolation and resource governance

## Phase 3 — Firewall Ingestion

- FortiGate native parser and semantic normalization
- FortiGate YAML adapter
- OPNsense and pfSense XML normalization
- Interface, route, gateway, SD-WAN, policy, object, NAT, and VPN graph
- Traffic-path explanation
- Golden fixtures and adversarial parsing

## Phase 4 — Cisco Collection Foundation

- Device enrollment and host-key pinning
- NPS/RADIUS service authentication
- Restricted command profiles
- 2960-family IOS collection
- 9200/9300/9500 IOS XE collection
- Command transcripts and protected evidence
- Daily/weekly health and 30-day comprehensive schedule

## Phase 5 — Cisco Semantic and Preventive Analysis

- Access-port and endpoint attribution
- Trunk, pruning, CDP/LLDP, STP, ACL, QoS, and port-channel analysis
- Resource, environment, stack, software, and counter trends
- Catalyst 9800 controller, AP, profile, tag, and client analysis
- Finding correlation and duplicate suppression

## Phase 6 — Projects, Changes, and Acceptance

- Complete project portfolio
- Two-person and multi-approver policies
- Pre/post evidence comparison
- Expected and unexpected difference disposition
- Rollback and emergency change handling
- Formal acceptance and closure

## Phase 7 — Topology and Diagrams

- Normalized graph
- Generated switch-port, Layer 2, VLAN, route, SD-WAN, firewall, and wireless views
- Draw.io-compatible generated sources
- Curated diagram lifecycle
- SVG/PDF publication and drift checks

## Phase 8 — External Integrations

- Production Zabbix sender delivery
- Zabbix provisioning boundary where justified
- OpenMetrics, syslog, webhook, and SIEM adapters
- Delivery outbox, retry, backpressure, and dead-letter handling

## Phase 9 — Production Security and Recovery

- Signed builds and provenance
- SBOM and package integrity
- Off-host logging and integrity anchors
- Backup protection and restoration validation
- Break-glass
- Trusted rebuild and compromise recovery
- Operational acceptance on representative hardware
