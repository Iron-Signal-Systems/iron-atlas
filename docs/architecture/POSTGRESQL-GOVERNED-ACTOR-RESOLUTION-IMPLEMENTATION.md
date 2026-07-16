# PostgreSQL Governed Actor Resolution Implementation

## Status

Phase 1 Step 3 bounded implementation candidate.

This checkpoint implements the governed actor-resolution half of the accepted
trusted-authentication contract. It does not implement an external
authentication provider, login flow, browser session, CSRF control, trusted
proxy, or production readiness.

## Accepted predecessor

The exact predecessor is the merged authentication-foundation boundary:

`c6ad0d8d5c6268e5bd850eae646bd2e21ed7f3f5`

## Boundary

A verified `authentication.Principal` is resolved using only:

- its normalized provider ID;
- its normalized stable provider subject;
- an active governed identity provider;
- exactly one governed external-identity mapping;
- an active governed Atlas actor;
- active Atlas role definitions; and
- role bindings valid at database transaction time.

No provider group, provider administrative status, request header, query
parameter, body field, or default actor becomes Atlas authority.

## Least-privilege database interface

`atlas_application` receives no broad `SELECT` privilege on `atlas.actor`,
`atlas.identity_provider`, `atlas.external_identity`, `atlas.role_definition`,
or `atlas.role_binding`.

The application may execute only:

`atlas.resolve_governed_actor(text, text)`

The function is `SECURITY DEFINER`, fixes its `search_path`, rejects malformed
or non-normalized identifiers, and returns no row for missing, inactive,
disabled, retired, or unmapped identity state.

## Go resolver

`internal/authentication/postgresql.Resolver` implements
`authentication.ActorResolver`.

It:

- distinguishes missing governed state from database unavailability;
- explicitly maps governed database role codes to Go authorization roles;
- rejects unknown and duplicate role codes;
- enforces bounded role count;
- returns an actor with an empty role set when authentication is valid but no
  current role binding exists; and
- does not manufacture a default role.

An empty role set authenticates the actor but authorizes no protected
operation.

## Nonclaims

This implementation does not prove:

- OIDC, SAML, or Active Directory protocol verification;
- session creation, rotation, revocation, or expiry;
- CSRF or replay defenses;
- trusted-proxy enforcement;
- production secret delivery;
- operational availability; or
- production readiness.
