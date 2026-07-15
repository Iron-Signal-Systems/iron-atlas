# PostgreSQL Migration and Ownership Model

## Status

Phase 1 Step 1 accepted under tag `phase-1-step-1-postgresql-governance-foundation-complete-v1`.

## Purpose

Iron Atlas uses an ordered manifest of immutable PostgreSQL migrations. Schema changes are applied by a dedicated migration identity, under an explicit schema-owner role, and recorded by content hash. The application service does not own the database, schema, tables, functions, or migrations.

## Ordered Manifest

`sql/schema/manifests/atlas.manifest` is the authoritative execution order. Every entry must:

- Use a three-digit, gap-free sequence.
- Refer to one migration file.
- Be immutable after acceptance.
- Begin a transaction and set bounded lock, statement, and idle-in-transaction timeouts.
- execute `SET LOCAL ROLE atlas_schema_owner` before creating or altering objects.
- Complete with `COMMIT`.

The runner holds a PostgreSQL advisory lock for the complete migration session, verifies previously recorded hashes, and rejects changed historical migrations.

## Role Topology

| Role | Login | Purpose |
|---|---:|---|
| `atlas_database_owner` | No | Owns the database only. |
| `atlas_schema_owner` | No | Owns the `atlas` schema and database objects. |
| `atlas_migrator` | Yes | Applies migrations by explicitly setting the schema-owner role. |
| `atlas_application` | Yes | Executes approved service functions and reads approved projections. |
| `atlas_readonly` | No | Read-only operational reporting role. |
| `atlas_auditor` | No | Read access to governed history and audit records. |
| `atlas_test_runner` | Yes, development only | Runs disposable database tests. |

Login credentials are not stored in Git. Production bootstrap creates roles without passwords; credential issuance and rotation are deployment responsibilities.

## Ownership Boundary

- The application role never owns database objects.
- The application role cannot apply migrations.
- The migrator is not the runtime application identity.
- Object ownership is assigned to `atlas_schema_owner`.
- `PUBLIC` receives no schema creation rights and no runtime table privileges.
- Runtime grants are applied separately from object creation by an authorized bootstrap administrator.

## Migration Execution

`tools/database/apply_migrations.sh`:

1. Reads the ordered manifest.
2. Calculates SHA-256 for every migration.
3. Opens one `psql` session.
4. Obtains an advisory lock.
5. Verifies hashes of already recorded migrations.
6. Applies only missing migrations.
7. Records version, filename, content hash, runner identity, and timestamp.
8. Releases the lock when the session exits.

Re-running the migration command must be idempotent. A historical checksum mismatch is a hard failure and is not automatically repaired.

## Historical Candidate

The Phase 0 file `000_initial_atlas_schema.sql` was a design candidate, not an accepted production migration. Phase 1 archives it under `sql/schema/candidates/` and excludes it from the executable manifest. The accepted Phase 0 tag preserves its historical form.

## Acceptance Limit

This step proves the migration framework and database-enforced governance behaviors in a disposable PostgreSQL cluster. It does not establish production credentials, high availability, backup recovery, application database integration, or production readiness.
