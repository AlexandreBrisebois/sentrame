# Independent challenge and current disposition

An independent GPT-6 Astra review challenged the initial capability draft against primary sources and the WFS example. Thirteen findings changed the planning contracts below. “Corrected” means documentation changed; runtime behavior, security, accessibility, performance and OTLP compatibility remain unimplemented and unqualified. The original review and its hashed baseline are preserved in the portable Library review package. No second broad review of the corrected package is claimed.

## Findings and corrections

| Finding | Current contract / remaining qualification |
| --- | --- |
| Receiver partial acceptance | Ordinary maintained receiver consumers do not provide per-item partial counts. Initially prevalidate and commit/reject each direct storage request atomically, with explicit permanent versus retryable statuses. Upstream queue/batch acknowledgement has different durability/envelope boundaries. Pin and qualify the actual path. |
| Shared agreement versus visit outcome | Client completion-screen observation is a visit outcome, distinct from committed vote, emitted reason and who caused shared agreement. Another participant, initial convergence, repeated hydration, departure and post-commit publication errors have separate cases; current telemetry is not assumed. |
| Identity collision | App + trusted configured producer namespace + signal + supported source identity scopes replay. Logical correlation, record/action identity, authority and causal ordering remain distinct. |
| Invented paths/rates | Hidden X in A-X-B cannot create a direct A-B edge. Declare cohorts, clipping, follow-up, repeated-step units, missing starts and mixed versions. |
| Deletion/replay | Use stable identifiable app/logical scopes and explicit correction/tombstone horizons. Mutable inferred journey IDs alone are insufficient; unknown IDs/missing timestamps cannot provide unlimited original-age guarantees. |
| Metrics arithmetic | Standard receipt remains a target; initial metrics use is bounded sanitized inspection. Arithmetic requires qualified stream identity, type, temporality, interval and reset/replay semantics. Cumulative 10→12 is not 22. |
| Seven-day baseline | User intent resolves this: selected-app accumulated seven-day traffic with a chosen live visit overlay, plus a clean-live toggle. Compatible filters and separate counts remain; it is not same-visit history. |
| Time-only SSE/races | Consistent snapshot/cursor revisions, bounded invalidation without new input, restart/reset and old-filter isolation. No stream-long DB transaction. |
| Diagnostic occurrence counting | Received records, qualified exception occurrences and affected scopes have distinct units. A grouping fingerprint cannot prove two carriers are one exception or two failures are duplicates. |
| Viewer mutation/disclosure | CSRF/origin/host protections, atomic invite redemption, rotation/revocation of active streams and inert untrusted text. Detail authorization is the disclosure boundary; a disabled copy button cannot prevent copying disclosed data. |
| Durability | Declare process/OS/power-loss scope and SQLite/storage assumptions. Restart verification does not establish power-loss durability; no destructive test is requested here. |
| Qualification dependency cycle | Separate G3a supplied envelope, G3b later authorized narrow tracer and G3c integrated qualification. Selected stack is not measured fit. |
| Tool-version evidence | Generated skills identify OpenSpec 1.14.3; the observed validation executable is 1.15.0. Structural validation proves document structure only. |

The [official receiver](https://github.com/open-telemetry/opentelemetry-collector/tree/main/receiver/otlpreceiver), [batch processor](https://github.com/open-telemetry/opentelemetry-collector/tree/main/processor/batchprocessor), [OTLP specification](https://opentelemetry.io/docs/specs/otlp/) and [SQLite synchronous documentation](https://www.sqlite.org/pragma.html#pragma_synchronous) support these qualifications. The research documents retain deeper source evidence and alternative-project analysis.

## Decisions that supersede historical review questions

The accepted [release scope](../planning/sentrame-release-scope.md) settles the five main choices: complete WFS release A with B retained; first-event/goal/next-event boundaries with 30-minute inactivity and per-tab browser isolation; later narrow one-live-session invitations; 30-day detail/12-month summaries; Svelte/Vite + Go + SQLite. A list-only tracer is an intermediate increment, not the final map release.

Later direction also establishes Home Server Compose/`.env` deployment behind shared Cloudflare Tunnel HTTPS; owner-emitted events and discovered goal selection; client-only/server-only/combined correlation; saved reference views; time-at-stage node sizing with traffic-weighted edges and timing safeguards; and discovery-led integration-agent instructions that propose missing events, warnings, errors and recovery before developer-approved edits. These supersede old terminology and unresolved ballots in the historical report. Public Tunnel gRPC and routed SSE remain separate qualification concerns; TLS does not replace app authorization.

## Remaining focused design review

Choose receiver packaging/pins, numerical request/queue limits, durability scope, lateness/skew and source-identity policy, actual event names/active goal and privacy enablement. Specify timing activity/fallback/slower percentile and visual scale; summary dimensions/distributions/precision/follow-up/correction/deletion/backup horizons; later invite lifetime/use; deployment origin settings and real resource envelope. Do not reopen settled product choices.

Before custom work, compare the same WFS questions with a periodic event report or existing OTel/analytics infrastructure. Stop or adapt if no decision changes, acceptable evidence requires forbidden identity, maintained alternatives cost less, or operating cost exceeds measured value. The [feasibility gates](../research/sentrame-feasibility.md) remain proportionate stop/go checks.

Planning reconciliation and a documentation-only draft PR are approved. No app/harness/Compose implementation, external instrumentation, canonical sync/archive, merge, deployment or production test is part of this package. All implementation tasks remain unchecked.
