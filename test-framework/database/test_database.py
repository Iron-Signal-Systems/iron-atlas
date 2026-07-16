#!/usr/bin/env python3
import os
import subprocess
import sys
import time

PSQL = os.environ["ATLAS_PSQL"]
BASE_ENV = os.environ.copy()


def sql(statement: str, user: str | None = None, actor: str | None = None, expect=True):
    env = BASE_ENV.copy()
    if user:
        env["PGUSER"] = user
    if actor:
        env["PGOPTIONS"] = f"-c atlas.actor_id={actor}"
    result = subprocess.run(
        [PSQL, "--no-psqlrc", "-X", "-v", "ON_ERROR_STOP=1", "-Atqc", statement],
        env=env,
        text=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    if expect and result.returncode != 0:
        raise AssertionError(f"SQL failed: {statement}\n{result.stderr}")
    if not expect and result.returncode == 0:
        raise AssertionError(f"SQL unexpectedly succeeded: {statement}")
    return result


def setup():
    sql("""
    SET ROLE atlas_schema_owner;

    INSERT INTO atlas.identity_provider(
        provider_id,
        provider_type,
        display_name,
        active
    ) VALUES
      ('dev', 'LOCAL_DEVELOPMENT', 'Disposable test provider', true),
      ('inactive-provider', 'OIDC', 'Inactive test provider', false)
    ON CONFLICT DO NOTHING;

    INSERT INTO atlas.actor(
        actor_id,
        display_name,
        actor_type,
        actor_status,
        disabled_at
    ) VALUES
      ('requester', 'Requester', 'HUMAN', 'ACTIVE', NULL),
      ('approver-a', 'Approver A', 'HUMAN', 'ACTIVE', NULL),
      ('approver-b', 'Approver B', 'HUMAN', 'ACTIVE', NULL),
      ('unauthorized', 'Unauthorized', 'HUMAN', 'ACTIVE', NULL),
      ('disabled-actor', 'Disabled Actor', 'HUMAN', 'DISABLED', transaction_timestamp()),
      ('no-role', 'No Role Actor', 'HUMAN', 'ACTIVE', NULL),
      ('expired-role', 'Expired Role Actor', 'HUMAN', 'ACTIVE', NULL),
      ('inactive-role', 'Inactive Role Actor', 'HUMAN', 'ACTIVE', NULL),
      ('unknown-role', 'Unknown Role Actor', 'HUMAN', 'ACTIVE', NULL)
    ON CONFLICT DO NOTHING;

    INSERT INTO atlas.external_identity(
        actor_id,
        provider_id,
        provider_subject
    ) VALUES
      ('requester', 'dev', 'subject-requester'),
      ('disabled-actor', 'dev', 'subject-disabled-actor'),
      ('no-role', 'dev', 'subject-no-role'),
      ('expired-role', 'dev', 'subject-expired-role'),
      ('inactive-role', 'dev', 'subject-inactive-role'),
      ('unknown-role', 'dev', 'subject-unknown-role'),
      ('unauthorized', 'inactive-provider', 'subject-inactive-provider')
    ON CONFLICT DO NOTHING;

    INSERT INTO atlas.role_definition(
        role_code,
        description,
        active
    ) VALUES
      ('INACTIVE_ROLE', 'Inactive resolver test role', false),
      ('CUSTOM_ROLE', 'Unsupported resolver test role', true)
    ON CONFLICT (role_code) DO NOTHING;

    INSERT INTO atlas.role_binding(
        actor_id,
        role_code,
        valid_from,
        valid_until,
        granted_by_actor_id,
        grant_reason
    ) VALUES
      ('requester', 'NETWORK_TECHNICIAN',
       transaction_timestamp() - interval '1 hour', NULL,
       NULL, 'test fixture'),
      ('approver-a', 'CHANGE_APPROVER',
       transaction_timestamp() - interval '1 hour', NULL,
       NULL, 'test fixture'),
      ('approver-b', 'CHANGE_APPROVER',
       transaction_timestamp() - interval '1 hour', NULL,
       NULL, 'test fixture'),
      ('expired-role', 'NETWORK_TECHNICIAN',
       transaction_timestamp() - interval '2 hours',
       transaction_timestamp() - interval '1 hour',
       NULL, 'expired resolver fixture'),
      ('inactive-role', 'INACTIVE_ROLE',
       transaction_timestamp() - interval '1 hour', NULL,
       NULL, 'inactive resolver fixture'),
      ('unknown-role', 'CUSTOM_ROLE',
       transaction_timestamp() - interval '1 hour', NULL,
       NULL, 'unsupported resolver fixture')
    ON CONFLICT DO NOTHING;
    """)

def assert_eq(actual, expected, name):
    if actual.strip() != expected:
        raise AssertionError(f"{name}: expected {expected!r}, got {actual!r}")
    print(f"PASS: {name}")


def main():
    setup()
    count = sql("SELECT count(*) FROM atlas.schema_migration;").stdout.strip()
    assert_eq(count, "7", "seven migrations recorded")

    resolved = sql(
        """
        SELECT actor_id || '|' || array_to_string(role_codes, ',')
        FROM atlas.resolve_governed_actor('dev', 'subject-requester');
        """,
        user="atlas_application",
    ).stdout.strip()
    assert_eq(
        resolved,
        "requester|NETWORK_TECHNICIAN",
        "governed actor resolver returns current Atlas roles",
    )

    for statement, name in (
        (
            "SELECT count(*) FROM atlas.resolve_governed_actor("
            "'dev', 'subject-disabled-actor');",
            "disabled actor resolution fails closed",
        ),
        (
            "SELECT count(*) FROM atlas.resolve_governed_actor("
            "'inactive-provider', 'subject-inactive-provider');",
            "inactive provider resolution fails closed",
        ),
        (
            "SELECT count(*) FROM atlas.resolve_governed_actor("
            "'dev', 'subject-unmapped');",
            "unmapped identity resolution fails closed",
        ),
        (
            "SELECT count(*) FROM atlas.resolve_governed_actor("
            "' dev', 'subject-requester');",
            "unnormalized provider resolution fails closed",
        ),
    ):
        assert_eq(
            sql(statement, user="atlas_application").stdout.strip(),
            "0",
            name,
        )

    for subject, actor_id, name in (
        (
            "subject-no-role",
            "no-role",
            "actor with no binding resolves with zero roles",
        ),
        (
            "subject-expired-role",
            "expired-role",
            "expired role binding is excluded",
        ),
        (
            "subject-inactive-role",
            "inactive-role",
            "inactive role definition is excluded",
        ),
    ):
        value = sql(
            "SELECT actor_id || '|' || cardinality(role_codes)::text "
            "FROM atlas.resolve_governed_actor("
            f"'dev', '{subject}');",
            user="atlas_application",
        ).stdout.strip()
        assert_eq(value, f"{actor_id}|0", name)

    sql(
        "SELECT actor_id FROM atlas.actor;",
        user="atlas_application",
        expect=False,
    )
    print("PASS: application cannot directly select governed actor table")

    sql(
        "SELECT atlas.current_actor_id();",
        user="atlas_application",
        actor="requester",
        expect=False,
    )
    print("PASS: internal actor helper is not directly executable by application role")

    sql("SELECT atlas.create_change_request('CHG-SELF', 'self approval test', 1);", user="atlas_application", actor="requester")
    requester_actor = sql(
        "SELECT requester_actor_id FROM atlas.change_request WHERE change_id='CHG-SELF';"
    ).stdout.strip()
    assert_eq(
        requester_actor,
        "requester",
        "governed change API resolves actor context unambiguously",
    )
    sql("SELECT atlas.record_approval('CHG-SELF', 'APPROVE', 'not allowed');", user="atlas_application", actor="requester", expect=False)
    print("PASS: requester self-approval rejected")

    sql("SELECT atlas.record_approval('CHG-SELF', 'APPROVE', 'independent approval');", user="atlas_application", actor="approver-a")
    status = sql("SELECT status FROM atlas.change_request WHERE change_id='CHG-SELF';").stdout.strip()
    assert_eq(status, "APPROVED", "independent approval changes status")

    sql("SELECT atlas.create_change_request('CHG-NOAUTH', 'authority test', 1);", user="atlas_application", actor="requester")
    sql("SELECT atlas.record_approval('CHG-NOAUTH', 'APPROVE', 'not authorized');", user="atlas_application", actor="unauthorized", expect=False)
    print("PASS: unauthorized approval rejected")

    # Direct application table writes are prohibited by grants.
    sql("INSERT INTO atlas.actor(actor_id, display_name, actor_type) VALUES ('bad', 'bad', 'HUMAN');", user="atlas_application", expect=False)
    print("PASS: application direct table write rejected")

    sql("SELECT atlas.create_change_request('CHG-TWO', 'two person concurrency', 2);", user="atlas_application", actor="requester")
    processes=[]
    for actor in ("approver-a", "approver-b"):
        env=BASE_ENV.copy(); env["PGUSER"]="atlas_application"; env["PGOPTIONS"]=f"-c atlas.actor_id={actor}"
        processes.append(subprocess.Popen(
            [PSQL, "--no-psqlrc", "-X", "-v", "ON_ERROR_STOP=1", "-Atqc",
             "SELECT atlas.record_approval('CHG-TWO','APPROVE','concurrent independent approval');"],
            env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
        ))
    results=[p.communicate() + (p.returncode,) for p in processes]
    if any(rc != 0 for _,_,rc in results):
        raise AssertionError(f"independent concurrent approvals failed: {results}")
    status = sql("SELECT status FROM atlas.change_request WHERE change_id='CHG-TWO';").stdout.strip()
    assert_eq(status, "APPROVED", "two independent concurrent approvals accepted")

    sql("SELECT atlas.create_change_request('CHG-DUP', 'duplicate concurrency', 2);", user="atlas_application", actor="requester")
    processes=[]
    for _ in range(2):
        env=BASE_ENV.copy(); env["PGUSER"]="atlas_application"; env["PGOPTIONS"]="-c atlas.actor_id=approver-a"
        processes.append(subprocess.Popen(
            [PSQL, "--no-psqlrc", "-X", "-v", "ON_ERROR_STOP=1", "-Atqc",
             "SELECT atlas.record_approval('CHG-DUP','APPROVE','duplicate race');"],
            env=env, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True,
        ))
    results=[p.communicate() + (p.returncode,) for p in processes]
    success=sum(1 for _,_,rc in results if rc == 0)
    failure=sum(1 for _,_,rc in results if rc != 0)
    if (success, failure) != (1,1):
        raise AssertionError(f"duplicate concurrency expected one success and one failure: {results}")
    state=sql("SELECT count(*) FROM atlas.change_approval_state WHERE change_id='CHG-DUP' AND current_decision='APPROVE';").stdout.strip()
    assert_eq(state, "1", "duplicate concurrent approval suppressed")

    sql("SET ROLE atlas_schema_owner; UPDATE atlas.approval_action SET reason='changed' WHERE change_id='CHG-SELF';", expect=False)
    print("PASS: append-only approval action update rejected")
    sql("SET ROLE atlas_schema_owner; DELETE FROM atlas.change_status_history WHERE change_id='CHG-SELF';", expect=False)
    print("PASS: append-only status history delete rejected")

    print("Database correctness result: PASS")


if __name__ == "__main__":
    try:
        main()
    except Exception as exc:
        print(f"FAIL: {exc}", file=sys.stderr)
        raise SystemExit(1)
