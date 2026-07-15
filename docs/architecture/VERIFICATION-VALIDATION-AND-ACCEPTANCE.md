# Verification, Validation, and Acceptance

## Separation

- **Verification:** the implementation matches its contract.
- **Validation:** the resulting behavior is suitable for the intended operational use.
- **Acceptance:** an authorized record freezes a tested boundary and its limitations.

A parser test does not validate a live network. A successful collection does not prove device health. A zero-finding result does not prove security.

## Test Classes

- Unit tests
- Parser fixture tests
- Golden normalization tests
- Negative and malformed-input tests
- Concurrency tests
- Authorization and approval-independence tests
- Integration tests
- Disposable PostgreSQL tests
- Device-lab tests
- Resource and performance observations
- Security and adversarial tests
- Upgrade, backup, restoration, and compromise-recovery tests
- Accessibility tests

## Acceptance Record

An acceptance record states:

- Exact commit and artifact hashes
- Scope
- Requirements and tests covered
- Results
- Resource observations
- Known warnings
- Unsupported behavior
- Security assumptions
- Operational limitations
- Reviewers and independent approvers
- Accepted tag
- Next work

Correctness and performance observations remain separate until representative evidence establishes defensible budgets.

## Reproducible Acceptance Invariant

No implementation step may be accepted unless a clean clone from the canonical GitHub repository can execute its applicable validation using only version-controlled project artifacts, declared and verifiable external toolchain requirements, disposable test environments, and explicitly supplied non-repository secrets.

Acceptance records identify the canonical commit, committed evidence path, evidence SHA-256 records, toolchain-requirements hash, and clean-clone result. Workstation-local paths and untracked scripts are not acceptance evidence.
