# Telemetry intake delta

## MODIFIED Requirements

### Requirement: Configured application attribution
The system SHALL bind each registered browser/server credential to exactly one stable app and producer role. Standard intake routes SHALL resolve the app from that credential, independently of service/resource claims. Optional metadata SHALL NOT be required for attribution. Service selectors SHALL NOT move data into another app, and ambiguous or unauthorized credentials SHALL reject intake.

#### Scenario: Several services for one app
- **WHEN** a configured browser service and API service map to WFS
- **THEN** both belong to the same registered app while retaining service identity

#### Scenario: Same service name in different apps
- **WHEN** two authorized routes receive the same service name for separate applications
- **THEN** their events, queries, diagnostics and counts remain separate

#### Scenario: Missing service metadata
- **WHEN** an authorized route binds a valid standard record to an app but its service or version is absent
- **THEN** the record is accepted with unknown metadata rather than invented values

#### Scenario: Same endpoint and different app credentials
- **WHEN** A and B send to the same standard OTLP signal URL with their respective credentials
- **THEN** each request reaches only its credential-bound app even when service, trace, span, session and event IDs overlap
- **AND** record attributes or claimed app IDs cannot override the credential's app or producer role

#### Scenario: Multiple products in a collector
- **WHEN** a Collector exports multiple products to Sentrame
- **THEN** requests are separated by the configured app credential before export; one authenticated request cannot distribute resource groups among apps by payload claims


## ADDED Requirements

### Requirement: Standard sender configuration and restart continuity
Registered apps SHALL use the existing standard OTLP endpoints, Authorization bearer header and supported wire behavior. Restart or registering another app SHALL NOT change existing sender settings. Credential rotation SHALL be a separate explicit owner action. Legacy credential import and single-app adoption are outside this greenfield change.

#### Scenario: Unchanged sender across restart
- **WHEN** a registered sender sends a qualified request before and after restart
- **THEN** the same endpoint, header and payload work for the same app and stable-identity replay stays idempotent

#### Scenario: Configure a newly registered app
- **WHEN** a second app uses an existing standard OTLP exporter
- **THEN** endpoint/protocol/header configuration with the generated credential is sufficient for supported reception, without a Sentrame SDK, app-ID tag, custom envelope or journey fields
- **AND** journey analysis still depends on its separately documented mapping/correlation prerequisites

#### Scenario: Rotation and intentional reassignment
- **WHEN** the owner explicitly rotates a credential or intentionally configures a producer for a different registered app
- **THEN** guidance identifies the required exporter credential change without requiring instrumentation changes or implying that previously stored data moves with the producer
