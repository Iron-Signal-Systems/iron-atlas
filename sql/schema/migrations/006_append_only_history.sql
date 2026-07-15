BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE TABLE atlas.decision_record (
    decision_record_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    decision_code text NOT NULL,
    subject_type text NOT NULL,
    subject_id text NOT NULL,
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    decision text NOT NULL,
    rationale text NOT NULL,
    recorded_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.audit_event (
    audit_event_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    event_type text NOT NULL,
    actor_id text REFERENCES atlas.actor(actor_id),
    subject_type text NOT NULL,
    subject_id text NOT NULL,
    event_data jsonb NOT NULL DEFAULT '{}'::jsonb,
    recorded_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE OR REPLACE FUNCTION atlas.reject_append_only_mutation()
RETURNS trigger
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
BEGIN
    RAISE EXCEPTION '% is append-only; % is prohibited', TG_TABLE_NAME, TG_OP;
END;
$$;

CREATE TRIGGER schema_migration_append_only
BEFORE UPDATE OR DELETE ON atlas.schema_migration
FOR EACH ROW EXECUTE FUNCTION atlas.reject_append_only_mutation();

CREATE TRIGGER change_status_history_append_only
BEFORE UPDATE OR DELETE ON atlas.change_status_history
FOR EACH ROW EXECUTE FUNCTION atlas.reject_append_only_mutation();

CREATE TRIGGER approval_action_append_only
BEFORE UPDATE OR DELETE ON atlas.approval_action
FOR EACH ROW EXECUTE FUNCTION atlas.reject_append_only_mutation();

CREATE TRIGGER decision_record_append_only
BEFORE UPDATE OR DELETE ON atlas.decision_record
FOR EACH ROW EXECUTE FUNCTION atlas.reject_append_only_mutation();

CREATE TRIGGER audit_event_append_only
BEFORE UPDATE OR DELETE ON atlas.audit_event
FOR EACH ROW EXECUTE FUNCTION atlas.reject_append_only_mutation();

ALTER FUNCTION atlas.reject_append_only_mutation() OWNER TO atlas_schema_owner;
ALTER TABLE atlas.decision_record OWNER TO atlas_schema_owner;
ALTER TABLE atlas.audit_event OWNER TO atlas_schema_owner;

COMMIT;
