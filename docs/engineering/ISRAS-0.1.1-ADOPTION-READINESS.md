# ISRAS 0.1.1 Adoption Readiness

## Status

Planning and migration-readiness record.

This document does not claim that Iron Atlas has adopted ISRAS 0.1.1.

## 1. Current Atlas state

Iron Atlas currently retains a project-owned repository-assurance framework
recorded from the earlier ISRAS v1.0.1 baseline.

That framework remains part of the Atlas repository until an explicit,
reviewable replacement is validated and accepted.

## 2. Target direction

Atlas intends to adopt an immutable release of the restarted **Iron Signal
Repository Assurance Standard (ISRAS)** through its pinned-project framework.

Atlas shall not:

- follow `engineering-standards/dev`;
- follow `engineering-standards/main` automatically;
- claim adoption from a prerelease;
- copy the new validator source into Atlas;
- add the standards tool to the Atlas application dependency graph;
- replace current validation before the replacement proves compatibility; or
- create a placeholder pin using unknown release digests.

## 3. Release prerequisite

Adoption work begins only after the target release has:

- a stable version;
- exact signed source commit;
- verified signed annotated tag;
- complete deterministic release artifact set;
- provenance;
- checksum manifest;
- published release;
- re-download and remote-byte verification;
- accepted project-pin schema; and
- accepted consuming-project adoption instructions.

## 4. Atlas inventory

The adoption candidate inventories and maps:

- `REPOSITORY-ASSURANCE.json`;
- existing project-owned validation;
- current source manifests;
- toolchain requirements;
- Go tool pinning;
- clean-clone validation;
- hosted workflows;
- phase gates;
- acceptance evidence;
- release and recovery documentation;
- project-specific exceptions; and
- commands required by Atlas.

Adoption is not permission to reorganize working application code merely to
match a reference layout.

## 5. Intended project-owned artifacts

When the released contract authorizes adoption, Atlas expects:

```text
.isras/project.json
tools/isras
project-owned command declarations
project-owned documentation
project-owned bounded exceptions
CI integration pinned to an immutable workflow commit
```

The exact contents and digests come only from the accepted release.

## 6. Candidate project commands

The Atlas pin should declare bounded commands equivalent to:

- repository validation;
- complete test framework;
- focused Go tests;
- Go race tests where applicable;
- Go vet;
- Go module verification;
- vulnerability analysis;
- disposable PostgreSQL tests;
- manifest validation;
- documentation validation; and
- clean-clone acceptance validation.

The final command names and argv must match the accepted ISRAS contract.

## 7. Migration stages

### Stage A — Readiness

- retain current Atlas validation;
- document the target release;
- identify command mapping;
- identify conflicts and exceptions;
- do not claim adoption.

### Stage B — Release verification

- verify tag and exact commit;
- verify all required assets;
- verify checksums and provenance;
- verify validator embedded identity;
- retain private verification evidence.

### Stage C — Candidate pin

- create a dedicated adoption branch;
- add the exact committed pin;
- add the release launcher or wrapper;
- declare bounded project commands;
- preserve existing validation in parallel;
- run read-only pin validation.

### Stage D — Compatibility validation

- run Atlas validation directly;
- run the release validator against `/src/iron-atlas`;
- compare local and hosted behavior;
- run clean-clone validation;
- validate historical accepted boundaries;
- record gaps and exceptions.

### Stage E — Acceptance

- synchronize architecture, requirements, validation, roadmap, limitations, and
  status;
- identify the exact pushed candidate;
- retain sanitized evidence;
- accept through an annotated milestone tag;
- remove deprecated project-owned copies only when explicitly authorized by a
  later change.

## 8. Fail-closed rules

Adoption fails when:

- release identity differs from the pin;
- origin differs from the pin;
- target commit differs from the pin;
- pin is modified or staged;
- required artifact digest differs;
- validator identity differs;
- a declared command differs;
- repository state violates the command contract;
- secret or evidence handling regresses;
- an accepted Atlas gate is weakened; or
- complete Atlas validation does not pass.

## 9. Current decision

The compromise blast-radius and IFI documentation may proceed independently.

The ISRAS pin must wait for a verified stable release and accepted adoption
boundary. This avoids delaying product architecture while also avoiding a false
or floating standards claim.
