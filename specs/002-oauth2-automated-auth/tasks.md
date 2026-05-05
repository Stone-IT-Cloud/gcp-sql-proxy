# Tasks: OAuth2 Automated Authentication

**Input**: Design documents from `/specs/002-oauth2-automated-auth/`
**Prerequisites**: plan.md (required), spec.md (required), research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are REQUIRED. Include table-driven unit tests and integration coverage for auth startup flows.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: OAuth authentication project scaffolding and dependencies

- [x] T001 Create auth package directory and base file at `internal/auth/auth.go`
- [x] T002 [P] Add OAuth2 dependencies in `go.mod` (`golang.org/x/oauth2`, `golang.org/x/oauth2/google`)
- [x] T003 [P] Add test file skeletons `tests/unit/auth_token_test.go` and `tests/integration/oauth_flow_test.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core auth primitives required by all user stories

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Define OAuth config builder and injectable client credential fields in `internal/auth/auth.go`
- [x] T005 Implement token path resolver (`~/.sql-proxy/token.json`) in `internal/auth/auth.go`
- [x] T006 Implement token load/decode helper from token file in `internal/auth/auth.go`
- [x] T007 Implement secure token save helper with owner-only permissions in `internal/auth/auth.go`
- [x] T008 Implement browser launcher helper for Linux, Windows, and macOS in `internal/auth/auth.go`
- [x] T009 Implement auth state generation and validation helper in `internal/auth/auth.go`
- [x] T010 Implement invalid-token rename/remove helper in `internal/auth/auth.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Authenticate Without External CLI (Priority: P1) 🎯 MVP

**Goal**: Complete first-run browser OAuth flow and return authenticated HTTP client without external CLI.

**Independent Test**: Run first-time auth without token and confirm callback code exchange, success response, token persistence, and authenticated client creation.

### Tests for User Story 1 (REQUIRED) ⚠️

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [x] T011 [P] [US1] Add table-driven unit tests for OAuth config construction and scope/redirect settings in `tests/unit/auth_token_test.go`
- [x] T012 [P] [US1] Add integration test for first-run browser auth flow and token persistence in `tests/integration/oauth_flow_test.go`
- [x] T013 [US1] Add integration test for callback success response and code capture lifecycle in `tests/integration/oauth_flow_test.go`

### Implementation for User Story 1

- [x] T014 [US1] Implement temporary callback server startup and callback handler in `internal/auth/auth.go`
- [x] T015 [US1] Implement auth URL generation with offline access and per-session state in `internal/auth/auth.go`
- [x] T016 [US1] Implement code exchange and token persistence path in `internal/auth/auth.go`
- [x] T017 [US1] Implement exported `GetClient(ctx)` first-run path returning authenticated HTTP client in `internal/auth/auth.go`
- [x] T018 [US1] Wire auth client creation into startup integration point in `cmd/sql-proxy/main.go`

**Checkpoint**: User Story 1 is fully functional and testable independently

---

## Phase 4: User Story 2 - Reuse Persisted Token Securely (Priority: P2)

**Goal**: Reuse valid token on startup while enforcing secure storage behavior.

**Independent Test**: Provide valid token file and verify startup uses token directly without browser flow and permissions remain restricted.

### Tests for User Story 2 (REQUIRED) ⚠️

- [x] T019 [P] [US2] Add unit tests for token file read/write and permission handling in `tests/unit/auth_token_test.go`
- [x] T020 [US2] Add integration test verifying valid-token startup skips browser flow in `tests/integration/oauth_flow_test.go`

### Implementation for User Story 2

- [x] T021 [US2] Implement persisted-token validation and reuse branch in `internal/auth/auth.go`
- [x] T022 [US2] Enforce owner-only token file permission checks and correction in `internal/auth/auth.go`
- [x] T023 [US2] Update startup auth call sequence to prefer valid persisted token path in `cmd/sql-proxy/main.go`

**Checkpoint**: User Stories 1 and 2 work independently and together

---

## Phase 5: User Story 3 - Recover From Invalid or Expired Token (Priority: P3)

**Goal**: Recover automatically from invalid credentials and secure callback edge cases.

