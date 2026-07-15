# Iron Atlas PostgreSQL

## Status

Phase 1 Step 1 implementation candidate. The database framework is not approved for production use.

## Layout

- `bootstrap/`: role creation and post-migration runtime grants
- `schema/manifests/atlas.manifest`: authoritative migration order
- `schema/migrations/`: executable immutable migrations
- `schema/candidates/`: archived, non-executable design candidates
- `tests/`: SQL fixtures and documentation used by the disposable database harness

## Apply to a Development Database

Set normal PostgreSQL connection environment variables and run:

```bash
./tools/database/apply_migrations.sh
```

An authorized bootstrap administrator must create roles before migrations and apply runtime grants afterward. The disposable test harness performs those steps automatically in an isolated cluster.

## Production Boundary

- Never store passwords in these files.
- Never run the application as the database or schema owner.
- Never grant the application role migration authority.
- Never edit an accepted migration; create a new migration.
- A changed historical migration checksum is a hard failure.
