# Minimal Arch Linux Deployment

## Runtime Goal

Run on a base Arch Linux system with only the packages necessary for the service, protected storage, PostgreSQL, TLS, systemd, and controlled administration.

## Initial Package Direction

Production runtime:

- Base Arch system
- Linux kernel and firmware
- systemd
- CA certificates
- OpenSSH
- PostgreSQL
- Signed Iron Atlas Go binaries

Development and administration:

- Go
- Python 3
- Git
- Vim
- tmux

Rust is not required initially. Node.js and npm are not required because the HTML5 interface is embedded and rendered by Go.

## Hardening

- Dedicated service and collector accounts
- No interactive shell for service identities
- systemd sandboxing
- Read-only system paths
- Explicit writable state paths
- Host firewall and management-network restriction
- TLS termination and certificate lifecycle
- Off-host logs and integrity records
- Protected backups and restore tests
- Disabled-at-rest break-glass process
- Reproducible build and signed artifact verification
- Package and `/etc` integrity baselines
- Trusted rebuild and compromise recovery
