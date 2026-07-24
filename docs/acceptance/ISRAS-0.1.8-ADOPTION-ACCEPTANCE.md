# ISRAS 0.1.8 Adoption Acceptance

## Decision

**FORMALLY ACCEPTED — EFFECTIVE WHEN THIS SIGNED ACCEPTANCE CHANGE ENTERS `dev`**

Atlas adopts the exact published ISRAS 0.1.8 `ISRAS-SD` Go profile as its
repository-assurance floor. This is a prospective acceptance of the exact
candidate below. It does not retroactively govern earlier commits or accept an
Atlas product phase.

## Accepted identities

- Candidate commit: `d7398968d1e900854345e0114c221e1f4e70deee`
- Candidate tree: `d06416839b6a43ddf24bab9682f02d3a402fdd32`
- Candidate signature status: good SSH signature
- Candidate signer principal: `kb2vhn@gmail.com`
- Candidate signer fingerprint:
  `SHA256:CiONRPnsf/rG0Ix5LJmYeoCVdI4d1kRfYtQfQTp/vDQ`
- Canonical repository: `https://github.com/Iron-Signal-Systems/atlas.git`
- Candidate branch: `work/adopt-isras-v0.1.8`
- Standards release: `isras-v0.1.8`
- Standards source commit: `0aed6ceb1db8558cb66cd81a7b6c84f2751b231e`
- Profile: `ISRAS-SD` with Go defaults
- Approved deviations: None

## Hosted evidence

All hosted checks below completed successfully against the exact candidate
commit on 2026-07-24:

| Boundary | Run | Job | Result |
| --- | ---: | ---: | --- |
| ISRAS pinned project, push | `30086984075` | `89461347729` | PASS |
| ISRAS pinned project, pull request | `30086986274` | `89461354777` | PASS |
| Atlas portable validation | `30086985760` | `89461353383` | PASS |
| Atlas repository and complete framework | `30086985772` | `89461353215` | PASS |

The retained ISRAS artifacts are:

- artifact `8594134018`, `isras-evidence-30086984075`,
  `sha256:8593a0c61c9639bd046b1682dcc422243b13ce6b97e77b53ee6c4a7d5caf0776`;
- artifact `8594142474`, `isras-evidence-30086986274`,
  `sha256:72073f3aed02a3d32c3e6e3c6828e9008226dc020cdd16cde3d4d101ed5ac9fc`.

GitHub reported both artifacts bound to candidate commit
`d7398968d1e900854345e0114c221e1f4e70deee`. Their configured retention ends
2026-08-23; the immutable run, job, artifact, and digest identities remain in
this record.

## Canonical clean-clone evidence

An isolated clone of
`https://github.com/Iron-Signal-Systems/atlas.git` at branch
`work/adopt-isras-v0.1.8` resolved to the exact candidate commit with a clean
tracked working tree before tool bootstrap. The repository-pinned tool
environment was then created and `tools/validation/validate_portable.sh`
passed.

That validation proved the canonical repository identity, required environment
and assurance artifacts, policy checks, module verification, ordinary and race
tests, whitespace and syntax checks, and the pinned `govulncheck` campaign.
The vulnerability campaign reported no known vulnerabilities at validation
time.

## Preserved stronger controls

ISRAS is a floor, not a replacement for Atlas controls. Atlas retains its
historical exact-boundary phase gates, isolated predecessor revalidation,
PostgreSQL security and concurrency campaigns, authentication and OIDC
campaigns, evidence controls, repository-specific validation, and formal
product-phase acceptance process.

The historical ISRAS 1.0.1 recorded boundary and merged but unaccepted 0.1.4
candidate remain reconstructable historical facts. Neither is relabeled by this
decision.

## Claim boundary

This decision establishes self-governed repository-assurance adoption. It does
not claim independent certification, production readiness, complete vendor
coverage, security accreditation, or completion or acceptance of any Atlas
product phase.

## Post-merge signing boundary

GitHub's merge operation produced merge commit `414c0c258dc32631d7e7363fa544dc2d01da8985`
without an SSH signature. That commit is retained as historical merge
topology, but it is not an ISRAS acceptance boundary. The signed commit
carrying this amendment is the authoritative post-merge `dev` boundary and is
the commit validated by the hosted ISRAS workflow after this change.
