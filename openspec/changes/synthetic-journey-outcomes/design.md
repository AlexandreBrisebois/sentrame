# Design

## Context

See proposal.md for motivation and scope. This is an initial review draft, not implementation authorization. Requirements describe proposed behavior; they do not certify existing behavior.

Sentrame has a README, license, local OpenSpec configuration and generated skills. No application, Taskfile, CI or development harness was found. The existing untracked `.agents/` and `openspec/` setup is user work and remains untouched outside this new change.

The name combines sens and trame: meaning, direction and the thread connecting events. The working tagline is “Find the meaning between events.” Neither naming nor this proposal establishes domain or trademark clearance.

### Candidate integration evidence

Source inspection through the GitHub connector on 2026-10-25 found:

- [Invitation](https://github.com/AlexandreBrisebois/whats-for-supper/blob/main/pwa/src/app/(auth)/invite/page.tsx), blob `6aa4f3a7f49c23c6c0a6b85c38f4d90fe4327445`: invitation can select a member automatically and defaults to Home. Manual member selection and Discovery are not mandatory next steps. Invitation URLs contain sensitive parameters.
- [Discovery](https://github.com/AlexandreBrisebois/whats-for-supper/blob/main/pwa/src/app/(app)/discovery/page.tsx), blob `98c77f8581745b8a977fe820abf98adde70c5332`: cards update optimistically; vote requests run afterward. Shared convergence and Planner navigation exist but are different from an individual's first contribution.
- [Vote service](https://github.com/AlexandreBrisebois/whats-for-supper/blob/main/api/src/RecipeApi/Services/DiscoveryService.cs), blob `cfb43b2a611e50b1a56047625f15b7eedbb5d447`: persistence precedes follow-up publications, so a later error need not imply that saving failed.

These observations are static source evidence from the preceding assessment, not runtime qualification, a complete telemetry audit or proof of deployed instrumentation. Recheck source before any separately approved WFS work.

## Goals / Non-Goals

**Goals:** A local synthetic producer, bounded ingestion, durable evidence, deterministic outcomes and one inspectable summary. Make completion and uncertainty useful before broadening the product question.

**Non-Goals:** Building a telemetry platform, automatic discovery of user intent, generalized funnel configuration, hosting infrastructure or instrumenting any real application in this change.

## Decisions

### Status and review baseline

User-established constraints: design before implementation; synthetic first; OTel foundation; default-off opt-in for the first integration; no personal content; no replay, fingerprinting, cross-site tracking or AI dependency. Short successful interactions must remain successes.

Proposed defaults in this draft: Discovery start to first persisted vote; direct local ingestion; seven-day retention; no cross-attempt identity; one fixed outcome rule. These defaults are concrete enough to review and revise, but are not recorded as separately approved user decisions. The request to build the initial spec authorizes these planning artifacts only.

### Foundation and signal roles

Use structured OTel log records as the product evidence carrier, with a small Sentrame schema. Reuse official protocol definitions, SDKs and exporters instead of writing a telemetry transport. Traces can later connect a product action to technical execution, but a human journey is not one long trace and sampled traces are not the counting source. Metrics summarize results and collection health; they cannot reconstruct attempts. No attempt IDs in metric labels.

OTLP standardizes transmission, not product meaning, consent, durable application commits or exactly-once processing. Sentrame owns evidence validation and outcome rules. Reference: [OTel signals](https://opentelemetry.io/docs/concepts/signals/) and [OTLP specification](https://opentelemetry.io/docs/specs/otlp/), consulted in the preceding assessment.

### Architecture comparison

| Concern | A: Direct local ingestion — proposed | B: Collector in front |
| --- | --- | --- |
| Ingestion | OTLP/HTTP protobuf logs endpoint using existing definitions | Standard OTLP receiver forwards to the same endpoint |
| Validation | Sentrame allowlist, versions and bounds | Collector filtering plus authoritative Sentrame validation |
| Storage | One transactional store | Same |
| Aggregation | Derive from retained event sets on read initially | Same |
| Interface | One summary and expandable evidence list | Same |
| Cost | Few processes; Sentrame owns endpoint behavior | Extra process/configuration; reusable routing and buffering |

Choose A as the draft baseline because it exercises the product with fewer operational parts. Keep the endpoint Collector-compatible. Language, database and UI framework are deliberately unselected; choose maintained components supporting this contract before implementation. No cloud provider is needed for local synthetic work.

No external queue or background aggregation service is justified by the assumed low volume. Acknowledge accepted records after storage commits. Use bounded producer retries for retryable failures, preserving IDs. If measured load or outage tolerance later requires buffering, first evaluate Collector sending queues and disk persistence; they still have capacity and loss limits. Reference: [Collector resiliency](https://opentelemetry.io/docs/collector/resiliency/).

### Event contract and trust

Draft profile: OTLP logs over HTTP/protobuf, one supported schema version, at most 100 records and 256 KiB uncompressed per request; reject compressed input initially. Use a standard exporter configured within these bounds. Bind services to loopback for the synthetic demonstration. Synthetic scope labels are not authentication and must never be presented as production access control.

Allowed envelope:

- Application scope from local configuration; schema version and fixed journey-definition version.
- Random UUID event ID and attempt ID; no reuse as person or device identity.
- OTel occurrence timestamp; server-assigned receipt timestamp.
- Event name from `attempt.started`, `action.requested`, `action.persisted`, `attempt.interrupted`, `attempt.failed`, `request.error`.
- Optional reason code from a documented finite enumeration; no free-text body, arbitrary attributes, URLs or resource enrichment.

Only the configured synthetic authoritative producer emits persistence and terminal-failure records. That convention is sufficient for fixtures, not proof of authority from arbitrary clients. A real integration needs a separately reviewed authenticated server-origin path. Persistence evidence must describe a committed operation, not a click, optimistic UI, generic HTTP status or client assertion.

Validate the complete OTLP envelope, including resources, scopes and record bodies. Reject unknown fields rather than retaining them for future analysis. For mixed valid/invalid batches, persist valid records and report rejected-record counts using OTLP partial-success semantics; never echo payloads. Retry transient request failures, not permanently rejected records. A conflicting duplicate is a rejected record; identical duplicates are successful no-ops.

Reject events older than seven days or more than five minutes in the future. Keep occurrence and receipt time separate. Timestamp order does not determine completion; elapsed duration is unavailable when inconsistent. These numeric limits are proposed review defaults.

### Outcome derivation and interface

Group by application and attempt; deduplicate by application and event ID. Apply rule version 1 to the retained evidence set:

1. Any authoritative persistence evidence: completed milestone.
2. Both terminal failure and interruption, without persistence: unknown with conflicting evidence.
3. Terminal failure only: failed.
4. Interruption only: interrupted.
5. Otherwise: unknown.

The rule is independent of arrival order. Later persistence can revise any incomplete result to completed; keep errors and interruptions visible as observations. For this slice, “interrupted” requires an explicit event; no WFS interruption control is presumed to exist. Activity after the milestone is outside the initial question, and absence of activity does not prove an exit.

Render the question, rule version, retention range and four outcome counts for attempts with a retained start. Show the started-attempt denominator explicitly. Show incomplete attempts separately, including persistence-only evidence; do not silently include them in conversion rates. Provide an evidence list with classification reason, occurrence/receipt times, conflicts and timing availability. Label outcomes revisable within retention and label the whole demo synthetic. No general dashboard builder is needed.

### Privacy, retention and deletion

Do not collect recipe content or IDs, names, household/member IDs, passphrases, invitation secrets, sensitive URLs, request bodies, IP addresses as analytics attributes, or fingerprinting data. Do not export existing application logs wholesale. Reject raw payload logging and review access logs when any real transport is introduced.

Attempt IDs permit correlation and are pseudonymous. Rotate between attempts and member changes; do not join across visits, sites or devices. This deliberately prevents longitudinal person-level analysis. Household administrator enablement alone is not assumed to establish participant consent. That consent model belongs to a separately approved integration design.

Expire events seven days after occurrence; exclude expired rows immediately at query time and physically purge daily. Derive on read so no separate aggregate retention is needed. Older partial attempts remain visibly incomplete until their own evidence expires. No backups or remote exports in the synthetic slice.

Deletion removes evidence and dependent views. Retain only a minimal application/attempt deletion marker for seven days from deletion to prevent in-window replay; discard it afterward. Disable an application scope on whole-scope deletion so ingestion cannot recreate it without explicit reconfiguration. The marker is pseudonymous metadata with a bounded purpose, not a hidden event archive. Future backups need an explicit expiry and restore-deletion policy before real data is allowed.

### Verification and harness

Reuse the existing OpenSpec installation/configuration. There is no Taskfile to reuse. After implementation approval, add only task entry points for spec validation, focused tests and the local synthetic demonstration. Do not copy WFS's harness.

Verify protocol rejection and acknowledgements, duplicate/conflict handling, ordering permutations, completion after errors, summary denominators, timing inconsistency, expiry and deletion replay. The end-to-end demonstration must pass through actual ingestion and storage, not load precomputed outcomes into a mock interface. Tests and demo are future work, not claimed evidence here.

## Risks / Trade-offs

- [Persisted action with lost telemetry] → Unknown remains possible; even a durable ingest cannot close the application's commit-to-emission gap. Evaluate an outbox only if a future real integration requires stronger coverage.
- [Selection bias] → Results describe observed opt-in attempts, not all users or value delivered to a population.
- [Seven-day window and rotating IDs] → Limited trend and return-visit analysis; intentional privacy trade-off for this slice.
- [Fixed milestone] → Demonstrates mechanics but cannot establish satisfaction or that Discovery improves household outcomes.
- [OTLP ecosystem varies by language] → Confirm exporter/profile compatibility when choosing the implementation stack; do not claim broad OTLP conformance from one synthetic producer.

## Migration Plan

No deployed system or data migration exists. This change remains planning-only until a subsequent implementation request. Future local implementation can be rolled back by stopping the demo and deleting synthetic storage. Deployment, publication and live WFS instrumentation remain separately authorized steps.

## Open Questions

- Which language, transactional database and minimal UI framework best fit the local implementation? The choice must preserve the defined contract and single-service shape.
- What exact labels and layout make outcome evidence easiest to read? Keep the fixed summary and evidence view; no additional analytics capability is implied.

Product scope, retention and ingestion proposals are documented above for iterative review rather than silently treated as final decisions.
