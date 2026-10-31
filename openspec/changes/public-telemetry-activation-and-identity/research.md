# Research evidence — 2026-10-31

Read-only research at `main`, HEAD `c7288d572811979fecc2f3f856a6a072d85395f1`. No live receiver, registry, credential file or deployed app was accessed. Code references below are relative to the repository root at that HEAD. External links were inspected on 2026-10-31; branch URLs are moving examples, not immutable dependency pins. Findings are architectural inferences unless labeled observed code.

## Executive findings

1. **Keyless for an installer is achievable; authenticated genuine installation is a separate goal.** A public app handle or public write-only token can remove manual key lookup. Downloading the same public permission at activation improves operations, not caller authenticity.
2. **Prefer stable app identity plus `service.version`, not a forever shared secret.** A durable public alias can resolve to current configuration, while independently revocable permission binds to one active app/generation. App name is display metadata and an unauthenticated version is a claim.
3. **Conditional recommendation:** for low-risk opt-in community diagnostics sent to the developer's central receiver, explore owner-published public configuration discovery with strict admission budgets and visibly public assurance. Keep protected server intake separate. If individual installation revocation is necessary, consider enrolled short-lived permission; if verified operator/workload attribution is necessary, require a trusted admission proof. Do not select either until the product questions in design.md are answered.
4. **Recovery is not revocation reversal.** Current complete registry/store backup restores identity and keys together; no existing public rebind API exists. A stale backup can restore revoked authority. The proposed alias needs owner-controlled recovery and a reliable current revocation/deletion record, or must stay disabled.

## Current Sentrame boundaries verified in source

| Observed code / contract | Implication for this exploration |
| --- | --- |
| `internal/app/registry.go:20–35,108–140`: maximum 32 entries; v1 registry stores opaque ID, generation, display name, state, storage and separate browser/server credentials. Validator demands 64-hex role credentials, unique names and tokens, fixed app storage path. | No public alias, enrollment issuer, installation record or revocation history exists in this schema. A stable public alias would be additional state, not a special app name or hand-edited key. |
| `registry.go:433–479`: exact bearer lookup returns only active app plus browser/server role; registration independently generates ID, generation and both tokens. | A key baked in a public artifact would still hold its configured role. Publishing the current server token is unsafe; an untrusted installation's server process does not make it trusted to the central operator. |
| `registry.go:481–516`: rotation replaces one role's credential and commits; rename leaves IDs/history/credentials intact. `151–208`: external-edit detection, atomic file publication and fail-closed uncertain post-rename state. | Existing keys survive restart and can be explicitly rotated, but the old key is not an alias that resolves to its replacement. No arbitrary owner-supplied key restore/import API was found. |
| `registry.go:211–269,292–409`: open-store identity matching, initialized marker, damaged/missing registry refusal, orphan rejection. `518–603`: disable by deleting state before purge and restart resumption. | Do not reconstruct authority from a claimed app name, database or old public identifier. Re-registering a name creates new authority. |
| `internal/app/apps_http.go:12–142`: session authorization, mutation Origin/CSRF checks; owner-only connection read returns both role tokens. Viewer sessions cannot access admin operations. | Public bootstrap must not reuse or expose the owner connection response. No public activation route is present in the inspected API dispatch. |
| `internal/app/http.go:236–280`: ingestion authenticates before proxying, supports `/v1/logs`, `/v1/traces`, `/v1/metrics`, finite concurrency, 1 MiB request and 4 MiB gzip expansion bounds. | Body/concurrency guards exist, but are not evidence of per-public-app rate, Sybil, enrollment or Internet-abuse protection. |
| `internal/app/intake.go:40–68,151–165,215–231,317–334`: consumers revalidate bearer, select app, preserve role-derived source and retain `service.name`/`service.version`. `store.go:42,81,211`: tombstones, page/byte and record guards. | Standard resource metadata already separates versions; it does not authorize routing. New public assurance must survive downstream analysis and replay identity rather than be labeled server. Storage bounds alone cannot preserve quality or fairness under spam. |
| `docs/MULTI-APP.md:33–44` and `openspec/specs/data-lifecycle/spec.md` (“Coordinated persistent state recovery”, “Deletion and replay protection”) | Stopped whole-directory coordinated backup; identity/generation mismatch fails; old snapshots require latest rotation/deletion/revocation reconciliation. No guarantee that old backups remember future revocations. |
| `registry_test.go:153–193,195–244,246–275,297–334`: source assertions cover faulted commit/restart, pending removal, missing registry/orphans, identity and external edits. `apps_test.go` includes rotation/intake concurrency coverage. | Tests were read, not run. They support understanding intended existing boundaries, not runtime qualification of a proposed activation system. |

Current shared intake concurrency is documented as eight; current per-app pre-write byte guard is 896 MiB (`docs/MULTI-APP.md`, `store.go:81`). This is a Home Server contract, not a capacity estimate for arbitrary public installations. Exact later admission rates, active-install counts, event volume, cardinality caps and disk/network budgets need owner input.

## Primary-source patterns and limits

### Sentry public DSN

