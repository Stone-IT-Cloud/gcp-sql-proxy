# Quickstart: GitHub Actions Automated Build Pipeline

## 1. Validate workflow file placement

- Confirm workflow exists at `.github/workflows/build.yml`.
- Confirm workflow name is `Build and Release`.

## 2. Verify trigger coverage

Check workflow trigger configuration includes:

- push to `main`
- pull request targeting `main`
- push tags matching `v*`
- manual `workflow_dispatch`

## 3. Validate matrix target coverage

Confirm build matrix includes all required targets:

- linux/amd64
- windows/amd64
- darwin/amd64
- darwin/arm64

And confirm runner selection follows matrix-provided `os` values.

## 4. Validate build output behavior

For each matrix target:

- build command uses `go build`
- `GOOS` and `GOARCH` are set from matrix fields
- binary output follows `gcp-db-proxy-<os>-<arch>`
- Windows output adds `.exe`
- linker flags include version metadata when available

## 5. Validate artifact upload behavior

- `actions/upload-artifact@v4` is used.
- one artifact is uploaded per matrix target.
- artifact names match target binary names.
- retention is set to 14 days.

## 6. Functional workflow checks

Execute or inspect runs for:

1. Pull request to `main`
2. Push to `main`
3. Push tag `vX.Y.Z`
4. Manual workflow dispatch

Expected outcomes:

- All runs start automatically for event-driven triggers.
- Matrix runs preserve per-target results.
- If one target fails, remaining targets still execute and workflow concludes failed overall.
- Successful targets always produce downloadable artifacts.

## 7. Local repository validation

Run the repository test suite before opening a PR:

```bash
go test ./...
```

Record the command result in PR notes or implementation summary.

## 8. Artifact retrieval timing check

For at least 10 successful workflow runs, measure time from opening the run page
to downloading a target artifact.

Success threshold:
- At least 95% of measured retrievals complete within 2 minutes.
