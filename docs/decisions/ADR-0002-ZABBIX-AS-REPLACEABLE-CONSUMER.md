# ADR-0002 — Zabbix as a Replaceable Consumer and Complementary Operational System

## Status

Accepted initial direction, clarified by the operational-system complement model.

## Decision

Retain Zabbix for mature continuous network monitoring, availability, performance history, triggers, graphing, alerting, escalation, maintenance, and proxy behavior.

Iron Atlas owns canonical infrastructure evidence, normalized identity, topology, governed findings, project, change, reconciliation, and acceptance records.

Atlas may deliver metrics and findings to Zabbix and may generate reviewable host reconciliation, templates, low-level discovery, maps, dashboards, maintenance recommendations, and report context through versioned adapters.

Atlas shall not represent generated Zabbix definitions as applied configuration. Any future Zabbix API write requires a separately accepted least-privileged, previewable, attributable, approval-aware, bounded, and validated provisioning boundary.

## Consequences

Atlas can add infrastructure context, preventive analysis, reconciliation, maps, dashboards, and reporting assistance without prematurely recreating Zabbix polling, trigger, graphing, escalation, maintenance, proxy, and availability-monitoring capabilities.

Zabbix remains replaceable as an integration destination because Atlas canonical records do not depend on the Zabbix schema.
