# MFA and Authentication-Assurance Requirements

## Status

Normative Phase 1 Step 3 requirements. Architecture alignment is accepted; the provider-neutral assurance-evidence implementation candidate is active and remains non-production.

## Authority boundary

User primary authentication and MFA occur at an approved external OpenID Connect identity provider.

Atlas:

- does not store user passwords;
- does not store provider TOTP seeds;
- does not generate provider MFA QR codes;
- does not implement ordinary local recovery codes;
- does not treat provider roles or groups as Atlas authorization; and
- does not permit silent administrator bypass.

Atlas remains responsible for deciding whether verified provider evidence satisfies Atlas policy.

## Assurance evidence

Atlas evaluates bounded issuer, provider and stable subject, `acr`, `amr`, `auth_time`, authentication age, provider profile, policy version, and requested role sensitivity. Presence of a claim is not proof that MFA occurred.

## Decisions

The result is satisfied, additional authentication required, phishing-resistant authentication required, stale authentication, unsupported assurance, ambiguous/conflicting evidence, or rejected. Only satisfied assurance may reach server-side session creation.

## MFA policy

Atlas policy declares exact accepted contexts and method sets. It requires production MFA, maximum authentication age, phishing-resistant authentication for governed high-impact roles or actions, downgrade prevention, exact policy-version binding, unknown-method rejection, and step-up that creates no stronger session until reauthentication succeeds.


## Provider-neutral assurance-evidence checkpoint

The active checkpoint uses only Atlas-controlled synthetic evidence. A
successful provider login is not MFA proof. Atlas rejects `acr` or `amr` unless
the same verified token includes explicit bounded `auth_time`, and accepted
method sets match exactly without ungoverned additional methods.

This checkpoint does not establish compatibility with any named provider.
Provider-specific semantics remain future evidence-backed work.

WebAuthn, FIDO2, passkeys, or hardware security keys are preferred phishing-resistant provider methods. Provider-managed RFC 6238 TOTP may be accepted only when an exact versioned policy permits it.

## Representative-provider compatibility

Before formal Step 3 acceptance, Atlas validates a disposable provider emulator and representative approved-provider profiles. Evidence covers discovery, authorization code with PKCE, issuer, audience, authorized party, signature, key rotation, time, exact `acr`/`amr`/`auth_time`, step-up, logout/provider-session behavior where applicable, outage, malformed responses, role-sensitive policy, and provider limitations.

Provider quirks are handled through versioned profiles, never permissive global parsing.

## Session lifecycle successor gates

Successor gates implement rotation, idle and absolute expiry, bounded session count and cleanup, logout, administrative revocation, actor/identity/provider/role invalidation, audit, CSRF, trusted proxy, and production wiring.

## Emergency and recovery access

Emergency access is not a local MFA bypass. A separate governed break-glass contract requires explicit request and reason, independent approval where possible, bounded time and scope, distinct emergency authority, strong external authentication, enhanced audit and notification, no silent role inheritance, immediate revocation, and post-use review.

Identity-provider recovery remains provider-owned. Atlas actor remapping or reactivation follows governed actor and separation-of-duty controls.

## Failure behavior

Atlas fails closed for missing, malformed, duplicate, stale, unknown, unsupported, conflicting, or insufficient assurance. No session cookie or persistent authenticated session is created.

## Required validation

Test downgrade, omitted and duplicate claims, unsupported contexts, method ambiguity, stale/future authentication time, clock skew, step-up, high-impact role policy, concurrent actor change, provider outage, key rotation, callback conflict, session handoff prevention, representative-provider compatibility, and secret redaction.

## Historical preservation

PR #15 and its signed post-merge boundary remain an implementation checkpoint. This alignment does not rewrite its code, gate, evidence, or nonclaims and is not formal Phase 1 Step 3 acceptance.
