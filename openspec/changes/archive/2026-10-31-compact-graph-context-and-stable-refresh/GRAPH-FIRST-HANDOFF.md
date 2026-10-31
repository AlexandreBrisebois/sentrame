# Graph-first refinement — final handoff 2026-10-31

Later closeout: Alex explicitly requested sync/archive with task5.2 still pending. See [ARCHIVE-NOTE.md](ARCHIVE-NOTE.md); this implementation checkpoint does not claim native acceptance.

**Approved refinement implemented and verified; 12/13 tasks complete.** Task5.2 (actual native200% zoom and physical EN/FR screen reader) remains pending. The prior inline-disclosure implementation/review records are historical; this is the current handoff. Existing change proposal/design/delta/tasks now describe the approved flyouts and compact header. No canonical sync/archive was performed; the old archive exception does not waive this change’s acceptance.

## Before / after on current code

Same isolated synthetic navigation fixture and CSS viewports, closed flyouts, no material warnings or inspector:

| Viewport / mode | Prior graph top | Final graph top | Final graph height |
|---|---:|---:|---:|
| 1440×900 EN Historical | 467.73px | **165px** | **686px** |
| 390×844 EN Historical | 841.59px | **200.19px** | **576.88px** |
| 390×844 FR Historical | 899.97px | **200.19px** | **576.88px** |
| 1440×900 EN settled Live+baseline | 533.13px | **165px** | **686px** |

1280×900 EN and FR also begin at165px. Document scrollWidth equals viewport width in every measured case. Phone graphs now appear in the initial viewport. Material warnings and active narrowing intentionally may add height; these results are not a promise to conceal them to meet a budget. Earlier measurements592px for Live were a loading state and were separated from the settled baseline in the review.

## What changed

- Compact masthead keeps app selection and Journeys/Traces direct; Account **⋯** contains settings, language and sign-out. Removed the repeated workspace heading and separate mode-navigation strip.
- First graph row retains all seven view buttons with selected state and stable IDs; on phones it is a horizontal strip, with keyboard focus revealing every destination. Time & filters remains direct; References and highlights are in the existing workspace menu.
- Second row retains short bounds, scope, current/L/B retained counts and active narrowing/clear chips. **Context**, **Help** and **Coverage** are anchored labeled native nonmodal flyouts, with bounded scrolling and44px controls. Exact windows, definitions, scope details, generic explanatory context and Live presentation controls are on demand. No denominator or query semantics changed.
- Removed standalone empty update/loading rows. Refresh/pause occupies an existing fixed-size action slot. Routine loading does not move the graph; meaningful updates/stale/error states remain distinguishable and announced. Material coverage/read/comparison warnings stay visible with flyouts closed.
- Preserved the previous repeat-label and evidence/camera fixes. The relocated References menu exposed a real dialog-focus gap; generic tool dialogs now fall back to persistent `workspace-menu-button` instead of the page body. Account-to-settings and graph/trace/example returns are verified.

Main code: `web/src/Workspace.svelte`, `Disclosure.svelte`, `CoverageDisclosure.svelte`, `InvestigationHelp.svelte`, shared styles/i18n and user guide. Prior InvestigationView/Modal fixes remain intact. No graph-layout/physics, backend, ingestion, persistence, authentication or lifecycle rule changed.

## Exact final verification

- Final `task check`: **PASS**, Go packages (internal/app98.45s), **145 frontend tests /0fail**, Svelte **0errors/0warnings**, **10strict canonical specs**, whitespace. Production build: PASS. Selected-change strict validation: PASS.
- **20/20 final browser checks** plus **3/3 extra checks**, **0page errors**. Chrome156.0.8078.74, dedicated temporary profiles and ephemeral loopback in-memory mock API. No real telemetry/receiver/protocol qualification.
- Native Enter/Space and Tab, Escape, explicit Close, native outside mouse dismissal, emulated touch, replacement by another flyout, viewport-clamped phone flyout bounds, focus return, account EN/FR/settings, menu-to-modal return, seven direct and seven saved-reference destinations passed. Open flyouts do not alter canvas top/height.
- Delayed700ms unchanged reads in Hierarchy/Force/Sankey: **0px jitter**, same SVG, selected node, focus, scroll, manual node position and camera. Changed revision with selected edge, explicit Refresh, stale failure/retry, superseded filters, unrelated deletion with an open example, unavailable chosen detail, focused-node deletion, trace return and access clearing passed.
- Coverage remains open through refresh. Source cap/partial evidence/list omissions remain visibly scoped while closed; details retain unknown collection completeness and unchanged ordered-journey denominators. Local graph100-node cap with10omitted remains visible. Phone Live/Baseline identities, pause/resume and active-version clearing passed.
- Repeat controls retain18px inputs and44px full-label targets; new flyout controls meet44px. EN/FR desktop/390 checks report no page overflow. Final desktop, phone EN/FR, warning/flyout and phone Live pixels were inspected. Native200% zoom, physical screen reader and real touch hardware were not performed.

Evidence root: `.task/graph-first-refinement/`. `result.json`, `extra-result.json`, `measurements.json`, `measure.log`, `verify.cjs`, `measure.cjs`, `final-check.log`, `build.log`, `final-source.json`, `final-state.json`, `task-progress.json` and screenshots. Browser server code reads an immutable `final-dist/` copy; source hashes were checked unchanged during verification. Earlier driver expectation corrections and the actual removed-menu focus defect are preserved in first-pass/before-menu-focus-fix evidence.

Final built/verified JS `index-D83fO1bZ.js`: SHA256 `90a04d5dfae948a7338b37a0f2208bea8b5dde75f4e424cdb6ff9ee3cdf88b86`. CSS `index-C67bGDGA.css`: SHA256 `11cc331e7a6775f572c983eb0a7736449ee40e8e8804d338f0512d089644780f`. Both browser result manifests match these final artifacts.

## Acceptance, preservation and release

Pending action: coordinate a native browser slot with WFS M7’s owner, then observe actual200% zoom and physical EN/FR screen-reader behavior of the changed controls/status (task5.2). **Pending, not denied; not attempted.** No native/shared UI slot was claimed. No waiver or physical-accessibility conformance claim is inferred from headless evidence. The old archive and its three user-owned manual checks remain untouched.

Same `main` / HEAD `44e4f1430ec3f923a5e5934a1dda0b624cf80bad`; index empty. All unrelated `.gitignore` content is hash-verified unchanged from entry. Existing dirty work is preserved; no staging/unstaging, commits, pushes or branch switches. Active planning files and approved frontend/docs/tests are the only owned source edits. No automatic-review denial or uncertain mutation remains.

Existing preview19010 was not redeployed and **does not show this result**; it remains on its older stage5 assets. All screenshots/measurements here use the isolated final build. No production/shared receiver, credential, WFS file/port/container or other task resource changed. All owned temporary Chrome/node processes and ephemeral listeners stopped; scoped process lookup returned no matches. Synthetic profiles/build copies/evidence remain for review. No running service remains owned. Parent/Alex own source review, preview/deployment decisions and pending native acceptance from this handoff.

Suggested commit subject: `Give journey graphs a compact header and accessible context flyouts`
