# Home Server deployment contract — planning only

## Established direction

The user chose a Home Server deployment using Docker Compose, an associated Compose file and `.env`, behind an existing shared Cloudflare Tunnel providing HTTPS. This document specifies that boundary only. It creates no Compose file, `.env`, container, tunnel, credential, DNS record, certificate or network change. The owner reports sufficient RAM and CPU options; exact hardware, numerical budgets and deployment fit remain unmeasured.

## Proposed small deployment contract

| Concern | Proposed contract to review before implementation |
| --- | --- |
| Delivery | A versioned Compose file and matching documented `.env` configuration surface form the setup. Pin compatible application/optional Collector images; static frontend assets ship with the app. Do not add a second frontend runtime by default. |
| Configuration | Document public HTTPS base URL, allowed origins/trusted proxies, application/producer bindings, viewer-auth secret inputs, local persistent storage, resource/intake/query/retention bounds and optional Collector routing. Missing/invalid required settings fail clearly without printing secrets; no production credential defaults. Exact variable names are future implementation detail. |
| Existing shared tunnel | Reuse the existing tunnel route/connector under an explicitly reviewed integration. Sentrame does not provision or take ownership of a second tunnel, its token, sibling routes or shared connector lifecycle. Count incremental tunnel load and shared failure impact. |
| Origin reachability | Describe how the existing connector reaches the Compose origin: an explicitly selected internal network/address and port. Do not assume `localhost` in one container means another container, publish database/admin ports, or silently attach networks. Exact network/hostname settings remain an integration choice. |
| TLS boundaries | Browser/exporter to Cloudflare uses the intended HTTPS route. Cloudflare-to-connector protection and connector-to-app origin transport are separate links. Select/document HTTP on the controlled local path or HTTPS with certificate validation; edge HTTPS does not prove the origin leg uses TLS. No certificate-validation bypass is prescribed. |
| Application access | Passphrase/cookie viewer grants, scoped invites, producer authorization and CSRF/origin checks remain Sentrame responsibilities. A tunnel or HTTPS alone is not app authentication. Cloudflare Access is not assumed configured and must not be substituted for the agreed app access model. |
| Persistence | Use a genuine local Home Server volume with documented ownership, durability, capacity, expiry/deletion and backup policy. Container replacement/restart must not silently discard evidence or revive deleted scope. Do not use a network-file SQLite topology by implication. |
| Operations | Document health/readiness, resource bounds, configuration reload/restart semantics, retention cleanup and rollback/data compatibility for the chosen version. No secret/raw-telemetry dumps in Compose or tunnel logs. Define the failure model; neither Compose restart policy nor tunnel encryption is a durability guarantee. |

The minimal candidate remains one app process/container with static assets and embedded storage/receiver if qualified, or an app plus maintained Collector when that is simpler. The already shared connector is additional infrastructure, not a free process or an implicitly managed service in Sentrame's Compose project.

## Qualify the routed path separately

| Path | Proposed qualification contract; not a passed test |
| --- | --- |
| Viewer HTTP and cookies | Confirm the real HTTPS hostname, secure cookies, allowed host/origin, trusted forwarded metadata and CSRF/redirect behavior. Untrusted forwarded headers cannot grant access or select another app. |
| Live SSE | Use `Content-Type: text/event-stream` and authenticated non-cacheable responses. Verify timely flush, heartbeat/idle handling, reconnect replay/reset, count consistency and active revocation through the shared route. Direct-origin success is insufficient. No Quick Tunnel substitute is assumed. |
| OTLP/HTTP | Qualify binary and JSON Protobuf, gzip/none, exact signal paths/content types, request/body/time limits, producer auth, explicit errors/retries and durable acknowledgement through the route using standard exporters. Browser CSP/CORS/lifecycle checks remain distinct from server exporter support. |
| OTLP/gRPC | Keep the receiver/origin target distinct from published-route support. Current Cloudflare guidance does not support Tunnel public-hostname gRPC; do not advertise that route as gRPC-compatible. A qualified private/origin or separately reviewed standard Collector-forwarding path can serve applicable producers. No alternate network route is created by this proposal. |

No SSE/OTLP path has been exercised here. The chosen route, shared policies, versions, proxy limits and actual exporter must be recorded in later qualification evidence. Do not infer browser gRPC support, security certification or privacy compliance from HTTPS.

## Targeted primary-source checks

- Cloudflare's [published-application protocols](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/routing-to-tunnel/protocols/) distinguish an HTTPS public request forwarded to an HTTP origin from an HTTPS origin. [Origin parameters](https://developers.cloudflare.com/cloudflare-one/networks/connectors/cloudflare-tunnel/configure-tunnels/origin-parameters/) govern connector-to-origin behavior, separately from browser HTTPS.
- Cloudflare's [gRPC limitations](https://developers.cloudflare.com/network/grpc-connections/) currently state that Tunnel supports gRPC through private subnet routing, while public-hostname deployments are unsupported. Ordinary Cloudflare reverse-proxy gRPC support is not proof of Tunnel support.
- Cloudflare's [Tunnel troubleshooting](https://developers.cloudflare.com/tunnel/troubleshooting/) states that streaming responses are buffered unless the origin uses `text/event-stream`. Its [Quick Tunnel limits](https://developers.cloudflare.com/tunnel/get-started/quick-tunnels/) exclude SSE. These facts guide the future test; they do not prove this shared tunnel's behavior.

The tunnel adds a processing boundary for telemetry and diagnostics. Sender minimization where supported, permitted attributes and logging policies should be reviewed before real export; filtering only after arrival cannot establish that fields never traversed the provider. No account settings, traffic or policies were inspected or changed during this documentation task.
