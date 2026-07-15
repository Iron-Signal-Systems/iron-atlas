# PostgreSQL Database Security Boundary

## Status

Phase 1 Step 1 implementation candidate.

## Security Objective

PostgreSQL must enforce the durable identity, authority, requester independence, approval, decision, and history invariants that cannot safely depend only on an HTML5 client or Go handler.

## Trusted Boundaries

- A future authenticated Go service maps an authenticated identity to an Iron Atlas actor.
- The service sets the transaction-local `atlas.actor_id` context before calling governed database functions.
- Caller-supplied JSON, form, query, or API fields never select the acting or approving identity.
- The database functions derive the actor from transaction context and compare it with durable actor and authority records.
- Direct database access is not an end-user interface.

Phase 1 Step 1 tests the database context contract. Production identity-provider integration remains outside this step.

## Governed Records

The initial durable boundary includes:

- Actors and external identities
- Role definitions and time-bounded role bindings
- Authority definitions and role-authority mappings
- Change requests and append-only status history
- Append-only approval actions and a mutable approval-state projection
- Decision records
- Audit events
- Migration history

## Two-Person Enforcement

Approval enforcement is performed inside PostgreSQL and applies across independent sessions:

- The requester cannot approve their own change.
- The acting approver must be active.
- The acting approver must hold `change.approve` authority through an active role binding.
- Approval functions derive the actor from `atlas.actor_id` transaction context.
- Duplicate active approvals by the same actor are rejected.
- Concurrent approvals serialize on the governed change row.
- The change reaches `APPROVED` only after the required number of distinct current approvals is met.
- Approval events remain append-only even when the current projection changes.

## Runtime Privilege Boundary

`atlas_application` receives:

- `CONNECT` to the Atlas database
- `USAGE` on the `atlas` schema
- `EXECUTE` on explicitly approved service functions
- `SELECT` on approved projections

It does not receive direct insert, update, delete, truncate, ownership, DDL, migration, role-management, database-creation, or bypass-RLS privileges.

## Append-Only Controls

Update and delete operations are rejected for:

- Schema migration history
- Change status history
- Approval actions
- Decision records
- Audit events

Mutable projections are not evidence. Their authoritative supporting events remain append-only.

## Administrative Reality

A PostgreSQL superuser can bypass ordinary database controls. Production deployment must therefore restrict superuser access, protect the host, retain off-host audit evidence, and define break-glass use. Those controls are staged for the production-security phase.

## Failure Behavior

Missing actor context, inactive actors, missing authority, self-approval, duplicate active approval, invalid state, changed migration checksums, and prohibited history mutation fail closed.
