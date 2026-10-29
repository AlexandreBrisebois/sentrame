# app-management Specification

## Purpose

Let an owner register and operate several named product workspaces in one installation without generating ingestion secrets or editing Sentrame environment files for each app.

## Requirements

### Requirement: Owner-managed named applications
The owner SHALL register, list and rename apps through the authenticated UI. Registration SHALL assign an immutable opaque identity and generation independently of the display name. One app SHALL contain multiple services. Registration SHALL NOT require a command line, environment-file edits or owner-supplied credentials.

#### Scenario: Register a second product
- **WHEN** the owner registers “Whats for Supper” alongside another app
- **THEN** both are selectable workspaces and the new app starts without another app's telemetry, mappings or saved views
- **AND** services appear within the app receiving their telemetry rather than creating apps automatically

#### Scenario: Invalid or duplicate name
- **WHEN** a name is blank, longer than 80 characters or duplicates an existing non-deleted app name after trimming and case folding
- **THEN** registration or rename is rejected with an inline explanation and no partial app or credential change

#### Scenario: Rename an active app
- **WHEN** the owner changes its display name
- **THEN** identity, credentials, telemetry, definitions, references and valid invitations keep their existing associations

### Requirement: Automatically generated ingestion credentials
Registration SHALL generate distinct cryptographically random browser and server bearer credentials with at least 256 bits of entropy each. Credentials SHALL be unique across registrations and roles. Success SHALL be reported only after the app store is ready and its credentials are durably persisted in the JSON registry. Generated app credentials SHALL NOT be stored in a database.

#### Scenario: Register without a shell
- **WHEN** the owner submits a valid app name in the UI
- **THEN** Sentrame creates both credentials, persists them and opens connection settings with clearly labeled browser and server roles

#### Scenario: Persistence failure
- **WHEN** registration cannot persist its registry or initialize its store
- **THEN** the UI reports failure, no new credentials authorize intake and existing apps keep their prior configuration

### Requirement: Private durable JSON registry
App metadata and ingestion credentials SHALL persist in a versioned `apps.json` under the persistent data directory. Writes SHALL be serialized and atomic with owner-only file access. Restart SHALL preserve credentials. Missing, malformed, unsupported or inconsistent existing registry state SHALL fail closed without regenerating credentials or silently rebuilding apps from telemetry databases.

#### Scenario: Restart and restricted access
- **WHEN** the service restarts with a valid registry
- **THEN** the same apps and credentials work without environment tokens, and unauthenticated users and temporary viewers cannot read the registry or credentials through static files, APIs or diagnostics

#### Scenario: Interrupted write or damaged registry
- **WHEN** a registry write is interrupted or the existing registry cannot be validated
- **THEN** an intact prior registry remains usable or startup reports a sanitized recovery error, with no partial registry accepted and no substitute secrets generated

#### Scenario: Registry missing beside managed stores
- **WHEN** managed app stores exist but `apps.json` is missing
- **THEN** startup requests restoration of the registry and does not treat the installation as fresh or repeat legacy credential import

#### Scenario: Registry missing after initialization
- **WHEN** an initialized installation loses its registry before any app has been created
- **THEN** startup detects the initialization marker and requires recovery rather than generating substitutes

### Requirement: Owner connection settings
The owner SHALL reopen app connection settings to explicitly reveal and copy persisted credentials and standard OTLP endpoint/header settings. Lists SHALL mask secrets. Server credentials SHALL be identified as server-only; browser credentials SHALL retain untrusted authority. Secrets SHALL NOT enter URLs, browser persistent storage, ordinary listings, telemetry or logs.

#### Scenario: Connect a standard exporter
- **WHEN** the owner opens server connection settings
- **THEN** the UI supplies standard OTLP/HTTP signal URLs and an Authorization bearer header, explains base versus signal endpoints and allows copying the existing credential without a shell
- **AND** it explains that the sender must be configured with these values and no Sentrame SDK or telemetry app attribute is required

#### Scenario: Return later for the same credential
- **WHEN** the owner returns after restart and explicitly reveals a credential
- **THEN** the originally persisted credential is available, without rotation or a database copy of the secret

