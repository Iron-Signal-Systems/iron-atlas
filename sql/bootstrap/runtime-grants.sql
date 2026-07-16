\set ON_ERROR_STOP on

REVOKE ALL ON SCHEMA atlas FROM PUBLIC;
REVOKE ALL ON ALL TABLES IN SCHEMA atlas FROM PUBLIC;
REVOKE ALL ON ALL SEQUENCES IN SCHEMA atlas FROM PUBLIC;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA atlas FROM PUBLIC;

GRANT USAGE ON SCHEMA atlas TO atlas_application, atlas_readonly, atlas_auditor;

GRANT EXECUTE ON FUNCTION atlas.create_change_request(text, text, integer) TO atlas_application;
GRANT EXECUTE ON FUNCTION atlas.record_approval(text, text, text) TO atlas_application;
GRANT EXECUTE ON FUNCTION atlas.resolve_governed_actor(text, text) TO atlas_application;

GRANT SELECT ON atlas.change_request, atlas.change_approval_summary TO atlas_application;
GRANT SELECT ON ALL TABLES IN SCHEMA atlas TO atlas_readonly;
GRANT SELECT ON atlas.schema_migration,
                atlas.change_status_history,
                atlas.approval_action,
                atlas.decision_record,
                atlas.audit_event
TO atlas_auditor;

ALTER DEFAULT PRIVILEGES FOR ROLE atlas_schema_owner IN SCHEMA atlas
    REVOKE ALL ON TABLES FROM PUBLIC;
ALTER DEFAULT PRIVILEGES FOR ROLE atlas_schema_owner IN SCHEMA atlas
    REVOKE ALL ON SEQUENCES FROM PUBLIC;
ALTER DEFAULT PRIVILEGES FOR ROLE atlas_schema_owner IN SCHEMA atlas
    REVOKE ALL ON FUNCTIONS FROM PUBLIC;
