# ISRAS 0.1.4 Adoption Candidate

## Status

**CANDIDATE — NOT FORMALLY ACCEPTED**

This change proposes that Iron Atlas adopt the exact published ISRAS 0.1.4
`ISRAS-SD` Go profile.

## Exact identities

- Atlas pre-candidate baseline: `56495408aa51d886a214859f814ebae1e30e2b92`
- Standards repository: `github.com/Iron-Signal-Systems/engineering-standards`
- Release: `isras-v0.1.4`
- Source commit: `c9345d6d731600df7bd4ba4a133c07265db55e5a`
- Profile: `ISRAS-SD` with Go defaults
- Runtime evidence directory: `.local/isras`

## Candidate boundary

The official release validator generated exactly:

- `.isras/project.json`;
- `.isras/adoption-verification.json`;
- `.isras/check-go-format`; and
- `.github/workflows/isras-validation.yml`.

The reusable workflow is pinned to the exact standards source commit. Existing
Atlas-specific validation, immutable historical gates, PostgreSQL campaigns,
authentication campaigns, evidence controls, and acceptance records remain in
force and are not replaced by ISRAS.

## Historical boundary

The previously recorded 1.0.1 repository-assurance baseline remains historical.
This candidate does not relabel earlier commits or claim that 0.1.4 governed
work completed before the future acceptance commit.

## Acceptance requirements

Formal adoption requires a separate signed acceptance change after this exact
candidate is merged and its exact pushed commit passes:

1. complete Atlas validation;
2. ISRAS commit validation;
3. the pinned hosted reusable workflow;
4. exact release and artifact verification; and
5. retained hosted evidence.

No independent certification, production readiness, or product-phase acceptance
is claimed.

## Module-consistency correction

ISRAS adoption exposed that `golang.org/x/oauth2 v0.36.0` is imported directly by the OIDC implementation but remained classified as an indirect dependency. The candidate applies the exact `go mod tidy` correction. The module version and `go.sum` remain unchanged; no dependency upgrade or runtime behavior change is claimed.
