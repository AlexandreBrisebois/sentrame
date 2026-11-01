# Proposal

## Why

Sentrame should explain client event journeys without requiring an owner to manage two telemetry roles, two keys or server-workflow concepts. Removing server telemetry as a product feature makes app setup and investigation simpler while retaining ordinary OpenTelemetry integration.

## What Changes

- **BREAKING:** Remove server-only and combined client/server journey support, trusted persistence goals, server credential controls and role selectors throughout active APIs, analysis and UI.
- Generate one durable telemetry key per registered app. Recommend retaining the existing browser key at an explicitly approved cutover and rejecting the old server key; no credential changes occur in this proposal.
- Keep standard OTLP signals and qualified transports for client observations, contextual discovery, safe diagnostics and client trace detail. Generic reception is not a promise of server-workflow analysis or truthful client claims.
- Keep optional goals, healthy loops, timing with honest limits, historical/live exploration, app isolation, temporary viewers, Tools/Outils and compact graph navigation.
- Introduce an explicit non-destructive compatibility boundary for legacy data, mappings, references, summaries and replay identities; never flatten mixed history into new client paths.
- Replace server integration recipes and agent instructions with one client-focused onboarding path using existing SDKs/exporters. Supersede the earlier server API/dependency expansion direction; retain independent investigation usability work.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `app-management`: one generated app key, connection settings and whole-key rotation.
- `telemetry-intake`: app-key attribution without producer roles; untrusted observations and compatible OTLP reception/replay.
- `journey-modeling`: client-only correlated journeys, role-free mappings and observation-based goals/timing.
- `journey-exploration`: remove role-oriented controls and server paths, preserve useful investigation and compact Tools.
- `trace-exploration`: retained client diagnostic context without a server dependency product.
- `journey-diagnostics`: client exception linkage without server execution claims.
- `data-lifecycle`: explicit non-destructive key/data cutover and non-comparable legacy evidence.
- `app-integration-guidance`: single-key client OTel onboarding, examples and qualification.

## Impact

Future implementation spans the registry, ingestion, replay/storage, derivation, protected APIs, Svelte types/state/views, synthetic fixtures, tests, scripts and integration/deployment/agent documentation. Sentrame's Go backend, SQLite persistence, receiver, authentication, Compose service and Home Server deployment remain required. No new collector tier, custom SDK or enrollment system is proposed.

This is a review candidate with recommended migration choices, not approval to apply them. See `inventory.md` for observed source and prior-proposal reconciliation, `design.md` for alternatives and decision gates, and `REVIEW-HANDOFF.md` for the short review. No code, data, credentials, canonical specs or earlier change files are changed by this proposal.
