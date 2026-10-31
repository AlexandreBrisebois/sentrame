# Compact graph implementation — 2026-10-31

This is the earlier inline-disclosure checkpoint. The approved tighter layout and current acceptance are in [GRAPH-FIRST-HANDOFF.md](GRAPH-FIRST-HANDOFF.md).

**Software implemented and verified; 8/9 tasks complete. Native acceptance task 5.2 remains pending.** Alex explicitly approved this reviewed implementation. No sync/archive was performed: the earlier archive exception applied to another change and does not waive this change’s unperformed checks.

## Result

- Removed the duplicate task dropdown. All seven view buttons and saved-reference destinations remain, with stable active-view focus fallbacks (`web/src/Workspace.svelte:347`). Trace/example return tests exposed and fixed the dialog’s connected-body fallback, with preventScroll on focus restoration (`web/src/Modal.svelte`).
- Added reusable native disclosures with keyboard/touch activation, Escape within the disclosure, explicit Close, focus return and bounded inline content (`web/src/Disclosure.svelte`, `InvestigationHelp.svelte`). Full question/unit/selection and applicable L/B/error explanations are available on demand.
- Made Coverage compact while preserving source/population/graph/record-list details, unknown collection completeness and authoritative denominator explanation. Actual source caps, graph/list omissions, partial/censored evidence and unknown metadata remain visible; graph-local display-cap warning remains in its existing footer (`CoverageDisclosure.svelte`). Comparison/stale failure restrictions remain outside help.
- Reserved stable refresh-status space and stopped explicit same-context Refresh from blanking the successful snapshot. Historical validation still polls at its existing cadence; changed revisions wait for explicit reconciliation unless retention/deletion invalidates evidence. Surviving selection, open examples and presentation state are retained; unavailable raw detail is removed and a disconnected focus target gets a valid successor (`Workspace.svelte:234`, `:372`). No graph physics/camera algorithm changes.
- Repeat checkboxes use intrinsic width and nonshrinking space with complete wrapping label spans in both timeline and sequence. Updated English/French disclosure labels and the user guide.

## Final verification

- `task check`: PASS (Go packages, 144 frontend tests at that checkpoint, Svelte zero errors/warnings, ten strict canonical specs, whitespace). After adding the independent-warning test, final `npm test`: **145/145 PASS**. Production `npm run build`: PASS with zero Svelte diagnostics. Strict selected-change validation and final whitespace check: PASS.
- **19/19 final-code isolated browser checks**, zero page errors, plus **1/1 additional saved-reference check** across all seven destinations. Synthetic in-memory API, ephemeral loopback listener, Chrome156 headless and separate temporary profiles; no real receiver/telemetry. First-load loading state was observed in the additional run.
- Delayed background reads in Hierarchy/Force/Sankey: **0px canvas shift**, stable SVG identity, selected node, focus, scroll, manually moved node position and panned camera. The prior investigation measured38.25px shift. Changed revision/explicit Refresh preserved a selected edge and mounted camera. Failed reads/retry, unrelated deletion with open evidence, unavailable chosen detail, removed focused graph target, trace return, removed invoker, superseded filters and access clearing passed.
- Disclosures passed Enter/Space, Tab, Escape, explicit Close, focus return and emulated-touch activation/close. Coverage stays open across refresh; local100-node cap with10omitted remains visible. Independent SSR warning cases cover source, partial starts, censoring, node/edge/list omissions, baseline and unavailable coverage; denominators are unchanged.
- EN/FR×1280/390, timeline and sequence: checkboxes **18px**, label targets **44px** high, full text readable, no page horizontal overflow. Desktop graph and EN/FR screenshots were visually inspected. These checks do not establish native zoom or physical assistive-technology conformance.

Evidence: `.task/compact-graph-review/{result.json,additional-result.json,verify.cjs,run.log,additional-run.log,final-identity.json,tasks-progress.json}` and final check/frontend/build logs. Screenshots include `graph-desktop.png`, `en-1280.png`, `fr-390.png`. Earlier harness-driver failures and the discovered dialog fallback defect are preserved in first/second/third-run records. Keyboard-driver failures were corrected with complete Enter/Space events; SVG activation and persistent CDP lexical variables were harness issues. Final assertions all pass.

Final built assets: JS `index-ByiRgp_s.js` SHA256 `57cb56f1d7ba9bca6468ca81231b6ceaa91af8639b1fcf2a264cac7baeba7fed`; CSS `index-Cd2thu0E.css` SHA256 `7cba5a1ad2a81a74a82d6e387975c31ed391f0e151e677e0590cee0eee0cd6df`. Browser records match these hashes; final-identity records all changed source hashes.

## Remaining acceptance and ownership

**Pending, not denied:** task5.2 requires actual native200% browser zoom and physical EN/FR screen-reader checks of the changed controls/status. No interactive browser slot was taken; coordinate with the separate WFS M6 owner before using shared native UI, then record actual observations. Headless widths/touch/keyboard checks do not substitute. No approval-review rejection or uncertain mutation remains. The apply skill requires honest task tracking: [openspec-apply-change](../../../../.agents/skills/openspec-apply-change/SKILL.md), “Only mark a task `- [x]` when its specified behavior is fully implemented”. No waiver is inferred from the old archive.

Branch/HEAD unchanged: `main`, `44e4f1430ec3f923a5e5934a1dda0b624cf80bad`. Entry index empty; final index empty. No staging/unstaging, commit/push, branch switch, canonical spec edits or old archive edits. A concurrent unrelated `.gitignore` simplification appeared after entry; it was not authored or reverted here and must be preserved. Changed owned source: Workspace, CoverageDisclosure, InvestigationView, Modal, styles/i18n, new Disclosure/InvestigationHelp and viewportRender tests, plus user-guide wording and this change’s artifacts.

All owned temporary Chrome/node verification processes and ephemeral servers are stopped; a scoped final process check found none. Temporary synthetic profiles/evidence remain for inspection. No WFS files/ports/containers/native session, production, credentials or shared receiver were used or changed. Existing19010 preview was not redeployed and still serves its earlier build. New assets were inspected only in the owned isolated browser. No local service remains owned by this implementation task. Parent/Alex own source review and remaining acceptance from this handoff.

Suggested commit subject: `Improve graph viewport space and historical refresh stability`
