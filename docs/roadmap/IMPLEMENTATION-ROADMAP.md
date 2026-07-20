# Implementation Roadmap

## Current Primary Focus

Iron Atlas is the primary active product-development effort.

The implementation roadmap preserves accepted foundation boundaries while delivering bounded vertical product slices that answer real network and security questions as early as safely possible.

> **Cisco, FortiGate, and BloodHound-derived identity context are the first core evidence sources. The product is the correlated answer, risk analysis, and decision evidence produced from them.**

Implementation may proceed through compatible parallel workstreams, but formal acceptance follows the phase order below.

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
- BloodHound-derived identity privilege and critical-asset context where evidence permits;
- governed identity-to-asset correlation;
- identity privilege paths;
- network reachability paths;
- combined identity and network attack paths;
- explicit separation of identity capability, network capability, and combined path;
- unresolved, ambiguous, stale, incomplete, and conflicting evidence;
- step-by-step supporting evidence; and
- evidence-backed result.

The BloodHound capability is part of the completed Atlas intelligence model. A query may return no identity context when approved BloodHound evidence is unavailable, incomplete, stale, or inapplicable. That absence shall remain explicit and shall not be represented as proof that no identity path exists.

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

---

## Phase 0 — Repository and Executable Baseline

**Status:** Accepted non-production baseline.

- Mission, architecture, requirements, and boundaries
- Go module and embedded HTML5 interface
- RBAC and requester-approval independence candidate
- Module registry and parser contracts
- Initial evidence-source candidates
- Unit tests and repository validation
- Minimal Arch Linux deployment material

Acceptance does not claim production authentication, persistence, collection, semantic completeness, or production readiness.

Historical Phase 0 validators and acceptance evidence remain unchanged.

---

## Phase 1 — PostgreSQL Foundation and Governed Identity

**Status:** Steps 1 and 2 accepted; Step 3 authenticated server-side session is the active bounded trusted-authentication candidate. Authentication assurance, MFA policy, optional local TOTP, session lifecycle, CSRF, trusted proxies, production wiring, and remaining production-foundation boundaries are incomplete.

### Accepted and active boundaries

- Manifest-driven migrations
- Production role and ownership topology
- Governed actors and external identities
- Role binding, authority, approval, decision, and audit records
- Database-enforced requester and approver independence
- Append-only history
- Disposable database and concurrency tests
- Replaceable PostgreSQL change-service implementation
- Least-privileged bounded connection pool
- Transaction-local actor context
- Persistent change creation and approval
- Rollback and pooled-connection isolation
- PostgreSQL-aware readiness
- Trusted production authentication
- Governed actor resolution
- Bounded sessions, cookies, CSRF, logout, and trusted proxies
- Provider-neutral authentication assurance, phishing-resistant MFA, and governed RFC 6238 TOTP fallback
- Credential delivery and rotation
- PostgreSQL TLS
- Foundational backup and recovery
- Production connection and resource budgets

Offline evidence, parser, canonical-model, fixture, and documentation work may continue in isolated workstreams without representing live collection or production authority as accepted.

Historical Phase 1 checkpoint and acceptance validators remain unchanged. Later work shall extend the accepted sequence rather than relabeling earlier history.

---

## Phase 2 — Evidence Intake, Protection, and Storage

Phase 2 establishes the vendor-independent evidence boundary used by Cisco, FortiGate, BloodHound-derived context, and future evidence sources.

### Step 2.1 — Evidence-bundle and provenance contract

- Versioned evidence-bundle contracts
- Source-system identity
- Collection and observation timestamps
- Parser and schema versions
- Evidence digests
- Required, optional, unsupported, failed, and unavailable records
- Deterministic validation
- Explicit configured-versus-observed distinction

### Step 2.2 — Intake and quarantine

- Manual intake
- Authenticated intake seam
- Structural validation
- Size and count limits
- Rejection classifications
- Quarantine
- No parser execution for rejected evidence

### Step 2.3 — Durable staging and duplicate detection

- Durable staging
- Transaction behavior
- Duplicate identification
- Replay handling
- Safe retry
- Interruption recovery

### Step 2.4 — Protected content-addressed storage

- Encrypted content-addressed storage
- Signed evidence bundles
- Integrity verification
- Immutable raw evidence
- Least-privileged evidence access
- Protected evidence metadata

### Step 2.5 — Classification and redaction

- Evidence classification
- Redaction
- Protected-field handling
- Original-versus-redacted lineage
- Sanitized logs and retained validation evidence

