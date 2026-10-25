# Proposal

## Why

Application builders need to know whether an observed journey reached a useful outcome, where explicit friction occurred, and where evidence is insufficient. Sentrame connects observable events to declared product questions without treating time spent as value or telemetry as an explanation of personal intent.

## What Changes

- Propose a synthetic-only, local end-to-end slice: accept a small set of structured OpenTelemetry events through OTLP, correlate an attempt, derive an evidence-based outcome, and present a summary with supporting evidence.
- Use “Did a Discovery attempt reach its first persisted vote?” as the initial example question. This is a proposed review baseline, not a claim about lifetime activation, satisfaction, or implemented WFS instrumentation.
- Distinguish completed milestone, explicit interruption, terminal failure and unknown evidence; preserve completion after a quick departure.
- Bound event fields, correlation, retention, deletion and delivery semantics before any real integration.
- Keep the artifacts open to iterative design. Implementation and any live integration require subsequent authorization.

## Capabilities

### New Capabilities

- `journey-evidence`: Accept, validate, correlate, retain and delete bounded product event evidence.
- `journey-outcomes`: Derive and explain outcomes for a declared journey milestone, including uncertainty and collection limitations.

### Modified Capabilities

None; no existing capability specifications are present.

## Impact

- Planning artifacts only in this change; no application code, live instrumentation, deployment, credentials or real telemetry.
- Existing repository contains README and license plus newly present OpenSpec configuration and generated skills. Preserve that local setup; no Taskfile or development harness exists.
- Future implementation needs a small ingestion service, transactional storage, deterministic derivation and a minimal insights view. Language, database and deployment provider remain unselected.
- Intended first users are application builders and small product teams. WFS is a source-checked candidate integration, not the product boundary or an architecture template.
- Out of scope: invitation funnel, cross-visit identity, session replay, fingerprinting, cross-site tracking, AI, arbitrary journey editors, technical monitoring dashboards and a distributed telemetry platform.
