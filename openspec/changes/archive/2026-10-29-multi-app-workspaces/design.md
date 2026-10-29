# Design

## Context

See `proposal.md` for motivation and `specs/` for the behavioral contracts. Inspection found:

- `cmd/sentrame/main.go` opens one `SENTRAME_DB`, holds a database lease and requires both producer environment tokens through `app.New`.
- `Service` owns one `Store`, two tokens, instance sessions, concurrency channels and a loopback OTel receiver. Both the public proxy and receiver consumer validate the bearer.
- `Store` schema 3 already persists an app ID and generation, but its evidence, settings, definitions, views, invitations, tombstones, summaries and boundaries all assume one app. Its initial definition is synthetic.
- Invitations already bind app identity and generation. Local invitation IDs and session references are only unique within the current singleton.
- `App.svelte` has a fixed “Synthetic workspace” label and unscoped `/api/*` requests. Persistent deployment storage already lives under `/data`.

These are observed source behaviors, not runtime qualification of this proposal.

## Goals / Non-Goals

**Goals:** Reuse the proven single-app derivation/store semantics behind explicit app selection; make credentials a UI-managed JSON concern; support greenfield registration without producer environment credentials; ensure every request and lifecycle operation resolves its app before reaching data.

**Non-Goals:** Schema migration from unsupported historical releases, cross-app trace stitching, identity federation, automatic external secret deployment, new SDK/exporter support, changes to journey correlation rules or production capacity claims. A missing correlation ID still limits journey analysis. Existing owner passphrase and instance transport settings remain deployment configuration.

## Decisions

### 1. One private JSON registry and one telemetry database per app

Use `<data-directory>/apps.json`, where the default data directory is the parent of the configured existing `SENTRAME_DB`. Keep that setting only to locate the parent data directory. Existing singleton files are never adopted or changed; a fresh directory is required. Registry version 1 stores a revision and entries with opaque app ID, registration generation, display name, state (`active` or `deleting`), storage locator and browser/server credentials. Credentials are recoverable secrets because the owner must be able to reveal them again. They never enter SQLite, including settings, audit records or invitation tables. App identity/generation also remain in each database for consistency checks.

New database paths are derived from opaque IDs under `apps/<id>/telemetry.db`; display names and request input never become filesystem paths. Only generated managed locators are accepted. Validate locators against those rules, reject symlinks escaping the data root, and never serve this directory as web content. Preserve one database schema and existing SQL app boundaries instead of adding app predicates to every table/query. A shared database with app columns was considered; the current singleton SQL, retention and derivation paths make a missed predicate a greater risk than a bounded number of SQLite handles.

Initial limit: 32 registered entries including pending deletion; reject further registration visibly. Use one bounded store manager, at most 32 open stores, one shared receiver, one instance session authority and global budgets. Reuse the current global intake/stream limits of eight each; introduce an eight-request query budget and serialize cleanup with round-robin bounded app work. Preserve existing per-app evidence/query/byte limits and document the resulting total disk ceiling. This is a bounded default, not a measured Home Server capacity guarantee.

Fresh apps receive an empty valid initial definition, not synthetic mappings or a synthetic goal. Demo tools explicitly register/select a demonstration app and seed only its definitions.

### 2. Durable registry commits govern activation

Own a registry/process lease as well as the existing database leases. Serialize read-modify-write changes; validate the entire candidate including distinct credentials and ID/path/generation consistency. Create the data directory with mode 0700 and registry/temp files with mode 0600 on supported Unix hosts; document equivalent owner-only ACL requirements elsewhere. Generate credentials with a checked cryptographic random source, not the current helper that ignores entropy read errors.

Write a same-directory private temporary JSON file, flush it, atomically rename it and synchronize the parent directory before publishing the new immutable in-memory registry/token index. No plaintext registry journals in SQLite. Refuse external concurrent edits; offline manual repair follows documented validation/restart steps. An error before replacement leaves the old registry authoritative. If persistence becomes uncertain after replacement, stop affected intake/administration and reconcile the committed file before resuming; do not report successful rotation or silently trust a speculative in-memory value.

Registration initializes the private app database first, then commits its registry entry. It becomes visible/authorized only after that commit. A crash before registry commit can leave an orphan directory, which is never auto-registered or exposed. Report such an interrupted-registration recovery condition; do not infer credentials from orphan stores. An existing empty registry is valid after deleting all apps and never triggers legacy adoption again.

Persist a private, non-secret `.apps-initialized` marker before first registry creation and synchronize its directory. If this marker exists but the registry does not, require restoration rather than creating substitute credentials. This also covers interrupted initialization with no managed stores. Interrupted first initialization also fails closed with explicit recovery guidance. Backups include this marker; it carries no credential values.

App deletion commits `deleting` and removes active credential authority first, closes sessions/streams, then performs bounded purge of only that app store (including WAL/SHM), followed by removal of its registry entry and credentials. Restart resumes pending deletion. Renaming leaves identity, generation and credentials unchanged. Rotation changes only the chosen role's secret; old requests must revalidate before storage commit so disable/delete/rotation cannot slip through a stale token index. Already committed requests retain their normal acknowledgement semantics.

### 3. Preserve standard sender URLs and payloads

