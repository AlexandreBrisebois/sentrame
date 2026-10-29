# Data lifecycle delta

## MODIFIED Requirements

### Requirement: Explicit greenfield cutover
The system SHALL initialize this multi-app release in a fresh data directory without importing legacy credentials or single-app data. Existing single-app files SHALL be preserved and rejected for automatic reuse. Unsupported schemas SHALL produce an explicit startup error without resetting data.

#### Scenario: Old test database selected
- **WHEN** initialization finds an existing single-app database
- **THEN** startup reports that a fresh directory and UI app registration are required, without modifying the old database

## ADDED Requirements

### Requirement: Greenfield persistent initialization
A new installation SHALL start with an empty private registry and no producer environment tokens. Owner registration SHALL generate all app credentials. Subsequent starts SHALL use JSON and SHALL NOT import environment values. Missing or inconsistent registry state SHALL fail closed.

#### Scenario: Fresh installation
- **WHEN** no registry, initialization marker, prior database or managed stores exist
- **THEN** owner login opens app registration with no synthetic app or producer tokens required

#### Scenario: Obsolete environment values
- **WHEN** a registered installation restarts with obsolete producer token environment values
- **THEN** only the JSON credentials authorize intake and no credentials are overwritten

### Requirement: Complete app state isolation
Every app SHALL own its evidence, discovery, mappings, definition versions, traces, journeys, diagnostics, references, health counters, revisions, summaries, deletion barriers and invitations. Equal local identifiers SHALL NOT join, deduplicate, resolve links or affect another app. Cleanup, replay and deletion SHALL remain app-scoped with existing retention and privacy semantics.

#### Scenario: Deliberately colliding identities
- **WHEN** A and B contain the same service, trace, span, event, session, view, definition, invitation and revision identifiers
- **THEN** reads, derivation, replay and lifecycle operations on A cannot disclose or change B's state
- **AND** a reference or trace link cannot resolve across the app boundary

#### Scenario: Delete or expire one app's evidence
- **WHEN** cleanup expires detail or the owner deletes a supported visit in A
- **THEN** A's own detail, aggregates and replay barriers reconcile under existing rules while B remains unchanged

### Requirement: Coordinated persistent state recovery
Backup/restore guidance SHALL treat the JSON registry and its referenced app stores as one protected state set. Restores SHALL preserve app-to-store identity, credentials, generations and deletion state, apply expiry before reads, and reject mismatched state. Registry failure SHALL NOT trigger reconstruction or regeneration from database contents.

#### Scenario: Restore mismatched files
- **WHEN** a restored registry identifies a generation or app that does not match its database
- **THEN** the affected app is unavailable with a sanitized recovery error and is not served under another registration

#### Scenario: Restore predating deletion or rotation
- **WHEN** an older backup predates app deletion, credential rotation or invitation revocation
- **THEN** restore guidance requires reconciliation against the latest deletion and credential state before enabling access
- **AND** an old backup alone is not claimed to preserve later revocations

### Requirement: Bounded multi-app operation
The installation SHALL enforce shared intake, query and streaming concurrency limits alongside per-app storage/query bounds. App registration SHALL have a documented finite limit, and cleanup SHALL schedule work fairly across apps. Adding apps SHALL NOT multiply receiver listeners, unbounded store handles or existing global resource budgets.

#### Scenario: One noisy app
- **WHEN** one app exceeds an intake or query bound while another is active
- **THEN** backpressure remains explicit and work stays within installation limits, with app-attributed health feedback and no cross-app data fallback
