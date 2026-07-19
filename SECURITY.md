# Security Policy

Iron Atlas is pre-alpha and is not ready for production use. No production
support or security service-level commitment is claimed.

## Repository assurance baseline

This repository records adoption of:

- Standard: Iron Signal Repository Assurance Standard (ISRAS)
- Version: `1.0.1`
- Signed tag: `isras-v1.0.1`
- Exact commit: `c379417720faa595fa5cb89a1dfdb2259d6cb95e`
- Adoption level: `RECORDED`
- Enforcement mode: observation

The version string alone does not prove the adopted standard identity. The
signed tag and exact commit are authoritative.

## Never commit

- Raw firewall configuration backups
- Cisco technical-support reports
- Unredacted command transcripts
- Passwords, tokens, shared secrets, private keys, or certificate private
  material
- Production IP inventories not approved for repository storage
- NPS, TACACS+, SNMP, RADIUS, VPN, or API credentials
- Unsanitized validation transcripts or environment dumps
- Production data or unrestricted logs

## Validation evidence

Only sanitized, checksummed validation evidence may be committed under
`validation/evidence/` or another explicitly governed evidence directory.
Validators must not enumerate ambient environment variables or persist
repository-external secrets.

## Reporting a vulnerability

Do not open a public issue containing exploit details, credentials, protected
data, sensitive infrastructure details, or instructions that materially
increase risk. Use a private, approved Iron Signal Systems reporting channel.

Include the affected commit or release, affected component, reproduction
conditions, potential impact, possible protected-data exposure, and known
mitigation.

## Response handling

Reports are acknowledged, triaged, contained, remediated, revalidated, and
disclosed according to risk. Response timing is risk-based and is not a
contractual service-level commitment.

## Governing principle

A successful parser or collection proves only that evidence was obtained and
interpreted to the tested extent. It does not prove a device is healthy,
secure, correctly configured, or synchronized with documentation.
## ISRAS 0.1.4 adoption candidate

Atlas is evaluating the exact published `isras-v0.1.4` release at source commit
`c9345d6d731600df7bd4ba4a133c07265db55e5a`. The generated project pin and
hosted workflow are a candidate until separate prospective acceptance. The
previously recorded 1.0.1 boundary remains historical and no production or
independent-review claim is added.
