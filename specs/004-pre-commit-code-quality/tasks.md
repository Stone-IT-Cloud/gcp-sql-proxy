# Tasks: Pre-commit Code Quality Gates

**Input**: Design documents from `/specs/004-pre-commit-code-quality/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are REQUIRED. Include unit tests (table-driven when applicable) and integration coverage for critical flows.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create baseline quality config files and contributor workflow entry points.

- [ ] T001 Create root pre-commit scaffold with `default_stages: [commit]` in `.pre-commit-config.yaml`
- [ ] T002 [P] Create initial Go linter policy file with required linter set in `.golangci.yml`
- [ ] T003 [P] Add pre-commit prerequisite and install instructions section in `README.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Establish shared hook behavior and verification baseline used by all user stories.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [ ] T004 Define common pre-commit repo version pins and baseline config comments in `.pre-commit-config.yaml`
- [ ] T005 Add contributor troubleshooting guidance for missing required tools in `README.md`
- [ ] T006 Create validation notes for local quality-gate verification steps in `specs/004-pre-commit-code-quality/quickstart.md`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - Block non-compliant commits (Priority: P1) 🎯 MVP

**Goal**: Ensure commits are blocked when formatting, linting, testing, or required-tool checks fail.

**Independent Test**: Install hooks, stage intentionally failing changes, run `pre-commit run --all-files`, and verify failures block commit with actionable output.

### Tests for User Story 1 (REQUIRED) ⚠️

- [ ] T007 [P] [US1] Add regression test case covering blocked startup on missing principal tooling guidance in `tests/unit/main_startup_test.go`
- [ ] T008 [P] [US1] Add integration validation flow for commit-blocking failure behavior in `tests/integration/precommit_failure_behavior_test.go`

### Implementation for User Story 1

- [ ] T009 [US1] Add local `go-test` hook (`go test -short ./...`) with commit-blocking behavior in `.pre-commit-config.yaml`
- [ ] T010 [US1] Add local `go-mod-tidy` hook with commit-blocking behavior in `.pre-commit-config.yaml`
- [ ] T011 [US1] Document commit-blocking semantics and remediation guidance in `README.md`

**Checkpoint**: User Story 1 should now block invalid commits with actionable guidance.

---

## Phase 4: User Story 2 - Standardize baseline repository hygiene checks (Priority: P2)

**Goal**: Enforce trailing whitespace, end-of-file newline, and YAML syntax checks in local commits.

**Independent Test**: Stage files with trailing whitespace, missing newline, and invalid YAML; verify hygiene hooks block/fix as configured.

### Tests for User Story 2 (REQUIRED) ⚠️

- [ ] T012 [P] [US2] Add hook configuration regression test fixture for baseline hygiene hooks in `tests/unit/testdata/precommit_baseline.yaml`
- [ ] T013 [P] [US2] Add integration check for whitespace/newline/YAML hook execution in `tests/integration/precommit_baseline_hooks_test.go`

### Implementation for User Story 2

- [ ] T014 [US2] Add `pre-commit-hooks` repository block with `trailing-whitespace`, `end-of-file-fixer`, and `check-yaml` hooks in `.pre-commit-config.yaml`
- [ ] T015 [US2] Document baseline hook behavior and expected auto-fix flow in `README.md`

**Checkpoint**: User Stories 1 and 2 both operate independently for commit gating and baseline hygiene.

---

## Phase 5: User Story 3 - Enforce Go-specific quality checks (Priority: P3)

**Goal**: Run Go formatting/import, linting, and short tests for staged Go-related changes while skipping Go hooks for non-Go-only commits.

**Independent Test**: Stage Go changes with lint/test issues and verify blocking behavior; stage non-Go-only changes and verify Go hooks are skipped.

### Tests for User Story 3 (REQUIRED) ⚠️

