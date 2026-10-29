# Greenfield cutover and OpenSpec reconciliation

## Owner clarification

The owner clarified during planning: “I have no users, this is greenfield no need for backwards compatibility.” This supersedes the initial request's obligation to preserve deployed integrations. No existing-user migration, legacy replay comparison, old schema decoding, API alias, dual write, feature-flag compatibility mode or automatic old-data conversion is required.

Optional `stage.name`, paired stage boundaries and presence remain intentional forward capabilities, not compatibility shims. Prospective definitions and stable saved references remain necessary for future captures after launch; they do not create a promise to load pre-change local data.

## Cutover plan for a later authorized implementation

1. Implement the new projection and API contract against an isolated fresh database with synthetic data. Store explicit schema/projection/computation versions from the start.
2. Preserve unrelated working files. Do not operate on `.sentrame` during this planning pass or point qualification at an existing review database.
3. Prove an old unsupported schema fails with clear fresh-directory guidance. No startup auto-reset and no silent interpretation using the new model.
4. When the owner adopts the implementation, explicitly stop the old local instance and choose a fresh data directory, or deliberately reset only the selected synthetic dataset using reviewed local tooling. Do not execute reset merely because this plan describes it.
5. Recreate test mappings/reference examples under the new contract; regenerate fixtures. No historic local result comparison or legacy-data import is an acceptance gate.
6. If rollback is needed, run an old binary with its matching old test directory or create another clean directory. No database downgrade guarantee, deployment action or publication follows passing tests.

## Canonical baseline and concurrent archive

At planning start, canonical specs were empty and `synthetic-journey-outcomes` had 24/30 tasks complete. During final checks, concurrent work archived that package to `openspec/changes/archive/2026-10-29-synthetic-journey-outcomes/` and synchronized eight canonical capabilities. All requirement blocks were compared with the inspected baseline and match. This proposal now adds unique requirements to five canonical owners, modifies one diagnostic requirement to remove the obsolete legacy-data gate, and introduces `trace-exploration`. It does not edit the archive or canonical files.

| Baseline owner | Existing obligation retained | Increment / clarification |
| --- | --- | --- |
| `telemetry-intake` | Standard receipt independent of business context; auth, atomicity, supported identities, metric receipt | Structural projection, contextual discovery, no new transport cell assumed |
| `journey-modeling` | First mapped / active goal / 30-minute inactivity; explicit scopes; healthy loops; source authority; versioned definitions | Standard span stages, hints and explicit selectors; native span timing supplements supported pairs |
| `journey-exploration` | Captured boundaries, graph layouts/counts, gaps, references, timing/table and accessible controls | Typed trace view distinct from journey graph; reference mode/policy context |
| `journey-diagnostics` | App-wide errors, scrubbed text and honest occurrence counts | Trace links retain existing context intent; a narrow MODIFIED delta replaces only legacy-projection replay with same-model replay |
| `data-lifecycle` | 30-day detail, 12-month summaries, bounds, deletion/privacy | Structural lifecycle and new-model replay; legacy diagnostic projection replay is not a greenfield requirement |
| `app-integration-guidance` | Discover source and existing SDK; reviewed minimum event plan; no proprietary transport | Progressive recipes/docs and qualification separation |
| `live-journeys`, `viewer-access` | Snapshot revision, bounded SSE, read-only expiry/revocation | Reuse existing boundary for trace reads; no new grants, baseline semantics or live push protocol |

The baseline's Stage occurrence timing requirement already permits standard OTel start/end evidence. This change makes native spans explicit without deleting the optional pair behavior. Baseline examples do not force paired logs for every operation. Its broad compatibility target includes gRPC; this change does not implement or advertise it. The baseline's diagnostic “Projection upgrade and replay” scenario was written for earlier retained test data; the diagnostic delta explicitly replaces it with same-model replay per owner clarification. Runtime replay of records accepted by the new model remains required.

## Apply and future synchronization order

- Apply against the current implementation and canonical baseline. Do not rerun archived completed work or mark deferred resource/roadmap items complete.
- Read this increment with canonical specs and later decision-register scope extensions. Update only this change's task progress during apply.
- A separately authorized future sync/archive adds these unique requirements to their existing owners and applies the complete diagnostic replacement block. Retain all unrelated canonical requirements; no duplicate capability names.
- No canonical sync/archive or archived-baseline edits happen in this planning pass. The archived design's historical no-application and duration-sized-node claims do not override source or later decisions.

Roadmap B stays deferred. A single-app trace list and contextual selector fields are part of this change, not a multi-app registration UI, arbitrary rule engine, richer release comparison or narrow invitation project.
