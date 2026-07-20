# Phase 1 Step 3 Requirements Traceability

## Status

Phase 1 Step 3 contract integrated. Authentication foundation, governed actor resolution, OIDC ID-token verification, authorization-code with PKCE, HTTP login and callback, and authenticated server-side session checkpoints are merged. Authentication assurance is the active bounded candidate; no local TOTP, completed session lifecycle, CSRF, trusted-proxy, or production authentication is accepted by this record.

## Accepted predecessor

- Tag:
  `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`
- Active implementation base:
  `6c912428a90b125f1b826729593e11ed914c12e9`
- Preserved invariant: authenticated actor context remains transaction-local in
  PostgreSQL and cannot leak across pooled connections.

## Traceability matrix

| Requirement | Contract focus | Planned enforcement | Planned evidence |
|---|---|---|---|
| `IA-AUTH-001` | UI is not authority | Policy and service checks | Handler and policy tests |
| `IA-AUTH-002` | Fail closed | Middleware and resolver | Negative and outage tests |
| `IA-AUTH-003` | Admin is not approval | Atlas role bindings | Separation tests |
| `IA-AUTH-004` | Trusted adapters only | Versioned adapter | Protocol tests |
| `IA-AUTH-005` | Unique provider-subject mapping | Governed resolver | Missing and ambiguity tests |
| `IA-AUTH-006` | Active provider and actor | Status checks | Disablement tests |
| `IA-AUTH-007` | Atlas roles authoritative | Role resolver | Claim escalation tests |
| `IA-AUTH-008` | Request data cannot select actor | Request context | Spoofing tests |
| `IA-AUTH-009` | Auth modes exclusive | Typed configuration | Startup tests |
| `IA-AUTH-010` | Bounded revocable sessions | Session service | Fixation and expiry tests |
| `IA-AUTH-011` | CSRF defense | CSRF middleware | Missing and replay tests |
| `IA-AUTH-012` | Explicit trusted proxy | Peer/header policy | Bypass tests |
| `IA-AUTH-013` | Secret redaction | Structured audit | Capture tests |
| `IA-AUTH-014` | Immutable and transaction-local actor | Context and PostgreSQL | Step 2 regression |
| `IA-AUTH-015` | Bounded lifecycle effects | Invalidation policy | Concurrent change tests |
| `IA-AUTH-016` | No unrestricted context | Boundary enforcement | Confused-deputy tests |
| `IA-AUTH-017` | Provider-neutral assurance and step-up | Assurance policy | Missing, ambiguous, downgrade, age, and step-up tests |
| `IA-AUTH-018` | Phishing-resistant MFA option with TOTP fallback | Provider and authenticator policy | WebAuthn/FIDO2 and RFC 6238 compatibility tests |
| `IA-AUTH-019` | Governed Atlas-local TOTP lifecycle | Encrypted authenticator service | Enrollment, replay, throttling, recovery, reset, and key-rotation tests |

## Existing governed schema

The accepted schema already contains:

- `atlas.actor`;
- `atlas.identity_provider`;
- `atlas.external_identity`;
- `atlas.role_definition`;
- `atlas.authority_definition`;
- `atlas.role_authority`; and
- `atlas.role_binding`.

Any Step 3 migration shall be limited to proven lifecycle, session,
authentication-event, or controlled service-API needs.

## Authentication foundation implementation status

- `IA-AUTH-008`: Implemented for the current HTTP boundary. Protected handlers consume only the actor from private server-side request context; body, form, query, path, and ordinary headers do not select it.
- `IA-AUTH-009`: Implemented for configuration. `development` and `production` are typed, mutually exclusive modes; the legacy boolean is rejected.
- `IA-AUTH-014`: Partially implemented. Request identity is private and immutable by copy, while the accepted transaction-local PostgreSQL actor boundary remains unchanged.
- `IA-AUTH-016`: Partially implemented. Nested identity middleware, unknown development roles, production development headers, and missing production adapters fail closed.

HTTP login and callback routes, production authenticator wiring, durable sessions, cookies, CSRF, trusted-proxy enforcement, lifecycle invalidation, authentication audit persistence, and representative-provider evidence remain required.


## HTTP login and callback implementation status

The active HTTP candidate partially implements `IA-AUTH-002`, `IA-AUTH-004`,
`IA-AUTH-008`, `IA-AUTH-009`, `IA-AUTH-013`, `IA-AUTH-015`, and
`IA-AUTH-016`.

It proves exact browser login and callback routes, a bounded secure browser
state-binding cookie, duplicate and unknown callback rejection, constant-time
state matching, callback issuer binding, provider-error cancellation, generic
failure responses, concurrent one-consumer behavior, and verified-principal
handoff.

It does not prove durable sessions, protected-request authentication, session
rotation, idle or absolute expiry, logout, administrative revocation, CSRF,
trusted proxies, production application wiring, authentication audit
persistence, representative-provider compatibility, or formal Step 3
acceptance.

## Authentication-assurance implementation status

The active candidate partially implements `IA-AUTH-002`, `IA-AUTH-005`, `IA-AUTH-006`, `IA-AUTH-008`, `IA-AUTH-013`, `IA-AUTH-015`, `IA-AUTH-016`, `IA-AUTH-017`, and `IA-AUTH-018`.

