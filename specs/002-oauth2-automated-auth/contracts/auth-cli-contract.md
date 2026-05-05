# CLI Auth Contract: OAuth2 Automated Authentication

## Purpose

Define required behavior for automated OAuth2 desktop authentication in the SQL proxy CLI.

## Inputs

- OAuth client credentials (injected by runtime config/environment/build variables).
- Existing token file at `~/.sql-proxy/token.json` (optional).
- Local OS/browser availability.

## Authentication Contract

1. On startup, system checks persisted token file:
   - If valid token exists: build authenticated HTTP client and continue without browser flow.
   - If missing/invalid token: begin interactive OAuth flow.
2. Interactive OAuth flow MUST:
   - Generate and store per-session `state`.
   - Start local callback server on `localhost` preferring port `8080`.
   - Fallback to another available localhost port when `8080` is occupied.
   - Open default browser (or print URL when auto-open fails).
   - Validate callback `state`; fail on missing/mismatch.
   - Exchange callback code for token and persist token.
   - Shut down callback server gracefully.

## Token Persistence Contract

- Token path: `~/.sql-proxy/token.json`
- Permissions: owner-only read/write (`0600` on Unix-like systems or platform-equivalent)
- Invalid token handling: rename or remove invalid token file before fresh auth

## Error Contract

User-facing errors MUST be clear and actionable for:

- callback port bind failures (including fallback failures)
- browser launch failures (must include manual URL guidance)
- callback without code or with OAuth error
- callback state mismatch
- token read/decode/save failures
- token exchange failures

## Output Contract

- Return an authenticated `*http.Client` ready for Cloud SQL dialer initialization.
- Authentication feature must not degrade existing IAM/private tunnel defaults in proxy startup.
