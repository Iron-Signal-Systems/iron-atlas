# Committed Validation Evidence

This directory contains sanitized evidence intentionally retained for implementation and acceptance decisions.

Each recorded run directory contains:

- `metadata.json` — boundary, commit, command, source, timestamp, and result;
- `environment.txt` — non-secret toolchain and host-class fingerprint;
- `validation.log` — sanitized command transcript;
- `summary.txt` — correctness and phase-gate summary lines; and
- `sha256sums.txt` — integrity records for the retained files.

Raw infrastructure evidence, credentials, environment dumps, database URLs, private keys, tokens, hostnames, and unapproved production identifiers are prohibited here.

A failed run may be committed when it is necessary to explain a defect or corrective change. An acceptance record must identify the passing canonical-clean-clone evidence used for the decision.
