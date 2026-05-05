# Quickstart: Cloud SQL Proxy and IAP Tunneling

## Prerequisites

- Feature 02 OAuth completes once so `~/.sql-proxy/token.json` exists **or** OAuth client env / build tags inject valid desktop client credentials.
- GCP roles (representative — final messaging implemented in code/tests):
  - `roles/cloudsql.client`
  - `roles/iap.tunnelResourceAccessor` when organization requires IAP-aware tunnel access
  - Ability to read instance metadata (`cloudsql.instances.get` / equivalent policy)
- Cloud SQL PostgreSQL instance reachable over **private networking** from Google’s connector infrastructure for your environment.

## Run the proxy

```bash
go run ./cmd/sql-proxy --instance PROJECT:REGION:INSTANCE --port 55432
```

Expect STDOUT instructions:

```
Host: 127.0.0.1
Port: 55432
User: you@example.com
Password: [LEAVE EMPTY]
```

Connect with `psql`:

```bash
psql "host=127.0.0.1 port=55432 dbname=postgres user=you@example.com sslmode=disable"
```

> **Note**: Exact `sslmode` requirements follow PostgreSQL client defaults for loopback; tighten if your client mandates TLS locally.

## Negative testing (operator view)

1. Revoke `cloudsql.client` temporarily → command MUST fail during pre-flight with DevOps guidance (not a hang).  
2. Disconnect network → expect **permission check unavailable** class message (FR-005a).  
3. Omit OAuth env / token → expect authentication failure before relay loop.

## Automated tests

```bash
go test ./...
```

Unit tests SHOULD cover instance string parsing, HTTP status → UX mapping, and relay lifecycle without live GCP.
