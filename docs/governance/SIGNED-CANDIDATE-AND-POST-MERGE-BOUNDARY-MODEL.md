# Signed Candidate and Post-Merge Boundary Model

## Status

Normative repository-governance alignment candidate.

## Purpose

Iron Atlas requires SSH-signed commits from an allowed principal for formal trust boundaries. GitHub-created merge commits do not carry the maintainer's personal SSH signature, so candidate validation, merge integration, and signed closure are distinct.

## Required sequence

```text
SSH-signed candidate commit
→ pull request validation on the exact candidate
→ GitHub merge commit
→ merged-tree validation
→ SSH-signed empty post-merge boundary commit
→ hosted validation on the exact signed boundary
```

## Candidate trust

The candidate is signed by an allowed principal, descends from the declared predecessor, contains the exact reviewed tree, passes applicable validation, and remains identifiable as a PR parent after merge.

## GitHub merge commit

A GitHub-created merge commit is an integration artifact and is not represented as personally SSH-signed by the maintainer. A signer-trust job may fail temporarily when triggered on that exact merge commit. The failure is retained as accurate evidence that the merge commit is not the formal signed boundary.

## Post-merge closure

The maintainer fast-forwards local `dev` to the canonical merge, validates the merged tree, creates an SSH-signed empty commit whose parent is the merge, pushes it to `dev`, and requires hosted validation—including signer trust—to pass on that exact commit.

The empty commit changes no tree content. It identifies the exact merged tree as the signed, validated repository boundary.

## Prohibitions

Do not disable signer verification, rewrite or force-push accepted history, label GitHub's merge commit as maintainer-signed, delete failed merge-run evidence, create the boundary before merged-tree validation, or start new bounded work before the signed boundary is pushed and green.

## Future ISRAS improvement

ISRAS may later classify a GitHub merge commit as an expected intermediate after verifying a trusted signed candidate parent and merge metadata. That classification must not replace the final signed post-merge boundary or weaken exact-target signer trust.

## Explicit nonclaims

A signed boundary proves signer identity and the tested repository state. It does not establish independent review, legal approval, production readiness, or correctness outside the validation scope.
