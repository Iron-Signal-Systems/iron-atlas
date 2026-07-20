# Authenticated Server-Side Session Implementation

## Status

Phase 1 Step 3 bounded implementation candidate.

This checkpoint follows the merged HTTP login and callback boundary and starts
from the exact SSH-signed post-merge boundary
`6c912428a90b125f1b826729593e11ed914c12e9`.

It establishes durable PostgreSQL-backed session creation and authentication.
It does not establish session rotation, activity refresh, bounded session-count
or cleanup policy, logout, administrative revocation workflow, CSRF,
trusted-proxy enforcement,
production application wiring, local TOTP enrollment, representative-provider
compatibility, formal Step 3 acceptance, or production readiness.

## Purpose

A successful OIDC callback shall not become a durable browser identity by
placing provider tokens or actor authority in the browser. Atlas instead issues
one opaque cryptographically random session identifier, stores only its
nonreversible digest, and re-resolves current governed identity for protected
requests.

## Trust path

```text
verified OIDC principal
→ current governed actor resolution
→ controlled PostgreSQL session creation
→ opaque Secure HttpOnly host-only cookie
→ digest-only lookup on a later request
→ current governed actor and role re-resolution
→ immutable request identity context
```

The provider proves identity. Atlas resolves the governed actor and current
roles. The session only preserves a bounded authentication reference; it does
not grant provider-supplied authority.

## Browser identifier

The browser cookie is named:

```text
__Host-iron_atlas_session
```

It shall contain exactly 32 bytes of cryptographic randomness encoded with
unpadded base64url. The cookie shall have:

- `Secure`;
- `HttpOnly`;
- `Path=/`;
- no `Domain` attribute;
- `SameSite=Lax`; and
- absolute expiration bounded by the server-side record.

The identifier shall not appear in URLs, logs, errors, audit text, retained
validation evidence, or persistent plaintext.

## PostgreSQL record

Migration `008_authenticated_session.sql` creates
`atlas.authenticated_session`. The record binds:

- a SHA-256 digest of the browser identifier;
- provider ID and stable provider subject;
- the resolved governed actor ID;
- durable provider and actor references that preserve the historical session
  tuple without blocking a governed external-identity remap;
- creation, authentication, and last-activity timestamps;
- idle and absolute expiry;
- revocation and rotation-lineage placeholders for successor checkpoints;
- normalized authentication context and methods;
- normalized MFA state and MFA authentication time; and
- the applicable Atlas authentication security-policy version.

The application role has no direct table access. It may use only:

- `atlas.create_authenticated_session(...)`; and
- `atlas.authenticate_session(bytea)`.

Both functions are `SECURITY DEFINER` routines with a fixed search path and
explicitly revoked `PUBLIC` access.

## Session creation

The verified-principal handoff resolves the principal through the current
governed actor resolver before session creation. PostgreSQL independently
requires the same active provider, exact external-identity mapping, active
actor, bounded lifetime, normalized assurance values, and exact actor binding.

Failure to prove this state creates no session and no cookie.

A successful creation redirects only to a fixed configured local path. The path
uses a narrow ASCII allowlist and rejects schemes, hosts, queries, fragments,
percent-encoded separators, backslashes, and protocol-relative forms.

The first session checkpoint establishes fixed creation-time idle and absolute
bounds. Sliding activity refresh, bounded active-session counts, expired-record
cleanup, rotation, logout, and administrative revocation are deliberately
deferred to the successor lifecycle checkpoint.

## Protected-request authentication

A protected request must present exactly one canonical session cookie. Atlas:

1. rejects missing, duplicate, malformed, noncanonical, or oversized values;
2. decodes the identifier only in process memory;
3. hashes the raw identifier with SHA-256;
4. performs the controlled PostgreSQL lookup;
5. rejects unknown, expired, revoked, inactive, or remapped state;
6. returns the verified provider principal with the session-bound actor ID; and
7. requires the existing middleware to re-resolve the current actor and roles.

The immutable resolved identity validation rejects a current actor that differs
from the actor bound when the session was created. This prevents an external
identity remap from silently transferring an existing browser session to a
different governed actor. The session table intentionally does not use a
composite foreign key to the mutable external-identity mapping; a governed
remap can proceed, while the controlled lookup immediately stops returning the
old session.

## Authentication assurance and future MFA

The session record is designed now to preserve bounded provider-neutral
assurance metadata without enforcing an MFA policy in this checkpoint.

The planned successor assurance gate will define accepted OIDC `acr`, `amr`,
`auth_time`, MFA-age, role-sensitive, and step-up policies. Upstream-provider
MFA is the primary model. Phishing-resistant WebAuthn, FIDO2, passkeys, or
hardware security keys are preferred for high-impact authority.

RFC 6238 TOTP is planned as a compatible fallback for applications such as
Google Authenticator, 1Password, Aegis, FreeOTP, and similar clients. Atlas-local
TOTP remains a separate gate because it requires encrypted shared-secret
storage, enrollment proof, replay prevention, rate limiting, recovery codes,
reset governance, encryption-key rotation, and durable audit evidence.

No absent assurance value is interpreted as proof that MFA occurred.
Provider roles or groups do not become Atlas roles.

## Failure classification

- Missing, malformed, duplicate, unknown, expired, revoked, or remapped session
  state returns generic authentication-required behavior.
- PostgreSQL or session-store unavailability returns generic authentication
  service unavailable behavior.
- Identifier collisions fail closed and do not replace an existing session.
- No error contains the raw identifier, digest, provider token, authorization
  code, or unnecessary identity details.

## Validation

The bounded campaign includes:

- secure cookie and local-redirect behavior;
- cryptographic identifier generation;
- digest-only persistence requests;
- exact cookie cardinality and canonical encoding;
- missing, malformed, unknown, expired, and outage behavior;
- actor-remapping rejection;
- current governed actor re-resolution;
- PostgreSQL controlled-function access;
- denial of direct application table reads;
- concurrent lookup and Go race detection;
- manifest, migration, documentation, and static-validator checks; and
- predecessor HTTP login and callback revalidation.

## Nonclaims

This candidate does not establish:

- sliding idle refresh;
- bounded active-session count and expired-record cleanup;
- session rotation;
- logout;
- administrative revocation workflow;
- CSRF protection;
- trusted-proxy enforcement;
- production `atlasd` authentication wiring;
- authentication audit persistence;
- MFA enforcement or step-up authentication;
- Atlas-local TOTP enrollment or recovery;
- representative-provider compatibility;
- formal Phase 1 Step 3 acceptance; or
- production readiness.
