# Feature Specification: Configuration and CLI Interface

**Feature Branch**: `001-config-cli-interface`
**Created**: 2026-05-05
**Status**: Draft
**Input**: User description: "Spec 01: Configuration and CLI Interface"

## Clarifications

### Session 2026-05-05

- Q: How should the app behave when `config.yaml` is missing vs malformed? → A: Missing file continues with defaults; malformed file exits with clear error and code 1.
- Q: How should the app behave when no instance is provided by flag or config? → A: Exit with a clear user-facing error and code 1.
- Q: What valid local port range should be accepted? → A: Validate and allow ports 1 through 65535; otherwise fail with clear error and code 1.
- Q: How should `.sql-proxy` directory permissions behave cross-platform? → A: Use mode `0755` on Unix-like systems and platform-equivalent defaults on Windows.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Start Proxy with Predictable Configuration (Priority: P1)

As a CLI user, I want the proxy to load settings from flags, config file, and defaults in a
predictable order so I can start a connection without manually setting every option each time.

**Why this priority**: Reliable startup behavior is the foundation for all other proxy usage.

**Independent Test**: Launch with different combinations of flag values, config file values, and
missing values; verify the resolved configuration follows precedence and the proxy starts with the
expected local port and target instance.

**Acceptance Scenarios**:

1. **Given** no config file and no CLI flags, **When** the user starts the app, **Then** the app
   uses default values and starts with local port `5432`.
2. **Given** a config file with `port: 6000`, **When** the user starts the app without `--port`,
   **Then** the app uses port `6000`.
3. **Given** a config file with `port: 6000` and CLI flag `--port 7000`, **When** the user starts
   the app, **Then** the app uses port `7000`.

---

### User Story 2 - Recover from Port Conflict Quickly (Priority: P2)

As a CLI user, I want a clear error when the selected local port is unavailable so I can fix the
issue immediately and retry.

**Why this priority**: Port conflicts are common in local development and can block usage.

**Independent Test**: Reserve a local port with another process, run the app on that port, and
verify it exits with code `1` and prints recovery guidance.

**Acceptance Scenarios**:

1. **Given** another process is already bound to `127.0.0.1:5432`, **When** the user starts the
   app using port `5432`, **Then** the app exits with code `1` and shows a user-friendly message
   describing the conflict and how to change the port via flag or config file.

---

### User Story 3 - Stop Cleanly on Termination Signals (Priority: P3)

As an operator, I want the proxy to shut down cleanly on interruption signals so local resources
and upstream sessions are not left in an inconsistent state.

**Why this priority**: Graceful shutdown improves reliability and avoids lingering listeners.

**Independent Test**: Start the app, send `SIGINT`/`SIGTERM`, and verify the listener closes and
the active operation context is cancelled before process exit.

**Acceptance Scenarios**:

1. **Given** the proxy is running, **When** a termination signal is received, **Then** the app
   closes the local listener, cancels active connection work, and exits cleanly.

---

### Edge Cases

- Home directory exists but the `.sql-proxy` folder does not exist yet.
- Config file is missing and defaults must be used without failing startup.
- Config file path exists but file is unreadable or malformed and startup must fail with clear guidance.
- User passes an invalid port value outside accepted range.
- User provides only one of required startup inputs (for example, missing instance identifier).
- A termination signal arrives while startup is still initializing resources.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST resolve the current user's home directory on all supported platforms.
- **FR-002**: System MUST create a `.sql-proxy` directory under the home directory when missing.
- **FR-002a**: On Unix-like systems, system MUST create `.sql-proxy` with directory mode `0755`.
- **FR-002b**: On Windows, system MUST apply platform-equivalent default permissions and continue.
- **FR-003**: System MUST use `~/.sql-proxy/config.yaml` as the default configuration file path.
- **FR-004**: System MUST apply default configuration values when no explicit values are provided.
- **FR-004a**: System MUST continue startup when the configuration file is missing by using defaults and
  any provided command-line inputs.
- **FR-004b**: System MUST fail startup with a clear user-facing message and exit code `1` when the
  configuration file exists but cannot be parsed.
- **FR-005**: System MUST support command-line inputs for local port and target instance identifier.
- **FR-005a**: System MUST require a target instance identifier from command-line inputs or configuration
  file before starting proxy initialization.
- **FR-005b**: System MUST validate local port input and accept only values in range `1-65535`.
- **FR-006**: System MUST resolve configuration precedence as:
  command-line inputs > configuration file values > defaults.
- **FR-007**: System MUST verify local port availability before initializing proxy connection flow.
- **FR-008**: System MUST return a user-friendly port conflict error without panic when the selected
  port is already in use.
- **FR-009**: System MUST include recovery guidance in port conflict errors, including how to change
  the port through command-line input or configuration file updates.
- **FR-010**: System MUST exit with status code `1` on port conflict.
- **FR-010a**: System MUST exit with status code `1` when required startup inputs are missing.
- **FR-010b**: System MUST exit with status code `1` when port validation fails.
- **FR-011**: System MUST listen for process termination signals and trigger graceful shutdown.
- **FR-012**: System MUST close the local TCP listener and cancel active connection operations before
  process exit when termination is requested.

### Key Entities *(include if feature involves data)*

- **Runtime Configuration**: Effective startup settings resolved from defaults, config file, and
  command-line inputs, including local port and target instance.
- **Configuration File**: Persistent user-defined values stored in `~/.sql-proxy/config.yaml`.
- **Listener Session**: Active local TCP listener state and lifecycle events (started, conflict,
  shutting down, closed).
- **Shutdown Signal Event**: External process signal that triggers coordinated teardown behavior.

## Non-Functional Requirements *(mandatory)*

- **NFR-001 (Cross-Platform)**: Feature MUST work on Windows, macOS, and Linux.
- **NFR-002 (Error UX)**: Error messages for startup failure MUST be understandable by non-expert
  users and include corrective action.
- **NFR-003 (Responsiveness)**: The app MUST detect and report port conflict immediately at startup
  before expensive connection setup begins.
- **NFR-004 (Reliability)**: Shutdown handling MUST release local listener resources on every
  termination path.
- **NFR-005 (Security)**: Local configuration directory and files MUST follow restrictive permission
  practices compatible with the platform, including `0755` for `.sql-proxy` on Unix-like systems.
- **NFR-006 (Security Baseline Preservation)**: This feature MUST NOT weaken existing IAM
  authentication defaults or private tunnel defaults used by the proxy runtime.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In validation tests, 100% of startup attempts resolve configuration using documented
  precedence rules.
- **SC-002**: In conflict tests, 100% of port-in-use launches fail fast with exit code `1` and
  include actionable guidance.
- **SC-003**: During shutdown tests, 100% of signal-triggered terminations release the listener and
  complete shutdown within 2 seconds.
- **SC-004**: In usability checks, at least 90% of users can correct a port conflict and relaunch
  successfully on the first retry.
- **SC-005**: In local startup checks, configuration validation (including malformed config detection,
  required input validation, and port availability check) completes in under 1 second.

## Assumptions

- The application runs in user contexts where home directory resolution is available.
- The feature scope includes only startup configuration and shutdown behavior, not full proxy data
  transfer or authentication flows.
- A valid target instance identifier is provided by users for normal operation, with format
  validation handled by existing or subsequent features; missing values are treated as startup errors.
- Platform-specific permission semantics may differ, but equivalent restrictive behavior is expected.