### Step 2.6 — Parser isolation and resource governance

- Parser isolation
- Cancellation and timeout
- CPU and memory bounds
- Storage and concurrency limits
- Oversized-input handling
- Parser-failure containment

### Step 2.7 — Recovery, corruption, and hostile evidence

- Backup and restoration
- Digest-mismatch detection
- Corruption handling
- Truncated and conflicting evidence
- Replay and duplication
- Malformed and adversarial bundles

### Phase 2 acceptance

Phase 2 is accepted only when evidence can be received, protected, stored, recovered, processed, and traced without relying on any particular vendor parser.

---

## Phase 3 — Cisco Offline Evidence and Normalization

Phase 3 establishes Cisco support around operating-system families, command behavior, device roles, and capability profiles rather than individual hardware models.

### Step 3.1 — Cisco evidence profiles and operating-system identity

- Sanitized versioned command and configuration bundles
- Device identity
- Operating-system family identification
- Device-role identification
- Software-release identification
- Command provenance
- Required, optional, unsupported, failed, and unavailable command records
- Partial, truncated, malformed, conflicting, and stale evidence handling

### Step 3.2 — Classic Cisco IOS normalization

Establish the parser and compatibility boundary for monolithic Classic IOS platforms.

- Running and startup configuration evidence
- Version, inventory, boot, licensing, and platform identity
- Interfaces and switch-port state
- VLANs, trunks, native VLANs, and allowed VLANs
- Port channels
- CDP and LLDP where supported
- Spanning Tree Protocol state
- MAC-address and ARP evidence
- Layer-3 interfaces
- Routes and next hops
- Access control lists
- Classic IOS syntax variations and legacy-output handling

Applicable Catalyst 2960 and similar systems are represented through compatibility profiles beneath the Classic IOS operating-system boundary.

### Step 3.3 — Cisco IOS XE normalization

Establish the parser and compatibility boundary for modular and programmable IOS XE platforms.

- Running and startup configuration evidence
- Version, package, install-mode, boot, licensing, and platform identity
- Chassis, stack, StackWise, supervisor, and member relationships
- Interfaces and switch-port state
- VLANs, trunks, native VLANs, pruning, and allowed VLANs
- Port channels and bundled-interface state
- CDP and LLDP
- Spanning Tree Protocol state
- MAC-address and ARP evidence
- Layer-3 interfaces
- VRFs and routing domains
- Routes and next hops
- Access control lists
- Software, process, resource, environmental, and counter evidence
- Structured output where safely and reliably supported
- IOS XE syntax, release, feature, and platform-capability variation

Applicable Catalyst 9200, 9300, 9500, and similar systems are represented through compatibility profiles beneath the IOS XE operating-system boundary.

### Step 3.4 — IOS XE wireless-controller role profile

Wireless-controller support remains an IOS XE role profile rather than a separate operating-system family.

- Controller and redundancy identity
- Access points
- WLANs
- Policy profiles
- Site, policy, and RF tags
- Flex profiles
- AP join and operational state
- Client and attachment relationships where evidence permits
- Configured-versus-observed wireless state
- Unsupported and incomplete relationship reporting

Applicable Catalyst 9800 systems are tracked as IOS XE wireless-controller compatibility profiles.

### Step 3.5 — Cross-OS canonical normalization

Normalize Classic IOS and IOS XE evidence into common Atlas records while preserving meaningful source-specific differences.

- Devices and logical devices
- Chassis, stacks, members, and supervisors
- Interfaces and subinterfaces
- VLANs and switched domains
- Trunks and allowed-VLAN state
- Port channels
- Neighbors
- Spanning-tree instances and roles
- MAC and ARP observations
- Layer-3 interfaces
- VRFs and routing domains
- Routes and next hops
- Access control lists
- Wireless entities and relationships
- Software, health, and diagnostic observations
- Evidence provenance, confidence, uncertainty, and unsupported state

### Step 3.6 — Capability and compatibility profiles

Each compatibility profile records:

- operating-system family;
- platform family;
- software release;
- device role;
- supported command set;
- supported normalized records;
- unsupported or degraded capabilities;
- known syntax variations;
- sanitized fixtures;
- parser-test coverage;
- validation status; and
- exact accepted commit.

A platform profile may be added without creating a new architectural phase when it remains within an accepted operating-system and normalization contract.

### Step 3.7 — Adversarial and resource validation

