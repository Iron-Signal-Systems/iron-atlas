BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE TABLE atlas.approval_action (
    approval_action_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    change_id text NOT NULL REFERENCES atlas.change_request(change_id),
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    decision text NOT NULL CHECK (decision IN ('APPROVE', 'REJECT', 'WITHDRAW')),
    reason text NOT NULL CHECK (btrim(reason) <> ''),
    recorded_at timestamptz NOT NULL DEFAULT transaction_timestamp()
);

CREATE TABLE atlas.change_approval_state (
    change_id text NOT NULL REFERENCES atlas.change_request(change_id),
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    current_decision text NOT NULL CHECK (current_decision IN ('APPROVE', 'REJECT', 'WITHDRAW')),
    last_approval_action_id bigint NOT NULL UNIQUE REFERENCES atlas.approval_action(approval_action_id),
    updated_at timestamptz NOT NULL DEFAULT transaction_timestamp(),
    PRIMARY KEY (change_id, actor_id)
);

CREATE VIEW atlas.change_approval_summary AS
SELECT
    cr.change_id,
    cr.requester_actor_id,
    cr.status,
    cr.required_approvals,
    count(*) FILTER (WHERE cas.current_decision = 'APPROVE')::integer AS approval_count,
    count(*) FILTER (WHERE cas.current_decision = 'REJECT')::integer AS rejection_count
FROM atlas.change_request cr
LEFT JOIN atlas.change_approval_state cas ON cas.change_id = cr.change_id
GROUP BY cr.change_id, cr.requester_actor_id, cr.status, cr.required_approvals;

CREATE OR REPLACE FUNCTION atlas.record_approval(
    requested_change_id text,
    requested_decision text,
    requested_reason text
)
RETURNS bigint
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
DECLARE
    current_actor_id_value text := atlas.current_actor_id();
    change_row atlas.change_request%ROWTYPE;
    action_id bigint;
    approvals integer;
BEGIN
    SELECT cr.* INTO change_row
    FROM atlas.change_request AS cr
    WHERE cr.change_id = requested_change_id
    FOR UPDATE;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'change not found';
    END IF;
    IF change_row.status NOT IN ('PENDING_APPROVAL', 'APPROVED') THEN
        RAISE EXCEPTION 'change is not in an approvable state';
    END IF;
    IF current_actor_id_value = change_row.requester_actor_id
       AND requested_decision = 'APPROVE' THEN
        RAISE EXCEPTION 'requester cannot approve own change';
    END IF;
    IF NOT atlas.actor_has_authority(current_actor_id_value, 'change.approve') THEN
        RAISE EXCEPTION 'actor lacks change.approve authority';
    END IF;
    IF requested_decision = 'APPROVE' AND EXISTS (
        SELECT 1
        FROM atlas.change_approval_state AS cas
        WHERE cas.change_id = requested_change_id
          AND cas.actor_id = current_actor_id_value
          AND cas.current_decision = 'APPROVE'
    ) THEN
        RAISE EXCEPTION 'actor already has an active approval';
    END IF;

    INSERT INTO atlas.approval_action(change_id, actor_id, decision, reason)
    VALUES (
        requested_change_id, current_actor_id_value,
        requested_decision, requested_reason
    )
    RETURNING approval_action_id INTO action_id;

    INSERT INTO atlas.change_approval_state(
        change_id, actor_id, current_decision, last_approval_action_id, updated_at
    ) VALUES (
        requested_change_id, current_actor_id_value,
        requested_decision, action_id, transaction_timestamp()
    )
    ON CONFLICT (change_id, actor_id) DO UPDATE SET
        current_decision = EXCLUDED.current_decision,
        last_approval_action_id = EXCLUDED.last_approval_action_id,
        updated_at = EXCLUDED.updated_at;

    IF requested_decision = 'REJECT' THEN
        UPDATE atlas.change_request AS cr
        SET status = 'REJECTED', updated_at = transaction_timestamp()
        WHERE cr.change_id = requested_change_id;

        INSERT INTO atlas.change_status_history(
            change_id, previous_status, new_status, actor_id, reason
        ) VALUES (
            requested_change_id, change_row.status, 'REJECTED',
            current_actor_id_value, requested_reason
        );
    ELSIF requested_decision = 'APPROVE' THEN
        SELECT count(*)::integer INTO approvals
        FROM atlas.change_approval_state AS cas
        WHERE cas.change_id = requested_change_id
          AND cas.current_decision = 'APPROVE';

        IF approvals >= change_row.required_approvals
           AND change_row.status <> 'APPROVED' THEN
            UPDATE atlas.change_request AS cr
            SET status = 'APPROVED', updated_at = transaction_timestamp()
            WHERE cr.change_id = requested_change_id;

            INSERT INTO atlas.change_status_history(
                change_id, previous_status, new_status, actor_id, reason
            ) VALUES (
                requested_change_id, change_row.status, 'APPROVED',
                current_actor_id_value,
                'required independent approvals satisfied'
            );
        END IF;
    END IF;

    RETURN action_id;
END;
$$;

ALTER FUNCTION atlas.record_approval(text, text, text) OWNER TO atlas_schema_owner;
ALTER TABLE atlas.approval_action OWNER TO atlas_schema_owner;
ALTER TABLE atlas.change_approval_state OWNER TO atlas_schema_owner;
ALTER VIEW atlas.change_approval_summary OWNER TO atlas_schema_owner;

COMMIT;
