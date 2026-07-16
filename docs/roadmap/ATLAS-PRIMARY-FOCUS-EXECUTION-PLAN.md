# Atlas Primary-Focus Execution Plan

## Purpose

This document governs the current execution focus for Iron Atlas.

Iron Atlas is the primary active product-development effort. Work must remain concentrated on moving the accepted non-production foundation toward a bounded, read-only, demonstrable infrastructure-intelligence, documentation, reconciliation, and integration-assistance product.

This plan applies only to Iron Atlas. It does not govern or include work for other projects.

## Relationship to the Implementation Roadmap

This document governs execution priority, capacity, representative-environment validation, and focus controls.

The [Implementation Roadmap](IMPLEMENTATION-ROADMAP.md) remains authoritative for phase order and accepted capability boundaries. This plan must not be used to skip a phase gate, weaken an accepted predecessor, or represent roadmap work as complete before formal acceptance.

## Planning Horizon

The current planning horizon is twelve months.

The objective is not to complete the entire long-term Atlas vision. The objective is to establish a minimum credible operational product that can be installed repeatedly, evaluated safely, validated against representative infrastructure, and used to produce defensible infrastructure evidence and reports.

## Current Development Environment

The current authoritative development environment is Arch Linux under Windows Subsystem for Linux with systemd enabled.

The current repository checkout is expected at:

```text
/src/iron-atlas
```

WSL Arch is an accepted development, automated-testing, PostgreSQL-integration, documentation, packaging, and clean-clone-validation environment when the applicable phase gate explicitly supports it.

WSL execution does not by itself prove native-host production behavior, boot-time recovery, hardware-backed isolation, host networking, storage-failure behavior, or representative deployment performance. Production-readiness claims must therefore include later validation on an accepted representative Arch Linux host or virtual machine in addition to WSL development evidence.

## Twelve-Month Product Objective

By the end of the current planning horizon, Iron Atlas should be able to:

- install on the current Arch Linux WSL development host and later on an accepted representative minimal Arch Linux deployment host;
- authenticate users through an accepted production identity boundary;
- receive authorized Cisco, Fortinet, Zabbix, Graylog, and approved documentation evidence without modifying source infrastructure;
- preserve evidence provenance, integrity, parser version, and collection context;
- normalize infrastructure inventory, interfaces, VLANs, trunks, neighbor relationships, port channels, routes, firewall objects, and selected policy records;
- identify a bounded set of high-confidence infrastructure conditions;
- reconcile collected infrastructure identity with selected Zabbix records;
- generate governed Graylog lookup, query, pipeline, stream, dashboard, and report context for syslog and SNMP traps;
- generate useful inventory, topology, preventive-health, discrepancy, integration, and change-review reports;
- operate under explicit resource, timeout, cancellation, retention, and recovery controls; and
- support a controlled representative-environment pilot without claiming production readiness beyond accepted evidence.

## Product Position

Atlas is an active complement to existing operational systems rather than a passive companion or a replacement for mature tools.

- Zabbix remains responsible for monitoring, availability, performance history, triggers, alerting, graphing, escalation, maintenance, and proxy behavior. Atlas supplies authoritative infrastructure identity, evidence age, topology, reconciliation, findings, and reviewable maps, dashboards, template, discovery, and report assistance.
- Graylog and other log platforms remain responsible for centralized syslog and SNMP-trap collection, indexing, retention, search, and investigation. Atlas supplies governed lookup, enrichment, site, device, interface, VLAN, wireless, and topology context plus reviewable queries, pipelines, streams, dashboards, and report definitions.
- Security Onion and other network-security platforms remain responsible for network-security monitoring, packet analysis, detection, and investigation. Atlas may supply governed asset and topology context.
- Cisco and Fortinet systems remain responsible for infrastructure operation and enforcement. Atlas observes, documents, compares, explains, and validates their state.
- Draw.io remains a supported human-editable diagram source and publication format. Atlas may generate separate Draw.io-compatible topology sources.
- Atlas is responsible for authoritative infrastructure evidence, normalized records, topology context, documentation, governed findings, change comparison, integration assistance, and formal acceptance history.

Atlas may consume approved evidence from existing systems and may generate recommendations, exports, maps, dashboards, queries, lookup data, templates, reports, and other reviewable artifacts that improve them.

