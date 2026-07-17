# Project Mission

## Mission

**Iron Atlas is an evidence-driven network and security intelligence platform built to give Network Operations Teams, Security Operations Teams, operational leaders, and change authorities fast, defensible answers about the environment.**

Atlas ingests configuration, operational, diagnostic, monitoring, logging, security, and documentation evidence. It correlates that evidence across vendors and systems to explain:

- what infrastructure exists;
- where an IP address, CIDR, subnet, VLAN, interface, device, route, policy, service, tunnel, or endpoint exists;
- how Layer 2, Layer 3, routing, firewall, NAT, VPN, wireless, and security relationships connect;
- what can reach what, through which controls, and with what uncertainty;
- which trust boundaries, exposures, dependencies, and attack paths exist;
- what changed from the previous accepted state;
- what a proposed change will affect;
- what risk exists if the change is approved;
- what risk remains if the change is denied or delayed;
- what evidence supports each conclusion; and
- where evidence is missing, stale, unsupported, incomplete, conflicting, inferred, or unknown.

## Governing Mindset

> **Do not make the user browse the network. Reconstruct the network and answer the question.**

Atlas shall not force an experienced engineer to repeat the manual work of opening dozens of vendor screens, command outputs, monitoring pages, diagrams, spreadsheets, tickets, and log searches merely to reconstruct one operational or security answer.

The interface shall support the work rather than become additional work.

## Primary Users

Atlas is designed for:

- Network Operations Teams;
- Security Operations Teams;
- senior network engineers and administrators;
- security engineers and red-team-oriented purple teamers;
- infrastructure technicians;
- incident responders;
- operational leaders;
- change authorities;
- reviewers and auditors; and
- infrastructure teams responsible for availability, security, and change.

The same evidence may be presented differently to different users, but the underlying facts, lineage, confidence, and uncertainty shall remain consistent.

## Purpose

Atlas maintains a governed relationship among:

- imported configuration;
- running and diagnostic state;
- monitoring and logging context;
- security telemetry and asset context;
- identity and privilege attack-graph context from BloodHound and approved collectors;
- curated documentation and diagrams;
- canonical infrastructure identity;
- topology and dependency;
- address, prefix, VLAN, route, policy, NAT, VPN, SD-WAN, ACL, wireless, and management-plane relationships;
- calculated reachability and attack-path context;
- findings and human disposition;
- current, proposed, accepted, and superseded state;
- change requests, approvals, denial decisions, implementation, rollback, and validation; and
- formal acceptance records.

## First Major Ingestion Points

The first major infrastructure-ingestion workstreams are Cisco and FortiGate.

Cisco evidence provides switching, routing, wireless, endpoint-attachment, topology, and health context.

FortiGate evidence provides policy, security-boundary, routing, NAT, VIP, VPN, SD-WAN, and management-exposure context.

These workstreams may be developed independently, but they converge into one vendor-neutral network and security model. Neither parser, collector, nor vendor module is the Atlas product by itself.

## Complementary-System Principle

Atlas complements mature operational systems.

- Monitoring systems retain monitoring and alerting responsibility.
- Logging systems retain collection, indexing, retention, and search responsibility.
- Security-monitoring platforms retain detection and investigation responsibility.
- BloodHound retains identity and privilege attack-graph responsibility, while Atlas correlates that context with network reachability, exposure, dependencies, and change impact.
- Vendor platforms retain configuration, operation, and enforcement responsibility.
- Diagramming systems retain curated diagram-authoring responsibility.

Atlas retains responsibility for evidence lineage, normalized identity, cross-system correlation, topology, reachability explanation, change impact, governed findings, decision support, validation, and acceptance history.

## Human Authority and Initial Read-Only Boundary

Atlas begins as a read-only evidence, intelligence, comparison, and decision-support platform.

It may produce:

- explanations;
- reports;
- diagrams;
- queries;
- recommendations;
- proposed configuration snippets;
- implementation plans;
- validation plans;
- rollback plans; and
- reviewable integration artifacts.

It shall not silently modify infrastructure or external systems.

Any future write or provisioning capability requires a separately accepted boundary with preview, attribution, authorization, bounded scope, approval awareness, idempotency where practical, reversibility where practical, failure isolation, and post-application validation.

## Leading Principles

> **Answer the question. Show the evidence. Preserve uncertainty. Govern the change. Validate the result. Preserve the record.**
