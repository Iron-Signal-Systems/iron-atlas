# Target Architecture

## Status

Normative target direction. Accepted implementation boundaries remain governed by acceptance records and the implementation roadmap. This document describes the product architecture and must not be read as a claim that all layers are implemented or production-ready.

## Product Flow

```text
Configuration · operational commands · diagnostics · monitoring · logs
security telemetry · identity attack graphs · diagrams · documentation · accepted change records
                                  ↓
                   Authorized evidence ingestion
                                  ↓
        Protected raw evidence, provenance, digest, and lineage
                                  ↓
          Vendor parsers, collectors, and compatibility profiles
                                  ↓
              Canonical network and security record model
                                  ↓
    Identity · prefix · Layer 2 · Layer 3 · policy · NAT · VPN · SD-WAN
      wireless · topology · dependency · health · change · acceptance
                                  ↓
        Correlation, reachability, attack-path, and impact engines
                                  ↓
     Query answers · findings · diagrams · comparisons · change packages
                                  ↓
           HTML5 interface · versioned API · governed exports
```

## Architectural Layers

### 1. Evidence Sources

Initial sources include:

- Cisco configuration and operational command bundles;
- FortiGate configuration and operational or diagnostic output;
- OPNsense and pfSense configuration evidence;
- Zabbix inventory and monitoring metadata;
- Graylog syslog and SNMP-trap context;
- Security Onion and other approved security-platform context;
- BloodHound identity and privilege graph context and approved SharpHound-derived evidence;
- Draw.io and curated documentation;
- approved project and change records; and
- human verification and disposition.

### 2. Evidence Ingestion and Protection

The ingestion boundary provides:

- authorization;
- classification;
- source identity;
- collection or import context;
- digest and integrity;
- immutable raw-evidence reference;
- parser eligibility;
- duplicate and quarantine behavior;
- size, time, and resource bounds; and
- explicit rejection or partial-evidence state.

### 3. Vendor Adapters and Collectors

Adapters preserve vendor semantics without allowing vendor-specific representations to define the canonical model.

Cisco and FortiGate are the first major ingestion workstreams.

Collectors are read-only by default and operate through fixed, versioned, bounded profiles.

### 4. Canonical Network and Security Model

The canonical model represents:

- organizations, sites, buildings, rooms, racks, and locations;
- devices, components, software, licenses, and logical systems;
- interfaces, links, port channels, circuits, tunnels, and neighbors;
- MAC, ARP, endpoint, and attachment observations;
- VLANs, broadcast domains, subnets, prefixes, zones, VDOMs, VRFs, and routing domains;
- connected, static, dynamic, default, policy, and SD-WAN routes;
- firewall objects, policies, ACLs, schedules, identities, NAT, VIPs, VPNs, and inspection controls;
- wireless controllers, APs, WLANs, profiles, sites, flex profiles, tags, radios, and selected clients;
- external identity principals, directory objects, identity-graph references, correlation evidence, and privilege-path context;
- evidence, observations, findings, uncertainty, and human disposition;
- projects, proposed state, changes, approvals, denials, implementation, rollback, validation, and acceptance; and
- historical and superseded state.

### 5. Correlation and Intelligence Engines

The intelligence layer includes:

- entity resolution;
- IP and prefix containment;
- longest-prefix selection;
- overlap and conflict detection;
- Layer-2 relationship reconstruction;
- Layer-3 path reconstruction;
- policy, ACL, NAT, VPN, and SD-WAN correlation;
- reachability explanation;
- return-path and asymmetry analysis;
- attack-path and trust-boundary context;
- dependency and blast-radius analysis;
- current-to-prior comparison;
- current-to-proposed comparison;
- risk-of-change and risk-of-no-change analysis; and
- evidence and confidence aggregation.

### 6. Governed Application Services

Application services provide:

- answer-first query;
- inventory and topology;
- findings and disposition;
- projects and changes;
- approval and denial;
- pre-change and post-change comparison;
- validation and rollback;
- formal acceptance;
- report and diagram generation; and
- external-system reconciliation and exports.

### 7. Interface and API

The user interface and API expose one consistent model.

The interface is:

- search- and query-first;
- pivot-driven;
- keyboard-operable;
- accessible;
- evidence-aware;
- explicit about stale, incomplete, uncertain, conflicting, queued, failed, proposed, approved, denied, implemented, validated, and accepted state; and
- designed to minimize navigation work.

## Process Direction

The target deployment separates:

- `atlas-api` — HTML5, API, query, and governed workflow service;
- `atlas-worker` — parsing, correlation, analysis, diagram, comparison, and delivery jobs;
- `atlas-ingest` — authenticated evidence intake and sequencing;
- `atlas-collector` — site-scoped read-only collection;
- PostgreSQL — authoritative normalized and governed records; and
- protected evidence storage — encrypted raw evidence addressed by content hash.

The current implementation may combine processes while boundaries are being established. Combined deployment does not erase logical trust or responsibility boundaries.

## Evidence State

Every material conclusion carries one or more evidence states:

- `CONFIGURED`;
- `OBSERVED`;
- `CALCULATED`;
- `INFERRED`;
- `UNKNOWN`; and
- `CONFLICTING`.

A configuration backup proves configured intent. It does not prove current link state, route installation, VPN state, SD-WAN member health, session state, monitoring success, or actual traffic.


## BloodHound and Identity-Graph Direction

Atlas integrates with BloodHound through explicit import and export boundaries.

- Atlas may import approved BloodHound query results, findings, graph exports, or versioned context bundles.
- Atlas may generate a versioned OpenGraph extension and payloads containing selected network, management-plane, exposure, and reachability context.
- Atlas shall not depend on BloodHound internal database schemas.
- Identity-to-asset correlation requires stable identifiers, confidence, provenance, and conflict handling.
- Structural topology edges and traversable attack edges remain distinct.
- Proposed Atlas state shall not be mixed with current BloodHound truth.
- Raw SharpHound and unredacted BloodHound evidence remain protected evidence and are prohibited from Git.

The normative integration model is defined in [BloodHound and Identity Attack-Graph Integration](BLOODHOUND-AND-IDENTITY-ATTACK-GRAPH-INTEGRATION.md).

## External-System Direction

External adapters operate through governed import and export boundaries.

```text
External evidence or metadata
            ↓
Atlas evidence boundary
            ↓
Canonical records and analysis
            ↓
Reviewable recommendation, context, query, map, report, or export
            ↓
Optional separately governed external application
            ↓
Post-application validation evidence
```

External systems retain responsibility for their mature operational functions. Atlas retains responsibility for evidence lineage, normalized identity, cross-system correlation, topology, reachability explanation, governed findings, change history, and acceptance history.

## Non-Negotiable Boundaries

- Collectors do not receive unrestricted PostgreSQL access.
- Raw protected evidence does not enter Git.
- Vendor-specific records do not become the canonical model.
- Monitoring or logging products do not become authorization sources of truth.
- External-system records do not silently become authoritative Atlas state.
- UI visibility is never treated as authorization.
- Identity-provider claims never directly become Atlas authority.
- A requester cannot independently approve the requester’s own governed infrastructure change.
- Parser and analyzer uncertainty remains visible.
- Configured intent is not presented as observed state.
- Calculated reachability is not presented without assumptions and evidence.
- Generated diagrams never overwrite curated diagrams.
- Atlas does not silently modify infrastructure.
- Future writes require a separately accepted, previewable, attributable, bounded, approval-aware, and validated boundary.
- Passing self-validation does not create independent review.
