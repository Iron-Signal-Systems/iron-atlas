# Zabbix Integration Contract

## Purpose

Deliver selected Atlas health, collection, finding, and freshness indicators to the existing Zabbix environment without making Zabbix the canonical Atlas database.

## Initial Sender Metrics

Examples:

```text
atlas.service.health
atlas.collection.last_success_unixtime
atlas.collection.overdue
atlas.collection.partial
atlas.device.findings.critical
atlas.device.findings.high
atlas.device.findings.moderate
atlas.device.config_drift
atlas.device.evidence_age_seconds
atlas.integration.delivery_failures
```

## Ownership

Atlas owns the source event, context, classification, finding state, and delivery history. Zabbix owns its item history, triggers, notifications, dashboards, maintenance, and escalation behavior.

## Delivery

The initial Go adapter implements the Zabbix sender protocol for pre-created trapper items. Later work may add TLS profiles, low-level discovery, template provisioning, and API-based lifecycle management through a separate versioned adapter.

## Failure

Failed delivery enters bounded retry and eventually dead-letter state. It does not discard the canonical event or block collection and analysis.

## Replacement

The same canonical event can be delivered to OpenMetrics, syslog, webhooks, or SIEM destinations without changing the source schema.
