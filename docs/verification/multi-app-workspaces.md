# Multi-app workspace qualification — 2026-10-29

The owner selected **greenfield installation** and will re-register apps. Legacy credential import, singleton history adoption and unchanged-sender upgrade qualification were removed from the authorized change. Old singleton files are preserved and rejected for automatic reuse. This report covers the working-tree implementation; nothing was committed, published or deployed to a user environment.

## Evidence ledger

| Check | Result | Scope and limits |
| --- | --- | --- |
| `task test` | Passed | Go suite, frontend tests and Svelte/TypeScript checks. Final equivalent Go/frontend runs also pass after the last changes; 49 frontend tests, zero Svelte errors/warnings. |
| `task test:viewer-access` | Passed | Existing invitation tests plus app collision/authorization/stream cases. |
| `task test:delivery` | Passed | Readiness/entrypoint tests and six Python startup/probe/private-local-configuration regressions. |
| Focused `go test -race ./internal/app -run '^TestRegistry\|^TestManaged\|^TestInvitationsManaged' -count=1` | Passed | Registry persistence, concurrent updates, managed routing and viewer paths. No race reported. |
| `openspec validate multi-app-workspaces --strict` | Passed | Greenfield delta contracts. |
| `openspec validate --specs --strict` | Passed | All 10 canonical capabilities after synchronization. |
| `git diff --check` | Passed | Tracked changes; new source files also checked during review. |
| Registry validation and faults | Passed | Checked entropy failure/duplicate output, 32 actual registrations plus rejected 33rd, Unicode case folding, invalid names/tokens/paths/symlinks, private modes, concurrent process lease, external edits, malformed/version/schema/identity rejection, write/file-sync/rename/directory-sync faults, restart reconciliation and uncertain commits failing closed. |
| Registration/removal recovery | Passed | Unready stores never activate; failed pre-commit stores remain unregistered; interrupted purge stays disabled and resumes on restart; empty registry remains empty; re-registration uses fresh authority. |
| Credential custody | Passed | Enumerated every SQLite table/cell for generated credentials. Registry-only authority and explicit owner-only secret retrieval; ordinary lists/session responses omit credentials. Known registered credentials in projected telemetry reject before commit. |
| Two-app HTTP isolation | Passed | Same standard URLs, colliding record/trace/span/service identities, browser/server roles, malicious app attributes/routing headers, independent stores, visit deletion/replay, references and invitation-ID collisions. Cross-app data/admin/stream requests and mismatched generation/cursors reject. |
| Standard OTLP wire matrix | Passed | Logs/traces/metrics × HTTP JSON/binary Protobuf × none/gzip × browser/server roles: 24 combinations, exact retained request bytes repeated before/after registry restart without producer env tokens, with B active. Supported stable log/span identities stay idempotent. Metrics remain receipt-only and are not claimed idempotent. |
| Real managed Go SDK exporter | Passed | `TestManagedStandardGoExporterNeedsOnlyEndpointAndHeader` sends ordinary named logs to two generated app credentials using the existing Go SDK/provider. No app attribute, correlation or custom identity required. This qualifies that log path only. Existing suite also exercises baseline Go logs/traces/metric exporters. |
| Maintained literal demonstration | Passed | 22 REST requests plus protected per-scope journey/timing/privacy assertions against a fresh disposable greenfield app. |
| Workspace and diagnostic fixtures | Passed | App-scoped workspace cohort: 20 common paths, one rare path, one unmapped event and stable replay; diagnostic requests preserve scrubbed messages/stacks and record identities. |
| Standard fixture verifier | Passed | 37-step standard/native trace, mappings, references, prospective definitions and journey assertions through actual durable protected reads, including inside disposable containers. These are wire fixtures, not .NET/Next.js exporter runs. |
| Real owner browser | Passed | Chrome Guest on macOS: keyboard registration in English/French, duplicate-name failure, generated masked connection settings, explicit reveal without logging secret values, clearing on dismissal/reopen, app switch from Traces to Journeys, explicit saved trace-reference link, reload, two independent tabs and empty-app guidance. |
| Narrow-screen browser | Passed | Chrome device emulation at iPhone SE **375 × 667**, both English and French registration/settings controls, readable labels/wrapping and modal scrolling. This is desktop emulation, not a physical phone. |
| Real viewer browser | Passed | Explicit one-use acceptance selects/names invited app, shows scope/absolute lifetime, opens Journeys, hides owner selector/settings, confines forged app links and clears historical workspace after revocation. |
| Browser expiry observation | Passed on retry | A second viewer was opened and its disposable server-side deadline advanced past expiry. After the native connection recovered, Chrome Guest showed “Viewer access ended”, “Access is closed” and confirmation that workspace content was cleared. A read-only fixture inspection confirmed one redeemed, expired, unrevoked app A grant. This verifies server-enforced expiry clearing in a historical view; it is an injected deadline, not a real-time 24-hour wait. |
| Container ARM64 | Passed | Local non-root read-only image, fresh registration, generated settings, two active app stores, actual OTLP/REST/SSE, unchanged credentials and exact native evidence after restart, damaged-registry fail-closed startup and protected coordinated archive restore. Disposable Docker volume; no NAS deployment. |
| Container AMD64 | Passed | Same disposable rehearsal under local Docker architecture emulation. Physical x86 hardware performance is not qualified. |
| Compose configuration | Passed | Both supplied files validated with a private disposable `.env`, no producer tokens. Actual NAS bind ownership/network deployment remains unrun. |
| Actual .NET / Next.js / browser OTel exporters | Not run | No new framework/exporter qualification in this change. Older reports remain separately scoped evidence. The browser checks above exercise Sentrame's UI. |
| Physical Home Server / Synology / Tunnel / external HTTPS | Not run | No production deployment, real tunnel or sender configuration was changed. |
| Physical power-loss, filesystem/ACL portability and 32-app load soak | Not run | Fault injection, bounded store tests and local Docker restart do not establish these guarantees. |
| Legacy unchanged-sender upgrade | Not applicable | Removed by the owner's explicit greenfield decision; no legacy credentials imported. |

