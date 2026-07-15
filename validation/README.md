# Validation Portability Boundary

Iron Atlas validation is a version-controlled project capability, not a workstation-local procedure.

The machine-readable external toolchain contract is in `toolchain-requirements.json`. Applicable validation scripts, disposable-environment builders, evidence recorders, phase gates, and evidence validators are committed with the implementation they validate.

Transient `latest` output under `test-framework/test-results/` remains local and replaceable. A log, summary, environment fingerprint, or machine-readable result that is deliberately retained as implementation or acceptance evidence must be sanitized and committed below `validation/evidence/`.

No implementation step may be accepted unless a clean clone from the canonical GitHub repository can execute its applicable validation using only:

- version-controlled Iron Atlas artifacts;
- declared and verifiable external toolchain requirements;
- disposable test environments created by repository scripts; and
- explicitly supplied non-repository secrets.
