# Tasks

All tasks are future implementation work, intentionally unchecked. Current research is source/history-backed; a fresh build succeeded but isolated UI observation was blocked before browser launch. Earlier native acceptance remains pending independently.

## 1. Exact-ID investigation

- [ ] 1.1 Add compact typed Find ID entry and synchronize existing exact scope filtering; verify session/operation, journey and trace actions, invalid/control/over-limit input, case-sensitive equality and no silent filter broadening, and document ID meanings in EN/FR help.
- [ ] 1.2 Implement protected exact journey/trace resolution independent of capped dropdowns; verify equal IDs across apps/generations, multiple tabs/journeys, viewer scope and revoked access without enumerating inaccessible targets.
- [ ] 1.3 Add type-labeled Copy/Open/Filter actions beside supported IDs and evidence-backed no-match recovery; verify invalid, retained-but-unmodeled, unknown/expired and capped states, with concise contextual guidance rather than guessed diagnoses.
- [ ] 1.4 Unify accepted query versus draft state and same-tab reload/back/forward restoration without new raw-ID URL exposure; verify app/sign-out/revoke/schema invalidation and stale response discard, and document the minimal local persistence policy.

## 2. Actionable matching and timeline

- [ ] 2.1 Carry graph criterion, exact journey and matched occurrences into matching/timeline/story entry; verify stage, transition and prefix examples open relevant occurrences without auto-selecting a representative journey, and document each view's question/next action.
- [ ] 2.2 Add Follow scope and previous/next match with repeat-group preservation; verify error → recovery, missing/unplaced diagnostics, grouped repeats and configured versus goal-free outcomes without invented success semantics.
- [ ] 2.3 Connect occurrence → request/diagnostic → back and safe-copy handoff; verify valid focus/scroll/selection restoration, local failures/expiry and inaccessible references using synthetic interaction tests.
- [ ] 2.4 Reconcile path-table inspect actions and saved/reference criterion/ID preferences; verify exact path membership, old reference defaults, authorization and retention invalidation, and update reference/help documentation with the resulting behavior.

## 3. Directed occurrence hierarchy

- [ ] 3.1 Project bounded uninterrupted prefix nodes with separate semantic-stage identity and edge provenance, reusing/reconciling existing prefix evidence; verify A-B-A-C, A-B-D/A-C-D and A-B versus X-B-C do not fabricate paths/rejoins, and document ordinal versus causal/time semantics.
- [ ] 3.2 Add staged expansion, cutoff/end/goal/gap markers, stage-wide cross-reference and exact-ID rare/deep-path access; verify 32-step/200-node/initial-branch limits do not alter authoritative counts or hide omitted support, with explicit EN/FR explanations.
- [ ] 3.3 Preserve semantic selection and occurrence criterion through Force/Sankey switches, manual vertical arrangement and structured keyboard/list alternative; verify direction, repeat identity, counting units and concurrent API uncertainty remain explicit.

## 4. Integrated task and accessibility acceptance

- [ ] 4.1 In a fresh identified isolated synthetic build, verify known ID → graph → relevant example → timeline occurrence → trace → return → reference/reload with cross-app, expiry, source-cap and failure fixtures; record actual outcomes and build identity without using stale19010.
- [ ] 4.2 Check desktop/390px EN/FR, compact canvas-start budget, 44px targets, keyboard disclosure/dismissal/focus and no permanent added header row; update guidance and record any actual regression separately from visual preference.
- [ ] 4.3 Run native200% zoom, EN/FR screen-reader and physical-touch tasks where available, with an uncoached investigation walkthrough; report passed/failed/not-run separately and do not close earlier pending checks based on mocked or build-only evidence.

## Workflow follow-up

- Obtain a separate implementation instruction after design review. Do not stage, commit, push, change branch, implement, sync canonical specs or archive during this planning request.
