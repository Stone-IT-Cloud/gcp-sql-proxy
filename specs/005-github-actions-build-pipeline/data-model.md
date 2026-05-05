# Data Model: GitHub Actions Automated Build Pipeline

## Entity: WorkflowDefinition

- **Description**: CI configuration that defines event triggers, build jobs, matrix strategy, and artifact upload policy.
- **Key Fields**:
  - `name`
  - `on` trigger set (`push`, `pull_request`, `workflow_dispatch`)
  - `jobs.build`
  - `jobs.build.strategy.matrix`
- **Validation Rules**:
  - Must include required triggers and branch/tag filters.
  - Must define build job using required action dependencies.

## Entity: BuildTarget

- **Description**: A single matrix entry representing one platform compilation target.
- **Key Fields**:
  - `os` (runner label)
  - `target_os` (GOOS)
  - `target_arch` (GOARCH)
- **Validation Rules**:
  - Must include the required four target combinations.
  - Must map consistently from runner to binary output naming.

## Entity: BinaryOutput

- **Description**: A compiled executable generated for a specific build target.
- **Key Fields**:
  - `base_name` (`gcp-db-proxy`)
  - `target_os`
  - `target_arch`
  - `extension` (`.exe` for Windows only)
- **Validation Rules**:
  - Filename must include OS and architecture.
  - Windows target output must include `.exe`.
  - Output naming must be collision-free across matrix targets.

## Entity: TargetArtifactPackage

- **Description**: Uploaded artifact representing exactly one target binary.
- **Key Fields**:
  - `artifact_name`
  - `artifact_path`
  - `retention_days`
- **Validation Rules**:
  - Must be one artifact per matrix target.
  - Retention days must be explicitly set and bounded.

## Entity: VersionMetadata

- **Description**: Build-time version information injected through linker flags.
- **Key Fields**:
  - `version_value`
  - `source_context` (tag, commit reference, or workflow metadata)
- **Validation Rules**:
  - Build should inject metadata when version context is present.
  - Metadata injection must not expose secrets.

## Relationships

- `WorkflowDefinition` contains multiple `BuildTarget` entries through matrix strategy.
- Each `BuildTarget` produces one `BinaryOutput`.
- Each `BinaryOutput` maps to one `TargetArtifactPackage`.
- `VersionMetadata` is applied during `BinaryOutput` creation when available.
