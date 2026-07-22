# Representative-Provider Evidence Foundation Testing

## Scope

This model validates the evidence-bundle contract, sanitized synthetic fixture,
digest binding, secret rejection, literal-claim preservation, deterministic
output, repository registration, and exact predecessor revalidation.

It does not test a live or named external identity provider and does not
establish compatibility.

## Positive cases

- strict schema version and closed object shapes;
- synthetic or controlled-sanitized evidence classification;
- observation-only claim status and compatibility claim fixed to false;
- exact provider software, version, configuration digest, and capture identity;
- HTTPS issuer and digest-bound discovery and JWKS artifacts;
- unique scenario and artifact identities;
- exact literal `acr`, ordered `amr`, and `auth_time`;
- scenario claims equal the referenced sanitized claims artifact;
- explicit redaction status and limitations; and
- deterministic validation output.

## Negative and hostile cases

The regression rejects compatibility claims, raw-artifact commitment,
credential or identity-bearing keys, JWT-shaped values, private-key material,
path traversal, duplicate artifact paths and scenario identifiers, missing or
mismatched digests, claims-artifact disagreement, interpreted MFA fields,
malformed types, and missing nonclaims.

## Evidence boundary

The committed fixture is synthetic and uses reserved example identifiers. It
proves the contract and validator, not provider behavior. Controlled provider
captures remain outside the repository until sanitization, review, digest
binding, and a separate evidence-backed change.

## Validation sequence

1. Confirm the signed provider-neutral boundary is an ancestor.
2. Revalidate its exact phase gate in an isolated checkout.
3. Validate the evidence schema and committed synthetic fixture.
4. Run hostile mutation regression and determinism checks.
5. Run the complete test framework and repository validation.
6. Preserve explicit compatibility, production, and acceptance nonclaims.


## Historical checkpoint isolation

Architecture-alignment and provider-neutral validators contain status assertions
that are correct only at their accepted predecessor boundaries. Successor
documentation must not retain stale active-checkpoint wording merely to satisfy
those historical validators.

The complete framework and repository validator execute those status-sensitive
checks from isolated clones of exact predecessor commits. The active foundation
validator evaluates the successor documentation and current checkpoint status.
