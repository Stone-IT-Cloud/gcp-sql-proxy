# Contract: Tunnel Success and PostgreSQL Guidance Output

## Scope

Defines required CLI output behavior when tunnel connectivity is established or re-established.

## Success Output Contract

1. A success message MUST be printed when tunnel-ready state is confirmed.
2. The message MUST include:
   - local host
   - local port
   - target instance context
   - PostgreSQL client command template
3. The message MUST be printed on every successful reconnect event.
4. The message MUST NOT be printed for failed startup attempts.

## Security Contract

1. No secret credential values (passwords/tokens) may be printed.
2. Required credential fields must be represented using placeholders.

## Failure Contract

1. Failure output MUST suppress success wording.
2. Failure output must remain actionable and unambiguous.

## Verification Contract

Feature is conformant when:

1. Successful startup prints readiness + guidance exactly when ready.
2. Successful reconnect prints the same full guidance block.
3. Failure paths do not print success messages.
4. Output includes no secret values and remains copy-friendly.
