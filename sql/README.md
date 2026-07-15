# PostgreSQL Boundary

The SQL directory is a planned authoritative persistence boundary. The Phase 0 Go service deliberately uses memory storage so the API, HTML interface, authorization model, module contracts, and tests can be established without prematurely declaring a production database contract.

The initial migration is a design candidate and is not accepted for production. It demonstrates actor, role, governed change, independent approval, evidence reference, canonical telemetry, and external delivery-outbox concepts.
