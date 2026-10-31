# Design — exploration, not an implementation decision

## Context

See proposal.md for motivation and research.md for dated primary evidence and verified current-code references. The design is included because identity, public exposure and recovery introduce security/migration ambiguity. Alex explicitly requested research with unresolved decisions: this document intentionally presents conditional architectures rather than resolving material product choices by assumption.

Current Sentrame authenticates standard OTLP with private-registry browser/server bearers. Owners may reveal, rotate and recover credentials from a coordinated private backup. No anonymous/public configuration endpoint, alias rebind, installation enrollment or trusted issuer is currently implemented. The separate Tools and two earlier exploratory changes are outside this scope.

## Goals / Non-Goals

**Goals:** remove manual key acquisition for consented open-source installations; preserve normal SDK/exporter settings; keep aggregation stable across releases; make caller assurance and operational cost explicit; define safe identity recovery.

**Non-goals:** authenticating arbitrary downloaded software by an embedded secret; hiding public signatures; generating keys or credentials during research; implementing endpoint/UI/registry changes; selecting an API schema, mandatory cloud platform or per-record signature protocol; enabling production traffic or evaluating exact WFS exports.

## Decisions and options

### Fixed exploration boundaries

These derive from the user request and existing contracts: distributed signing material is public; public configuration grants ingestion only; it never exposes the protected server bearer, reads, invitations or owner APIs. Name/version claims cannot create or take over an app. Authority must remain distinguishable through derivation and summaries. Consent and bounded failure are required. These boundaries do not choose an enrollment architecture.

### Separate identity, metadata and permission

A potential model has five independent concepts, with storage/schema deliberately undecided:

| Concept | Lifetime and purpose | Authority |
| --- | --- | --- |
| Registered app identity | Stable product workspace across releases and rotations | Owner-controlled binding; current registry ID is a possible anchor |
| Intake generation | Dataset/replay boundary; changes for explicit replacement | Never crossed merely by a display-name match or cached configuration |
| Public activation handle | Durable discoverable alias to published configuration | Routing only; possession does not prove ownership |
| Ingestion permission | Current scoped write capability; revocable/rotatable independently | Public/unverified or admitted assurance, never automatic server trust |
| Service/version and optional installation reference | Version comparison and, if needed, installation revocation | Claims unless independently verified; not person identity or a secret |

Keep one app for several releases/services, bounded `service.version` metadata and unknown/missing versions. Do not allocate one registry app or credential per version. A single shared public token also cannot yield reliable unique installation counts. Persistent installation references are correlatable data; collect them only if required by an approved purpose, disclose their lifetime and avoid fingerprinting.

### Compare complete options

| Option | Installer friction and receiver | What is authenticated | Limits, rotation and recovery |
| --- | --- | --- | --- |
| A. Anonymous intake with public app route/identifier | Very low after consent; developer operates central receiver | At most selected public routing namespace | Anyone can spoof app/version/events; no per-install revocation; disable route or change public handle. Changing handles disrupts baked old clients. |
| B. Baked long-lived public write-only token/DSN | Very low; central receiver | Possession of a public write capability, not genuine installation | Simpler ordinary header setting; public token can be scraped and reused forever if never revoked. Rotate/disable with client update or bootstrap. Durable alias recovery must not revive revoked token. |
| C. Activation fetch of current public configuration | Low, consent plus one HTTPS discovery; central receiver | Configured receiver authenticity via TLS; caller still anonymous | Stable handle can survive endpoint/token changes; new network/cache/refresh dependency. Attackers can fetch each replacement token, so rotation cannot eject them. Still needs budgets and emergency disable. |
| D. Per-installation key/enrollment | Low-to-medium if automatic; may need approval/login or deployment bootstrap | Fresh key proves possession; subsequent signed enrollment/refresh can prove continuity | Can revoke/quota one enrolled identity; open enrollment is Sybil-prone and not genuine-app proof. Protected local key storage, expiry, refresh and lost-key recovery add state/operational cost. Self-generated keys alone offer weak attribution. |
| E. Trusted-server-signed enrollment | Higher unless installer already has an accepted identity; central receiver plus protected issuer | Approved issuer's admission claim and optionally installation key binding | Can bind app/generation/audience/scope/expiry and individually revoke. Trust comes from operator approval, account/deployment identity or attestor, never app name. Issuer compromise, policy and key rotation become obligations. Honest events are still not proven. |
| F. Installer-owned Sentrame/Collector | Highest receiver setup, or modest if already operated; installer owns data/access | Local receiver's configured authorization | Stronger control over destination, retention and access, no built-in central aggregation across independent receivers. Existing ordinary exporter settings apply; setup can later automate provisioning without public central credentials. |

