# Tasks

## 1. JSON app registry and private credentials

- [x] 1.1 Implement a versioned registry under the persistent data directory with app identity/generation/name/state, validated storage locators, checked random browser/server credentials and a 32-entry limit; verify validation, entropy failure, duplicate names/tokens, path escape and symlink cases with focused Go tests.
- [x] 1.2 Implement serialized atomic JSON persistence, owner-only permissions, process lease and immutable runtime publication; verify concurrent updates, write/rename/sync failures and restart behavior with filesystem fault tests, including uncertain commits failing closed.
- [x] 1.3 Add a bounded app-store manager preserving existing SQLite semantics and identity checks, with empty definitions for newly registered apps; verify isolation of two stores, inconsistent identity rejection and absence of generated credential values from all SQLite tables.
- [x] 1.4 Implement registration activation and disabled-before-purge deletion with restart recovery; verify pre-commit orphan stores never authorize intake, interrupted deletion stays disabled, re-registration gets fresh authority and deleting all apps preserves a valid empty registry.
- [x] 1.5 Document registry shape, storage locations, file/ACL requirements and recovery states in the deployment guide; verify examples match the registry validator and never contain working credentials.

## 2. Greenfield startup and persistent state

- [x] 2.1 Add greenfield startup regression coverage: initialize an empty registry without producer environment tokens, preserve any old single-app files by rejecting their reuse, and verify a fresh app begins with empty valid definitions.
- [x] 2.2 Update `cmd/sentrame` startup to use the private registry and bounded app stores; verify identity/generation consistency, malformed or unsupported registry/schema rejection and restart without producer environment tokens.
- [x] 2.3 Persist a non-secret initialization marker and ignore obsolete producer environment values; verify first-run UI readiness, missing registry fails closed and an empty registry never imports old data or credentials.
- [x] 2.4 Update deployment/local startup scripts, `.env.example`, Docker/Compose guidance and backup/restore instructions for greenfield registration and UI-managed credentials; verify configuration/script tests and a disposable coordinated backup/restore rehearsal.

## 3. Route intake and lifecycle work by app

- [x] 3.1 Replace singleton token/store attribution in the public intake and loopback consumer with registry-resolved immutable app/generation/role context; verify same `/v1/*` endpoints, no required app fields, forged resource/header isolation and browser/server authority tests.
- [x] 3.2 Synchronize credential rotation/disable with final authorization before commit and add per-role rotation persistence; verify old credentials reject after success, failed pre-commit rotation preserves old authority, other apps/roles are unaffected and in-flight requests cannot commit after disabling their authority.
- [x] 3.3 Apply existing derivation, trace/error/discovery/reference, expiry, summary, revision and replay operations to the selected store; verify paired-app fixtures with colliding local IDs, cross-app link rejection, isolated visit deletion and retained-history behavior across restart.
- [x] 3.4 Share receiver/listeners and global intake/stream/query budgets, bound store handles and schedule fair cleanup; verify over-limit responses and simultaneous app workloads remain bounded without multiplying existing global channels per app.
- [x] 3.5 Update `CONTRACT.md` and integration/deployment guides with credential-based routing, Collector separation and measured versus unqualified limits; verify all documented standard signal paths remain unchanged and no manual new-token generation is required.

## 4. App administration and viewer authorization

- [x] 4.1 Add owner app list/register/rename/connection/rotation/removal APIs with masked lists and explicit secret retrieval; verify method, owner, origin/CSRF and no-store behavior, direct viewer rejection and no secret disclosure through errors, logs or session payloads.
- [x] 4.2 Introduce a shared app resolver for app-scoped data/configuration APIs and reject ambiguous old unscoped data routes; verify owner two-tab behavior and forged app IDs across every protected snapshot, trace, diagnostic, reference, deletion and SSE route.
- [x] 4.3 Bind viewer sessions/redemption/revocation to app ID, generation and local invitation ID; verify colliding invitation IDs, atomic one-use redemption, restart, expiry, cross-app denial and stream closure within five seconds without affecting other apps.
- [x] 4.4 Scope live cursor IDs/cache keys/reference URLs to app identity/generation; verify mismatched cursors and reference IDs never resolve across apps and retained references reopen in their app.
- [x] 4.5 Update user access/sharing documentation and API fixtures; verify the app name/scope/lifetime appear in the invitation flow and documented temporary viewers cannot access app administration or credentials.