The maintained [JavaScript DSN parser/serializer](https://github.com/getsentry/sentry-javascript/blob/develop/packages/core/src/utils/dsn.ts) defaults to a public representation without the deprecated password component. Its [envelope URL builder](https://github.com/getsentry/sentry-javascript/blob/develop/packages/core/src/api.ts) sends the DSN public key as `sentry_key`. This is a concrete open-source example of deliberately distributed ingestion configuration. Sentrame should copy the separation of public routing and administration, not Sentry's envelope format or query-string credential placement.

Sentry's [client-key API](https://docs.sentry.io/api/projects/create-a-new-client-key/) requires a separate scoped administrative bearer and supports a key-level error rate limit. This illustrates a management boundary and abuse control; it does not establish genuine event origin or promise protection from every telemetry-abuse class. These source observations do not establish identical DSN access rights across every Sentry endpoint; Sentrame's public permission must explicitly have no read/admin access.

### Application Insights and Azure Monitor

Microsoft's [connection-string documentation](https://learn.microsoft.com/en-us/azure/azure-monitor/app/connection-strings) describes an instrumentation key as a resource association identifier rather than a security key. A connection string identifies destination/configuration; concealing it in a binary cannot turn it into producer identity. This is a useful analogy for a public stable app handle.

[Microsoft Entra authenticated ingestion](https://learn.microsoft.com/en-us/azure/azure-monitor/app/azure-ad-authentication) is a different mode: identity and resource-scoped publishing authorization are configured, and local authentication can be disabled. The Application Insights JavaScript SDK is listed as unsupported for that mode. This is not anonymous zero-setup ingestion for arbitrary third-party machines.

The current [Azure Monitor OTLP Collector guide](https://learn.microsoft.com/en-us/azure/azure-monitor/containers/opentelemetry-protocol-ingestion) documents ordinary OTLP exporters with Azure authentication and explicitly requires Entra authentication for that Collector path. Do not infer that an Application Insights instrumentation key alone authenticates native Azure Monitor OTLP. The guide has version-specific extension examples; no Azure adoption or credential setup is proposed here.

### Ordinary OpenTelemetry settings

The [OTLP exporter specification](https://opentelemetry.io/docs/specs/otel/protocol/exporter/) defines base/signal-specific endpoints, headers, protocol, compression and TLS options. A bootstrap response can feed these through app startup configuration or a standard Collector; this is an integration inference, not a built-in universal discovery feature. Ordinary exporters do not acquire Sentrame activation configuration by themselves.

The maintained [Go HTTP exporter options](https://github.com/open-telemetry/opentelemetry-go/blob/main/exporters/otlp/otlptrace/otlptracehttp/options.go) provide endpoint and header configuration. Dynamic refresh still needs qualification per language/version; immutable startup options may require exporter recreation or a supported Collector configuration path. No custom SDK fork is needed for static endpoint/bearer settings, and no automatic hot reload is claimed.

The [service semantic conventions](https://opentelemetry.io/docs/specs/semconv/resource/service/) define `service.version` as a component version string; format is not prescribed. Sentrame's multi-service app identity should remain registry-owned, with service versions as bounded optional metadata. One product release may contain several independently versioned services; product-release mapping remains an owner choice, not a new mandatory OTLP attribute.

### Signed enrollment and proof of possession

[SPIRE configuration](https://spiffe.io/docs/latest/deploying/configuring/) distinguishes node attestation, workload attestation and trust domains. Its join-token path requires an out-of-band server-supplied single-use token before issuing identity. The maintained [SPIRE agent API source](https://github.com/spiffe/spire/blob/main/pkg/server/api/agent/v1/service.go) validates/consumes join tokens or performs attestor challenge-response before signing a CSR. This illustrates the missing trust bootstrap: signing a client-generated public key does not make arbitrary signup trustworthy. SPIRE itself is likely too much machinery for simple community telemetry and is not recommended as a dependency here.

[GitHub OIDC](https://docs.github.com/en/actions/concepts/security/openid-connect) and [cloud-provider configuration](https://docs.github.com/actions/deployment/security-hardening-your-deployments/configuring-openid-connect-in-cloud-providers) show short-lived signed issuer claims checked against configured workflow trust. That can authorize a controlled CI/deployment environment; a downloaded open-source image does not inherit the workflow's identity. Never bake a CI token into a release. A release signature verifies an artifact for a verifier; it does not make arbitrary requests authentic or telemetry truthful.

[RFC 9449](https://datatracker.ietf.org/doc/html/rfc9449) distinguishes sender-constrained tokens from bearer possession and specifies freshness, request binding, nonce and replay considerations. It demonstrates key possession, not genuine app provenance or honest measurements. Any selected enrollment could borrow a standard pattern for its control plane; requiring per-OTLP-request DPoP would not be satisfied by ordinary static exporter headers and is outside the recommended basic path.

## What a signature over HTTPS establishes

TLS authenticates the configured receiver and protects transport when correctly validated; it does not authenticate an arbitrary installer. A signature with signing material in a repository/container is reproducible by anyone. A fresh locally generated key can establish continuity of one key holder, but an attacker can generate many keys. A trusted-server signature establishes that the issuer approved specified claims under its admission policy; assurance is only as strong as that policy and protection of the issuer. None proves a real user action, honest event, or exact install count. Public key *verification* material can safely ship; private signing authority cannot safely be shared publicly.

CORS constrains cooperating browsers, not scripts, servers or forged Origin headers. Neither CORS, an app name, a version string nor rate limiting is producer authentication.

## Evidence limitations

This is a source/docs-only exploration, not a security audit or load test. No credentials, enrollment keys or signatures were generated, and no production endpoints or shared resources were exercised. Sentry developer overview returned an unsupported Markdown fetch error; the maintained SDK source and administrative API above supply the used evidence instead. No local server/UI was started; the prior ephemeral loopback EPERM was not retried. External main/develop links may change; recheck and pin relevant revisions if an implementation is authorized.
