# Quickstart: OAuth2 Automated Authentication

## Prerequisites

- Feature branch: `002-oauth2-automated-auth`
- OAuth desktop client credentials available to the app runtime
- Browser available (or ability to manually open URL)

## 1) First-run interactive authentication

1. Ensure `~/.sql-proxy/token.json` does not exist.
2. Start the proxy.
3. Confirm browser consent URL is opened automatically (or displayed for manual open).
4. Complete consent.
5. Confirm browser callback displays success message.
6. Confirm token file is created at `~/.sql-proxy/token.json` with restricted permissions.

## 2) Persisted-token reuse

1. Keep valid token file in place.
2. Restart proxy.
3. Confirm no browser consent flow is triggered.
4. Confirm authenticated client is created and passed to dialer path.

## 3) Invalid-token recovery

1. Corrupt token file contents.
2. Start proxy.
3. Confirm invalid token file is renamed or removed.
4. Confirm fresh browser auth flow is triggered and new token is persisted.

## 4) Callback port conflict recovery

1. Bind another process to localhost `8080`.
2. Start proxy without valid token.
3. Confirm auth flow still proceeds using fallback localhost callback port.

## 5) Callback state validation

1. Trigger OAuth flow and simulate callback with missing or mismatched state.
2. Confirm authentication fails with clear security-related error.
3. Confirm no valid authenticated client is produced.

## 6) Run tests

```bash
go test ./...
```

## 7) Validate credential injection behavior

1. Start the app without OAuth client credentials configured.
2. Confirm startup path remains available and does not crash.
3. Configure OAuth client credentials and rerun to validate interactive auth behavior.
