# ISRAS 0.1.8 Adoption Candidate

## Status

**CANDIDATE — HOSTED VALIDATION AND FORMAL ACCEPTANCE PENDING**

This change proposes that Atlas adopt the exact published ISRAS 0.1.8
`ISRAS-SD` Go profile. It replaces the unaccepted 0.1.4 candidate as the active
standards candidate without retroactively relabeling earlier Atlas commits.

## Exact identities

- Atlas pre-candidate baseline: `9e5812f4a1184724022feb521b1e5f159c604654`
- Standards repository: `github.com/Iron-Signal-Systems/engineering-standards`
- Release: `isras-v0.1.8`
- Source commit: `0aed6ceb1db8558cb66cd81a7b6c84f2751b231e`
- Profile: `ISRAS-SD` with Go defaults
- Runtime evidence directory: `.local/isras`
- Approved deviations: None

## Generated boundary

The published release validator verified the signed release, exact six-asset
inventory, SHA-256 and SHA-512 manifests, provenance bindings, and execution
authorization before generating:

- `.isras/project.json`;
- `.isras/adoption-verification.json`;
- `.isras/check-go-format`; and
- `.github/workflows/isras-validation.yml`.

Atlas retains every stronger project-specific validation, historical phase
gate, PostgreSQL campaign, authentication campaign, evidence control, and
canonical clean-clone acceptance requirement. ISRAS does not replace those
controls.

## Signer-rotation reason

ISRAS 0.1.8 adds the current Arch development-host signing key to the governed
release-bound signer inventory. The older 0.1.4 pin correctly rejects commits
signed by that rotated key. Atlas must upgrade explicitly; it does not alter its
own allowed-signer inventory or weaken exact principal and fingerprint checks.

## Canonical repository identity prerequisite

The candidate also incorporates the two already signed commits that transition
the active repository identity from `Iron-Signal-Systems/iron-atlas` to
`Iron-Signal-Systems/atlas`. Hosted portable validation checks the canonical
origin from `REPOSITORY-ASSURANCE.json`; retaining the former active identity
would make the 0.1.8 candidate non-reproducible in the renamed GitHub
repository. Historical evidence keeps the former identity where it describes
the repository that was actually validated.

## Historical boundary

The merged 0.1.4 candidate remains reconstructable and unaccepted. The
previously recorded 1.0.1 repository-assurance boundary also remains historical.
This candidate makes no claim that either release governed work completed
before a future signed acceptance commit.

## Required acceptance evidence

Formal adoption requires:

1. a signed, pushed exact candidate commit;
2. complete Atlas validation and every committed ISRAS project command;
3. isolated canonical clean-clone validation;
4. the immutable 0.1.8 hosted reusable workflow;
5. retained hosted evidence with exact run, job, artifact, and digest; and
6. a separate signed acceptance-only change.

This candidate does not establish formal standards adoption, independent
certification, production readiness, or acceptance of any Atlas product phase.
