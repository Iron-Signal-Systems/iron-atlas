# Phase 1 Step 3 Acceptance Record Template

## Record Status

- Decision:
- Product: Iron Atlas
- Repository: `Iron-Signal-Systems/atlas`
- Canonical repository: `https://github.com/Iron-Signal-Systems/atlas.git`
- Branch: `dev`
- Acceptance date:
- Accepted tag:

## Predecessor

- Accepted Phase 1 Step 2 tag:
  `phase-1-step-2-go-postgresql-runtime-and-identity-context-complete-v1`
- Accepted Phase 1 Step 2 `dev` merge boundary:
  `1a750f7de791f567184c6f48e18eaec2933b8a14`
- Isolated predecessor revalidation:

## Candidate

- Candidate implementation commit:
- HTTP login and callback checkpoint:
- Authenticated server-side session checkpoint:
- Authentication assurance and MFA checkpoint:
- Provider-neutral assurance-evidence checkpoint:
- Representative-provider evidence-foundation checkpoint:
- Representative-provider compatibility checkpoint:
- Session lifecycle and revocation checkpoint:
- Emergency and recovery access checkpoint:
- Repository-complete evidence boundary commit:
- Candidate Git archive SHA-256:
- Toolchain requirements SHA-256:
- Go version:
- Authentication adapter and version:
- Provider emulator or representative provider:
- PostgreSQL versions:
- Non-secret host-class fingerprint:

## Accepted Scope

- Trusted authentication adapter:
- Verified provider identity normalization:
- Governed external-identity and actor resolution:
- Atlas-owned role-binding resolution:
- Immutable request-context identity:
- Bounded server-side sessions:
- Cookie security:
- Authentication assurance and MFA policy:
- Provider-owned MFA assurance and compatibility:
- Emergency and recovery access policy:
- CSRF:
- Trusted proxy:
- Audit and redaction:
- Transaction-local PostgreSQL actor propagation:
- Readiness and failure behavior:

## Protocol Security Profile

- Issuer:
- Audience and authorized-party policy:
- Signature algorithms:
- Key discovery and rotation:
- State:
- Nonce:
- PKCE:
- Redirect URI policy:
- Clock skew:
- Maximum sizes:
- Replay protections:

## Session and Browser Security Profile

- Session identifier:
- Persistent representation:
- Cookie attributes:
- Idle lifetime:
- Absolute lifetime:
- Rotation:
- Revocation:
- CSRF mechanism:
- Origin policy:
- Cleanup and resource bounds:
- Authentication context and methods:
- MFA age and step-up policy:
- Preferred phishing-resistant provider authenticators:
- Provider-managed TOTP compatibility, when permitted:
- Emergency and recovery access:

## Trusted Proxy Profile

- TLS termination:
- Trusted peer addresses or socket:
- Accepted forwarded headers:
- Direct-backend bypass prevention:
- Redirect host and scheme validation:

## Validation Evidence

- Local implementation gate:
- Canonical clean-clone validation:
- Canonical clone commit:
- Applicable validator:
- Committed evidence:
- Evidence checksums:
- Predecessor revalidation:
- Repository validation:
- Test framework:
- Protocol-emulator tests:
- Positive cases:
- Negative and adversarial cases:
- Identity-ambiguity cases:
- Session and CSRF cases:
- Trusted-proxy cases:
- Concurrency and race cases:
- PostgreSQL identity-isolation regression:
- Secret-redaction result:
- Correctness result:
- Resource observation:
- Performance thresholds:

## Reproducibility Statement

Confirm that the exact pushed commit was validated from a clean canonical clone
using version-controlled artifacts, declared and verified tooling, disposable
test environments, explicitly supplied non-repository secrets, and a fully
identified provider-emulator or representative-provider boundary.

## Review and Approval

- Requester and implementer:
- Independent reviewer:
- Conflicts checked:
- Temporary development exception, when applicable:
- Production-boundary exception prohibited:

## Explicit Exclusions

- Collector or device credential delivery and rotation
- PostgreSQL TLS and certificate deployment
- Backup and restoration
- High availability
- Live infrastructure collection
- Protected evidence intake and storage
- Automated remediation
- Production performance budgets
- Production readiness

## Decision

Record the exact tested authentication, actor-resolution, session, CSRF,
trusted-proxy, audit, and transaction-local identity boundary. State no broader
claim.

## Exact Next Work

State the next Phase 1 boundary without combining it into this acceptance.

## Authentication assurance checkpoint evidence

- Exact `acr`, `amr`, and `auth_time` normalization:
- Versioned assurance policy:
- Required MFA enforcement before session creation:
- High-impact phishing-resistant policy behavior:
- Stale-authentication step-up behavior:
- Callback issuer and metadata hardening:
- Unsatisfied-assurance no-cookie evidence:
- Concurrent and hostile validation:
- Explicit provider compatibility and production nonclaims:

## Provider-neutral assurance-evidence checkpoint

- Exact signed implementation base:
- Atlas-controlled provider-emulator evidence:
- Successful-login-without-MFA result:
- `acr` or `amr` without `auth_time` rejection:
- Exact method-set enforcement:
- Unknown and additional method rejection:
- Stale-authentication step-up result:
- Named-provider compatibility explicitly not claimed:

## Representative-provider evidence-foundation checkpoint

- Exact signed provider-neutral predecessor:
- Evidence schema and validator version:
- Synthetic fixture result:
- Controlled provider software, version, and configuration digest:
- Sanitization and review status:
- Discovery and JWKS artifact digests:
- Literal `acr`, `amr`, and `auth_time` observations:
- Secret, identity, JWT, private-key, and path-traversal rejection:
- Deterministic validation result:
- Compatibility claim explicitly false:
