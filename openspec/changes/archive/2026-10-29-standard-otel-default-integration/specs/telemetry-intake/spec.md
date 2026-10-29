# Telemetry Intake

## ADDED Requirements

### Requirement: Progressive evidence availability
The system SHALL distinguish receipt/discovery, trace operations, configured business meaning, explicitly correlated multi-request journeys and optional browser presence. Missing higher-level evidence SHALL disable only dependent behavior, never invalidate otherwise supported OTLP.

#### Scenario: Standard telemetry only
- **WHEN** valid supported OTLP arrives without `stage.name`, session, stage occurrence or presence fields
- **THEN** it remains accepted and discoverable with available operation/error evidence, while unavailable journey/presence capabilities explain their missing prerequisites

### Requirement: Minimized structural receipt
The system SHALL retain supported native span times, identity, parent, kind/status, ordered events and bounded link targets, plus standard log context. Projection SHALL preserve exact timestamps and producer namespaces without retaining raw OTLP, arbitrary attributes or link payloads.

#### Scenario: Standard span structure
- **WHEN** a span carries start/end times, parent ID, events and links
- **THEN** the protected retained projection exposes supported structure and exact times, distinguishes omitted fields and does not require synthetic boundary events

#### Scenario: Structural bound exceeded
- **WHEN** more than 32 links appear on one otherwise acceptable span
- **THEN** at most 32 link targets are retained with an omission count, without rejecting the record solely for optional link truncation or implying complete structure

### Requirement: Contextual safe discovery
Discovery SHALL distinguish service, signal/carrier, span/event name, trusted producer and retained field availability. It SHALL expose unknown/omitted values honestly. Mappings SHALL use only supported privacy-projected fields; metrics SHALL remain receipt-only.

#### Scenario: Same name across carriers
- **WHEN** a browser log and server span have the same name
- **THEN** discovery separates their carrier/source context and permits intentional owner selectors without presenting the browser record as server evidence

#### Scenario: Route and secret fields
- **WHEN** telemetry contains an approved route template, HTTP method/status, raw URL/query and an authorization value
- **THEN** only supported validated selectors survive projection, and secrets/raw URL content remain absent from storage, preview, errors, copies and live reads

### Requirement: Credential-derived authority
Producer authority SHALL derive from authenticated intake configuration and SHALL NOT come from telemetry attributes, stage labels, trace IDs or claimed parentage. Browser forwarding SHALL preserve browser authority.

#### Scenario: Forged persistence claim
- **WHEN** browser-authenticated telemetry claims server source, `stage.name=voting.persist` and a trusted span's parent ID
- **THEN** it remains browser evidence and cannot satisfy a server-only persistence goal or suppress an equal-ID server record
