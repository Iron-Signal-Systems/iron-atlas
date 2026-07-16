## Purpose

Describe the bounded change and why it is needed.

## Scope

State what is included and explicitly prohibited.

## Architecture, authority, and trust boundary

- [ ] Architecture and trust boundaries are unchanged, or exact changes are
      documented.
- [ ] New authority, identities, privileges, routes, migrations, or
      dependencies are listed.
- [ ] Prohibited work remains absent.

## Implementation and tests

- [ ] Source changes are committed.
- [ ] Tests, hostile cases, fixtures, and expected outcomes are committed.
- [ ] Applicable toolchain and test-framework checks pass.
- [ ] The applicable Iron Atlas phase gate is identified and passes.

## Reproducible validation

- [ ] `./tools/validation/validate_portable.sh` passes.
- [ ] Applicable canonical or specialized validation is identified.
- [ ] Historical predecessor handling is correct.
- [ ] No untracked, ignored, workstation-only, or undeclared project input is
      required.
- [ ] Acceptance-bound work includes or schedules exact canonical clean-clone
      verification.

## Documentation and evidence

- [ ] Documentation is synchronized in this change set.
- [ ] Retained logs and results are sanitized and checksummed.
- [ ] Warnings, limitations, non-claims, status, and next work are updated.
- [ ] No credential, protected data, or unrestricted raw evidence is included.

## Change and approval record

State required approvals and separation-of-duties handling.

## Acceptance impact

State whether this change affects an accepted boundary, creates a candidate, or
is ordinary post-acceptance development.

## Known limitations
