# OIDC HTTP Login and Callback Implementation

## Status

Phase 1 Step 3 bounded implementation candidate.

The OIDC discovery, ID-token verification, authorization-code exchange, state,
nonce, and PKCE transaction components are merged predecessors. This checkpoint
adds the browser-facing HTTP login and callback boundary and a verified-principal
handoff seam.

This checkpoint is not a durable authenticated-session implementation and is
not formal Phase 1 Step 3 acceptance.

## Accepted working base

- Canonical branch: `dev`
- Required base commit:
  `28ec1eab5b5c4e69731e9b0a79fe6105beab316d`
- Preserved protocol predecessor:
  `validate_phase1_step3_oidc_authorization_code_pkce.sh`

## Purpose

The HTTP boundary binds one browser initiation to one OIDC authorization-code
callback without allowing request-controlled actor selection, redirect
selection, callback replay, provider-detail reflection, or state-cookie
confusion.

The boundary produces only a verified `authentication.Principal`. A later
bounded session package will consume that principal and create a server-side
session.

## Routes

The checkpoint defines exactly two browser routes:

```text
GET /auth/login
GET /auth/callback
```

No state-changing behavior is accepted on another HTTP method.

The routes are implemented as a standalone handler and are deliberately not
wired into `atlasd` or the protected application middleware by this checkpoint.

## Login boundary

`GET /auth/login`:

1. rejects request query parameters;
2. begins one bounded OIDC authorization transaction;
3. receives the exact authorization URL, state, and expiry from the accepted
   authorization-code flow;
4. validates that the destination is an exact HTTPS URL;
5. sets one short-lived browser-binding cookie;
6. returns an empty-body redirect to the exact generated authorization URL; and
7. emits no-cache and no-referrer response headers.

The route never accepts a request-supplied return location.

## State-binding cookie

The cookie name is:

```text
__Host-iron_atlas_oidc_state
```

The cookie is:

- `Secure`;
- `HttpOnly`;
- `SameSite=Lax`;
- `Path=/`;
- host-only, with no `Domain`;
- bounded to the preauthentication transaction lifetime; and
- replaced by each new login initiation.

The callback clears the cookie on every success and failure path.

The raw state exists only in the browser redirect and short-lived browser
cookie. The server-side preauthentication store retains only the SHA-256 state
digest.

## Callback boundary

`GET /auth/callback`:

1. bounds the complete query string;
2. parses the query without ignoring parse errors;
3. rejects unknown parameters;
4. requires exactly one `state`;
5. requires exactly one of `code` or `error`;
6. rejects duplicate security-sensitive fields;
7. requires exactly one state cookie;
8. validates canonical 256-bit state encoding;
9. compares callback state to cookie state in constant time;
10. verifies the optional callback issuer against the exact configured issuer;
11. atomically consumes error callbacks so they cannot be replayed;
12. completes code exchange and ID-token verification through the accepted
    authorization-code flow; and
13. passes only the verified normalized principal to the injected
    `VerifiedPrincipalHandler`.

Provider `error_description`, `error_uri`, authorization codes, state values,
cookies, tokens, and internal dependency errors are not reflected to the
browser.

## Supported callback fields

The bounded parser recognizes only:

- `state`;
- `code`;
- `iss`;
- `session_state`;
- `error`;
- `error_description`; and
- `error_uri`.

Unknown fields, including request-controlled redirect fields, fail closed.

## Provider-error cancellation

The authorization-code flow adds a bounded `Cancel` operation that atomically
consumes a valid preauthentication transaction without performing token
exchange.

This is used for provider-declared callback failures and issuer mismatch.
Repeated cancellation or later completion of the same state fails closed.

## Verified-principal handoff seam

The callback does not construct an Atlas actor and does not create a session.

It invokes:

```go
type VerifiedPrincipalHandler interface {
    ServeVerifiedPrincipal(
        http.ResponseWriter,
        *http.Request,
        authentication.Principal,
    )
}
```

The injected handler is trusted internal code. The next bounded package will
implement this seam using a server-side session service. Until then, the HTTP
handler remains unwired from the production application.

## Error classification

Browser responses are deliberately small and generic:

- invalid, malformed, mismatched, replayed, or provider-denied callbacks:
  `401`;
- unavailable provider, token endpoint, randomness source, or transaction
  store: `503`;
- unsupported method: `405`;
- invalid login request: `400`.

Responses do not enumerate accounts, identities, providers, or internal
failure details.

## Hostile and concurrency campaign

The checkpoint tests:

- secure state-cookie attributes;
- empty redirect bodies;
- missing and duplicate callback values;
- duplicate cookies;
- malformed and oversized callbacks;
- unknown callback parameters;
- state mismatch before exchange;
- issuer mismatch;
- provider-declared errors;
- provider outage;
- token-flow failure classification;
- state and code redaction;
- unsafe methods;
- insecure issuer configuration;
- exact one-consumer behavior under concurrent callbacks; and
- exact one-consume behavior for provider-error cancellation.

## Resource boundary

Correctness tests remain separate from resource observations.

The implementation bounds:

- authorization URL length;
- callback query length;
- individual callback values;
- provider error text;
- cookie lifetime; and
- callback parameter cardinality.

Performance thresholds remain `NOT_EVALUATED`.

## Explicit nonclaims

This checkpoint does not establish:

- durable or restart-surviving preauthentication state;
- server-side authenticated sessions;
- session cookies;
- protected-route authentication;
- session rotation;
- idle or absolute session expiry;
- logout;
- administrative revocation;
- CSRF protection;
- trusted-proxy enforcement;
- external-origin reconstruction;
- production application wiring;
- authentication audit persistence;
- production credential delivery;
- representative-provider compatibility;
- formal Phase 1 Step 3 acceptance; or
- production readiness.
