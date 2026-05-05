# Research: OAuth2 Automated Authentication

## Decision 1: OAuth2 Desktop Flow with Local Callback

- **Decision**: Use OAuth2 authorization code flow for desktop apps with a localhost callback listener,
  strict `state` generation/validation, and offline access request.
- **Rationale**: Matches user requirement to remove external CLI dependency and keeps auth user-driven,
  secure, and recoverable.
- **Alternatives considered**:
  - Device authorization flow: rejected because browser callback UX was explicitly required.
  - Service account auth: rejected because user-scoped interactive auth is required.

## Decision 2: Callback Port Strategy

- **Decision**: Default callback to `http://localhost:8080`, but auto-select another available localhost
  port when `8080` is occupied and continue auth for that session.
- **Rationale**: Preserves deterministic default while improving reliability on developer machines.
- **Alternatives considered**:
  - Hard fail on `8080` conflict: rejected due to avoidable user friction.
  - Infinite retry on `8080`: rejected due to uncertain completion behavior.

## Decision 3: Token Persistence and File Security

- **Decision**: Persist tokens at `~/.sql-proxy/token.json` with owner-only read/write permissions (`0600`
  on Unix-like systems, platform-equivalent restrictions elsewhere).
- **Rationale**: Aligns with constitutional security discipline and least-privilege local secrets handling.
- **Alternatives considered**:
  - In-memory token only: rejected because restart usability requires persistence.
  - System keychain-only storage: deferred for future enhancement due to scope and portability complexity.

## Decision 4: Invalid Token Recovery

- **Decision**: On invalid/expired persisted token, rename or remove the invalid token file before
  starting fresh interactive auth.
- **Rationale**: Prevents repeated failure loops and keeps recovery deterministic for users.
- **Alternatives considered**:
  - Keep invalid file untouched: rejected due to recurring startup errors and ambiguous state.
  - Sidecar invalid marker: rejected as unnecessary complexity for current scope.

## Decision 5: Browser Launch Behavior

- **Decision**: Attempt default browser launch using host-native commands for Linux, Windows, and macOS;
  if auto-open fails, print the auth URL for manual copy/paste and continue.
- **Rationale**: Keeps flow cross-platform and resilient to host restrictions.
- **Alternatives considered**:
  - Fail immediately on browser launch error: rejected because manual URL open is viable.
  - Shell-script wrapper dependencies: rejected to minimize external dependencies.