### Requirement: Explicit credential rotation
The owner SHALL rotate either credential through an intent-protected UI action that warns of the sender configuration update. A successful durable rotation SHALL immediately reject the previous credential for subsequent intake and SHALL preserve the app's identity, history and other role's credential. Rename, selection and restart SHALL NOT implicitly rotate credentials.

#### Scenario: Rotate only the server role
- **WHEN** the owner confirms server credential rotation
- **THEN** the new credential survives restart, the old one is rejected and browser ingestion is unaffected
- **AND** copyable settings explain how to update the sender's exporter configuration without changing instrumentation

#### Scenario: Rotation cannot commit
- **WHEN** the JSON update fails before replacing the committed registry
- **THEN** rotation reports failure and the previously committed credential remains authoritative

#### Scenario: Rotation persistence is uncertain
- **WHEN** replacement occurred but persistence cannot be confirmed
- **THEN** affected intake remains unavailable until committed state is reconciled, and the UI reports recovery is required without claiming rotation succeeded

### Requirement: App-aware workspace navigation
The workspace SHALL show the selected app name in place of the synthetic badge and allow owners to select or register an app. App identity SHALL be part of workspace URLs. Switching SHALL close old streams, clear old app data and editors, and reject late responses from previous selections. Empty installations SHALL show registration; empty apps SHALL show connection guidance.

#### Scenario: Switch during a slow query
- **WHEN** the owner switches from A to B while A has in-flight reads or live updates
- **THEN** no A response, graph, diagnostic, editor, reference or stream is rendered or submitted as B
- **AND** new requests explicitly address B and cross-app filters/selections are cleared

#### Scenario: Reload a selected workspace
- **WHEN** an authorized owner reloads an app URL or opens its saved reference
- **THEN** that app is selected before loading data, without silently using whichever app another tab selected

#### Scenario: Keyboard and language support
- **WHEN** the owner registers or selects an app using a keyboard or either supported language
- **THEN** labels, validation, focus handling and status feedback remain usable in English and French

### Requirement: Journeys as the default app view
App workspaces SHALL default to Journeys for owners and temporary viewers when no explicit exploration target is supplied. Journeys SHALL appear before Traces in exploration navigation. Opening or switching apps SHALL NOT inherit another app's Traces selection. Explicit trace navigation and saved trace references SHALL still open Traces. Missing journey evidence SHALL show guidance within Journeys rather than automatically switching to Traces.

#### Scenario: Open or switch apps
- **WHEN** an owner opens an app URL without an explicit exploration target or switches apps while viewing Traces
- **THEN** the selected app opens in Journeys with Journeys selected first in exploration navigation

#### Scenario: Viewer opens an invitation
- **WHEN** a temporary viewer opens its invited app without an explicit exploration target
- **THEN** Journeys is the selected view

#### Scenario: Explicit trace exploration
- **WHEN** an authorized user selects Traces or opens an app-scoped trace link or saved trace reference
- **THEN** Traces opens for that app and the default does not override the explicit target

#### Scenario: No journey evidence yet
- **WHEN** an app has only technical traces or no telemetry and its workspace opens without an explicit exploration target
- **THEN** Journeys remains selected with connection or mapping/correlation guidance and Traces stays available for explicit selection

### Requirement: Explicit app removal
The owner SHALL delete an explicitly identified app after a named confirmation. Deletion SHALL durably disable its credentials and revoke its viewer access before purging app-owned data. A later registration with the same name SHALL receive a new identity, generation and credentials. Failures SHALL leave the app inaccessible with an actionable pending-removal state.

#### Scenario: Delete one of two apps
- **WHEN** the owner confirms deletion of A
- **THEN** A stops intake and reads, its retained data and credentials are removed under the deletion policy, and B remains unchanged
- **AND** A's credentials and invitations never authorize a newly registered app with the same name

#### Scenario: Stop during removal
- **WHEN** a crash interrupts data purge after disabling A
- **THEN** restart keeps A disabled and resumes bounded purge without reviving its credentials
