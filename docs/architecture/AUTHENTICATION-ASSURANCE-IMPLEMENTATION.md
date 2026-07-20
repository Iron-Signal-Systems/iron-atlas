# Authentication Assurance Implementation

## Status

Phase 1 Step 3 bounded implementation candidate. This checkpoint establishes provider-neutral authentication-assurance normalization and policy enforcement between the verified OIDC principal and authenticated-session creation. It does not establish local TOTP enrollment, WebAuthn enrollment, session lifecycle completion, production application wiring, formal Step 3 acceptance, or production readiness.

## Purpose

Iron Atlas must not issue a fully authenticated server-side session merely because an OIDC provider returned a valid signed identity token. Atlas separately governs whether the verified authentication event satisfies the configured assurance policy for the resolved Atlas actor.

The enforced path is:

```text
verified OIDC principal
→ current governed actor resolution
→ provider-neutral assurance policy
→ accepted provider MFA or stronger assurance
→ policy-version binding
→ authenticated server-side session
```

Unsatisfied assurance does not create a session cookie or persistent authenticated-session record. It returns a generic additional-authentication-required result for a future Atlas-local TOTP, WebAuthn, passkey, hardware-key, or provider step-up boundary.

## Provider-neutral claims

The OIDC verifier normalizes:

- `acr` into one bounded authentication-assurance context;
- `amr` into a bounded, duplicate-free list of authentication methods;
- `auth_time` into the primary authentication time.

The verifier does not infer MFA merely because `acr` or `amr` exists. Provider claims remain evidence. Atlas policy determines whether an exact configured context or method set is accepted.

Provider roles, groups, directory memberships, and authorization claims do not become Atlas roles or authority.

## Policy model

One versioned policy defines:

- whether MFA is mandatory;
- maximum accepted authentication age;
- maximum clock skew;
- exact accepted MFA contexts;
- exact accepted MFA method sets;
- roles requiring phishing-resistant assurance;
- exact accepted phishing-resistant contexts; and
- exact accepted phishing-resistant method sets.

Method-set evaluation is explicit subset matching. Unknown, partial, malformed, stale, or ambiguous claims do not satisfy MFA.

The policy produces one bounded outcome:

- `satisfied`;
- `additional_authentication_required`;
- `step_up_required`;
- `phishing_resistant_required`; or
- `denied`.

Only `satisfied` may reach the authenticated-session service.

## Defense in depth

The assurance service resolves the current governed actor before policy evaluation. The session service resolves the actor again immediately before persistence, protecting the boundary from actor remapping, disabling, or role changes between policy evaluation and session creation.

The authenticated-session service also requires:

- MFA assurance marked successful;
- a nonzero MFA authentication time; and
- exact equality between the principal assurance policy version and the session service policy version.

The session service does not manufacture or overwrite assurance success.

## HTTP callback hardening

The login and callback handler retains the exact issuer validated during construction. Callback issuer comparison uses that immutable value rather than asking a replaceable flow for the issuer again.

Provider error metadata is accepted only on error responses. `session_state` is accepted only on successful authorization-code callbacks. Conflicting metadata fails closed.

## Security properties

This candidate proves:

- exact `acr` and `amr` normalization;
- duplicate-sensitive-claim rejection;
- explicit provider-MFA acceptance policy;
- no MFA inference from unknown claims;
- role-sensitive phishing-resistant policy;
- stale-authentication step-up decisions;
- no session handoff for unsatisfied assurance;
- policy-version binding before session creation;
- generic browser failure responses;
- no session cookie on failure;
- stateless concurrent policy evaluation; and
- preservation of the accepted governed actor and server-side session boundaries.

## Deliberate exclusions

This checkpoint does not implement or accept:

- TOTP secret generation;
- QR-code enrollment;
- TOTP verification or replay counters;
- recovery codes;
- authenticator reset or replacement;
- WebAuthn, FIDO2, passkeys, or hardware security keys;
- session rotation or sliding idle refresh;
- logout or administrative revocation workflow;
- CSRF protection;
- trusted-proxy enforcement;
- production `atlasd` wiring;
- representative-provider compatibility;
- authentication audit persistence;
- formal Phase 1 Step 3 acceptance; or
- production readiness.

## Next boundary

The successor checkpoint is the governed RFC 6238 TOTP lifecycle: encrypted per-authenticator secrets, provisional enrollment, standards-compatible QR enrollment, proof of possession, 30-second verification, atomic replay prevention, rate limits, one-time recovery codes, governed reset, key rotation, and durable audit evidence.
