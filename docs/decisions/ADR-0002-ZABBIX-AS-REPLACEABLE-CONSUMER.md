# ADR-0002 — Zabbix as a Replaceable Consumer

## Status

Accepted initial direction.

## Decision

Retain Zabbix for mature continuous network monitoring and deliver relevant Atlas metrics and findings through a versioned adapter. Iron Atlas owns canonical evidence, findings, project, change, and acceptance records.

## Consequences

Atlas can add context and preventive analysis without prematurely recreating Zabbix polling, trigger, graphing, escalation, maintenance, proxy, and availability-monitoring capabilities.
