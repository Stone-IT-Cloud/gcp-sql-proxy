# Tasks: Configuration and CLI Interface

**Input**: Design documents from `/specs/001-config-cli-interface/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are REQUIRED. Include table-driven unit tests and integration coverage for startup flows.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and baseline structure

- [x] T001 Initialize Go module metadata in `go.mod`
- [x] T002 Create project directories `cmd/sql-proxy/`, `internal/config/`, `tests/unit/`, and `tests/integration/`
- [x] T003 [P] Add dependency requirements for `github.com/spf13/viper` and `github.com/spf13/pflag` in `go.mod`
- [x] T004 [P] Create CLI contract baseline in `specs/001-config-cli-interface/contracts/cli-contract.md` for implementation traceability notes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core primitives required by all user stories

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T005 Implement configuration types and constants in `internal/config/config.go`
- [x] T006 Implement config directory resolution and creation logic in `internal/config/config.go`
- [x] T007 Implement Viper bootstrap (defaults + config path + read behavior) in `internal/config/config.go`
- [x] T008 Implement CLI flag registration and binding to Viper in `internal/config/config.go`
- [x] T009 Implement startup input validators (instance required, port range `1-65535`) in `internal/config/config.go`
- [x] T010 [P] Implement shared user-facing error formatting helpers in `internal/config/config.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Start Proxy with Predictable Configuration (Priority: P1) 🎯 MVP

**Goal**: Resolve startup configuration deterministically from flags, config file, and defaults.

**Independent Test**: Running the CLI with combinations of missing config, config values, and flags must produce expected effective values and reject missing/invalid inputs.

### Tests for User Story 1 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T011 [P] [US1] Add table-driven unit tests for precedence resolution in `tests/unit/config_init_test.go`
- [x] T012 [P] [US1] Add table-driven unit tests for missing instance and invalid port validation in `tests/unit/validation_test.go`
- [x] T013 [US1] Add integration startup test for defaults and config-file precedence in `tests/integration/startup_flow_test.go`

### Implementation for User Story 1

- [x] T014 [US1] Implement `Init()` orchestration for config loading and validation in `internal/config/config.go`
- [x] T015 [US1] Implement missing-config non-fatal behavior and malformed-config fatal behavior in `internal/config/config.go`
- [x] T016 [US1] Create CLI entrypoint skeleton that calls configuration initialization in `cmd/sql-proxy/main.go`
- [x] T017 [US1] Wire resolved configuration values (port and instance) into startup path in `cmd/sql-proxy/main.go`

**Checkpoint**: User Story 1 is fully functional and testable independently

---

## Phase 4: User Story 2 - Recover from Port Conflict Quickly (Priority: P2)

**Goal**: Detect local port conflicts early and return actionable guidance without panic.

**Independent Test**: When target port is pre-occupied, process must fail with exit code `1` and clear remediation message.

### Tests for User Story 2 (REQUIRED) ⚠️

- [x] T018 [P] [US2] Add unit tests for port conflict error messaging in `tests/unit/validation_test.go`
- [x] T019 [US2] Add integration test for bind failure and exit behavior in `tests/integration/startup_flow_test.go`

### Implementation for User Story 2

- [x] T020 [US2] Implement preflight listener bind check on `127.0.0.1:<port>` in `cmd/sql-proxy/main.go`
- [x] T021 [US2] Implement user-friendly port-conflict output with flag/config remediation guidance in `cmd/sql-proxy/main.go`
- [x] T022 [US2] Implement explicit exit code `1` path for bind conflicts in `cmd/sql-proxy/main.go`

**Checkpoint**: User Stories 1 and 2 work independently and together

---

## Phase 5: User Story 3 - Stop Cleanly on Termination Signals (Priority: P3)

**Goal**: Shut down listener and cancel runtime context cleanly on `SIGINT`/`SIGTERM`.

**Independent Test**: Running process must release listener and exit cleanly when sent termination signals.

### Tests for User Story 3 (REQUIRED) ⚠️

