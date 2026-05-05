<!--
Sync Impact Report
- Version change: 0.0.0 -> 1.0.0
- Modified principles:
  - Placeholder Principle 1 -> I. Idiomatic Go and Project Structure
  - Placeholder Principle 2 -> II. Reliable Error Handling and CLI UX
  - Placeholder Principle 3 -> III. Concurrency Safety and Lifecycle Management (NON-NEGOTIABLE)
  - Placeholder Principle 4 -> IV. Testability and Context-Driven Design
  - Placeholder Principle 5 -> V. Security and Dependency Discipline
- Added sections:
  - Operational Constraints
  - Development Workflow and Quality Gates
- Removed sections:
  - None
- Templates requiring updates:
  - ✅ updated: .specify/templates/plan-template.md
  - ✅ updated: .specify/templates/spec-template.md
  - ✅ updated: .specify/templates/tasks-template.md
  - ⚠ pending (not present): .specify/templates/commands/*.md
  - ✅ reviewed (no change required): README.md
- Follow-up TODOs:
  - TODO(RATIFICATION_DATE): original ratification date is unknown; replace placeholder when known.
-->
# Project Constitution: GCP Cloud SQL CLI Proxy

## Core Principles

### I. Idiomatic Go and Project Structure
Code MUST follow idiomatic Go conventions and pass `gofmt`, `go vet`, and linter checks.
Source layout MUST use `cmd/` for entrypoints and `internal/` for non-public implementation.
Exported identifiers MUST be descriptive, local variables SHOULD remain concise, and unnecessary
getters/setters MUST be avoided unless required for abstraction boundaries.
Rationale: Consistent style and structure reduce review friction and long-term maintenance cost.

### II. Reliable Error Handling and CLI UX
All errors MUST be handled explicitly with wrapping (`fmt.Errorf("%w", err)`) and inspected using
`errors.Is`/`errors.As` where branching is needed. Panics MUST NOT be used for expected runtime
failures. CLI errors MUST be user-friendly, concise, and actionable, without exposing stack traces
or internal details. Control flow SHOULD prefer early returns to keep logic simple.
Rationale: Predictable error behavior and clear output are essential for operator trust and support.

### III. Concurrency Safety and Lifecycle Management (NON-NEGOTIABLE)
Proxy runtime components MUST manage goroutine lifecycles explicitly and terminate cleanly on
disconnect, cancellation, or shutdown. Bidirectional stream copying and local TCP serving MUST use
safe synchronization patterns and MUST avoid leaks. Long-running operations MUST accept
`context.Context` and honor cancellation.
Rationale: The proxy is connection-heavy; lifecycle bugs degrade stability and can leak resources.

### IV. Testability and Context-Driven Design
Features MUST be modular and follow single-responsibility boundaries. Tests MUST be table-driven
where applicable and use Go's `testing` package as the default. Unit and integration coverage MUST
verify core proxy flows, auth behavior, and graceful shutdown semantics. New long-running logic MUST
be context-aware and testable without live cloud dependencies when feasible.
Rationale: Testable design allows safe iteration across auth, networking, and platform differences.

### V. Security and Dependency Discipline
Sensitive data MUST NEVER be hardcoded. Local credential and token artifacts MUST enforce restricted
permissions (for example `0600` on Unix-like systems) or platform-equivalent protection. The proxy
MUST prefer Zero Trust posture: IAM authentication and Private IP tunneling by default. Dependencies
MUST remain minimal, use Go modules, avoid vendoring, and prioritize standard library plus officially
supported Google packages.
Rationale: This tool abstracts complex cloud networking and must preserve strong security guarantees.

## Operational Constraints

- The CLI MUST remain cross-platform for Windows, macOS, and Linux.
- User-facing output MUST be clear and readable; color can be used when terminal support exists, but
  functionality MUST NOT depend on color.
- The implementation MUST hide OAuth/IAP/IAM complexity from end-users while preserving transparent,
  auditable behavior for maintainers.
- Performance-sensitive paths (listener setup, proxy loop, shutdown) MUST avoid unnecessary abstractions.

## Development Workflow and Quality Gates

- Every change MUST include corresponding tests or a documented rationale when tests are not applicable.
- Pull requests MUST confirm: formatting/lint pass, test pass, context cancellation behavior, and no
  secret leakage through code or logs.
- Reviewers MUST reject changes that violate core principles unless a documented constitutional exception
  is approved under Governance.
- Complexity additions MUST include explicit justification and simpler rejected alternatives.

## Governance

This constitution supersedes local conventions for this repository.

Amendment process:
- Propose changes in a pull request that includes rationale, impacted principles, and migration steps.
- Obtain maintainer approval before merge.
- Update dependent templates and guidance docs in the same change when principles are affected.

Versioning policy:
- MAJOR: backward-incompatible governance changes or principle removals/redefinitions.
- MINOR: new principle/section or materially expanded guidance.
- PATCH: wording clarifications without semantic governance change.

Compliance review expectations:
- Every implementation plan MUST pass a Constitution Check before design and before implementation.
- Every tasks artifact MUST map required quality, testing, and security work to executable tasks.
- Periodic audits SHOULD validate that templates remain aligned with this constitution.

**Version**: 1.0.0 | **Ratified**: TODO(RATIFICATION_DATE): original adoption date is unknown | **Last Amended**: 2026-05-05