The real UI checks preceded the final generation-header/known-credential projection guards. Frontend type/tests and real HTTP regressions cover those final guards; a complete final browser rerun is not claimed. The final focused expiry observation passed after native access recovered.

The wire matrix preserves generated request bytes and authentication values through an in-process registry/service restart. The HTTP test listener can use a new ephemeral port; container tests preserve the reachable port through restart. These are distinct from an actual external exporter deployment qualification.

## Failures found and resolved

- The first matrix assertion incorrectly expected metric receipt to deduplicate. It now verifies stable log/span replay separately from unchanged receipt-only metrics.
- Literal fixture tests initially expected old unscoped browser APIs. Their expected URLs now include the selected app.
- Strict delta validation detected a renamed existing scenario in the greenfield edit. Its canonical scenario name was preserved, and strict validation passes.
- The old demonstration fixture put an unkeyed privacy sentinel inside diagnostic free text, which the pre-existing scrubbed-diagnostic policy intentionally retains. Its diagnostic secrets now use keyed password values, preserving the sentinel assertion without changing the receiver's diagnostic policy. A fresh disposable demonstration passed all assertions. Fixture tools also explicitly restore the synthetic demo definition for a new demo generation after a standard-mapping exercise; retained histories stay pinned to their captured definitions.
- Native browser creation through the browser host was unavailable; real UI checks used native Chrome Guest successfully. Later, the forced-expiry observation lost native windows; the connection recovered on retry and the expired workspace clearing was observed. No verification blocker remains.

## Recovery rehearsal and boundaries

The container script starts without producer environment values, registers apps through the owner API, sends standard telemetry, restarts with the same persistent volume and compares settings and protected evidence. It stops the container, exports the entire data directory to a mode-0600 archive with numeric ownership, deliberately damages `apps.json`, verifies nonzero fail-closed startup, restores the coordinated archive and rechecks credentials/evidence.

This is a stopped archive round trip, not live SQLite backup or a physical crash test. Keep registry, marker, every store and WAL/SHM state together. Restoring older authority can restore revoked tokens or invitations: reconcile the latest rotation/deletion/revocation state before enabling access. Old binaries cannot read this registry; rollback requires a matching old snapshot/image and cannot preserve later registered apps automatically. See [the operating and recovery guide](../MULTI-APP.md).
