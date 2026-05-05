# gcp-db-proxy (`gcp-sql-proxy`)

CLI proxy that listens on localhost and tunnels PostgreSQL connections to [Google Cloud SQL](https://cloud.google.com/sql) using the Cloud SQL Go connector with IAM database authentication. Use it when your application expects a TCP host/port on your machine rather than connecting to the instance IP directly.

**Requirements:**

- Go **1.25.7+** if you build from source (see [`go.mod`](./go.mod))
- A GCP project with a Cloud SQL **PostgreSQL** instance
- An **OAuth 2.0 Client ID** (Desktop app type) embedded at build time (see below)
- Your Google account added as an OAuth **test user** while the OAuth consent screen is in **Testing**, or the app published appropriately

---

## Table of contents

- [Installation](#installation)
- [Configuration](#configuration)
- [OAuth credentials and building](#oauth-credentials-and-building)
- [Running the proxy](#running-the-proxy)
- [Connecting PostgreSQL clients](#connecting-postgresql-clients)
- [CI binaries and GitHub Actions secrets](#ci-binaries-and-github-actions-secrets)
- [Contributing](#contributing)
- [Pre-commit quality gates](#pre-commit-quality-gates)
- [Build and release workflow](#build-and-release-workflow)

---

## Installation

### Option A: Download a CI-built binary (recommended for users)

1. Open the repository’s **Actions** tab and select the latest **Build and Release** run (or a successful run for your branch).
2. Download the workflow **artifact** for your OS and CPU (for example `gcp-db-proxy-darwin-arm64`).
3. Unpack if needed, rename or move the binary to `gcp-db-proxy`, and ensure it is executable:

   ```bash
   chmod +x gcp-db-proxy
   ```

CI builds embed OAuth client ID and secret from repository secrets `API_KEY_ID` and `API_KEY_SECRET` (see [CI binaries](#ci-binaries-and-github-actions-secrets)). If those secrets are not set on a fork, binaries will lack OAuth credentials and the proxy will **not** complete authentication.

### Option B: Install from source (developers)

```bash
git clone https://github.com/Stone-IT-Cloud/gcp-sql-proxy.git
cd gcp-sql-proxy
```

Build without OAuth (useful only for developing non-auth paths):

```bash
go build -o gcp-db-proxy ./cmd/sql-proxy
```

Build **with OAuth** embedded (needed for normal tunnel startup):

```bash
GOOS=darwin GOARCH=arm64 go build \
  -ldflags "-X 'github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth.OAuthClientID=YOUR_CLIENT_ID' -X 'github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth.OAuthClientSecret=YOUR_CLIENT_SECRET'" \
  -o gcp-db-proxy ./cmd/sql-proxy
```

Adjust `GOOS`/`GOARCH` for your platform (`linux/amd64`, `windows/amd64`, etc.).

---

## Configuration

### Config file location

Defaults and instance can live in **`~/.sql-proxy/config.yaml`** (created next to the config directory on first run).

Example:

```yaml
port: 5432
instance: your-project:us-east1:your-instance
private_ip: false   # optional; default is public IP
```

### Command-line flags (override file)

| Flag | Short | Description |
|------|-------|-------------|
| `--instance` | `-i` | Cloud SQL **connection name**: `PROJECT_ID:REGION:INSTANCE_ID` |
| `--port` | `-p` | Local TCP port (default `5432`) |
| `--public-ip` | | Use **public** IP for the connector (default if neither IP flag is set) |
| `--private-ip` | | Use **private** IP (requires network path to the instance private IP) |

Do **not** pass both `--public-ip` and `--private-ip`.

### Verifying the instance name

```bash
gcloud sql instances describe INSTANCE_NAME --project=PROJECT_ID --format='value(connectionName)'
```

Use that `connectionName` value as `--instance`.

---

## OAuth credentials and building

1. In **Google Cloud Console** → **APIs & Services** → **Credentials**, create an OAuth client of type **Desktop app**.
2. Configure the **OAuth consent screen**; while in **Testing**, add every Google account that will sign in under **Test users**.
3. Embed **Client ID** and **Client secret** at build time with `-ldflags` (see [Installation](#installation)).

On first successful run with valid credentials, the binary stores tokens under **`~/.sql-proxy/token.json`** with restrictive permissions where supported.

Without embedded OAuth credentials the process may listen on localhost but **not** dial Cloud SQL; use a binary built with `OAuthClientID` and `OAuthClientSecret` set.

---

## Running the proxy

Examples:

```bash
# Default: public IP to Cloud SQL
./gcp-db-proxy --instance your-project:us-east1:your-instance

# Explicit public IP (same behavior as default)
./gcp-db-proxy --instance your-project:us-east1:your-instance --public-ip

# Private IP path (VPC / hybrid connectivity required)
./gcp-db-proxy --instance your-project:us-east1:your-instance --private-ip

# Custom local port
./gcp-db-proxy --instance your-project:us-east1:your-instance --port 15432
```

### Startup output

When the tunnel path is established, stdout includes:

- Confirmation that the tunnel is ready (and again on **each successful reconnect** to a client session)
- **IP connectivity mode** (`public` or `private`)
- Target instance connection name
- Local host and port
- A **`psql`-style template** using placeholders (`<db_name>`, `<db_user>`, `<db_password>`) — **credentials are never printed**

If OAuth is missing while the listener still binds, stdout explains that the listener is active but the tunnel is not established.

Graceful shutdown: **Ctrl+C** (SIGINT/SIGTERM).

---

## Connecting PostgreSQL clients

The proxy uses **IAM authentication** toward Cloud SQL for the tunnel. Database **roles and passwords** are still PostgreSQL-owned: connect with an existing Postgres user your team created (often a built-in Cloud SQL user), not necessarily your Gmail address unless you explicitly use IAM DB users.

Typical **`psql`** through the tunnel (local host/port from the startup message):

```bash
PGPASSWORD='your-database-password' psql \
  "host=127.0.0.1 port=5432 dbname=your_database user=your_db_user sslmode=disable"
```

If authentication fails:

- Confirm `dbname`, `user`, and `password` against your app’s secrets or `gcloud sql users list`.
- Prefer **IAM auth for the connector**, **password auth for Postgres** unless your org uses IAM DB users end-to-end.

---

## CI binaries and GitHub Actions secrets

The workflow `.github/workflows/build.yml` builds cross-platform artifacts and passes linker flags for version/commit and OAuth variables.

Maintain these **repository secrets** (Settings → Secrets and variables → Actions):

| Secret | Purpose |
|--------|---------|
| `API_KEY_ID` | OAuth 2.0 Client ID embedded as `internal/auth.OAuthClientID` |
| `API_KEY_SECRET` | OAuth client secret embedded as `internal/auth.OAuthClientSecret` |

### Triggers

- Push to **`main`**
- Pull requests targeting **`main`**
- Tags matching **`v*`**
- **workflow_dispatch** (manual run)

Artifacts are retained **14 days**; matrix jobs use **`fail-fast: false`** so all targets stay visible even if one fails.

---

## Contributing

1. **Fork / branch** — Follow your team’s branching rules; feature work often uses spec-driven branches (`specs/00x-...`).
2. **Issues and specs** — Large changes align with specs under `specs/` and project rules in [`.cursor/rules/`](./.cursor/rules/).
3. **Commit style** — Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `docs:`, `chore:`, …).
4. **Tests** — Run the full suite before opening a PR:

   ```bash
   go test ./...
   ```

5. **Quality** — Enable [pre-commit](#pre-commit-quality-gates); PRs should pass lint/tests and avoid logging or printing secrets.

### Project layout

- `cmd/sql-proxy/` — CLI entrypoint
- `internal/auth/` — OAuth and Cloud SQL token sources
- `internal/config/` — Flags and `~/.sql-proxy/config.yaml`
- `internal/proxy/` — Access checks, connector dial plan, relay, tunnel messages
- `tests/unit/`, `tests/integration/` — Go tests

---

## Pre-commit quality gates

This repository uses [pre-commit](https://pre-commit.com/) for formatting, linting, and fast tests.

### Prerequisites

- Python 3 and `pre-commit`
- Go on `PATH`
- `golangci-lint` and `goimports` on `PATH`

```bash
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Setup

```bash
pre-commit install --hook-type pre-commit
pre-commit run --all-files
```

- Baseline hooks: whitespace, EOF, YAML.
- Go hooks run when Go or module files are touched.
- **Missing required tools are treated as failures** — install `goimports` and `golangci-lint` (see Prerequisites) and rerun hooks.

---

## Build and release workflow

- Workflow file: `.github/workflows/build.yml`
- **Targets:** `linux/amd64`, `windows/amd64`, `darwin/amd64`, `darwin/arm64`
- **Binary names:** `gcp-db-proxy-<os>-<arch>` (Windows adds `.exe`)
- **Linker flags:** `main.version`, `main.commit`, and OAuth fields as documented above

For local release-style builds matching CI:

```bash
VERSION="$(git describe --tags --always)"
COMMIT="$(git rev-parse HEAD)"
go build \
  -ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} \
    -X 'github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth.OAuthClientID=YOUR_ID' \
    -X 'github.com/Stone-IT-Cloud/gcp-sql-proxy/internal/auth.OAuthClientSecret=YOUR_SECRET'" \
  -o gcp-db-proxy ./cmd/sql-proxy
```

---

## Security notes

- **Never commit** OAuth client secrets or production database passwords.
- **Rotate** any secret that appeared in logs, chat, or CI output.
- Repository **`.gitignore`** includes `.env*` and common editor artifacts; keep local credentials out of version control.
