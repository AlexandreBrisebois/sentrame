# Standard OTel implementation qualification

Local implementation started 2026-10-29. All databases used by Go seam tests are fresh `t.TempDir()` SQLite files. Existing `.sentrame` data has not been opened or reset. Wire fixtures are not .NET, Next.js, browser, container or external-route qualification.

## Evidence ledger

- Baseline `go test ./... -count=1`: passed before implementation (app 29.896s).
- RED `go test ./internal/app -run '^TestStandardOTelStructuralProjection$' -count=1`: executed actual loopback receiver → SQLite → protected snapshot; nine records accepted; failed because projection policy, native end/parent/kind/status/flags, event index, bounded links/drop counts and HTTP selectors were absent. No setup error.
- `go test ./cmd/fixtures -run '^TestStandardFixture' -count=1`: passed; shared readable standard wire source retains native structure and exact duration, absent end and no custom journey attributes.

## Qualification scope

This is a chronological ledger: earlier not-run statements describe their checkpoint, and schema1/2 or computation v2 results are historical iterations. The final matrix below identifies current-model evidence per seam. No deployment, publication or production collection occurred. Home Server capacity requires owner-supplied workload and resource budgets.

- RED `go test ./internal/app -run '^TestStandardOTel(NativeReplay|RejectsOld|HTTPProjection)' -count=1`: changed native end incorrectly returned200; old schema attempted initialization and failed without usable fresh-directory guidance. These are behavior failures. HTTP privacy rejection assertions already passed because the old projection omitted all HTTP fields.
- GREEN `go test ./internal/app -run 'TestStandardOTel|TestDiagnosticSameModel|TestEvidenceRead' -count=1`: passed (2.402s), including native structural fields, bounded links, safe HTTP fields, atomic changed-end/parent/event/removal/link/status rejection, forbidden sentinels absent in stored rows/protected reads, unsupported schema preservation, exact projection after close/reopen and same-byte replay. This qualifies direct local HTTP binary uncompressed wire assertions only.
- Intermediate full-suite run failed the two old compatibility tests (legacy diagnostic upgrade and numeric stored JSON). Replaced them with same-model conflict and rejection assertions as required by the greenfield cutover; this intermediate failure is not hidden as a new-model regression.

- Full Go suite after greenfield test reconciliation: passed (`go test ./... -count=1`, app31.612s).
- Task2.1 RED: `go test ./internal/app -run '^TestTraceExplorer' -count=1` reached authenticated reads and failed on absent `/api/traces` (404 instead of200). `node --test web/src/trace.test.ts` with a minimal unimplemented row function failed expected child depth1 versus0. The earlier missing-module invocation was setup failure, not the behavioral red evidence.

- GREEN task2.2: `go test ./internal/app -run 'TestTrace|TestStandardOTel' -count=1` passed (3.865s), including actual two-request ingestion of2,001 spans with2,000-record detail/truncation, viewer GET allowlist, revocation, ambiguous cross-source parents and exact native timing exclusions. `node --test web/src/trace.test.ts` passed2 tests; frontend build passed with0 Svelte errors/warnings.
- Task2.3 literal qualification: first run failed400 because the generator used generic protojson on hexadecimal OTLP IDs. Corrected to maintained OTLP codecs; strengthened fixture test to assert16/8-byte IDs. This was a fixture bug, not receiver support.
- A restart comparison initially included global revision and failed when normal cleanup invalidated the revision. Corrected the check to compare all retained detail with revision removed; no product behavior was suppressed. Global revision is allowed to advance independently of a particular trace.
- GREEN: `/tmp/sentrame-standard-qualify.py` built fixture command and executed all10 literal `rest/standard-otel.http` requests against a fresh separate service at `http://127.0.0.1:18096`: unauthenticated detail401, login200, native span/log exports200, protected list/detail200, exact binary replay200, equivalent JSON span replay200, gzip binary span replay200, stable protected native detail. SIGTERM/actual process restart using the same fresh SQLite file preserved retained detail. No direct DB preload. This temporary orchestration script supplied private random tokens without printing them; the durable generator/verifier lives in `cmd/fixtures`.
- Native Chrome Guest review of the running isolated service: sign-in, trace selection, tree sibling rows, exact nanoseconds, missing/invalid timing, typed relationships and timeline control visibly verified. Browser access host was unavailable; native Chrome provided the fallback. Actual browser telemetry exporter and VS Code REST Client interaction remain not-run. No accessibility conformance claim.

Qualification host: macOS-26.7.2-arm64-arm-64bit-Mach-O; go version go1.27.3 darwin/arm64. Local direct HTTP origin only.

Prepared fixture trace: `2e73cc84ad9f38d5434c37e6b1b1165f`. Retained disposable data: `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-standard-otel-z3izlhjw`.

| Generated body | SHA-256 |
| --- | --- |
| standard-logs.bin | `a855d6943f021c71113d7aa49eee4efb19b84af5ea5be9f982562a910d76a0f9` |
| standard-logs.bin.gz | `8fae5a756cc5dbd966b16c386bfcf690e2b95786ec560ae8e4455da0821a9aee` |
| standard-logs.json | `d1dfc4420d4e06167668e782841b01e8cac2c0f2fc5ad5628341dd6edbc7b012` |
| standard-traces.bin | `33152c3d6b8d65fbc97ff3b4f0b46bda4af0aee32bfbb8aa80c19677bc438165` |
| standard-traces.bin.gz | `794a266c60e993547989507056faedef137b9f1355c876f96d6b2c3022dc6003` |
| standard-traces.json | `360d06814927a2a0d293ee362815e0f225f51fc5305217909ce55394b8714839` |

- Task3.1 RED: contextual preview matched2 `vote.committed` records instead of1 because carrier was ignored; definition activation rejected disjoint services. `go test ./internal/app -run '^TestStandardMapping' -count=1` failed these assertions after successful real ingestion/login. Frontend contextual-mapping test failed blank service instead of retained `voting-api`. Goal completion-boundary, ambiguity and forbidden condition-field checks were added before implementing the overlay.

