# Trusted Authentication and Governed Actor Resolution Testing

## Status

Phase 1 Step 3 test contract integrated. Authentication foundation, governed resolver, OIDC ID-token, authorization-code with PKCE, and HTTP login and callback campaigns are merged; authenticated-session tests are active.

These tests are required for the final executable candidate. This document does
not claim that a production adapter or session implementation exists.

## Authentication foundation implementation checkpoint

The merged foundation checkpoint adds executable tests for typed mode parsing, default mode selection, legacy-setting rejection, bounded development headers, duplicate and unknown role rejection, immutable context copies, production rejection of development headers, missing-adapter fail-closed behavior, adapter/resolver composition, nested middleware rejection, public health/readiness paths, and query-string actor spoofing.

These tests do not substitute for later provider-protocol, actor-resolution, authorization-code, session, CSRF, replay, logout, key-rotation, or trusted-proxy campaigns.


## HTTP login and callback implementation campaign

The active bounded campaign covers:

- exact login and callback methods and paths;
- rejection of login query parameters;
- empty-body provider redirects;
- secure host-only state-cookie attributes;
- callback query and value bounds;
- exact callback parameter cardinality;
- duplicate cookie rejection;
- unknown parameter rejection;
- state mismatch before exchange;
- constant-time state comparison;
- exact optional issuer binding;
- provider-error cancellation;
- provider outage classification;
- state, code, cookie, provider-message, and internal-error redaction;
- exactly one concurrent verified-principal handoff; and
- exact one-time cancellation.

The campaign does not cover durable sessions, session cookies, protected
requests, rotation, expiry, logout, revocation, CSRF, trusted proxies,
production application wiring, audit persistence, or representative-provider
compatibility.

## Authenticated-session implementation campaign

The active bounded campaign covers:

- secure host-only session-cookie attributes;
- narrow fixed local-path redirect validation;
- 32-byte cryptographic identifier generation and canonical base64url encoding;
- SHA-256 digest-only persistence requests;
- controlled PostgreSQL creation and lookup routines;
- denial of direct application session-table access;
- exact cookie cardinality and malformed-value rejection;
- missing, unknown, expired, actor-mismatched, and outage behavior;
- current governed actor and role re-resolution;
- external-identity actor-remapping rejection;
- normalized assurance metadata and policy-version retention;
- concurrent lookups under the race detector; and
- predecessor HTTP login and callback revalidation.

The campaign does not cover sliding activity refresh, bounded session-count
or cleanup policy, rotation, logout, administrative revocation workflow, CSRF,
trusted proxies, production wiring,
MFA enforcement, local TOTP, representative providers, or formal acceptance.

## Authentication assurance and MFA successor campaigns

The planned assurance campaign will test accepted and rejected `acr`, `amr`,
`auth_time`, MFA-age, step-up, role-sensitive, missing, ambiguous, downgraded,
and provider-outage state. No absent claim may be interpreted as successful
MFA.

The planned local TOTP campaign will test Google Authenticator, 1Password,
Aegis, FreeOTP, and equivalent RFC 6238 clients through standard provisioning;
encrypted secret storage; enrollment confirmation; time-window and replay
behavior; rate limiting; concurrent verification; recovery-code one-time use;
reset separation of duties; actor disablement; secret redaction; backup and
restoration; and encryption-key rotation.

## Test layers

### Pure contract tests

Test normalized provider identities, actor resolution, Atlas role selection,
session policy, CSRF policy, and error classification without HTTP where
possible.

### HTTP boundary tests

Exercise callbacks, cookies, authenticated requests, logout, mutations,
security headers, trusted-proxy handling, and rejection of development identity
headers in production mode.

### PostgreSQL integration tests

Use disposable PostgreSQL to prove unique provider-subject mapping, provider
and actor status, current role bindings, session revocation when persisted, and
continued transaction-local actor isolation.

### Protocol-emulator tests

Use a disposable provider emulator or deterministic signed fixtures to test
issuer, audience, signature, key rotation, expiry, state, nonce, PKCE, redirect,
and replay behavior without a production dependency.

## Required positive cases

- one verified provider subject resolves to exactly one active actor;
- current Atlas role bindings load deterministically;
- authenticated reads require Atlas authority;
- governed mutations use the same actor in policy and PostgreSQL context;
- login rotates preauthentication state;
- logout revokes the session;
- idle and absolute expiry fail closed;
- role-binding expiry removes authority;
- bounded provider key rotation succeeds; and
- audit events contain no secrets.

## Protocol adversarial cases

- unknown or inactive issuer;
- wrong audience or authorized party;
- invalid signature;
- prohibited or `none` algorithm;
- unknown or stale key;
- expired, premature, or implausibly issued token;
- missing or duplicate stable subject;
- state, nonce, PKCE, or redirect mismatch;
- authorization-code replay;
- malformed or duplicate security-sensitive claims;
- oversized token, claim, header, or callback body; and
- clock skew outside the accepted bound.

