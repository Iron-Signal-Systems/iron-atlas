# Provider-Neutral OIDC Assurance Evidence Testing

## Scope

This campaign validates the bounded provider-neutral assurance-evidence
checkpoint. All claim values are produced by Atlas-controlled fixtures. No live
or named external identity provider is required or represented.

## Required positive cases

- a signed provider-emulator token with no assurance claims remains a valid
  non-MFA principal;
- exact synthetic `acr`, exact synthetic `amr`, and explicit bounded
  `auth_time` reach the existing assurance policy;
- an exact configured method set is accepted;
- accepted assurance retains the exact policy version; and
- concurrent read-only verification and evaluation remain race-free.

## Required negative and hostile cases

- successful login without assurance evidence never proves MFA;
- `acr` without `auth_time` is rejected;
- `amr` without `auth_time` is rejected;
- `acr` and `amr` without `auth_time` are rejected;
- a partial method set is rejected;
- an additional ungoverned method is rejected;
- an unknown context or method is rejected;
- stale authentication requires step-up;
- future authentication time is denied;
- duplicate and malformed assurance claims are rejected; and
- unsatisfied assurance reaches neither cookie issuance nor session creation.

## Validation sequence

1. static contract validation;
2. focused Go formatting;
3. focused `go test -race`;
4. focused `go vet`;
5. module verification;
6. vulnerability analysis;
7. complete test framework;
8. repository validation; and
9. SSH-signed candidate commit.

## Nonclaims

Passing this campaign does not establish interoperability with any named
provider, live provider configuration, hosted MFA behavior, phishing-resistant
provider behavior, completed session lifecycle, formal Phase 1 Step 3
acceptance, or production readiness.