- Task3 GREEN: `go test ./internal/app -run '^TestStandardMapping' -count=1` passed; route-template condition regression first failed400 then passed after using the same projection grammar in validation. Full `go test ./... -count=1` passed (app34.319s); all43 frontend tests and build passed with0 errors/warnings.
- Extended literal REST run passed14 requests/assertion steps including actual owner preview/activation, explicit definition preview with exactly one trusted `Vote saved` span-event goal, native operation durations and unchanged receipt-time default overlay. The success assertion inspects configured outcome evidence, not export acceptance alone. Actual service stop/restart preserved retained detail. Fresh synthetic directory/trace/hashes for this run follow.
- Native Chrome review passed contextual service/carrier/source defaults, non-persisting v2 preview, `Vote saved · v2 · Configured event goal`, two native Persistence operation rows and match-preview control. No actual exporter/browser instrumentation was used. First reload briefly displayed a blank document during browser painting; console showed only the expected unauthenticated session401; the page subsequently rendered and sign-in passed.

Latest trace `3ee096fdef46c63c126f6474b4e287e5`; disposable directory `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-standard-otel-fvoes2k3`.

| Body | SHA-256 |
| --- | --- |
| standard-logs.bin | `81f4630a40f504eae0560d84b4219b4414de5fd64e5ced7095061ed86894be40` |
| standard-logs.bin.gz | `e3553a4e887e0b4bb86c90f1450dadb7191cc5830d6cb6e2bfba631316525175` |
| standard-logs.json | `9760b583823034688ae006f573576afaa33ca5c0dde9c72c722c5e00b5b22b62` |
| standard-traces.bin | `37a4ab90618295bd6417297e6dcddd8b331f440f97a5b4dc1a8828305666b17d` |
| standard-traces.bin.gz | `2e8701e80c9f422dd0de472d7c92ba5326418220ae0ce01991af4d8ffd7df49f` |
| standard-traces.json | `7d58fceedbcf0e3bcc8e2b8eac10e255afcfd6e64ecfcb34ca8347fedf5517b3` |

- Task4.1 RED: `go test ./internal/app -run '^TestStageHints' -count=1` failed absent contextual `stageHints` (a second assertion also encountered the absent field; no setup failure). Explicit owner precedence/history test already passed. Frontend acceptance test failed `vote.committed` versus explicitly accepted literal `voting.persist`. Tests cover empty/whitespace/nonstring/unsafe hints, separate labels, non-inheritance and browser forgery.

- Task4.2 GREEN: `go test ./internal/app -run '^TestStageHints' -count=1` passed (1.532s); frontend trace tests4/4 and Svelte checks0 errors/warnings. Discovery retains at most32 distinct emitted labels per contextual selector, separately disclosing omitted suggestions. Only unassigned active selectors receive suggestion actions; accepting builds a literal owner draft, not a tag-dependent rule. Actual receiver/storage/read tests prove explicit owner `Vote saved` overrides `voting`, invalid/absent values give no hints, child events do not inherit a span label, distinct labels remain separate and later definitions/browser claims cannot rename prior captures.

- Task4.3 GREEN:25 literal REST steps/assertions passed through a fresh actual service. Enriched spans preserve **Vote saved** over emitted `voting`; generic overlapping rule returns400; a saved overview reference retains its original definition; accepting `voting.retry` activates a literal activity rule only for future captures. Future `changed.retry` emissions keep the owner's `voting.retry` label. Exact enriched replay preserves evidence identities, structure, duration and receipt-time overlay; active suggestion availability is intentionally separate configuration context. Actual process restart preserves retained detail. Trace-mode references are not qualified by this overview-reference check; that is a later increment.
- Additional native Chrome check: sent one controlled synthetic `vote.reviewed` correlated log under server authority, selected **Use emitted suggestion: voting.review**, verified literal Stage with empty condition fields, Preview1 match/1 missing journey scope with trace availability, and activated the prospective owner draft. The old receipt-time record did not silently gain the new label. This is an application UI check, not browser telemetry-exporter qualification.
- Stage hint bounds check passed through actual log intake/read:34 distinct labels yield32 suggestions plus2 omissions, preserving34 received records. Focused tests passed (app4.223s; hint bounds1.767s). Local docs fence/link review found a pre-existing stale README baseline link after archival; corrected it to the existing archive. No archive/spec-sync action occurred.
- Task4.4 recipes use existing .NET Activity/ILogger scope and existing Next.js span/named log APIs with no new provider/SDK/wrapper. Naming/precedence/preview/agent instructions updated. These remain source-checked examples; actual .NET/Next exporters still not-run.

Task4 prepared trace `e5c060f07c911ebef492f504028be49e`, enriched/future traces in the private manifest, disposable directory `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-standard-otel-csj4ubst`.

| Body | SHA-256 |
| --- | --- |
| standard-enriched.bin | `218e42da7f3cda28b90567221db6a56f2e8590e7f8ddbaab197562d65b3b9748` |
| standard-future.bin | `d5b9c1fe32437baf36f542042e302ce178048b75a3490f2c42ffca577dfa9cdd` |
| standard-logs.bin | `9fe7985a07e52af62ffbab64572c14f22ac97a0858e636f7c1a4462b5de5fe44` |
| standard-logs.bin.gz | `4f346d89b9aaaad1c0637b59bfbc3714333058df29301919e7edd67fd57b3b4f` |
| standard-logs.json | `bb206bc789e75d28ff7deb163e628c40ffad3a1909fa9d3194b5127db347c6d1` |
| standard-traces.bin | `6d1db72693fff656fa476499e595a4a034b9827ea7ae1021fc6530274ce68f6f` |
| standard-traces.bin.gz | `8484ffcddd366ac9e143216f6c52d5f4be1ad894df62e7d72644d6d4e87b4f58` |
| standard-traces.json | `2c7139d90b72f7bcb07355929569b792677014045524bb66f4ab23ec48ad0e4e` |
| standard-ui-suggestion.json | `3ba0d0e68a77953e104bd8b6b4a54fd2dd365edb0eee962b0be047d18d0c3aaf` |

