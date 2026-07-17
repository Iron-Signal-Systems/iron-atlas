# Solo-Developer Operating Model

## Purpose

Provide a serious, sustainable engineering model for Iron Atlas while it is developed primarily by one person.

The model protects truth, security, reproducibility, and maintainability without importing team-scale ceremony that blocks product learning or misrepresents independent assurance.

## Relationship to Engineering Standards

The Iron Signal Systems `engineering-standards` repository is a source of reusable practices, validator code, templates, and maturity guidance.

Atlas does not silently inherit or automatically depend on that repository.

Atlas adopts standards through explicit, reviewable project changes.

The governing relationship is:

```text
Iron Atlas
    owns product vision
    owns architecture
    owns development sequence
    owns project-specific tests
    owns acceptance criteria
    owns release decisions

engineering-standards
    provides reusable practices
    provides optional validator tooling
    provides maturity profiles
    provides templates and guidance
    does not dictate product architecture
    does not silently alter Atlas
```

## Proportional-Assurance Principle

> **Adopt the discipline that protects truth, safety, reproducibility, and maintainability. Do not adopt ceremony that provides no proportional engineering value.**

## Truthful Status

Permitted development-assurance descriptions include:

- `EXPLORATORY`;
- `DEVELOPMENT`;
- `SELF-VALIDATED`;
- `EXTERNAL-REVIEW-PENDING`;
- `INDEPENDENTLY-REVIEWED`;
- `ACCEPTED-MILESTONE`; and
- `RELEASED`.

Self-authored and self-tested work is self-validated. Commit signing proves attribution and integrity; it does not create independent review.

## Work Classes

### Exploratory

Used to learn a format, device, protocol, or product direction.

Required:

- no secrets or protected evidence in Git;
- no production-readiness claim;
- clear limitations;
- bounded branch or worktree;
- tests for reusable behavior where practical; and
- preservation of useful findings.

Exploratory work need not carry full milestone ceremony.

### Functional Candidate

Used for reusable behavior intended to remain.

Required:

- implementation;
- focused tests;
- negative cases;
- documentation synchronization;
- visible unsupported behavior;
- formatting and static validation;
- full applicable test suite;
- failure log; and
- reviewable diff.

### Sensitive Candidate

Used for authentication, authorization, cryptography, secrets, audit, database boundaries, migrations, collection authority, evidence protection, privileged execution, provisioning, deployment security, or recovery.

Requires all functional-candidate controls plus:

- threat and abuse cases;
- adversarial tests;
- fail-closed behavior;
- secret and log review;
- concurrency or race testing where applicable;
- rollback and recovery consideration;
- exact boundary documentation; and
- stronger acceptance evidence.

### Accepted Milestone

Used for a durable project checkpoint.

Required:

- exact commit identification;
- synchronized code, tests, documentation, requirements, status, and limitations;
- clean-clone validation;
- sanitized retained evidence;
- truthful self-validation status;
- signed commit and annotated tag where available;
- rollback or recovery instructions; and
- explicit next work.

### Release

Adds:

- stable version;
- release notes;
- package and provenance records;
- vulnerability review;
- supported-platform statement;
- installation, upgrade, removal, backup, and recovery validation; and
- explicit release authorization.

## Required Practices

Atlas development shall:

1. use Git as authoritative source history;
2. keep protected infrastructure evidence and secrets out of Git;
3. preserve exact test and validation source;
4. identify the exact candidate being validated;
5. add meaningful tests for behavior changes;
6. keep parser and analyzer uncertainty visible;
7. synchronize material implementation and documentation changes;
8. create a usable local log for failed validation;
9. avoid automated commit, push, tag, merge, reset, or destructive cleanup without explicit authorization;
10. label self-validation truthfully;
11. preserve known limitations and unsupported environments; and
12. avoid weakening an accepted security boundary merely to make a later candidate pass.

## Recommended Practices

- signed commits;
- annotated signed milestone tags;
- clean-clone milestone validation;
- Go formatting, tests, vet, build, module verification, race testing where applicable, and vulnerability scanning;
- checksums for retained evidence;
- resource observations for parsers and collectors;
- deterministic fixtures;
- hostile, malformed, truncated, oversized, and conflicting input tests;
- cancellation and timeout testing;
- branch and worktree isolation; and
- concise operator-facing failure messages.

## Parallel Workstreams

A single developer may maintain multiple bounded workstreams when they remain isolated.

Examples:

- Cisco ingestion;
- FortiGate ingestion;
- product-vision documentation;
- authentication foundation;
- cross-vendor intelligence integration.

The rule is:

> **One active acceptance candidate per bounded workstream.**

Each workstream has:

- one branch or worktree;
- one declared scope;
- one accepted predecessor;
- one candidate;
- one validation boundary;
- one evidence set; and
- one explicit integration point.

Cross-workstream behavior is accepted through a separate integration candidate.

## What Must Not Be Claimed

A passing self-validation run does not establish:

- independent review;
- certification;
- regulatory compliance;
- production readiness;
- complete security;
- complete vendor coverage;
- absence of vulnerabilities;
- correct interpretation of every operational environment; or
- authorization to use protected operational evidence.

## Adoption of Future Standards

A later `engineering-standards` update is adopted only through a normal Atlas change that records:

- source standard version;
- source commit;
- imported files or practices;
- project-specific modifications;
- compatibility impact;
- validation results; and
- limitations.

Standards adoption shall not block bounded product experiments unless the experiment crosses a sensitive or accepted boundary requiring the standard.
