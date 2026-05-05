# Implementation Plan: Configuration and CLI Interface

**Branch**: `001-config-cli-interface` | **Date**: 2026-05-05 | **Spec**: `specs/001-config-cli-interface/spec.md`
**Input**: Feature specification from `/specs/001-config-cli-interface/spec.md`

## Summary

Implement deterministic startup configuration for the SQL proxy CLI using defaults, a user config
file, and command-line flags with strict precedence and validation. Startup must fail gracefully for
malformed config, missing required instance input, invalid ports, and occupied local ports while
supporting clean signal-based shutdown.

## Technical Context

**Language/Version**: Go (stable current release via Go modules)
**Primary Dependencies**: `github.com/spf13/viper`, `github.com/spf13/pflag`
**Storage**: Local filesystem (`~/.sql-proxy/config.yaml`)
**Testing**: Go `testing` package with table-driven unit tests + integration tests for startup flows
**Target Platform**: Windows, macOS, Linux
**Project Type**: CLI application
**Performance Goals**: Validate startup inputs and detect port conflicts in under 1 second locally
**Constraints**: Must bind to `127.0.0.1:<port>`; valid port range `1-65535`; no panic-based flow;
must preserve existing IAM authentication and private tunnel defaults
**Scale/Scope**: Single-process local CLI runtime; one listener instance per execution

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
specs/001-config-cli-interface/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── cli-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
cmd/
└── sql-proxy/
    └── main.go

internal/
└── config/
    └── config.go

tests/
├── integration/
│   └── startup_flow_test.go
└── unit/
    ├── config_init_test.go
    └── validation_test.go
```

**Structure Decision**: Use a standard Go CLI layout with `cmd/sql-proxy` for entrypoint orchestration
and `internal/config` for configuration lifecycle and precedence logic. Tests are split into unit and
integration suites to validate both deterministic config resolution and runtime startup behavior.

## Phase 0: Research Plan

- Confirm Viper + pflag precedence and binding strategy for flag > config > defaults.
- Confirm cross-platform behavior for creating `~/.sql-proxy` and handling permissions.
- Confirm robust signal handling and listener shutdown patterns for Go CLI processes.

## Phase 1: Design Outputs

- `data-model.md`: Runtime configuration, validation constraints, and listener lifecycle model.
- `contracts/cli-contract.md`: CLI flags, config schema, error contract, and exit-code contract.
- `quickstart.md`: Setup, run commands, failure-mode checks, and shutdown validation flow.

## Post-Design Constitution Re-Check

- [x] No constitutional violations introduced by design artifacts.
- [x] Testing and security expectations remain explicit and measurable.

## Complexity Tracking

No constitution gate violations requiring justification.
