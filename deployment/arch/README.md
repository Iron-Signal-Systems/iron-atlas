# Arch Linux Deployment Baseline

The initial deployment target is a minimal Arch Linux system managed by systemd.

## Production Runtime

The production host should contain only packages required for the operating system, PostgreSQL, TLS trust, controlled remote administration, and Iron Atlas binaries. Go and Python are build or administrative dependencies and need not remain on a production appliance when signed binaries and validation tooling are delivered separately.

Rust, Node.js, npm, and a separate Zabbix sender package are not required by the initial design.

## Filesystem Boundary

```text
/etc/iron-atlas/                 root-owned configuration
/var/lib/iron-atlas/             service state
/var/lib/iron-atlas/evidence/    encrypted raw evidence boundary
/var/log/iron-atlas/             local operational logs when journald is insufficient
/usr/local/bin/iron-atlas        service binary
```

Raw evidence should preferably use a dedicated encrypted filesystem or protected remote object store. It must not share the Git working tree.

## Service Identity

Run Iron Atlas under a dedicated locked system account with no interactive shell. Collectors should use separate service identities and separate secrets from the web/API service.
