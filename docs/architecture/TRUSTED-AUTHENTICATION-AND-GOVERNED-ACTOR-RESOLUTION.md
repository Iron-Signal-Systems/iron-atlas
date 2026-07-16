# Trusted Authentication and Governed Actor Resolution

## Status

Phase 1 Step 3 normative contract integrated; authentication-foundation implementation candidate active.

This document defines the production-authentication and governed
actor-resolution boundary. It does not claim that a production identity
provider, login flow, session service, trusted-proxy deployment, or production
authentication implementation is accepted.

The accepted predecessor is Phase 1 Step 2. The current accepted `dev` merge
boundary is `1a750f7de791f567184c6f48e18eaec2933b8a14`, and the immutable
predecessor tag is
`phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`.

## Purpose

Iron Atlas shall derive every production request identity from a trusted,
verified authentication result and then resolve that identity through
Atlas-governed records.

Authentication proves that a provider identity was verified. Actor resolution
determines which governed Atlas actor that identity represents. Authorization
determines what that actor may do. These decisions shall remain distinct.

No authentication adapter, provider claim, proxy, session, service,
administrator, or accumulated privilege may create an unrestricted execution
context or propagate authority across trust boundaries.

## Existing accepted foundation

Phase 1 Steps 1 and 2 already provide:

- governed `atlas.actor` records with active, disabled, and retired states;
- governed `atlas.identity_provider` records;
- unique `(provider_id, provider_subject)` external identities;
- Atlas-owned role definitions, authorities, and time-bounded role bindings;
- a least-privileged PostgreSQL runtime;
- transaction-local actor context;
- pooled-connection identity isolation; and
- database-enforced governed-change controls.

Step 3 extends this foundation. It shall not replace or weaken it.


## Authentication foundation implementation checkpoint

The first bounded implementation candidate now provides:

- typed `development` and `production` authentication modes;
- a dedicated middleware that establishes identity before protected handlers;
- private immutable request context for the normalized principal and resolved actor;
- a controlled development-only header adapter with bounded parsing;
- unconditional rejection of development identity headers in production mode;
- future `Authenticator` and `ActorResolver` interfaces;
- fail-closed protected routes when no production adapter is configured; and
- public health, readiness, and static-asset paths that do not manufacture an actor.

No external provider adapter is accepted. No provider-backed actor lookup, session, CSRF, replay, logout, trusted-proxy, or production deployment boundary is implemented by this checkpoint.

## Terminology

- **Authentication adapter:** A versioned component that verifies one supported
  authentication protocol and produces a normalized verified identity.
- **Provider subject:** The provider-stable identifier for an authenticated
  identity. A display name or mutable email address is not the durable subject.
- **External identity:** The governed provider-and-subject mapping to one Atlas
  actor.
- **Resolved actor:** One active Atlas actor plus current Atlas role bindings,
  placed into immutable server-side request context.
- **Session:** Bounded, revocable server-side authentication state referenced
  by an opaque browser cookie.

## Required trust chain

```text
untrusted request
    ↓
trusted transport and explicitly configured proxy boundary
    ↓
authentication adapter verifies protocol evidence
    ↓
normalized provider ID and stable provider subject
    ↓
active governed provider lookup
    ↓
unique governed external-identity lookup
    ↓
active governed Atlas actor lookup
    ↓
current Atlas role-binding lookup
    ↓
immutable server-side request context
    ↓
Atlas authorization policy
    ↓
transaction-local PostgreSQL actor context for governed mutations
```

Skipping, combining, or reversing these decisions is prohibited.

## Authentication adapter boundary

The first executable adapter may use OpenID Connect Authorization Code flow
with PKCE. The architecture shall remain provider-neutral behind a versioned
adapter interface.

A production adapter shall verify all applicable protocol properties:

- trusted issuer;
- intended audience and authorized party;
- signature and permitted algorithm;
- token and authorization-code expiry;
- not-before and issued-at bounds;
- state, nonce, and PKCE binding;
- redirect URI;
- provider key identity and bounded key rotation;
- maximum token, claim, header, and callback sizes;
- required stable subject;
- authentication time and assurance metadata when policy requires it; and
- replay resistance.

Atlas shall not collect external-directory passwords through this boundary.
Simple LDAP bind and request-supplied identity headers are not production
authentication mechanisms.

## Governed actor resolution

A verified provider result is not an Atlas actor until:

1. The provider exists and is active.
2. A stable provider subject is present.
3. Exactly one external identity matches provider plus subject.
4. The mapped actor exists.
5. The actor is `ACTIVE`.
6. Current Atlas role bindings are loaded.
7. No conflicting or ambiguous state exists.

Missing, unmapped, inactive, duplicated, expired, malformed, or ambiguous state
shall fail closed.

Provider groups, roles, administrative status, and other claims shall not
directly become Atlas authority. Atlas role bindings remain authoritative.

No default production actor or default production role is permitted.

## Request identity boundary

The resolved actor shall be placed into immutable server-side request context.
Handlers and services may read it but shall not replace it.

The following shall never select the production actor or roles:

- request bodies;
- forms;
- query parameters;
- path parameters;
- ordinary request headers;
- browser local storage; or
- database session state left by an earlier request.

The existing `X-Iron-Atlas-Actor` and `X-Iron-Atlas-Roles` headers are
development-only. Production mode shall ignore and reject attempts to use them
for identity selection.

Development identity and production authentication modes shall be mutually
exclusive. Production mode shall not silently fall back to development mode.

## Session boundary

Browser sessions shall use an opaque, cryptographically random identifier in a
cookie. Provider access, ID, and refresh tokens shall not be stored in browser
local storage.