**Independent Test**: Supply invalid token and callback errors/state mismatch; verify invalid token cleanup, fresh re-auth, and graceful callback server shutdown.

### Tests for User Story 3 (REQUIRED) ⚠️

- [x] T024 [P] [US3] Add unit tests for invalid-token rename/remove behavior in `tests/unit/auth_token_test.go`
- [x] T025 [P] [US3] Add integration test for callback state mismatch failure path in `tests/integration/oauth_flow_test.go`
- [x] T026 [US3] Add integration test for callback port `8080` fallback to available localhost port in `tests/integration/oauth_flow_test.go`

### Implementation for User Story 3

- [x] T027 [US3] Implement invalid/expired token detection and cleanup before re-auth in `internal/auth/auth.go`
- [x] T028 [US3] Implement strict callback state mismatch failure behavior in `internal/auth/auth.go`
- [x] T029 [US3] Implement callback server port fallback when `8080` is occupied in `internal/auth/auth.go`
- [x] T030 [US3] Implement callback server graceful shutdown on success/failure paths in `internal/auth/auth.go`

**Checkpoint**: All user stories are independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final hardening, documentation, and regression coverage

- [x] T031 [P] Add/refresh exported GoDoc comments for auth package APIs in `internal/auth/auth.go`
- [x] T032 [P] Update OAuth quickstart verification steps from implementation details in `specs/002-oauth2-automated-auth/quickstart.md`
- [x] T033 Validate that auth feature does not degrade IAM/private tunnel defaults in startup path via integration assertions in `tests/integration/oauth_flow_test.go`
- [x] T034 Run full test suite and capture sign-off evidence with `go test ./...`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies; starts immediately
- **Foundational (Phase 2)**: Depends on Setup completion; blocks all user stories
- **User Stories (Phase 3-5)**: Depend on Foundational completion
  - US1 delivers MVP authentication capability
  - US2 depends on token persistence behavior introduced in US1
  - US3 depends on US1/US2 auth flows for recovery and fallback handling
- **Polish (Phase 6)**: Depends on completion of desired user stories

### User Story Dependencies

- **US1 (P1)**: Can start after Phase 2; no dependency on other stories
- **US2 (P2)**: Starts after US1 token lifecycle foundation
- **US3 (P3)**: Starts after US1/US2 establish callback and token reuse flows

### Within Each User Story

- Tests MUST be written and fail before implementation
- Config/state helpers before callback server runtime wiring
- Token persistence before token reuse and recovery logic
- Callback handling before integration into startup path

### Parallel Opportunities

- Setup tasks marked `[P]` can run in parallel
- Foundational helper tasks on distinct concerns can be parallelized after base auth file exists
- US1 unit and integration tests can run in parallel
- US3 unit invalid-token tests and state mismatch tests can run in parallel

---

## Parallel Example: User Story 1

```bash
# Run US1 test tasks together:
Task: "T011 Add table-driven unit tests for OAuth config construction in tests/unit/auth_token_test.go"
Task: "T012 Add integration test for first-run browser auth flow in tests/integration/oauth_flow_test.go"

# Implement core auth flow wiring:
Task: "T014 Implement temporary callback server startup and callback handler in internal/auth/auth.go"
Task: "T015 Implement auth URL generation with offline access and per-session state in internal/auth/auth.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Setup + Foundational phases
2. Complete US1 flow end-to-end
3. Validate first-run auth without external CLI
4. Demo authenticated client integration path

### Incremental Delivery

1. Deliver US1: first-run browser authentication
2. Deliver US2: persisted-token secure reuse
3. Deliver US3: invalid-token and callback recovery hardening
4. Complete polish and full regression pass

### Parallel Team Strategy

With multiple developers:

1. Developer A: Foundation + callback runtime flow
2. Developer B: Token persistence + reuse path
3. Developer C: Recovery/fallback edge cases and integration tests
4. Shared final polish for docs and security baseline validation

---

## Notes

- All tasks follow required checklist format with IDs, optional `[P]`, `[USx]` labels, and exact file paths.
- Keep task-to-commit traceability tight for easier PR review.
- Validate each story independently at checkpoints before proceeding.
