# Viewer access delta

## ADDED Requirements

### Requirement: Instance owner and app-scoped access
Owner authentication SHALL authorize instance app administration. Each app data request SHALL explicitly resolve its target app independently of mutable browser/session selection. App-wide temporary viewer grants SHALL authorize only their bound app identity and generation. App registration, lists, secrets, rotation and deletion SHALL require owner authorization; mutations SHALL retain origin/CSRF protection.

#### Scenario: Two owner tabs
- **WHEN** an owner opens A and B in separate tabs
- **THEN** each tab's requests continue to use its explicit app context and neither tab changes the other's target through session-wide selection

#### Scenario: Viewer forges a different app ID
- **WHEN** a viewer invited to A substitutes B in a route, query, reference, trace URL or stream request
- **THEN** the server denies the request without revealing B's existence or data
- **AND** omission of app context does not select another app or expand the grant

#### Scenario: Numeric invitation IDs collide
- **WHEN** A and B have invitations with the same local invitation ID
- **THEN** redemption, revocation and session validation resolve the app as well as invitation identity, and revoking one does not affect or authorize the other

#### Scenario: Viewer attempts app administration
- **WHEN** a temporary viewer directly requests app listing, registration, connection secrets, rotation, rename or deletion
- **THEN** server authorization denies the operation independently of UI visibility

### Requirement: App-aware invitation experience
Invitation creation and acceptance SHALL name the target app and existing read-only scope and lifetime. A redeemed viewer SHALL open only that app, without an owner app selector or registration controls. App disable, deletion, generation change, expiry and revocation SHALL deny new reads and close affected streams within five seconds without affecting other apps.

#### Scenario: Invitation opens its intended app
- **WHEN** a valid app-wide invitation for A is explicitly redeemed
- **THEN** the workspace selects A before data loads and explains read-only whole-app access without exposing ingestion credentials

#### Scenario: App deleted while viewer is paused
- **WHEN** A is deleted while its viewer has a historical or paused view open
- **THEN** the next session check clears A's displayed content and open streams close within five seconds
- **AND** B's grants and streams remain valid
