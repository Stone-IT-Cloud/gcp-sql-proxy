# Research: Configuration and CLI Interface

## Decision 1: Configuration Precedence with Viper and pflag

- **Decision**: Use Viper defaults for baseline values, load `config.yaml` if present, then bind and
  parse pflags so final effective values are resolved as `flags > config file > defaults`.
- **Rationale**: This directly matches feature requirements and keeps behavior deterministic for users.
- **Alternatives considered**:
  - Manual merge logic without Viper: rejected due to higher maintenance and parsing complexity.
  - Flag-only configuration: rejected because persistent user preferences are required.

## Decision 2: Missing vs Malformed Config Behavior

- **Decision**: Missing config file is non-fatal and startup continues with defaults/flags; malformed
  config is fatal with clear error message and exit code `1`.
- **Rationale**: Missing file is an expected first-run state, while malformed config indicates a user
  action that must be corrected to avoid hidden misconfiguration.
- **Alternatives considered**:
  - Always fallback on malformed config: rejected because it can hide configuration mistakes.
  - Auto-rewrite malformed config: rejected to avoid destructive behavior on user files.

## Decision 3: Port Validation and Conflict Detection

- **Decision**: Validate ports strictly in range `1-65535`, then perform an early bind check on
  `127.0.0.1:<port>` before proxy initialization.
- **Rationale**: Input validation and preflight binding provide fast, actionable feedback and prevent
  wasted setup work.
- **Alternatives considered**:
  - Rely only on listen errors without range validation: rejected due to inconsistent UX.
  - Block privileged ports by default: rejected because requirement allows full valid range.

## Decision 4: Signal-Driven Graceful Shutdown

- **Decision**: Use signal notifications for `SIGINT` and `SIGTERM`, trigger listener closure and
  context cancellation, and terminate cleanly.
- **Rationale**: Ensures predictable cleanup and aligns with reliability requirements.
- **Alternatives considered**:
  - Immediate exit on signal: rejected because it can leak resources and skip cleanup.
  - Polling loop shutdown checks: rejected due to unnecessary complexity.

## Decision 5: Cross-Platform Directory Permissions

- **Decision**: Create `~/.sql-proxy` with mode `0755` on Unix-like systems; use platform-equivalent
  defaults on Windows while preserving secure practical behavior.
- **Rationale**: Matches clarified spec while respecting OS permission model differences.
- **Alternatives considered**:
  - Enforce exact POSIX bits on Windows: rejected as non-portable and unreliable.
  - Always use OS defaults on all platforms: rejected because Unix mode expectation is explicit.
