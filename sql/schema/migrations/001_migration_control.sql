BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE SCHEMA IF NOT EXISTS atlas AUTHORIZATION atlas_schema_owner;
REVOKE ALL ON SCHEMA atlas FROM PUBLIC;

CREATE TABLE atlas.schema_migration (
    migration_version integer PRIMARY KEY CHECK (migration_version > 0),
    migration_filename text NOT NULL UNIQUE,
    content_sha256 text NOT NULL CHECK (content_sha256 ~ '^[0-9a-f]{64}$'),
    applied_by text NOT NULL,
    applied_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

ALTER TABLE atlas.schema_migration OWNER TO atlas_schema_owner;
COMMENT ON TABLE atlas.schema_migration IS
    'Append-only record of applied Iron Atlas schema migrations.';

COMMIT;
