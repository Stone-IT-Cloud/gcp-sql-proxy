# Phase 0 Research: GitHub Actions Automated Build Pipeline

## Decision 1: Use one matrix job with dynamic runner selection

- **Decision**: Use a single `build` job with `strategy.matrix` including runner label (`ubuntu-latest`, `windows-latest`, `macos-latest`) and target OS/arch fields.
- **Rationale**: A unified matrix keeps workflow maintenance low while ensuring platform-specific execution and clear per-target status.
- **Alternatives considered**:
  - Separate jobs per platform (rejected: duplicated steps and more maintenance overhead).
  - Static single-runner cross-compilation only on Linux (rejected: reduced confidence for platform-specific shell/file handling).

## Decision 2: Continue matrix execution on failures and fail overall workflow

- **Decision**: Configure matrix behavior to keep running all targets even if one fails, while preserving overall failed workflow status.
- **Rationale**: Maintainers get full target visibility in a single run, speeding triage and reducing reruns.
- **Alternatives considered**:
  - Fail-fast behavior (rejected: hides potential additional failures).
  - Trigger-dependent behavior (rejected: inconsistent operational expectations).

## Decision 3: Use per-target artifact packaging

- **Decision**: Upload one artifact per matrix target using a deterministic artifact name aligned with binary filename.
- **Rationale**: Improves retrieval clarity and troubleshooting by mapping artifacts directly to target outputs.
- **Alternatives considered**:
  - Single merged artifact package (rejected: harder target-level retrieval).
  - OS-level grouped artifacts (rejected: weaker target precision).

## Decision 4: Use Go version from `go.mod`

- **Decision**: Configure `actions/setup-go` with `go-version-file: go.mod`.
- **Rationale**: Keeps CI toolchain aligned with repository-declared Go version and reduces version drift.
- **Alternatives considered**:
  - Hardcoded Go version in workflow (rejected: higher maintenance and drift risk).

## Decision 5: Include workflow dispatch and tag-triggered release pathways

- **Decision**: Support automatic runs on main pushes, pull requests, tag pushes (`v*`), and manual dispatch.
- **Rationale**: Covers continuous validation, release-like events, and on-demand operational reruns.
- **Alternatives considered**:
  - Only push/PR triggers (rejected: misses release and manual operations use cases).
  - Manual-only release workflows (rejected: weak automation guarantees).

## Decision 6: Set artifact retention to 14 days

- **Decision**: Use 14-day artifact retention for build outputs.
- **Rationale**: Balances retrievability for debugging/release workflows with storage control.
- **Alternatives considered**:
  - 7 days (rejected: potentially too short for some review cycles).
  - Default retention (rejected: less explicit storage governance).
