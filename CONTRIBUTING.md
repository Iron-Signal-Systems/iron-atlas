# Contributing

1. Work from the `dev` branch unless a release process states otherwise.
2. Keep modules replaceable and dependency direction explicit.
3. Do not place credentials or raw infrastructure evidence in Git.
4. Add or update tests for behavior changes.
5. Run `python3 tools/validation/validate_toolchain.py`, `./test-framework/run_all.sh`, and the applicable phase gate before review.
6. Use change records and independent approval for governed production-impacting changes.
7. Keep documentation, tests, validation, retained evidence, status, and next-work statements synchronized.
8. Do not rely on workstation-only scripts, services, logs, or undeclared packages.
9. Acceptance requires the exact pushed commit to pass from a clean canonical GitHub clone.
