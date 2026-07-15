# Target Architecture

## Status

Normative target direction. Phase 1 Step 1 is accepted; the current executable includes a Phase 1 Step 2 PostgreSQL runtime candidate.

## Layers

```text
HTML5 user interface and versioned API
                 ↓
Governed application services
  inventory · topology · projects · changes · approvals · findings
                 ↓
Canonical normalized records and decision-supporting evidence
                 ↓
Parser and analyzer module contracts
                 ↓
Vendor adapters and collection profiles
                 ↓
Signed ingestion boundary and protected raw evidence storage
                 ↓
Authorized read-only devices and imported configuration backups
```

Replaceable external adapters consume canonical telemetry and records:

```text
Canonical Atlas telemetry/outbox
        ├── Zabbix sender adapter
        ├── OpenMetrics adapter
        ├── Syslog adapter
        ├── Webhook adapter
        └── Future monitoring systems
```

## Process Direction

The target deployment separates:

- `atlas-api`: HTML5, API, query, and governed workflow service
- `atlas-worker`: parsing, analysis, diagram, and delivery jobs
- `atlas-ingest`: authenticated evidence intake and sequencing
- `atlas-collector`: site-scoped read-only collection
- PostgreSQL: authoritative normalized state and governed records
- Protected evidence store: encrypted raw evidence by content hash

The current candidate still combines the UI/API service in one process. Phase 1 Step 2 adds persistent change workflow through a least-privileged PostgreSQL pool, but it is not the final trust boundary.

## Non-Negotiable Boundaries

- Collectors do not receive unrestricted database access.
- Raw evidence is not stored in Git.
- Vendor-specific representations do not become the canonical model.
- External monitoring products do not become authorization or historical sources of truth.
- UI visibility is never treated as authorization.
- A requester cannot independently approve the requester’s own governed change.
- Parser uncertainty is visible and never silently converted into certainty.
- Generated diagrams never overwrite curated diagrams.