- Malformed and hostile command output
- Truncated command output
- Paging and terminal artifacts
- Duplicate and reordered sections
- Conflicting configuration and operational evidence
- Unknown software releases
- Unsupported commands
- Oversized evidence
- Cancellation and timeout
- CPU and memory governance
- Deterministic output
- Golden fixtures for each accepted profile
- Regression across Classic IOS and IOS XE

### Phase 3 acceptance

Phase 3 acceptance establishes offline Cisco evidence and normalization only. It does not establish semantic topology conclusions, unrestricted live collection, device modification, or production readiness.

---

## Phase 4 — FortiGate Offline Evidence and Normalization

Phase 4 establishes the first firewall evidence and normalization boundary.

### Step 4.1 — Native configuration ingestion

- FortiGate native configuration parsing
- Version and platform identity
- VDOM and virtual-system identity
- Provenance and source location
- Partial and unsupported sections

### Step 4.2 — Supported YAML ingestion

- Supported FortiGate YAML forms
- Bounded maintained decoding
- Native-layout detection
- Explicit unsupported layouts
- Configured-versus-observed distinction

### Step 4.3 — Native and YAML equivalence

- Canonical equivalence where support is claimed
- Explicit format-specific differences
- Deterministic normalized output
- No equivalence claim for unsupported or ambiguous semantics

### Step 4.4 — Interfaces, VDOMs, zones, and objects

- Interfaces and VLANs
- Zones and VDOM boundaries
- Address and service objects
- Object groups
- Unresolved and ambiguous references

### Step 4.5 — Routes, policies, NAT, and VIPs

- Routes and gateways
- Policy order
- Policy references
- NAT
- VIPs
- Disabled-versus-absent configuration

### Step 4.6 — VPN, SD-WAN, management, and local-in

- VPN relationships
- SD-WAN members and rules
- Management exposure
- Local-in policy
- Trust-boundary context

### Step 4.7 — Operational evidence and runtime uncertainty

- Operational and diagnostic evidence profiles
- Runtime-health uncertainty
- Configured-versus-observed state
- Stale and unavailable runtime evidence

### Step 4.8 — Adversarial and resource validation

- Malformed and hostile input
- Oversized and excessive structures
- Duplicate keys and unsupported YAML features
- Truncation and conflicting evidence
- Cancellation, timeout, CPU, and memory governance
- Sanitized golden fixtures

### Phase 4 acceptance

Phase 4 is accepted when supported FortiGate evidence can be deterministically normalized into the Atlas canonical model with policy order, provenance, uncertainty, and unsupported behavior preserved.

OPNsense, pfSense, and other firewall platforms may later be added as compatibility extensions. They are not prerequisites for initial core Atlas acceptance.

---

## Phase 5 — BloodHound Identity Context and Asset Correlation

BloodHound-derived identity context is part of the core Atlas product boundary.

Iron Atlas does not replace BloodHound or independently recreate BloodHound's identity attack graph. Atlas imports bounded, approved BloodHound context and correlates it with network, infrastructure, exposure, and change evidence.

### Step 5.1 — Offline BloodHound evidence contract

- Versioned offline BloodHound query-result bundles
- Approved SharpHound-derived evidence records
- Source-system and collection provenance
- Schema and parser versions
- Evidence digest and integrity
- Classification and redaction
- Malformed, oversized, incomplete, and conflicting input handling

### Step 5.2 — Identity and privilege-path normalization

- Directory principals
- Groups and computers
- Critical-asset context
- Identity privilege-path references
- Source-system identifiers
- Evidence freshness and uncertainty

### Step 5.3 — Atlas asset-correlation candidates

- Atlas assets
- Correlation candidates
- Matching evidence
- Confidence
- Ambiguity and conflict
- No silent promotion of uncertain matches

### Step 5.4 — Governed correlation decisions

- Accepted correlations
- Rejected correlations
- Ambiguous correlations
- Actor, decision, and history records
- Reviewable evidence
- Reversible decisions

### Step 5.5 — Correlation freshness and history

- Stale identity evidence
- Superseded correlations
- Conflicting source identities
- Correlation history
- Re-evaluation behavior

### Step 5.6 — Atlas OpenGraph identity extension

- Versioned Atlas OpenGraph extension
- Reviewable import and export payloads
- Network and identity relationships
- Evidence references
- Uncertainty and unsupported state
- Source-system lineage

### Step 5.7 — Privacy, hostile input, and resource governance

