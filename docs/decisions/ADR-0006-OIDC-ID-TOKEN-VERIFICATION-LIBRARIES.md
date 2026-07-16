# ADR-0006: OIDC ID-Token Verification Libraries

## Status

Accepted for the bounded Phase 1 Step 3 ID-token verification candidate.

## Context

OIDC discovery, JOSE parsing, asymmetric signature verification, audience and
issuer validation, and remote JWKS refresh are security-sensitive protocol
functions. Reimplementing those functions inside Iron Atlas would create a
larger custom cryptographic and protocol surface.

The repository already permits deliberately selected, pinned Go dependencies
where the dependency replaces high-risk custom infrastructure.

## Decision

Use `github.com/coreos/go-oidc/v3` at `v3.19.0`.

Its pinned dependency graph supplies the required OAuth 2.0 and JOSE
implementation. Exact module versions and checksums are retained in `go.mod`
and `go.sum`.

Iron Atlas adds its own stricter boundary around the library:

- exact HTTPS issuer and endpoint validation;
- explicit asymmetric algorithm allowlist;
- nonce verification;
- issued-at, not-before, authentication-time, and token-age bounds;
- authorized-party validation;
- duplicate sensitive field rejection;
- size and timeout limits;
- stable-subject normalization;
- access-token hash verification when present;
- outage classification; and
- deterministic provider-emulator, key-rotation, hostile, race, and concurrency
  tests.

## Rejected alternatives

### Custom JOSE and OIDC implementation

Rejected because it would create an unnecessary custom cryptographic and
protocol parser requiring substantially more independent assurance.

### Trusting provider claims without Atlas checks

Rejected because the dependency does not replace Atlas policy. In particular,
nonce verification and local protocol limits remain the caller's
responsibility.

### Provider-specific SDK

Rejected because Step 3 requires a provider-neutral adapter boundary and a
disposable emulator before representative-provider claims.

## Consequences

The repository gains a small, pinned security dependency graph and must retain
module verification, vulnerability scanning, update review, and hostile
protocol tests.

This ADR does not accept browser sessions, PKCE transaction persistence, CSRF,
trusted proxies, or production authentication.
