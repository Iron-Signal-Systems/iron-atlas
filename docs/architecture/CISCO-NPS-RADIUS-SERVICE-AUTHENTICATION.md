# Cisco NPS/RADIUS Service Authentication

## Purpose

Authenticate the Atlas Cisco collection service through Microsoft Windows Server NPS and Active Directory while enforcing the actual read-only command boundary on each Cisco platform.

## Separation

Use different Active Directory groups and NPS policies for:

- Endpoint IEEE 802.1X access
- Human device administration
- Read-only human investigation
- Automated Atlas collection

802.1X group membership never grants Cisco SSH management access.

Suggested groups:

```text
GG-NETWORK-8021X-USERS
GG-NETWORK-8021X-COMPUTERS
GG-NETWORK-DEVICE-ADMINS
GG-NETWORK-DEVICE-READONLY
GG-ATLAS-CISCO-COLLECTORS
```

## Service Identity

The collector uses a dedicated nonhuman domain account. It is not a Domain Administrator, general network administrator, email account, or interactive workstation account. Its secret remains in an approved secret manager and is never written to the repository, device manifest, logs, or evidence bundle.

A conventional service account may be required for password-based RADIUS authentication. Do not assume that a gMSA can be presented as an ordinary Cisco SSH password credential.

## Authorization Flow

```text
Atlas collector
  → SSH to enrolled Cisco device
  → Cisco RADIUS request to NPS
  → NPS evaluates device-management policy and AD group
  → NPS returns approved Cisco privilege attributes
  → Device assigns restricted parser view or privilege level
  → Collector executes versioned read-only profile
```

NPS/RADIUS authenticates and assigns session-level privilege. The Cisco device enforces the local command set. This is not TACACS+-equivalent centralized per-command authorization.

## Required Denials

Enrollment testing must prove the service identity cannot:

- Enter configuration mode
- Modify interfaces, VLANs, routing, AAA, users, or logging
- Save configuration
- Reload devices or processes
- Install software
- Delete or format storage
- Enable unrestricted debugging
- Copy arbitrary files
- Use an enable password or unrestricted fallback

## Failure

NPS or Active Directory failure causes the collection to fail closed. The collector does not fall back to a human account, shared administrator account, unrestricted local account, or break-glass credential.

## Accounting

Retain NPS session evidence, device AAA evidence, and the Atlas-signed command transcript. Ordinary RADIUS accounting is not treated as complete command accounting.
