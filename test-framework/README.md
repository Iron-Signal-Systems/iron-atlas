# Test Framework

Run all current checks:

```bash
./test-framework/run_all.sh
```

Phase 1 requires PostgreSQL development tools because the test framework creates a disposable local cluster. Phase 1 Step 2 also requires access to the pinned Go module dependencies during the initial `go mod download` or `go mod tidy` operation.

The disposable database runner executes the SQL governance suite and the Go PostgreSQL integration packages. Mutable `latest` results are written under `test-framework/test-results/` and remain local. A transcript or result deliberately retained for implementation or acceptance is recorded with `tools/validation/record_validation_evidence.sh`, sanitized, checksummed, validated, and committed under `validation/evidence/`.

The framework also executes a self-contained phase-gate regression test proving that a failing or missing historical predecessor validator cannot be masked by successful temporary-clone cleanup.

Phase 1 Step 3 authentication coverage includes the contract, middleware
foundation, governed actor resolution, OIDC ID-token verification, and the
bounded authorization-code and PKCE transaction candidate. Browser routes,
durable sessions, CSRF, logout, and trusted-proxy campaigns remain later work.

Correctness results remain separate from PostgreSQL version, elapsed time, temporary database size, identity-isolation operation counts, and other resource observations. Formal acceptance links immutable records to committed evidence and requires a passing clean clone of the exact canonical GitHub commit.
