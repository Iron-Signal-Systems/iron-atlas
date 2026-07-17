# Implementation Roadmap

## Current Primary Focus

Iron Atlas is the primary active product-development effort.

The implementation roadmap preserves accepted foundation boundaries while delivering vertical product slices that answer real network and security questions as early as safely possible.

> **Cisco and FortiGate are the first two major evidence sources. The product is the correlated answer, risk analysis, and decision evidence produced from them.**

## Vertical Product Slices Across Phases

The phase sequence remains authoritative, but implementation within and across compatible phases should produce bounded vertical slices.

### Slice A — IP, CIDR, Subnet, and VLAN Intelligence

- evidence ingestion;
- canonical identity;
- prefix containment and longest-prefix match;
- VLAN and interface placement;
- routes and gateways;
- policy and ACL references;
- NAT, VIP, VPN, and SD-WAN relationships;
- evidence and unknowns.

### Slice B — Reachability and Attack-Path Explanation

- source and destination context;
- Layer-2 and Layer-3 path;
- route selection;
- firewall policy and ACL evaluation;
- NAT;
- VPN and SD-WAN;
- return path;
- trust-boundary crossings;
- optional BloodHound identity privilege and critical-asset context;
- explicit separation of identity capability, network capability, and combined path;
- evidence-backed result.

### Slice C — Change Impact and Decision Support

- prior, current, proposed, and post-change state;
- paths and dependencies changed;
- blast radius;
- risk of approval;
- risk of denial or delay;
- engineering plan;
- leadership decision summary;
- validation, rollback, and acceptance.

A slice may begin with partial evidence. Unsupported and unknown behavior must remain explicit.

## Phase 0 — Repository and Executable Baseline

**Status:** Accepted non-production baseline.

- Mission, architecture, requirements, and boundaries
- Go module and embedded HTML5 interface
- RBAC and requester-approval independence candidate
- Module registry and parser contracts
- Initial FortiGate, OPNsense, pfSense, Cisco, and Zabbix candidates
- Unit tests and repository validation
- Minimal Arch deployment material

Acceptance does not claim production authentication, persistence, collection, semantic completeness, or production readiness.

## Phase 1 — PostgreSQL Foundation and Governed Identity

**Status:** Steps 1 and 2 accepted; Step 3 and remaining production boundaries incomplete.

### Step 1

- Manifest-driven migrations
- Production role and ownership topology
- Governed actors and external identities
- Role binding, authority, approval, decision, and audit records
- Database-enforced requester and approver independence
- Append-only history
- Disposable database and concurrency tests

### Step 2

- Replaceable PostgreSQL change-service implementation
- Least-privileged bounded connection pool
- Transaction-local actor context
- Persistent change creation and approval
- Rollback and pooled-connection isolation
- PostgreSQL-aware readiness

### Step 3 and Remaining Work

- Trusted production authentication
- Governed actor resolution
- Bounded sessions, cookies, CSRF, logout, and trusted proxies
- Credential delivery and rotation
- PostgreSQL TLS
- Backup and recovery
- Production connection and resource budgets

Offline evidence, parser, canonical-model, and documentation work may continue in isolated workstreams without representing live collection or production authority as accepted.

## Phase 2 — Evidence Intake and Storage

- Versioned evidence-bundle contracts
- Manual and authenticated intake
- Durable staging
- Signed evidence bundles
- Encrypted content-addressed storage
- Duplicate detection and quarantine
- Redaction and classification
- Parser isolation
- Cancellation, timeout, and resource governance
- Provenance from receipt through normalized output

## Phase 3 — Cisco Offline Evidence and Normalization Foundation

- Sanitized versioned command bundles
- Catalyst 9300L/9300 access-switch profiles
- Catalyst 9500 core and distribution profiles
- Catalyst 9800 controller profiles
- Compatibility profiles for Catalyst 9200 and 2960 families after first-value support
- Device, stack, interface, VLAN, trunk, neighbor, port-channel, spanning-tree, route, ACL, and wireless normalization
- Running configuration and diagnostic evidence
- Unsupported, malformed, truncated, conflicting, and partial-state handling
- Golden fixtures, adversarial parsing, cancellation, and resource governance

## Phase 4 — Cisco Semantic Analysis, Topology, and Restricted Collection

- Layer-2 and Layer-3 topology
- VLAN and trunk analysis
- Native-VLAN, pruning, STP, and port-channel consistency
- Endpoint attachment and neighbor relationships
- Route and ACL context
- Wireless controller, AP, WLAN, profile, site, flex, and tag relationships
- Resource, environmental, software, stack, and counter trends
- Zabbix reconciliation
- Graylog enrichment and query context
- Draw.io-compatible generated topology
- Restricted read-only live collection only after offline and security-boundary acceptance

## Phase 5 — Firewall Ingestion and Traffic-Boundary Analysis

- FortiGate native parser and semantic normalization
- FortiGate supported YAML ingestion
- FortiGate operational and diagnostic evidence profiles
- OPNsense and pfSense XML normalization
- Interface, VLAN, zone, VDOM, route, gateway, policy, ACL-equivalent, object, NAT, VIP, VPN, and SD-WAN graph
- Policy order and traffic-boundary explanation
- Configured-versus-observed distinction
- Management and local-in exposure
- Runtime-health uncertainty
- Golden fixtures and adversarial parsing
- Correlation with Cisco-derived topology

## Phase 6 — Query, Projects, Changes, and Acceptance

- Global answer-first query
- IP, CIDR, subnet, and VLAN intelligence
- Reachability explanation
- Attack-path and trust-boundary context
- Dependency and blast-radius analysis
- Project portfolio
- Two-person and multi-approver policies
- Pre-change and post-change comparison
- Risk of approval and risk of denial
- Director-facing and engineering-facing change packages
- Expected and unexpected difference disposition
- Rollback and emergency change handling
- Formal acceptance and closure

## Phase 7 — Topology, Diagrams, and Interface Completion

- Normalized graph
- Switch-port, Layer-2, VLAN, route, SD-WAN, firewall, VPN, wireless, dependency, and attack-path views
- Draw.io-compatible generated sources
- Curated diagram lifecycle
- SVG and PDF publication
- Drift checks
- Accessible answer workspace and pivots

## Phase 8 — External Integrations

- Production Zabbix sender delivery
- Zabbix reconciliation, map, dashboard, template, discovery, and reporting assistance
- Graylog lookups, enrichment, queries, pipelines, streams, dashboards, and report assistance
- Security-platform asset and topology context
- BloodHound query-result and identity-context import
- Versioned Atlas OpenGraph extension and payload export
- Governed identity-to-asset correlation
- Combined identity, network, exposure, and change-impact analysis
- Optional least-privileged BloodHound API adapter only after offline acceptance
- OpenMetrics, syslog, webhook, and SIEM adapters
- Delivery outbox, retry, backpressure, and dead-letter handling
- Separately accepted provisioning boundaries where justified

## Phase 9 — Production Security, Recovery, and Representative Acceptance

- Signed builds and provenance
- SBOM and package integrity
- Off-host logging and integrity anchors
- Backup protection and restoration validation
- Break-glass
- Trusted rebuild and compromise recovery
- Representative-host deployment validation
- Controlled read-only pilot
- Operational acceptance based on evidence

## Acceptance Rule

No phase or slice is complete until:

- implementation;
- tests;
- requirements;
- architecture;
- limitations;
- evidence;
- validation;
- status;
- next work; and
- exact accepted commit

describe the same boundary.

Self-validation shall not be represented as independent review.
