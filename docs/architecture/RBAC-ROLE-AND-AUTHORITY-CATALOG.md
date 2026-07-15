# RBAC Role and Authority Catalog

## Principle

Roles organize eligible capabilities. They do not replace exact operation, resource, scope, purpose, risk, session, approval, and current-state evaluation.

## Initial Roles

### Viewer

May view nonrestricted dashboards and inventory allowed by data classification. Cannot collect evidence, request changes, approve changes, or administer modules.

### Network Technician

May inspect assigned infrastructure, run approved collection profiles, investigate access ports, review redacted evidence, and request changes. Cannot approve the technician’s own change.

### Network Administrator

May perform governed network administration, design and request changes, administer collection coverage, and approve changes only when independent and within authority.

### Network Security

May review firewall, ACL, management-plane, authentication, and security-control state; approve security-relevant changes when independent; and review audit records.

### Change Approver

May approve or reject changes within an exact delegated scope. This role alone does not grant implementation authority or platform administration.

### Auditor

May read governed history, evidence lineage, approval, acceptance, exception, and audit records. It does not grant operational modification authority.

### Platform Administrator

May administer Iron Atlas deployment and module configuration. It does not automatically grant infrastructure change approval, raw-evidence access, or unrestricted database authority.

## Scope Dimensions

Role bindings should be constrained by organization, site, device group, device, capability, operation, environment, data classification, time, and delegation lineage.

## Prohibited Accumulation

No ordinary role accumulation creates an unrestricted account. High-impact operations may require independent approval, step-up authentication, exact target binding, short-lived authority, and a durable decision record.
