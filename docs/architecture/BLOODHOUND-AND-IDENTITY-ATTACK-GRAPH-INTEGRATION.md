# BloodHound and Identity Attack-Graph Integration

## Purpose

Iron Atlas shall integrate with BloodHound and approved identity-graph evidence so network reachability, infrastructure exposure, identity privilege, and proposed change impact can be evaluated together.

BloodHound and SharpHound are not replacements for Atlas, and Atlas is not a replacement for BloodHound.

- SharpHound collects Active Directory relationship evidence for BloodHound.
- BloodHound models identity and privilege relationships and performs attack-path analysis.
- Atlas models network, firewall, wireless, routing, reachability, dependency, operational state, and governed change evidence.
- The integration correlates those models without erasing source ownership, uncertainty, or product boundaries.

The intended result is an answer such as:

> A principal can compromise a server through an identity relationship; the server is located in this subnet and VLAN; this network path and firewall policy permit the required service; this management boundary is exposed; and this proposed change would create, widen, narrow, or remove the combined path.

## Governing Principle

> **BloodHound explains identity privilege paths. Atlas explains network and control paths. Together they explain whether an identity-derived opportunity is reachable, exposed, operationally significant, and affected by change.**

Neither system silently becomes authoritative for the other system's domain.

## Integration Directions

### 1. BloodHound and SharpHound Evidence into Atlas

Atlas may ingest approved, bounded identity-graph context derived from:

- SharpHound collection artifacts;
- BloodHound query-result exports;
- BloodHound findings or saved-query results;
- BloodHound node and relationship exports;
- versioned Atlas-specific context bundles produced from approved BloodHound queries; and
- a separately accepted, least-privileged BloodHound API adapter where the deployed edition and documented interface support it.

The initial integration should prefer offline or manually exported evidence over direct database access.

Atlas shall not directly depend on BloodHound's internal PostgreSQL, Neo4j, or application schemas. Internal database access would create avoidable coupling, bypass application authorization, and make Atlas dependent on implementation details outside its control.

### 2. Atlas Network and Security Context into BloodHound

Atlas may generate a BloodHound OpenGraph extension definition and matching data payloads representing selected Atlas records and relationships.

The export may include supported, high-confidence context for:

- organizations and sites;
- network devices and managed endpoints;
- IP addresses and prefixes;
- VLANs and security zones;
- management planes and administrative services;
- routed and firewall-controlled reachability;
- external exposure;
- VPN and remote-access boundaries;
- jump hosts and management networks;
- critical infrastructure assets; and
- approved identity-to-asset correlations.

The OpenGraph export is a reviewable integration artifact. It does not make BloodHound the canonical Atlas store and does not prove that every exported relationship is currently exploitable.

### 3. Correlated Answers inside Atlas

Atlas shall combine identity-graph context with network and security evidence to answer questions such as:

- Which users, groups, computers, or service identities can reach network-device management services?
- Which identity paths terminate on systems exposed through a firewall, VPN, wireless network, or management VLAN?
- Which Tier Zero or otherwise critical identity assets are reachable from less-trusted network segments?
- Which compromised endpoint would provide the broadest network pivot opportunity?
- Which network controls prevent an identity path from becoming an operational attack path?
- Which firewall, ACL, route, VLAN, VPN, or SD-WAN change would create or remove a path to a domain controller, certificate authority, hypervisor, backup system, security platform, or management plane?
- Which privileged identities administer a device, and through which workstation, jump host, subnet, route, and policy is that administration possible?
- Which identity and network evidence is current, stale, incomplete, conflicting, or unsupported?

## Canonical Boundary

Atlas shall retain canonical records for network and security context. Imported BloodHound records remain external identity-graph observations with provenance.

BloodHound remains responsible for:

- Active Directory, Entra ID, and supported OpenGraph identity and privilege relationships;
- BloodHound-native attack-path semantics;
- BloodHound-native pathfinding, queries, privilege zones, findings, and remediation behavior; and
- compatibility between BloodHound and its supported collectors.

Atlas remains responsible for:

- network and firewall evidence;
- IP, CIDR, VLAN, interface, route, ACL, policy, NAT, VPN, SD-WAN, wireless, topology, and dependency models;
- configured-versus-observed state;
- network reachability and trust-boundary explanation;
- cross-system correlation;
- proposed-state and change-impact analysis;
- director-facing and engineering-facing decision packages; and
- Atlas evidence, uncertainty, approval, validation, and acceptance history.

## Identity Correlation

Cross-system identity linkage is security-sensitive and shall be explicit.

Preferred correlation evidence includes:

- stable directory object identifiers;
- security identifiers;
- tenant or domain identity;
- fully qualified names with source context;
- device certificates or other governed machine identity;
- approved asset identifiers;
- time-bounded IP and endpoint observations; and
- human-reviewed mappings.

Atlas shall not silently merge records based only on a short hostname, display name, mutable IP address, or other weak identifier.

Every correlation shall retain:

- source records;
- matching method;
- confidence;
- observation time;
- conflicts;
- reviewer disposition where required; and
- supersession history.

## OpenGraph Model Direction

An Atlas OpenGraph extension should use a dedicated namespace and versioned schema.