- Task5.1 RED: actual scoped native intake/read tests failed native goals closing at start (splitting parallel work), late earlier parent reselecting a newer definition, missing end completing a goal, and trace references rejecting mode. Mixed native/page-pair/presence regression was added before integration; existing lifecycle/presence/deletion suites remain relevant but do not by themselves qualify native spans. New native deletion/replay and browser-no-tab/source separation seam assertions supplement them. All data is synthetic and isolated.

- Task5.2 GREEN: actual scoped receiver/storage/protected reads pass native goal-end closure, parallel starts, missing-end exclusion, source/tab separation, scope deletion/replay prevention, expiry/no resurrection and identity-free native summary inspection. A protected whole-UTC-day summary read also verifies retained native basis after detail expiry. Late earlier parent preserves definition2 after3 activates. Native/page pair/presence bases stay separate; competing elapsed policies for the same occurrence yield no completed samples until one is explicitly selected. Timing conflict grouping is bounded by retained occurrences, without a long-lived membership index.
- Full Go suite after native capture integration passed (`go test ./... -count=1`, app37.253s); focused native/reference/presence suite passed. Exact journey/step/timing boundaries now serialize as decimal strings. Schema is2 because detail-only capture-definition metadata was added; schema0/1 is unsupported, with no migration/reset.
- Task5.3 GREEN:37-step literal run passed on schema2 with a fresh database. Separate failed/retry native request traces join only via explicit random scope/tab; a trusted commitment survives publication error and recovery. Before the browser-screen goal the journey remains active despite commitment/HTTP200; after its separately emitted observation it completes/recovered under the original definition4 despite5 activating. Native samples200/300ms, paired page1s and synthetic observed-active800ms remain separate; goal replay does not create another journey. A trace-mode reference pins definition4 and native200ms operation timing. Actual process restart preserves native detail.
- Native Chrome review passed reopening **Native voting request** as a trace reference at pinned v4, switching to chronological Journeys, two Vote processing occurrences, one Vote saved, one completed journey and native inspector median250ms with2 samples/p90 unavailable. This is UI/wire evidence; synthetic presence observations are not actual browser visibility collection. No .NET/Next exporter or deployment cell is filled by these results.

Task5 trace `b11e20de6a60ad561b92800335d10bfa`; fresh schema2 directory `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-standard-otel-pumpylg_`.

| Body | SHA-256 |
| --- | --- |
| standard-enriched.bin | `bc6ff570fe307acab7673e12cc7988190dbd22e1ec01aeeb81fb20526c5e7ff7` |
| standard-future.bin | `24167689be0b989356f79abb31a00908758fae26a14c95c8d5388545cf9fc5e9` |
| standard-journey-browser.bin | `84b2ec3b05544f17ffe9911743a2b91ace79084502b0011a2d98dee5c85f689e` |
| standard-journey-goal.bin | `98be8966beb9362f51b936959e84f7019a765d4e430663339bbe36650717afbf` |
| standard-journey-server.bin | `e023536ecadb3d034e18e0d10b2c83b92d876eb64ae117d8b010b1b32e79ee32` |
| standard-journey-traces.bin | `fb0045cd463ce7b5d1c8e9b2e86ef0e2cfebe7257907f62399eaa9d911c4cf59` |
| standard-logs.bin | `ab3b4842b2f4aaea6a27eda3196171ad627e7d237a84e2a753e1a83d1af74266` |
| standard-logs.bin.gz | `90e74a63a1bf46685d11b47ca475b7f26f6c34018cc7d8731ed105fcf2481917` |
| standard-logs.json | `53ee1ad4c5f128f1b70424b2e922fe3d602346fc0b4bf083756cf2f40b70ade7` |
| standard-traces.bin | `e3ae614e9341ee129b1dacc37684baee4c0d9e005def6b47ca237264e776a5c3` |
| standard-traces.bin.gz | `d5257decdb63fdafd2c7777022abfd2c75277e6081457d458c3fc983e3dada2f` |
| standard-traces.json | `d629966e1c0de6f6041994f1bff1bc3552222018ba1d00c6dbcb51d5de49cb53` |

- Exporter preflight RED: actual wire test for `POST /votes/{voteId}` exposed safe route names becoming unnamed, raw URL operation names surviving, absent omission disclosure and mapping validation rejecting the template. Corrected receiver/mapping name policy to preserve bounded approved HTTP templates, omit raw URLs/queries/dynamic numeric paths and expose `nameState`, projection omission counts and SDK resource/scope drops separately. Focused projection/hint/mapping suite passed (11.210s). This tightened projection is `structural-v2`; schema3 explicitly rejects prior local schemas without reset. Earlier schema1/2 results above are historical iteration evidence, not final-model support.
- Task6.1 actual .NET exporter GREEN: disposable ASP.NET Core sample compiled with0 warnings/errors using SDK10.0.102, runtime.NET10.0.2 and exact OTel1.20.1 Hosting/OTLP/ASP.NET Core/HTTP packages plus packages.lock.json. Source inventory: one hosted tracing configuration and one logging provider, one reused ActivitySource; no duplicate SDK/wrapper or external application edits. `/vote/fail`503, `/vote/after`503 after actual synthetic file persistence, `/vote/ok`200. Exported37 records/8 spans through HTTP Protobuf **none** to a fresh schema3 receiver. Native route `POST /vote/{mode}`, three distinct persistence spans with exact times/parents, two named `vote.committed` logs, named persistence/publication errors with exception message/stack, span/log stage hints, exact log context and no private sentinel were asserted via protected reads. No journey was fabricated. Producer SIGTERM flushed normal batching; an actual receiver stop/restart preserved all native fields. This does not prove abrupt-kill or power-loss durability.

