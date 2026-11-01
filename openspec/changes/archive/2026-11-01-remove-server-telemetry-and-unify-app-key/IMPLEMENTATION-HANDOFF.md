# Implementation handoff — 2026-11-01

## Authorization and selected contracts

Alex's explicit request to apply `remove-server-telemetry-and-unify-app-key`, following the review of the retained-key and safe-history recommendations, authorizes this local implementation. It does not authorize live migration, credential operations, deployment, commits, pushes, canonical synchronization or archival. No such operation was performed.

The selected design keeps qualified standard OTLP/client diagnostics and Sentrame's backend, uses the existing browser key as the one app key at an explicit future offline cutover, and preserves incompatible legacy evidence without presenting it as a client path.

| Contract | Local implementation |
| --- | --- |
| Connection | Owner-only `/api/apps/{id}/connection` returns `telemetryKey`, `endpoint`, `name`; no browser/server fields |
| Rotation | Owner/CSRF-protected `/rotate` with `{}` replaces the whole app key; any `role` field, including empty/null, rejects without rotation |
| Registry | Version 2, one 256-bit generated `telemetryKey`, initialization marker 2; old registry requires explicit offline conversion |
| Store/projection | SQLite schema 4; `client-structural-v3`; unsupported schemas fail closed, without resetting data |
| Derivation | `journey-client-v5`; journeys require explicit visit/tab; server-only and combined derivation removed |
| Retired filters | `source`, `role`, `producer` query keys reject explicitly; mapping/goal source selectors reject; legacy browser-only definitions get an analysis-only compatible projection |
| Compatibility | Snapshot includes a localized legacy-history notice and excluded-record count; affected historical references return explicit incompatibility rather than broaden or return fake zero |

## What changed

The registry/service/ingestion path now has one write-only app key. The UI has one masked connection card and role-free mappings. Client journey loops, optional goals, standard span structure/timing, errors and protected trace detail remain. Tools/Outils and the seven direct journey views are preserved.

`cmd/migrate-client` is an offline converter, never startup automation. Preflight validates supported registry/store identities, bounded counts and definitions. Explicit `--apply --backup-confirmed` stages coordinated private backups and new stores before publication. It retains the existing browser key, retires the server key, uses a fail-closed interrupted marker and an unsupported authority version during publication, verifies staged digests and supports deterministic resume. It preserves evidence, definitions, references, summaries and deletion state. This CLI was exercised only by synthetic tests; no real installation was converted.

Legacy mixed visit records are excluded by persisted record IDs before query/display caps. Hashed temporary scope barriers prevent expiry or late arrivals from reconstructing a mixed path; continuations across that boundary are explicitly partial. New client replay aliases only the old browser identity domain, independently of the rotating key. Old server IDs never join that alias. Incompatible definitions get a new prospective empty definition for owner review, not widened rules or replacement goals. Inseparable old summaries remain stored and non-comparable.

Active guides, REST examples, generators and smoke/qualification helpers use one key and client observations. The previous developer integration evidence is preserved as `docs/verification/historical-server-integration.md`. Earlier independent proposal files and archived UI work are unchanged; the absent server-expansion proposal remains superseded by direction, without claiming its unavailable files were reviewed.

## Verified evidence

- `GOCACHE=/tmp/sentrame-go-cache task check`: passed after fixture assertions were aligned with the client examples. Go packages, 152 frontend tests, Svelte type checks (zero errors/warnings), all ten canonical spec validations and whitespace checks passed. The final acceptance rerun after the cutover-continuation addition is recorded in `/tmp/sentrame-task-check-acceptance.log`.
- Added focused migration/history tests pass, including original-record preservation, no fabricated mixed transition, persisted exclusions after expiry, delayed old observations, explicit partial continuations, compatible pure-client references, replay aliasing, interrupted stages/resume, corrupt staging, writer leases, identity mismatch and newer authority markers.
- `TestMigratedAppStandardExporterReplayAcrossRegistryRestart` passes using a real standard OTel Go exporter and the retained key: new client reception/journey, retired-key rejection, unchanged outer endpoint and idempotent replay after registry/receiver restart. This is an isolated synthetic regression, not a real data cutover.
- Updated `demo`, `workspace`, `standard` and `diagnostics` commands passed against a new loopback instance with private generated files under `/tmp`, through actual OTLP, durable storage, derivation and protected reads. Their fixed REST requests, replay/count/goal/timing/trace/error assertions passed.
- Local `linux/arm64` image `sentrame:client-review` built successfully using the repository Dockerfile. `scripts/container-smoke.py` passed against that image: non-root read-only runtime, readiness/static UI, protected reads, OTLP, 37-step native-model verifier, SSE, exact persistent-volume trace evidence after restart, damaged-registry rejection and coordinated backup/restore. Only disposable uniquely named test containers/volumes were created/removed.
- Fixture Go tests and six Python helper tests passed; Python checks used `/opt/homebrew/bin/python3`. No remote/Home Server/Tunnel verification was attempted.
- Native Edge review on the disposable app: three explicit-tab client journeys; one masked app key; English/French labels; native toolbar-confirmed 200% zoom with scroll-reachable key actions; Escape restores Account focus; keyboard Tools/Outils opens with preserved groups and direct views. No real credential was revealed/copied/rotated.
- Strict validation of this change passes. [SCENARIO-EVIDENCE.md](SCENARIO-EVIDENCE.md) indexes all 103 delta scenarios to the relevant test/documentation/observation evidence families.

Initial Go runs could not bind ephemeral loopback sockets inside the sandbox. The approved elevated synthetic runs completed; those initial environmental failures are not reported as product failures. Default Mac Go-cache writes were also restricted, so test caches used `/tmp`. No automatic approval rejection occurred.

## Remaining qualification and source limitation

Follow-up 2026-11-01: actual local browser SDK/forwarder validation is now complete. See [the browser report](../../../docs/verification/client-browser-otel.md). Four logs/two native spans formed one completed client journey; standard binary replay stayed single; receiver process restart preserved all six records exactly; the same browser providers exported a seventh record after restart with no settings change. Origin/body limits and key absence in public assets passed. Task 6.3's selected local browser/container qualification is complete; deployment/device cells remain separate and unqualified.

- Full native 390 CSS-pixel/physical-touch verification and exact phone canvas-start measurement remain incomplete. The earlier emulation scaling/cropping prevented reliable assessment. A fresh UI clone was prepared from only the synthetic SDK data for another attempt, but native state inspection was blocked because the Mac locked and automatic unlock failed. Manual unlock is needed to continue software checks; physical touch still needs a real device. No responsive pass or uncertain real mutation is claimed.
- Deployed Home Server/Tunnel routes, other browser/encoding/compression cells, physical devices and public abuse/load capacity remain unqualified. The new passed cell is specifically native Edge logs/traces through a local same-origin forwarder, HTTP Protobuf without compression. It does not establish direct cross-origin/browser CORS support or deployed-route/public capacity.
- `openspec/changes/server-api-journeys-and-dependency-timing/` remains unavailable. Exact file-by-file reconciliation is blocked by missing source, not a denial or pending approval. This limits historical reconciliation only: the new direction, actual code/deltas and synthetic regression tests define and verify the implemented behavior without depending on that missing proposal. Its server-expansion direction is superseded; independent investigation/activation proposals remain untouched.
- Migration of real retained state and retirement of real server keys are intentionally not performed. A separate explicit offline operational authorization and current private coordinated backup are required.

All changes remain local on `main`; no staging, commit, push, deployment, canonical sync or archive. The remaining incomplete acceptance/source tasks remain unchecked in `tasks.md`.
