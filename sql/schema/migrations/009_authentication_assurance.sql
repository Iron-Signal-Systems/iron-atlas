BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

UPDATE atlas.authenticated_session
   SET revoked_at = transaction_timestamp(),
       revocation_reason = 'MFA assurance policy introduced'
 WHERE NOT mfa_authenticated
   AND revoked_at IS NULL;

ALTER TABLE atlas.authenticated_session
    ADD CONSTRAINT authenticated_session_requires_mfa
    CHECK (
        revoked_at IS NOT NULL
        OR (mfa_authenticated AND mfa_authenticated_at IS NOT NULL)
    );

COMMIT;
