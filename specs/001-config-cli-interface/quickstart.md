# Quickstart: Configuration and CLI Interface

## Prerequisites

- Go installed and available in PATH
- Repository checked out on branch `001-config-cli-interface`

## 1) Prepare directories

```bash
mkdir -p cmd/sql-proxy internal/config tests/unit tests/integration
```

## 2) Add dependencies

```bash
go get github.com/spf13/viper github.com/spf13/pflag
```

## 3) Run with defaults

```bash
go run ./cmd/sql-proxy --instance "project:region:instance"
```

Expected behavior:
- Creates `~/.sql-proxy` if missing
- Uses `~/.sql-proxy/config.yaml` when present
- Uses default port `5432` if not provided elsewhere

## 4) Verify precedence

```bash
go run ./cmd/sql-proxy --instance "project:region:instance" --port 7000
```

Expected behavior:
- CLI flag value overrides config file and defaults

## 5) Validate error flows

- **Missing instance**: run without `--instance` and without config value; expect clear error + exit code `1`.
- **Invalid port**: run `--port 70000`; expect clear error + exit code `1`.
- **Port conflict**: occupy target port with another process and start app; expect clear conflict guidance + exit code `1`.
- **Malformed config**: corrupt `~/.sql-proxy/config.yaml`; expect clear parse error + exit code `1`.

## 6) Validate graceful shutdown

1. Start the proxy process.
2. Send `SIGINT` (`Ctrl+C`) or `SIGTERM`.
3. Confirm listener closes and process exits cleanly.

## 7) Run tests

```bash
go test ./...
```

## 8) Validate security baseline preservation

- Run integration tests and confirm no startup/config path introduces IAM authentication or private
  tunnel override keys.
