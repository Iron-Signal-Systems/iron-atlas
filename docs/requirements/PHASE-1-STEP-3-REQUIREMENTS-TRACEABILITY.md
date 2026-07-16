# Phase 1 Step 3 Requirements Traceability

## Status

Phase 1 Step 3 contract candidate. No executable production authentication is
accepted by this record.

## Accepted predecessor

- Tag:
  `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`
- Accepted `dev` merge boundary:
  `1a750f7de791f567184c6f48e18eaec2933b8a14`
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