## Governed-resolution cases

- unknown or inactive provider;
- unmapped, duplicate, remapped, or ambiguous external identity;
- missing, disabled, or retired actor;
- missing, expired, premature, revoked, or conflicting role binding;
- provider role or group escalation;
- actor disablement during concurrent requests; and
- identity remapping while sessions are active.

## Request and proxy spoofing

- `X-Iron-Atlas-Actor` in production mode;
- `X-Iron-Atlas-Roles` in production mode;
- actor or roles in body, form, query, or path;
- forged forwarded address, scheme, or host;
- direct backend access bypassing the trusted proxy;
- host-header redirect injection; and
- duplicate or conflicting proxy headers.

## Session and CSRF cases

- preauthentication session fixation;
- identifier reuse after rotation;
- missing secure cookie attributes;
- identifier in URL, response, or log;
- plaintext identifier at rest;
- idle and absolute expiry;
- logout replay;
- administrative revocation;
- concurrent use while revocation occurs;
- missing, mismatched, expired, or replayed CSRF state;
- cross-origin mutation;
- mutation through a safe method; and
- unsupported content type.

## Confused-deputy assertions

Every privileged operation shall prove:

- the actor came only from immutable server-side context;
- permission came only from Atlas roles;
- PostgreSQL received the same actor;
- requester and approver independence still applies; and
- errors cannot produce a default actor or privilege fallback.

## Redaction assertions

Search logs, test output, retained evidence, panic text, and HTTP responses for
passwords, authorization codes, raw access or refresh tokens, raw session
cookies, CSRF secrets, private keys, client secrets, credential-bearing
database URLs, and full authentication assertions.

## Concurrency and resource observations

Run race and deterministic concurrency tests for multiple requests on one
session, rotation, logout, actor disablement, role revocation, identity
remapping, provider disablement, key rotation, and PostgreSQL pool reuse.

Record login latency, authenticated-request latency, key-cache behavior,
session-store operations, memory under malformed input, goroutines,
connections, cleanup cost, audit volume, cancellation, and timeouts separately
from correctness.

Performance thresholds remain `NOT_EVALUATED` until representative same-host
runs establish defensible budgets.

## Nonclaims

Passing Step 3 does not prove credential delivery, PostgreSQL TLS, backup,
high availability, live collection, protected evidence storage, all providers,
all proxies, or production readiness.

## Governed actor-resolution implementation campaign

The governed resolver checkpoint adds tests for:

- active provider, external identity, actor, and current-role resolution;
- inactive providers and disabled actors;
- unmapped and non-normalized subjects;
- expired bindings and inactive role definitions;
- unsupported and duplicate role codes;
- denial of direct application-role table reads;
- database unavailability classification;
- concurrent pooled resolution isolation; and
- Go race detection.

These tests do not exercise OIDC, browser sessions, CSRF, logout, replay
defense, or trusted-proxy behavior.

## OIDC ID-token verification implementation campaign

The bounded provider-emulator campaign now covers:

- exact TLS discovery and advertised protocol capabilities;
- valid RS256 verification through a discovered remote JWKS;
- wrong issuer, audience, authorized party, nonce, signature, and algorithm;
- expired, future, and stale token time state;
- future not-before state;
- missing and unnormalized stable subjects;
- duplicate security-sensitive claims;
- malformed and oversized tokens;
- remote key rotation and unknown-key refresh;
- provider outage classification;
- concurrent read-only verification under the race detector; and
- insecure or unbounded verifier configuration.
- valid, missing, and mismatched access-token hashes;
- valid and future authentication-time claims;
- duplicate protected JWT header fields;
- unknown-key rejection while the provider remains responsive; and
- non-success JWKS response classification.

That verifier campaign does not itself cover authorization-code replay, PKCE
transaction state, browser sessions, cookies, CSRF, logout, or trusted-proxy
behavior.

## OIDC authorization-code and PKCE transaction campaign

The active bounded campaign covers:

- authorization URLs with exact client ID, redirect URI, `openid` scope, state,
  nonce, PKCE challenge, and `S256`;
- absence of the client secret and PKCE verifier from the authorization URL;
- SHA-256 state-digest storage;
- successful code exchange and verified-principal production;
- exact redirect URI, client authentication, code, and PKCE verifier exchange;
- unknown, malformed, expired, and replayed state rejection;
- invalid authorization-code classification;
- provider-outage classification;
- token-response size bounds;
- bounded store capacity and expired-state cleanup;
- cryptographic-randomness failure;
- secret-redacted errors; and
- exactly one concurrent consumer and token exchange under the race detector.

The campaign does not cover HTTP login or callback handlers, browser state
cookies, durable restart-surviving transactions, authenticated sessions, CSRF,
logout, trusted proxies, governed actor wiring, or production credential
delivery.