.NET disposable source: `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-dotnet-otel-sj6xowoi/sample`; host info: Version:           10.0.102; OS Version:  26.7; RID:         osx-arm64; Version:      11.0.0-rc.3.26522.148; Architecture: arm64.

| Source / lock | SHA-256 |
| --- | --- |
| Program.cs | `cbf028e3857660efa6095c05317b208a204db9e85f85fcd5302096679f091948` |
| Sample.csproj | `91b5636eed3be562fa3d73fa0b59eb5fbc75e35732efa7577414f16b1c8de1b7` |
| global.json | `e7ac16311eecea3797bfa553cfaf77c023c89f4ee5d6ad1df402f9b10f896e35` |
| packages.lock.json | `385206a65c8cd343ebebedc05aa522de227b6b49da2da946fd1a5677dc458726` |

.NET qualification: direct Mac loopback HTTP Protobuf none for logs/traces **passed**. JSON/gzip/metrics/Edge/browser/cross-service/Tunnel/Collector/container/framework deployment cells **not run** by this sample. Official API/package references remain links in the developer guide; source review alone is not the pass evidence.

- Task6.2 actual Next exporter GREEN: disposable production App Router build on Next16.4.3, React/ReactDOM19.2.0 and Node26.12.0; exact OTel API1.10.0, SDK/log/resource3.0.0 and Node/HTTP Protobuf exporters0.300.0 in package-lock.json. One Node-only instrumentation path reported initialization count1; no @vercel/otel/duplicate provider. Real `/api/vote/fail`503, `/api/vote/after`503 after actual synthetic file persistence, `/api/vote/ok`200 and unhandled `/api/error`500 exercised the real Next request-error hook. HTTP Protobuf **none** export retained55 records/49 spans, three native persistence operations with parent IDs/times, two named commitments with exact log context, controlled exception message/stack, optional hints and no private sentinels/tokens. Explicit local shutdown/flush awaited SDK shutdown before process termination. Actual receiver stop/restart preserved all native fields; ordinary unscoped telemetry created no journey. This qualifies that explicit shutdown path, not automatic serverless/SIGKILL flushing.

Next disposable source: `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-next-otel-3qu_v2ft/sample`.

| Source / lock | SHA-256 |
| --- | --- |
| app/api/vote/[mode]/route.js | `a8ba47180ff54003d4bcd4d2d8241f1850746a707ee68fb2839f8fc7ec820dad` |
| instrumentation.js | `cb1c4fd2d500f11bd2a7c32a5b7b79a3ffe792a30810448ded84a6b54032af64` |
| instrumentation.node.js | `b94b9ba8b214b8637fd32233bfe8207d2d4d7782caefb818d11968c149ca37d6` |
| package-lock.json | `fc4c73805e63a71b6c7f81dfbc3119cb8d73777232ce78cb5a9e9f33caf4d421` |
| package.json | `7fe686a48d70e91f73dd8d2f6a520a3d1ec89693e38310375dab2cece633d61f` |

Next cell: direct Mac Node App Router HTTP Protobuf none logs/traces **passed**. Gzip/JSON/metrics/Edge/browser/cross-service/.NET↔Next/Tunnel/Collector/container deployment cells **not run** by this sample. No cross-service call was used; sharing a trace string alone was not treated as propagation proof. TypeScript recipe snippets remain source-checked where the disposable sample used JavaScript.

- Native browser presence follow-up RED: `go test ./internal/app -run '^TestNativeBrowserIntervalUsesOnlyBrowserPresence' -count=1` failed at the protected-read assertion: a completed native page interval had no observed-active basis despite scoped presence received through actual OTLP ingestion/storage. GREEN: the focused native/pair/presence suite passed (0.681s) after adding a separate browser-only presence basis for complete native intervals contained within the capture. Authenticated server operations cannot borrow browser observations even when scope/tab/service strings match. The regression also asserts exact active duration `600000000` ns rather than the elapsed interval. Derived computation is now `journey-native-v3`, with `presence-cap60-browser-v2`; earlier v2 exporter results remain historical evidence for their stated unscoped native fields, not qualification of this presence computation. This wire regression does not qualify an actual browser SDK.

- Task 6.3 actual browser exporter **passed**: Chrome 156.0.8078.49 Guest on macOS ARM64, disposable Next 16.4.3 production sample, Node 26.12.0, OTel API 1.10.0 / SDK trace+logs+resources+core 3.0.0 / OTLP exporters 0.300.0. Native browser SDK HTTP Protobuf **none** exported through reviewed same-origin `http://localhost:18098/api/otel/{signal}` under a server-held browser token to loopback receiver schema 3 / structural-v2 / journey-native-v3. The server producer used its independent server credential. No provider started before Begin, one browser provider per signal and one Node initialization; actual static bundles and rendered HTML contained no server/browser token or owner passphrase. Forwarder wrong-Origin 403 and streaming >1 MiB 413 passed; actual Chrome CORS and CSP block messages were observed through native UI control. Browser control was initially **blocked** (`cgWindowNotFound`); user restored the window, and it then passed. Disposable setup failures (external node_modules symlink outside Turbopack root; mistaken separate api-logs import under API 1.10) were corrected; these are setup failures, not red behavior tests.

Actual UI: Begin → fail-before-persistence → save-on-retry (200 shown) → forged commitment → hidden tab/CORS check → return → Stop → CSP check → leave unfinished span. Protected reads asserted two distinct pressed spans, W3C browser press→Node request parents and exact server commit log→persistence span context, one failed persistence, one real synthetic persisted vote, browser forged commitment remaining browser/unmapped, actual active/hidden/idle observations and explicit completion. Page native elapsed 36017000000 ns and observed-active 7814000000 ns stayed separate; no server timing borrowed browser presence. An unended span did not appear after navigation even though its correlated creation log flushed; no invented live span. Controlled forceFlush succeeded; abrupt kill, mobile/background suspension, beacon reliability and lossless lifecycle delivery remain **not run**. App-origin forwarder adds Node request telemetry (293 total records, 28 browser records); volume must be reviewed per app.

