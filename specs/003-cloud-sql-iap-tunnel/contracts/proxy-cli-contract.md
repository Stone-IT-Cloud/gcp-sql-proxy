# Contract: Cloud SQL Proxy Tunnel CLI

## Purpose

Define operator-visible behavior for bridging `127.0.0.1:<port>` to Cloud SQL PostgreSQL via the Cloud SQL Go connector with IAM authentication and private IP defaults.

## Preconditions

1. Configuration resolves `--instance` / config file to canonical `project:region:instance`.
2. OAuth Feature 02 artifacts yield a refreshable user credential with scopes sufficient for SQL Admin + IAM login token minting (see `research.md`).
3. Local port is available on loopback (existing port conflict messaging remains).

## Startup Sequence (normative)

| Step | Behavior |
|------|----------|
| S1 | Parse settings (`internal/config`). |
| S2 | Acquire authenticated credential via `auth.GetClient` — **hard failure** if interactive / persisted auth cannot proceed (this feature requires an authenticated user). |
| S3 | Resolve **`principalEmail`** — **hard failure** before listening if missing (FR-010a). |
| S4 | Bind TCP listener on `127.0.0.1:<activePort>`. |
| S5 | Run `VerifyAccess` (`sqladmin.instances.get` semantics) — failures classified per FR-005 / FR-005a. |
| S6 | Construct Cloud SQL dialer with IAM auth + private IP defaults (see plan + research). |
| S7 | Print `ConnectionInstructions` block to STDOUT. |
| S8 | Enter accept loop until context cancellation or fatal listener error. |

## Runtime Contract

### Accept / Dial

- Each `Accept` spawns an isolated handler (goroutine) that:
  - Calls `Dial` with the process context (or derived timeout context if added during implementation — must honor parent cancellation).
  - Performs bidirectional copy using `io.Copy` pattern (one direction may run in helper goroutine).
  - Always closes both connections on completion.

### Shutdown

- Signal handling closes listener; in-flight relays observe closed connections and exit without panicking.
- Dialer `Close` invoked during graceful shutdown path (best-effort after accept loop exits).

## Outputs

### Successful instructions (STDOUT)

Minimum required lines (order flexible but all must appear):

```
Host: 127.0.0.1
Port: <port>
User: <principalEmail>
Password: [LEAVE EMPTY]
```

(No sensitive tokens, refresh material, or JWTs.)

### Errors (STDERR)

- Missing / failed auth: actionable, single-paragraph guidance.  
- Permission denied (403-class): include DevOps escalation AND mention likely roles `roles/cloudsql.client` + `roles/iap.tunnelResourceAccessor` without pretending definitive diagnosis when uncertain.  
- Permission **check unavailable**: distinct wording from denial — instruct retry / VPN / outage suspicion.  
- Forbidden: raw HTTP payloads, stack traces.

## Compatibility Notes

- Instance identifier MUST remain `project:region:instance` matching connector expectations.
- Future flags (explicit DB name, SSL mode) are **out of scope** unless spec amended.
