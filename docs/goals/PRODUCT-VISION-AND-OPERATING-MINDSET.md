# Product Vision and Operating Mindset

## Canonical Product Statement

**Iron Atlas is an evidence-driven network and security intelligence platform built to give Network Operations Teams, Security Operations Teams, operational leaders, and change authorities fast, defensible answers about the environment.**

## The Problem Atlas Solves

Infrastructure truth is usually fragmented across:

- firewall and switch interfaces;
- running configurations;
- technical-support bundles;
- routing and forwarding tables;
- VPN and SD-WAN diagnostics;
- monitoring systems;
- logging systems;
- packet and security platforms;
- identity and privilege attack graphs;
- diagrams;
- spreadsheets;
- tickets;
- change records; and
- engineer knowledge.

An experienced engineer can often reconstruct the answer, but doing so is slow, repetitive, error-prone, difficult to review, and hard to preserve.

Atlas exists to perform that correlation once, preserve the evidence, and present the answer in a form that can be understood and challenged.

## Product Mindset

### Answer-first

Atlas begins with the user’s question, not the vendor’s navigation tree.

A query may begin with:

- an IP address;
- a CIDR or subnet;
- a VLAN;
- a device;
- an interface;
- a route;
- a policy;
- a service or port;
- a source and destination;
- an observed problem;
- a proposed change; or
- a suspected attack path.

### Evidence-driven

Every material conclusion identifies:

- the evidence source;
- collection or import time;
- artifact digest;
- parser or analyzer version;
- source location or command;
- evidence state;
- confidence;
- conflicts;
- limitations; and
- unknowns.

### Cross-vendor

Vendor syntax and identifiers remain available as evidence and attributes, but vendor-specific representations do not become the canonical model.

### Operational and security views of the same environment

Network availability and security exposure are not separate realities.

Atlas correlates:

- switching;
- routing;
- wireless;
- firewall policy;
- ACLs;
- NAT;
- VPN;
- SD-WAN;
- management-plane controls;
- monitoring;
- logging;
- security telemetry;
- identity privilege and attack-graph context;
- dependencies; and
- change history.

### Decision support, not unexplained scoring

Atlas explains why a condition matters. It does not replace accountable engineering or management judgment with an unexplained risk score.

### Visible uncertainty

Configured intent is not automatically current operational truth.

Atlas distinguishes:

- `CONFIGURED`;
- `OBSERVED`;
- `CALCULATED`;
- `INFERRED`;
- `UNKNOWN`; and
- `CONFLICTING`.

### Human authority

Atlas initially recommends, explains, compares, and validates. Humans authorize and execute infrastructure changes.

## The Atlas Product Test

A feature has product value when it helps a qualified user obtain a defensible answer faster, with less manual correlation and better evidence.

The primary product test is:

> **Does this capability save a senior engineer from manually correlating multiple devices, interfaces, VLANs, routes, objects, policies, command outputs, diagrams, logs, and monitoring screens to reach the same conclusion?**

A parser may be necessary without being independently valuable. A report may be accurate without being useful. A dashboard may be attractive without answering a question.

Progress is measured by trustworthy answers and improved decisions, not by parser count, screen count, document count, or raw record volume.

## Required Answer Shape

A complete Atlas answer should provide, when applicable:

1. **Identity** — what the subject is.
2. **Placement** — where it exists physically and logically.
3. **Address and prefix** — containing networks, longest-prefix match, ranges, overlap, and conflict.
4. **Layer 2** — VLANs, access ports, trunks, port channels, neighbors, and spanning tree.
5. **Layer 3** — interfaces, gateways, routing domains, routes, and next hops.
6. **Controls** — firewall policy, ACL, NAT, VPN, SD-WAN, schedules, identities, and inspection.
7. **Reachability** — what can communicate, why, and under which assumptions.
8. **Exposure and attack path** — trust-boundary crossings, identity privilege paths, network pivots, and the evidence connecting them.
9. **Dependencies and blast radius** — what relies on the subject.
10. **Change** — what differs from prior, proposed, or accepted state.
11. **Evidence** — exact support for each conclusion.
12. **Confidence and unknowns** — what Atlas cannot prove.

## Change-Decision Vision

Atlas shall support both engineering and leadership decisions.

A director-facing change view explains:

- the problem;
- the proposed outcome;
- operational and security benefits;
- risk of approval;
- risk of denial or delay;
- affected services and departments;
- outage and rollback expectations;
- recommendation;
- confidence; and
- the decision requested.

The engineering and security view explains:

- exact devices and configuration involved;
- current and proposed state;
- topology and traffic-path effects;
- routes, VLANs, policies, ACLs, NAT, VPNs, and dependencies affected;
- network and identity attack paths created or removed;
- implementation sequence;
- validation;
- rollback;
- post-change monitoring; and
- source evidence.

Both views are generated from the same governed evidence and analysis.

## First Cross-Vendor Intelligence Slice

The first complete cross-vendor slice should accept an IP address, CIDR, or VLAN and answer:

- Where is it?
- Which subnet or prefix contains it?
- Which VLAN, switch, port, trunk, SVI, gateway, VDOM, VRF, routing domain, and zone participate?
- Which route is selected?
- Which firewall policies, ACLs, NAT, VIPs, VPNs, and SD-WAN relationships affect it?
- What can it reach?
- What can reach it?
- Which trust boundaries can it cross?
- Which approved BloodHound identity and privilege paths apply to the asset, and does the network permit the required service?
- What depends on it?
- What evidence proves the answer?
- What remains unknown?
- How would a proposed change alter the answer?

Partial support is acceptable when limitations are explicit. Silent guessing is not.

## Anti-Patterns

Atlas shall avoid:

- reproducing a vendor’s 55-click navigation model;
- presenting raw normalized data as though it were an answer;
- hiding unsupported evidence;
- treating configuration as proof of live state;
- treating monitoring metadata as canonical authority;
- presenting inferred topology as observed fact;
- silently rewriting accepted history;
- automatically declaring every unusual configuration incorrect;
- requiring a large JavaScript framework without demonstrated operational value;
- treating BloodHound paths as proof of packet reachability or packet reachability as proof of identity compromise;
- importing organizational process that blocks proportional single-developer progress; and
- mistaking engineering ceremony for engineering assurance.