- [ ] T016 [P] [US3] Add configuration-focused unit test fixture for Go-scoped hook file filters in `tests/unit/testdata/precommit_go_scoped.yaml`
- [ ] T017 [P] [US3] Add integration validation for Go-file scoping and Go hook execution in `tests/integration/precommit_go_hooks_scope_test.go`

### Implementation for User Story 3

- [ ] T018 [US3] Add `golangci-lint` pre-commit repository block and hook configuration in `.pre-commit-config.yaml`
- [ ] T019 [US3] Configure required linters (`errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `gofmt`) in `.golangci.yml`
- [ ] T020 [US3] Scope Go-specific local hooks to staged Go module/source/test files in `.pre-commit-config.yaml`
- [ ] T021 [US3] Update quick verification commands for Go and non-Go commit paths in `specs/004-pre-commit-code-quality/quickstart.md`
- [ ] T022 [US3] Add explicit Go import-normalization hook configuration and verification notes in `.pre-commit-config.yaml`

**Checkpoint**: All user stories are independently functional and verified.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final verification, consistency checks, and readiness for implementation handoff.

- [ ] T023 [P] Run full local validation (`pre-commit run --all-files` and `go test ./...`) and record outcomes in `specs/004-pre-commit-code-quality/quickstart.md`
- [ ] T024 Verify configuration aligns with contract requirements in `specs/004-pre-commit-code-quality/contracts/pre-commit-contract.md`
- [ ] T025 [P] Confirm documentation setup path satisfies <10 minute onboarding target in `README.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - starts immediately.
- **Foundational (Phase 2)**: Depends on Setup completion - blocks all user story work.
- **User Stories (Phase 3+)**: Depend on Foundational completion.
  - Preferred sequence: US1 (MVP) → US2 → US3.
  - US2 and US3 can proceed in parallel after US1 if team capacity allows.
- **Polish (Phase 6)**: Depends on completion of all selected user stories.

### User Story Dependencies

- **US1 (P1)**: Starts after Foundational phase; no dependency on other stories.
- **US2 (P2)**: Starts after Foundational phase; independent but complements US1 commit gating.
- **US3 (P3)**: Starts after Foundational phase; depends on shared config structure and can be implemented independently of US2 content.

### Within Each User Story

- Tests are authored before implementation tasks and validated against expected failing/passing behavior.
- Configuration tasks precede documentation updates for the same story.
- Story-level independent verification must pass before moving to next priority.

### Parallel Opportunities

- Phase 1 tasks marked `[P]` can run concurrently.
- Foundational doc/check tasks can run in parallel once base config exists.
- In each user story, fixture-based test tasks marked `[P]` can run in parallel with integration test scaffolding.
- After Foundational completion, US2 and US3 can be staffed concurrently.

---

## Parallel Example: User Story 3

```bash
# Parallel test scaffolding:
Task: "Add configuration-focused unit test fixture in tests/unit/testdata/precommit_go_scoped.yaml"
Task: "Add integration validation in tests/integration/precommit_go_hooks_scope_test.go"

# Parallel implementation/documentation after hook config:
Task: "Configure required linters in .golangci.yml"
Task: "Update quick verification commands in specs/004-pre-commit-code-quality/quickstart.md"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 and Phase 2.
2. Deliver Phase 3 (US1) to enforce commit-blocking quality checks.
3. Validate independent US1 behavior before expanding scope.

### Incremental Delivery

1. Add US1 for core commit gate behavior.
2. Add US2 for baseline hygiene enforcement.
3. Add US3 for Go-specific quality controls and scoping.
4. Finish with Phase 6 cross-cutting verification.

### Parallel Team Strategy

1. One developer finalizes shared setup/foundation.
2. Then split:
   - Developer A: US2 baseline hygiene.
   - Developer B: US3 Go quality/scoping.
3. Rejoin for final polish and end-to-end validation.

---

## Notes

- All tasks use the required checklist format with task ID, optional `[P]`, optional `[US#]`, and explicit file path.
- User story tasks are independently testable and mapped to corresponding specification priorities.
- Keep commit messages scoped to logical task groups for easier review.
