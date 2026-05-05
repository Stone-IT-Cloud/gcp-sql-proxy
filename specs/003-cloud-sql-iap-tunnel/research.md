# Research: Cloud SQL Proxy and IAP Tunneling

## R-001 — Cloud SQL connector authentication with Feature 02 OAuth tokens

**Decision**: Initialize `cloudsqlconn.NewDialer` with **both** IAM database authentication and **explicit dual token sources** suited to user OAuth refresh tokens: `cloudsqlconn.WithIAMAuthNTokenSources(adminTokenSource, iamLoginTokenSource)` where:

- `adminTokenSource` mints access tokens with scopes required for SQL Admin operations (must include `https://www.googleapis.com/auth/sqlservice.admin`; Cloud SQL connector docs also expect `https://www.googleapis.com/auth/cloud-platform` for the admin-side source when IAM DB auth is enabled).
- `iamLoginTokenSource` mints tokens scoped only to `https://www.googleapis.com/auth/sqlservice.login` for the IAM login leg.

Additionally pass `cloudsqlconn.WithHTTPClient(httpClient)` **only when** `httpClient` is the OAuth-configured Admin API transport the project already uses — but note: **`WithIAMAuthNTokenSources` already registers `googleapi` token transport for Admin**; redundant HTTP client layering should be avoided. Prefer deriving both `oauth2.TokenSource` values from the persisted refresh token (`oauth2.Token` from Feature 02) via two `oauth2.Config` instances that differ **only by `Scopes`** (same Client ID / secret).

**Rationale**: `cloudsqlconn` rejects incompatible combinations when IAM auth is enabled unless IAM login tokens are wired correctly (`WithIAMAuthNTokenSources` / `WithIAMAuthNCredentials`). Passing only `WithHTTPClient` without setting credentials pushes the dialer toward Application Default Credentials, which violates the Feature 02 authenticated-user requirement.

**Alternatives considered**:

- **`WithHTTPClient` alone** — Rejected for primary auth path (does not reliably supply IAM login token provider user credentials).
- **Service account JSON on disk** — Rejected for this CLI product direction (violates authenticated user story and Feature 02 contract).

---

## R-002 — IAM auth + private IP together

**Decision**: Enable IAM DB auth globally on the dialer via `cloudsqlconn.WithIAMAuthN()` and force connector dials over private addressing via `cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP())`.

**Rationale**: Matches spec FR-002/FR-003 and organizational Zero Trust stance in the constitution. Public connectors like `cloudsqlconn.WithPublicIP` / `WithAutoIP` are not acceptable defaults.

**Alternatives considered**:

- Auto IP fallback — Rejected (weakens connectivity guarantees and violates private-IP mandate).

---

## R-003 — Pre-flight permission checks via SQL Admin API

**Decision**: Implement `VerifyAccess` using `sqladmin.NewService(ctx, option.WithHTTPClient(authedHTTP))` followed by `Instances.Get(projectID, instanceName)` where `(projectID, _, instanceName)` are parsed from `project:region:instance`.

Interpret outcomes as:

| Outcome | User-facing posture |
|---------|---------------------|
| HTTP 403 or explicit permission denied payloads | Fail with DevOps escalation message referencing likely missing roles **including** `roles/cloudsql.client` and `roles/iap.tunnelResourceAccessor` (wording SHOULD NOT claim a single-role root cause unless API message supports it). |
| HTTP 404 / invalid instance identifiers | Fail with clear “instance not found or not accessible” guidance (distinct from IAM/IAP remediation). |
| Transport errors, timeouts, ambiguous 5xx | Fail with **permission check unavailable** (spec FR-005a) — do not silently continue into client dial timeouts. |

**Rationale**: `instances.get` is a lightweight entitlement probe aligned with operational expectations; aligns with UX requirement to surface access problems before ambiguous network failures.

**Alternatives considered**:

- **Skip pre-flight entirely** — Rejected (conflicts FR-004/FR-005).
- **Separate IAP API probe first** — Deferred to implementation if distinguishable diagnostics are feasible without expanding scope prematurely (Cloud SQL connectivity may indirectly exercise IAP prerequisites).

---

## R-004 — Resolved email for STDOUT instructions

**Decision**: After obtaining the OAuth-backed HTTP client / token sources, retrieve a human-readable mailbox using Google’s OAuth2 **`userinfo.email`** verification flow (`https://www.googleapis.com/oauth2/v3/userinfo` with bearer token) guarded by timeouts; if unavailable, fallback to validating `openid` **`email` claim** when present on ID tokens Feature 02 may reuse.

**Rationale**: Spec FR-010/FR-010a requires deterministic operator guidance; SMTP-style email aligns with Postgres IAM authentication expectations Cloud SQL operators document.

**Alternatives considered**:

- Print opaque numeric subject — Rejected as primary (not acceptable SQL client UX).
- Require manual `--user` flag — Rejected unless future spec explicitly widens CLI contract.

---

## R-005 — OAuth scope expansion for IAM-enabled connector

**Decision**: Extend the desktop OAuth scope set (Feature 02) to include **`https://www.googleapis.com/auth/cloud-platform`** alongside the existing Cloud SQL admin scope so minted user tokens satisfy `WithIAMAuthNTokenSources` administrative requirements while still issuing a separate login-scoped source for `sqlservice.login`.

**Rationale**: Cloud SQL connector documentation for `WithIAMAuthNTokenSources` lists both `sqlservice.admin` and `cloud-platform` for the administrative token source; Feature 02 currently requests only `sqlservice.admin`, which is insufficient for that contract.

**Alternatives considered**:

- **Drop `cloud-platform` scope** — Rejected unless Google support / connector maintainers confirm narrower scopes; default remains expansion to meet published guidance.

---

## R-006 — IAP tunneling wording vs connector mechanics

**Decision**: Maintain product language “IAP + private IP posture” while implementing via supported Cloud SQL Go connector transports. Operational docs SHOULD instruct IAM grants (`iap.tunnelResourceAccessor`, `cloudsql.client`, instance-level permissions) consistent with organizational DevOps templates.

**Rationale**: “IAP tunneling” manifests as entitlement and network path prerequisites; exhaustive protocol tracing is out-of-scope for the specification but entitlement messaging MUST remain actionable.

**Alternatives considered**:

- Embed BeyondCorp TCP forwarding CLI — Rejected (unnecessary divergence from Cloud SQL connector path for this MVP).
