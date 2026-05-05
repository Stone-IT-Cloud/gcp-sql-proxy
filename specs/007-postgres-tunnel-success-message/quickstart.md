# Quickstart: PostgreSQL Tunnel Success Message

## 1. Start the tunnel successfully

Run the CLI with valid configuration and instance settings.

Expected:
- Tunnel startup reaches ready state.
- CLI prints "connection established" style confirmation.

## 2. Validate instruction payload

Confirm success message includes:
- local host
- local port
- target instance context
- PostgreSQL command template with placeholders

## 3. Validate failure behavior

Run with intentionally invalid startup configuration.

Expected:
- No success message is printed.
- Failure guidance is shown instead.

## 4. Validate reconnect behavior

Trigger transient disconnect/reconnect scenario.

Expected:
- Each successful reconnect prints the full success + PostgreSQL guidance block again.

## 5. Validate security constraints

Inspect output for sensitive values.

Expected:
- No passwords or tokens in output.
- Placeholders present for user-supplied credential fields.
