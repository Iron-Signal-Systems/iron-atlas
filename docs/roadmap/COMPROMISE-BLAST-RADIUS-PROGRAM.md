# Compromise Blast-Radius Program

## Status

Planned cross-phase product program.

This program does not replace the current Cisco, FortiGate, authentication, or
evidence-foundation workstreams.

## 1. Objective

Deliver a defensible answer to:

> **For a defined compromised subject and time window, what did it affect, what
> could it affect, what was prevented, what remains unknown, and who must act?**

## 2. Dependency direction

```text
Cisco and FortiGate evidence
    -> network identity, topology, routing, policy, and reachability

BloodHound and governed identity context
    -> identity privilege and control paths

Iron File Intelligence
    -> file identity, access, classification, activity, and audit coverage

Monitoring, logging, EDR, and security platforms
    -> health, event, detection, process, session, and packet context

Atlas canonical model
    -> correlation, time reconstruction, dependency, and blast-radius explanation
```

## 3. Non-disruption rule

Cisco and FortiGate remain the first major infrastructure evidence sources.

The blast-radius program begins with contracts and sanitized fixtures while
those workstreams continue. It must not claim an end-to-end result before the
required infrastructure model and evidence sources exist.

## 4. Program gates

### BR-0 — Product and authority contract

- define true blast radius;
- define impact categories;
- define external authority;
- accept IFI separation;
- define nonclaims;
- synchronize mission, architecture, requirements, roadmap, and testing.

### BR-1 — Incident, subject, and time model

- incident identity;
- compromised-subject identity;
- time-window uncertainty;
- pre-incident and incident state;
- containment timeline;
- accepted evidence cutoff;
- historical and superseded results.

### BR-2 — External context-bundle intake

- strict bundle contract;
- signature and revocation;
- sequence and replay;
- bounded parsing;
- raw bundle preservation;
- normalized candidates;
- source independence;
- sanitized fixtures.

### BR-3 — Network blast radius

- subject placement;
- Layer-2 and Layer-3 paths;
- route and policy evaluation;
- management-plane exposure;
- observed and potential communications;
- prevented paths;
- stale and unknown state.

This gate depends on the applicable Cisco and FortiGate model.

### BR-4 — Identity and combined attack paths

- governed principal correlation;
- BloodHound context;
- identity-only capability;
- network-only capability;
- combined traversable path;
- credential, certificate, and session references;
- conflict and uncertainty.

### BR-5 — IFI data radius

- IFI endpoint and principal correlation;
- classification snapshots;
- observed file activity;
- potential data access;
- audit coverage;
- copy/disclosure state;
- data owner and review roles.

This gate requires a stable IFI context contract. It does not require Atlas to
access IFI raw files or database tables.

### BR-6 — Infrastructure control and operational dependency radius

- administrative control relationships;
- management systems;
- deployment and automation dependencies;
- business services;
- departments and sites;
- redundancy and alternate paths;
- single points of failure;
- expected operational effect.

### BR-7 — Governance, answer workspace, and reports

- required owners and reviewers;
- responsibility mapping;
- answer-first incident workspace;
- evidence pivots;
- canonical JSON result;
- leadership summary;
- engineering and security detail;
- accessible presentation.

### BR-8 — Representative acceptance

- manually verified multi-system scenarios;
- positive and negative cases;
- stale, missing, conflicting, and delayed evidence;
- source outages;
- replay and revocation;
- false-positive and false-negative accounting;
- correctness and resource summaries;
- clean-clone validation;
- exact accepted commit and evidence.

## 5. First bounded scenario

The first scenario uses sanitized fixtures for:

- one compromised workstation;
- one governed user;
- one switch attachment and VLAN;
- one routed path;
- one firewall permit and one deny;
- one identity-control path;
- one IFI classified-file access summary;
- one operational dependency;
- one missing or stale source; and
- one configured data owner.

The result must show observed, potential, prevented, and unknown impact.

## 6. Workstream placement

Recommended parallel workstreams:

- Cisco evidence;
- FortiGate evidence;
- trusted authentication;
- canonical query and reachability;
- blast-radius contracts and fixtures;
- IFI integration contract; and
- ISRAS adoption readiness.

Only one active acceptance candidate exists per workstream. Cross-workstream
acceptance occurs through an explicit integration candidate.

## 7. Initial implementation sequence

1. documentation and contracts;
2. schemas and sanitized fixtures;
3. bundle validator;
4. incident and subject records;
5. time-window engine;
6. network-radius calculation;
7. identity correlation;
8. IFI context import;
9. dependency and responsibility mapping;
10. answer workspace;
11. representative acceptance.

## 8. Measures

- time saved answering a real incident question;
- percentage of conclusions with exact evidence;
- manually verified path accuracy;
- manually verified affected-asset accuracy;
- false-positive and false-negative results;
- unresolved identity rate;
- stale and missing evidence visibility;
- resource cost;
- usefulness to incident responders and engineers;
- usefulness to data owners and compliance reviewers; and
- limitations reported honestly.