It proves bounded `acr`, `amr`, and `auth_time` normalization, duplicate and malformed assurance rejection, explicit versioned provider-MFA policy, role-sensitive phishing-resistant outcomes, maximum authentication age, generic additional-authentication responses, no session handoff for unsatisfied assurance, and exact policy-version and MFA enforcement in the server-side session boundary.

It does not prove Atlas-local TOTP, QR enrollment, recovery codes, WebAuthn, session rotation, logout, administrative revocation workflow, CSRF, trusted proxies, production wiring, representative-provider compatibility, formal Step 3 acceptance, or production readiness.

## Authenticated-session implementation status

The active candidate partially implements `IA-AUTH-002`, `IA-AUTH-005`,
`IA-AUTH-006`, `IA-AUTH-008`, `IA-AUTH-010`, `IA-AUTH-013`, `IA-AUTH-015`, and
`IA-AUTH-016`.

It proves opaque cryptographic browser identifiers, SHA-256 digest-only
persistence, controlled PostgreSQL creation and lookup, secure host-only cookie
attributes, exact cookie cardinality, fixed idle and absolute validity bounds,
current governed actor and role re-resolution, actor-remapping rejection,
least-privileged table isolation, assurance metadata retention, generic failure
classification, concurrent lookup, and database-outage behavior.

It does not prove sliding activity refresh, bounded session-count or cleanup
policy, session rotation, logout, administrative revocation workflow, CSRF,
trusted proxies, production wiring,
authentication audit persistence, MFA enforcement, local TOTP, representative-
provider compatibility, or formal Step 3 acceptance.

The planned successor checkpoints explicitly include authentication assurance
and MFA policy followed by TOTP enrollment, verification, and recovery. The
assurance gate governs provider-supplied `acr`, `amr`, `auth_time`, MFA age, and
step-up policy. The local TOTP gate governs encrypted secrets, compatible
applications, enrollment, replay prevention, throttling, recovery codes, reset,
and audit evidence.

## Required hostile classes

- invalid issuer, audience, signature, algorithm, key, redirect, state, nonce,
  PKCE, time, or replay state;
- inactive provider;
- missing, duplicate, remapped, or ambiguous external identity;
- disabled or retired actor;
- absent, expired, revoked, or conflicting role binding;
- actor or role selection through request-controlled data;
- development-header use in production mode;
- session fixation, rotation, expiry, logout, and revocation;
- CSRF missing, mismatch, expiry, and replay;
- trusted-proxy bypass and forwarded-header spoofing;
- provider outage, key rotation, and excessive clock skew;
- concurrent identity, actor, role, and session changes;
- secret and personal-data leakage; and
- continued PostgreSQL pooled-identity isolation.

## Final evidence requirements

The final acceptance record shall identify the adapter, provider emulator or
representative provider, toolchain, protocol profile, session policy, CSRF
policy, proxy topology, hostile classes, test counts, correctness result,
resource observations, limitations, predecessor revalidation, and exact
canonical clean-clone result.

## Governed actor-resolution implementation status

The PostgreSQL governed actor-resolution candidate partially implements
IA-AUTH-004, IA-AUTH-005, IA-AUTH-006, IA-AUTH-007, IA-AUTH-008,
IA-AUTH-009, IA-AUTH-014, and IA-AUTH-016.

It proves active-provider enforcement, exact provider-subject mapping, active
actor enforcement, current Atlas role loading, explicit role-code translation,
least-privileged database access, and fail-closed missing or unsupported state.
Provider protocol verification, sessions, CSRF, trusted proxies, and formal
Step 3 acceptance remain unimplemented.

## OIDC ID-token verification implementation status

The OIDC ID-token verification candidate partially implements `IA-AUTH-002`,
`IA-AUTH-004`, `IA-AUTH-008`, `IA-AUTH-009`, `IA-AUTH-013`, and
`IA-AUTH-016`.

It proves exact HTTPS discovery, JWKS-backed asymmetric signature verification,
issuer, audience, authorized party, expiry, issued-at, not-before, nonce,
stable-subject, access-token-hash, duplicate sensitive-field, key-rotation,
outage, malformed-input, race, and concurrency handling.

That verifier checkpoint did not implement authorization-code exchange or
PKCE transaction state. The successor candidate now implements bounded code
exchange and one-time in-memory transaction handling. Browser routes, cookies,
durable sessions, CSRF, logout, trusted proxies, authentication audit
persistence, representative-provider evidence, and formal Step 3 acceptance
remain unimplemented.

## OIDC authorization-code and PKCE implementation status

The authorization-code and PKCE candidate partially implements `IA-AUTH-002`,
`IA-AUTH-004`, `IA-AUTH-008`, `IA-AUTH-009`, `IA-AUTH-013`, `IA-AUTH-015`, and
`IA-AUTH-016`.

It proves cryptographic state, nonce, and PKCE verifier generation; SHA-256 state
digests; discovered `S256`; exact redirect and token endpoints; explicit client
authentication; bounded lifetime, capacity, code, and response sizes; atomic
one-time state consumption; invalid-code and outage classification; verified
principal production; replay rejection; exactly one concurrent consumer; and
secret-redacted errors.

It does not prove browser state cookies, HTTP callbacks, durable
restart-surviving transactions, authenticated sessions, CSRF, logout, trusted
proxies, production credential delivery, representative-provider compatibility,
or formal Step 3 acceptance.
