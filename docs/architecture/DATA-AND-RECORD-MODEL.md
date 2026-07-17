# Data and Record Model

## Purpose

Define the vendor-neutral records required to answer network, security, change, and evidence questions while preserving vendor detail and source lineage.

## Canonical Record Families

### Organization and Location

- Organization
- Site
- Building
- Floor
- Room
- Rack
- Geographic or logical location

### Device and Platform

- Device
- Logical device
- Stack or cluster
- Component
- Hardware
- Software
- Image and boot state
- License
- Lifecycle state
- HA relationship

### Interface and Link

- Physical interface
- Logical interface
- Subinterface
- SVI
- Loopback
- Tunnel
- Port channel or LAG
- Member interface
- Circuit
- Link
- Neighbor relationship
- Administrative and operational state

### Address and Prefix

- IP address
- Prefix or CIDR
- Subnet
- Host range
- Gateway
- Broadcast address where applicable
- Address family
- Prefix containment
- Longest-prefix candidate
- Overlap
- Duplicate
- Conflict
- Reserved or special-purpose classification
- Address object and group membership

### Layer 2

- VLAN
- Voice VLAN
- Access-port membership
- Trunk
- Native VLAN
- Allowed VLAN
- Active VLAN
- Pruned VLAN
- Broadcast domain
- MAC observation
- Endpoint attachment
- Port channel
- Spanning-tree instance, root, role, state, and topology-change observation

### Layer 3 and Routing

- VDOM
- VRF
- Routing domain
- Zone
- Connected route
- Static route
- Default route
- Dynamic route
- Policy route
- SD-WAN route or selection context
- Next hop
- Administrative preference
- Metric or priority
- Route installation observation
- Return-path relationship

### Security and Traffic Control

- Firewall address, service, schedule, user, group, and security-profile object
- Firewall policy and evaluation order
- ACL and attachment
- Local-in or management-plane control
- Source NAT
- Destination NAT
- VIP
- Address pool
- Port translation
- VPN
- Tunnel selector
- SD-WAN zone, member, rule, health check, SLA, preference, and fallback
- Inspection, logging, shaping, and authentication control

### Wireless

- Wireless controller
- HA relationship
- Access point
- Radio
- WLAN
- Policy profile
- Site tag
- Policy tag
- RF tag
- Flex profile
- Selected client observation

### Identity and Privilege Graph Context

- External identity principal
- Directory user, group, computer, service identity, and tenant or domain reference
- Identity privilege or control relationship
- BloodHound node, edge, path, finding, query, and privilege-zone reference
- Atlas asset correlation
- Correlation method, confidence, conflict, and human disposition
- Current, proposed, historical, and superseded identity-path state

### Evidence and Lineage

- Evidence bundle
- Artifact
- Content digest
- Source device or system
- Collection or import context
- Command or file source
- Parser run
- Analyzer run
- Parser and analyzer version
- Warning
- Unsupported state
- Redaction and classification state
- Evidence-to-record lineage

### Analysis

- Entity resolution
- Prefix containment result
- Route decision
- Policy decision
- NAT decision
- Traffic path
- Reachability result
- Identity attack-path context
- Network attack-path context
- Combined identity and network attack-path result
- Trust-boundary crossing
- Dependency
- Blast-radius result
- Difference
- Finding
- Observation
- Trend
- Exception
- Remediation
- Risk statement
- Confidence and uncertainty

### Project, Change, and Acceptance

- Project
- Requirement
- Current state
- Proposed state
- Change
- Approval policy
- Approval
- Denial
- Revision request
- Implementation step
- Validation step
- Rollback step
- Evidence snapshot
- Expected difference
- Unexpected difference
- Human disposition
- Acceptance
- Exception
- Supersession

### Integration and Telemetry

- Canonical metric
- Health event
- External-system identity
- External graph source, collector, extension, schema, and query identity
- Reconciliation result
- Generated recommendation
- Export artifact
- Delivery destination
- Outbox record
- Delivery attempt
- Applied-state evidence

## Relationships

Relationships are first-class records when they carry evidence, time, confidence, direction, order, or scope.

Examples include:

- address contained by prefix;
- VLAN present on interface;
- interface member of port channel;
- interface connected to neighbor;
- route selects next hop and egress;
- policy references object;
- policy permits or denies traffic;
- NAT translates one identity to another;
- VPN carries prefix;
- SD-WAN rule selects member;
- endpoint observed on switch port;
- directory principal controls or administers external asset;
- identity path correlates with network reachability;
- service depends on path;
- proposed change affects object;
- finding supported by evidence; and
- acceptance supersedes prior state.

## Identity

Every durable canonical record receives a stable identifier.

Vendor IDs, names, policy numbers, interface names, serial numbers, and object names are attributes and evidence. They are not the sole canonical identity.

Entity resolution retains:

- exact matches;
- aliases;
- source identities;
- merge decisions;
- split decisions;
- conflicts;
- confidence; and
- human disposition.

## Scope

Every record and relationship identifies applicable scope, including where relevant:

- organization;
- site;
- device;
- logical device;
- VDOM;
- VRF;
- routing domain;
- zone;
- address family;
- time; and
- accepted-state version.

## Time

Distinguish:

- source-device time;
- collector-received time;
- ingestion time;
- parser time;
- analyzer time;
- observation-valid time;
- proposal time;
- approval or denial time;
- implementation time;
- validation time;
- acceptance time; and
- supersession time.

## Evidence State

Material records and conclusions distinguish:

- `CONFIGURED`;
- `OBSERVED`;
- `CALCULATED`;
- `INFERRED`;
- `UNKNOWN`; and
- `CONFLICTING`.

Multiple evidence states may contribute to one answer.

## History

Material changes retain lineage.

The newest observation does not silently overwrite accepted state.

Atlas preserves:

- prior observation;
- current observation;
- prior accepted state;
- proposed state;
- actual post-change state;
- expected and unexpected differences;
- correction;
- supersession; and
- acceptance.

## Storage

PostgreSQL is the planned authoritative store for normalized and governed records.

Raw evidence remains in protected content-addressed storage and is referenced by immutable hashes and storage references.

Git contains code, contracts, schemas, sanitized fixtures, documentation, validation logic, and sanitized retained evidence only.

## Identity Attack-Graph Records

Atlas shall retain identity-graph context as sourced external evidence rather than silently converting it into Atlas authority.

Identity-graph records retain:

- external graph source and edition;
- collector and collector version;
- domain, tenant, and environment identity;
- principal and directory-object reference;
- privilege or control relationship reference;
- saved query, finding, privilege zone, or path result;
- Atlas asset correlation;
- correlation method and confidence;
- current, proposed, historical, and superseded state; and
- unsupported, ambiguous, stale, conflicting, or incomplete evidence.

Identity records imported from BloodHound do not become Atlas authorization records. Atlas governed actors and roles remain separate from observed directory privilege.

## Graph Relationship Semantics

Every graph relationship shall declare whether it is:

- structural;
- observational;
- control or policy;
- reachability;
- identity privilege;
- dependency;
- attack-relevant; or
- proposed.

Only explicitly modeled, evidence-supported attack-relevant relationships may be exported as traversable attack edges. Physical adjacency, a route, address containment, or broad network reachability shall not automatically imply compromise capability.
