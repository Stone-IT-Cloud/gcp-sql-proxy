# Contract: Repository Pre-commit Quality Configuration

## Scope

Defines mandatory repository artifacts and behavioral expectations for local commit validation.

## Required Artifacts

1. `.pre-commit-config.yaml` at repository root.
2. `.golangci.yml` at repository root.
3. Contributor setup documentation in `README.md` (or equivalent docs entry).

## Configuration Contract

### `.pre-commit-config.yaml`

- MUST define `default_stages` including `commit`.
- MUST include a repo entry for `pre-commit-hooks` with:
  - `trailing-whitespace`
  - `end-of-file-fixer`
  - `check-yaml`
- MUST include a repo entry for `golangci-lint` with `golangci-lint` hook.
- MUST include a `repo: local` entry that defines:
  - `go-mod-tidy`
  - `go-test` (`go test -short ./...`)
- MUST enforce commit failure when any hook fails.
- MUST scope Go-specific hooks to staged Go-related files where applicable.

### `.golangci.yml`

- MUST enable at least:
  - `errcheck`
  - `gosimple`
  - `govet`
  - `ineffassign`
  - `staticcheck`
  - `gofmt`

## Behavioral Contract

1. If a required hook tool is missing locally, commit validation MUST fail with actionable install guidance.
2. If no staged Go-related files exist, baseline hooks still run while Go-specific hooks are skipped.
3. Hook feedback MUST identify failing checks clearly enough for contributor self-remediation.
4. Validation runtime should satisfy project target (95% of commit attempts complete within 120 seconds in standard contributor environments).

## Verification Contract

The feature is considered conformant when:

1. Configuration files exist with required entries.
2. A commit with intentionally malformed YAML or whitespace issues is blocked.
3. A commit introducing a Go lint or short-test failure is blocked when Go files are staged.
4. A non-Go-only commit does not run Go-specific hooks but still enforces baseline hooks.
