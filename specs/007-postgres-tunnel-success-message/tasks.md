# Tasks: PostgreSQL Tunnel Success Message

**Input**: Design documents from `/specs/007-postgres-tunnel-success-message/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED. Include unit tests (table-driven when applicable) and integration
coverage for critical flows.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and baseline scaffolding for tunnel messaging feature delivery.

- [X] T001 Create feature test fixture for successful tunnel-ready output in `tests/unit/testdata/tunnel_success_message_ready.txt`
- [X] T002 Create feature test fixture for reconnect success output in `tests/unit/testdata/tunnel_success_message_reconnect.txt`
- [X] T003 [P] Create feature test fixture for startup failure output in `tests/unit/testdata/tunnel_success_message_failure.txt`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared output primitives and event hooks required before implementing user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Add reusable success-message payload formatter for tunnel-ready events in `internal/proxy/tunnel_success_message.go`
- [X] T005 [P] Add reusable PostgreSQL command-template builder with placeholder enforcement in `internal/proxy/postgres_client_template.go`
- [X] T006 [P] Add startup/reconnect event wiring seam for emitting user-facing success blocks in `internal/proxy/proxy_runtime.go`
- [X] T007 Add foundational unit tests for formatter/template helpers in `tests/unit/tunnel_success_message_format_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - Confirm tunnel readiness (Priority: P1) 🎯 MVP

**Goal**: Show a clear success message only when tunnel-ready state is actually reached.

**Independent Test**: Start a successful tunnel and verify readiness confirmation is emitted exactly when the ready event occurs; verify startup failure emits no success message.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T008 [P] [US1] Add unit test for initial ready-event success emission in `tests/unit/tunnel_ready_success_event_test.go`
- [X] T009 [P] [US1] Add integration test proving startup failure suppresses success message in `tests/integration/tunnel_startup_failure_message_test.go`

### Implementation for User Story 1

- [X] T010 [US1] Emit ready-state success confirmation from runtime ready event in `internal/proxy/proxy_runtime.go`
- [X] T011 [US1] Route startup failure path to failure-only guidance without success text in `internal/proxy/proxy_runtime.go`
- [X] T012 [US1] Ensure CLI startup flow surfaces readiness message in `cmd/sql-proxy/main.go`

**Checkpoint**: User Story 1 is fully functional and independently testable.

---

## Phase 4: User Story 2 - Configure PostgreSQL client quickly (Priority: P2)

**Goal**: Include complete PostgreSQL connection parameters and command template in success output.

**Independent Test**: Successful startup output includes local host/port, target instance context, and a copy-ready PostgreSQL command template.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T013 [P] [US2] Add unit test for command template content and argument ordering in `tests/unit/postgres_client_template_test.go`
- [X] T014 [P] [US2] Add integration test for startup output payload completeness in `tests/integration/tunnel_success_instruction_payload_test.go`

### Implementation for User Story 2

- [X] T015 [US2] Populate local host/port and instance context in success payload builder in `internal/proxy/tunnel_success_message.go`
- [X] T016 [US2] Include copy-ready PostgreSQL command template in emitted success block in `internal/proxy/tunnel_success_message.go`
- [X] T017 [US2] Update user documentation with success output field meanings in `README.md`

**Checkpoint**: User Stories 1 and 2 both work independently.

---

## Phase 5: User Story 3 - Receive actionable and safe client guidance (Priority: P3)

**Goal**: Keep success instructions copy-friendly, secret-safe, and reprinted on every successful reconnect.

**Independent Test**: Reconnect emits full instruction block again; output contains placeholders instead of secret values.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T018 [P] [US3] Add unit test enforcing placeholder substitution for credentials in `tests/unit/tunnel_success_message_security_test.go`
- [X] T019 [P] [US3] Add integration test for full guidance reprint on reconnect success in `tests/integration/tunnel_reconnect_success_message_test.go`

### Implementation for User Story 3

- [X] T020 [US3] Enforce secret redaction and placeholder rendering in success message payload in `internal/proxy/tunnel_success_message.go`
- [X] T021 [US3] Emit full success guidance on each successful reconnect event in `internal/proxy/proxy_runtime.go`
- [X] T022 [US3] Add narrow-terminal friendly line wrapping for instruction output in `internal/proxy/tunnel_success_message.go`

**Checkpoint**: All user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final hardening and full-flow verification across stories.

- [X] T023 [P] Add end-to-end integration test for full startup-to-reconnect message lifecycle in `tests/integration/tunnel_success_message_lifecycle_test.go`
- [X] T024 Run and fix full Go test suite for this feature in `tests/unit/` and `tests/integration/`
- [X] T025 Validate quickstart scenarios and align docs wording in `specs/007-postgres-tunnel-success-message/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user stories.
- **User Stories (Phase 3+)**: Depend on Foundational completion.
- **Polish (Phase 6)**: Depends on completion of all targeted user stories.

### User Story Dependencies

- **User Story 1 (P1)**: Starts after Foundational; no dependency on other stories.
- **User Story 2 (P2)**: Starts after Foundational; can run independently of US1 but validates integrated output fields.
- **User Story 3 (P3)**: Starts after Foundational; integrates reconnect behavior with prior payload components.

### Within Each User Story

- Tests first, then implementation.
- Payload/model helpers before runtime wiring changes.
- Runtime changes before documentation updates.

### Parallel Opportunities

- T003 can run in parallel with T001-T002.
- T005 and T006 can run in parallel after T004 starts the shared contract shape.
- T008 and T009 can run in parallel.
- T013 and T014 can run in parallel.
- T018 and T019 can run in parallel.

---

## Parallel Example: User Story 2

```bash
Task: "Add unit test for command template content and argument ordering in tests/unit/postgres_client_template_test.go"
Task: "Add integration test for startup output payload completeness in tests/integration/tunnel_success_instruction_payload_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phases 1 and 2.
2. Deliver Phase 3 (US1) and validate readiness/failure behavior.
3. Stop and validate before expanding instruction payload scope.

### Incremental Delivery

1. Deliver US1 for confidence in readiness timing.
2. Deliver US2 for complete connection instructions.
3. Deliver US3 for security/redaction and reconnect guarantees.
4. Finish with Polish verification and quickstart alignment.

### Parallel Team Strategy

1. One engineer completes foundational formatter/runtime seams (Phase 2).
2. After foundation, engineers split across US1, US2, and US3 tests/implementation.
3. Final merge through Phase 6 integrated lifecycle validation.

---

## Notes

- [P] tasks touch different files and can be executed concurrently.
- [USx] labels map tasks directly to independently testable stories.
- Keep success messaging deterministic: ready/reconnect only, never on failures.
