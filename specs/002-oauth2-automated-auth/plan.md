# Implementation Plan: OAuth2 Automated Authentication

**Branch**: `002-oauth2-automated-auth` | **Date**: 2026-05-05 | **Spec**: `specs/002-oauth2-automated-auth/spec.md`
**Input**: Feature specification from `/specs/002-oauth2-automated-auth/spec.md`

## Summary

Introduce automated desktop OAuth2 authentication for the SQL proxy so users can authenticate without
external CLI dependencies. The feature covers secure token persistence, browser-driven consent,
callback handling with strict state validation, invalid-token recovery, and integration via an
authenticated HTTP client for Cloud SQL dialer initialization.

## Technical Context

**Language/Version**: Go (stable current release via Go modules)  
**Primary Dependencies**: `golang.org/x/oauth2`, `golang.org/x/oauth2/google`  
**Storage**: Local filesystem token store (`~/.sql-proxy/token.json`)  
**Testing**: Go `testing` package with table-driven unit tests and integration tests  
**Target Platform**: Linux, Windows, macOS  
**Project Type**: CLI application  
**Performance Goals**: Valid token startup path reuses credentials without browser flow in under 2 seconds  
**Constraints**: Token file permissions must be owner-only; callback state must be strictly validated;
local callback defaults to `8080` with dynamic localhost fallback if occupied  
**Scale/Scope**: Single-user local auth session flow per process execution

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Code structure uses `cmd/` and `internal/` with idiomatic Go conventions.
- [x] Error handling strategy uses wrapped errors and user-friendly CLI messages; no panic-based flow.
- [x] Concurrency design documents goroutine lifecycle, cancellation, and disconnection cleanup.
- [x] Testing plan includes table-driven unit tests and required integration coverage.
- [x] Security controls include IAM/private tunnel defaults and secure local token handling.
- [x] Dependency plan is minimal and justified (standard library first, then approved packages).

## Project Structure

### Documentation (this feature)

```text
specs/002-oauth2-automated-auth/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── auth-cli-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── sql-proxy/
    └── main.go

internal/
├── auth/
│   └── auth.go
└── config/
    └── config.go

tests/
├── integration/
│   └── oauth_flow_test.go
└── unit/
    └── auth_token_test.go
```

**Structure Decision**: Add an `internal/auth` package for OAuth2 flow and token management while
keeping startup orchestration in `cmd/sql-proxy/main.go`. This preserves modular responsibilities and
supports isolated unit testing for token logic plus integration testing for callback and browser flow.

## Phase 0: Research Plan

- Confirm OAuth2 desktop flow patterns for Go CLI apps with localhost callback servers.
- Confirm cross-platform browser-launch strategy and fallback behavior.
- Confirm secure token-file persistence and invalid-token recovery best practices.

## Phase 1: Design Outputs

- `data-model.md`: OAuth session, token record, callback session, and authenticated-client context model.
- `contracts/auth-cli-contract.md`: auth startup contract, callback behavior, error contract, and outputs.
- `quickstart.md`: local flow validation steps for first-run auth, token reuse, invalid-token recovery.

## Post-Design Constitution Re-Check

- [x] No constitutional violations introduced by design artifacts.
- [x] Testing and security requirements remain explicit and enforceable.

## Complexity Tracking

No constitution gate violations requiring justification.