A server-side session record shall bind at minimum:

- a nonreversible digest of the session identifier;
- provider ID and provider subject;
- resolved actor ID;
- creation, authentication, and last-activity times;
- idle and absolute expiry;
- revocation state and reason;
- rotation lineage;
- required authentication assurance metadata; and
- the applicable security-policy version.

Sessions shall have:

- `Secure`, `HttpOnly`, and an accepted `SameSite` policy;
- bounded idle and absolute lifetimes;
- rotation after login and material authentication-state changes;
- logout and administrative revocation;
- bounded invalidation after provider, identity, actor, or role changes;
- no identifier in URLs, logs, or persistent plaintext;
- bounded count and cleanup behavior; and
- constant-time secret comparisons where applicable.

## CSRF and browser boundary

State-changing browser requests shall require an accepted CSRF defense in
addition to authentication. Cookie policy alone is not sufficient.

The design shall include:

- origin or trusted-site validation;
- a cryptographically bound CSRF value;
- no state changes through safe HTTP methods;
- explicit content-type handling;
- bounded request bodies; and
- rejection of missing, malformed, expired, replayed, or mismatched CSRF state.

## Trusted proxy and transport boundary

The direct peer is untrusted unless it is inside an explicitly configured
trusted-proxy boundary.

Production deployment shall define:

- which process terminates client TLS;
- which proxy addresses or Unix-domain sockets are trusted;
- which forwarded headers are accepted;
- how forwarded identity headers are prohibited or protected;
- how original scheme and host are validated;
- how redirect URIs avoid host-header injection; and
- how direct access that bypasses the trusted proxy is denied.

Ordinary proxy headers shall not carry an Atlas actor or Atlas roles.

## Authorization and PostgreSQL context

Authentication does not grant authority.

Atlas authorization shall evaluate current Atlas role bindings. Platform or
identity-provider administration shall not automatically grant change approval.

For governed PostgreSQL mutations, the resolved actor ID shall continue to be
bound only within the database transaction. Session-scoped database actor
identity remains prohibited.

## Lifecycle and outage behavior

The implementation shall define fail-closed behavior for:

- provider disablement;
- actor disablement or retirement;
- external-identity remapping;
- role-binding grant, expiry, or revocation;
- signing-key rotation;
- session revocation;
- authentication-policy changes;
- excessive clock skew; and
- provider, database, or session-store unavailability.

Caching is permitted only with explicit maximum age, invalidation, and stale
state disposition. An outage shall not create a default identity or extend
authority without an accepted policy.

## Audit and privacy

Authentication and session events shall be auditable without recording
passwords, authorization codes, raw tokens, cookie values, CSRF secrets, client
secrets, private keys, or unnecessary personal data.

The audit model shall distinguish successful login, bounded failure categories,
logout, expiry, revocation, actor-resolution failure, ambiguous identity,
role-resolution failure, replay rejection, and policy or key-version changes.

Failure responses shall avoid identity and account enumeration.

## Required implementation seams

Step 3 implementation shall introduce replaceable boundaries for:

- protocol verification;
- provider configuration;
- external-identity and actor resolution;
- role-binding resolution;
- server-side session persistence;
- immutable request-context injection;
- CSRF verification;
- authentication audit events; and
- readiness observations.

HTTP handlers shall depend on the normalized resolved-actor contract rather
than provider-specific claims.

## Explicit exclusions

Step 3 does not accept:

- collector or device credential delivery and rotation;
- PostgreSQL TLS certificate deployment;
- backup and restoration;
- high availability;
- evidence intake or protected evidence storage;
- live infrastructure collection;
- automated remediation;
- production performance budgets; or
- production readiness.

## Acceptance boundary

Step 3 may be accepted only after the exact pushed candidate passes:

- unit, integration, race, and concurrency tests;
- spoofed, malformed, replayed, expired, ambiguous, and disabled-state tests;
- session fixation, rotation, expiry, logout, and revocation tests;
- CSRF and trusted-proxy tests;
- secret-redaction checks;
- continued PostgreSQL actor-isolation tests;
- complete repository validation;
- isolated revalidation of accepted Step 2; and
- canonical clean-clone validation with declared and verified tooling.

A disposable provider emulator may prove the generic protocol boundary.
Representative-provider validation remains a separately identified environment
claim.

## PostgreSQL governed actor-resolution implementation checkpoint

The active bounded implementation candidate adds
`atlas.resolve_governed_actor(text, text)` and a Go PostgreSQL
`authentication.ActorResolver`.

The application role receives function execution only and does not receive
broad `SELECT` access to governed identity or role tables. Resolution requires
an active provider, one external identity mapping, an active actor, active role
definitions, and role bindings valid at database transaction time. Unknown
database role codes fail closed in Go.

This checkpoint does not implement a production authentication adapter,
sessions, CSRF, trusted-proxy enforcement, or production readiness.

## OIDC ID-token verification implementation checkpoint

The active bounded implementation candidate adds provider-neutral OIDC
discovery, remote JWKS signature verification, exact issuer and audience
validation, authorized-party enforcement, explicit asymmetric algorithm
allowlisting, expiry, issued-at, not-before, nonce, stable-subject,
access-token-hash, duplicate sensitive-field, key-rotation, outage, race, and
concurrency enforcement.

This checkpoint deliberately stops before authorization-code exchange and
preauthentication transaction persistence because those controls must be
designed together with the later session, replay, cookie, CSRF, logout, and
trusted-proxy boundaries. It does not yet implement a production
`authentication.Authenticator`.
