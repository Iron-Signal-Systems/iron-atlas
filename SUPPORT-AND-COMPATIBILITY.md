# Iron Atlas Support and Compatibility Policy

## Product status

Iron Atlas is pre-alpha. It has no supported production release and makes no
availability, security-response, or operational service-level commitment.

## Repository assurance compatibility

Iron Atlas is pinned to:

- ISRAS version: `1.0.1`
- signed tag: `isras-v1.0.1`
- exact commit: `c379417720faa595fa5cb89a1dfdb2259d6cb95e`
- profile: `go-documentation-generation`
- adoption level: `RECORDED`

An ISRAS update requires a separately reviewed Iron Atlas change with
applicable validation. Iron Atlas must not follow the standards repository's
floating `dev` branch.

## Development and release branches

- `dev` contains active Iron Atlas development.
- `main` is the Iron Atlas release branch and is not advanced merely because
  `dev` changes.
- Accepted Iron Atlas tags and historical source boundaries remain immutable.
- A branch name or version string alone does not prove acceptance.

## Compatibility claims

Compatibility is limited to environments, formats, integrations, and
dependencies explicitly declared and validated by the applicable Iron Atlas
phase boundary. Undeclared platforms and production readiness are not implied.
## ISRAS 0.1.4 adoption candidate

Atlas is evaluating the exact published `isras-v0.1.4` release at source commit
`c9345d6d731600df7bd4ba4a133c07265db55e5a`. The generated project pin and
hosted workflow are a candidate until separate prospective acceptance. The
previously recorded 1.0.1 boundary remains historical and no production or
independent-review claim is added.