Candidate node kinds include:

- `IA_Site`;
- `IA_NetworkDevice`;
- `IA_ManagedEndpoint`;
- `IA_IPAddress`;
- `IA_Prefix`;
- `IA_VLAN`;
- `IA_SecurityZone`;
- `IA_ManagementPlane`;
- `IA_Service`;
- `IA_VPNBoundary`; and
- `IA_CriticalAsset`.

Candidate structural relationships include:

- contained by site;
- attached to interface or VLAN;
- assigned an address;
- contained by prefix;
- routed by device;
- protected by firewall or ACL;
- connected through VPN; and
- associated with a management plane.

Candidate attack-relevant relationships include:

- can reach administrative service;
- can authenticate to service;
- can administer device;
- exposes service to zone or prefix;
- can pivot to managed endpoint; and
- crosses trust boundary.

Structural relationships and attack-relevant relationships shall not be treated as equivalent.

A relationship shall be traversable in an attack graph only when its security meaning is explicit, evidence-supported, bounded, and tested. Atlas shall not turn every physical link, route, broad subnet relationship, or possible packet path into an attack edge.

## Current, Proposed, and Historical State

Current, proposed, and historical state shall remain distinguishable.

Atlas shall not upload proposed-state relationships into a production BloodHound graph in a way that can be mistaken for current truth.

Proposed-state analysis should use one of the following accepted boundaries:

- Atlas-native what-if analysis;
- a separately identified OpenGraph data source that cannot be confused with current state;
- an isolated BloodHound test instance or database; or
- a review-only export that is not ingested into the authoritative operational graph.

Change records shall identify whether an identity or network path is:

- present in accepted current state;
- observed in current operational evidence;
- calculated from current configuration;
- proposed;
- removed by the proposal;
- newly created by the proposal;
- uncertain; or
- conflicting.

## Collection and Security Boundary

SharpHound and BloodHound evidence may reveal highly sensitive identity, privilege, session, service, certificate, computer, and domain relationships.

Therefore:

- collection requires explicit authorization and declared scope;
- the compatible SharpHound release shall be selected for the deployed BloodHound version;
- SharpHound execution shall remain outside the Atlas server unless a separately accepted collector boundary justifies otherwise;
- collector credentials and tokens shall remain outside Git and retained logs;
- raw SharpHound archives and unredacted BloodHound exports shall remain outside Git;
- evidence shall receive classification, integrity, retention, and access controls;
- collection methods, target domains, hosts, timeouts, output limits, and stop conditions shall be bounded;
- live collection shall not be hidden inside a routine Atlas query;
- Atlas shall not expand collection scope automatically; and
- Atlas shall not represent BloodHound or SharpHound evidence as complete when the collector scope was partial.

## Initial Delivery Sequence

### Step 1 — Documentation and Contract

- Define the integration boundary.
- Define sensitive-data handling.
- Define stable identifiers and correlation rules.
- Define an Atlas identity-graph context bundle.
- Define required answer and change-impact use cases.

### Step 2 — Offline BloodHound Context Import

- Import sanitized fixtures and approved query-result bundles.
- Preserve BloodHound version, collector version, query identity, source time, and digest.
- Normalize only the identity context required by bounded Atlas questions.
- Keep unsupported nodes and relationships visible.

### Step 3 — Atlas OpenGraph Export

- Define the versioned Atlas OpenGraph extension schema.
- Export a small, reviewed set of Atlas nodes and edges.
- Validate structural and traversable relationship semantics.
- Test BloodHound ingestion using sanitized fixtures.
- Provide reviewable saved queries for the exported model.

### Step 4 — Cross-Domain Questions

Deliver at least these end-to-end answers:

1. Which identity principals can reach network-management services?
2. Which critical identity assets are exposed through supported network paths?
3. Which compromised computer can pivot into which protected segments?
4. Which proposed network change creates or removes an identity-to-infrastructure attack path?

### Step 5 — Governed Automation

Only after offline acceptance:

- evaluate a least-privileged BloodHound API adapter;
- define scheduled export and import behavior;
- add bounded retries, backpressure, and failure isolation;
- preserve per-run provenance and version compatibility; and
- prevent integration failure from blocking unrelated Atlas functions.

## Non-Goals

The initial integration shall not:

- replace BloodHound's identity attack graph;
- reimplement SharpHound;
- launch unrestricted domain collection;
- query BloodHound internal databases directly;
- treat every routed path as exploitable;
- silently create identity correlations;
- upload proposed state as current state;
- represent a BloodHound path as proof of packet reachability;
- represent packet reachability as proof of identity compromise; or
- automatically change Active Directory, firewall, switch, VPN, or BloodHound state.

## Acceptance Evidence

A bounded integration candidate requires:

- sanitized fixtures;
- exact BloodHound and collector compatibility records;
- schema and payload validation;
- deterministic correlation tests;
- false-match and ambiguous-match tests;
- graph-size and performance observations;
- sensitive-value and log-redaction tests;
- current-versus-proposed separation tests;
- manually verified example paths;
- documented false positives, false negatives, unknowns, and unsupported relationships; and
- exact source evidence for each Atlas conclusion.
