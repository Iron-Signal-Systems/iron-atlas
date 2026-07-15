# Testing and Adversarial Validation Model

## Current Phase 0 Tests

- Go compilation and unit tests
- RBAC separation
- Requester cannot self-approve
- Duplicate concurrent approval suppression
- HTML response and security headers
- FortiGate hierarchy parsing
- Malformed FortiGate rejection
- Cisco trunk endpoint-attribution exclusion
- Trunk required-evidence findings
- Zabbix sender packet encoding and response decoding
- Static repository structure and documentation-link validation

## Required Future Tests

- Malformed, truncated, oversized, compressed, encrypted, and hostile input
- Parser hangs and resource exhaustion
- Vendor version variations
- Unknown sections and forward compatibility
- SSH timeout, disconnect, host-key change, and command rejection
- NPS outage and fail-closed behavior
- Concurrent collection and device-load limits
- Evidence signature, sequence, replay, and duplicate handling
- Approval race, withdrawal, supersession, and reciprocal-approval protection
- Pre/post change mismatch
- Topology loops and ambiguous endpoint evidence
- Zabbix outage, slow destination, retry, and dead-letter behavior
- PostgreSQL privilege and independent-connection concurrency
- Backup, restore, upgrade, and compromise recovery
- Keyboard, screen-reader, zoom, contrast, and responsive UI behavior

## Resource Observation

Record CPU, memory, I/O, database, evidence size, parser duration, command duration, queue depth, external-delivery latency, and host fingerprints separately from correctness. Begin observation-only; do not create unsupported performance gates.
