# Portable Validation and Canonical Repository Acceptance

## Invariant

No implementation step may be accepted unless a clean clone from the canonical GitHub repository can execute its applicable validation using only version-controlled project artifacts, declared and verifiable external toolchain requirements, disposable test environments, and explicitly supplied non-repository secrets.

The canonical repository is `https://github.com/Iron-Signal-Systems/iron-atlas.git`; the active development branch is `dev`.

## Accepted Application

Phase 1 Step 2 is the first accepted application of this invariant. Its local and canonical clean-clone evidence is committed under `validation/evidence/phase-1-step-2/`, and the accepted boundary is frozen under tag `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`.

## Required Boundary

Every implementation step must commit, in the same change set:

- source, migrations, configuration examples, and documentation;
- every applicable validator, phase gate, helper, and disposable-environment script;
- machine-readable external toolchain requirements;
- dependency integrity records such as `go.sum`;
- sanitized validation evidence deliberately retained for implementation or acceptance; and
- an acceptance record identifying the exact commit and evidence used.

A workstation path, shell history, untracked helper, locally initialized permanent database, undeclared package, or private mutable state cannot be part of the proof.

## Result Classes

`test-framework/test-results/` contains replaceable local `latest` output and remains ignored. Results deliberately retained as evidence move through `tools/validation/record_validation_evidence.sh` into `validation/evidence/`, where they are sanitized, checksummed, validated, reviewed, and committed.

Raw infrastructure evidence and secrets remain prohibited from Git.

## Canonical Clean-Clone Gate

The repository-provided verifier:

1. resolves the canonical remote `dev` commit;
2. requires it to equal the expected full commit;
3. clones the canonical GitHub repository into a new temporary directory;
4. fetches accepted tags;
5. validates the declared external toolchain;
6. downloads and verifies pinned Go modules;
7. runs the applicable repository validator from the clean clone; and
8. preserves the validator exit status through cleanup.

Acceptance is withheld until that process passes for the pushed commit.

## Secrets

Secrets are supplied only through an explicitly documented external channel. Validators must not enumerate the ambient environment or write secret values into transcripts. Evidence recording sanitizes common credential forms and then performs a separate committed-evidence secret scan.
