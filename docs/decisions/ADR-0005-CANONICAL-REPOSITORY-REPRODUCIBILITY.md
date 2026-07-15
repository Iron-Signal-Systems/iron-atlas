# ADR-0005 — Canonical Repository Reproducibility

## Status

Accepted for all Phase 1 Step 2 and later implementation and acceptance work.

## Decision

Iron Atlas treats validation portability as an acceptance requirement. A step is not accepted from a developer workstation alone. Its exact pushed commit must pass the applicable validator from a clean clone of the canonical GitHub repository.

Validation scripts, phase gates, helpers, toolchain requirements, disposable test-environment builders, dependency integrity records, and deliberately retained sanitized evidence are version-controlled project artifacts.

## Consequences

- A fresh development box can clone or pull the repository, install the declared toolchain, and execute the same gates.
- Workstation-only scripts and retained logs are defects, not convenience files.
- Mutable local `latest` results remain diagnostic only and cannot support acceptance.
- Canonical verification requires network access to GitHub and dependency sources unless an independently verified mirror is explicitly governed later.
- Secrets remain external and must be explicitly supplied without entering Git or retained logs.
- Historical accepted tags remain immutable; earlier boundaries that predate this rule are documented honestly rather than rewritten.
