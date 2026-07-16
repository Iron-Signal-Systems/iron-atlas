# Iron Atlas Acceptance Records

Acceptance records freeze explicitly tested repository boundaries, evidence,
limitations, security assumptions, approval context, accepted tags, and exact
next work.

## Accepted records

- [Phase 0 Step 1 — Repository and executable baseline](PHASE-0-STEP-1-ACCEPTANCE-RECORD.md)
- [Phase 0 Step 1 acceptance errata](PHASE-0-STEP-1-ACCEPTANCE-ERRATA.md)
- [Phase 1 Step 1 — PostgreSQL governance foundation](PHASE-1-STEP-1-ACCEPTANCE-RECORD.md)
- [Phase 1 Step 2 — Go PostgreSQL runtime, identity-context, and portable-validation boundary](PHASE-1-STEP-2-ACCEPTANCE-RECORD.md)

## Templates

- [Phase 0 acceptance record template](PHASE-0-ACCEPTANCE-RECORD-TEMPLATE.md)
- [Phase 1 Step 1 acceptance record template](PHASE-1-STEP-1-ACCEPTANCE-RECORD-TEMPLATE.md)
- [Phase 1 Step 2 acceptance record template](PHASE-1-STEP-2-ACCEPTANCE-RECORD-TEMPLATE.md)
- [Phase 1 Step 3 acceptance record template](PHASE-1-STEP-3-ACCEPTANCE-RECORD-TEMPLATE.md)

## Evidence retention

This directory retains sanitized, durable acceptance records and evidence
digests. Large, sensitive, or restricted logs remain in approved evidence
storage or CI artifacts. Committed records identify external evidence by
SHA-256 digest, source commit, runner identity, and durable location.

## Governing rule

No implementation step may be accepted unless a clean clone from the canonical
GitHub repository can execute its applicable validation using only
version-controlled project artifacts, declared and verifiable external
toolchain requirements, disposable test environments, and explicitly supplied
non-repository secrets.

Acceptance does not expand the proven boundary. Historical tags remain
immutable; later-discovered record defects are corrected through explicit
errata rather than by moving or rewriting an accepted tag.

Operational infrastructure changes remain subject to independent two-person
approval according to the change-management architecture.
