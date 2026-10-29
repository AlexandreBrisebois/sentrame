# App integration guidance delta

## MODIFIED Requirements

### Requirement: Practical independent server recipes
The guides SHALL cover existing .NET and Next.js Node SDK/provider reuse, logs/traces, correct HTTP endpoints/protocol/auth, UI app registration, generated credential retrieval and secret storage, standard cross-service propagation, named outcomes, exceptions and controlled error messages. Browser instrumentation SHALL have a separate trust/lifecycle path.

#### Scenario: Existing provider and browser bundle
- **WHEN** an agent finds existing OTel initialization and a Next.js Client Component
- **THEN** it reuses the existing provider, keeps server credentials out of browser code, validates Node versus Edge constraints and introduces no mandatory proprietary wrapper

## ADDED Requirements

### Requirement: Configuration-only app onboarding
User/developer/deployment guides, producer contract and agent instructions SHALL explain UI registration, JSON credential persistence and standard exporter configuration. Registered senders SHALL need no changes during restart. Legacy adoption is outside this greenfield change. New-app setup SHALL require no Sentrame SDK, app-ID field or instrumentation rewrite. Secret handling SHALL distinguish Sentrame's registry from each sender's own configuration.

#### Scenario: Owner adds a product
- **WHEN** the owner follows onboarding for a second app
- **THEN** registration generates credentials and the UI supplies copyable OTLP endpoint/header settings without editing Sentrame's environment or running a secret-generation command
- **AND** guidance explains configuring the external sender's existing exporter or secret store with those values

#### Scenario: Greenfield deployment
- **WHEN** the owner follows the deployment guide
- **THEN** guidance uses a fresh persistent directory, UI registration and generated credentials without producer environment values or legacy import

#### Scenario: Browser and server settings
- **WHEN** the owner connects a browser frontend and server to the same app
- **THEN** guidance assigns their distinct generated credentials to their existing supported trust paths and never places server credentials in client code
- **AND** multi-app support is not presented as qualification of a new browser exporter, CORS route, transport or deployed Tunnel

### Requirement: Sender upgrade qualification
Qualification SHALL demonstrate actual unchanged qualified OTLP requests before and after restart without producer environment tokens, and while a second app is active. Evidence SHALL separate wire compatibility from real exporter and deployed-route qualification. Static checks or mocked UI tests SHALL NOT establish sender compatibility.

#### Scenario: Exercise unchanged and new senders
- **WHEN** the multi-app implementation is qualified locally
- **THEN** the report records unchanged endpoint/header/payload reception and stable-ID replay for the first registered app, plus isolated reception using a second generated app credential
- **AND** it records JSON persistence, absence of database credential copies and explicit passed/failed/blocked/not-run results for exporter, browser, container and Tunnel checks