Continue accepting the qualified combinations on `/v1/logs`, `/v1/traces`, `/v1/metrics` with `Authorization: Bearer <token>`. A token-index lookup resolves immutable app ID, generation and role. Use that context for the selected store throughout parsing, projection and commit. The loopback consumer must independently validate the bearer against the same registry, not trust a caller-supplied app header. No producer-controlled resource attribute or internal routing header can widen authority.

Do not add an app-specific URL requirement, app-ID resource field, proxy adapter or SDK. This preserves ordinary .NET, Next.js and Collector exporter configuration without a proprietary routing adapter. Existing OTLP HTTP body/compression/atomicity/error/durability rules still apply. Multiple services may share their app's server role; distinct producer namespaces beyond the existing browser/server roles remain deferred. A Collector forwarding multiple products needs separate credential-bound exports and must not merge their requests under one credential. Existing browser forwarding keeps its supported untrusted role; this does not introduce or qualify direct cross-origin browser intake.

### 4. Explicit browser API context and server-enforced grants

Keep instance login/logout/session and invitation redemption at their global endpoints. Add owner app administration at `/api/apps` and app connection/rename/rotation/removal routes under `/api/apps/{appId}`. Move data/configuration/invitation administration routes under `/api/apps/{appId}/...`; a shared resolver authorizes app identity/generation before calling existing domain handlers. Missing app context never falls back to a mutable selected app. Unscoped old data routes return an actionable error rather than choosing an arbitrary app. Browser API compatibility is outside the sender compatibility promise.

Owner sessions stay instance-wide. Viewer session keys include app ID, generation and local invitation ID, preventing collisions between per-app databases. Global redemption finds a bearer digest across the bounded active app set without exposing app lists; actual one-use consumption remains transactional in the app store. Viewer allowlists run after app resolution and before data reads. Existing expiry, CSRF, cookie, bearer-fragment and five-second revocation contracts remain. Ingestion credentials are unavailable to viewers, ordinary app-list responses and session payloads.

Use app-qualified SSE IDs/cursors (app ID, generation, revision) and reject mismatches; cache keys and saved references carry app context. A link to another app remains unresolved rather than searching that app for a matching trace/span/reference ID.

### 5. Workspace selector and connection flow

Route workspace state with an app ID in the browser URL and show its display name in the header. Owners see “Add app” and app settings. New installs begin at registration; successful registration opens connection settings. Use the existing visual language, English/French translation system and accessible dialog/select patterns. Keep registration to one required name field. Settings reveal each masked credential only on explicit owner action, offer separate copy controls and show reachable endpoint guidance without placing secrets in URLs or local storage. Clear revealed values on leaving settings or signing out.

Every frontend load captures app ID and a selection epoch. Switching aborts reads, closes SSE and dialogs/editors, clears app-derived state and starts the new app's load. Late results are ignored unless both identity and epoch still match. Preserve only instance preferences such as language; reset app-specific filters/selection so another app is not accidentally hidden or edited. A viewer starts directly in the invited app and has no owner selector. Root bookmarks select the sole app automatically, or show the app chooser when several exist.

Change the current `explorationMode` initial value from `trace` to `journey` and place Journeys before Traces in navigation. Opening an app without an explicit exploration target, switching apps and ordinary invitation entry select Journeys. Do not carry a prior app's Traces mode across selection. An explicit app-scoped trace target or saved trace reference takes precedence, and the user can still select Traces directly. Apps without journey evidence stay in Journeys with useful connection/mapping/correlation guidance; technical trace availability does not silently change the selected mode.

### 6. Greenfield qualification

The owner explicitly chose greenfield installation and will re-register apps. Do not import legacy credentials or adopt a singleton database. Refuse an existing single-app database at initialization without modifying it, and direct the owner to a fresh private data directory. Producer environment variables have no authority. Exercise generated browser/server credentials before and after restart and while another app is active, using the qualified OTLP encoding/compression matrix. Verify stable replay identity and app isolation. Keep actual exporter, browser, Docker and Tunnel qualification distinct from wire tests.

## Risks / Trade-offs

- [Recoverable credentials in JSON] → Restrict owner/file access, mask UI, use no-store secret responses and protect backups. File storage meets the requested boundary but is not encryption at rest.
- [Registry and SQLite cannot share a transaction] → Registry-controlled activation/deletion states, identity checks, private orphan stores, crash-injection tests and explicit failure recovery.
- [Per-app stores multiply file and cleanup costs] → Registration/handle limits, shared concurrency limits and fair cleanup; measure before Home Server qualification.
- [Old binaries cannot operate the multi-app registry] → Roll back from a coordinated pre-upgrade snapshot; do not run an old binary against only one current app database and imply all apps are covered.
- [Old backups can restore revoked authority] → Back up the complete state set, reconcile the latest deletion/revocation state before enabling restored access, and rotate credentials if that state is unavailable.

## Initialization and recovery plan

1. Select a fresh persistent data directory; preserve old single-app data separately. This change does not migrate data or credentials. Register apps again through the owner UI.
2. Initialize a private marker and empty versioned registry without producer environment credentials. Register each app through the UI and configure its existing standard exporter with the generated bearer credential.
3. On restart validate JSON, app identity/generation and each managed locator. Never reconstruct credentials from SQLite or import environment values.
4. Back up the complete stopped data directory including registry, marker and all stores. Restore it as one coordinated state set. Reconcile rotations, deletions and invitation revocations before enabling restored access.
5. Old binaries cannot operate this registry. A rollback requires a matching old snapshot/image and does not preserve newly registered apps or later revocations.
