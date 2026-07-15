# External-System-Independent Telemetry

## Decision

Iron Atlas owns canonical health, metric, finding, delivery-intent, and operational records independently of Zabbix or any other monitoring product.

Zabbix is an important initial delivery destination and may remain the primary infrastructure-monitoring interface. It is not Atlas’s authorization source, historical source of truth, or canonical schema.

## Initial Adapters

- Native Go Zabbix sender-protocol adapter
- Future Zabbix API adapter for governed provisioning where justified
- OpenMetrics-compatible endpoint
- Syslog
- Webhook
- SIEM delivery

## Delivery Contract

Each delivery contains:

- Integration contract version
- Destination ID
- Canonical metric or event ID
- Device, site, module, and finding context
- Value and timestamp
- Classification
- Attempt count and next attempt
- Delivery result

Retries are bounded with backpressure and dead-letter handling. Destination failure does not block collection, parsing, change management, or canonical recording.

## Zabbix Direction

The first implementation sends values to existing Zabbix trapper items using the sender protocol directly from Go, avoiding a required `zabbix_sender` package. TLS support, host/item provisioning, templates, low-level discovery, and long-term ownership boundaries remain later implementation work.

Atlas may eventually provide deeper configuration and evidence context than Zabbix, but it should not attempt to replace mature continuous availability polling, alerting, graphing, and escalation until those capabilities are deliberately implemented and proven.
