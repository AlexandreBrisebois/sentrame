# Header tools review — 2026-10-31

## Decision for Alex

Recommend **Tools ▾ / Outils ▾**, a grouped anchored disclosure in the existing graph toolbar. Keep the app selector dedicated to app switching; keep seven graph views and Journeys/Traces directly accessible. One reasonable alternative is an overlaid drawer: easier to grow and read long labels, but covers more graph and adds modal complexity. Details and full destination mapping: [design.md](design.md). No implementation is authorized or performed.

## Current source and observed UI

Entry and exit: branch `main`; HEAD `a2d2ee0cf7c0759fcc4ac7789a566f03ddd8dc05` (`updated docs`). Initially no dirty or staged files. Exit has only this new untracked proposal directory; tracked diff and staged diff are empty. Earlier journey and compact-header archives remain unchanged, including outstanding native zoom, screen-reader and uncoached walkthrough acceptance.

Relevant source:

- `web/src/Workspace.svelte:325–332`: masthead app selector, separate Journeys/Traces, Account/App settings/languages/sign out.
- `web/src/Workspace.svelte:354–360`: seven graph views, Historical/Live and current hamburger; graph highlight plus nine actions with no visual grouping. Configuration, save and invitation actions are owner-only.
- `web/src/Workspace.svelte:165–168,322,484`: menu actions close menu before dialogs; Escape only clears menu state; modal fallback is the menu trigger.
- `web/src/style.css:19,216,225`: menu layout, trigger sizing and bounded height.
- `web/src/Disclosure.svelte` and `Modal.svelte`: neighboring native disclosure and modal focus/dismissal models.
- `openspec/specs/journey-exploration/spec.md`, `app-management/spec.md`, `viewer-access/spec.md`: compact graph, reversible/accessibility and scoped navigation constraints.
- `docs/design/visual-theme/README.md`: retained forest-green, warm-paper visual language. No local designer skill or applicable AGENTS.md was found.

The current menu is below the app selector at the right of the graph toolbar, not immediately to its left. Findings concern the actual current menu; do not reintroduce older placement.

Fresh-browser evidence, same results across EN/FR at1280×900 and390×844:

| Check | Observation |
| --- | --- |
| Layout | Canvas starts165px desktop and200.1875px phone; panel does not shift it. No page horizontal overflow. |
| Trigger/targets | Trigger44×44px; desktop single-line actions42px high; phone actions44px or taller. |
| Keyboard | Enter opens; Tab enters native Highlight select; Escape closes but leaves BODY focused. |
| Outside click | Clicking brand leaves menu open. |
| Mutual dismissal | Opening Context leaves both Context and menu open. |
| Dialog return | References opens its existing dialog; Escape closes it and returns focus to `workspace-menu-button`. |
| Touch | Emulated phone tap opens/closes trigger. This is not physical-device acceptance. |
| State/labels | Localized accessible trigger name and expanded state exist. No controlled-panel relationship, title or Close. Highlight has current value; discovered-event state is conveyed only by action wording. No menu destination current mark; main views separately retain selected states. |

The list mixes graph state, saved-view actions, diagnostics and administration, making scanning harder. The hamburger glyph renders small and offers no visible tool meaning. No clipping was observed: French wrapping makes the current panel506.69px high at390px, ending608.69px within844px. The proposal addresses growth/zoom proactively without claiming a present clipping defect.

## Reproducibility and boundaries

Scratch evidence (ignored by Git): `.task/header-menu-review/` contains `review.cjs`, `result.json` with source/asset SHA256 identity and browser version, `run.log`, four `*-menu.png` screenshots and fresh `dist/`. The harness reuses only synthetic fixture construction from prior local review code and serves a fresh build; it does not contact a real receiver. Build command: `npm run build -- --outDir ../.task/header-menu-review/dist` from `web/`; Svelte check0 errors/0 warnings and Vite build passed. Assets: `index-C67bGDGA.css`, `index-D83fO1bZ.js`. All browser requests outside the owned loopback origin were blocked. Profile isolated, browser and owned ephemeral server stopped in cleanup. No persistent credential changes, containers, shared resources or user browser tabs were modified.

Initial sandbox loopback bind was denied with `listen EPERM`; the scoped escalated run was approved and completed. No unresolved denied/pending/uncertain action remains. Existing localhost19010 was not used.

Validation: `openspec validate simplify-workspace-tools --strict` passed; `openspec status --change simplify-workspace-tools` reports4/4 planning artifacts complete. All implementation tasks remain unchecked. This is a current-UI review, not validation of an implemented replacement. Native200% zoom, physical phone touch, native EN/FR screen-reader announcements and real-user usability remain untested; tasks preserve those obligations. Viewer action inventory was source-reviewed, not run with a viewer session. No full application test suite or backend qualification is implied.

## Handoff

Read proposal.md and design.md with Alex; obtain a new implementation instruction before applying. Do not stage, commit, push, switch branches, sync canonical specs or archive this proposal during review. Runtime, API, data and older acceptance records are outside this change. The focused tasks and delta are ready for review.
