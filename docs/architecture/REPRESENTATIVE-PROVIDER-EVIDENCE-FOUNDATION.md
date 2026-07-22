# Representative-Provider Evidence Foundation

## Status

Bounded Phase 1 Step 3 implementation candidate. This checkpoint extends the
signed provider-neutral OIDC assurance-evidence boundary
`e7824049852855f15d26686600fc42802b8a38ff`.

It establishes a sanitized evidence contract and deterministic offline
validation boundary. It is not representative-provider compatibility, live
provider integration, formal Phase 1 Step 3 acceptance, or production readiness.

## Purpose

Representative-provider compatibility cannot be asserted from source code,
vendor documentation, successful login, or remembered provider behavior.
Compatibility work begins with controlled, reproducible, sanitized observations
that preserve literal provider output without turning observations into policy
or interoperability claims.

The foundation path is:

```text
disposable controlled provider configuration
→ raw capture retained outside Git
→ explicit sanitization and review
→ digest-bound observation bundle
→ deterministic repository validation
→ later provider-specific evidence review
→ separately governed compatibility profile or rejection
```

## Evidence classes

The bundle contract permits only `synthetic` evidence for repository fixtures
and hostile tests, or `controlled-sanitized` evidence produced from an
identified disposable provider configuration after explicit review.

Every bundle is `observation-only` and sets `compatibility_claim` to `false`.

## Required provider and capture identity

A bundle records an opaque provider label, provider software identity, exact
version, configuration SHA-256, capture timestamp, disposable-environment
declaration, capture-tool identity and version, and the sanitized issuer.
Provider identity is reproducibility metadata; it grants no semantic meaning to
`acr`, `amr`, `auth_time`, endpoints, authenticators, or session behavior.

## Literal assurance observations

Each scenario preserves authentication success, literal `acr`, ordered `amr`,
literal `auth_time`, a digest-bound sanitized claims artifact, purpose, and
limitations. The bundle contains no Atlas policy result, MFA-satisfied boolean,
phishing-resistant classification, strength ranking, compatibility result, or
provider-name inference.

## Digest and path binding

Every artifact is named explicitly, constrained to the bundle directory,
required to use a sanitized filename, and bound by SHA-256. Discovery and JWKS
references must resolve to artifacts with the matching kind. Scenario claims
must exactly match their referenced sanitized artifact. Duplicate paths,
duplicate scenario identifiers, traversal, missing files, digest mismatch, and
interpretation fields fail closed.

## Prohibited repository material

The validator rejects raw JWT-shaped values, private-key material, and
credential or identity-bearing keys including passwords, secrets, access or
refresh tokens, ID tokens, authorization codes, PKCE verifiers, cookies, raw
subjects, usernames, email addresses, session identifiers, TOTP seeds, and
recovery codes.

Raw provider traffic, browser storage, cookies, token responses, credentials,
private keys, and unredacted identities remain outside Git. Sanitization does
not convert a raw capture into a compatibility claim.

## Deliberate non-wiring and nonclaims

This checkpoint adds no network capture, provider installation, provider
profile, runtime adapter, production configuration, authentication semantic
mapping, session behavior, or emergency-access behavior.

No named-provider compatibility gate exists in this checkpoint.
Representative-provider compatibility remains future evidence-backed work after
a controlled provider configuration produces reviewed evidence that can be
compared with the signed provider-neutral baseline.


## Validation succession

Accepted predecessor validators remain immutable evidence of their own repository
state. Status-sensitive architecture-alignment and provider-neutral checks run in
isolated exact-commit clones rather than against successor documentation. The
current foundation gate still exact-revalidates the complete signed
provider-neutral boundary before evaluating this candidate.
