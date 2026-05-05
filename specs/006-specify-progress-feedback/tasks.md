# Tasks: Specify Progress Feedback

**Input**: Design documents from `/specs/006-specify-progress-feedback/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED. Include unit tests (table-driven when applicable) and integration
coverage for critical flows.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish fixtures and baseline test harness for progress-feedback behavior.

- [X] T001 Create baseline fixture for phase progress output in `tests/unit/testdata/specify_progress_phase_output.txt`
- [X] T002 Create baseline fixture for failure progress output in `tests/unit/testdata/specify_progress_failure_output.txt`
- [X] T003 [P] Create baseline fixture for liveness and skipped events in `tests/unit/testdata/specify_progress_liveness_output.txt`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Add shared progress-event abstractions and renderer fallback mechanism.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Add shared progress-event model and helpers for `/speckit-specify` phases in `internal/specify/progress_events.go`
- [X] T005 [P] Add progress renderer interface supporting rich and plain-text modes in `internal/specify/progress_renderer.go`
- [X] T006 [P] Add fallback transition logic from rich to plain-text rendering in `internal/specify/progress_fallback.go`
- [X] T007 Add foundational unit tests for event model and renderer fallback in `tests/unit/specify_progress_foundation_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - See live execution status for `/speckit-specify` (Priority: P1) 🎯 MVP

**Goal**: Emit clear start/finish phase updates and final completion summary during normal runs.

**Independent Test**: Execute `/speckit-specify` happy path and verify phase-based progress updates plus final summary output.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T008 [P] [US1] Add unit tests for per-phase started/completed events in `tests/unit/specify_progress_phase_events_test.go`
- [X] T009 [P] [US1] Add integration test for normal-run phase visibility and completion summary in `tests/integration/specify_progress_happy_path_test.go`

### Implementation for User Story 1

- [X] T010 [US1] Emit progress events for major `/speckit-specify` phases in `internal/specify/specify_command.go`
- [X] T011 [US1] Emit pre-hook start/finish progress transitions in `internal/specify/specify_hooks.go`
- [X] T012 [US1] Emit completion message summarizing generated artifacts in `internal/specify/specify_command.go`

**Checkpoint**: User Story 1 is fully functional and independently testable.

---

## Phase 4: User Story 2 - Understand failures with context (Priority: P2)

**Goal**: Provide failed-phase context and actionable remediation guidance when execution stops.

**Independent Test**: Force failure during a known phase and verify output contains failed phase and next-step guidance.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T013 [P] [US2] Add unit test for failure-event payload requirements in `tests/unit/specify_progress_failure_event_test.go`
- [X] T014 [P] [US2] Add integration test for phase-specific failure guidance output in `tests/integration/specify_progress_failure_guidance_test.go`

### Implementation for User Story 2

- [X] T015 [US2] Add failure-event emission with failed-phase identification in `internal/specify/specify_command.go`
- [X] T016 [US2] Add actionable next-step guidance mapping by failure phase in `internal/specify/specify_errors.go`
- [X] T017 [US2] Ensure failure output excludes stack traces and sensitive values in `internal/specify/specify_errors.go`

**Checkpoint**: User Stories 1 and 2 both work independently.

---

## Phase 5: User Story 3 - Receive consistent feedback cadence across long-running steps (Priority: P3)

**Goal**: Emit periodic liveness updates, explicit skipped-step events, and plain-text fallback continuity.

**Independent Test**: Run long-running and skipped-step scenarios to verify 10-second heartbeat cadence, skipped event reason, and resilient fallback behavior.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T018 [P] [US3] Add unit test for 10-second heartbeat scheduler behavior in `tests/unit/specify_progress_heartbeat_test.go`
- [X] T019 [P] [US3] Add unit test for explicit skipped-event structure in `tests/unit/specify_progress_skipped_event_test.go`
- [X] T020 [P] [US3] Add integration test for rich-render failure fallback to plain text in `tests/integration/specify_progress_fallback_test.go`

### Implementation for User Story 3

- [X] T021 [US3] Implement long-running phase heartbeat emission at <=10-second intervals in `internal/specify/specify_progress_loop.go`
- [X] T022 [US3] Emit explicit skipped events with step identity and reason in `internal/specify/specify_hooks.go`
- [X] T023 [US3] Wire automatic fallback to plain-text progress output without aborting execution in `internal/specify/progress_fallback.go`

**Checkpoint**: All user stories are independently functional.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final consistency checks and documentation alignment across all stories.

- [X] T024 [P] Add integration test covering full lifecycle (phase, skip, heartbeat, completion/failure) in `tests/integration/specify_progress_lifecycle_test.go`
- [X] T025 Run and fix full Go test suite for this feature in `tests/unit/` and `tests/integration/`
- [X] T026 Validate quickstart scenarios and align terminology in `specs/006-specify-progress-feedback/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user stories.
- **User Stories (Phase 3+)**: Depend on Foundational completion.
- **Polish (Phase 6)**: Depends on completion of all targeted user stories.

### User Story Dependencies

- **User Story 1 (P1)**: Starts after Foundational; no dependency on other stories.
- **User Story 2 (P2)**: Starts after Foundational; validates failure path over US1 phase model.
- **User Story 3 (P3)**: Starts after Foundational; extends runtime loop/fallback behavior and remains independently testable.

### Within Each User Story

- Write tests first, then implement.
- Phase/event model changes before UI/output rendering changes.
- Rendering fallback logic before lifecycle integration verification.

### Parallel Opportunities

- T003 can run in parallel with T001-T002.
- T005 and T006 can run in parallel after T004 defines shared event structures.
- T008 and T009 can run in parallel.
- T013 and T014 can run in parallel.
- T018, T019, and T020 can run in parallel.

---

## Parallel Example: User Story 3

```bash
Task: "Add unit test for 10-second heartbeat scheduler behavior in tests/unit/specify_progress_heartbeat_test.go"
Task: "Add unit test for explicit skipped-event structure in tests/unit/specify_progress_skipped_event_test.go"
Task: "Add integration test for rich-render failure fallback to plain text in tests/integration/specify_progress_fallback_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phases 1 and 2.
2. Deliver Phase 3 (US1) and validate normal progress visibility.
3. Stop and validate before failure/cadence expansions.

### Incremental Delivery

1. Deliver US1 for baseline progress visibility.
2. Deliver US2 for actionable failure context.
3. Deliver US3 for liveness cadence, skip signaling, and fallback resilience.
4. Finish with lifecycle polish and quickstart validation.

### Parallel Team Strategy

1. One engineer delivers foundational event/rendering layer (Phase 2).
2. Then split story work across US1, US2, and US3 with test-first sequencing.
3. Merge and stabilize with lifecycle integration test in Phase 6.

---

## Notes

- [P] tasks touch different files and can be executed concurrently.
- [USx] labels map tasks directly to independently testable stories.
- Keep progress output concise, explicit, and resilient under failure/cancellation.
