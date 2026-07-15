# Evidence Ingestion and Protection

## Evidence Types

- Firewall configuration backup
- Cisco technical-support report
- Platform-specific supplemental report
- Read-only command transcript
- Device inventory export
- Diagram source
- Pre-change collection
- Post-change collection
- Incident collection
- Manual authorized import

## Intake Flow

```text
Collector or file import
      ↓ signed authenticated batch
Staging and ingestion service
      ├── authenticate source
      ├── enforce size and type limits
      ├── hash content
      ├── sequence and deduplicate
      ├── persist raw evidence durably
      ├── invoke isolated parser
      ├── record parser version and warnings
      └── publish normalized candidate records
```

Collectors have no direct PostgreSQL access. Endpoint-local storage is only a bounded encrypted outage buffer.

## Raw Evidence

Raw evidence is sensitive infrastructure material and may contain credentials, password hashes, internal addresses, device names, SNMP communities, VPN details, certificate material, AAA configuration, and security policy.

It shall be:

- Encrypted in transit and at rest
- Content-addressed by SHA-256
- Access-controlled by role
- Stored outside Git
- Retained under policy
- Exportable for authorized vendor support
- Covered by backup, restoration, integrity, and compromise-recovery procedures

## Parser Safety

Parsers operate with bounded memory, bounded input size, timeouts, no network access by default, and no authority to modify source devices. Unsupported sections are retained and reported, not discarded.