Durability: `/tmp/sentrame-browser-verify.py` used actual protected reads, stopped/reopened the receiver against the same synthetic database, and compared 43 trace records plus typed relationships exactly. Trace `f08592822ac4e46a36c9f7a296600cca`; retained disposable sample `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-browser-final-4kb7ee1r/sample`. Agent checklist walked on this sample: source/lock/provider inventory, minimum approved emitter boundaries, privacy/authority review, receipt before interpretation, optional hints, scoped presence only for its question, failure/retry/forgery/flush and evidence. No proprietary wrapper or external-app edits. This qualifies browser→Next parent propagation only; .NET↔Next and other deployment routes remain not run.

| Browser sample source / lock | SHA-256 |
| --- | --- |
| package.json | `eee17f06407b06865fdc40641a8a0dadef82f98929a58b80b001dc54dbce8599` |
| package-lock.json | `e276472f19e787cee00d0519d4455d4630db0c23df8a39b63dffb51a68d6ad73` |
| app/web-telemetry.js | `86d1f9c1ae2af00813fcca6954819016b74899b8dda42cc3161ee5074aed5f70` |
| app/voting.js | `63ddd974429b219430efde8dce4514ef0f56bf44db0122448035a5f542042cbe` |
| app/api/otel/[signal]/route.js | `f522826d9e3c6bdac6b918324befcd7cc9a992331383e2b037e2ddf616f9c752` |
| next.config.mjs | `8b755ea71e40bd39afab77f58df8cecf250d76b0309f6862637a0d7225a004cb` |

- Final tree-order regression RED: `npm --prefix web test` failed the new skewed-start assertion (44 passed, 1 failed): an explicit child appeared above its parent because tree rows retained timestamp order. GREEN: 45 tests passed after grouping children under only resolved retained parents, with bounded cycle handling; original timestamps remain untouched. Frontend production build passed with zero Svelte errors/warnings. This is UI transformation coverage, separate from the native Chrome runtime observations.
- Task 6.4 documentation reconciliation **passed**: all concrete matrix files reviewed/updated, including documentation index, synthetic event plan, NAS/DockerHub and environment/Compose comments. Published rc.1 explicitly predates this model; no supporting release number/image is invented. `python3 /tmp/sentrame-doc-check.py`: 14 files, 103 local links and balanced fences. Both Compose files passed `docker compose --env-file .env.example [ -f compose.cloudflare.yaml ] config --no-env-resolution --quiet`, without resolving secret runtime files. Current identity is schema3 / structural-v2 / journey-native-v3, with capture-time interpretation and browser-only presence policy.

## Final acceptance checkpoint — 2026-10-29

Build is an uncommitted working tree based on `42d7c79cd72a83753abe7034ddc19b5d140029d2`, receiver binary SHA-256 `994e3cb5cfc558bd3fc909a87fcf26192f206e1ae3b19894e0070e582412c510`. Schema **3**, projection **structural-v2**, journey computation **journey-native-v3**, native timing **native-inclusive-v1**, presence **presence-cap60-browser-v2**. Source and sample locks/hashes above identify the exercised content; no commit, push, publication, deployment, canonical sync or archive was performed by this implementation. Original unrelated work, including the already archived/synced baseline, was preserved. Roadmap B remains deferred.

Final `task test` **passed**: all Go packages (app 38.493s), 45 frontend tests, Svelte check zero errors/warnings. `npm --prefix web run build` **passed**. `task spec:validate` **passed**: this change strict-valid, eight canonical specs valid, zero failed. `git diff --check`, documentation local-link/fence checks, four Python smoke-probe regressions and both Compose syntax/config checks **passed**. These are local software checks, not framework or route qualifications.

Final literal Mac run: `python3 /tmp/sentrame-standard-qualify.py` → maintained `go run ./cmd/fixtures standard`, **37/37 steps passed**, actual receiver/storage/derivation/protected trace and journey reads, then actual receiver stop/restart against the same fresh database and exact retained trace comparison (revision excluded). Trace `75bec1fdfc11a9d8ecb2ee854f85fd9a`; disposable directory `/var/folders/7x/450tvgfj1lg2gx9p89txzx0w0000gn/T/sentrame-standard-otel-p8c1dphc`. No direct DB preload, mocked API or precomputed graph. Binary-none logs, binary-none/JSON-none/binary-gzip traces were exercised in the literal story; generated bodies alone do not qualify unused cells. The HTTP protocol suite separately exercises receipt encodings/metrics and Go SDK none/gzip; it does not qualify .NET/Next/browser compression variants.

| Final Mac generated body | SHA-256 |
| --- | --- |
| standard-enriched.bin | `eb083f58ca8cf150a93c21729419edf13e0b0268dd669686933b86286027ea55` |
| standard-future.bin | `4aaf87df2dc53049a094794386a87dbccca80591259a4c5362a086bb8e782a9c` |
| standard-journey-browser.bin | `7a62b5bf8e61043b78cbeaa19983334c94867dc1f1e8208d0c5cb8e7c10dd295` |
| standard-journey-goal.bin | `9f36b21e116c2de0f6c0f815ee71f2695d572829032a4339817c37643b4ed266` |
| standard-journey-server.bin | `fbba0f0d2c87fd165a08de820b2cd67cc470fda9a730998f6c9f68b97a5659d9` |
| standard-journey-traces.bin | `e2a8ddad632147c3e04bf3e4c49b4a14b5a05428623f9c2277b5df91d083e183` |
| standard-logs.bin | `d8bbf317d4a1669bef7d64df0d1b1f3a3d3fb2442f02130b362f1b778e5adce9` |
| standard-logs.bin.gz | `bc70c7bbfd17d490e664374662a72c9cf7d0767834e3e565ee4dea871fed11d0` |
| standard-logs.json | `8a727669793ed7d9f10701474bd0a3586b6a5c02da2b542c88ef6a9ffd208404` |
| standard-traces.bin | `be98d6e73db2f26f4089879d88084e09a8150e3605aeaa29862696b53fb4d23c` |
| standard-traces.bin.gz | `6838175ee6b52929e8761663b3fe4f9402faa890e2af1a3ea1c49d52521c63fa` |
| standard-traces.json | `39fd6b96dfc34366165b00feedd6b39dafca482002061c782f4bf3dc218c6360` |

