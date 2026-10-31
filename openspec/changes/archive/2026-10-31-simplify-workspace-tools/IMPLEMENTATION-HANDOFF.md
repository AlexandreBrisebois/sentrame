# Compact Tools implementation — 2026-10-31

## Outcome and ownership

Alex approved the recommended compact design before implementation. Replaced the hamburger in its existing graph-toolbar slot with **Tools ▾ / Outils ▾** and three groups: View, Investigate and Manage this app. Existing destinations, role restrictions and separate app selection remain intact. Saved views names the existing reference dialog; discovered-event visibility is an explicit checked checkbox. Selected graph views and highlight values retain their actual state; action launchers do not pretend to be selected pages.

Branch remains `main`, HEAD `a2d2ee0cf7c0759fcc4ac7789a566f03ddd8dc05`. Entry had only this proposal directory untracked; index was clean. No staging, commit, push, branch switch, canonical-spec sync or archive performed. Earlier October31 archives and their outstanding acceptance records remain unchanged. This supersedes the review-only outcome in REVIEW-HANDOFF.md; that document remains baseline evidence.

## Changed files

- `web/src/WorkspaceTools.svelte` (new): localized grouped controls, selected-app title, existing action callbacks, owner/viewer inventory, bounded wrapping presentation and44px targets.
- `web/src/Disclosure.svelte`: optional compact heading/trigger identity and focus-exit behavior; native `popovertarget` association makes repeated toggles and mutual dismissal reliable. Tools gets explicit Escape/Close return and closes on focus exit without stealing focus. Existing disclosures retain their default width/content and gain the native trigger association.
- `web/src/Workspace.svelte`: reuse Tools, remove legacy menu state, preserve dialog return to the Tools trigger when an action was in the now-hidden panel, reset panel on app/generation/role changes. Existing footer/inspector invokers retain ordinary modal return behavior.
- `web/src/style.css`: remove obsolete hamburger/popover/lens-only rules; no new header row.
- `web/src/i18n.ts`: Tools, scope title, three groups and Saved views in French.
- `web/src/workspaceToolsRender.test.ts` (new, isolated temporary Vite cache cleaned after use): owner/viewer destination inventory, explicit selected graph state, unsupported/empty highlighting, French task parity and accessible disclosure relationships.
- `docs/USER-GUIDE.md`, `docs/LOCAL-REVIEW.md`: current Tools routes and interactions; Account → App settings remains distinct.
- This active OpenSpec change: task progress and this handoff. Canonical specs are unchanged.

## Verification

Final code was built with `npm --prefix web run build -- --outDir ../.task/tools-implementation/dist --emptyOutDir`. Svelte check:0 errors/0 warnings; Vite build passed. Final assets: `index-C_flMSuQ.css`, `index-CfC0RfpK.js`. No stale localhost19010 assets were used or replaced.

`node .task/tools-implementation/verify.cjs` passed9 browser scenarios against those assets in isolated Chrome156.0.8078.74:

1. EN1280×900 geometry, keyboard, state, dismissal and panel switching.
2. FR1280×900 equivalent behavior.
3. EN390×844 equivalent behavior and touch emulation.
4. FR390×844 equivalent behavior and touch emulation.
5. All8 owner action buttons reach their existing dialogs, close Tools and return focus to its trigger.
6. Long French app scope at390×480 wraps; bounded scrolling reaches every action and Close.
7. Viewer inspection destinations remain; owner writes absent;401 revocation removes open Tools and clears access.
8. Empty and failed synthetic reads retain reachable tools and evidence explanation.
9. App switching during a delayed old-app read removes the prior panel and retains new-app scope after that response.

The four layout cases also verify Enter/Space, Tab into controls, Shift+Tab exit, Tab beyond the last action, Escape from inside, Close, repeated pointer/touch toggles, outside pointer/touch dismissal without stealing clicked-button focus, persisted checkbox/highlight state, all seven direct views revealed on keyboard focus and mutual replacement with Account/Context/Help/Coverage. Covered overlay targets were activated by keyboard; no forced pointer clicks through an overlay were used.

Canvas starts remain165px desktop and200.1875px phone EN/FR. Opening Tools causes zero displacement. Panel width320px remains inside the viewport; all interactive row targets meet44×44px and no page horizontal overflow occurs. Screenshots were visually inspected; French labels and the full app scope wrap, group headings and Close remain readable. Nine scenarios emitted no page errors and made zero API mutations. These are synthetic UI checks, not backend/production qualification.

Scratch evidence: `.task/tools-implementation/{browser-result.json,browser.log,verify.cjs,build.log,final-aggregate-isolated.log,final-frontend-serial.log,final-aggregate.log,dist/,en-1280.png,fr-1280.png,en-390.png,fr-390.png,fr-short-long-app.png}`. Browser results include final source/asset SHA256 identity; hashes were compared with current source after execution. Earlier failed focus/toggle checks drove the related-target focus-exit and native-trigger fixes; only final outcomes count as acceptance.

Final `task check` passed against final code (`.task/tools-implementation/final-aggregate-isolated.log`): all Go packages including `internal/app`,148/148 frontend tests, Svelte check0 errors/0 warnings,10/10 strict canonical specs, and whitespace checks. Strict active-change validation also passes. An earlier aggregate passed; a subsequent concurrent run stalled in the existing viewport SSR worker and surfaced `ENOTEMPTY` in shared Vite deps after that owned worker was terminated. All148 tests then passed serially. The new render test now uses and cleans its own temporary Vite cache; the final full aggregate rerun passed with that isolation. Existing render tests can still emit non-failing HMR-port warnings; no unrelated test harness was changed. Native200% browser zoom, native EN/FR screen-reader announcements and physical phone touch were not performed; task3.2 remains unchecked. Headless widths, browser key events and touch emulation do not close those obligations or any prior archive's pending checks. Keep this change active; do not archive without completing acceptance or Alex explicitly authorizing that exception.

## Cleanup and safety

Owned preview used an ephemeral127.0.0.1 listener with synthetic in-memory API responses. Browser routing blocked every other origin. The isolated browser/profile and owned listener were closed in `finally`; existing listener inventory is unchanged from entry. No shared receiver, production telemetry, real credentials, persistent credentials, containers or user browser tabs were touched. No service or preview is left running by this task; scratch build/evidence files remain for review. Initial sandbox loopback restriction was resolved through approved scoped escalation; no unresolved denied, pending or uncertain tool action remains.

## Suggested commit message

`Improve compact workspace Tools navigation and dismissal`

## Archive decision

On 2026-10-31 Alex explicitly authorized archiving with task 3.2 incomplete. Native 200% zoom, EN/FR screen-reader operation and physical phone touch acceptance remain pending.