## 5. Registration, selection and connection UI

- [x] 5.1 Replace the synthetic badge with app selection and registration, name validation, rename and empty-state guidance; verify frontend tests plus actual keyboard/narrow-screen checks in both English and French, including valid empty initial definitions.
- [x] 5.2 Add masked/reveal/copy browser/server settings and explicit per-role rotation/named deletion confirmations; verify successful/failed registration and rotation, feedback, cleared revealed secrets on exit, no secret browser persistence and clear sender-configuration instructions.
- [x] 5.3 Make app ID explicit in browser routes and all data requests; abort/epoch-guard in-flight responses, close old streams/editors and reset app-derived state on selection; verify delayed A responses after switching to B, two tabs, reload/deep links, app deletion and sign-out clearing.
- [x] 5.4 Restrict viewer UI to its invited app and preserve historical/paused access-expiry behavior; verify browser scenarios for forged links, expiration/revocation and absence of owner controls.
- [x] 5.5 Update README/user/developer/agent guides and `.NET`/Next.js examples to use UI-generated settings and existing providers; verify no Sentrame SDK/app-ID requirement, server secrets stay out of browser recipes and restart asks for no sender edits.
- [x] 5.6 Make Journeys the initial app exploration mode and first navigation option, resetting to Journeys on app switching while preserving explicit trace targets; verify owner/viewer entry, reload without a target, switching from another app's Traces, trace links/saved references and journey empty states without automatic mode changes. Update user guidance to describe Journeys as the default.

## 6. Integration acceptance and evidence

- [x] 6.1 Adapt fixture generation and REST Client requests to explicit app registration/selection, with synthetic mappings seeded only in a demonstration app; verify existing demo journey/diagnostic/standard-OTLP assertions and a new two-app end-to-end isolation scenario.
- [x] 6.2 Execute generated browser/server URLs, headers and payloads, repeat them after restart without producer env tokens and while B is active; verify durable reception and stable-ID log/span replay and receipt-only metric behavior within the existing qualified encoding/compression matrix, plus independent B reception.
- [x] 6.3 Run repository `task test`, `task test:viewer-access`, `task test:delivery` and strict change validation, resolving failures introduced by this change; record exact results and keep actual exporter, browser, container, Home Server and Tunnel evidence separate.
- [x] 6.4 Exercise a disposable container with the persistent mount and verify greenfield registration, JSON/store persistence through restart, sender continuity and protected backup/restore; publish a qualification report with passed/failed/blocked/not-run cells and recovery limitations.

## Scope decision

The owner explicitly selected greenfield installation and will re-register apps. Legacy credential import, supported single-app adoption and unchanged-sender upgrade qualification are no longer required. Existing single-app files must remain untouched; use a fresh directory.

## Workflow follow-up

- Synchronize canonical specs after implementation and qualification as authorized. Do not publish or deploy.

## Verification status

Canonical specs are synchronized. See [the qualification ledger](../../../../docs/verification/multi-app-workspaces.md) for exact passed, resolved-failure, blocked, not-run and not-applicable evidence.

All 29 tasks are verified. On retry, the native browser connection recovered and the historical viewer displayed “Viewer access ended” and confirmed workspace clearing after the disposable grant deadline was advanced past expiry. The grant was redeemed and unrevoked, distinguishing expiry from revocation. Browser invitation entry, forged app links and historical-view revocation clearing also passed; automated expiry/session/stream checks passed. Actual external exporters, Home Server and Tunnel checks remain separately not run, as recorded in the ledger.