### Q01–Q14 result by actual seam

“Passed” below applies to the named local checks, not all producers/deployments. Go seam tests create fresh SQLite stores, use authenticated real HTTP ingestion and protected reads where applicable; pure derivation/UI helpers are supplementary and never used as exporter qualification.

| ID | Status | Evidence and limit |
| --- | --- | --- |
| Q01 | Passed | Structural projection and protected trace tests plus final literal standard request: names/service/carrier/native spans useful without stage/scope; no invented journey. |
| Q02 | Passed | `TestTraceExplorerStandardProtectedStructure`, literal links/events/logs assertions, actual Chrome tree/timeline; typed parent/link/context, repeated siblings, missing/ambiguous targets. Skewed tree grouping regression separately passed. |
| Q03 | Passed | `TestStandardOTelStructuralProjection`, native mapping/goal tests and literal exact decimal duration/timestamps; .NET/Next/browser actual native exports retain identity/parent/times. |
| Q04 | Passed | Actual-ingestion 2,001-span bound test; missing/invalid end and SDK/projection omission reads; browser unended operation absent after navigation. Sampling completeness remains unknown. |
| Q05 | Passed | Contextual mapping preview/activation and safe route condition tests; literal owner mapping with no hint; earlier native Chrome contextual editor/preview. Unsupported raw fields stay unavailable. |
| Q06 | Passed | `TestStageHints*` and literal optional hint/override/conflict/accepted suggestion/history checks; native Chrome explicit suggestion acceptance. |
| Q07 | Passed | Literal failure/retry/trusted commitment/post-commit publication error/recovery and separate screen goal; .NET/Next actual persistence and after-commit failure; browser forge stays browser. Generic 2xx never substituted for persistence. |
| Q08 | Passed | Native scope/deletion/separation and lifecycle tests; literal explicit browser/server scope and actual browser sample. Missing browser context remains accepted-unassigned; no identity/proximity joins. |
| Q09 | Passed | Real receiver atomic conflict and exact replay tests; literal binary/JSON/gzip replay; actual process and container bind-storage reopen comparisons. Anonymous records remain distinct. |
| Q10 | Passed | Real late-earlier-parent capture test, literal definition switch while active, preserved native occurrence/timing/history and trace reference; native Chrome pinned-reference review. |
| Q11 | Passed | Real scoped native browser presence regression, literal optional paired/presence story, actual Chrome active/hidden/idle with native elapsed and observed-active separate. Backend operation timing cannot borrow browser presence; missing observation is unknown. |
| Q12 | Passed | Projection/privacy and native scope separation tests; actual exporter sentinels absent from durable protected reads; actual browser bundle/HTML credential scan, forged browser commitment and distinct producer tokens. No blanket production privacy certification. |
| Q13 | Passed | Actual protected trace allowlist/revocation/CSRF/owner-mutation tests and invitation/SSE suite; literal unauthenticated trace 401; actual container protected API/SSE. No external HTTPS/tunnel route inferred. |
| Q14 | Passed | Unsupported-schema preservation, real expiry/deletion/native reference and summary-boundary tests, expired replay rejection. Fresh disposable stores only; no old-data conversion or reset. |

### Producer and environment qualification

| Producer / environment | Status | Exact exercised scope / remaining limits |
| --- | --- | --- |
| Mac direct wire | Passed | Final 37-step literal story plus full Go seam suites. Standard native traces: HTTP Protobuf none/gzip and JSON none; literal logs: Protobuf none. Other receipt encodings use the separate protocol tests, not native fixture inference. |
| Mac .NET actual | Passed | ASP.NET Core SDK 10.0.102 / runtime 10.0.2 / OTel 1.20.1; actual traces/logs, errors, post-commit failure, shutdown and restart; direct HTTP Protobuf none. No .NET 11, JSON/gzip, browser or cross-service claim. |
| Mac Next Node actual | Passed | Next 16.4.3 / Node 26.12.0 / API 1.10.0 / SDK 3.0.0 / exporters 0.300.0; production App Router actual operations/logs/error hook/explicit shutdown/restart, direct HTTP Protobuf none. Edge/serverless/TypeScript runtime not run. |
| Actual Chrome browser | Passed | Exact version/source hashes above; app-origin forwarding, collection choice, no bundle secrets, CSP/CORS, namespace, actual presence, W3C browser→Node parents, controlled flush and unended-span loss. Protobuf none only. Mobile/other browser engines/abrupt termination not run. |
| .NET↔Next | Not run | No actual cross-service calls in these samples; browser→Node propagation does not fill this cell. |
| Linux ARM64 local container | Passed | Native ARM64 execution on Docker Desktop Linux/aarch64; non-root read-only image, 37 literal steps, static app, protected traces/SSE and exact bind-storage trace after restart. Wire fixtures only; framework exporters inside container not run. |
| Linux AMD64 local container | Passed on fresh rerun; initial startup failed | Emulated on ARM64 Docker Desktop. First process exited with SQLite `disk I/O error (8714)` before readiness. Fresh full rerun passed 37 literal steps, protected reads/SSE and exact bind-storage restart. Startup fault cause remains unresolved; native AMD64 not run. |
| Home Server / NAS hardware | Not run | Need owner-selected hardware/ISA/filesystem/workload and CPU/memory/disk/query/cleanup budgets. No capacity claim. |
| Shared Tunnel HTTPS | Not run | Auth/content paths/gzip/bounds/retries/exporter receipt and SSE flushing/reconnect/revocation must be exercised independently. |
| Collector | Not run | Separate queue, authority and acknowledgement boundary; direct receipt does not qualify it. |
| gRPC origin/public route | Not run | Outside this HTTP change; not advertised. |
| Production / WFS | Not run | No external application modified or contacted; no production telemetry collected. |
| VS Code REST Client extension | Not run | Literal `.http` requests were executed by the maintained verifier; extension UI itself not exercised. |

