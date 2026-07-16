# Trusted Authentication and Governed Actor Resolution Testing

## Status

Phase 1 Step 3 test contract integrated; authentication-foundation tests active.

These tests are required for the final executable candidate. This document does
not claim that a production adapter or session implementation exists.

## Authentication foundation implementation checkpoint

The current candidate adds executable tests for typed mode parsing, default mode selection, legacy-setting rejection, bounded development headers, duplicate and unknown role rejection, immutable context copies, production rejection of development headers, missing-adapter fail-closed behavior, adapter/resolver composition, nested middleware rejection, public health/readiness paths, and query-string actor spoofing.

These tests do not substitute for provider-protocol, database-backed actor-resolution, session, CSRF, replay, logout, key-rotation, or trusted-proxy test campaigns.


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
