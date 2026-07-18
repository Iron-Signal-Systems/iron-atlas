# ADR-0008 — Iron File Intelligence Remains an External Authoritative System

## Status

Accepted architectural direction; implementation not yet accepted.

## Context

Iron Atlas requires authoritative file identity, classification, effective
access, file activity, audit coverage, and forensic-lineage context to explain
the data dimension of a compromise blast radius.

Iron File Intelligence is a separate Iron Signal Systems product designed to
produce that evidence.

Combining IFI into Atlas, sharing one database, or allowing direct cross-product
table access would create release coupling, authority confusion, oversized
runtime identities, duplicated raw evidence, and pressure for Atlas to become a
generic all-in-one platform.

## Decision

IFI remains a separate product and authoritative system.

Atlas integrates through signed, versioned, minimized context bundles or
governed evidence references.

Atlas does not query IFI PostgreSQL, depend on IFI internal schemas, or become
IFI's primary interface.

IFI does not require Atlas to collect, classify, sign, seal, store, or accept
IFI evidence.

## Consequences

Positive consequences:

- independent product focus;
- independent release and recovery;
- explicit source authority;
- smaller service identities;
- minimized data duplication;
- replaceable integration;
- clearer failure isolation;
- preserved evidence lineage; and
- Atlas remains focused on correlation and explanation.

Costs:

- a formal integration contract is required;
- identity correlation must be governed;
- both products must preserve stable external identifiers;
- version compatibility must be tested;
- cross-product revocation and signing keys must be managed; and
- some users will pivot between products for detailed evidence.

## Rejected alternatives

### One shared database

Rejected because it couples migrations, authorization, releases, recovery, and
internal schemas.

### Direct Atlas access to IFI source systems

Rejected because Atlas must not acquire domain, endpoint, file-share, or source
collection authority.

### Copy all IFI raw evidence into Atlas

Rejected because it duplicates a high-volume sensitive evidence store and
turns Atlas into a second DLP or SIEM platform.

### Rebuild IFI capabilities inside Atlas

Rejected because it creates a half-complete all-in-one platform and weakens
both products.

## Validation direction

The boundary is validated through:

- contract fixtures;
- signature and revocation tests;
- sequence and replay tests;
- minimized-data tests;
- identity correlation;
- source outage;
- stale evidence;
- authorization;
- evidence pivots; and
- independent product operation.
