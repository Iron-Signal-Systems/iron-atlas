# Target Architecture

## Status

Normative target direction. Phase 1 Steps 1 and 2 are accepted. The Phase 1 Step 3 authentication foundation, governed actor resolver, and bounded OIDC discovery, JWKS, and ID-token verification checkpoints are integrated. Authorization-code exchange, PKCE transaction persistence, sessions, CSRF, trusted-proxy enforcement, and production authentication are not accepted.

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

Replaceable external adapters operate in both governed directions:

```text
Approved external evidence and metadata
        ├── Cisco and Fortinet evidence
        ├── Zabbix inventory and monitoring metadata
        ├── Graylog syslog and SNMP-trap context
        ├── Security-platform asset references
        └── Curated documentation and Draw.io sources
                         ↓
              Atlas evidence and adapter boundary
                         ↓
       Canonical normalized records and governed findings
                         ↓
Atlas context, outbox, recommendations, and generated artifacts
        ├── Zabbix metrics, reconciliation, maps, dashboards,
        │   template and discovery suggestions, and report context
        ├── Graylog lookup data, queries, pipelines, streams,
        │   dashboards, and report definitions
        ├── Security-platform asset and topology context
        ├── Draw.io-compatible generated topology sources
        ├── OpenMetrics
        ├── Syslog and webhooks
        └── Future versioned adapters
```

External systems retain responsibility for their mature operational functions. Atlas retains responsibility for canonical evidence, normalized infrastructure identity, topology, governed findings, change history, and formal acceptance history.

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
- Atlas complements mature monitoring, logging, security, diagramming, and vendor systems rather than recreating them without evidence of need.
- External-system records do not silently become authoritative Atlas state.
- Future external-system writes or provisioning require a separately accepted, previewable, attributable, bounded, and validated integration boundary.
- UI visibility is never treated as authorization.
- Identity-provider claims never directly become Atlas authority; production requests resolve through governed Atlas actors and role bindings.
- Request-controlled bodies, forms, queries, paths, and ordinary headers never select the production actor.
- A requester cannot independently approve the requester’s own governed change.
- Parser uncertainty is visible and never silently converted into certainty.
- Generated diagrams never overwrite curated diagrams.
