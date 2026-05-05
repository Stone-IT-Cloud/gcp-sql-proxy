# Implementation Plan: Pre-commit Code Quality Gates

**Branch**: `004-pre-commit-code-quality` | **Date**: 2026-05-05 | **Spec**: [`specs/004-pre-commit-code-quality/spec.md`](./spec.md)
**Input**: Feature specification from `/specs/004-pre-commit-code-quality/spec.md`

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Introduce repository-level local quality gates using pre-commit so every commit is validated for baseline file hygiene, Go formatting/import consistency, linting, and short unit tests. The implementation adds a root pre-commit configuration, a golangci-lint ruleset, and contributor setup guidance while keeping Go-specific checks conditional on staged Go-related files to balance quality and runtime cost.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.25.7, YAML-based tool configuration  
**Primary Dependencies**: pre-commit, pre-commit-hooks, golangci-lint pre-commit hook  
**Storage**: N/A (configuration-only feature; files in repository)  
**Testing**: `go test -short ./...` via pre-commit local hook and existing Go test suites  
**Target Platform**: Contributor development environments on Windows, macOS, and Linux  
**Project Type**: Go CLI application  
**Performance Goals**: At least 95% of commit attempts complete pre-commit validation within 120 seconds  
**Constraints**: Missing required tooling must block commit with install guidance; Go-specific checks run only for staged Go-related files  
**Scale/Scope**: Repository-wide pre-commit enforcement for all contributors and all local commits

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] Code structure uses `cmd/` and `internal/` with idiomatic Go conventions.
- [x] Error handling strategy uses wrapped errors and user-friendly CLI messages; no panic-based flow.
- [x] Concurrency design documents goroutine lifecycle, cancellation, and disconnection cleanup.
- [x] Testing plan includes table-driven unit tests and required integration coverage.
- [x] Security controls include IAM/private tunnel defaults and secure local token handling.
- [x] Dependency plan is minimal and justified (standard library first, then approved packages).

### Post-Design Re-Check

- [x] Design keeps existing `cmd/` and `internal/` boundaries intact (config/docs-only additions).
- [x] Error experience is documented via actionable hook failure and setup guidance requirements.
- [x] Concurrency lifecycle remains unaffected (no runtime goroutine changes introduced).
- [x] Validation strategy includes local short-test gate and repository verification steps.
- [x] Security/dependency posture is preserved through minimal, explicit tooling additions.

## Project Structure

### Documentation (this feature)

```text
specs/004-pre-commit-code-quality/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
cmd/
└── sql-proxy/

internal/
├── auth/
└── proxy/

tests/
├── integration/
└── unit/

.pre-commit-config.yaml
.golangci.yml
README.md
```

**Structure Decision**: Keep the existing single Go CLI project layout and add repository-root quality configuration files plus README updates; no new runtime packages are required.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
