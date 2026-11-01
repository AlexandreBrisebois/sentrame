# App management delta


## MODIFIED Requirements

### Requirement: Automatically generated ingestion credentials
Registration SHALL generate one cryptographically random telemetry bearer key with at least 256 bits of entropy. Keys SHALL be unique across registrations and grant only ingestion permission. Success SHALL be reported only after the app store is ready and its credentials are durably persisted in the JSON registry. Generated app credentials SHALL NOT be stored in a database.

#### Scenario: Register without a shell
- **WHEN** the owner submits a valid app name in the UI
- **THEN** Sentrame creates one key, persists it and opens connection settings with a single Telemetry key entry

#### Scenario: Persistence failure
- **WHEN** registration cannot persist its registry or initialize its store
- **THEN** the UI reports failure, no new credentials authorize intake and existing apps keep their prior configuration

### Requirement: Owner connection settings
The owner SHALL reopen app connection settings to explicitly reveal and copy the persisted app key and standard OTLP endpoint/header settings. Lists SHALL mask secrets. The key SHALL be described as write-only app attribution for untrusted observations, never permission to read data, administer the app or prove persistence. Secrets SHALL NOT enter URLs, browser persistent storage, ordinary listings, telemetry or logs.

#### Scenario: Connect a standard exporter
- **WHEN** the owner opens app connection settings
- **THEN** the UI supplies standard OTLP/HTTP signal URLs and an Authorization bearer header, explains base versus signal endpoints and allows copying the existing credential without a shell
- **AND** it explains that the sender must be configured with these values and no Sentrame SDK or telemetry app attribute is required

#### Scenario: Return later for the same credential
- **WHEN** the owner returns after restart and explicitly reveals a credential
- **THEN** the originally persisted credential is available, without rotation or a database copy of the secret

### Requirement: Explicit credential rotation
The owner SHALL rotate the app key through an intent-protected UI action that warns of the sender configuration update. A successful durable rotation SHALL immediately reject the previous credential for subsequent intake and SHALL preserve the app's identity and history. Rename, selection and restart SHALL NOT implicitly rotate credentials.

#### Scenario: Rotate the app key
- **WHEN** the owner confirms app-key rotation
- **THEN** the new credential survives restart, the old one is rejected and all senders using the previous key must update
- **AND** copyable settings explain how to update the sender's exporter configuration without changing instrumentation

#### Scenario: Rotate only the server role
- **WHEN** an old caller submits a role-specific server rotation request
- **THEN** it receives an explicit unsupported-role error, with no change to the single app key or intake authority

#### Scenario: Rotation cannot commit
- **WHEN** the JSON update fails before replacing the committed registry
- **THEN** rotation reports failure and the previously committed credential remains authoritative

#### Scenario: Rotation persistence is uncertain
- **WHEN** replacement occurred but persistence cannot be confirmed
- **THEN** affected intake remains unavailable until committed state is reconciled, and the UI reports recovery is required without claiming rotation succeeded
