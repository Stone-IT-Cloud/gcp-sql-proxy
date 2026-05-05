# Phase 0 Research: Pre-commit Code Quality Gates

## Decision 1: Use official pre-commit ecosystem hooks for baseline checks

- **Decision**: Use `pre-commit/pre-commit-hooks` for `trailing-whitespace`, `end-of-file-fixer`, and `check-yaml`.
- **Rationale**: These hooks are widely adopted, low maintenance, and directly match required baseline hygiene checks.
- **Alternatives considered**:
  - Local custom scripts for whitespace/YAML checks (rejected: unnecessary maintenance burden).
  - CI-only validation (rejected: does not enforce quality before commit).

## Decision 2: Use official golangci-lint pre-commit hook

- **Decision**: Integrate `golangci/golangci-lint` through its pre-commit hook and define active linters in `.golangci.yml`.
- **Rationale**: This provides stable local lint enforcement aligned with repository policy and avoids ad hoc local lint wrappers.
- **Alternatives considered**:
  - Running only `go vet` as lint gate (rejected: insufficient static analysis coverage).
  - Local custom linter script (rejected: duplicates official integration and increases maintenance cost).

## Decision 3: Implement local hooks for `go mod tidy` and `go test -short ./...`

- **Decision**: Add a `repo: local` hook block for module hygiene and short unit tests, with Go-specific file targeting where applicable.
- **Rationale**: Local hooks capture project-specific commands and keep commit-time validation directly aligned with current repository workflow.
- **Alternatives considered**:
  - Only run tests in CI (rejected: allows avoidable breakage into shared branches).
  - Run full test suites on every commit (rejected: excessive local feedback time).

## Decision 4: Make missing required tools commit-blocking with remediation guidance

- **Decision**: Treat missing required tooling as a failed commit validation state and provide setup guidance in project documentation.
- **Rationale**: Blocking behavior guarantees consistent enforcement; guidance minimizes friction for new contributors.
- **Alternatives considered**:
  - Warn and continue commit (rejected: inconsistent quality gate behavior).
  - Interactive prompt to bypass (rejected: weakens deterministic enforcement).

## Decision 5: Gate Go-specific checks on staged Go-related changes

- **Decision**: Execute Go-specific hooks only when staged files include Go module/source/test files.
- **Rationale**: Preserves enforcement for relevant changes while reducing unnecessary runtime for non-Go commits.
- **Alternatives considered**:
  - Always run all Go checks (rejected: increases commit latency for unrelated changes).
  - Run formatting always, lint/tests conditionally (rejected: added complexity with limited incremental value).
