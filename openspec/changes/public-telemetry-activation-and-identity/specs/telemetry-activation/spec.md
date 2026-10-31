# Telemetry activation — exploratory delta

## Purpose

Define candidate behavioral boundaries for opt-in telemetry activation in publicly distributed applications while separating public write permission, stable application identity and trusted enrollment. This delta is exploratory, not implemented or approved for deployment; receiver ownership and enrollment choices remain open in design.md.

## ADDED Requirements

### Requirement: Disclosed opt-in activation
Activation SHALL require explicit installer/operator consent before telemetry or activation requests are sent. The flow SHALL disclose the receiver operator, destination, minimized fields and collection purpose, allow refusal without impairing core app use, and allow later disablement without a Sentrame account.

#### Scenario: First start or consent refusal
- **WHEN** an installer starts the app or declines telemetry
- **THEN** no activation request or telemetry is sent, and core app functionality remains available

#### Scenario: Disable after activation
- **WHEN** the operator disables telemetry
- **THEN** collection/export stops and queued telemetry is discarded under the disclosed local queue policy, without claiming historical server data has been erased

### Requirement: Public configuration has write-only scope
If public activation is offered, it SHALL return only intentionally public ingestion configuration for an owner-published app binding. App name/version and public signing material SHALL NOT authorize privileged enrollment, registry reads, server credentials, data reads or administration.

#### Scenario: Name and version bootstrap
- **WHEN** an unauthenticated caller supplies a registered public app identifier, app name and version
- **THEN** the response contains at most the published write-only configuration; supplied names and versions do not create an app, select privileged roles or expose private settings

#### Scenario: Public value used on protected API
- **WHEN** a caller uses a public identifier, ingestion token or distributed signature to request telemetry reads or owner operations
- **THEN** access is denied independently of whether public intake is enabled

### Requirement: Stable identity independent of credentials and versions
Release aggregation SHALL bind to an authorized stable app identity and its active intake generation, with service/version retained as metadata. Rotation SHALL preserve the binding. Names and version claims SHALL NOT redirect records across apps or create credentials or workspaces per release.

#### Scenario: Multiple versions and credentials
- **WHEN** releases using different permitted credentials report versions 1 and 2 for the same active app generation
- **THEN** both appear under that app, with version comparisons available without a shared credential

#### Scenario: Forged names and version cardinality
- **WHEN** a caller claims another app's name or sends many invented versions
- **THEN** the authorized binding remains unchanged and bounded metadata policies prevent unbounded version dimensions, without describing the accepted version as verified provenance

### Requirement: Explicit producer assurance
Publicly authorized telemetry SHALL remain distinguishable from protected server evidence throughout ingestion, deduplication, journeys and summaries. HTTPS, a public distributed signing key, a self-generated installation key or a signed downloadable binary SHALL NOT alone establish an authentic installation or truthful server persistence.

#### Scenario: Forged server event with public authority
- **WHEN** public intake claims server source, persistence, another producer's event ID or a trusted parent span
- **THEN** it retains public assurance and cannot satisfy a trusted server-only goal or suppress protected producer evidence

#### Scenario: Installation key without trusted enrollment
- **WHEN** a client proves possession of its own newly generated key without additional admission evidence
- **THEN** assurance describes key possession or installation continuity only, and does not label the caller a genuine release or verified operator

### Requirement: Standard OTLP export after activation
Activation SHALL supply ordinary supported OTLP endpoint, protocol and header settings to an existing SDK/exporter or Collector. No Sentrame SDK fork, custom OTLP envelope or request-signature extension SHALL be required for basic export. Configuration refresh limitations SHALL be documented for qualified clients.

#### Scenario: Use an existing exporter
- **WHEN** consented activation returns valid ingestion settings
- **THEN** a qualified standard exporter can send unchanged OTLP using those settings; base and signal-specific endpoint behavior remains explicit

#### Scenario: Exporter cannot refresh in place
- **WHEN** the selected exporter cannot reload headers dynamically
- **THEN** integration guidance describes bounded reconfiguration or a standard Collector path without claiming environment variables are automatically hot reloaded

### Requirement: Bounded activation and public intake
Public activation and intake SHALL enforce finite request, decompression, concurrency, enrollment, rate, cardinality and storage budgets with explicit backpressure and an operator disable/revocation path. Untrusted requests SHALL NOT create unbounded registry entries or bypass protected intake budgets. CORS SHALL NOT be treated as authentication.

#### Scenario: Spam and enrollment flood
- **WHEN** a caller floods bootstrap, enrolls many keys or submits high-cardinality OTLP
- **THEN** work and state remain within declared budgets, excess requests receive bounded rejection/backpressure and protected app authority remains unchanged

#### Scenario: Non-browser caller forges Origin
- **WHEN** a caller omits or forges browser Origin headers
- **THEN** public scope and all admission limits still apply independently of CORS policy

### Requirement: Revocation-safe identity recovery
Owner-authorized recovery SHALL distinguish restoring identity/configuration from reviving revoked permission. Rebinding SHALL preserve app/generation, deletion barriers and latest revocation state or fail closed. Public identifiers SHALL NOT authorize rebind, overwrite existing ownership or resurrect a deleted app.

#### Scenario: Stale backup or cached token
- **WHEN** a restore predates rotation/deletion or a client presents revoked cached configuration
- **THEN** authority is reconciled against current lifecycle state before intake is enabled; absent reliable state, recovery stays disabled rather than restoring old authority

#### Scenario: Same name after deletion
- **WHEN** an owner creates a replacement app with a previously deleted name
- **THEN** it has a separate identity/generation and old public aliases, credentials and queued batches do not automatically target it

#### Scenario: Authorized same-identity rebind
- **WHEN** an owner requests recovery of a published alias to its existing active app
- **THEN** identity/store consistency and current permission state are verified, and rebinding does not migrate unrelated history or silently change receiver ownership

### Requirement: Bounded offline activation and refresh
Activation SHALL have bounded timeouts, retry/backoff and cache lifetime, preserve app/generation and destination binding, and leave core app use independent of receiver availability. Rejection/revocation SHALL NOT trigger endless refresh. Expired configuration SHALL stop export rather than silently broaden destination or authority.

#### Scenario: Offline startup
- **WHEN** a consented app starts without network access
- **THEN** core use continues, only policy-valid cached settings are eligible for export, and local buffering/retries remain bounded with visible loss/completeness limits

#### Scenario: Consent or binding changed
- **WHEN** consent is withdrawn or cached settings name a different receiver, app or generation
- **THEN** old queued telemetry is not silently replayed under the new binding

### Requirement: Enrollment trust and replay boundaries
If enrolled activation is selected, its contract SHALL identify the trusted issuer/admission proof, scope, expiry and revocation. Signed enrollment SHALL bind issuer, audience, app/generation, installation key and freshness with bounded replay state. Key possession SHALL NOT substitute for admission proof or per-record delivery deduplication.

#### Scenario: Replayed or misbound enrollment
- **WHEN** an enrollment proof is reused, expired, issued for another receiver/app or bound to another installation key
- **THEN** enrollment rejects without creating authority or disclosing credentials, under the selected one-use/freshness policy

#### Scenario: Trusted issuer admission
- **WHEN** an accepted protected issuer authorizes an installation
- **THEN** recorded assurance names the issuer and verified admission claims, without implying every telemetry event or claimed version is independently attested
