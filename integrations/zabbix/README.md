# Zabbix Adapter

The current code is a native Go implementation of the Zabbix sender packet format for delivery to pre-created trapper items.

The target Zabbix integration complements rather than replaces Zabbix. Future separately versioned work may add:

- device, interface, inventory, and monitoring reconciliation;
- template and low-level discovery recommendations;
- topology-derived maps;
- dashboard and report definitions;
- TLS;
- delivery outbox behavior;
- governed API provisioning; and
- post-provisioning readback and validation.

Generated recommendations and definitions remain distinct from applied Zabbix configuration. Production provisioning requires a separately accepted least-privileged and approval-aware boundary.
