# Canonical Clean-Clone Validation

## Local Candidate

Run the complete framework and phase gate from a clean working tree.

```bash
./test-framework/run_all.sh
./tools/validation/phase-gates/validate_phase1_step2.sh
```

## Record Retained Evidence

The recorder runs the command while the repository is clean, sanitizes the transcript, writes metadata and an environment fingerprint, computes SHA-256 records, validates the evidence, and only then moves the run into the repository.

```bash
./tools/validation/record_validation_evidence.sh   phase-1-step-2/local-phase-gate   ./tools/validation/phase-gates/validate_phase1_step2.sh
```

Review and commit the resulting `validation/evidence/...` directory when the run is part of the implementation or acceptance record.

## Canonical Repository Verification

After the exact commit is pushed to `origin/dev`, run:

```bash
commit="$(git rev-parse HEAD)"
IRON_ATLAS_VALIDATION_SOURCE=canonical-clean-clone ./tools/validation/record_validation_evidence.sh   phase-1-step-2/canonical-clean-clone   ./tools/validation/verify_canonical_clone.sh   "$commit"   ./tools/validation/phase-gates/validate_phase1_step2.sh
```

The verifier refuses to test a commit other than the exact current canonical `dev` commit.

## New Box

```bash
git clone https://github.com/Iron-Signal-Systems/iron-atlas.git
cd iron-atlas
git switch dev
python3 tools/validation/validate_toolchain.py
go mod download
go mod verify
./test-framework/run_all.sh
```

No repository secret is required for the current disposable Phase 1 Step 2 boundary.
