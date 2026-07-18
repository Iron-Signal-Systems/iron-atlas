# Query, Reachability, and Change-Impact Model

## Purpose

Define how Atlas converts multi-source infrastructure evidence into fast, defensible answers for network operations, security operations, leadership, and governed change.

## Query Subjects

Atlas queries may begin with:

- IP address;
- CIDR or subnet;
- VLAN;
- MAC address;
- hostname or endpoint;
- site;
- device;
- interface;
- zone;
- VDOM;
- VRF or routing domain;
- route or next hop;
- firewall policy;
- ACL;
- NAT or VIP;
- VPN or tunnel;
- SD-WAN rule, zone, member, or health check;
- protocol, service, or port;
- source and destination;
- finding;
- dependency;
- change;
- BloodHound principal, computer, group, privilege zone, finding, or attack path; or
- evidence artifact.

- incident;
- compromised user, account, endpoint, server, certificate, application,
  service, network device, or file;
- containment action; or
- Iron File Intelligence object, classification, or activity reference.

## Answer Contract

Every answer returns:

- subject identity;
- applicable scope;
- answer summary;
- supporting relationships;
- evidence;
- evidence state;
- confidence;
- assumptions;
- conflicts;
- unsupported areas;
- age;
- unknowns; and
- available pivots.

Atlas shall not present a calculated or inferred answer as observed fact.

## IP and CIDR Intelligence

Given an IP address or CIDR, Atlas should determine where evidence permits:

- address family;
- canonical representation;
- exact object matches;
- containing prefixes;
- longest-prefix match;
- subnet mask;
- network and broadcast address where applicable;
- usable range;
- gateway candidates;
- duplicate or overlapping prefixes;
- VLAN and Layer-2 domain;
- interfaces, SVIs, zones, VDOMs, VRFs, and routing domains;
- routes;
- address groups and policy use;
- NAT, VIP, VPN, and SD-WAN relationships;
- endpoint attachment;
- dependencies;
- exposure;
- changes; and
- findings.

## VLAN Intelligence

Given a VLAN ID or name, Atlas should determine:

- sites and devices where it exists;
- access interfaces;
- voice-VLAN use;
- trunks carrying it;
- native-VLAN use;
- allowed, active, and pruned state;
- port channels;
- spanning-tree instances, roots, roles, and state;
- SVI or gateway;
- associated prefixes;
- firewall zone or routed boundary;
- wireless use;
- endpoint observations;
- mismatches;
- dependencies; and
- evidence age.

## Route Selection

Route analysis preserves scope and vendor behavior.

A route answer considers:

1. source routing domain;
2. destination prefix;
3. connected routes;
4. static routes;
5. dynamic routes;
6. policy routes;
7. SD-WAN rules;
8. longest-prefix match;
9. administrative preference, distance, priority, and metric;
10. next-hop reachability;
11. route installation where observed;
12. egress interface or zone;
13. fallback;
14. tunnel state;
15. return path; and
16. uncertainty.

Configured routes and observed forwarding state are distinguished.

## Reachability Query

A reachability query includes, where known:

- source address, object, subnet, endpoint, interface, zone, site, or identity;
- destination address, object, subnet, endpoint, interface, zone, site, or identity;
- protocol;
- source port;
- destination port;
- time or schedule;
- application or internet-service context;
- routing domain;
- current, prior, proposed, or accepted state; and
- operational evidence time.

## Reachability Evaluation

Atlas evaluates:

1. source identity and attachment;
2. source prefix, VLAN, gateway, and routing domain;
3. ingress interface and zone;
4. Layer-2 forwarding relationships where required;
5. policy and ACL order;
6. object, group, service, schedule, identity, and profile expansion;
7. policy route or route-map behavior;
8. SD-WAN match, strategy, SLA, preference, and fallback;
9. longest-prefix route and next-hop reachability;
10. egress interface, zone, tunnel, or member;
11. source and destination translation;
12. inspection, logging, shaping, and authentication controls;
13. destination attachment or terminal network;
14. return path;
15. asymmetry;
16. runtime-health dependencies;
17. conflicts and unsupported behavior; and
18. evidence for every step.

## Reachability Result

The result is not limited to yes or no.

Permitted outcomes include:

- `PERMITTED`;
- `DENIED`;
- `NO_ROUTE`;
- `NO_POLICY`;
- `PARTIAL`;
- `CONDITIONAL`;
- `UNKNOWN`;
- `CONFLICTING`; and
- `UNSUPPORTED`.

The result explains why.

## Attack-Path Context

Attack-path analysis uses the same evidence and path model but emphasizes:

- externally reachable services;
- management-plane exposure;
- broad source or destination scope;
- high-risk services;
- trust-boundary crossings;
- segmentation gaps;
- east-west movement;
- identity or authentication dependencies;
- bypass of expected inspection;
- VPN and remote-access reach;
- chained pivots;
- asymmetric visibility;
- stale controls;
- change-created paths; and
- uncertainty.

Atlas presents reviewable attack-path context, not an unexplained verdict.

## Identity-Aware Reachability

When approved BloodHound context is available, Atlas shall distinguish three separate questions:

1. **Identity capability** — what the principal can control or compromise according to the identity graph.
2. **Network capability** — whether the required protocol and service are reachable through the observed or calculated network path.
3. **Combined operational path** — whether the identity and network evidence together support a credible path, including assumptions and unknowns.

A BloodHound path does not prove packet reachability. A reachable service does not prove usable credentials or identity privilege. Atlas reports the combined result only when the evidence supports both sides.

Supported questions should include:

- Which principals can reach network-management services?
- Which identity paths terminate on externally exposed or weakly segmented assets?
- Which critical directory assets are reachable from this VLAN or subnet?
- Which compromised computer provides a network pivot into protected segments?
- Which network control blocks or permits a BloodHound-derived path?
- Which proposed change creates or removes the combined path?

## Dependency and Blast Radius

Given a device, interface, route, VLAN, subnet, policy, tunnel, circuit, service, or change, Atlas should identify:

- direct dependents;
- transitive dependents;
- alternate paths;
- redundancy assumptions;
- single points of failure;
- services and departments affected;
- security controls affected;
- monitoring and logging dependencies;
- expected degradation;
- unknown dependencies; and
- evidence quality.

## Change-Impact Analysis

## Compromise Blast-Radius Analysis

A compromise analysis requires a defined subject, incident time window, scope,
and accepted evidence cutoff.

It reports separately:

- observed activity;
- potential capability;
- prevented activity or paths;
- unknown impact;
- network radius;
- identity and privilege radius;
- data radius;
- infrastructure-control radius;
- operational radius;
- governance and notification radius;
- containment effects;
- conflicts; and
- additional evidence required.

Iron File Intelligence remains authoritative for detailed file and data
evidence. Atlas consumes governed IFI context and correlates it with network,
identity, dependency, and operational evidence.

Potential capability is never presented as observed activity. Missing evidence
is never presented as proof that no activity occurred.

See
[Compromise Blast-Radius and Incident-Impact Intelligence](COMPROMISE-BLAST-RADIUS-AND-INCIDENT-IMPACT-INTELLIGENCE.md)
and [Blast-Radius Result Contract](BLAST-RADIUS-RESULT-CONTRACT.md).

## Change-Impact Analysis

Atlas compares at least:

- current observed or configured state;
- prior accepted state;
- proposed state;
- expected post-change state; and
- actual post-change state.

The analysis identifies:

- records added, removed, enabled, disabled, or modified;
- paths created, removed, or altered;
- routes selected differently;
- policy or ACL behavior changed;
- NAT or VIP behavior changed;
- VPN or SD-WAN behavior changed;
- trust boundaries widened or narrowed;
- dependencies affected;
- network and identity attack paths created or removed;
- monitoring and logging effects;
- expected and unexpected differences;
- validation requirements;
- rollback triggers; and
- residual risk.

## Risk of Approval and Risk of Denial

For every material proposed change, Atlas distinguishes:

### Risk of approval

- outage;
- misrouting;
- asymmetric traffic;
- policy bypass;
- excessive access;
- loss of inspection;
- performance or capacity effect;
- dependency failure;
- rollback complexity; and
- implementation uncertainty.

### Risk of denial or delay

- continuing outage or degradation;
- continuing exposure;
- unsupported state;
- operational workaround;
- staff burden;
- failed compliance objective;
- accumulating technical debt;
- reduced resilience;
- inability to deliver a required service; and
- future emergency-change pressure.

## Evidence and Confidence

Every path step and change-impact conclusion retains:

- evidence artifact;
- source system or device;
- command or source location;
- time;
- digest;
- parser and analyzer version;
- configured, observed, calculated, inferred, unknown, or conflicting state;
- confidence;
- assumptions; and
- human verification.

## Initial Scope

The first complete slice accepts an IP address, CIDR, or VLAN and correlates the Cisco and FortiGate evidence currently supported.

The slice may report partial results, but it shall clearly identify missing Cisco commands, unsupported FortiGate sections, absent live state, stale evidence, and unresolved relationships.

Completeness shall grow through representative evidence and explicit fixtures rather than silent assumptions.
