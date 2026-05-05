# Contract: GitHub Actions Cross-Platform Build Workflow

## Scope

Defines the required behavior and structure for repository build automation at `.github/workflows/build.yml`.

## Required Workflow Contract

1. Workflow file location MUST be `.github/workflows/build.yml`.
2. Workflow name MUST be `Build and Release`.
3. Trigger contract MUST include:
   - `push` on branch `main`
   - `pull_request` targeting branch `main`
   - `push` on tags matching `v*`
   - `workflow_dispatch`

## Build Job Contract

1. A `build` job MUST exist.
2. Job execution MUST use matrix-driven target definitions with dynamic runner assignment.
3. Required matrix target combinations:
   - Linux amd64
   - Windows amd64
   - Darwin amd64
   - Darwin arm64
4. Matrix failure behavior MUST continue remaining targets and report workflow failure if any target fails.

## Build Step Contract

1. Repository checkout MUST use `actions/checkout@v4` (or newer compatible major).
2. Go setup MUST use `actions/setup-go@v5` with `go-version-file: go.mod`.
3. Build command MUST use `go build`.
4. Build command MUST pass `GOOS` and `GOARCH` from matrix target fields.
5. Build output naming MUST be `gcp-db-proxy-<target_os>-<target_arch>` and append `.exe` for Windows.
6. Linker flags MUST support version metadata injection when source version context is present.

## Artifact Contract

1. Artifact upload MUST use `actions/upload-artifact@v4` (or newer compatible major).
2. Exactly one artifact MUST be uploaded per matrix target.
3. Artifact name MUST match the binary naming pattern for target clarity.
4. `retention-days` MUST be explicitly configured to 14.

## Verification Contract

The workflow is conformant when:

1. All required triggers are present and valid.
2. A workflow run on pull request, push to main, tag push, and manual dispatch starts successfully.
3. Successful matrix targets produce correctly named binaries.
4. Artifact uploads are present per target with configured retention.
5. A simulated failing target still allows other targets to complete while workflow ends as failed overall.
6. Failed target logs are clearly separated by matrix target so maintainers can identify the failing OS/architecture without ambiguity.
