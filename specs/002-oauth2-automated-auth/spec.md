# Feature Specification: OAuth2 Automated Authentication

**Feature Branch**: `002-oauth2-automated-auth`  
**Created**: 2026-05-05  
**Status**: Draft  
**Input**: User description: "Spec 02: OAuth2 Automated Authentication"

## Clarifications

### Session 2026-05-05

- Q: How should the app behave if localhost callback port `8080` is already in use? → A: Auto-select an available localhost port for this auth session and continue.
- Q: Should OAuth callback `state` be mandatory and strictly validated? → A: Yes, generate and strictly validate `state`; fail if missing or mismatched.
- Q: How should invalid persisted token files be handled before re-authentication? → A: Rename or remove invalid token file before starting a new auth flow.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Authenticate Without External CLI (Priority: P1)

As a proxy user, I want the app to authenticate directly with Google Cloud through a desktop OAuth flow
so I can connect without installing or running a separate CLI tool.

**Why this priority**: Authentication is a hard prerequisite for any connection workflow.

**Independent Test**: Run the app on a machine with no existing token and no external cloud CLI, complete
browser consent once, and confirm authentication succeeds and a token is persisted.

**Acceptance Scenarios**:

1. **Given** no stored token, **When** the user starts the app, **Then** the app launches a browser-based
   consent flow and completes authentication.
2. **Given** OAuth consent is completed, **When** the callback is received, **Then** the app exchanges the
   authorization code and stores token credentials for future sessions.

---

### User Story 2 - Reuse Persisted Token Securely (Priority: P2)

As a returning user, I want the app to reuse a previously stored valid token so startup is fast and
does not require repeated browser interaction.

**Why this priority**: Re-auth prompts on every startup would significantly degrade usability.

**Independent Test**: Place a valid token at the configured token path and start the app; confirm the
app uses stored credentials and skips launching browser consent.

**Acceptance Scenarios**:

1. **Given** a valid stored token exists, **When** the app starts, **Then** it uses that token to create
   an authenticated HTTP client.
2. **Given** token persistence is enabled, **When** token data is written, **Then** file access permissions
   are restricted to owner read/write only.

---

### User Story 3 - Recover From Invalid or Expired Token (Priority: P3)

As a user with stale or invalid credentials, I want the app to automatically fall back to a fresh browser
flow so I can recover without manual token file troubleshooting.

**Why this priority**: Token expiry and revocation are common and should not block users permanently.

**Independent Test**: Provide an expired or invalid token and start the app; confirm browser consent is
triggered, a new token is saved, and local temporary callback server shuts down.

**Acceptance Scenarios**:

1. **Given** stored token data is invalid or expired, **When** the app starts, **Then** it initiates a new
   local OAuth callback flow automatically.
2. **Given** callback handling completes, **When** new credentials are acquired, **Then** the temporary
   local callback server shuts down cleanly.

---

### Edge Cases

- The local callback port is already in use when authentication is required.
- Callback flow MUST continue by selecting an available localhost port when default port `8080` is occupied.
- Browser auto-open command is unavailable or fails on the host OS.
- Callback request arrives without an authorization code or with an OAuth error response.
- Callback request arrives with missing or mismatched `state` value.
- Token file exists but has incorrect file permissions.
- Token file contains invalid credentials and should not be retried unchanged.
- User closes browser before consent completes and callback is never received.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST authenticate users through a desktop three-legged OAuth2 authorization flow.
- **FR-002**: System MUST request Cloud SQL Admin access scope for authentication.
- **FR-003**: System MUST use localhost callback endpoint `http://localhost:8080` for OAuth redirects.
- **FR-004**: System MUST request offline access so long-lived refresh credentials can be issued.
- **FR-005**: System MUST store token credentials at `~/.sql-proxy/token.json`.
- **FR-006**: System MUST enforce restricted token file permissions equivalent to owner-only read/write.
- **FR-007**: System MUST check for a persisted token at startup and reuse it when valid.
- **FR-008**: System MUST trigger a new OAuth flow when no valid persisted token is available.
- **FR-008a**: When persisted token data is invalid, system MUST rename or remove the invalid token file
  before initiating a new OAuth flow.
- **FR-009**: System MUST start a temporary local callback server on port `8080` when interactive auth is required.
- **FR-009a**: If port `8080` is unavailable, system MUST choose an available localhost port for the
  current authentication session and continue the flow.
- **FR-010**: System MUST attempt to open the system default browser for consent on Linux, Windows, and macOS.
- **FR-011**: System MUST capture callback authorization code, exchange it for credentials, and persist the token.
- **FR-011a**: System MUST generate an OAuth `state` value per interactive auth session and validate it
  strictly on callback.
- **FR-011b**: System MUST fail authentication if callback `state` is missing or does not match the
  generated session `state`.
- **FR-012**: System MUST present a success response in the browser after callback completion.
- **FR-013**: System MUST shut down the temporary callback server gracefully after authentication completes or fails.
- **FR-014**: System MUST provide an authenticated HTTP client output for downstream Cloud SQL dialer initialization.

### Key Entities *(include if feature involves data)*

- **OAuth Session**: In-progress authorization attempt including redirect state, callback code, and result status.
- **Token Record**: Persisted credential set used to build authenticated HTTP clients across app restarts.
- **Callback Server Session**: Temporary local listener lifecycle for OAuth redirect handling.
- **Authenticated Client Context**: HTTP client instance configured with valid OAuth credentials for dialer usage.

## Non-Functional Requirements *(mandatory)*

- **NFR-001 (Cross-Platform)**: Authentication flow MUST work on Linux, Windows, and macOS.
- **NFR-002 (Security)**: Token persistence MUST never expose credentials with permissions broader than owner-only read/write.
- **NFR-003 (Reliability)**: Temporary callback server MUST always terminate after flow completion or terminal failure.
- **NFR-004 (Error UX)**: Authentication failures MUST return clear, actionable user guidance.
- **NFR-005 (Recovery)**: Invalid or expired token states MUST recover through automated re-auth without manual config edits.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of first-run auth tests complete successfully without requiring any external cloud CLI.
- **SC-002**: 100% of valid persisted-token startup tests skip browser consent and produce an authenticated client.
- **SC-003**: 100% of invalid/expired-token tests initiate re-auth and persist refreshed credentials.
- **SC-004**: 100% of token file permission checks confirm owner-only read/write protection after token save.

## Assumptions

- OAuth client credentials and consent screen are pre-configured for desktop use in the target Google Cloud project.
- Users have access to a browser-capable environment when interactive authentication is required.
- Callback port `8080` is the default and may require fallback handling if occupied.
- Cloud SQL dialer initialization consumes a standard authenticated HTTP client abstraction.
