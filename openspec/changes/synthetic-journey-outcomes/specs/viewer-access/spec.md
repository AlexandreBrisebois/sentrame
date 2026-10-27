# Viewer Access

## Purpose

Allow a small team to view journey evidence and invite another viewer while keeping access credentials and permissions separate from tracked application sessions. This is a proposed contract.

Release A includes owner passphrase access. Invitations are deferred from that release, with the narrow invitation design below already selected. Broader read-only app-workspace access remains retained roadmap B, not an implicit expansion of these grants.

## ADDED Requirements

### Requirement: Passphrase login with independent access sessions
The system SHALL protect viewer access through passphrase login and server-validated cookie sessions. Viewer credentials and identifiers SHALL NOT enter tracked telemetry or become application visit identifiers. Cookie, login-rate, timeout and revocation policies SHALL be documented before real use, and protected query/SSE APIs SHALL enforce them.

#### Scenario: Viewer logs in
- **WHEN** a passphrase is accepted
- **THEN** an independent viewer session is established and no application visit is created, resumed or joined by that login

#### Scenario: Cookie absent or invalid
- **WHEN** an unauthenticated request addresses a protected query, export, diagnostics or SSE endpoint directly
- **THEN** the server denies access without exposing data through a public shell or cached response

### Requirement: Scoped cookie-based viewer invitations
A deferred invitation SHALL grant read-only access to one live anonymous visit/session through cookie-based access. The owner SHALL explicitly choose inclusion of the app's seven-day baseline and exception details. Invitees SHALL NOT edit configuration. Grants SHALL expire and be revocable; exact lifetime/use policy SHALL be explicit. Broader app-workspace sharing SHALL remain a separate later mode.

#### Scenario: Another viewer watches the same visit
- **WHEN** an invitation restricted to an app and live visit is redeemed with explicit baseline/detail inclusions
- **THEN** the new viewer can observe that authorized view through its own access session without receiving tracked-app household identity or the inviter's credentials

#### Scenario: Optional baseline and exception details
- **WHEN** the owner has not granted a baseline or exception-detail inclusion
- **THEN** the invite cannot reveal that material through toggles, saved-view queries, diagnostic APIs, URLs or replay; enabling it requires an explicit grant

#### Scenario: Expired or revoked invite
- **WHEN** a viewer tries to redeem an expired/revoked invite or reuse a one-use invite
- **THEN** access is denied with a plain-language recovery path and the target visit's existence is not disclosed to unauthorized callers

#### Scenario: Concurrent redemption
- **WHEN** two clients redeem a one-use invite simultaneously
- **THEN** atomic consumption creates at most one permitted redemption and the other request is denied

### Requirement: Enforced least-scope permissions
The system SHALL distinguish viewing/detail disclosure, copy/export actions, definition editing, registration and invite administration. Permissions SHALL apply to every API, SSE replay and cache read. A live-only grant SHALL NOT widen through filters or cursors. Copy-button permission SHALL NOT imply that already displayed diagnostic content cannot be manually copied.

#### Scenario: Viewer attempts a configuration change
- **WHEN** a view-only session requests a mapping or retention change
- **THEN** the server refuses the change independently of hidden UI controls

#### Scenario: Revocation during streaming
- **WHEN** an active grant is revoked
- **THEN** affected streams are closed and subsequent snapshot, replay and copy requests are denied

### Requirement: Intent-protected mutations and inert rendering
Cookie-authenticated mutations SHALL enforce an explicit CSRF/origin policy and SHALL NOT use destructive GET operations. Session scope and rotation SHALL address the selected deployment, including sibling origins. Telemetry and diagnostic text SHALL render as inert untrusted content after privacy filtering.

#### Scenario: Forged owner request from a sibling origin
- **WHEN** another same-site application attempts a cookie-authenticated owner mutation without the required origin/CSRF proof
- **THEN** the mutation is denied even if SameSite permits the cookie to accompany it

#### Scenario: Markup and session rotation
- **WHEN** diagnostic text contains markup or an access session/passphrase is rotated
- **THEN** text cannot execute markup/scripts and the chosen rotation policy invalidates the affected old query and active-stream access

### Requirement: Credential-safe invitation handling
The system SHALL avoid exposing passphrases, invitation bearers or viewer cookies through telemetry, access logs, referrers, diagnostics or shared view state. Invitation redemption SHALL exchange a bounded bearer for the scoped session, remove the bearer from subsequent navigation and use an allowed local destination.

#### Scenario: Share a view after redemption
- **WHEN** a viewer copies the resulting view URL
- **THEN** it contains view state without an embedded credential or implicit new permission grant

#### Scenario: Unsafe redirect
- **WHEN** an invitation requests an unapproved external redirect
- **THEN** redemption does not navigate to it or disclose the bearer there
