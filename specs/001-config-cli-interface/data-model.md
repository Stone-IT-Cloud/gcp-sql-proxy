# Data Model: Configuration and CLI Interface

## Entity: RuntimeConfig

Represents effective startup configuration after precedence resolution.

### Fields

- `port` (integer): local TCP port used for listener bind.
- `instance` (string): Cloud SQL instance identifier.
- `config_file_path` (string): absolute path to active config file.
- `source_port` (enum: `default`, `config`, `flag`): origin of effective port value.
- `source_instance` (enum: `config`, `flag`): origin of effective instance value.

### Validation Rules

- `port` MUST be in range `1-65535`.
- `instance` MUST be non-empty at startup.
- `config_file_path` MUST point to `~/.sql-proxy/config.yaml`.

## Entity: ConfigDirectory

Represents managed local settings directory.

### Fields

- `path` (string): `~/.sql-proxy`.
- `exists` (boolean): whether directory existed before startup.
- `permission_mode` (string): expected mode (`0755` on Unix-like systems).

### Validation Rules

- Directory MUST be created if missing.
- Unix-like systems MUST apply mode `0755`.
- Windows MUST apply platform-equivalent defaults and continue.

## Entity: ListenerSession

Represents local bind and shutdown lifecycle.

### Fields

- `address` (string): `127.0.0.1:<port>`.
- `state` (enum): `initializing`, `listening`, `conflict`, `shutting_down`, `closed`.
- `exit_code` (integer): process exit code for terminal state.

### State Transitions

- `initializing -> listening` when port bind succeeds.
- `initializing -> conflict` when bind fails due to port in use.
- `listening -> shutting_down` when signal is received.
- `shutting_down -> closed` after listener close + context cancellation complete.

## Entity: StartupError

Represents user-visible startup validation failures.

### Fields

- `kind` (enum): `missing_instance`, `invalid_port`, `malformed_config`, `port_conflict`.
- `message` (string): user-facing actionable error text.
- `exit_code` (integer): always `1` for this feature scope.

### Validation Rules

- Messages MUST include clear remediation guidance.
- Errors MUST NOT expose stack traces or panic output.
