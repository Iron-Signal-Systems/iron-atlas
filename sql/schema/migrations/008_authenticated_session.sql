BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE TABLE atlas.authenticated_session (
    session_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    identifier_digest bytea NOT NULL UNIQUE
        CHECK (octet_length(identifier_digest) = 32),
    provider_id text NOT NULL,
    provider_subject text NOT NULL,
    actor_id text NOT NULL,
    created_at timestamptz NOT NULL,
    authenticated_at timestamptz NOT NULL,
    last_activity_at timestamptz NOT NULL,
    idle_expires_at timestamptz NOT NULL,
    absolute_expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    revocation_reason text,
    rotation_parent_digest bytea,
    authentication_context text,
    authentication_methods text[] NOT NULL DEFAULT ARRAY[]::text[],
    mfa_authenticated boolean NOT NULL DEFAULT false,
    mfa_authenticated_at timestamptz,
    security_policy_version text NOT NULL,
    CONSTRAINT authenticated_session_provider_fk
        FOREIGN KEY (provider_id)
        REFERENCES atlas.identity_provider(provider_id),
    CONSTRAINT authenticated_session_actor_fk
        FOREIGN KEY (actor_id)
        REFERENCES atlas.actor(actor_id),
    CHECK (provider_id = btrim(provider_id)
        AND octet_length(provider_id) BETWEEN 1 AND 256),
    CHECK (provider_subject = btrim(provider_subject)
        AND octet_length(provider_subject) BETWEEN 1 AND 512),
    CHECK (actor_id = btrim(actor_id)
        AND octet_length(actor_id) BETWEEN 1 AND 256),
    CHECK (authenticated_at <= created_at + interval '2 minutes'),
    CHECK (last_activity_at >= created_at),
    CHECK (idle_expires_at > last_activity_at),
    CHECK (absolute_expires_at > created_at),
    CHECK (idle_expires_at <= absolute_expires_at),
    CHECK ((revoked_at IS NULL AND revocation_reason IS NULL)
        OR (revoked_at IS NOT NULL
            AND revocation_reason IS NOT NULL
            AND revocation_reason = btrim(revocation_reason)
            AND octet_length(revocation_reason) BETWEEN 1 AND 256)),
    CHECK (rotation_parent_digest IS NULL
        OR octet_length(rotation_parent_digest) = 32),
    CHECK (authentication_context IS NULL
        OR (authentication_context = btrim(authentication_context)
            AND octet_length(authentication_context) BETWEEN 1 AND 256)),
    CHECK (cardinality(authentication_methods) <= 16),
    CHECK ((mfa_authenticated AND mfa_authenticated_at IS NOT NULL)
        OR (NOT mfa_authenticated AND mfa_authenticated_at IS NULL)),
    CHECK (mfa_authenticated_at IS NULL
        OR mfa_authenticated_at <= created_at + interval '2 minutes'),
    CHECK (security_policy_version = btrim(security_policy_version)
        AND octet_length(security_policy_version) BETWEEN 1 AND 128)
);

ALTER TABLE atlas.authenticated_session OWNER TO atlas_schema_owner;

REVOKE ALL ON TABLE atlas.authenticated_session FROM PUBLIC;
REVOKE ALL ON SEQUENCE atlas.authenticated_session_session_id_seq FROM PUBLIC;

