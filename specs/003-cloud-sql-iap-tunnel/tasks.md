# Tasks: Cloud SQL Proxy and IAP Tunneling

**Input**: Design documents from `/specs/003-cloud-sql-iap-tunnel/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/
**Tests**: Tests are REQUIRED. Include unit tests (table-driven when applicable) and integration coverage for critical flows.
**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create required package skeleton and update module dependencies.

- [X] T001 Create `internal/proxy/proxy.go` with exported `Start(ctx context.Context, listener net.Listener, instance string, httpClient *http.Client) error` signature
- [X] T002 Create `internal/proxy/access.go` with exported `VerifyAccess(ctx context.Context, client *http.Client, instance string) error` signature
- [X] T003 Create `internal/proxy/instance.go` helper scaffold for parsing `project:region:instance` (stub implementation returning errors)
- [X] T004 Create `internal/proxy/errors.go` with typed errors + `UserFacingError(err error) string` stub (implementation deferred)
- [X] T005 Create `internal/proxy/dialer_plan.go` with a helper that computes whether IAM auth and Private IP must be used (stub)
- [X] T006 Update `go.mod` to add required modules for `cloud.google.com/go/cloudsqlconn` and `google.golang.org/api/sqladmin/v1beta4`, then run `go mod tidy`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core utilities required by all user stories (parsing, email resolution, error UX, dialer defaults).

⚠️ CRITICAL: No user story implementation should start until this phase is complete.

- [X] T007 Implement instance parsing in `internal/proxy/instance.go` with strict validation for `project:region:instance`
- [X] T008 Implement principal email resolution and Cloud SQL connector IAM token sources for connection instructions in `internal/auth/principal_email.go` and `internal/auth/cloudsql_token_sources.go` (new file) as `PrincipalEmail(ctx context.Context, httpClient *http.Client) (string, error)` and `CloudSQLConnectorTokenSources(ctx context.Context) (adminTokenSource, iamLoginTokenSource oauth2.TokenSource, err error)`, including updating the OAuth scope set used for Admin-token minting to satisfy `WithIAMAuthNTokenSources` requirements (align with `research.md` R-005).
- [X] T009 Implement user-facing proxy error formatting in `internal/proxy/errors.go` as `UserFacingError(err error) string`, including DevOps contact message for missing Cloud SQL/IAP roles, permission check unavailable guidance (FR-005a), and missing human-readable principal email guidance (FR-010a)
- [X] T010 Implement dialer defaults logic in `internal/proxy/dialer_plan.go` so the helper indicates IAM database auth is enabled by default and Private IP is forced by default
- [X] T011 Create instruction formatting helper in `internal/proxy/instructions.go` that formats stdout lines exactly as `Host: 127.0.0.1`, `Port: <port>`, `User: <email>`, and `Password: [LEAVE EMPTY]`

---

## Phase 3: User Story 1 - Start Secure Local Tunnel (Priority: P1) 🎯 MVP

**Goal**: Start a local TCP listener on `127.0.0.1:<port>`, accept connections, dial Cloud SQL using the connector, and relay traffic bidirectionally; print operator-ready connection instructions after successful initialization.

**Independent Test**: Run the proxy start logic in a unit/integration harness that uses injected fake dialer connections to verify:
  1) Accept loop handles a single local client,
  2) Relay forwards bytes in both directions,
  3) Stdout instruction block contains the required host/port/user/password lines.

### Tests for User Story 1 (REQUIRED) ⚠️

> NOTE: Write these tests first, ensure they fail before implementation by keeping stubs temporarily if needed.

- [X] T012 [US1] Create unit test `tests/unit/proxy_start_relay_test.go` that verifies bidirectional `io.Copy` relay for at least 2 concurrent accepted local connections using an injected fake dialer, including one failing dial path, and asserts the other active relay session continues (SC-004 containment).
- [X] T013 [US1] Create unit test `tests/unit/proxy_start_instructions_test.go` that verifies `internal/proxy/instructions.go` output formatting includes `Host`, `Port`, `User` (email), and `Password: [LEAVE EMPTY]`
- [X] T014 [US1] Create unit test `tests/unit/proxy_start_shutdown_test.go` that verifies the accept loop terminates on `ctx.Done()` / listener close without deadlocking

### Implementation for User Story 1

- [X] T015 [US1] Implement `Start` in `internal/proxy/proxy.go` to initialize `cloudsqlconn.NewDialer` with IAM auth (`cloudsqlconn.WithIAMAuthN()`), force Private IP (`cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP())`), and wire IAM-enabled token sources via `cloudsqlconn.WithIAMAuthNTokenSources(adminTokenSource, iamLoginTokenSource)`, ensuring initialization honors `ctx` cancellation.
- [X] T016 [US1] In `internal/proxy/proxy.go`, implement accept loop using `listener.Accept()` in a `for` loop and exit only when `ctx` is canceled or listener is closed; otherwise continue on accept errors
- [X] T017 [US1] In `internal/proxy/proxy.go`, implement per-connection goroutine handler that calls `dialer.Dial(ctx, instance)` and closes both local and remote connections when finished
- [X] T018 [US1] In `internal/proxy/proxy.go`, implement per-connection bidirectional relay that runs `io.Copy(local, remote)` in a goroutine and `io.Copy(remote, local)` in the handler goroutine, ensuring both directions completion leads to connection closure
- [X] T019 [US1] In `internal/proxy/proxy.go`, ensure one relay failure does not crash the proxy or stop accepting new connections (contain errors per-session)
- [X] T020 [US1] Implement principal email usage in `Start` by calling `internal/auth.PrincipalEmail(ctx, httpClient)` (from `T008`) and returning the missing-email error (FR-010a) if it cannot be resolved
- [X] T021 [US1] Print the final instruction block to stdout after successful dialer init in `internal/proxy/proxy.go` using `internal/proxy/instructions.go`
- [X] T022 [US1] Wire feature path in `cmd/sql-proxy/main.go` by, after `net.Listen` succeeds, obtaining authenticated HTTP client via `auth.GetClient(ctx)`, running `proxy.VerifyAccess(ctx, httpClient, settings.Instance)` before starting accept loop, printing `proxy.UserFacingError(err)` to stderr and returning non-zero on failure, and calling `proxy.Start(ctx, listener, settings.Instance, httpClient)` on success

---

## Phase 4: User Story 2 - Enforce IAM and Private Connectivity Defaults (Priority: P1)

**Goal**: Ensure IAM DB authentication and Private IP tunneling defaults are always used, and verify this behavior with unit tests focused on dialer initialization intent.

**Independent Test**: Unit test the dialer plan decision logic and the Start dialer initialization path using an injected dialer constructor to confirm IAM + Private IP are selected.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T023 [US2] Create unit test `tests/unit/dialer_plan_test.go` that asserts `internal/proxy/dialer_plan.go` indicates IAM auth ON and Private IP ON as default behavior
- [X] T024 [US2] Create unit test `tests/unit/proxy_dialer_init_defaults_test.go` that injects a fake `cloudsqlconn.NewDialer` wrapper and verifies Start requests IAM auth + Private IP (no public IP fallback)

### Implementation for User Story 2

- [X] T025 [US2] Refactor `internal/proxy/proxy.go` to centralize dialer initialization selection via `internal/proxy/dialer_plan.go` so the unit tests can validate behavior deterministically
- [X] T026 [US2] Add a small internal injection point in `internal/proxy/proxy.go` (e.g., a package-level var) so tests can replace the dialer constructor without impacting production behavior

---

## Phase 5: User Story 3 - Receive Actionable Permission Feedback (Priority: P2)

**Goal**: Implement pre-flight permission checks (`sqladmin.instances.get`) and map permission outcomes to specific, human-readable errors:
 - missing Cloud SQL/IAP roles -> DevOps contact guidance
 - permission check unavailable (API/network uncertainty) -> distinct unavailable guidance (FR-005a)
 - unknown/404-like instance issues -> clear instance-not-found guidance

**Independent Test**: Unit test VerifyAccess using an injected instances-get function or stubbed SQL Admin transport so you can validate mapping for:
 - 403 permission denial
 - 404 instance invalid/not accessible
 - transport error / timeout => permission check unavailable

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T027 [US3] Create unit test `tests/unit/access_verify_permissions_test.go` that verifies `VerifyAccess` returns an IAM/IAP-missing error when instances.get yields a 403
- [X] T028 [US3] Create unit test `tests/unit/access_verify_unknown_test.go` that verifies `VerifyAccess` returns a permission-check-unavailable error for transport errors / ambiguous failures
- [X] T029 [US3] Create unit test `tests/unit/access_verify_instance_not_found_test.go` that verifies `VerifyAccess` returns an instance-not-found style error for 404/invalid instance identifier
- [X] T030 [US3] Create unit test `tests/unit/proxy_user_facing_error_test.go` that verifies `internal/proxy/errors.go` produces role-specific DevOps messaging and does not include stack traces or secret material

### Implementation for User Story 3

- [X] T031 [US3] Implement `VerifyAccess` in `internal/proxy/access.go` to parse `project:region:instance` using `internal/proxy/instance.go` and call `sqladmin.NewService` + `Instances.Get` for the parsed instance
- [X] T032 [US3] In `internal/proxy/access.go`, classify outcomes into typed errors that `internal/proxy/errors.go` maps to user-facing messages: 403 => missing Cloud SQL/IAP roles (FR-005), transport/ambiguous failures => permission check unavailable (FR-005a), 404/invalid identifier => instance-not-found guidance
- [X] T033 [US3] In `internal/proxy/access.go`, add a test injection point (instances-get function wrapper) so unit tests can simulate outcomes without live cloud dependencies
- [X] T034 [US3] Ensure `cmd/sql-proxy/main.go` uses `proxy.UserFacingError(err)` for all `VerifyAccess` failures and exits non-zero

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories: robustness, security hardening, and shutdown cleanup verification.

- [X] T035 Add integration test `tests/integration/proxy_tunnel_shutdown_test.go` that starts the proxy with a short-lived context and verifies shutdown closes the listener and terminates relay goroutines
- [X] T036 Validate cancellation safety in `internal/proxy/proxy.go` by ensuring per-connection handlers exit promptly on `ctx.Done()` and do not leak goroutines
- [X] T037 Ensure user-facing error messages remain actionable and non-technical by adding/expanding tests in `tests/unit/proxy_user_facing_error_test.go` for representative error variants
- [X] T038 Update `specs/003-cloud-sql-iap-tunnel/quickstart.md` if needed to reflect the exact stdout format and failure classes used by this feature
- [X] T039 Run `gofmt` across all touched Go files and ensure `go test ./...` passes

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup (Phase 1): No dependencies
- Foundational (Phase 2): Depends on Setup completion; blocks all user stories
- User Stories (Phase 3+): Depends on Foundational completion
- Polish (Phase 6): Depends on all desired user stories being complete

### User Story Dependencies

- User Story 1 (US1): Can start after Phase 2 and delivers the working localhost relay with instructions.
- User Story 2 (US2): Can start after Phase 2 and hardens security defaults via tests.
- User Story 3 (US3): Can start after Phase 2 and delivers pre-flight permission UX mapping.
