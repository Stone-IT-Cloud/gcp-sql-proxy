# Data Model: Pre-commit Code Quality Gates

## Entity: PreCommitConfig

- **Description**: Repository root pre-commit manifest defining stages, hook repositories, and hook entries.
- **Key Fields**:
  - `default_stages`: list of commit lifecycle stages where hooks run.
  - `repos`: ordered list of hook repository declarations.
- **Validation Rules**:
  - Must include baseline hygiene hooks.
  - Must include golangci-lint hook configuration.
  - Must include local Go hooks for module hygiene and short tests.
  - Must ensure Go-specific hooks are scoped to staged Go-related files where required.

## Entity: HookDefinition

- **Description**: Single executable pre-commit check or formatter with pass/fail behavior.
- **Key Fields**:
  - `id`: unique hook identifier.
  - `entry`: command to execute (local hooks).
  - `language`: hook runtime mode.
  - `types`/`files`: file-match scope controlling execution.
  - `pass_filenames`: indicates whether staged filenames are forwarded to command.
- **Validation Rules**:
  - Missing required tool dependencies cause hook failure with actionable output.
  - Hooks that modify files must signal failure until contributor re-stages changes.

## Entity: LintPolicy

- **Description**: Repository lint policy configuration for Go static analysis.
- **Key Fields**:
  - `linters.enable`: list of active Go linters.
  - Optional execution tuning fields (timeout/concurrency) if needed later.
- **Validation Rules**:
  - Must enable minimum required linters: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, and `gofmt`.

## Entity: ContributorSetupGuide

- **Description**: Documentation content that defines required local setup and troubleshooting for commit hooks.
- **Key Fields**:
  - prerequisites
  - installation steps
  - activation command
  - verification command
  - common failure remediation
- **Validation Rules**:
  - Must include instructions sufficient for setup completion within 10 minutes.
  - Must include explicit guidance for missing-tool failures and how to resolve them.

## Relationships

- `PreCommitConfig` contains multiple `HookDefinition` entries.
- `PreCommitConfig` references `LintPolicy` through golangci-lint integration.
- `ContributorSetupGuide` documents how contributors satisfy `PreCommitConfig` prerequisites and resolve `HookDefinition` failures.
