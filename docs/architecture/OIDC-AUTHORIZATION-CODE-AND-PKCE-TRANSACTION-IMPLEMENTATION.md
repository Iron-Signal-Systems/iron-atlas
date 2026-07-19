# OIDC Authorization-Code and PKCE Transaction Implementation

## Status

Phase 1 Step 3 bounded implementation candidate.

This checkpoint adds provider-neutral authorization-code exchange, PKCE S256,
short-lived one-time preauthentication transactions, and verified principal
production. It does not add HTTP login or callback routes, browser cookies,
durable restart-surviving sessions, CSRF enforcement, logout, trusted-proxy
enforcement, governed actor resolution wiring, or production readiness.

## Implementation predecessor

The exact repository base for this checkpoint is the merged `dev` commit:

`36394c917a7c60350f229fc80df2066a0c132681`

The formally accepted predecessor remains Phase 1 Step 2 under:

`phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`

The OIDC discovery, JWKS, and ID-token verification checkpoint remains a required
implementation predecessor and is revalidated by this candidate.

## Purpose

The authorization response must be bound to the exact browser initiation that
created it. Atlas shall not accept an authorization code merely because it came
from a configured token endpoint.

This checkpoint binds:

- one cryptographically random state value;
- one cryptographically random OIDC nonce;
- one cryptographically random PKCE verifier;
- the exact configured HTTPS redirect URI;
- one creation time;
- one short absolute expiry; and
- one atomic one-time consume operation.

## Begin contract

`AuthorizationCodeFlow.Begin` creates 256-bit random state, nonce, and PKCE
verifier values.

The returned authorization URL contains:

- `response_type=code`;
- the exact configured client ID;
- the exact configured HTTPS redirect URI;
- mandatory `openid` scope;
- the random state;
- the random nonce;
- an S256 PKCE challenge; and
- `code_challenge_method=S256`.

The client secret and PKCE verifier are never placed in the authorization URL.

The raw state is returned to the future browser boundary but is not retained by
the transaction store. The store retains only a SHA-256 state digest.

## Preauthentication transaction store

The candidate introduces a replaceable atomic store contract and a bounded
in-memory implementation.

The memory implementation:

- has an explicit maximum entry count;
- removes expired entries during create and consume operations;
- rejects duplicate state digests;
- consumes and deletes a transaction atomically;
- allows only one concurrent consumer to succeed;
- fails closed for missing, expired, replayed, malformed, or unavailable state;
- does not retain the raw state value; and
- loses every outstanding transaction on process restart.

Restart invalidation is intentional for this checkpoint. Durable,
restart-surviving preauthentication storage is not accepted because protecting a
persisted PKCE verifier requires the later credential and secret-protection
boundary.

## Complete contract

`AuthorizationCodeFlow.Complete`:

1. validates canonical 256-bit base64url state;
2. validates the bounded opaque authorization code;
3. atomically consumes the state transaction before contacting the provider;
4. uses the exact stored PKCE verifier;
5. uses the exact configured HTTPS redirect URI;
6. exchanges the code only with the discovered HTTPS token endpoint;
7. uses the explicitly compatible discovered client-authentication method;
8. prohibits token-endpoint redirects;
9. bounds the token response body;
10. requires a bearer access token and an ID token;
11. passes the stored nonce and returned access token to the existing verifier;
12. returns only a normalized `authentication.Principal`; and
13. never returns raw provider tokens, the code, state, nonce, verifier, or
    client secret.

Consumption before exchange means a transient provider failure requires a new
login attempt. This favors replay resistance over retrying an already presented
authorization code.

## Discovery requirements

The existing OIDC discovery boundary now retains:

- the exact authorization endpoint;
- the exact token endpoint;
- advertised token-endpoint authentication methods; and
- advertised PKCE challenge methods.

The flow requires advertised `S256` support.

A confidential-client configuration requires either
`client_secret_basic` or `client_secret_post`. A public-client configuration
requires advertised `none`. Unsupported or ambiguous client-authentication
configuration fails at construction.

## Provider and error behavior

The token exchange distinguishes:

- invalid request, authorization code, state, replay, or provider assertion;
- provider or network unavailability;
- token-endpoint server error;
- client-credential configuration failure;
- cancellation or timeout; and
- oversized provider response.

No error returned by the flow includes the authorization code, raw state,
nonce, PKCE verifier, access token, ID token, refresh token, or client secret.

## Concurrency and replay

The race campaign proves that concurrent completion attempts against the same
state allow exactly one token exchange. Every other consumer fails before
contacting the provider.

A completed, failed, or replayed callback does not recreate the transaction.

## Security boundaries preserved

- Exact HTTPS issuer, authorization, token, and JWKS endpoints remain required.
- ID-token signature, issuer, audience, authorized party, time, nonce,
  stable-subject, and access-token-hash checks remain enforced.
- Provider claims do not become Atlas roles.
- This component returns a provider-neutral principal, not an Atlas actor.
- PostgreSQL actor context remains transaction-local.
- Development identity headers remain prohibited in production mode.
- Raw authentication secrets remain prohibited from logs, responses, Git, and
  retained evidence.

## Validation campaign

The checkpoint requires:

- static contract validation;
- valid authorization URL and PKCE S256 assertions;
- exact redirect and client-authentication assertions;
- successful code exchange and principal verification;
- invalid-code and provider-outage classification;
- unknown, expired, malformed, and replayed state rejection;
- concurrent single-consumer proof;
- bounded store capacity and expired-state cleanup;
- oversized token-response rejection;
- randomness-failure handling;
- secret-redaction assertions;
- Go race testing;
- Go vet;
- module verification;
- vulnerability analysis;
- full repository regression;
- predecessor phase-gate revalidation; and
- complete repository validation.

## Nonclaims

This checkpoint does not establish:

- HTTP `/login` or callback handlers;
- browser state binding through a cookie;
- a production `authentication.Authenticator`;
- governed actor-resolution wiring into the application;
- durable restart-surviving preauthentication storage;
- server-side authenticated sessions;
- session fixation, rotation, idle expiry, absolute expiry, logout, or
  revocation;
- CSRF enforcement;
- trusted-proxy or host-header enforcement;
- production credential delivery or secret rotation;
- representative-provider compatibility;
- formal Phase 1 Step 3 acceptance; or
- production readiness.

## HTTP successor checkpoint

A separate bounded HTTP login and callback component now consumes this
authorization-code flow.

The authorization-code component itself remains cookie-free and route-free. The
HTTP successor owns browser state-cookie binding, callback cardinality, issuer
binding, generic browser responses, and verified-principal handoff.

Neither component creates a durable authenticated Atlas session at this
checkpoint.
