# PostgreSQL Bootstrap Contracts

Bootstrap scripts are run by an authorized PostgreSQL administrator, not by the Iron Atlas application.

1. `production-role-contract.sql` creates the production role topology without passwords.
2. Create the database owned by `atlas_database_owner`.
3. Grant database connect to the required login roles.
4. Run migrations as `atlas_migrator`.
5. Run `runtime-grants.sql` as an authorized administrator.

`development-roles.sql` is only for disposable tests and local development clusters.
