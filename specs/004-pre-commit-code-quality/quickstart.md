# Quickstart: Pre-commit Code Quality Gates

## 1. Install prerequisites

- Install Python + `pre-commit`.
- Ensure Go toolchain is installed and available in `PATH`.
- Ensure `golangci-lint` is available in `PATH` (or install via project-approved method).

## 2. Install Git hooks

From repository root:

```bash
pre-commit install --hook-type pre-commit
```

## 3. Validate configuration

Run all hooks against the repository:

```bash
pre-commit run --all-files
```

Expected outcome:
- Baseline file hygiene hooks run successfully.
- GolangCI-Lint executes for Go-related scope.
- Local hooks run `go mod tidy` and `go test -short ./...` according to configured file targeting.

## 4. Verify commit-time behavior

1. Create a controlled formatting or YAML issue and confirm commit is blocked.
2. Create a controlled Go lint or short-test failure and confirm commit is blocked when Go files are staged.
3. Attempt a non-Go-only commit and confirm baseline hooks run while Go-specific hooks are skipped.

## 5. Troubleshooting

- If commit fails due to missing tooling, install the missing dependency and re-run:

```bash
pre-commit run --all-files
```

- If a hook auto-fixes files, re-stage the modified files before committing again.
