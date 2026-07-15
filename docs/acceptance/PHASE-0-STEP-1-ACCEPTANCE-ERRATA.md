# Phase 0 Step 1 Acceptance Errata

## Status

This errata supplements, but does not modify or move, the immutable accepted tag
`phase-0-repository-and-executable-baseline-complete-v1`.

## Corrected Record-Generation Fields

- Candidate implementation commit: `e74e4e0dcab239483327b5b329ce1f396c9837fa`
- Candidate short commit: `e74e4e0dcab2`
- Deterministic Git archive SHA-256: `eada5280cd9bd1d8c416764c6405af683b5030f9e6e2b26a267bfa59141eb26d`
- Archive command: `git archive --format=tar --prefix=iron-atlas-phase0-step1/ e74e4e0dcab239483327b5b329ce1f396c9837fa`

## Reason

The Phase 0 acceptance record generation command produced an empty short-commit
field and recorded the SHA-256 of an empty stream rather than the Git archive.
The candidate commit itself, validation decision, acceptance commit, and
annotated tag remain unchanged.

## Governance

Historical acceptance evidence is corrected through an explicit later errata.
The accepted tag is not rewritten, moved, or replaced.
