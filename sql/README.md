# Iron Atlas PostgreSQL

## Status

Phase 1 Step 1 migration and governance foundation accepted under tag `phase-1-step-1-postgresql-governance-foundation-complete-v1`. Phase 1 Step 2 consumes that accepted runtime surface and does not modify accepted migrations.

The database framework is not approved for production use.

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

## Go Runtime Boundary

The Step 2 Go adapter connects only as `atlas_application`, reads approved projections, and executes only:

- `atlas.create_change_request(text, text, integer)`
- `atlas.record_approval(text, text, text)`

Acting identity is supplied through transaction-local `atlas.actor_id` context. The runtime does not edit the manifest or apply migrations.

## Production Boundary

- Never store passwords or connection URLs in these files.
- Never run the application as the database or schema owner.
- Never grant the application role migration authority.
- Never edit an accepted migration; create a new migration.
- A changed historical migration checksum is a hard failure.
- Production credential delivery, TLS provisioning, backup, restoration, and high availability require separate acceptance.
