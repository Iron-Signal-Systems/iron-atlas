# External Evidence Context-Bundle Contract

## Status

Planning contract for external systems, including Iron File Intelligence.

## 1. Purpose

Define a product-independent, versioned, signed, bounded envelope for importing
normalized external context without granting the external system Atlas
authority or granting Atlas direct access to the external system's database.

## 2. Envelope

A bundle identifies:

- contract family and version;
- bundle ID;
- source product and deployment ID;
- organization and scope;
- source sequence;
- predecessor bundle digest;
- source-generated time;
- earliest and latest record-valid time;
- record count;
- uncompressed and transferred byte counts;
- compression profile where used;
- payload digest;
- per-record aggregate root;
- source software release;
- source policy or ruleset identity;
- source accepted-state identity;
- signing key identity;
- signature algorithm and signature;
- encryption recipient identity where applicable; and
- bundle classification and handling requirements.

## 3. Record requirements

Every record identifies:

- record type and version;
- external record ID;
- source entity identities;
- valid time or interval;
- source observation or calculation state;
- canonical payload;
- canonical SHA-256;
- evidence references;
- confidence where applicable;
- coverage limitations;
- supersession or correction identity; and
- handling classification.

## 4. Allowed record families

Initial families may include:

- external asset reference;
- external principal reference;
- certificate reference;
- session reference;
- process reference;
- file-object reference;
- content-identity reference;
- classification snapshot;
- activity operation;
- effective-access summary;
- audit-coverage state;
- security detection reference;
- monitoring-health reference;
- dependency reference;
- governance responsibility; and
- evidence locator.

A source-specific extension must not silently redefine a canonical Atlas field.

## 5. Cryptographic verification

Atlas verifies:

- trusted source identity;
- accepted key purpose;
- certificate path where X.509 is used;
- revocation evidence;
- signature;
- payload digest;
- aggregate root;
- record digests;
- sequence;
- predecessor;
- count;
- lengths; and
- contract version.

A bundle remains provisional until verification succeeds.

## 6. Sequence and replay

Atlas distinguishes:

- first accepted sequence;
- expected next sequence;
- duplicate replay;
- idempotent duplicate;
- conflicting duplicate;
- sequence gap;
- predecessor mismatch;
- signed reset;
- source epoch change; and
- source identity change.

A legitimate reset is a signed record that references the prior known chain
state and explains the reset authority and cause.

## 7. Revocation behavior

Allowed outcomes include:

- `VALID`;
- `REVOKED`;
- `EXPIRED`;
- `NOT_YET_VALID`;
- `UNKNOWN_ISSUER`;
- `WRONG_PURPOSE`;
- `REVOCATION_UNKNOWN`;
- `REVOCATION_STALE`;
- `KEY_COMPROMISE_SUSPECTED`; and
- `CA_COMPROMISE_SUSPECTED`.

Ciphertext may be retained provisionally under policy when revocation evidence
is temporarily stale. It must not become authoritative accepted context until
the required certificate state is established.

## 8. Bounded intake

Each source profile defines:

- maximum bundle bytes;
- maximum record count;
- maximum individual record size;
- maximum nesting;
- allowed compression;
- decompression ratio;
- processing timeout;
- memory budget;
- accepted record types;
- accepted time range;
- backlog maximum; and
- quarantine behavior.

## 9. Data minimization

Bundles contain only context required by the accepted integration.

Sensitive raw evidence remains in the authoritative source unless a separate
retention decision authorizes protected duplication.

## 10. Persistence

Atlas stores:

- immutable raw accepted bundle reference;
- bundle validation outcome;
- source and sequence;
- digests and signature identity;
- parser and adapter release;
- normalization candidates;
- rejected or unsupported records;
- conflicts;
- coverage state; and
- acceptance receipt.

A newer adapter creates new derived records; it does not rewrite the original
bundle.

## 11. Failure isolation

One invalid bundle, record family, source, or adapter does not block unrelated
sources or erase prior accepted state.

## 12. Initial formats

The initial reference format is canonical UTF-8 JSON for the manifest and
NDJSON for records.

YAML is not an initial authoritative bundle format.