- [x] T023 [P] [US3] Add unit tests for signal-handling shutdown coordination in `tests/unit/validation_test.go`
- [x] T024 [US3] Add integration test for graceful signal-triggered shutdown in `tests/integration/startup_flow_test.go`

### Implementation for User Story 3

- [x] T025 [US3] Implement signal channel setup and notification wiring in `cmd/sql-proxy/main.go`
- [x] T026 [US3] Implement shutdown goroutine to close listener and cancel context in `cmd/sql-proxy/main.go`
- [x] T027 [US3] Implement clean termination messaging and shutdown sequencing in `cmd/sql-proxy/main.go`

**Checkpoint**: All user stories are independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Harden feature quality across all stories

- [x] T028 [P] Update quickstart validation notes for final CLI behavior in `specs/001-config-cli-interface/quickstart.md`
- [x] T029 [P] Add/refresh documentation comments for exported configuration APIs in `internal/config/config.go`
- [x] T030 Run full test suite and capture results for feature sign-off using `go test ./...`
- [x] T031 Validate cancellation and goroutine cleanup behavior under shutdown/retry scenarios in `tests/integration/startup_flow_test.go`
- [x] T032 Validate user-facing startup errors for clarity and consistency in `tests/integration/startup_flow_test.go`
- [x] T033 Add regression checks that startup/config changes do not alter IAM auth and private tunnel defaults in `tests/integration/startup_flow_test.go`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; starts immediately
- **Foundational (Phase 2)**: Depends on Setup completion; blocks all user stories
- **User Stories (Phase 3-5)**: Depend on Foundational completion
  - US1 is MVP and should be delivered first
  - US2 depends on US1 startup baseline
  - US3 depends on active listener lifecycle from US1/US2
- **Polish (Phase 6)**: Depends on completion of selected user stories

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2; no dependency on other stories
- **US2 (P2)**: Starts after US1 startup path is in place
- **US3 (P3)**: Starts after listener lifecycle is available (US1 + US2 paths)

### Within Each User Story

- Tests MUST be written and fail before implementation
- Configuration and validation logic before entrypoint integration
- Entry point startup before conflict handling
- Conflict handling before shutdown orchestration

### Parallel Opportunities

- Setup tasks marked `[P]` run in parallel
- Foundational tasks T010 can run in parallel with T007/T008 after core constants exist
- US1 tests T011 and T012 run in parallel
- US2 test and messaging work can parallelize after bind check is drafted
- US3 test and shutdown implementation can parallelize after signal wiring exists

---

## Parallel Example: User Story 1

```bash
# Run US1 test tasks together:
Task: "T011 Add table-driven unit tests for precedence resolution in tests/unit/config_init_test.go"
Task: "T012 Add table-driven unit tests for missing instance and invalid port validation in tests/unit/validation_test.go"

# Implement config wiring once tests are in place:
Task: "T014 Implement Init() orchestration for config loading and validation in internal/config/config.go"
Task: "T016 Create CLI entrypoint skeleton that calls configuration initialization in cmd/sql-proxy/main.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 and Phase 2
2. Complete Phase 3 (US1)
3. Validate precedence, defaults, and input validation behavior
4. Demo MVP startup flow before moving to conflict/shutdown enhancements

### Incremental Delivery

1. Deliver US1 (predictable configuration)
2. Deliver US2 (port conflict recovery)
3. Deliver US3 (graceful shutdown)
4. Execute Phase 6 polish and full regression suite

### Parallel Team Strategy

With multiple developers:

1. Developer A: Foundation + US1 config flow
2. Developer B: US2 conflict handling after US1 skeleton merges
3. Developer C: US3 shutdown handling after listener lifecycle is available
4. Joint polish pass for docs and integration validation

---

## Notes

- All tasks follow the required checklist format with IDs, optional `[P]`, `[USx]` labels, and file paths.
- Keep commit scope aligned to task boundaries for clean traceability.
- Stop at each story checkpoint and validate independently before proceeding.
