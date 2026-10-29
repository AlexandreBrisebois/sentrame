# App workspaces

This release uses a **fresh persistent directory**. The owner chose greenfield registration: it does not import single-app history or environment credentials. Preserve any old directory separately, start with a new one and register your apps again. Startup refuses a singleton database without modifying it.

Sign in as the owner, choose **Add app**, enter a unique name (1–80 characters) and register. Sentrame generates separate random browser/server credentials and opens **App settings**. Each app begins with empty valid mappings. Journeys is the first and default view, including after switching apps and ordinary invitation entry. Empty Journeys explains connection/mapping/correlation prerequisites. Choose Traces explicitly for native operations; `/?app=<id>&trace=<trace>` and `/?app=<id>&reference=<trace-reference>` preserve explicit trace targets.

Use **App settings** to rename, reveal/copy connection headers, rotate one role or remove an app after entering its name. Rotation rejects the old credential as soon as success returns; update that role's external senders afterward. Renaming preserves identity and credentials. Removal disables intake and viewer access before purging that app's files. A pending removal can be retried or resumed by restart. Re-registering a name creates new authority.

## Sender configuration

Configure the app's **existing standard OpenTelemetry exporter** with its reachable OTLP/HTTP endpoint, HTTP Protobuf protocol and `Authorization: Bearer <generated-role-credential>`. A base endpoint adds `/v1/logs`, `/v1/traces` or `/v1/metrics`; signal-specific settings take the full signal URL. There is no Sentrame SDK, required app-ID attribute or instrumentation rewrite. Missing journey correlation limits business analysis, not standard reception. Metrics remain receipt-only and do not promise idempotent replay.

All apps share the same `/v1/*` URLs. The bearer selects exactly one app and browser/server role. Resource fields, service names and routing headers cannot select another app. A Collector exporting several products must use separate credential-bound export requests/pipelines. Do not merge their resources under one bearer. Keep server credentials in the sender's server secret store; browser paths retain untrusted authority through the supported app-origin forwarder. Multi-app support does not qualify direct cross-origin browser intake. Copying settings does not configure external senders automatically.

Producer environment tokens have no authority in the Sentrame process. Only the owner passphrase, origin, listen address and persistent-directory location are deployment inputs. Example variable names in .NET/Next.js recipes describe the **sender's own configuration**, not the Sentrame container environment. Existing providers should be reused.

## Storage and private access

The data directory is the parent of `SENTRAME_DB` (default `.sentrame/workspaces/telemetry.db`; Docker `/data`). The configured filename is an initialization guard against accidental reuse of old singleton data. App databases live only at `apps/<opaque-id>/telemetry.db`. No app name becomes a path. Schema 3 store semantics are reused independently for each app.

`apps.json` is version 1 with `version`, monotonically increasing `revision` and an `apps` array. Each entry has `id`, `generation`, `name`, `state` (`active` or `deleting`), `storage`, `browserCredential` and `serverCredential`. IDs, generations and generated credentials use checked 32-byte cryptographic randomness encoded as 64 hex characters. Credentials are recoverable owner secrets and never copied into SQLite. File storage is private access control, not encryption at rest.

The exact initial validator-compatible registry is:

```json
{"version":1,"revision":1,"apps":[]}
```

Do not hand-author working credential entries. Registration generates them. Unix data directories are mode 0700; registry, marker, temporary files and leases are mode 0600. Use equivalent owner-only ACLs on other hosts; the current process lease uses Unix flock. Docker bind folders must be owned by UID/GID 10001 with mode 0700. Never serve the data directory as web content. Protect backups to the same standard. Symlinks at or within the data root, unsafe locators and mismatched database identities are rejected.

Registry changes hold a process lease and serialize private temp-file write → file sync → atomic rename → parent sync → runtime publication. Intake keeps authorization locked through its SQLite commit. Before-replacement errors leave old authority intact. An uncertain post-replacement commit disables registry authority until restart validates the committed file. Do not retry with guessed credentials. External edits while running are refused; stop the process before offline repair.

## Recovery and coordinated backups

Stop Sentrame before copying its **whole data directory**: `apps.json`, `.apps-initialized`, every app database and any SQLite WAL/SHM files. Keep the matching image/configuration and owner passphrase separately protected. A consistent stopped archive is preferable to copying live SQLite files independently. Restore the complete state set with owner-only modes and matching ownership, then start and verify identity, retention and credential continuity.

- Missing/damaged/unsupported registry or an initialization marker without the registry: restore a coordinated backup. Startup never reconstructs app authority from databases or environment values.
- An interrupted first initialization leaves `.apps-initialized`. For a known empty disposable installation only, stop the service and deliberately choose another fresh directory; otherwise restore.
- An unregistered orphan directory from a failed pre-commit registration never authorizes intake. Stop the service, inspect and quarantine the orphan outside `apps/`, then restart. Do not guess an entry or credentials from its database.
- A `deleting` entry remains inaccessible. Retry named removal or restart to resume purge before final entry/credential removal. Preserve this state in backups.
- A database identity/generation mismatch or unsupported schema fails startup. Restore matching registry/stores; never rewrite IDs to make mismatched data appear valid.
- A backup predating rotation, deletion or invitation revocation can restore old authority. Reconcile the latest state before enabling access; rotate credentials and revoke affected invitations if that state is unavailable.

Old binaries cannot operate this registry. Rollback requires a matching old snapshot and old image. An old snapshot cannot include newly registered apps, later evidence, rotations or revocations. Preserve those separately before rollback; there is no singleton downgrade migration.

## Authorization and limits

Owner sessions are instance-wide. Every browser data/configuration request names `/api/apps/<app-id>/…`; two tabs never change each other's app. Old unscoped data routes return a selection error. Viewer sessions bind app ID, generation and local invitation ID. Invitations name the selected app, read-only scope and lifetime; viewers cannot list/register apps or obtain credentials. Expiry, revocation and removal deny subsequent reads and close streams on the five-second heartbeat. Previously copied information cannot be recalled. SSE cursors include app ID, generation and revision and reject another app's cursor. Generated browser links carry `generation`; data requests carry it in `X-Sentrame-App-Generation` (SSE uses the query parameter). Mismatched generations reject access. Known registered credential values in projected telemetry reject the batch before SQLite commit.

There are at most **32** registry entries/open app stores, including pending removals. One receiver and shared limits of eight intake requests, eight streams and eight data queries apply across apps. Cleanup visits one app each second in round-robin order. Existing per-app record, query, SQLite page and 896 MiB pre-write byte guards remain; 32 times that byte guard is 28 GiB before bounded write/WAL/metadata overhead. This is a configured bound, not a measured aggregate storage ceiling or NAS capacity guarantee. See the [qualification report](verification/multi-app-workspaces.md) for exact checked and unqualified environments.
