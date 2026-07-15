BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE TABLE atlas.role_definition (
    role_code text PRIMARY KEY,
    description text NOT NULL,
    active boolean NOT NULL DEFAULT true
);

CREATE TABLE atlas.authority_definition (
    authority_code text PRIMARY KEY,
    description text NOT NULL,
    governed boolean NOT NULL DEFAULT true
);

CREATE TABLE atlas.role_authority (
    role_code text NOT NULL REFERENCES atlas.role_definition(role_code),
    authority_code text NOT NULL REFERENCES atlas.authority_definition(authority_code),
    PRIMARY KEY (role_code, authority_code)
);

CREATE TABLE atlas.role_binding (
    role_binding_id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    actor_id text NOT NULL REFERENCES atlas.actor(actor_id),
    role_code text NOT NULL REFERENCES atlas.role_definition(role_code),
    valid_from timestamptz NOT NULL DEFAULT transaction_timestamp(),
    valid_until timestamptz,
    granted_by_actor_id text REFERENCES atlas.actor(actor_id),
    grant_reason text NOT NULL,
    CHECK (valid_until IS NULL OR valid_until > valid_from),
    UNIQUE (actor_id, role_code, valid_from)
);

INSERT INTO atlas.role_definition(role_code, description) VALUES
    ('NETWORK_TECHNICIAN', 'Network technician operational role'),
    ('NETWORK_ADMINISTRATOR', 'Network administrator role'),
    ('NETWORK_SECURITY', 'Network security review role'),
    ('CHANGE_APPROVER', 'Independent governed-change approver'),
    ('AUDITOR', 'Read-only governance and evidence auditor')
ON CONFLICT (role_code) DO NOTHING;

INSERT INTO atlas.authority_definition(authority_code, description) VALUES
    ('change.request', 'Create a governed change request'),
    ('change.approve', 'Approve or reject a governed change'),
    ('change.implement', 'Implement an approved change'),
    ('change.accept', 'Accept validated change results'),
    ('audit.read', 'Read governed history and audit evidence')
ON CONFLICT (authority_code) DO NOTHING;

INSERT INTO atlas.role_authority(role_code, authority_code) VALUES
    ('NETWORK_TECHNICIAN', 'change.request'),
    ('NETWORK_ADMINISTRATOR', 'change.request'),
    ('NETWORK_ADMINISTRATOR', 'change.implement'),
    ('NETWORK_SECURITY', 'change.approve'),
    ('CHANGE_APPROVER', 'change.approve'),
    ('AUDITOR', 'audit.read')
ON CONFLICT DO NOTHING;

ALTER TABLE atlas.role_definition OWNER TO atlas_schema_owner;
ALTER TABLE atlas.authority_definition OWNER TO atlas_schema_owner;
ALTER TABLE atlas.role_authority OWNER TO atlas_schema_owner;
ALTER TABLE atlas.role_binding OWNER TO atlas_schema_owner;

COMMIT;
