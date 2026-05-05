# Tasks: GitHub Actions Automated Build Pipeline

**Input**: Design documents from `/specs/005-github-actions-build-pipeline/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED. Include unit tests (table-driven when applicable) and integration coverage for critical flows.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Prepare workflow location, baseline job skeleton, and shared naming variables.

- [X] T001 Create workflow directory and base workflow file at `.github/workflows/build.yml`
- [X] T002 [P] Add workflow-level naming and metadata comments in `.github/workflows/build.yml`
- [X] T003 [P] Add release/build workflow documentation section in `README.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Define trigger and matrix foundation required by all user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Add trigger configuration for push main, pull_request main, tag `v*`, and `workflow_dispatch` in `.github/workflows/build.yml`
- [X] T005 Add base `build` job with matrix target definitions and dynamic runner selection in `.github/workflows/build.yml`
- [X] T006 Add matrix failure behavior configuration (continue targets, fail overall) in `.github/workflows/build.yml`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - Build binaries automatically on repository events (Priority: P1) 🎯 MVP

**Goal**: Ensure all required repository events execute the build workflow and compile binaries across matrix targets.

**Independent Test**: Trigger workflow on pull request and manual dispatch, verify all matrix jobs start and complete with per-target status.

### Tests for User Story 1 (REQUIRED) ⚠️

- [X] T007 [P] [US1] Add workflow trigger coverage fixture in `tests/unit/testdata/build_workflow_triggers.yml`
- [X] T008 [P] [US1] Add integration workflow trigger/matrix validation test in `tests/integration/build_workflow_trigger_matrix_test.go`

### Implementation for User Story 1

- [X] T009 [US1] Add checkout and Go setup steps (`actions/checkout@v4`, `actions/setup-go@v5` with `go-version-file`) in `.github/workflows/build.yml`
- [X] T010 [US1] Add `go build` step with matrix-derived `GOOS`/`GOARCH` and version ldflags support in `.github/workflows/build.yml`
- [X] T011 [US1] Document supported workflow trigger paths and expected matrix behavior in `README.md`

**Checkpoint**: User Story 1 workflow automation is functional and independently testable.

---

## Phase 4: User Story 2 - Produce uniquely named cross-platform binaries (Priority: P2)

**Goal**: Guarantee deterministic, non-overlapping binary names per target, including Windows extension handling.

**Independent Test**: Execute matrix run and confirm output names include target OS/architecture and `.exe` for Windows.

### Tests for User Story 2 (REQUIRED) ⚠️

- [X] T012 [P] [US2] Add binary naming fixture expectations in `tests/unit/testdata/build_workflow_binary_names.yml`
- [X] T013 [P] [US2] Add integration check for target-specific binary naming in `tests/integration/build_workflow_binary_naming_test.go`

### Implementation for User Story 2

- [X] T014 [US2] Add output filename generation logic for `gcp-db-proxy-<target_os>-<target_arch>` in `.github/workflows/build.yml`
- [X] T015 [US2] Add Windows `.exe` suffix branching logic in `.github/workflows/build.yml`
- [X] T016 [US2] Update quick validation commands for output naming verification in `specs/005-github-actions-build-pipeline/quickstart.md`

**Checkpoint**: User Stories 1 and 2 now produce valid platform-specific build outputs.

---

## Phase 5: User Story 3 - Store build outputs as downloadable artifacts with retention control (Priority: P3)

**Goal**: Upload one artifact per matrix target with clear naming and explicit retention policy.

**Independent Test**: Run workflow and verify each target uploads a distinct artifact with 14-day retention.

### Tests for User Story 3 (REQUIRED) ⚠️

- [X] T017 [P] [US3] Add artifact upload contract fixture in `tests/unit/testdata/build_workflow_artifacts.yml`
- [X] T018 [P] [US3] Add integration check for per-target artifact upload and retention in `tests/integration/build_workflow_artifacts_test.go`

### Implementation for User Story 3

- [X] T019 [US3] Add `actions/upload-artifact@v4` step with per-target artifact naming in `.github/workflows/build.yml`
- [X] T020 [US3] Set `retention-days: 14` and target-specific artifact path wiring in `.github/workflows/build.yml`
- [X] T021 [US3] Add workflow artifact retrieval guidance in `README.md`
- [X] T022 [US3] Align contract verification examples with per-target artifact behavior in `specs/005-github-actions-build-pipeline/contracts/workflow-build-contract.md`

**Checkpoint**: All user stories are independently functional and verified.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, documentation consistency, and implementation readiness checks.

- [X] T023 [P] Run repository validation (`go test ./...`) and record outcomes in `specs/005-github-actions-build-pipeline/quickstart.md`
- [X] T024 Validate workflow YAML formatting and trigger syntax with local checks in `.github/workflows/build.yml`
- [X] T025 [P] Verify quickstart end-to-end flow against final workflow behavior in `specs/005-github-actions-build-pipeline/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user stories.
- **User Stories (Phase 3+)**: Depend on Foundational completion.
  - Preferred order: US1 → US2 → US3.
  - US2 and US3 can run in parallel after US1 if capacity allows.
- **Polish (Phase 6)**: Depends on all implemented user stories.

### User Story Dependencies

- **US1 (P1)**: Starts after Foundational completion; provides baseline build execution.
- **US2 (P2)**: Depends on US1 build pipeline scaffolding but remains independently testable.
- **US3 (P3)**: Depends on US1/US2 output flow and artifact naming contracts, while maintaining independent validation.

### Within Each User Story

- Tests are authored before implementation and validated against expected behavior.
- Workflow structure tasks precede docs updates within each story.
- Story-level checkpoints must pass before moving to next priority.

### Parallel Opportunities

- Setup tasks marked `[P]` can run in parallel.
- Test fixture and integration validation tasks marked `[P]` per story can run in parallel.
- After Foundational phase, US2 and US3 can proceed in parallel if teams are split.
- Final documentation and validation tasks marked `[P]` can run concurrently.

---

## Parallel Example: User Story 3

```bash
# Parallel test preparation:
Task: "Add artifact upload contract fixture in tests/unit/testdata/build_workflow_artifacts.yml"
Task: "Add integration check in tests/integration/build_workflow_artifacts_test.go"

# Parallel post-implementation updates:
Task: "Add workflow artifact retrieval guidance in README.md"
Task: "Align contract verification examples in specs/005-github-actions-build-pipeline/contracts/workflow-build-contract.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup and Foundational phases.
2. Deliver US1 trigger + build matrix execution path.
3. Validate US1 independently before extending scope.

### Incremental Delivery

1. Complete shared foundation and event trigger setup.
2. Add US1 for automated build execution.
3. Add US2 for deterministic binary naming.
4. Add US3 for per-target artifact retention-governed uploads.
5. Finish with cross-cutting validation and documentation polish.

### Parallel Team Strategy

1. One engineer finalizes Setup + Foundational workflow scaffolding.
2. Then split:
   - Engineer A: US2 naming and output verification.
   - Engineer B: US3 artifact upload and retention controls.
3. Rejoin for final validation and documentation pass.

---

## Notes

- All tasks comply with checklist format requirements and include explicit file paths.
- Story labels are present for all user-story implementation and testing tasks.
- Keep commits scoped by phase or logical task clusters for easier review.
