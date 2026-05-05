# gcp-sql-proxy

## Pre-commit Quality Gates

This repository uses [pre-commit](https://pre-commit.com/) to enforce formatting,
linting, and fast unit tests before every commit.

### Prerequisites

- Python 3 with `pre-commit` installed
- Go toolchain available in `PATH`
- `golangci-lint` and `goimports` available in `PATH`

Install missing Go tools:

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Setup

Install hooks once per clone:

```bash
pre-commit install --hook-type pre-commit
```

Validate full hook execution:

```bash
pre-commit run --all-files
```

### Hook Behavior

- Baseline hygiene hooks check whitespace/newline/YAML validity.
- Go-specific hooks run only when staged files include Go source/module files.
- Any failed hook blocks commit creation.
- Missing required tools are treated as failures; install the missing tool and rerun.

If hooks auto-fix files, re-stage changed files before committing again.

## Build and Release Workflow

The repository includes a cross-platform GitHub Actions build workflow at
`.github/workflows/build.yml`.

### Trigger Paths

- Push to `main`
- Pull requests targeting `main`
- Tag pushes matching `v*`
- Manual dispatch (`workflow_dispatch`)

### Build Outputs

- Matrix targets: `linux/amd64`, `windows/amd64`, `darwin/amd64`, `darwin/arm64`
- Binary naming format: `gcp-db-proxy-<target_os>-<target_arch>`
- Windows binaries include `.exe`
- Build metadata is injected through linker flags when run context provides version values

### Artifacts

- One uploaded artifact per matrix target
- Artifact name matches target binary naming
- Retention: 14 days

### Failure Visibility

- Matrix jobs continue even when one target fails (`fail-fast: false`)
- Workflow result is marked failed if any target fails