CREATE OR REPLACE FUNCTION atlas.create_authenticated_session(
    requested_identifier_digest bytea,
    requested_provider_id text,
    requested_provider_subject text,
    requested_actor_id text,
    requested_authenticated_at timestamptz,
    requested_idle_lifetime_seconds integer,
    requested_absolute_lifetime_seconds integer,
    requested_authentication_context text,
    requested_authentication_methods text[],
    requested_mfa_authenticated boolean,
    requested_mfa_authenticated_at timestamptz,
    requested_security_policy_version text
)
RETURNS TABLE (
    created_at timestamptz,
    last_activity_at timestamptz,
    idle_expires_at timestamptz,
    absolute_expires_at timestamptz
)
LANGUAGE sql
VOLATILE
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
    INSERT INTO atlas.authenticated_session AS created_session (
        identifier_digest,
        provider_id,
        provider_subject,
        actor_id,
        created_at,
        authenticated_at,
        last_activity_at,
        idle_expires_at,
        absolute_expires_at,
        authentication_context,
        authentication_methods,
        mfa_authenticated,
        mfa_authenticated_at,
        security_policy_version
    )
    SELECT requested_identifier_digest,
           ip.provider_id,
           ei.provider_subject,
           a.actor_id,
           transaction_timestamp(),
           requested_authenticated_at,
           transaction_timestamp(),
           transaction_timestamp()
               + make_interval(secs => requested_idle_lifetime_seconds),
           transaction_timestamp()
               + make_interval(secs => requested_absolute_lifetime_seconds),
           requested_authentication_context,
           COALESCE(requested_authentication_methods, ARRAY[]::text[]),
           requested_mfa_authenticated,
           requested_mfa_authenticated_at,
           requested_security_policy_version
      FROM atlas.identity_provider AS ip
      JOIN atlas.external_identity AS ei
        ON ei.provider_id = ip.provider_id
      JOIN atlas.actor AS a
        ON a.actor_id = ei.actor_id
     WHERE octet_length(requested_identifier_digest) = 32
       AND requested_provider_id IS NOT NULL
       AND requested_provider_subject IS NOT NULL
       AND requested_actor_id IS NOT NULL
       AND requested_provider_id = btrim(requested_provider_id)
       AND requested_provider_subject = btrim(requested_provider_subject)
       AND requested_actor_id = btrim(requested_actor_id)
       AND octet_length(requested_provider_id) BETWEEN 1 AND 256
       AND octet_length(requested_provider_subject) BETWEEN 1 AND 512
       AND octet_length(requested_actor_id) BETWEEN 1 AND 256
       AND ip.provider_id = requested_provider_id
       AND ip.active
       AND ei.provider_subject = requested_provider_subject
       AND ei.actor_id = requested_actor_id
       AND a.actor_id = requested_actor_id
       AND a.actor_status = 'ACTIVE'
       AND requested_authenticated_at IS NOT NULL
       AND requested_authenticated_at <= transaction_timestamp() + interval '2 minutes'
       AND requested_idle_lifetime_seconds BETWEEN 1 AND 86400
       AND requested_absolute_lifetime_seconds BETWEEN 1 AND 86400
       AND requested_idle_lifetime_seconds <= requested_absolute_lifetime_seconds
       AND (requested_authentication_context IS NULL
            OR (requested_authentication_context = btrim(requested_authentication_context)
                AND octet_length(requested_authentication_context) BETWEEN 1 AND 256))
       AND cardinality(COALESCE(requested_authentication_methods, ARRAY[]::text[])) <= 16
       AND NOT EXISTS (
            SELECT 1
              FROM unnest(COALESCE(requested_authentication_methods, ARRAY[]::text[])) AS methods(method)
             WHERE method IS NULL
                OR method <> btrim(method)
                OR octet_length(method) NOT BETWEEN 1 AND 64
                OR method ~ '[[:cntrl:]]'
       )
       AND cardinality(COALESCE(requested_authentication_methods, ARRAY[]::text[])) = (
            SELECT count(DISTINCT method)::integer
              FROM unnest(COALESCE(requested_authentication_methods, ARRAY[]::text[])) AS methods(method)
       )
       AND requested_mfa_authenticated IS NOT NULL
       AND ((requested_mfa_authenticated
             AND requested_mfa_authenticated_at IS NOT NULL)
            OR (NOT requested_mfa_authenticated
                AND requested_mfa_authenticated_at IS NULL))
       AND (requested_mfa_authenticated_at IS NULL
            OR requested_mfa_authenticated_at <= transaction_timestamp() + interval '2 minutes')
       AND requested_security_policy_version IS NOT NULL
       AND requested_security_policy_version = btrim(requested_security_policy_version)
       AND octet_length(requested_security_policy_version) BETWEEN 1 AND 128
    RETURNING created_session.created_at,
              created_session.last_activity_at,
              created_session.idle_expires_at,
              created_session.absolute_expires_at;
$$;

CREATE OR REPLACE FUNCTION atlas.authenticate_session(
    requested_identifier_digest bytea
)
RETURNS TABLE (
    provider_id text,
    provider_subject text,
    actor_id text,
    authenticated_at timestamptz,
    created_at timestamptz,
    last_activity_at timestamptz,
    idle_expires_at timestamptz,
    absolute_expires_at timestamptz,
    authentication_context text,
    authentication_methods text[],
    mfa_authenticated boolean,
    mfa_authenticated_at timestamptz,
    security_policy_version text
)
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
    SELECT s.provider_id,
           s.provider_subject,
           s.actor_id,
           s.authenticated_at,
           s.created_at,
           s.last_activity_at,
           s.idle_expires_at,
           s.absolute_expires_at,
           s.authentication_context,
           s.authentication_methods,
           s.mfa_authenticated,
           s.mfa_authenticated_at,
           s.security_policy_version
      FROM atlas.authenticated_session AS s
      JOIN atlas.identity_provider AS ip
        ON ip.provider_id = s.provider_id
       AND ip.active
      JOIN atlas.external_identity AS ei
        ON ei.provider_id = s.provider_id
       AND ei.provider_subject = s.provider_subject
       AND ei.actor_id = s.actor_id
      JOIN atlas.actor AS a
        ON a.actor_id = s.actor_id
       AND a.actor_status = 'ACTIVE'
     WHERE requested_identifier_digest IS NOT NULL
       AND octet_length(requested_identifier_digest) = 32
       AND s.identifier_digest = requested_identifier_digest
       AND s.revoked_at IS NULL
       AND transaction_timestamp() < s.idle_expires_at
       AND transaction_timestamp() < s.absolute_expires_at;
$$;

ALTER FUNCTION atlas.create_authenticated_session(
    bytea,
    text,
    text,
    text,
    timestamptz,
    integer,
    integer,
    text,
    text[],
    boolean,
    timestamptz,
    text
) OWNER TO atlas_schema_owner;

ALTER FUNCTION atlas.authenticate_session(bytea)
    OWNER TO atlas_schema_owner;

REVOKE ALL ON FUNCTION atlas.create_authenticated_session(
    bytea,
    text,
    text,
    text,
    timestamptz,
    integer,
    integer,
    text,
    text[],
    boolean,
    timestamptz,
    text
) FROM PUBLIC;

REVOKE ALL ON FUNCTION atlas.authenticate_session(bytea)
FROM PUBLIC;

COMMIT;
