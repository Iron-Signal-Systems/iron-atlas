# Change Management and Two-Person Control

## Governing Principle

A governed infrastructure change requires at least two independently accountable people: a requester/implementer and an authorized approver who is not the requester or directly affected actor.

Higher-risk operations may require more than one independent approver, including network-security approval or separate organizational authority.

## Lifecycle

```text
Draft
  ↓
Discovery and evidence
  ↓
Requirements and impact
  ↓
Implementation and rollback plan
  ↓
Pending independent approval
  ↓
Approved
  ↓
Pre-change collection
  ↓
Implementation
  ↓
Post-change collection and validation
  ↓
Documentation reconciliation
  ↓
Formal acceptance
  ↓
Closed
```

## Independence Rules

An approval is invalid when:

- The approver is the requester.
- The approver is the directly affected identity where policy prohibits it.
- The same effective actor is represented through multiple accounts or delegated paths.
- The required authority was expired, revoked, or out of scope.
- Reciprocal approval arrangements violate policy.
- Incompatible duties are held for the operation.
- The approval occurred outside the current policy stage or after supersession.

## Roles

Initial role concepts include:

- `viewer`
- `network_technician`
- `network_administrator`
- `network_security`
- `change_approver`
- `auditor`
- `platform_administrator`

No accumulated ordinary role set grants unrestricted authority. Platform administration does not automatically grant change approval.

## Durable Records

Retain:

- Request and requester
- Business and technical justification
- Affected sites and devices
- Risk and outage impact
- Current-state evidence
- Target-state design
- Implementation steps
- Validation plan
- Rollback plan and rollback decision
- Approval actions and authority context
- Implementation transcript
- Post-change evidence
- Expected and unexpected differences
- Documentation reconciliation
- Acceptance and exceptions

Material records use correction and supersession rather than silent rewriting.
