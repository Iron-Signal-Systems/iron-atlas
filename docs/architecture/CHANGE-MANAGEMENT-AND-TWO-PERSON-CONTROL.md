# Change Management and Two-Person Control

## Purpose

Use Atlas evidence and analysis to help engineers, security staff, operational leaders, and change authorities understand, approve, deny, implement, validate, roll back, accept, and later reconstruct infrastructure changes.

## Governing Principle

A governed infrastructure change requires at least two independently accountable people:

- a requester or implementer; and
- an authorized approver who is not the requester and is not disqualified by the applicable separation-of-duties policy.

Higher-risk changes may require additional independent technical, security, operational, or organizational approval.

## Decision Support Principle

A change request shall explain both:

- **the risk of approving and implementing the change**; and
- **the risk of denying, delaying, or leaving the existing condition unchanged**.

A decision authority should not be forced to approve or deny a technical request without understandable operational, security, availability, usability, cost, and rollback context.

## Lifecycle

```text
Observed condition or requested outcome
                  ↓
Discovery and evidence
                  ↓
Current-state reconstruction
                  ↓
Proposed target state
                  ↓
Reachability, attack-path, dependency, and blast-radius analysis
                  ↓
Implementation, validation, and rollback plan
                  ↓
Independent review and decision
       ┌──────────┴──────────┐
     denied             approved
       ↓                    ↓
reason and residual     pre-change evidence
risk retained                ↓
                        implementation
                              ↓
                    post-change evidence
                              ↓
                  validation and comparison
                       ┌──────┴──────┐
                    rollback      accept
                       ↓              ↓
                 validation     documentation
                       ↓         reconciliation
                    closure            ↓
                                  formal acceptance
                                        ↓
                                      closed
```

## Director and Change-Authority View

The decision-facing view explains:

- the problem or requested outcome;
- why the change is needed;
- operational and security benefit;
- risk if approved;
- risk if denied or delayed;
- affected departments, services, sites, systems, and users;
- expected outage, degradation, or maintenance effect;
- implementation window;
- rollback readiness and expected recovery time;
- evidence quality;
- known unknowns;
- recommendation and confidence; and
- the exact decision requested.

Technical detail remains available but does not overwhelm the decision summary.

## Engineering and Security View

The technical view explains:

- affected devices and systems;
- current and proposed configurations;
- topology before and after;
- VLANs, CIDRs, routes, policies, ACLs, NAT, VIPs, VPNs, SD-WAN, wireless, and dependencies affected;
- traffic paths created, removed, or changed;
- trust boundaries crossed;
- network and identity attack paths created, widened, narrowed, or removed;
- BloodHound or other identity-graph context affecting critical assets, privileged identities, and administrative paths;
- availability, redundancy, and failure-domain effects;
- implementation order;
- commands or reviewable configuration snippets;
- pre-change checks;
- validation commands and expected results;
- rollback commands and rollback triggers;
- post-change monitoring; and
- evidence supporting each conclusion.

## Change Package

A governed change package retains:

- request and requester;
- business and technical justification;
- problem statement;
- current-state evidence;
- target-state design;
- affected scope;
- operational benefit;
- security benefit;
- risk of approval;
- risk of denial or delay;
- outage and usability impact;
- dependency and blast-radius analysis;
- reachability and attack-path comparison;
- implementation plan;
- validation plan;
- rollback plan;
- approval policy;
- approval, denial, or revision decisions;
- authority context;
- pre-change evidence;
- implementation transcript;
- post-change evidence;
- expected and unexpected differences;
- human disposition;
- documentation reconciliation;
- acceptance or rollback; and
- exceptions and supersession.

## Independence Rules

An approval is invalid when:

- the approver is the requester;
- the approver is the directly affected identity where policy prohibits it;
- the same effective actor is represented through multiple accounts or delegated paths;
- required authority is expired, revoked, inactive, or out of scope;
- reciprocal approval arrangements violate policy;
- incompatible duties are held for the operation;
- the decision occurred outside the applicable stage;
- the change package materially changed after approval without required reapproval; or
- the decision applies to a superseded change revision.

## Roles

Initial role concepts include:

- `viewer`;
- `network_technician`;
- `network_administrator`;
- `network_security`;
- `change_approver`;
- `operational_leader`;
- `auditor`; and
- `platform_administrator`.

No accumulated ordinary role set grants unrestricted authority.

Platform administration does not automatically grant infrastructure-change approval.

## Implementation Boundary

Atlas initially remains read-only with respect to infrastructure.

Atlas may generate:

- proposed configuration;
- commands;
- implementation sequence;
- validation commands;
- rollback commands;
- diagrams;
- reports; and
- evidence packages.

A human engineer executes the change through an authorized operational process.

Any future Atlas-driven provisioning requires a separately accepted automation boundary and does not inherit authority merely because a change record is approved.

## Validation and Acceptance

A change is not accepted merely because implementation commands completed.

Acceptance requires:

- post-change evidence;
- expected-state confirmation;
- unexpected-difference review;
- required service and security tests;
- monitoring review;
- rollback decision where applicable;
- documentation reconciliation;
- human disposition; and
- formal acceptance by the required authority.

## Historical Integrity

Material records use correction and supersession rather than silent rewriting.

Atlas shall be able to answer:

- Why was this route, policy, VLAN, ACL, VPN, or exception created?
- Who requested and approved it?
- What risk was understood?
- What evidence existed?
- What was implemented?
- What validation passed or failed?
- Was rollback required?
- Which state was formally accepted?
