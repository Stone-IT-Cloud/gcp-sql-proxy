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