A, B and C have effectively the same producer-authenticity ceiling. C is an operational refinement to B/A. D with no admission gate is managed pseudonymous intake, not trusted enrollment. E can be proportionate for controlled deployments, but cannot promise arbitrary open-source installs are authentic. F is a different product/data-ownership choice, not a transparent substitute for central product analytics.

### Conditional practical recommendation

If Alex wants a developer-operated central receiver for small-volume, low-risk community diagnostics, prefer exploring **C: opt-in public configuration discovery by a stable owner-published handle**, with a public write-only permission and isolated abuse budgets/assurance. This best fits the requested app activation experience and configuration continuity. B is a simpler offline-friendly fallback when changing endpoint/token requires a release and public-token abuse is acceptable. Do not advertise any token as irrevocable or guaranteed forever.

If individual installation revocation is needed, evaluate D with an explicit admission policy and disclose its Sybil limits. If decisions require verified operator/workload provenance, use E with protected issuer/admission evidence or retain F/private deployment configuration. Do not add cryptography merely to make anonymous telemetry appear trusted.

A gateway is one candidate for separating public admission from the Home Server; it is not selected. It must not forward all accepted public records under the existing protected server token and thereby launder their assurance. Direct new public authority or a separately constrained gateway path would require explicit changes to intake/derivation contracts and qualification.

### HTTPS activation sketch, deliberately not an API contract

After consent, an app could send its published handle plus bounded app/service name and version to a known HTTPS configuration origin. Owner registration/publishing happens separately; untrusted lookup cannot auto-register or reserve a name. The response could carry a configuration revision, validity interval, active app/generation binding, approved OTLP endpoint/protocol/header settings, public scope/assurance and policy information. Exact method, URL, field names, token type and cache duration remain undecided.

A public response contains only intentionally public material. Enrolled responses require independent admission evidence and must not be shared through public/CDN caches. Rate limits cover negative lookups as well as known handles; generic errors should not expose private app inventory. Changes to destination/operator/data policy require consent review; bootstrap is not arbitrary remote configuration or authorization to redirect to any URL. Avoid redirects or credentials across origins, downgrade to HTTP, unbounded response size and arbitrary executable config.

Feed accepted settings into ordinary OTel initialization or a supported Collector. A tiny app lifecycle integration may fetch and apply settings; that is not a custom SDK fork. SDK-specific environment encoding and signal endpoint semantics still require tests. Native HTTP/Protobuf is the current Sentrame starting point; public HTTPS does not qualify gRPC, direct browser CORS or the shared tunnel route.

### Cache, offline and enrollment replay

Proposed candidate behavior: keep bounded last-known-valid configuration bound to consent, receiver, app and generation; refresh with timeout, jitter and backoff; do not delay core app startup. Use only unexpired settings when discovery is unavailable. No initial configuration means export stays disabled. Disk buffering is optional and bounded by age/bytes, with minimized payload and explicit loss reporting. Expiry/revocation suspends export; a 401/403 should trigger at most bounded refresh rather than a tight loop, and overload respects retry guidance. Exact TTL, refresh lead time, maximum offline interval and drop policy are product/operational decisions.

A stable public handle automatically discovering replacement public permission improves update delivery, but an attacker can do the same. If the product wants to exclude a revoked installation, public refresh must not hand it fresh privileged permission. Public alias disablement and per-install revocation are different controls.

For D/E, enrollment proofs should bind trusted issuer/key, audience/receiver, app/generation, installation public key, scope and expiry. One-use enrollment authorization, a fresh challenge or bounded nonce/identifier replay cache must be selected with limits and crash/retry semantics. Idempotent enrollment retry must return the same authorized result or a clear recoverable status, not mint unlimited identities after a lost response. A client-signed fresh challenge proves that key's possession; trusted admission still comes separately. Private issuer keys stay on the trusted server; only verification material is distributed.

These are control-plane concerns. A short-lived scoped bearer can use ordinary exporter headers; standard mTLS or a supported Collector auth extension may be alternatives where qualified. Per-record/per-request DPoP or a custom signing transport is not required. Enrollment anti-replay does not deduplicate OTLP. Existing signal-specific source identity and generation/tombstone horizons still apply; metrics remain receipt-only. Retrying a batch under a replacement generation must never make old data new.

## Short threat / abuse analysis

Assets are private data/admin authority, app identity/history, trusted-source interpretation, resource budgets and installer privacy. Adversaries include arbitrary Internet callers, another app's installer, modified clients and compromised admitted installations. The distributed app and payload metadata are outside the operator's trust boundary.

