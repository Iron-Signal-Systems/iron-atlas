\set ON_ERROR_STOP on
\ir production-role-contract.sql

DO $bootstrap$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'atlas_test_runner') THEN
        CREATE ROLE atlas_test_runner LOGIN NOINHERIT NOSUPERUSER NOCREATEDB NOCREATEROLE NOREPLICATION NOBYPASSRLS;
    END IF;
END
$bootstrap$;
