# CLI Contract: Configuration and CLI Interface

## Command

`sql-proxy [flags]`

## Flags

| Flag | Alias | Type | Required | Description |
|------|-------|------|----------|-------------|
| `--port` | `-p` | integer | No | Local listener port. Valid range: `1-65535`. |
| `--instance` | `-i` | string | Yes* | Target Cloud SQL instance identifier. Required after config resolution. |

\* `--instance` may be omitted only if provided in `~/.sql-proxy/config.yaml`.

## Config File Contract

- **Path**: `~/.sql-proxy/config.yaml`
- **Format**: YAML
- **Supported keys**:
  - `port` (integer)
  - `instance` (string)

## Precedence Contract

1. CLI flags
2. Config file values
3. Defaults

Default values:
- `port: 5432`
- `instance`: no default (must be provided by flag or config file)

## Error Contract

All startup validation failures return:
- exit code `1`
- user-friendly message with a corrective action

Error categories:
- malformed config file
- missing required instance
- invalid port (outside `1-65535`)
- port conflict on `127.0.0.1:<port>`

## Signal Handling Contract

- Handle `SIGINT` and `SIGTERM`.
- On signal: close local listener, cancel active context, and exit cleanly.
