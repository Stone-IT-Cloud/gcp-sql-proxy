# Data Model: PostgreSQL Tunnel Success Message

## Entity: TunnelReadyEvent

- **Description**: Runtime signal indicating the tunnel is ready to accept PostgreSQL client traffic.
- **Key Fields**:
  - `local_host`
  - `local_port`
  - `instance_identifier`
  - `event_type` (`initial_ready` or `reconnect_ready`)
- **Validation Rules**:
  - Must only occur after successful connectivity establishment.
  - Must suppress emission in failure states.

## Entity: SuccessMessagePayload

- **Description**: User-facing message block emitted on each successful ready event.
- **Key Fields**:
  - `status_line`
  - `connection_parameters`
  - `psql_command_template`
  - `security_placeholders`
- **Validation Rules**:
  - Must include host and port.
  - Must include target context identifier.
  - Must not include secret values.

## Entity: PostgreSQLClientTemplate

- **Description**: Copy-ready command template for PostgreSQL client connection through tunnel.
- **Key Fields**:
  - `host_flag`
  - `port_flag`
  - `db_name_placeholder`
  - `user_placeholder`
- **Validation Rules**:
  - Must be syntactically valid and copy-friendly.
  - Placeholders must be explicit for user-supplied values.

## Relationships

- `TunnelReadyEvent` triggers one `SuccessMessagePayload`.
- `SuccessMessagePayload` includes one `PostgreSQLClientTemplate`.
