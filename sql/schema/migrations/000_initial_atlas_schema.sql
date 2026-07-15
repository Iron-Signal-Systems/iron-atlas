BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';

CREATE SCHEMA IF NOT EXISTS atlas;

CREATE TABLE atlas.actor (
    actor_id text PRIMARY KEY,
    display_name text NOT NULL,
    actor_type text NOT NULL CHECK (actor_type IN ('HUMAN', 'SERVICE')),
    active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.role_definition (
    role_code text PRIMARY KEY,
    description text NOT NULL
);

CREATE TABLE atlas.role_binding (
    role_binding_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    role_code text NOT NULL REFERENCES atlas.role_definition(role_code),
    valid_from timestamptz NOT NULL DEFAULT transaction_timestamp(),
    valid_until timestamptz,
    CHECK (valid_until IS NULL OR valid_until > valid_from),
    UNIQUE (actor_id, role_code, valid_from)
);

CREATE TABLE atlas.change_request (
    change_id text PRIMARY KEY,
    title text NOT NULL,
    requester_actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    risk text NOT NULL CHECK (risk IN ('LOW','MODERATE','HIGH','CRITICAL')),
    status text NOT NULL CHECK (status IN ('DRAFT','PENDING_APPROVAL','APPROVED','REJECTED','IMPLEMENTING','VALIDATING','ACCEPTED','CLOSED')),
    required_approvals integer NOT NULL CHECK (required_approvals >= 1),
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.approval_action (
    approval_action_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    change_id text NOT NULL REFERENCES atlas.change_request(change_id),
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    decision text NOT NULL CHECK (decision IN ('APPROVE','REJECT','WITHDRAW')),
    reason text NOT NULL,
    recorded_at timestamptz NOT NULL DEFAULT transaction_timestamp(),
    UNIQUE (change_id, actor_id, decision)
);

CREATE TABLE atlas.evidence_bundle (
    evidence_id text PRIMARY KEY,
    device_id text NOT NULL,
    collection_type text NOT NULL,
    content_sha256 text NOT NULL CHECK (content_sha256 ~ '^[0-9a-f]{64}$'),
    classification text NOT NULL,
    storage_reference text NOT NULL,
    parser_version text,
    collected_at timestamptz NOT NULL,
    accepted_at timestamptz
);

CREATE TABLE atlas.telemetry_event (
    telemetry_event_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    metric_name text NOT NULL,
    metric_value text NOT NULL,
    labels jsonb NOT NULL DEFAULT '{}'::jsonb,
    observed_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.integration_outbox (
    outbox_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    destination_id text NOT NULL,
    contract_version text NOT NULL,
    payload jsonb NOT NULL,
    state text NOT NULL CHECK (state IN ('PENDING','DELIVERED','FAILED','DEAD_LETTER')),
    attempt_count integer NOT NULL DEFAULT 0,
    available_at timestamptz NOT NULL DEFAULT transaction_timestamp(),
    created_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE OR REPLACE FUNCTION atlas.enforce_requester_approval_independence()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    requester text;
BEGIN
    SELECT requester_actor_id INTO requester
    FROM atlas.change_request
    WHERE change_id = NEW.change_id
    FOR SHARE;

    IF NEW.decision = 'APPROVE' AND NEW.actor_id = requester THEN
        RAISE EXCEPTION 'requester cannot approve own change';
    END IF;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_approval_requester_independence
BEFORE INSERT ON atlas.approval_action
FOR EACH ROW
EXECUTE FUNCTION atlas.enforce_requester_approval_independence();

COMMIT;
