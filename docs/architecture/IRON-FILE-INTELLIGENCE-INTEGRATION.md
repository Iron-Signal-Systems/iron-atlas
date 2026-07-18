# Iron File Intelligence Integration

## Status

Normative target direction and planning contract.

Iron File Intelligence, its production integration, and the Atlas adapter
described here are not yet implemented or accepted.

## 1. Purpose

Define how Iron Atlas consumes governed context from the separate Iron Signal
Systems product **Iron File Intelligence (IFI)**.

IFI supplies authoritative file, classification, access, activity, audit
coverage, and forensic-lineage context. Atlas correlates that context with
network, identity, infrastructure, monitoring, security, dependency, change,
and incident evidence.

## 2. Product separation

Canonical repositories:

```text
Iron-Signal-Systems/iron-file-intelligence
Iron-Signal-Systems/iron-atlas
```

Each product owns its own:

- product mission;
- architecture;
- repository;
- database;
- migrations;
- raw evidence;
- runtime identities;
- authorization;
- release lifecycle;
- validation;
- acceptance history; and
- recovery boundary.

IFI remains useful without Atlas. Atlas remains useful without IFI, although
data-related incident answers may be less complete.

## 3. Authority boundary

IFI remains authoritative for:

- filesystem object identity;
- paths and observations;
- security descriptors and access evidence;
- content identity;
- classifications;
- file activity;
- audit coverage;
- acquisition and extraction lineage;
- IFI cryptographic receipts; and
- detailed forensic evidence.

Atlas remains authoritative for Atlas-owned:

- cross-system entity resolution;
- canonical infrastructure identity;
- network topology;
- reachability;
- combined identity and network path calculation;
- infrastructure dependencies;
- incident correlation;
- blast-radius calculation;
- uncertainty aggregation;
- governance responsibility mapping; and
- Atlas result acceptance.

An imported IFI record does not silently become Atlas-observed truth. Atlas
retains its external origin and accepted IFI state.

## 4. Prohibited coupling

Atlas shall not:

- query IFI PostgreSQL directly;
- depend on IFI internal tables;
- share IFI database roles;
- ingest IFI credentials;
- require IFI source-system credentials;
- ingest raw file content by default;
- ingest complete extracted document text by default;
- copy every IFI Windows event;
- become IFI's primary user interface;
- alter IFI classifications;
- rewrite IFI accepted history; or
- require synchronized product releases for unrelated changes.

IFI shall not require Atlas to collect, classify, sign, seal, anchor, store, or
review IFI evidence.

## 5. Integration modes

### 5.1 Evidence reference

Atlas stores:

- source product;
- external record ID;
- accepted time;
- digest;
- evidence state;
- authorized locator;
- minimized normalized summary; and
- adapter identity.

Detailed evidence remains in IFI.

### 5.2 Signed context bundle

IFI produces a signed, versioned, minimized context bundle containing facts
required for Atlas correlation.

This is the normal integration mode.

### 5.3 Protected raw evidence

Atlas retains protected raw IFI evidence only when a separately accepted
durability, investigation, recovery, or legal-preservation requirement
justifies duplication.

## 6. Context eligible for Atlas

A context bundle may contain:

- IFI source-system identity;
- endpoint identity;
- directory SID and governed principal references;
- filesystem object references;
- content identity references;
- classification snapshots;
- effective-access summaries;
- observed file operations;
- potential-access summaries;
- copy-correlation state;
- audit-coverage state;
- deletion and absence-confirmation state;
- source and event times;
- data-owner references;
- compliance-review references;
- evidence digests;
- authorized evidence locators;
- IFI software, policy, contract, and acceptance identities;
- sequence and predecessor digest; and
- IFI signature.

## 7. Data minimization

The default Atlas bundle excludes:

- source file bytes;
- extracted full text;
- every ACE;
- complete raw security descriptors;
- complete raw Windows event XML;
- endpoint private keys;
- recipient private keys;
- source credentials;
- database credentials; and
- unrestricted personal or regulated content.

Atlas receives only what is necessary to perform governed correlation and
explanation.

## 8. Intake validation

Atlas validates:

- configured IFI source identity;
- contract version;
- bundle signature;
- signer certificate or key identity;
- certificate status and revocation evidence where applicable;
- sequence;
- predecessor digest;
- bundle digest;
- record count and length;
- per-record digests or aggregate root;
- replay, duplicate, gap, reset, and conflict behavior;
- accepted IFI state;
- evidence locator format;
- organization and scope;
- time bounds; and
- resource limits.

Validation failure creates an immutable intake outcome. It does not erase prior
accepted IFI context or block unrelated Atlas operations.

## 9. Identity correlation

IFI and Atlas identity correlation may use:

- governed stable UUIDs;
- source-system IDs;
- machine certificate identity;
- directory SID;
- domain or tenant identity;
- stable filesystem object identity;
- content digest;
- accepted hostname alias;
- MAC or endpoint attachment evidence;
- IP address with valid-time constraints; and
- human disposition.

Mutable names and addresses are evidence, not sole identity.

## 10. Time and classification

Atlas preserves separately:

- classification applicable at event time;
- classification accepted by IFI at bundle generation;
- classification known to Atlas at ingestion;
- later superseding classification; and
- current classification.

Historical activity is never rewritten merely because a policy or
classification changed later.

## 11. Authorization and pivots

Atlas authorization governs whether a user may:

- see minimized IFI context;
- see file paths;
- see classification detail;
- open an authorized IFI evidence pivot;
- request a protected report; or
- include IFI context in an export.

Atlas UI visibility does not grant IFI authority.

## 12. Failure and independence

IFI integration failure:

- does not erase Atlas evidence;
- does not block Cisco or FortiGate ingestion;
- does not block unrelated queries;
- marks data-radius answers partial or unknown;
- records the last accepted bundle and freshness;
- exposes backlog, gap, and revocation state; and
- prevents stale data from being presented as current.

## 13. Initial integration slice

The first bounded slice uses sanitized offline fixtures and answers:

- which Atlas asset corresponds to an IFI endpoint;
- which governed principal corresponds to an IFI actor;
- which classified data categories were observed accessed;
- which additional data was potentially accessible;
- which network path permitted the file-server access;
- which controls denied other paths;
- which evidence coverage gaps remain; and
- which configured owners require review.

No live API or automated IFI collection is required for this first slice.

## 14. Acceptance direction

Acceptance requires:

- a stable IFI context contract;
- signed fixture bundles;
- invalid signature and revocation cases;
- replay, gap, duplicate, and conflicting bundle cases;
- minimized-data validation;
- identity merge and split cases;
- event-time and classification-time cases;
- source outage and stale context;
- authorization and evidence-pivot tests;
- deterministic blast-radius outputs;
- resource bounds; and
- exact IFI and Atlas release identities.
