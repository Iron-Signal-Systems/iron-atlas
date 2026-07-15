# Arch Linux Deployment Baseline

The initial deployment target is a minimal Arch Linux system managed by systemd.

## Production Runtime

The production host should contain only packages required for the operating system, PostgreSQL client connectivity where locally required, TLS trust, controlled remote administration, and Iron Atlas binaries. Go and Python are build or validation dependencies and need not remain on a production appliance when signed binaries and validation tooling are delivered separately.

Rust, Node.js, npm, and a separate Zabbix sender package are not required by the initial design.

## Filesystem Boundary

```text
/etc/iron-atlas/                 root-owned configuration and protected environment
/var/lib/iron-atlas/             service state
/var/lib/iron-atlas/evidence/    encrypted raw evidence boundary
/var/log/iron-atlas/             local operational logs when journald is insufficient
/usr/local/bin/iron-atlas        service binary
```

Raw evidence should preferably use a dedicated encrypted filesystem or protected remote object store. It must not share the Git working tree.

## PostgreSQL Runtime Configuration

The Step 2 service candidate reads:

- `IRON_ATLAS_CHANGE_STORE=postgresql`
- `IRON_ATLAS_DATABASE_URL`
- `IRON_ATLAS_DATABASE_MAX_CONNECTIONS`
- `IRON_ATLAS_DATABASE_MIN_CONNECTIONS`

The systemd candidate reads `/etc/iron-atlas/iron-atlas.env` when present. That file should be root-owned, mode `0600`, excluded from Git, and populated through the deployment secret process. Production credential delivery, certificate provisioning, and rotation are not accepted in Step 2.

## Service Identity

Run Iron Atlas under a dedicated locked system account with no interactive shell. Collectors should use separate service identities and separate secrets from the web/API service.
