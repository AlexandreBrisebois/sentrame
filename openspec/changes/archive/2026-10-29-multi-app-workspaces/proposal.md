# Proposal

## Why

The owner needs to track several real products in one Sentrame installation. The current synthetic workspace, singleton store and environment-managed producer tokens make onboarding another app an operational task and leave app separation implicit.

## What Changes

- Register named apps in the owner UI, each containing one or more standard OpenTelemetry services. Generate distinct browser/server ingestion credentials automatically and persist them in a private JSON file on disk, never in SQLite.
- Replace the hard-coded “Synthetic workspace” badge with an app selector and provide registration, connection settings, rename, credential rotation and app deletion flows.
- Open app workspaces in Journeys by default, with Journeys first in exploration navigation and Traces available through explicit selection or trace links.
- Isolate telemetry, discovery, definitions, traces, journeys, errors, references, summaries, deletion state and invitations by stable app identity.
- Use greenfield registration with no legacy data or credential import, as explicitly selected by the owner. Preserve standard OTLP endpoints and wire payloads. New apps need ordinary exporter configuration only; no Sentrame SDK, app attribute or custom envelope.
- Keep owner login instance-wide; temporary viewers remain restricted to their invited app.
- Update integration/deployment documentation, fixtures and qualification evidence for registration and sender continuity across restart.

This proposal is planning only. It excludes cross-app graphs/correlation, organizations, per-app owner accounts, automatic apps from service names, new transports, and narrower visit-only viewer grants. Application exporter configuration still needs to receive the generated credentials; Sentrame does not edit external applications or their secret stores.

## Capabilities

### New Capabilities

- `app-management`: UI app registration and selection, JSON-persisted generated credentials, connection settings and app lifecycle.

### Modified Capabilities

- `telemetry-intake`: Credential-selected app routing and unchanged standard sender compatibility.
- `data-lifecycle`: Greenfield persistent initialization, app-isolated durable state, deletion and coordinated backup/restore.
- `viewer-access`: Instance owner authority and enforceable per-app temporary viewer access.
- `app-integration-guidance`: UI-generated connection settings and configuration-only sender onboarding/restart guidance.

## Impact

Touches `cmd/sentrame/main.go`, service/intake/router/session/store/lifecycle seams under `internal/app`, the Svelte workspace and translations, Docker/local startup configuration, fixtures/REST requests, and user/developer/deployment contracts. Proposed storage is a versioned `apps.json` registry with one existing-model SQLite store per app. The existing supported schema is reused; existing singleton files remain untouched and require a fresh directory. App-scoped browser API URLs will change; standard OTLP `/v1/*` URLs remain unchanged. No external service or new SDK is required.