The default mode is read, reconcile, generate, recommend, export, or deliver. Any future write or provisioning integration requires a separately accepted, least-privileged, previewable, attributable, bounded, approval-aware, and validated boundary.

The first product boundary is an infrastructure-assessment, documentation, and integration-assistance system, not a network controller or automated-remediation platform.

## Capacity Assumption

Planning assumes approximately forty to forty-four focused Atlas hours per week when operational responsibilities permit.

Up to twenty hours per week may occur in an explicitly authorized operational environment. That time must be limited to Atlas work that provides legitimate operational value and complies with employer authorization, security, confidentiality, records, acceptable-use, and intellectual-property requirements.

The plan must retain unscheduled capacity for operational incidents, research, validation failures, documentation synchronization, and recovery. Weekly targets are planning controls, not permission to bypass safety or acceptance requirements.

## Work Boundaries

### Authorized operational-environment work

Appropriate Atlas work in an authorized operational environment includes:

- documenting infrastructure questions and recurring operational pain points;
- preparing approved and sanitized parser fixtures;
- comparing Atlas output with known-good Cisco, Fortinet, Zabbix, Graylog, and documentation state;
- manually validating inventory, topology, and findings;
- recording unsupported syntax and partial-evidence conditions;
- evaluating reports against real administrator workflows;
- measuring collection cost and operational usefulness;
- developing deployment, removal, recovery, and operating procedures; and
- conducting approved read-only lab or shadow evaluation.

Operational-environment work must not include unauthorized access, collection, storage, disclosure, or publication of protected infrastructure information.

### Independent product work

Independent product work includes:

- reusable Go implementation;
- PostgreSQL migrations and service boundaries;
- generic parsers and normalizers;
- test fixtures that contain no protected operational data;
- positive, negative, adversarial, concurrency, recovery, and resource tests;
- installer and systemd work;
- the embedded HTML5 interface;
- product documentation;
- validation tooling and phase gates;
- clean-clone validation; and
- reusable reports and export formats.

### Prohibited repository content

The public repository must not contain:

- raw operational configurations;
- credentials or secrets;
- private keys or certificates;
- unredacted technical-support output;
- protected addresses, VPN details, or security policy;
- employer-specific findings or reports;
- screenshots containing protected infrastructure information; or
- evidence whose publication has not been explicitly authorized.

## Current Execution Sequence

The accepted predecessor is Phase 1 Step 2. Later work must preserve that accepted boundary and proceed in bounded steps.

The Phase 1 Step 3 contract, authentication foundation, governed actor resolver, and bounded OIDC discovery, JWKS, and ID-token verification checkpoints are integrated. Authorization-code exchange with PKCE S256 and bounded one-time in-memory preauthentication transactions is the active implementation candidate. HTTP login/callback routes, cookies, durable sessions, CSRF, logout, trusted-proxy enforcement, governed actor wiring, and production authentication remain unaccepted.

### Stage 1 — Complete the Phase 1 production foundation

Target outcomes:

1. Production authentication and external-identity integration boundary.
2. Production credential delivery and rotation boundary.
3. Verified PostgreSQL TLS and certificate-deployment boundary.
4. Database backup and restoration test boundary.
5. Production connection, queue, worker, storage, timeout, and resource budgets.

No live infrastructure collection should be accepted before the applicable identity, credential, TLS, recovery, and resource controls are established.

### Stage 2 — Evidence intake and protected storage

Target outcomes:

1. Canonical versioned evidence-bundle format.
2. Manual and authenticated evidence intake.
3. Durable staging and explicit receipt state.
4. Signed evidence bundles.
5. Content-addressed immutable evidence storage.
6. Duplicate detection and quarantine.
7. Redaction and classification status.
8. Parser isolation, cancellation, and resource governance.
9. Complete provenance from receipt through normalized output.

Manual sanitized evidence intake should be accepted before live recurring collection.

### Stage 3 — Cisco offline evidence and normalization foundation

Offline sanitized command-bundle support precedes restricted live collection.

The first infrastructure-value slice prioritizes:

1. Catalyst 9300L and Catalyst 9300 access switching.
2. Catalyst 9500 core and distribution switching.
3. Catalyst 9800 Wireless LAN Controllers.

These systems provide broad organizational visibility because users, devices, phones, servers, access points, VLANs, and network paths depend on them.

Initial evidence profiles should cover the commands needed for:

- platform, software, hardware, license, and stack inventory;
- interface state, descriptions, switchport mode, errors, and selected counters;
- VLAN and voice-VLAN inventory;
- trunk, native VLAN, allowed, active, and pruned VLAN state;
- CDP and LLDP neighbors;
- port-channel state and membership;
- spanning-tree state;
- IP-interface and management identity;
- selected environmental and resource state;
- selected running-configuration semantics;
- Catalyst 9800 controller, HA, AP, WLAN, profile, and tag relationships; and
- explicit unsupported, malformed, truncated, conflicting, and partial evidence.

The accepted offline path must prove evidence receipt, digest, provenance, classification, parser version, deterministic normalization, resource governance, human verification, and report usefulness.

### Stage 4 — Cisco semantic, topology, and complementary-integration slice

Initial analysis and integration priority:

1. Device, hardware, software, license, and stack inventory.
2. Interface state, descriptions, errors, and counter deltas.
3. CDP and LLDP relationships.
4. VLAN existence and use.
5. Trunk allowed, active, and pruned VLAN state.
6. Port-channel membership and consistency.
7. Native-VLAN consistency.
8. Spanning-tree root, role, state, and topology-change observations.
9. Catalyst 9800 controller, HA, AP, WLAN, profile, and tag relationships.
10. Zabbix host, interface, template, discovery, map, dashboard, and documentation reconciliation.
11. Graylog lookup and enrichment context for syslog and SNMP traps.
12. Reviewable Graylog query, pipeline, stream, dashboard, and report suggestions.
13. Draw.io-compatible topology generation.
14. Documentation discrepancies.
15. Finding correlation and duplicate suppression.

Deep ACL, QoS, complete wireless-client attribution, and broad endpoint attribution follow a trustworthy inventory and topology foundation.

Restricted live collection, when separately authorized and accepted, must use:

- pinned SSH host keys;
- restricted service authentication;
- fixed command profiles;
- no configuration mode;
- bounded concurrency;
- per-command and per-device timeouts;
- cancellation;
- schedule jitter;
- protected transcripts; and
- complete evidence provenance.

### Stage 5 — Fortinet policy and traffic-boundary slice

The first Fortinet release should prioritize:

- device metadata;
- interfaces and zones;
- addresses and address groups;
- services and service groups;
- firewall policies and evaluation order;
- routes and gateways;
- enabled and disabled state;
- referenced and unreferenced objects;
- parser warnings and unsupported syntax; and
- high-confidence configuration-consistency findings.

Fortinet adds policy and traffic-boundary context to the Cisco-derived inventory and topology foundation.

VPN, SD-WAN, advanced NAT, and complete traffic-path explanation remain later increments unless required by an accepted step.

### Stage 6 — First operational report

The first operationally useful Atlas report should include:

- executive summary;
- collection coverage and age;
- unsupported and incomplete evidence;
- device, hardware, and software inventory;
- site and topology views;
- interfaces, trunks, VLANs, and port channels;
- neighbor relationships;
- spanning-tree observations;
- Fortinet object and policy observations when the applicable slice is available;
- Zabbix reconciliation plus map, dashboard, template, discovery, and report recommendations;
- Graylog lookup, enrichment, query, dashboard, and report context;
- findings by confidence and severity;
- documentation discrepancies;
- evidence source, timestamp, digest, parser version, and collection context;
- changes since the previous accepted collection; and
- human disposition state.

Every finding must identify what Atlas observed, why it may matter, what evidence supports it, the confidence level, and what a human should verify.

### Stage 7 — Controlled representative-environment pilot

The first pilot must be explicitly authorized and bounded.

A suitable initial scope is:

- one noncritical site or lab boundary;
- one Catalyst 9300L/9300 access switch or stack;
- one Catalyst 9500 core or distribution device or logical pair where authorized;
- one Catalyst 9800 controller or approved offline controller bundle where available;
- Zabbix reconciliation for the same equipment;
- Graylog context generation for approved syslog or SNMP-trap examples;
- one Fortinet configuration export only after the Fortinet slice is ready;
- offline or manual intake before live collection;
- no automated changes;
- no automated remediation;
- explicit collection windows and stop conditions;
- documented installation and removal; and
- no publication of protected evidence or findings.

Pilot success requires:

- clean installation and removal on the accepted pilot host;
- separate identification of WSL development results and representative-host results;
- no observed infrastructure modification;
- no unacceptable resource impact;
- complete evidence provenance;
- manually verified inventory and topology results;
- recorded false positives, false negatives, unsupported cases, and uncertainty;
- useful and explainable findings;
- successful backup and recovery procedures;
- measurable administrator time saved or decision quality improved; and
- an explicit statement of what the pilot does and does not prove.

## Quarterly Milestones

### Quarter 1 — Finish the production foundation

- Complete the remaining Phase 1 boundaries.
- Preserve Phase 1 Step 2 acceptance unchanged.
- Establish accepted production identity, credential, TLS, recovery, and resource-control contracts.

### Quarter 2 — Establish evidence intake and Cisco first value

- Accept the evidence-bundle and protected-storage boundaries.
- Accept parser isolation and provenance.
- Deliver the first bounded offline Catalyst 9300L/9300, 9500, and 9800 evidence, normalization, and topology slice.
- Establish initial Zabbix reconciliation and Graylog context generation.

### Quarter 3 — Establish restricted Cisco collection and Fortinet boundary context

- Add authorized restricted Cisco collection only after offline acceptance.
- Deliver normalized device, interface, VLAN, trunk, neighbor, stack, port-channel, spanning-tree, and wireless records.
- Deliver reviewable Zabbix maps, dashboards, template, discovery, and reporting assistance.
- Deliver reviewable Graylog lookup, query, dashboard, and report assistance.
- Deliver the first bounded Fortinet policy and traffic-boundary slice.

### Quarter 4 — Deliver reporting and a controlled pilot

- Deliver the first operational report.
- Reconcile selected collected records with Zabbix.
- Execute the controlled representative-environment pilot.
- Record findings, limitations, operational cost, and next-product decisions.

## Weekly Operating Rhythm

A normal week should include:

- requirements, contracts, and architecture;
- implementation;
- unit and integration tests;
- negative, adversarial, concurrency, recovery, and resource testing;
- representative-evidence validation;
- documentation synchronization;
- acceptance-candidate review; and
- next-step planning.

Operational validation should not replace implementation testing, and implementation testing should not replace manual comparison with representative evidence.

## Focus Controls

### One active implementation step

At any time, Atlas should have:

- one accepted predecessor;
- one current candidate step;
- one documented scope;
- one acceptance gate;
- one bounded change set;
- one authoritative development branch; and
- synchronized documentation, implementation, tests, validation, evidence, and acceptance records.

Multiple later phases may be researched, but only the current step may drive implementation acceptance.

### Required first-product backlog

The first-product backlog includes:

- production identity;
- credential delivery and rotation;
- database TLS;
- backup and recovery;
- resource budgets;
- evidence intake and storage;
- Cisco 9300L/9300, 9500, and 9800 offline parsing and normalization;
- Cisco restricted collection after offline acceptance;
- inventory and topology;
- Zabbix reconciliation, maps, dashboards, template, discovery, and report assistance;
- Graylog lookup, enrichment, query, pipeline, stream, dashboard, and report assistance;
- Fortinet parsing, policy normalization, and traffic-boundary context;
- defensible findings;
- operational reporting; and
- a controlled read-only pilot.

### Deferred backlog

The following are deferred until the first product boundary is proven:

- automated configuration deployment;
- automated remediation;
- full network-controller behavior;
- full CMDB replacement;
- full ticketing-system behavior;
- broad multi-vendor coverage;
- replacement of Zabbix, Graylog, Security Onion, or vendor management systems;
- unrestricted user-defined device commands;
- AI-generated production changes; and
- any feature that weakens evidence provenance, two-person control, or source-system safety.

## Progress Measures

Progress is measured by accepted capability and operational usefulness, not by document count or feature count.

Primary measures include:

- accepted steps completed without weakening predecessors;
- representative fixtures covered;
- parser determinism and unsupported-case handling;
- manually verified inventory and topology accuracy;
- false-positive and false-negative accounting;
- collection and parser resource observations;
- installation, removal, backup, and restoration success;
- report usefulness to an infrastructure operator;
- administrator time saved;
- operational decisions improved by Atlas evidence; and
- unresolved risks and limitations stated explicitly.

## Review Cadence

This plan should be reviewed at least monthly and at every phase boundary.

A review may change sequencing when evidence justifies it, but changes must be synchronized with the implementation roadmap, requirements, architecture, testing model, phase gates, and acceptance records that are affected.

The plan must not be used to represent unaccepted work as complete or production-ready.