No current check is left **blocked**: restored browser control resolved that prerequisite. AMD64 has an executed initial **failure** and an independently recorded fresh **passed** rerun, not a blocked/not-run cell. Historical failures and reruns remain in this ledger; the unresolved intermittent startup fault is a material limitation.

AMD64 failure investigation: a disposable statically compiled probe using the exact current modernc.org/sqlite dependency ran under the same AMD64 image/emulation and a new bind directory. Both inspect-first and WAL-first sequences (user_version/schema reads, WAL/FULL setup, create/insert/read) passed. This does not reproduce the receiver startup failure or establish its cause. A fresh full smoke rerun is recorded separately; the first failure is retained. No storage durability policy or dependency was weakened to pass qualification. Library context: [modernc SQLite changelog](https://gitlab.com/cznic/sqlite/-/blob/master/CHANGELOG.md) describes code 8714 as SQLITE_IOERR_IN_PAGE, which alone is not root-cause evidence.

### Local container run identities and exact commands

Docker Engine 29.9.2, Docker Desktop Linux/aarch64 on the ARM64 Mac. ARM64 runtime is native ISA inside the Linux VM; AMD64 runtime is emulated (the precise emulator implementation was not independently identified). Both builds used Go 1.27.3 Alpine 3.24 and Node 26.13.0 Alpine 3.24, with builder manifest digests captured by BuildKit. Local images only; no registry publication.

```sh
docker buildx build --platform linux/arm64 --load -t sentrame:otel-review-arm64 .
IMAGE_UNDER_TEST=sentrame:otel-review-arm64 PLATFORM_UNDER_TEST=linux/arm64 python3 scripts/container-smoke.py
docker buildx build --platform linux/amd64 --load -t sentrame:otel-review-amd64 .
IMAGE_UNDER_TEST=sentrame:otel-review-amd64 PLATFORM_UNDER_TEST=linux/amd64 python3 scripts/container-smoke.py
```

| Platform | Local Docker image ID | Result |
| --- | --- | --- |
| linux/arm64 | `sha256:712672acabeca97ba6da3fc988b5219d63527e0d73ba9bd63b478351bc5a884e` | Passed full smoke |
| linux/amd64 | `sha256:b7296cd4bba59d38f5c4de38646cad0b22c475469f3db38e7ee4bbd5b0652257` | Passed full smoke on fresh rerun; initial startup failed |

Every successful platform run used a new synthetic bind folder and private random credentials; cleanup removed only its own container and temporary data. Each ran the maintained literal verifier through real ingestion, storage, derivation and protected reads, then reopened the same bind storage and compared native trace operations/events/logs/typed relationships exactly with revision excluded. The container checks do not qualify framework exporters running inside Linux, native NAS capacity, power loss or a Tunnel.

| ARM64 generated body | SHA-256 |
| --- | --- |
| standard-enriched.bin | `ec976dd1a66f856a833d014d08a2cbc61647621e670c1bd5aa5c1e984e0b30a1` |
| standard-future.bin | `821f039cea923393c24b563fdca515b4a13b02663cc1341a4a9a9abda2df5c22` |
| standard-journey-browser.bin | `8020c42aaf4ad09b7ee091d22129fb6ccf47fe3a45c1b57309d1ccfe2f1158af` |
| standard-journey-goal.bin | `02d0e06f4f30dbe75a4f580cd17dd0b29692d32d753af71105dc0f716cce5e06` |
| standard-journey-server.bin | `77efd202b6355ce1f7512a0604079a63214220c5dbbb0629e8e996dc51084dc5` |
| standard-journey-traces.bin | `35ee21cf4790150e153409b152a5f9027cf365f934447b965aa4767623b98603` |
| standard-logs.bin | `9a5d7097df29ca01a447f4df05a7268aacb34b7183224e783107e0891c0f7748` |
| standard-logs.bin.gz | `6c54886502dda14d9b5474b2fefa7959c0f05cbeee6bbc04fa1de2950748f636` |
| standard-logs.json | `dc1140523d4973cc3bd576aa7ee66d716f44364c5ac9c92dffb65d1b3d8031c3` |
| standard-traces.bin | `5a68b15b9ed1f7a38993c53de8e6a99a890b79701d6cc3b2b9670c4e3432989d` |
| standard-traces.bin.gz | `24ec2721f35b7504b6b1ac486585d48d58942e170da7c2bb9fb64acf6e47b4a4` |
| standard-traces.json | `6720c14f3f944f853b5d798452180d638fcd6a482907b481974a82b46b06f311` |

| AMD64 fresh rerun generated body | SHA-256 |
| --- | --- |
| standard-enriched.bin | `dec2592c33b4f7e808916bc39910a1a02e96cc433a47c07e74e5dcf99a983b9e` |
| standard-future.bin | `062f8f9b6c474f4d05ad5c67492f303151c65c2c175fb05a48f6a7859c36c323` |
| standard-journey-browser.bin | `665b2c650046e30080c769e135f832ccb9e18a3a5b413443ec3ae219606b38fb` |
| standard-journey-goal.bin | `afaf669dd8383b089ed185b44a923c56ed077c919fcde025e910e4589c6da8c4` |
| standard-journey-server.bin | `581e3b04904822deba868cc1b5b5eb5d7278f0063d2f47064ed026ec49f39c66` |
| standard-journey-traces.bin | `921ec86d57706f61953730a0f4e22d270d0c9c42fabfdfb8f00058e376793dc4` |
| standard-logs.bin | `d5c6c78a506671d21851fda168206af00d117426a311251d3acfeaab80733701` |
| standard-logs.bin.gz | `99a811b127f2680fa4f316ac4d76cbd0ba0fad3fa7c5496e89c795c824361268` |
| standard-logs.json | `de8b4b17d43be6a20c018896605136b75df42288826f960546e8d0e5ed9809c1` |
| standard-traces.bin | `89965c93ad1259fb76a81b9dbdddddd81236f0b481b82510345ba2b37f3a1ffd` |
| standard-traces.bin.gz | `9a6da77928a17ed315c3bc9c030a38f0383fbd197ad592d1accb58ab6c4b0d51` |
| standard-traces.json | `e6d36f00caa6c1ffccf4c4ac8e7b6ad75a577234ba8e810478c2cd66b44b7185` |

### Final workflow review

The five capability levels remain progressive: (1) ordinary standard receipt/discovery, (2) native trace structure and timing, (3) owner-configured meaning without custom labels, (4) optional literal stage hints with owner precedence, (5) explicit scope/tab/presence only for richer lifecycle or observed-active questions. Trace trees never become chronological sibling journeys; native inclusive operation time never becomes dwell or attention. The voting story preserves failed attempt, retry, trusted persistence, post-commit friction/recovery and separately configured screen completion. Forward definitions/references retain captured meaning.

Home Server capacity still requires owner workload/retention/concurrency and CPU/memory/disk budgets; those are prerequisites, not measured support. Shared Tunnel, Collector, gRPC, external applications, .NET↔Next calls and unlisted transport/runtime variants remain explicitly not run. Roadmap B is deferred. Final native Chrome review of the current build displayed default Traces, capture-time interpretation, exact times/durations, parent-grouped siblings, missing/invalid end exclusions and separate SDK/projection omission disclosure.

`go vet ./...` passed. No final suite remains failed; the initial AMD64 startup failure remains unresolved as an intermittent runtime limitation despite a subsequent full passed run.

### Changed-file inventory

Includes authorized edits layered on existing documentation work. Pre-existing canonical/archive changes are excluded because this implementation did not edit them.

```text
.env.example
CONTRACT.md
README.md
Taskfile.yml
cmd/fixtures/main.go
cmd/fixtures/standard_journey.go
cmd/fixtures/standard_journey_test.go
cmd/fixtures/standard_otel.go
cmd/fixtures/standard_otel_test.go
compose.cloudflare.yaml
compose.yaml
docs/DEVELOPER-GUIDE.md
docs/DOCKERHUB.md
docs/LOCAL-REVIEW.md
docs/NAS.md
docs/README.md
docs/USER-GUIDE.md
docs/integration/AGENT-INTEGRATION.md
docs/integration/OTEL-CONTRACT.md
docs/integration/synthetic-demo.md
docs/integration/synthetic-event-plan.md
docs/planning/sentrame-decisions.md
docs/verification/standard-otel.md
internal/app/diagnostic_privacy.go
internal/app/diagnostics_test.go
internal/app/discovery.go
internal/app/evidence_json.go
internal/app/evidence_json_test.go
internal/app/http.go
internal/app/intake.go
internal/app/invitations.go
internal/app/journey.go
internal/app/lifecycle.go
internal/app/native_journey_test.go
internal/app/native_presence_test.go
internal/app/query.go
internal/app/reference.go
internal/app/stage_hints_test.go
internal/app/standard_mapping_test.go
internal/app/standard_otel_test.go
internal/app/store.go
internal/app/trace.go
internal/app/trace_exploration_test.go
internal/app/types.go
openspec/changes/standard-otel-default-integration/tasks.md
rest/fixtures/standard-accepted-hint.json
rest/fixtures/standard-conflict.json
rest/fixtures/standard-hints.json
rest/fixtures/standard-journey-mapping.json
rest/fixtures/standard-journey-next-mapping.json
rest/fixtures/standard-journey.json
rest/fixtures/standard-mapping.json
rest/fixtures/standard-otel.json
rest/standard-otel.http
scripts/container-smoke.py
web/src/App.svelte
web/src/Traces.svelte
web/src/i18n.ts
web/src/trace.test.ts
web/src/trace.ts
web/src/types.ts
web/src/workspace.ts
```

Final tracking: **23/23 tasks complete**, confirmed by `openspec instructions apply --change standard-otel-default-integration --json` returning `all_done`. Post-checklist strict validation passed again (change plus eight canonical specs); Go formatting, changed-file whitespace/existence and standard JSON syntax checks passed. Owned disposable review services were stopped; retained host sample databases were not reset. No specifications were synced or archived.

## Separately authorized specification sync and archive — 2026-10-29

After implementation, the owner explicitly requested sync followed by archive. Seven capability deltas were reconciled into canonical specs: 30 added requirements and one modified diagnostic requirement; `trace-exploration` is new. Existing titles/purposes, unrelated requirements, optional paired timing, viewer/live boundaries and roadmap B were preserved. All nine canonical specs and the active change passed strict validation before archival. Delta-to-main comparison and a second merge verified complete, idempotent synchronization. The package moved byte-for-byte to `openspec/changes/archive/2026-10-29-standard-otel-default-integration/`; its artifacts and 23/23 completed tasks are preserved. Documentation links now target the archive, and `task spec:validate` validates the canonical specs rather than a stale active-change name. Historical commands/status statements above describe their original implementation checkpoint.