- Protected identity evidence
- Redacted retained output
- Hostile and malformed records
- Oversized graphs
- Cancellation and resource limits

### Phase 5 acceptance

Phase 5 is accepted when Atlas can safely import approved identity context, correlate it with Atlas assets, preserve ambiguity, and explain what the evidence does and does not prove.

Direct BloodHound API access and automated collector orchestration remain optional modules. Offline evidence support is the required core boundary.

---

## Phase 6 — Cross-Source Canonical Graph and Semantic Analysis

Phase 6 combines accepted Cisco, FortiGate, and BloodHound-derived records into a common evidence-backed model.

### Step 6.1 — Graph identity and lifecycle

- Stable canonical identities
- Source references
- Graph versions
- Replacement and supersession behavior
- Deterministic construction

### Step 6.2 — Layer-2 graph

- VLANs and switched domains
- Interfaces and trunks
- Port channels
- Neighbors and endpoint attachment
- Spanning-tree relationships

### Step 6.3 — Layer-3 graph

- Subnets and prefixes
- Layer-3 interfaces
- VRFs, VDOMs, and routing domains
- Routes and next hops
- Gateway relationships

### Step 6.4 — Firewall and traffic-boundary graph

- Zones and trust boundaries
- Policies and ACLs
- Address and service objects
- NAT and VIPs
- VPN and SD-WAN relationships

### Step 6.5 — Identity and critical-asset graph

- Identity principals
- Critical assets
- Governed identity-to-asset relationships
- Identity privilege references
- Explicit identity-versus-network separation

### Step 6.6 — State and time model

- Prior, current, proposed, and post-change state
- Configured and observed state
- Stale and superseded state
- Collection-time and effective-time distinctions

### Step 6.7 — Conflict, uncertainty, and provenance

- Conflicting evidence
- Unsupported state
- Unknown state
- Confidence
- Source-specific provenance
- No invented relationships

### Step 6.8 — Semantic consistency analysis

- VLAN and trunk consistency
- Native-VLAN and pruning analysis
- Spanning-tree relationships
- Port-channel consistency
- Route-selection context
- Policy and ACL context
- NAT relationships
- Configured-versus-observed reconciliation
- Identity and network relationship boundaries

### Phase 6 acceptance

Phase 6 is accepted when Atlas can construct a deterministic, evidence-backed, cross-source graph without erasing uncertainty or inventing relationships.

---

## Phase 7 — Query, Reachability, and Attack-Path Intelligence

Phase 7 establishes the primary answer-first product capability.

### Step 7.1 — IP, CIDR, subnet, and VLAN intelligence

Atlas answers, where evidence permits:

- where an address, prefix, or VLAN exists;
- which prefix contains an address;
- longest-prefix match;
- overlap and conflict;
- switch, interface, trunk, and VLAN placement;
- gateway and routing-domain context;
- route selection;
- policy and ACL references;
- NAT, VIP, VPN, and SD-WAN relationships;
- dependencies;
- supporting evidence; and
- unknown or unsupported state.

### Step 7.2 — Route and policy explanation

- Selected route and why
- Candidate and rejected routes
- Policy and ACL evaluation
- Object and service resolution
- Evidence and uncertainty

### Step 7.3 — Forward reachability

- Source and destination context
- Protocol and port
- Layer-2 path
- Layer-3 path
- Route selection
- Policy and ACL evaluation
- NAT, VPN, and SD-WAN
- Trust-boundary crossings

### Step 7.4 — Return path and asymmetry

- Return route
- Return policy
- Reverse translation
- Asymmetric-path uncertainty
- Missing and stale evidence

### Step 7.5 — Identity privilege-path explanation

- Identity privilege paths
- Critical-asset context
- Governed correlations
- Evidence freshness
- Ambiguous and unresolved paths

### Step 7.6 — Combined identity and network attack path

- Identity capability without network capability
- Network capability without identity privilege
- Combined capability
- Trust-boundary crossings
- Step-by-step evidence

### Step 7.7 — Evidence-backed answer explanation

Every answer distinguishes:

- facts;
- inferred relationships;
- assumptions;
- unknowns;
- stale evidence;
- conflicts;
- unsupported behavior; and
- source evidence.

### Step 7.8 — Accuracy and hostile-query campaign

- Verified positive and negative cases
- False-positive accounting
- False-negative accounting
- Malformed and excessive queries
- Cancellation and timeout
- Query resource limits
- Deterministic results

### Phase 7 acceptance

