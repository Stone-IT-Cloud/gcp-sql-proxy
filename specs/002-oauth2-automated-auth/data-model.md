# Data Model: OAuth2 Automated Authentication

## Entity: OAuthSession

Represents one interactive desktop OAuth authorization attempt.

### Fields

- `session_id` (string): unique identifier for auth attempt.
- `state` (string): cryptographically random value for callback validation.
- `auth_url` (string): generated consent URL used by browser launch.
- `redirect_url` (string): localhost callback URL used for this run.
- `status` (enum): `initialized`, `awaiting_callback`, `exchanged`, `failed`, `cancelled`.
- `error_reason` (string, optional): terminal failure context.

### Validation Rules

- `state` MUST be generated per session and validated on callback.
- `redirect_url` MUST use localhost and include selected callback port.

## Entity: TokenRecord

Represents persisted OAuth credential data.

### Fields

- `access_token` (string)
- `refresh_token` (string)
- `token_type` (string)
- `expiry` (timestamp)
- `path` (string): `~/.sql-proxy/token.json`
- `permission_mode` (string): expected owner-only read/write permissions.

### Validation Rules

- Stored token file MUST use restricted permissions.
- Invalid token data MUST trigger file rename/removal prior to fresh auth.

## Entity: CallbackServerSession

Represents temporary local HTTP server lifecycle for OAuth callbacks.

### Fields

- `requested_port` (integer): default `8080`.
- `active_port` (integer): actual selected port (may differ if fallback used).
- `status` (enum): `starting`, `listening`, `received_callback`, `shutting_down`, `closed`, `failed`.
- `received_code` (string, optional): auth code from callback.
- `received_state` (string, optional): state from callback.

### State Transitions

- `starting -> listening` when local server binds.
- `starting -> listening` with fallback port when `8080` unavailable.
- `listening -> received_callback` on valid callback request.
- `received_callback -> shutting_down -> closed` after code exchange completion.

## Entity: AuthenticatedClientContext

Represents output consumed by Cloud SQL dialer initialization.

### Fields

- `token_source` (enum): `persisted_valid`, `fresh_exchange`.
- `http_client_ready` (boolean)
- `dialer_integration_status` (enum): `pending`, `attached`, `failed`.

### Validation Rules

- Context MUST be created only after valid token acquisition.
- Output MUST provide an authenticated HTTP client object for dialer setup.
