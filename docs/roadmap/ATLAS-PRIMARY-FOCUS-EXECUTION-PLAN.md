# Atlas Primary-Focus Execution Plan

## Purpose

Govern the current development focus for Iron Atlas without allowing platform plumbing, vendor parsers, or engineering ceremony to obscure the product objective.

Iron Atlas is the primary active product-development effort.

## Product Objective

Iron Atlas shall become an evidence-driven network and security intelligence platform that gives Network Operations Teams, Security Operations Teams, operational leaders, and change authorities fast, defensible answers about the environment.

The objective is not to collect the most records or implement the most vendor syntax.

The objective is to answer useful questions with evidence.

## Product Mindset

> **Do not make the user browse the network. Reconstruct the network and answer the question.**

## Current Development Context

The authoritative local repository is expected at:

```text
/src/iron-atlas
```

Arch Linux under Windows Subsystem for Linux is an accepted development and automated-testing environment where the applicable gate supports it.

WSL results do not by themselves establish native-host production behavior, boot recovery, systemd isolation, storage failure behavior, host firewall behavior, representative performance, or production compatibility.

## Current Foundation Status

Phase 0 is an accepted non-production executable baseline.

Phase 1 Steps 1 and 2 are accepted non-production PostgreSQL governance and runtime boundaries.

Phase 1 Step 3 trusted-authentication work remains incomplete.

Cisco and FortiGate ingestion work is active product research and implementation. It shall remain clearly distinguished from formally accepted production boundaries.

## First Major Product Workstreams

### Cisco Evidence Workstream

Purpose:

- ingest offline configuration and operational command evidence;
- normalize devices, interfaces, VLANs, trunks, neighbors, port channels, spanning tree, routes, ACLs, wireless relationships, health, and diagnostics;
- preserve provenance, command support, partial state, and uncertainty; and
- provide Layer-2, Layer-3, attachment, topology, and health context.

Offline evidence precedes restricted live collection.

### FortiGate Evidence Workstream

Purpose:

- ingest native configuration, supported YAML, and operational or diagnostic evidence;
- normalize interfaces, zones, VDOMs, addresses, groups, services, policies, routes, NAT, VIPs, VPNs, SD-WAN, and relevant security controls;
- preserve policy order, runtime uncertainty, unsupported sections, and unresolved relationships; and
- provide routing, policy, translation, VPN, SD-WAN, and trust-boundary context.

Configuration ingestion precedes restricted live diagnostics or collection.

### Identity Attack-Graph Integration Workstream

Purpose:

- import bounded, approved BloodHound context and SharpHound-derived evidence;
- correlate directory principals and computers with Atlas assets using governed identity matching;
- distinguish identity privilege from packet reachability;
- generate a versioned Atlas OpenGraph extension and reviewable payloads;
- answer combined identity, network, exposure, and change-impact questions; and
- keep BloodHound responsible for its identity attack graph rather than rebuilding it inside Atlas.

Offline exports and sanitized fixtures precede API automation or collector orchestration.

### Product-Vision and Query Workstream

Purpose:

- maintain the answer-first product boundary;
- define the canonical query and answer model;
- define the cross-vendor record model;
- define UI behavior;
- define change-impact and decision support; and
- prevent vendor modules from becoming isolated products.

### Foundation Workstream

Purpose:

- complete required identity, credential, TLS, recovery, resource, evidence-protection, and deployment boundaries without blocking safe offline parser and model development.

## Parallel Workstream Rule

A single developer may maintain multiple bounded workstreams in separate branches or worktrees.

Each workstream shall have:

- one accepted predecessor;
- one declared candidate;
- one scope;
- one validation boundary;
- one evidence set; and
- one next step.

Only one acceptance candidate may be active within a workstream.

Cross-workstream behavior is accepted through an explicit integration candidate.

## First Cross-Vendor Intelligence Slice

The first complete Atlas intelligence slice accepts an IP address, CIDR, or VLAN and answers, where evidence permits:

- Where is it?
- Which prefix contains it?
- Which VLAN, switch, port, trunk, SVI, gateway, VDOM, VRF, routing domain, and zone participate?
- Which route is selected and why?
- Which firewall policies, ACLs, NAT rules, VIPs, VPNs, and SD-WAN relationships affect it?
- What can it reach?
- What can reach it?
- Which trust boundaries and attack paths exist?
- What depends on it?
- What evidence supports the answer?
- What remains unknown, stale, incomplete, conflicting, or unsupported?
- How would a proposed change alter the answer?