Phase 7 is accepted when Atlas can produce defensible answers with explicit evidence, uncertainty, unsupported behavior, and reproducible reasoning.

Absence of evidence shall never be represented as proof of absence.

---

## Phase 8 — Projects, Change Impact, and Decision Support

Phase 8 turns Atlas intelligence into governed engineering and leadership decisions.

### Step 8.1 — Project and proposed-state model

- Project portfolio
- Proposed-state records
- Scope and ownership
- Current-to-proposed linkage

### Step 8.2 — Change approval governance

- Two-person and multi-approver policies
- Requester and approver independence
- Conflicting and stale decisions
- Governed emergency handling

### Step 8.3 — Difference and dependency analysis

- Prior, current, proposed, and post-change comparison
- Expected and unexpected differences
- Dependency changes
- Paths created, removed, or altered

### Step 8.4 — Blast radius and risk

- Network blast radius
- Identity blast radius
- Risk of approval
- Risk of denial or delay
- Evidence and assumptions

### Step 8.5 — Engineering and leadership packages

- Engineering implementation plan
- Validation plan
- Rollback plan
- Director-facing decision summary
- Evidence and limitations

### Step 8.6 — Post-change validation and closure

- Post-change evidence
- Expected and unexpected difference disposition
- Rollback or acceptance
- Formal closure
- Retained decision evidence

### Phase 8 acceptance

Phase 8 is accepted when a proposed change can be evaluated, approved, implemented, validated, compared, rolled back, and formally closed using evidence-backed records.

---

## Phase 9 — Topology, Diagrams, and Accessible Interface

Phase 9 completes the core human-facing product experience.

### Step 9.1 — Answer-first workspace

- Global answer-first query
- Evidence inspection
- Uncertainty and unsupported-state presentation
- Keyboard-first operation
- Accessible navigation and pivots

### Step 9.2 — Network and security views

- Switch-port and Layer-2 views
- VLAN and trunk views
- Layer-3 and route views
- Firewall and trust-boundary views
- VPN and SD-WAN views
- Wireless views
- Dependency views
- Identity and attack-path views
- Change-impact views

### Step 9.3 — Generated and curated diagrams

- Draw.io-compatible generated sources
- Curated diagram lifecycle
- Diagram provenance
- Stable identifiers
- Drift detection

### Step 9.4 — Publication

- SVG publication
- PDF publication
- Evidence and version references
- Accessible generated output

### Step 9.5 — Accessibility and interface governance

- WCAG 2.1 Level AA target behavior
- Keyboard navigation
- Screen-reader behavior
- Visual and non-color-only status
- Interface resource budgets
- Bounded rendering and query behavior

### Phase 9 acceptance

Phase 9 is accepted when users can obtain, understand, navigate, export, and verify Atlas answers without manually browsing disconnected vendor records.

---

## Phase 10 — Restricted Read-Only Collection and Evidence Refresh

Phase 10 introduces live source interaction only after offline evidence, parser, normalization, and semantic boundaries are accepted.

### Step 10.1 — Collector security foundation

- Read-only collection contracts
- Source allowlisting
- Command allowlisting
- Least-privileged credentials
- Credential isolation
- Host-key and certificate verification
- Explicit prohibition on source modification

### Step 10.2 — Source collectors

Initial core collection scope may include:

- supported Classic IOS sources;
- supported IOS XE sources; and
- supported FortiGate sources.

Each source collector remains bounded by accepted offline evidence contracts.

### Step 10.3 — Source protection and resource governance

- Bounded concurrency
- Cancellation and timeout
- Rate limiting
- Collection windows
- Source-system load protection
- Partial collection handling
- Stale-evidence handling

### Step 10.4 — Reproducible collected evidence

- Evidence-bundle generation
- Offline replay
- Deterministic normalization
- Complete collection audit
- Failure isolation

### Step 10.5 — Hostile source conditions

- Slow and interrupted responses
- Authentication failure
- Changed host keys and invalid certificates
- Unexpected prompts
- Paged and oversized output
- Command refusal
- Source reboot
- Concurrency pressure
- Attempted command escalation

### Phase 10 acceptance

Phase 10 is accepted when Atlas can collect bounded read-only evidence without modifying source systems or weakening offline evidence reproducibility.

BloodHound offline evidence remains sufficient for core acceptance. Direct BloodHound API access is an optional module.

---

## Phase 11 — Production Security, Recovery, and Representative Deployment

Phase 11 establishes the production operating boundary.

