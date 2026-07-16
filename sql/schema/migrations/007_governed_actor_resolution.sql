BEGIN;
SET LOCAL lock_timeout = '5s';
SET LOCAL statement_timeout = '1min';
SET LOCAL idle_in_transaction_session_timeout = '1min';
SET LOCAL ROLE atlas_schema_owner;

CREATE OR REPLACE FUNCTION atlas.resolve_governed_actor(
    requested_provider_id text,
    requested_provider_subject text
)
RETURNS TABLE (
    actor_id text,
    role_codes text[]
)
LANGUAGE sql
STABLE
SECURITY DEFINER
SET search_path = pg_catalog, atlas
AS $$
    SELECT a.actor_id,
           COALESCE(
               array_agg(DISTINCT rb.role_code ORDER BY rb.role_code)
                   FILTER (WHERE rd.role_code IS NOT NULL),
               ARRAY[]::text[]
           ) AS role_codes
    FROM atlas.identity_provider AS ip
    JOIN atlas.external_identity AS ei
      ON ei.provider_id = ip.provider_id
    JOIN atlas.actor AS a
      ON a.actor_id = ei.actor_id
    LEFT JOIN atlas.role_binding AS rb
      ON rb.actor_id = a.actor_id
     AND rb.valid_from <= transaction_timestamp()
     AND (
         rb.valid_until IS NULL
         OR rb.valid_until > transaction_timestamp()
     )
    LEFT JOIN atlas.role_definition AS rd
      ON rd.role_code = rb.role_code
     AND rd.active
    WHERE requested_provider_id IS NOT NULL
      AND requested_provider_subject IS NOT NULL
      AND requested_provider_id = btrim(requested_provider_id)
      AND requested_provider_subject = btrim(requested_provider_subject)
      AND octet_length(requested_provider_id) BETWEEN 1 AND 256
      AND octet_length(requested_provider_subject) BETWEEN 1 AND 512
      AND ip.provider_id = requested_provider_id
      AND ip.active
      AND ei.provider_subject = requested_provider_subject
      AND a.actor_status = 'ACTIVE'
    GROUP BY a.actor_id;
$$;

ALTER FUNCTION atlas.resolve_governed_actor(text, text)
    OWNER TO atlas_schema_owner;

REVOKE ALL
ON FUNCTION atlas.resolve_governed_actor(text, text)
FROM PUBLIC;

COMMIT;
