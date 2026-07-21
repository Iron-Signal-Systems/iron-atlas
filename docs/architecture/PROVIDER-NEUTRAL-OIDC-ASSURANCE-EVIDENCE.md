# Provider-Neutral OIDC Assurance Evidence

## Status

Bounded Phase 1 Step 3 implementation candidate. This checkpoint extends the
merged authentication-assurance boundary from the exact signed architecture
alignment evidence-closure base
`2347d21f779768f40496a93cb1d9140cc3b6e0ce`.

It is not representative-provider compatibility, formal Phase 1 Step 3
acceptance, or production readiness.

## Purpose

Atlas must distinguish a successful external-provider login from evidence that
the login satisfies Atlas authentication-assurance policy. A valid signature,
issuer, audience, nonce, and subject prove a bounded OIDC identity assertion;
they do not prove MFA.

The provider-neutral path is:

```text
verified signed OIDC ID token
→ literal bounded acr, amr, and auth_time evidence
→ exact versioned Atlas assurance policy
→ satisfied, additional authentication required, step-up required, or denied
```

## Atlas-controlled evidence only

This checkpoint uses the existing disposable Atlas-controlled OIDC provider
emulator and synthetic values reserved for tests. It does not encode or claim
behavior for any named commercial, cloud, or self-hosted provider.

Synthetic examples such as `urn:iron-atlas:test:mfa`, `pwd`, and `otp` are test
fixtures, not interoperability claims and not production policy guidance.

## Mandatory correlation

`acr` and `amr` are accepted as assurance evidence only when the same verified
ID token contains an explicit, bounded `auth_time`. Atlas rejects a token that
presents `acr` or `amr` without `auth_time`.

A token with no assurance claims may still establish a verified provider
identity, but it remains non-MFA evidence and cannot satisfy a policy requiring
MFA.

## Exact governed method sets

An accepted method set is exact. A configured set such as `pwd` plus `otp` does
not accept a third ungoverned method merely because the configured values are a
subset. Stronger, additional, renamed, or unknown methods require an explicit
new policy version and evidence-backed review.

Context acceptance remains exact and versioned. No claim is accepted through
substring, case-insensitive, prefix, suffix, or vendor-name inference.

## Failure behavior

Atlas fails closed for:

- assurance claims without explicit `auth_time`;
- stale or future authentication time;
- missing, malformed, duplicate, excessive, or unnormalized assurance claims;
- unknown contexts;
- partial, additional, or unknown method sets;
- policy-version mismatch; and
- any attempt to treat login success as MFA proof.

Unsatisfied evidence creates no stronger session and no MFA assurance record.

## Deliberate non-wiring and nonclaims

This checkpoint preserves the existing provider-neutral verifier and assurance
service seams. It does not add provider profiles, production configuration,
provider-hosted enrollment, local credentials, local TOTP, recovery codes,
WebAuthn implementation, session rotation, logout, administrative revocation,
CSRF, trusted-proxy enforcement, emergency access, or production application
wiring.

Representative-provider work begins only after sanitized evidence exists from
a controlled provider configuration and its observed behavior can be compared
to this provider-neutral baseline.