### Step 11.1 — Build and supply-chain integrity

- Reproducible builds
- Signed builds
- Build provenance
- SBOM
- Dependency and package integrity

### Step 11.2 — Runtime hardening

- Hardened service identities
- Host and service isolation
- Secret delivery and rotation
- PostgreSQL TLS
- Host and network hardening

### Step 11.3 — Logging and integrity

- Off-host logging
- Integrity anchors
- Protected audit and validation evidence
- Failure visibility

### Step 11.4 — Backup, restore, and break-glass

- Protected backups
- Restoration validation
- Break-glass lifecycle
- Expiry and revocation
- Complete audit

### Step 11.5 — Trusted rebuild and compromise recovery

- Trusted rebuild
- Compromise recovery
- Credential and key replacement
- Evidence integrity review
- Recovery documentation

### Step 11.6 — Representative deployment and failure behavior

- Representative-host deployment
- Service startup and boot recovery
- Storage and network failure behavior
- Resource budgets
- Performance observations
- Upgrade and rollback

### Phase 11 acceptance

Phase 11 is accepted when Atlas can be securely installed, started, operated, backed up, restored, rebuilt, upgraded, and recovered on representative production hardware.

---

## Phase 12 — Controlled Pilot and Operational Acceptance

Phase 12 proves whether Atlas provides useful and defensible answers in a representative authorized environment.

### Step 12.1 — Pilot authorization and baseline

- Explicitly authorized environment
- Read-only restrictions
- Approved evidence sources
- Environment baseline
- Privacy and handling rules
- Removal and rollback plan

### Step 12.2 — Pilot deployment and source safety

- Exact deployed build
- Signed build provenance
- Source-system impact controls
- Complete operational logging
- No infrastructure modification

### Step 12.3 — Answer-verification campaign

Manually verify:

- asset identity;
- IP and prefix;
- VLAN placement;
- routing;
- policy;
- NAT;
- reachability;
- identity correlation; and
- attack paths.

### Step 12.4 — Accuracy accounting

- False positives
- False negatives
- Unsupported answers
- Unknown answers
- Stale evidence
- Correlation ambiguity

### Step 12.5 — Evidence and operational quality

- Evidence completeness
- Evidence lineage
- Resource impact
- Administrator time saved
- Decision quality
- Engineering package usefulness
- Leadership package usefulness

### Step 12.6 — Recovery and residual risk

- Restoration exercise
- Operational recovery
- Residual-risk record
- Pilot limitations
- Explicit statement of what the pilot does not prove

### Phase 12 acceptance

Phase 12 is the formal core Iron Atlas operational-acceptance boundary.

Core operational acceptance must remain valid when no optional integration modules are installed.

---

## Optional Post-Core Modules

Optional modules are maintained outside the mandatory Phase 0–12 acceptance sequence.

They may be added after the core product establishes the necessary architecture and extension contracts.

Optional modules include:

- Zabbix integration;
- Graylog integration;
- other SIEM integrations;
- OpenMetrics export;
- syslog delivery;
- webhook delivery;
- optional read-only APIs;
- automated BloodHound API access;
- additional Cisco compatibility profiles;
- OPNsense compatibility;
- pfSense compatibility;
- additional firewall vendors;
- other network vendors; and
- separately governed provisioning capabilities.

Each optional module requires its own:

- requirements;
- architecture;
- threat model;
- privilege boundary;
- data-classification rules;
- configuration contract;
- credentials;
- failure isolation;
- resource limits;
- retry and backpressure behavior where applicable;
- compatibility matrix;
- tests;
- validation evidence;
- limitations;
- release lifecycle;
- acceptance record; and
- exact accepted commit.

Delivery outboxes, retry behavior, backpressure, and dead-letter handling are required only for modules that perform outbound or asynchronous delivery.

No core phase gate may require every optional module to exist, be installed, or pass validation.

Failure or absence of an optional module must not prevent core Iron Atlas operation.

---

## Phase and Gate Rule

Each phase uses:

1. a phase-entry contract gate;
2. bounded implementation gates;
3. one phase integration gate; and
4. one formal phase-acceptance gate.

A workstream has:

- one accepted predecessor;
- one declared candidate;
- one scope;
- one validation boundary;
- one evidence set; and
- one next step.

Only one acceptance candidate may be active within a workstream.

Cross-workstream behavior is accepted through an explicit integration candidate.

Historical gates remain frozen and are revalidated at their exact accepted commits rather than weakened for later repository state.

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
