# Phase 0 Research: PostgreSQL Tunnel Success Message

## Decision 1: Emit success guidance from authoritative ready event

- **Decision**: Print success + PostgreSQL guidance only when the tunnel-ready event is confirmed, not during preflight.
- **Rationale**: Prevents false-ready messaging and aligns instructions with real availability.
- **Alternatives considered**:
  - Print early after listener bind (rejected: can mislead when upstream connection fails).
  - Print only at process start (rejected: ignores reconnect behavior).

## Decision 2: Reprint full instruction block on each successful reconnect

- **Decision**: Re-emit full connection established message and PostgreSQL configuration guidance on every successful reconnect.
- **Rationale**: Matches clarified requirement and helps users recover after transient failures.
- **Alternatives considered**:
  - First-success only guidance (rejected: insufficient for reconnection scenarios).
  - Short reconnect-only message (rejected: incomplete user guidance).

## Decision 3: Use safe placeholders for credential-sensitive parameters

- **Decision**: Include placeholders (e.g., `<db_user>`, `<db_password>`) instead of runtime secret values.
- **Rationale**: Preserves usability without leaking credentials or tokens.
- **Alternatives considered**:
  - Print resolved credentials if available (rejected: security risk).
  - Omit credential fields entirely (rejected: weaker onboarding clarity).

## Decision 4: Provide command template + explicit host/port context

- **Decision**: Output host, port, and a ready-to-copy PostgreSQL client command template in each success block.
- **Rationale**: Reduces user setup errors and accelerates first successful connection.
- **Alternatives considered**:
  - Free-form text instructions only (rejected: less actionable).
  - Documentation link only (rejected: extra context switch for users).
