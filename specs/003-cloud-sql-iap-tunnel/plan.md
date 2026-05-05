# Implementation Plan: Cloud SQL Proxy and IAP Tunneling

**Branch**: `003-cloud-sql-iap-tunnel` | **Date**: 2026-05-05 | **Spec**: `specs/003-cloud-sql-iap-tunnel/spec.md`  
**Input**: Feature specification from `/specs/003-cloud-sql-iap-tunnel/spec.md` (plus implementation outline from `/speckit-plan` command)

## Summary

Deliver a localhost TCP bridge to Cloud SQL for PostgreSQL that uses the authenticated Google identity from Feature 02, initializes the official Cloud SQL Go connector with IAM database authentication and private IP as defaults, performs fast pre-flight access checks (`sqladmin.instances.get`) before accepting clients, relays bytes bidirectionally per connection under full `context.Context` cancellation, prints operator-ready SQL client instructions including the resolved principal email with empty-password guidance (per clarifications FR-005a / FR-010a).

## Technical Context

**Language/Version**: Go 1.25.x (repo `go.mod`)  
**Primary Dependencies**:  
- `cloud.google.com/go/cloudsqlconn` — Cloud SQL Go connector dialer (`NewDialer`, `Dial`)  
- `google.golang.org/api/sqladmin/v1beta4` — SQL Admin REST client for pre-flight (`Instances.Get`)  
- `google.golang.org/api/googleapi` — Detect HTTP status for permission/classification UX  
Existing: `internal/auth`, `internal/config`, `golang.org/x/oauth2`, `spf13/viper` / `pflag`

**Storage**: Configuration and tokens remain in Feature 02 paths (`~/.sql-proxy/`); tunnel feature does not add new persisted stores.

**Testing**: Go `testing`; table-driven units for parse/error mapping; integration tests mocked via `httptest` / fake SQL Admin handlers or build tags where cloud calls are unavoidable.

**Target Platform**: Linux, Windows, macOS (loopback TCP only).

**Project Type**: CLI application (`cmd/sql-proxy`).

**Performance Goals**: Pre-flight succeeds or fails distinctly before arbitrary client-visible dial timeouts align with SC-001 (≤ ~5 s for happy path startup in well-connected environments); per-connection isolation per SC-004.

**Constraints**:  
- MUST use IAM DB auth by default (`WithIAMAuthN()`).  
- MUST use private IP for dials (`WithDefaultDialOptions(WithPrivateIP())`).  
- MUST pass OAuth-derived credentials into connector per Feature 02 (connector requires explicit token sourcing when IAM auth is enabled — see research.md).  
- MUST fail startup if authenticated principal email cannot be resolved (spec FR-010a).  
- MUST fail startup when permission pre-flight outcome is unknown (network/API ambiguity) — spec FR-005a.  
- No stack traces or secret material in STDERR/STDOUT messages.

**Scale/Scope**: Single-process, multi-connection local relay; IAM / Cloud SQL quotas apply upstream.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Code structure uses `cmd/` and `internal/` with idiomatic Go conventions.
- [x] Error handling strategy uses wrapped errors and user-friendly CLI messages; no panic-based flow.
- [x] Concurrency design documents goroutine lifecycle, cancellation, and disconnection cleanup.
- [x] Testing plan includes table-driven unit tests and required integration coverage.
- [x] Security controls include IAM/private tunnel defaults and secure local token handling.
- [x] Dependency plan is minimal and justified (official Google modules + existing auth/config).

## Project Structure

### Documentation (this feature)

```text
specs/003-cloud-sql-iap-tunnel/
├── plan.md              # This file
├── research.md          # Phase 0
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1
├── contracts/
│   └── proxy-cli-contract.md
└── tasks.md             # Phase 2 (/speckit-tasks — not produced here)
```

### Source Code (repository root)

```text
cmd/
└── sql-proxy/
    └── main.go

internal/
├── auth/
│   └── auth.go                     # Extend: principal email + token sources for connector
├── config/
│   └── config.go
└── proxy/
    ├── proxy.go                     # Dialer lifecycle, Accept loop, relay goroutines
    └── access.go                    # VerifyAccess pre-flight helpers (optional split)

tests/
├── integration/
│   └── proxy_tunnel_test.go          # Relay / shutdown harnesses (expand as needed)
└── unit/
    ├── proxy_relay_test.go           # Parsing, error UX mapping
    └── access_verify_test.go         # Permission error classification
```

