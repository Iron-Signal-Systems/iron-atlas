# Compromise Blast-Radius and Incident-Impact Intelligence

## Status

Normative target direction and planning contract.

This document does not claim that cross-system incident correlation, Iron File
Intelligence integration, complete reachability, complete dependency mapping, or
production incident-response functionality is implemented or accepted.

## 1. Objective

Iron Atlas shall explain the most complete evidence-backed blast radius for a
defined subject, incident time window, and accepted evidence state.

Atlas shall not claim omniscience. It shall distinguish:

- what accepted evidence proves occurred;
- what available identity, permission, reachability, trust, and dependency
  evidence shows was possible;
- what an effective control prevented;
- what remains unknown because evidence is absent, stale, incomplete,
  unsupported, conflicting, or outside collection coverage; and
- which assumptions materially affect the conclusion.

## 2. Central question

For a compromised user, endpoint, server, account, certificate, application,
network device, service, or file:

> **What could the compromised subject reach, access, alter, impersonate,
> administer, disrupt, disclose, destroy, or use as a path to another system—and
> what does the evidence prove actually happened?**

## 3. Required scope

Every blast-radius calculation identifies:

- incident identity;
- compromised subject or subjects;
- subject type;
- incident start and end;
- uncertainty in the incident window;
- organization, site, tenant, domain, network, and service scope;
- accepted-state version or evidence cutoff;
- calculation release;
- evidence sources;
- coverage state;
- assumptions;
- conflicts; and
- result digest.

A result without a defined subject, time window, and evidence cutoff is not a
complete blast-radius result.

## 4. Impact classes

### 4.1 Observed impact

Activity directly supported by accepted evidence, such as:

- successful logon;
- process execution;
- network connection;
- file access;
- file modification or deletion;
- permission change;
- remote administration;
- infrastructure configuration change;
- certificate use;
- credential use;
- service disruption; or
- confirmed data movement.

### 4.2 Potential impact

Capability supported by accepted identity, permission, reachability, trust,
dependency, or control evidence but not proven to have been exercised.

Potential impact must never be presented as observed activity.

### 4.3 Prevented impact

An attempted or technically relevant action that an accepted control denied or
prevented.

The result identifies the control, enforcement point, time applicability, and
evidence supporting the prevention claim.

### 4.4 Unknown impact

Impact that cannot be resolved because evidence is:

- missing;
- stale;
- delayed;
- incomplete;
- unsupported;
- conflicting;
- outside audit coverage;
- outside retention;
- unavailable because a source failed; or
- not collectible under the accepted authority boundary.

Missing evidence does not prove that no activity occurred.

## 5. Blast-radius dimensions

Atlas calculates related dimensions rather than one unexplained score.

### 5.1 Network radius

- reachable devices, addresses, services, and segments;
- Layer-2 and Layer-3 paths;
- routing decisions;
- firewall and ACL decisions;
- NAT, VIP, VPN, and SD-WAN effects;
- return paths;
- trust-boundary crossings;
- management-plane exposure;
- observed communications; and
- denied or unavailable paths.

### 5.2 Identity and privilege radius

- users, groups, computers, service identities, certificates, and sessions;
- direct and transitive authority;
- local and remote administration;
- delegated control;
- identity-provider and directory relationships;
- BloodHound-derived context;
- credential or session exposure; and
- identity paths that are also network-traversable.

Identity capability, network capability, and the combined operational path
remain separate conclusions.

### 5.3 Data radius

- file and data-object identities;
- accepted classification at event time;
- effective-access summaries;
- observed reads, writes, renames, moves, permission changes, encryption, and
  deletion;
- potential access;
- copy or disclosure correlation;
- data owner;
- regulated or contractual handling context;
- audit-coverage limitations; and
- authoritative Iron File Intelligence evidence references.

### 5.4 Infrastructure-control radius

- switches;
- routers;
- firewalls;
- wireless controllers;
- hypervisors;
- storage;
- backup systems;
- certificate authorities;
- monitoring and logging systems;
- endpoint-management and deployment systems;
- automation identities; and
- management APIs.

Reachability to a management service does not prove administrative authority.
Administrative authority does not prove that the required service was
reachable.

### 5.5 Operational radius

