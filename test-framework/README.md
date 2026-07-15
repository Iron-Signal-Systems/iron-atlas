# Test Framework

Run all current checks:

```bash
./test-framework/run_all.sh
```

Phase 1 requires PostgreSQL development tools because the test framework creates a disposable local cluster. Generated results are written under `test-framework/test-results/` and remain excluded from Git.

Correctness results remain separate from PostgreSQL version, elapsed time, temporary database size, and other resource observations. Formal acceptance evidence belongs in immutable acceptance records, not mutable `latest` files.
