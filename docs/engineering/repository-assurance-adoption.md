# Iron Atlas Repository Assurance Adoption

Iron Atlas records adoption of the Iron Signal Repository Assurance Standard.

## Exact standard identity

- Standard version: `1.0.1`
- Signed release tag: `isras-v1.0.1`
- Exact accepted commit: `c379417720faa595fa5cb89a1dfdb2259d6cb95e`
- Standard repository: `Iron-Signal-Systems/engineering-standards`
- Standard source-manifest SHA-256:
  `90fd3eb3d19a0d4d846bedf8b9657454538e81d3cb14fd62e8775f9aa206c7c1`
- Adoption date: `2026-07-16`
- Profile: `go-documentation-generation`

The signed tag and exact commit are authoritative. Iron Atlas does not adopt
from the standards repository's floating `dev` branch.

## Governing rule

A change is complete only when its exact pushed commit can be reconstructed,
validated, and evidenced from the canonical repository using declared
environments and committed project-owned assets.

## Native-first boundary

Portable validation runs directly on approved hosts. Canonical and specialized
validation may use native hosts or disposable virtual machines. Containers are
optional and are not the sole validation path unless an accepted deployment
model requires them.

## Adoption status

- Adoption level: `RECORDED`
- Required checks: observation mode
- Direct pushes to `dev`: prohibited by policy
- Independent human review: not claimed
- Production readiness: not claimed
- Security evaluation: not established by this adoption

Observation mode collects results without making the new ISRAS workflow a
required repository rule.

## Template compatibility note

The accepted v1.0.1 repository-assurance template carried the prior `1.0.0`
value in its generated manifest. Iron Atlas deliberately records `1.0.1` while
retaining the exact accepted v1.0.1 commit pin. The publisher-oriented generic
support text was also replaced with an Iron Atlas-specific pre-alpha support
policy.

## Pinned Go security tool

The `go-documentation-generation` profile requires Go vulnerability analysis. Iron Atlas pins `golang.org/x/vuln/cmd/govulncheck@v1.6.0` in `tools/go-tools.lock.json`, including the module and module-file checksums.

The bootstrap verifies the downloaded module identity before installation and verifies the installed binary's embedded package path, version, and checksum. Portable validation uses `.isras-go-tools/bin` instead of an ambient workstation installation.

`govulncheck` queries the current Go vulnerability database. “No vulnerabilities found” is a time-bounded observation for the tested source, scanner version, and database state; it is not proof that the repository is vulnerability-free.

## Next maturity work

Before advancing beyond `RECORDED`, Iron Atlas must separately review and
register its accepted historical phase checkpoints, validate the pushed
adoption commit from a clean clone, observe the portable workflow across
approved systems, and define any repository-specific exceptions.