- applications;
- business services;
- departments;
- sites;
- users;
- public or institutional functions;
- upstream and downstream dependencies;
- redundancy and alternate paths;
- single points of failure;
- expected degradation; and
- recovery dependencies.

### 5.6 Governance and notification radius

Atlas identifies configured review and responsibility relationships, including:

- system owner;
- application owner;
- service owner;
- data owner;
- information-security authority;
- privacy authority;
- legal or compliance authority;
- records authority;
- contract owner;
- regulatory or program-specific authority;
- insurance incident contact; and
- executive incident authority.

Atlas does not make a legal notification decision. It identifies affected
evidence and the governed roles that must review it.

## 6. Evidence state

Every material fact or conclusion retains one or more states:

- `CONFIGURED`;
- `OBSERVED`;
- `CALCULATED`;
- `INFERRED`;
- `UNKNOWN`;
- `CONFLICTING`;
- `STALE`;
- `INCOMPLETE`;
- `UNSUPPORTED`; and
- `PREVENTED`.

The result must not flatten these states into one confidence number.

## 7. Time reconstruction

Atlas reconstructs:

- prior accepted state;
- state immediately before the incident;
- state at incident start;
- state during the incident;
- observed activity;
- configuration and identity changes;
- containment actions;
- the time each containment action became effective;
- state at incident end;
- post-incident validation; and
- current state.

The newest observation never overwrites the historical state that applied
during the incident.

## 8. Relationship model

Every material incident relationship records:

- source entity;
- destination entity;
- capability or relationship;
- direction;
- protocol, port, right, permission, or trust requirement;
- valid time;
- accepted-state version;
- impact class;
- evidence state;
- evidence references;
- confidence;
- correlation method;
- assumptions;
- conflicts;
- coverage limitations; and
- supersession.

Structural relationships and traversable compromise paths remain distinct.

## 9. Correlation quality

Allowed correlation quality includes:

- `CONFIRMED`;
- `PROBABLE`;
- `POSSIBLE`;
- `UNRESOLVED`; and
- `CONFLICTING`.

Atlas records why the correlation received that state.

Short hostname, mutable IP address, display name, or temporal proximity alone
shall not silently merge identities.

## 10. Required answer shape

A complete answer presents:

1. subject identity;
2. incident window;
3. evidence cutoff and coverage;
4. observed activity;
5. potential reach and authority;
6. prevented paths;
7. identity and privilege exposure;
8. data and classification exposure;
9. infrastructure-control exposure;
10. operational dependencies;
11. containment timeline and remaining exposure;
12. responsible owners and required reviewers;
13. exact supporting evidence;
14. confidence, assumptions, conflicts, stale evidence, and unknowns; and
15. additional evidence required.

## 11. Example summary

```text
Incident: Compromise of WS-044
Window: 2026-07-18 08:42–09:31
Coverage: strong server and network evidence; partial client process coverage

Observed
- The governed user authenticated from WS-044.
- A process contacted ADMIN01 over WinRM.
- Accepted IFI evidence reports classified file reads and one deletion on FS01.

Potential
- WS-044 could reach 31 servers.
- The identity could administer 12 workstations.
- Nine targets were both identity-accessible and network-reachable.

Prevented
- RDP to the management network was denied.
- The backup subnet was not reachable under the applicable firewall state.

Unknown
- Client process telemetry was unavailable for seven minutes.
- One legacy share lacked complete file-audit coverage.
- Available evidence does not prove credential extraction or external exfiltration.
```

## 12. Product boundary

Atlas remains an evidence-correlation, calculation, explanation, and
decision-support platform.

Atlas does not become:

- an EDR;
- a SIEM;
- a packet-analysis platform;
- a file-classification engine;
- a duplicate Iron File Intelligence evidence store;
- a legal decision engine;
- an automated containment controller; or
- a generic single pane of glass.

## 13. Acceptance direction

Representative acceptance requires:

- deterministic contract fixtures;
- manually verified incident scenarios;
- positive and negative paths;
- stale, missing, delayed, conflicting, and replayed evidence;
- time-window changes;
- identity-only and network-only paths;
- combined traversable paths;
- IFI data-access context;
- dependency and owner mapping;
- explicit false-positive and false-negative accounting;
- separate correctness and resource summaries;
- complete evidence lineage; and
- exact accepted commit and analyzer identity.
