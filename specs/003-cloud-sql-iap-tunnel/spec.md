# Feature Specification: Cloud SQL Proxy and IAP Tunneling

**Feature Branch**: `003-cloud-sql-iap-tunnel`  
**Created**: 2026-05-05  
**Status**: Draft  
**Input**: User description: "Spec 03: Cloud SQL Proxy and IAP Tunneling"

## Clarifications

### Session 2026-05-05

- Q: How should startup behave when pre-flight permission checks cannot confirm access due to API/network uncertainty? → A: Fail startup with a clear "permission check unavailable" error.
- Q: What should connection instructions show for the database user when printing startup guidance? → A: Always show the email from the active authenticated identity; fail startup if it cannot be determined.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Start Secure Local Tunnel (Priority: P1)

As a developer, I want the proxy to open a secure local tunnel to a Cloud SQL PostgreSQL instance over private connectivity so I can connect from my workstation without exposing public database access.

**Why this priority**: Tunnel establishment is the core business value; without it the feature provides no usable outcome.

**Independent Test**: Start the proxy with a valid instance and authenticated session, open a local database client to the advertised host and port, and confirm traffic is relayed through the secure tunnel.

**Acceptance Scenarios**:

1. **Given** valid runtime configuration and valid authentication, **When** the user starts the proxy, **Then** the app opens a local listener and initializes a secure tunnel to the target Cloud SQL instance.
2. **Given** a client connects to the local listener, **When** the connection is accepted, **Then** the app establishes a corresponding remote database connection and relays traffic bidirectionally until either side closes.

---

### User Story 2 - Enforce IAM and Private Connectivity Defaults (Priority: P1)

As a platform owner, I want IAM Database Authentication and private connectivity to be enforced by default so every connection follows organizational security policy.

**Why this priority**: Security controls are mandatory guardrails and must be active for every session.

**Independent Test**: Run the proxy with default settings and verify that startup refuses any mode that would bypass IAM authentication or private connectivity.

**Acceptance Scenarios**:

1. **Given** proxy startup with default settings, **When** the tunnel configuration is prepared, **Then** IAM database authentication is enabled for all remote database dials.
2. **Given** proxy startup with default settings, **When** the tunnel configuration is prepared, **Then** remote database dials are restricted to private connectivity paths.

---

### User Story 3 - Receive Actionable Permission Feedback (Priority: P2)

As an operator, I want pre-flight permission validation and clear failure guidance so I can quickly identify missing access instead of troubleshooting generic timeout errors.

**Why this priority**: Clear diagnostics reduce support burden and speed incident resolution, but depend on the core tunnel flow.

**Independent Test**: Start the proxy with credentials missing required roles and confirm the command fails fast with role-specific guidance instead of a generic network timeout.

**Acceptance Scenarios**:

1. **Given** missing instance read permission, **When** the user starts the proxy, **Then** startup fails before dial attempts with a clear message describing the missing permission.
2. **Given** missing Cloud SQL client or IAP tunnel access roles, **When** the user starts the proxy, **Then** startup returns a human-readable message that advises contacting DevOps for role assignment.

---

### Edge Cases

- Local port binding fails because the requested loopback port is already in use.
- The authenticated identity is valid for OAuth but lacks one or more Cloud SQL or IAP permissions.
- Pre-flight checks pass but the remote tunnel path is interrupted during active session traffic.
- Local client disconnects abruptly while upstream copy goroutines are still active.
- Upstream database side closes first while local client remains connected.
- Startup is interrupted by process cancellation while listeners or relay goroutines are being initialized.
- The active authenticated session does not expose a resolvable human-readable email for connection instructions.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST initialize a Cloud SQL connector session for each proxy runtime using the authenticated HTTP client and the authenticated credentials derived from the existing authentication flow to satisfy IAM database authentication.
- **FR-002**: System MUST enforce IAM database authentication as the default behavior for all Cloud SQL connection attempts.
- **FR-003**: System MUST enforce private connectivity as the default network path for all Cloud SQL connection attempts.
- **FR-004**: System MUST perform a pre-flight permission check before opening the proxy serving loop to confirm required Cloud SQL instance access is available.
- **FR-005**: If required access roles for Cloud SQL client connectivity or IAP tunneling are missing, system MUST fail fast with a specific, human-readable remediation message advising the user to contact DevOps.
- **FR-005a**: If pre-flight checks cannot determine permission status due to dependency errors or network/API uncertainty, system MUST fail startup before dial attempts and return a clear "permission check unavailable" guidance message.
- **FR-006**: System MUST start a local TCP listener on loopback address `127.0.0.1` and the configured active port.
- **FR-007**: For every accepted local connection, system MUST create a corresponding Cloud SQL connection through the secure tunnel.
- **FR-008**: System MUST relay traffic bidirectionally between local and remote connections until completion or failure, and release both sides promptly afterward.
- **FR-009**: System MUST isolate each client relay flow so one failed client session does not terminate other active sessions.
- **FR-010**: On successful startup, system MUST print connection guidance that includes host `127.0.0.1`, the active local port, the human-readable user email derived from the active authenticated identity, and password instruction `[LEAVE EMPTY]`.
- **FR-010a**: If the human-readable user email cannot be determined from the active authenticated identity, system MUST fail startup before opening the local listener and MUST return a clear error describing the missing identity detail.

### Key Entities *(include if feature involves data)*

- **Proxy Session**: Runtime state for one proxy process including listener address, active port, authenticated identity, and target instance.
- **Permission Check Result**: Structured outcome of startup access validation including pass/fail status and missing permission or role context.
- **Client Relay Session**: Pairing of one accepted local connection and one remote Cloud SQL connection with bidirectional traffic lifecycle state.
- **Connection Instructions**: User-facing connection metadata (host, port, username, password guidance) emitted after successful initialization.

## Non-Functional Requirements *(mandatory)*

- **NFR-001 (Security Baseline)**: Proxy startup MUST preserve secure defaults by requiring IAM-based identity and private connectivity for every session.
- **NFR-002 (Error UX)**: Permission and startup failures MUST be actionable, concise, and free of low-level stack traces.
- **NFR-003 (Resilience)**: Per-connection relay failures MUST be contained to the affected session and MUST NOT crash the proxy process.
- **NFR-004 (Concurrency & Cancellation)**: Listener and relay operations MUST support graceful shutdown and context-driven cancellation.
- **NFR-005 (Operational Clarity)**: Successful startup output MUST provide unambiguous connection instructions that can be used directly by database clients.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of validated startup runs with proper access open a usable local tunnel within 5 seconds of command launch.
- **SC-002**: 100% of startup runs with missing required access roles fail before connection timeout and include role-specific remediation guidance.
- **SC-003**: 100% of successful startups display complete connection instructions including host, active port, user email resolved from the active authenticated identity, and empty-password guidance.
- **SC-004**: 100% of parallel client session tests confirm one relay failure does not terminate other active sessions.

## Assumptions

- Users run the proxy in environments where the existing authentication feature can produce a valid authenticated identity.
- The target Cloud SQL instance is PostgreSQL and reachable through approved private connectivity controls in the selected project and network.
- DevOps owns IAM and IAP role assignment and is the escalation path for missing access.
- Connection instruction output is consumed directly by operators or local SQL client tooling without additional transformation.
