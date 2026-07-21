# Monitoring, Alerting, and Evidence-Freshness Model

## Status

Normative architecture alignment candidate.

## Purpose

Atlas complements monitoring platforms rather than replacing their alerting, graphing, escalation, and high-frequency telemetry. Atlas must still report whether its own services and evidence are trustworthy, current, degraded, or blocked.

## Independent dimensions

1. **Service health** — whether the process or component functions.
2. **Service readiness** — whether a declared operation can run safely.
3. **Evidence freshness** — whether accepted evidence is recent enough for a source, record class, and question.

A healthy process may hold stale evidence. A ready API may be unable to answer one question with current evidence.

## Freshness policy

A source- and class-specific policy declares expected interval, warning age, stale age, unusable age when applicable, clock-skew limit, required record classes, partial-evidence handling, and escalation route. One global threshold for all evidence is prohibited.

## Freshness states

- `CURRENT`
- `AGING`
- `STALE`
- `UNUSABLE`
- `MISSING`
- `PARTIAL`
- `CONFLICTING`
- `UNKNOWN`

## Answer propagation

Every evidence-dependent answer exposes accepted snapshot identity, source, record class, collection and observation times, freshness, missing or failed required records, conflicts, inference versus direct evidence, and the effect on confidence. Silence is never proof that a route, policy, identity path, dependency, or exposure does not exist.

## Monitoring integrations

Zabbix, Graylog, Security Onion, and vendor systems remain authoritative for their accepted functions. Atlas may expose scheduler health, collection duration, queue saturation, candidate outcome, accepted-snapshot age, parser failures, dependency outages, database readiness, and resource observations.

## Alerting boundary

Atlas alerts on Atlas-owned conditions or accepted derived findings. Alerts are deduplicated, rate limited, severity governed, attributable to evidence, maintenance-aware, explicit about stale or incomplete inputs, and separate from acceptance of the underlying evidence.

## Degraded operation

Atlas may remain available when one source or adapter is unavailable, a candidate fails while prior accepted evidence remains valid, or a nonessential integration fails. The limitation remains visible in status and answers.

## Required validation

Test freshness boundaries, clock skew, missing timestamps, stale snapshots, partial bundles, source disagreement, integration outage, duplicate alerts, maintenance windows, alert storms, accepted-state preservation, and query propagation.
