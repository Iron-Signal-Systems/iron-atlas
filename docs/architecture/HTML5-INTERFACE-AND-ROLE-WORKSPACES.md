# HTML5 Interface and Role Workspaces

## Purpose

Provide an accessible, responsive, answer-first HTML5 interface that helps Network Operations Teams, Security Operations Teams, operational leaders, change authorities, reviewers, and auditors obtain fast, defensible answers without reproducing vendor navigation complexity.

## Governing Principle

> **Do not make the user browse the network. Reconstruct the network and answer the question.**

The interface shall not require a user to traverse dozens of menus to correlate information that Atlas already possesses.

## Primary Interaction

The primary interaction is a global search and query control that accepts:

- IP addresses;
- CIDRs and subnets;
- VLAN IDs and names;
- devices and hostnames;
- interfaces and ports;
- routes and next hops;
- firewall policies and ACLs;
- services, protocols, and ports;
- VPNs, tunnels, and SD-WAN members;
- sites and locations;
- findings and changes; and
- source-to-destination reachability questions.

Examples:

```text
10.20.30.45
10.20.30.0/24
VLAN 240
TenGigabitEthernet1/0/48
tcp/445
can USER-LAN reach DB-CLUSTER on tcp/5432
show overlapping prefixes
show external management exposure
show paths from guest wireless to server networks
what depends on this uplink
what changes if this route is removed
```

## Answer Workspace

A complete answer may contain:

1. Identity
2. Placement
3. Address and prefix
4. Layer 2
5. Layer 3
6. Routing decision
7. Firewall and ACL controls
8. NAT and exposure
9. VPN and SD-WAN
10. Reachability
11. Trust-boundary and attack-path context
12. Dependencies and blast radius
13. Changes from prior, proposed, or accepted state
14. Findings and risk
15. Source evidence
16. Confidence, conflicts, and unknowns

17. Incident time window and evidence cutoff
18. Observed, potential, prevented, and unknown impact
19. Data exposure and IFI evidence references
20. Containment timeline and remaining exposure
21. Responsible owners and required reviewers

The workspace shows the most relevant answer first and permits direct pivots into supporting details.

## Pivot Model

Every meaningful entity is pivotable.

A subnet may pivot to:

- containing or overlapping prefixes;
- VLANs;
- gateways and interfaces;
- devices and sites;
- routes;
- firewall policies;
- ACLs;
- NAT and VIP relationships;
- VPN and SD-WAN use;
- endpoints;
- findings;
- dependencies;
- changes; and
- evidence.

A policy may pivot to:

- source and destination objects;
- service and schedule;
- ingress and egress;
- routes;
- NAT;
- inspection and logging;
- matching traffic paths;
- affected subnets and systems;
- changes; and
- evidence.

A pivot shall preserve context so the user does not repeatedly restart the investigation.

## Role Perspectives

Role perspectives influence emphasis and available governed actions; they do not create separate facts.

### Network Operations

Emphasize:

- device and interface state;
- VLANs, trunks, port channels, spanning tree, neighbors, and wireless;
- subnets, gateways, routing, and dependencies;
- endpoint attachment;
- health and capacity;
- current incidents and findings;
- planned maintenance and change; and
- validation status.

### Security Operations

Emphasize:

- policy and ACL behavior;
- management-plane exposure;
- NAT, VIP, VPN, and internet exposure;
- trust-boundary crossings;
- attack and pivot paths;
- broad or shadowed access;
- segmentation;
- logging and inspection;
- control drift;
- relevant security-platform context; and
- high-risk changes.

### Operational Leadership and Change Authorities

### Incident Response and Compromise Review

Emphasize:

- compromised subject and incident window;
- observed actions;
- potential reach and authority;
- prevented paths;
- identity and session exposure;
- classified data exposure;
- infrastructure-control exposure;
- dependencies and affected services;
- containment status;
- evidence coverage and gaps;
- required system, data, security, privacy, compliance, and executive reviewers;
  and
- exact evidence pivots.

### Operational Leadership and Change Authorities

Emphasize:

- problem and requested decision;
- business and operational effect;
- security and availability effect;
- risk of approval;
- risk of denial or delay;
- affected departments, sites, and services;
- expected outage or degradation;
- rollback readiness;
- evidence quality;
- recommendation and confidence; and
- decision history.

### Reviewers and Auditors

Emphasize:

- evidence provenance;
- parser and analyzer lineage;
- decision and approval history;
- implementation and validation records;
- accepted and superseded state;
- exceptions;
- uncertainty; and
- historical reconstruction.

## State Presentation

The interface shall distinguish:

- configured;
- observed;
- calculated;
- inferred;
- unknown;
- conflicting;
- stale;
- incomplete;
- unsupported;
- proposed;
- approved;
- denied;
- implementing;
- validated;
- rolled back;
- accepted; and
- superseded.

A successful visual state shall not be shown before authoritative confirmation.

## Accessibility and Interaction

- WCAG 2.1 Level AA is the initial target.
- Keyboard operation is a functional requirement.
- Semantic HTML, readable focus state, and screen-reader support are required.
- Color is not the sole carrier of meaning.
- Tables, diagrams, findings, and path explanations provide textual equivalents.
- High-density operational screens remain readable and responsive.
- Common work has a direct path.
- Destructive or governed actions require explicit confirmation and authority.

## Implementation Direction

The initial interface uses Go `html/template`, embedded static assets, semantic HTML5, and minimal JavaScript.

A large single-page application framework is not required for the first implementation and shall not be introduced without demonstrated operational benefit.

The interface is a presentation and workflow boundary. It does not create canonical truth, manufacture authority, or silently convert uncertainty into certainty.
