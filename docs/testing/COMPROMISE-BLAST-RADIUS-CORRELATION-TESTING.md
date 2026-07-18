# Compromise Blast-Radius Correlation Testing

## Status

Planning test model.

## 1. Objective

Prove that Atlas can correlate external evidence into deterministic,
time-bounded, evidence-backed incident-impact results without inventing
authority, hiding uncertainty, or coupling to external databases.

## 2. Test layers

### Contract tests

- bundle schema;
- result schema;
- canonicalization;
- digest;
- signature;
- revocation evidence;
- sequence;
- predecessor;
- replay;
- duplicate and conflict behavior.

### Entity-correlation tests

- stable exact match;
- alias match;
- SID and domain match;
- certificate identity;
- endpoint attachment;
- content identity;
- hostname collision;
- mutable IP reuse;
- merge;
- split;
- unresolved;
- conflicting human disposition.

### Time tests

- pre-incident state;
- incident start and end;
- clock uncertainty;
- source delay;
- configuration change during incident;
- identity change during incident;
- containment becoming effective;
- classification changed after access;
- late-arriving evidence;
- superseding result.

### Network tests

- permitted path;
- denied path;
- no route;
- return-path failure;
- NAT;
- VPN;
- SD-WAN;
- management plane;
- stale firewall state;
- conflicting topology;
- unsupported vendor behavior.

### Identity tests

- identity path without network reachability;
- network reachability without identity authority;
- combined traversable path;
- disabled account;
- expired role binding;
- certificate or session exposure;
- conflicting BloodHound and directory context.

### IFI tests

- observed classified-file access;
- potential access without observed read;
- deletion and absence confirmation;
- audit coverage gap;
- stale IFI bundle;
- invalid IFI signature;
- revoked IFI signer;
- classification at event time;
- minimized bundle content;
- authorized evidence pivot.

### Dependency and governance tests

- direct dependency;
- transitive dependency;
- alternate path;
- single point of failure;
- unknown owner;
- multiple data owners;
- required review mapping;
- no automatic legal notification conclusion.

## 3. Hostile tests

- oversized bundle;
- excessive records;
- decompression expansion;
- deeply nested input;
- malformed UTF-8;
- duplicate JSON keys;
- unsupported contract;
- signature substitution;
- certificate purpose confusion;
- stale revocation;
- sequence reset without authority;
- cross-organization transplant;
- record replay under another source;
- path explosion;
- graph cycle;
- resource exhaustion;
- cancellation;
- parser crash;
- partial PostgreSQL failure.

## 4. Determinism

The same accepted evidence set, analyzer release, policy release, subject, and
time window must produce the same semantic result digest.

Order of independent bundle arrival must not change the final accepted semantic
result after all bundles are present.

## 5. Reference fixture campaign

The first fixture campaign contains:

- workstation `WS-044`;
- a governed user;
- switch, VLAN, subnet, and firewall evidence;
- one observed WinRM connection;
- one denied RDP path;
- one identity administration path;
- one IFI sensitive-file read;
- one potential IFI data-access set;
- one operational dependency;
- one audit gap; and
- one configured privacy reviewer.

Expected output identifies observed, potential, prevented, and unknown impact.

## 6. Quality accounting

Retain:

- expected affected entities;
- reported affected entities;
- false positives;
- false negatives;
- expected unknowns;
- incorrectly asserted certainty;
- unresolved identities;
- conflict handling;
- explanation quality;
- evidence completeness; and
- human disposition.

## 7. Resource observations

Keep correctness and resources separate.

Record:

- total duration;
- phase duration;
- CPU;
- memory;
- PostgreSQL activity;
- graph nodes and edges;
- path candidates;
- bundle bytes and records;
- cache use;
- result size;
- cancellation time; and
- host and toolchain fingerprint.

Thresholds begin as observation-only until representative baselines exist.

## 8. Acceptance evidence

An accepted campaign binds:

- exact Atlas commit;
- exact external fixture digests;
- analyzer and policy release;
- database migration identity;
- toolchain;
- test source;
- result digest;
- correctness summary;
- resource summary;
- limitations; and
- sanitized retained evidence.