| Threat or failure | Candidate boundary/control | Residual risk |
| --- | --- | --- |
| Scraped public token, shared signing material, forged app/version/server events | Public assurance, strict role separation, binding independent of metadata; protected server credential never returned | Public observations remain spoofable; signatures using public private material add no assurance |
| App-name squatting, cross-app IDs and dedup poisoning | Owner-controlled publish/rebind, opaque app binding, isolated public namespace through summaries/links/replay | A public app's own charts can still be polluted; one token cannot partition attacker from real installer |
| Spam, distributed quota drain and starvation | Per-app and global request/byte/event rates, fair concurrency, finite enrollment/cache state, emergency off switch; IP/network heuristics supplementary only | Sybil/distributed senders bypass per-install/IP limits; central operator may choose to lose legitimate telemetry to contain cost |
| Cardinality, tiny-record floods, decompression/parser/CPU/disk pressure | Finite request/expanded bytes, record/attribute/string/bucket bounds; unique service/version/event/installation dimension caps; schema policy before expensive work; daily storage/retention budget | Current finite storage limits do not provide fair admission or establish sufficient public capacity; malicious fields require sanitization in every UI/log path |
| Replayed enrollment or copied key/token | Expiry/audience/app/key binding, one-use admission/freshness with bounded replay state, issuer revocation | Authorized or compromised clients can still fabricate fresh events; bearer theft works until expiry/revocation |
| Old backup, deleted app name or stale queued configuration | Owner-reviewed identity rebind; current deny/deletion state; coordinated stores/registry and recovery generation checks | Latest revocations cannot be inferred from an old snapshot; fail closed if reconciliation is unavailable |
| Consent drift or privacy leakage | No pre-consent discovery, minimization before export/persistence, no person/fingerprint fields, destination policy checks; finite queue and disable | Opt-in installation telemetry is selected, incomplete and not population conversion; previous disclosed copies cannot be recalled |

Numerical budgets and deployment location must be agreed against expected workload before exposing intake. An unguessable *public* identifier is not a capacity control; CORS, sampling and retention are supplementary measures, not authentication or hard admission quotas. No load/security qualification was performed.

## Migration and recovery candidates — not deployment instructions

No migration runs in this change. Existing senders and owner/viewer APIs remain unchanged. A later selected design would need an explicit compatibility delta, versioned registry/alias state, consent UI, assurance propagation and per-app admission policy. Do not repurpose the current browser token without checking browser-specific analysis semantics, or publish a server token merely because an installer runs a backend.

For coordinated restore, retain registry, stores, generation/deletion state, any alias ownership/bindings, enrolled issuer state and latest revocation/configuration revision under appropriate protection. Public aliases need integrity/ownership protection even if not confidential; private role/issuer/installation material needs confidentiality. A stable external configuration hostname is also an availability/ownership dependency, not a promise guaranteed by an alias.

Recovery cases to decide explicitly:

- Same active identity and matching complete backup: reconcile latest deny state before enabling, preserve history and supported replay barriers.
- Lost key with intact identity: owner may issue replacement permission under current policy; restore identity without reactivating revoked key.
- Lost registry or mismatched stores: follow existing fail-closed coordinated restore; public bootstrap cannot reconstruct authority.
- Deleted/tombstoned identity or stale snapshot: keep disabled; same-name new registration is a separate generation. Do not silently retarget old aliases or old queues.
- Endpoint/operator change: owner-authorized binding and consent/destination policy review; no automatic cross-owner history merge.
- Rollback: disable new public intake first; preserve newer deny state. Older code cannot be assumed to understand new registry state. Restoring an old image/snapshot alone is not revocation-safe rollback.

## Product decisions required before an implementation proposal

These questions intentionally remain unresolved and block architecture selection, not completion of this research:

1. **Whose receiver should receive telemetry:** the app developer's central Sentrame, each installer's own Sentrame/Collector, or an explicit choice? Who operates and pays for admission, storage and incidents?
2. **What scale and trust are expected:** approximate active installs, events/bytes per day, burst rate and budget; community diagnostics or evidence for business/persistence decisions? Can arbitrary public callers consume the entire public quota?
3. **What opt-in flow is acceptable:** installer/admin consent UI, deployment flag or both; which exact minimized fields/purpose; withdrawal/queue handling; whether destination/policy changes need renewed consent?
4. **What attribution confidence is sufficient:** public app/version claims, continuity of a pseudonymous installation, approved operator or attested workload? Is per-install revocation worth an enrollment account/approval/trusted deployment step?
5. **What does “restore/rebind forever” mean:** retaining the stable product/config alias after disaster, rebinding a matching active generation, moving receivers/owners, or reactivating a previously revoked token? The last operation is distinct and must not be implicit. Which current revocation/deletion authority survives loss of the main registry?
6. For the chosen model, what offline/cache interval and bounded loss are acceptable, and does a product release map to one service version or several independently versioned services?

Public API paths, exact cache values, registry schema, token format, gateway placement and issuer technology remain unselected. Do not infer consent or a design decision from these suggested defaults or from planning validation.
