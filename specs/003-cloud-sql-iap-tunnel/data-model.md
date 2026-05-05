# Data Model: Cloud SQL Proxy and IAP Tunneling

## Entities

### 1. `ProxySession`

Represents one running CLI process bridging localhost TCP to Cloud SQL until shutdown.

| Field | Description | Validation |
|-------|-------------|------------|
| `listenerAddr` | Host/port bound locally | MUST be `127.0.0.1:<port>` |
| `instanceCN` | Cloud SQL connection name `project:region:name` | Non-empty segments; canonical separator `:` |
| `principalEmail` | Operator-visible mailbox for IAM DB auth | MUST be non-empty before serving (FR-010a) |
| `dialer` | Cloud SQL connector abstraction | Initialized once per session; closed on shutdown |
| `ctx` | Shutdown/cancellation scope | Honors process signals / cancellation |

**State transitions**

`InitConfig → AuthReady → PreflightOK → Listening → (per client) Relaying → Shutdown`

Invalid transitions: `PreflightOK → Listening` forbidden if `principalEmail` unresolved; `Listening` forbidden if pre-flight indeterminate (FR-005a).

---

### 2. `PermissionProbeResult`

Structured outcome of startup checks before accept loop.

| Field | Description |
|-------|-------------|
| `status` | `allowed` / `denied` / `unknown` |
| `httpStatus` | Optional upstream HTTP code when applicable |
| `message` | Stable internal reason token (not necessarily user-visible raw text) |

Rules:

- `403` + permission semantics → `denied`.
- Missing instance / access not found → distinct `denied` variant (not identical copy to IAM denial).
- Network / 5xx / non-classified → `unknown` (must trigger FR-005a messaging path).

---

### 3. `ClientRelaySession`

One accepted local socket paired with one upstream `Dial` result.

| Field | Description |
|-------|-------------|
| `local` | Accepted loopback connection |
| `remote` | Cloud SQL connector `net.Conn` |
| `startedAt` | Observability hook (optional metric logging) |

Lifecycle:

1. Accept local.
2. Dial remote (honor context).
3. Start bidirectional `io.Copy`.
4. On any direction completion / error, close both sides idempotently.

Invariant: failures here MUST NOT abort unrelated relay sessions (FR-009 / SC-004).

---

### 4. `ConnectionInstructions`

Emitted to STDOUT after initialization succeeds.

| Field | Required | Notes |
|-------|----------|-------|
| Host | yes | literal `127.0.0.1` |
| Port | yes | Numeric active port from listener |
| Username | yes | Mirrors `principalEmail` |
| Password hint | yes | Exactly `[LEAVE EMPTY]` (spec FR-010 wording) |
