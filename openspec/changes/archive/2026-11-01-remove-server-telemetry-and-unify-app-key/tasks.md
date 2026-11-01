# Tasks

Future implementation checklist for the recommended review candidate. Every item is intentionally unchecked. Decisions A–C in `design.md` and a separate implementation instruction are prerequisites; planning completeness is not approval to modify credentials, data or code. If Alex chooses another migration or receipt boundary, revise deltas/design/tasks before apply.

## 1. Review and compatibility contracts

- [x] 1.1 Record Alex's explicit choices for standard signals/client Traces, retained versus fresh app key and legacy-history policy; verify proposal/deltas/design agree with the recorded choices before code edits.
- [ ] 1.2 Obtain the missing `server-api-journeys-and-dependency-timing` files if available and produce a requirement disposition list; verify server expansion is superseded and independent investigation/activation concerns remain separately scoped, without applying those proposals.
- [x] 1.3 Lock the one-key connection/rotation JSON contract, retired-role error behavior, registry/store/computation versions and supported managed-state transition; verify a reviewed contract table covers old requests and unsupported state.

Historical-only limitation for 1.2: the missing proposal is unavailable; no implementation correctness depends on it. Its scope is superseded, but exact file reconciliation is not claimed.

## 2. Registry and app connection

- [x] 2.1 Implement one-key fresh registration and durable role-free lookup/rotation in registry/service construction; adapt registry/apps tests and verify entropy uniqueness, restart, failed/uncertain publication, authority-lock races, removal and app isolation.
- [x] 2.2 Update owner connection/rotation APIs and their tests; verify only one key is returned, old role inputs reject explicitly, viewers/unauthenticated callers cannot retrieve keys and origin/CSRF checks still protect mutations.
- [x] 2.3 Update `Apps.svelte`, app types/helpers, EN/FR key wording and app tests; verify masked/reveal/copy/confirm/cancel/error states use one key and rotation explains all-sender updates. Update MULTI-APP connection guidance in the same change.

## 3. Existing-state transition and replay

- [x] 3.1 Implement offline bounded preflight/staging for supported managed two-key state; verify disposable fixtures cover unchanged originals, private coordinated backup requirements, app/generation mismatch, unsupported schemas and secret-free counts/reports. Document the Home Server transition and explicit authorization boundary.
- [x] 3.2 Implement versioned single-key publication and restart markers under the approved key policy; verify interruption before/after each durable boundary keeps intake stopped or selects one verified authority, retired-key rejection and pure-client sender continuity without reading actual deployment secrets.
- [x] 3.3 Preserve compatible old client identity via explicit aliases and maintain tombstones/revocation/deletion state; verify unchanged log/span/event replay, rotation, equal old role IDs, cross-app collisions, conflicting content and deleted generations cannot duplicate or resurrect contributions.
- [x] 3.4 Implement shared legacy eligibility before scan caps for all snapshot/discovery/timeline/trace/error/live/baseline reads; verify mixed A→B→C never becomes A→C, inseparable summaries are non-comparable, compatible pure client captures retain meaning and no result mislabels unsupported history as zero. Document historical limits.
- [x] 3.5 Implement reference/definition compatibility preview and prospective role-free activation; verify old references remain stored, unsupported filters fail explicitly, source-disjoint rules cannot broaden/overlap silently and persistence goals require owner review rather than automatic replacement.
- [x] 3.6 Verify coordinated recovery/rollback with interrupted-migration and stale-backup fixtures; demonstrate older-version rejection and latest retired/rotated/deleted authority stays denied, or recovery remains disabled. Add the concrete recovery procedure to NAS/developer guidance.

## 4. Intake and client journey semantics

- [x] 4.1 Remove active source-role assignment from intake/types/discovery while retaining qualified OTLP signals and privacy projection; adapt protocol/standard-exporter/identity tests and verify generic SERVER-kind spans, absent optional metadata and unsupported wire cells behave honestly without new protocol claims.
- [x] 4.2 Implement client visit/tab eligibility, role-free mapping/service/carrier selectors and client-only default examples; adapt navigation/native journey/mapping tests and verify missing context remains received/unassigned, goals stay optional, rule conflicts block and no server-only or combined path is derived.
- [x] 4.3 Preserve native/paired/estimated/observed-active timing and client diagnostics under new identity domains; adapt trace/diagnostic/presence tests and verify late parents, overlapping spans, error→recovery, censored timing and unsupported persistence claims retain honest units and limits.
- [x] 4.4 Update OTEL-CONTRACT, CONTRACT and developer/agent integration instructions alongside the intake/model changes; verify one ordinary client exporter path, scoped optional enrichment, public-key limitations and supported forwarder behavior agree with executable synthetic examples.

## 5. Investigation and documentation cleanup

- [x] 5.1 Remove role controls/labels from workspace mapping, discovery, evidence, errors, trace detail and copied context; update frontend types/state and focused tests. Verify client filters remain useful and obsolete query/reference state has a localized recovery path.
- [x] 5.2 Verify Tools/Outils, all seven direct journey views, optional goals, exact scope navigation, timeline/diagnostic return paths and historical/live populations remain usable; run relevant workspace/graph/investigation/reference/reconnect tests and document any compatibility limitations locally.
- [ ] 5.3 Perform EN/FR keyboard, focus return, 390px and native 200% zoom review; verify 44px targets, readable controls, stable refresh and approximately165px desktop/201px phone normal-fixture canvas start. Record observed versus deferred checks without inventing accessibility qualification.
  Remaining 5.3: a fresh local UI retry was blocked by the locked Mac (automatic unlock failed). Native 200%/EN/FR/focus checks already passed; exact phone geometry and physical touch remain unverified.

- [x] 5.4 Update `cmd/fixtures`, REST requests/JSON fixtures, local helper scripts and qualification assertions for one client key and observation-based goals; verify fresh disposable generation/replay never selects production state or prints keys and update the synthetic demo/event-plan documentation.
- [x] 5.5 Audit README, USER-GUIDE, DEVELOPER-GUIDE, MULTI-APP, NAS, docs index, active planning/agent docs and configuration examples for obsolete server telemetry instructions; verify a residual `rg` audit distinguishes historical evidence/legacy decoders/Sentrame backend language from active removed features. Preserve archived reports with explicit historical status.

## 6. Integrated qualification

- [x] 6.1 Run `task check` and focused synthetic migration/recovery tests; verify Go/frontend behavior, types, canonical spec validation and whitespace, and report exact failed/blocked/not-run checks.
- [x] 6.2 Qualify fresh and migrated disposable apps end to end through the retained real standard client exporter: ingest, commit, restart/replay, discover, client journey/goal/error/timing, history, live and second-app isolation. Verify old server key rejection and protected read/admin denial with an app key.
- [x] 6.3 Qualify the selected actual browser/forwarder path and disposable container deployment; separate Home Server/Tunnel checks and existing connector ownership. Verify receiver limits/backpressure and document measured bounds without claiming public capacity or route support from local wire success.

  Local browser logs/traces through the disposable same-origin forwarder and the ARM64 container cell passed. See `docs/verification/client-browser-otel.md` and the implementation handoff. Home Server/Tunnel, physical devices and public capacity remain explicitly unqualified and outside this local operation.
- [x] 6.4 Produce a review matrix mapping every delta scenario to tests or observed acceptance evidence; verify no application changes, unapproved state operations or earlier-proposal implementations are hidden within the removal scope.

## Workflow follow-up

- Review completed implementation and evidence before any separately authorized deployment or live data/credential transition.
- Sync canonical specs and archive only after actual implementation and repository acceptance; this planning session does neither.
