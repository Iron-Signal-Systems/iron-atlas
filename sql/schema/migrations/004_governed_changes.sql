BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE OR REPLACE FUNCTION atlas.current_actor_id()
RETURNS text
LANGUAGE plpgsql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
DECLARE
    resolved_actor_id text;
BEGIN
    resolved_actor_id := nullif(current_setting('atlas.actor_id', true), '');
    IF resolved_actor_id IS NULL THEN
        RAISE EXCEPTION 'atlas.actor_id transaction context is required';
    END IF;
    IF NOT EXISTS (
        SELECT 1
        FROM atlas.actor AS a
        WHERE a.actor_id = resolved_actor_id
          AND a.actor_status = 'ACTIVE'
    ) THEN
        RAISE EXCEPTION 'active actor context is required';
    END IF;
    RETURN resolved_actor_id;
END;
$$;

CREATE OR REPLACE FUNCTION atlas.actor_has_authority(
    requested_actor_id text,
    requested_authority text,
    as_of timestamptz DEFAULT transaction_timestamp()
)
RETURNS boolean
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
    SELECT EXISTS (
        SELECT 1
        FROM atlas.role_binding rb
        JOIN atlas.role_authority ra USING (role_code)
        JOIN atlas.role_definition rd USING (role_code)
        JOIN atlas.authority_definition ad USING (authority_code)
        WHERE rb.actor_id = requested_actor_id
          AND ra.authority_code = requested_authority
          AND rd.active
          AND rb.valid_from <= as_of
          AND (rb.valid_until IS NULL OR rb.valid_until > as_of)
    );
$$;

CREATE TABLE atlas.change_request (
    change_id text PRIMARY KEY,
    title text NOT NULL CHECK (btrim(title) <> ''),
    requester_actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    risk text NOT NULL CHECK (risk IN ('LOW', 'MODERATE', 'HIGH', 'CRITICAL')),
    status text NOT NULL DEFAULT 'PENDING_APPROVAL'
        CHECK (status IN ('DRAFT', 'PENDING_APPROVAL', 'APPROVED', 'REJECTED', 'IMPLEMENTING', 'VALIDATING', 'ACCEPTED', 'CLOSED', 'CANCELLED')),
    required_approvals integer NOT NULL CHECK (required_approvals >= 1 AND required_approvals <= 10),
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp(),
    updated_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.change_status_history (
    change_status_history_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    change_id text NOT NULL REFERENCES atlas.change_request(change_id),
    previous_status text,
    new_status text NOT NULL,
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    reason text NOT NULL,
    recorded_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE OR REPLACE FUNCTION atlas.create_change_request(
    requested_change_id text,
    requested_title text,
    requested_approvals integer DEFAULT 1
)
RETURNS text
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
DECLARE
    current_actor_id_value text := atlas.current_actor_id();
BEGIN
    IF NOT atlas.actor_has_authority(current_actor_id_value, 'change.request') THEN
        RAISE EXCEPTION 'actor lacks change.request authority';
    END IF;
    INSERT INTO atlas.change_request(
        change_id, title, requester_actor_id, risk, status, required_approvals
    ) VALUES (
        requested_change_id, requested_title, current_actor_id_value,
        'MODERATE', 'PENDING_APPROVAL', requested_approvals
    );
    INSERT INTO atlas.change_status_history(
        change_id, previous_status, new_status, actor_id, reason
    ) VALUES (
        requested_change_id, NULL, 'PENDING_APPROVAL',
        current_actor_id_value, 'change requested'
    );
    RETURN requested_change_id;
END;
$$;

ALTER FUNCTION atlas.current_actor_id() OWNER TO atlas_schema_owner;
ALTER FUNCTION atlas.actor_has_authority(text, text, timestamptz) OWNER TO atlas_schema_owner;
ALTER FUNCTION atlas.create_change_request(text, text, integer) OWNER TO atlas_schema_owner;
ALTER TABLE atlas.change_request OWNER TO atlas_schema_owner;
ALTER TABLE atlas.change_status_history OWNER TO atlas_schema_owner;

COMMIT;