Partial results are acceptable when limitations are explicit.

## Execution Sequence

### Stage 1 — Preserve and Complete Required Foundation Boundaries

Continue the accepted Phase 1 sequence for:

- trusted production authentication;
- governed actor resolution;
- credential delivery and rotation;
- PostgreSQL TLS;
- backup and recovery; and
- production resource budgets.

Offline parser, model, fixture, and documentation work may continue in isolated workstreams because it does not exercise live collection or production authority.

### Stage 2 — Establish Offline Evidence Contracts

Deliver:

- versioned Cisco command bundles;
- versioned FortiGate configuration and diagnostic bundles;
- versioned BloodHound context bundles and approved SharpHound-derived evidence records;
- provenance and digest;
- classification;
- parser version;
- protected-input boundaries;
- deterministic normalization;
- resource and cancellation controls;
- malformed, truncated, oversized, and conflicting input handling; and
- sanitized fixtures.

### Stage 3 — Build the Minimum Canonical Cross-Vendor Model

Prioritize records required by the first intelligence slice:

- devices;
- interfaces;
- VLANs;
- subnets and CIDRs;
- zones, VDOMs, VRFs, and routing domains;
- routes and next hops;
- policies and ACLs;
- address and service objects;
- NAT and VIPs;
- VPN and SD-WAN;
- topology and attachment;
- evidence and uncertainty; and
- accepted and proposed state;
- external identity principals, privilege-path references, and governed asset correlations.

### Stage 4 — Deliver IP, CIDR, and VLAN Intelligence

Provide query results for:

- containment;
- longest-prefix match;
- overlap and conflict;
- VLAN and interface placement;
- Layer-2 and Layer-3 relationships;
- routes;
- policy and ACL use;
- NAT and VPN use;
- dependencies;
- findings; and
- evidence.

### Stage 5 — Deliver Reachability and Attack-Path Explanation

Provide source-to-destination analysis with:

- protocol and port;
- route selection;
- policy and ACL evaluation;
- NAT;
- VPN and SD-WAN;
- return-path uncertainty;
- trust-boundary crossings;
- identity privilege, network reachability, and combined attack-path context; and
- step-by-step evidence.

### Stage 6 — Deliver Change-Impact and Decision Packages

Provide:

- current-to-prior comparison;
- current-to-proposed comparison;
- paths created, removed, or altered;
- dependency and blast-radius analysis;
- risk of approval;
- risk of denial or delay;
- director-facing decision summary;
- engineering implementation and validation plan;
- rollback;
- post-change evidence; and
- acceptance record.

### Stage 7 — Controlled Representative Pilot

Use an explicitly authorized, read-only, bounded environment.

Pilot success requires:

- manually verified answers;
- false-positive and false-negative accounting;
- visible unsupported and unknown state;
- no infrastructure modification;
- complete evidence lineage;
- acceptable resource impact;
- measurable administrator time saved;
- improved decision quality;
- useful director and engineering change packages; and
- explicit limits on what the pilot proves.

## Product Measures

Primary measures include:

- time required to answer a real engineer question before and after Atlas;
- percentage of answer steps supported by evidence;
- manually verified identity, prefix, VLAN, route, policy, and path accuracy;
- false-positive and false-negative accounting;
- unresolved and unsupported evidence;
- parser and collector cost;
- usefulness of findings;
- usefulness of change-impact explanations;
- decision quality improved;
- rework avoided;
- accepted capability completed without weakening predecessor boundaries; and
- limitations stated honestly.

## Engineering-Practice Boundary

Atlas follows the [Solo-Developer Operating Model](../engineering/SOLO-DEVELOPER-OPERATING-MODEL.md).

The `engineering-standards` repository provides adoptable practices and tooling. It does not silently govern Atlas or block bounded product work.

## Deferred Until Justified

- automated network changes;
- automated remediation;
- full controller behavior;
- unrestricted device commands;
- broad multi-vendor coverage before Cisco and FortiGate slices prove value;
- replacement of Zabbix, Graylog, Security Onion, BloodHound, vendor systems, or a full CMDB;
- AI-generated production changes;
- unexplained risk scoring; and
- features that weaken evidence provenance, authorization, approval independence, source-system safety, or uncertainty.
