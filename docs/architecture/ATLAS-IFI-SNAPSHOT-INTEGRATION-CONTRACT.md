# Atlas–Iron File Intelligence Snapshot Integration Contract

## Status

Normative architecture alignment candidate.

## Purpose

Iron File Intelligence (IFI) and Iron Atlas remain separate products with separate authority, data, credentials, lifecycle, and acceptance boundaries.

```text
IFI authoritative state
→ immutable signed purpose-limited evidence snapshot
→ Atlas IFI import adapter
→ validated Atlas correlation records
```

## Authority

IFI is authoritative for file observations, file-intelligence records, evidence collection and lifecycle, IFI classifications and provenance, IFI-side access control and retention, and export-field meaning.

Atlas is authoritative for import validation, Atlas correlation identity, infrastructure reachability and exposure, network/security dependencies, change impact, answer confidence and limitations, and Atlas acceptance history.

## Prohibited coupling

Atlas must not:

- query the IFI PostgreSQL database directly;
- share IFI schemas or service credentials;
- mount or read unrestricted IFI storage;
- assume IFI availability for startup or unrelated queries;
- copy unrestricted file content or unnecessary sensitive metadata;
- write back to IFI;
- issue IFI collection or lifecycle decisions;
- reinterpret IFI classifications without preserving source meaning;
- inherit IFI administrative or forensic authority; or
- treat absence of an IFI snapshot as proof that no relevant file exists.

## Snapshot contract

An IFI snapshot is immutable, purpose limited, minimally disclosed, versioned, schema identified, source/scope identified, time identified, content digested, signed by an accepted IFI export signer, bounded in records and bytes, explicit about omissions and redactions, and bound to export-policy and producer-version identity.

Raw file content, credentials, unrestricted paths, personal data, and forensic material are excluded unless a separate requirement, classification, approval, transport, and storage boundary is accepted.

## Import validation

Atlas verifies signer, signature, digest, schema, producer version, Atlas audience and purpose, source and scope, time and freshness, size and record limits, replay, duplicate state, cardinality, redaction and classification markings, and unsupported fields. A valid signature does not by itself make the snapshot accepted evidence.

## Availability and failure

IFI is optional. Atlas continues unrelated capabilities when IFI is unavailable. Prior accepted IFI-derived evidence ages normally, IFI-dependent answers report missing or stale context, no empty snapshot replaces accepted evidence, and retry follows scheduled-ingestion rules.

## Correlation and write-back

Every Atlas-derived correlation retains the exact IFI snapshot identity. IFI records and Atlas infrastructure identities remain distinct and are joined through explicit governed correlation. No write-back exists in this contract.

## Required validation

Test unknown signer, bad signature, digest mismatch, wrong audience, wrong scope, unsupported schema, oversized snapshot, replay, duplicate records, redaction violations, stale data, outage, partial import, atomic rejection, and cross-source lineage.
