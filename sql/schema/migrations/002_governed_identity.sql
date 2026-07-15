BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE TABLE atlas.actor (
    actor_id text PRIMARY KEY,
    display_name text NOT NULL CHECK (btrim(display_name) <> ''),
    actor_type text NOT NULL CHECK (actor_type IN ('HUMAN', 'SERVICE')),
    actor_status text NOT NULL DEFAULT 'ACTIVE'
        CHECK (actor_status IN ('ACTIVE', 'DISABLED', 'RETIRED')),
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp(),
    disabled_at timestamptz,
    CHECK ((actor_status = 'ACTIVE' AND disabled_at IS NULL)
        OR (actor_status <> 'ACTIVE' AND disabled_at IS NOT NULL))
);

CREATE TABLE atlas.identity_provider (
    provider_id text PRIMARY KEY,
    provider_type text NOT NULL CHECK (provider_type IN ('ACTIVE_DIRECTORY', 'OIDC', 'SAML', 'LOCAL_DEVELOPMENT')),
    display_name text NOT NULL,
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.external_identity (
    external_identity_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    provider_id text NOT NULL REFERENCES atlas.identity_provider(provider_id),
    provider_subject text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp(),
    UNIQUE (provider_id, provider_subject),
    UNIQUE (actor_id, provider_id, provider_subject)
);

ALTER TABLE atlas.actor OWNER TO atlas_schema_owner;
ALTER TABLE atlas.identity_provider OWNER TO atlas_schema_owner;
ALTER TABLE atlas.external_identity OWNER TO atlas_schema_owner;

COMMIT;