**Structure Decision**: Add `internal/proxy` owning dialer initialization, Accept loop, and bidirectional relays; extend `internal/auth` with email resolution + IAM-compatible token sourcing for Cloud SQL connector; keep thin orchestration in `cmd/sql-proxy/main.go`.

## Complexity Tracking

No constitution gate violations requiring justification.

## Phase 0: Research Plan

- Confirm connector option composition for IAM auth **with OAuth user credentials** (`WithIAMAuthNTokenSources` vs credentials-only helpers) and interplay with SQL Admin HTTP client configuration.
- Map `sqladmin.instances.get` outcomes (403, 404, 5xx, transport errors) to spec FR-005 / FR-005a user-facing UX.
- Document PostgreSQL + IAM username expectations for printed instructions versus connector authentication concerns.

**Output**: `research.md` (complete — no NEEDS CLARIFICATION remain before implementation).

## Phase 1: Design Outputs

- `data-model.md`: Proxy session, permission check, relay session, connection instructions.
- `contracts/proxy-cli-contract.md`: Startup ordering, errors, STDOUT instruction format, shutdown semantics.
- `quickstart.md`: Local validation flow and operator expectations.

## Implementation Outline (from command input)

1. **Initialize proxy package**  
   - Create `internal/proxy/` with `proxy.go` (and optional `access.go` for `VerifyAccess` cleanness).

2. **Pre-flight check**  
   - `VerifyAccess(ctx context.Context, client *http.Client, instance string) error`  
   - Parse `project:region:name` from `instance`; call SQL Admin `Instances.Get(project, name)` using a service constructed with `option.WithHTTPClient(client)` (or shared transport consistent with auth).  
   - Classify `403` as missing Cloud SQL / IAM access — message MUST mention DevOps escalation and reference likely roles (`roles/cloudsql.client`, `roles/iap.tunnelResourceAccessor`) without leaking internal URLs.  
   - Surfaces that are **not** definitive permission denial (network reset, DNS, 5xx, unexpected transport) → fail with **permission check unavailable** (FR-005a).

3. **Dialer initialization**  
   - Exported `Start(ctx context.Context, listener net.Listener, instance string, httpClient *http.Client) error` (signature may grow with explicit `oauth2.TokenSource` pairs — see `research.md` R-001 / R-005).  
   - `cloudsqlconn.NewDialer(ctx, …)` with:  
     - `cloudsqlconn.WithIAMAuthN()`  
     - `cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP())`  
     - `cloudsqlconn.WithIAMAuthNTokenSources(adminTokenSource, iamLoginTokenSource)` (avoid relying on `httpClient` alone — it does not satisfy IAM login token requirements).  
     - Optional `cloudsqlconn.WithHTTPClient(httpClient)` **only** when it matches the same OAuth transport and does not duplicate/conflict with token source configuration.  
   - After successful init, print connection instructions (host `127.0.0.1`, active port, email, `[LEAVE EMPTY]` password line).

4. **Connection loop**  
   - `for { conn, err := listener.Accept(); … }` with `ctx` cancellation closing listener (existing `main` pattern).  
   - Per connection: `dialer.Dial(ctx, instance)` (or child context with timeout if needed — document in implementation), `defer` close both conns, `go io.Copy(local, remote)` + `io.Copy(remote, local)` in handler goroutine, ensure errors do not kill process (log to STDERR concisely).

5. **Wire `cmd/sql-proxy/main.go`**  
   - After `config.Init()` and successful `net.Listen("tcp", "127.0.0.1:port")`:  
     - Require successful `auth.GetClient(ctx)` for this feature path (missing OAuth client credentials or token failure → user-facing fatal error, not silent continue).  
     - Resolve user email for instructions; fail before serving if missing (FR-010a).  
     - Call `proxy.VerifyAccess` then `proxy.Start`.  
   - Preserve signal-driven shutdown: closing listener unblocks Accept; ensure dialer `Close` on exit.

## Post-Design Constitution Re-Check

- [x] IAM + private IP defaults explicit in design and contract.  
- [x] Context cancellation and per-connection isolation captured.  
- [x] Error UX remains non-technical and avoids secret leakage.  
- [x] Test strategy covers permission mapping and relay lifecycle without mandating live cloud in default unit runs.
