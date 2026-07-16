# OIDC ID-Token Verification Implementation

## Status

Phase 1 Step 3 bounded implementation candidate.

This checkpoint implements provider discovery, JWKS-backed signature
verification, and normalized ID-token verification. It does not implement a
browser login route, authorization-code exchange, PKCE transaction storage,
cookies, durable sessions, CSRF, logout, trusted-proxy enforcement, or
production readiness.

## Accepted predecessor

The exact predecessor is the merged governed actor-resolution checkpoint:

`3ad3220c51179d3772d90da7f1025c4d41382922`

## Dependency decision

Iron Atlas uses:

- `github.com/coreos/go-oidc/v3` at `v3.19.0`; and
- its pinned OAuth 2.0 and JOSE dependency graph.

This avoids implementing JOSE signature validation, discovery processing, and
remote key-set refresh from scratch. Exact versions and checksums remain pinned
in `go.mod` and `go.sum`, and the candidate remains subject to `govulncheck`.

## Verification boundary

The verifier accepts only an explicitly configured provider profile:

- normalized Atlas provider ID;
- exact HTTPS issuer;
- bounded client ID;
- explicit asymmetric signing-algorithm allowlist;
- bounded HTTP timeout;
- bounded clock skew;
- bounded token age; and
- bounded token size.

Provider discovery must advertise:

- an exact matching issuer;
- HTTPS authorization, token, and JWKS endpoints;
- authorization-code response support;
- at least one subject type; and
- at least one permitted signing algorithm.

## ID-token checks

The verifier enforces:

- exactly three JWT segments;
- bounded protected header and payload;
- no duplicate security-sensitive JWT header or claim;
- no `none`, HMAC, unknown, or unconfigured signing algorithm;
- provider signature through the discovered JWKS;
- exact issuer;
- expected audience;
- required authorized party when multiple audiences exist;
- matching authorized party when present;
- expiry;
- issued-at presence, freshness, and future bound;
- not-before bound when present;
- nonce presence and constant-time equality;
- normalized stable subject;
- authentication-time future bound when present;
- access-token hash verification when `at_hash` is present; and
- normalized provider-neutral `authentication.Principal` output.

## Key rotation and outage behavior

An unknown key ID may trigger the dependency library's bounded remote-key-set
refresh. A valid replacement key is accepted after refresh. Network timeout or
key-provider unavailability is classified separately from an invalid
authentication assertion.

No outage manufactures a default principal or extends authority.

Because the pinned OIDC library flattens its outer signature-verification error,
Iron Atlas wraps the remote key set and records provider-dependency failure in
request-local verification state before that wrapping occurs. A completed key
refresh that cannot verify the signature remains an invalid assertion; network,
timeout, malformed-key-set, and non-success provider responses remain
dependency-unavailable failures. This classification never permits a principal.

The provider-emulator campaign separately proves that an unknown key while the
provider remains responsive is an invalid assertion, while network loss and a
non-success JWKS response are dependency-unavailable conditions. It also
directly verifies valid, missing, and mismatched access-token hashes,
authentication-time selection and bounds, and duplicate protected-header
rejection.

## Nonclaims

This checkpoint does not establish:

- an HTTP login or callback route;
- state or authorization-code replay storage;
- PKCE verifier persistence or code exchange;
- a production `authentication.Authenticator`;
- browser cookies or sessions;
- session fixation, rotation, expiry, logout, or revocation controls;
- CSRF defenses;
- trusted-proxy enforcement;
- representative-provider compatibility; or
- production readiness.
