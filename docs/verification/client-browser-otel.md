# Single-key client browser interoperability — 2026-11-01

**Passed locally:** actual native Microsoft Edge → standard OTel browser SDK logs/traces → disposable same-origin forwarder → changed Sentrame receiver → SQLite → protected journey/trace/error reads. This is a new single-key/schema-4 cell, not reuse of the historical two-role qualification.

## Scope and setup

All services, apps, keys, passphrases, visit IDs and data were synthetic, under `/tmp/sentrame-browser-client-w895lkxv`. The receiver used `127.0.0.1:18196`; the forwarder bound an ephemeral loopback port. A second disposable UI clone used `127.0.0.1:18197`. No real app, retained data, deployed key, Home Server, shared Tunnel or security setting was accessed or changed. All test services and launchers were stopped after verification.

The browser sample used the existing installed official OTel packages: API 1.10.0; SDK trace/logs/resources 3.0.0; HTTP Protobuf trace/log exporters 0.300.0. Esbuild bundled their browser paths. No new dependency, SDK fork, custom telemetry protocol or Sentrame producer wrapper was introduced. The tiny Node forwarder was test infrastructure, not server telemetry instrumentation or a production gateway recommendation.

Browser exporters used ordinary same-origin signal URLs. The forwarder held the one synthetic app key privately and supplied its Authorization header to the receiver; neither that key nor the owner passphrase appeared in the HTML or browser bundle. No browser role or server role existed.

## Observed results

| Check | Result |
| --- | --- |
| Loaded page before Begin | Zero OTLP requests, no provider initialization |
| Explicit Begin | One initialization; subsequent controls used the same providers |
| Standard exporter flush | Browser UI reported successful logs/traces flush; actual Protobuf bytes were accepted, with neither JSON nor gzip encoding |
| Retained observations | Four logs and two native client spans; all six records attributed to `client` and the registered app |
| Journey/goal | One completed journey: Home observed → Search observed → Goal observed; `journey-client-v5` |
| Error and continuation | Controlled `SyntheticClientError` retained and found through protected diagnostic reads; it did not terminate the journey before the observed goal |
| Native structure | Search span parent matched the page span; trace IDs matched and native end exceeded start |
| Replay | Reposting the actual SDK-exported binary batches left six records and one completed goal |
| Receiver process restart | All six retained record IDs, names, exact times, ends and native identities remained equal |
| Browser export after restart | The same running SDK/providers and URLs exported another error without reconfiguration: seven retained records, still one completed journey/goal |
| Forwarder wrong Origin | HTTP 403; no upstream forwarding |
| Forwarder body >1 MiB | HTTP 413; no upstream forwarding |
| Direct receiver attempt | Rejected in the browser under the sample's self-only `connect-src` CSP; the same-origin route remained usable |
| Secrets | Synthetic app key and owner passphrase absent from public HTML/bundle and protected telemetry projection |

The machine became locked before fresh responsive UI acceptance could finish. The attempted native browser state read reported that automatic unlock failed and manual unlock was needed. No further UI action was attempted while locked. This did not invalidate the completed SDK/exporter and protected-read checks, which had already passed.

## Evidence and reproducibility

[Sanitized result](client-browser-otel-result.json) records package versions, counts, source/bundle hashes and restart results. The disposable sample contains `browser.js`, `server.mjs`, `public/index.html`, `public/app.js`, captured standard SDK protobuf bodies and a credential-free result file. Private configuration and generated data are intentionally not copied into the repository.

Local orchestration/verification scripts remain `/tmp/sentrame-browser-client-start.py`, `/tmp/sentrame-browser-client-verify.py`, `/tmp/sentrame-browser-client-restart.py` and `/tmp/sentrame-browser-client-ui-start.py`. They are disposable test helpers, not a maintained deployment or integration product. The first launcher exited before browser connection; keeping the launcher alive corrected that setup issue before the passed run.

## Limits and user-owned follow-up

This qualifies the stated local browser logs/traces and same-origin HTTP Protobuf/no-compression path. It does not qualify a deployed forwarder, direct cross-origin export, browser metrics, other encodings/compression, other browsers, mobile suspension, abrupt termination/beacon delivery, lossless background flushing, public abuse capacity or a production quota. The 1 MiB boundary is this test forwarder's explicit limit, not a new global Sentrame limit.

Home Server/Tunnel deployment checks require separate operational authorization. Native phone/touch checks require an available unlocked environment/physical device. Existing Go protocol/backpressure and local container qualification remain separate evidence; local SDK success is not a claim about those unrun deployment/device cells.
