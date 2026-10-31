# Design

## Context

Review baseline: clean `main`, HEAD `a2d2ee0cf7c0759fcc4ac7789a566f03ddd8dc05`, 2026-10-31. See proposal.md for motivation and REVIEW-HANDOFF.md for evidence. No AGENTS.md exists in this checkout or its ancestor chain; repository `.agents/skills` contains OpenSpec workflows but no designer skill. `docs/design/visual-theme/README.md` supplies the canonical visual direction.

The reported menu “left of the app selector” is not in that position in current source or fresh screenshots. `Workspace.svelte:325–360` puts app selection in the forest-green masthead and the hamburger at the right of the graph-view toolbar below it. Do not restore old positioning based on older screenshots. Its purpose is a local workspace tool launcher, not app switching or primary page navigation. Journeys/Traces and seven graph views already have selected-state navigation.

Current menu: Highlight selector; References; Show/Hide discovered events; Unmapped events; Mappings & goal (owner); Save current view (owner); Temporary viewers (owner); Errors & exceptions; Data lifecycle; About this evidence. Actions open existing dialogs except Highlight and discovered-event visibility. Account currently holds App settings, EN/FR and Sign out. Footer separately exposes Errors, Table, Journeys & evidence.

Observed problems: a small hamburger glyph communicates neither purpose nor scope; nine ungrouped actions follow an unrelated selector; “References” and “Save current view” are separated; no heading or Close; the trigger exposes expanded state but no controlled-panel relationship; Escape removes focused content without recovering focus; outside click does not dismiss; another flyout can coexist. The menu itself has no active-destination indicator, although modal titles and graph-view selected states exist elsewhere. Desktop actions measure 42px; phone actions are at least44px. Phone French wrapping lengthens the panel; the current panel fits the tested viewport and produces no page horizontal overflow, so do not claim clipping was observed.

## Goals / Non-Goals

**Goals:** Make tools recognizable and scannable; make scope and state clear; keep the current compact rows and graph viewport; provide consistent keyboard/touch behavior.

**Non-Goals:** New navigation routes, persistent sidebars, extra masthead rows, changed evidence semantics, API changes, credential operations, implementation, or acceptance closure of archived work.

## Decisions

### Recommendation: labeled grouped Tools disclosure

Replace the hamburger in its existing toolbar slot with **Tools ▾ / Outils ▾**. Keep Historical/Live immediately before it and the horizontally scrollable seven-view strip before that. Use a nonmodal anchored panel with a title **Tools for [selected app] / Outils pour [application]**, Close/Fermer, quiet group headings and subtle separators. Retain the forest-green masthead, paper/surface colors, sans labels and amber keyboard focus.

| Group EN / FR | Controls in order | Behavior |
| --- | --- | --- |
| View / Affichage | Highlight / Mettre en évidence; Show discovered events / Afficher les événements découverts | Existing native select with current option; labeled checkbox/switch with checked state. Hide unsupported highlights under existing rules. Changes stay open and apply only to current graph context. |
| Investigate / Explorer | Saved views / Vues enregistrées (existing References); Save current view… / Enregistrer la vue actuelle… (owner); Errors & exceptions; About this evidence | Existing reference/save/diagnostic/evidence dialogs. Adjacent saved-view actions explain their relationship. Selecting closes Tools before the dialog opens. |
| Manage this app / Gérer cette application | Unmapped events; Mappings & goal (owner); Temporary viewers (owner); Data lifecycle | Existing destinations and exact authorization. Data lifecycle stays available to viewers with its existing read-only inspection and restrictions; headings never grant authority. |

Retain existing translations for unchanged labels. No nested submenu. Do not move App settings from Account in this minimal change: document that it manages connection/app identity, while Tools manages evidence/investigation. App dropdown continues to switch app and never executes a tool. Scope title identifies the full name even when the masthead truncates a long app name. App change/access loss closes Tools and clears obsolete targets under existing lifecycle rules.

Use normal button/select/checkbox semantics in a labeled disclosure/dialog region, not ARIA menu/menuitem semantics with an improvised arrow-key model. Visible label, accessible name, expanded state and controlled panel must agree. Enter/Space opens while focus stays on trigger; Tab enters Close then controls in reading order, Shift+Tab reverses; leaving the panel by Tab closes it without stealing focus. Escape/Close returns to trigger; outside pointer dismissal leaves focus at the clicked target. Opening Account, Help, Coverage or Context replaces Tools and vice versa. Destination dialogs receive focus and restore it to Tools on close because the original action has been removed. Provide visible focus and no focus trap in the nonmodal panel.

Use at least44×44px interactive targets; panel width up to320px bounded by viewport minus24px, internal vertical scroll with reachable Close, wrapping French labels, no page overflow. Expand the labeled trigger horizontally within the existing row; allow the view strip to yield width and retain focus-revealing horizontal scroll. No new header row. Empty/normal historical fixture must keep canvas start at165px desktop and approximately200px at390px EN/FR (allow1px rounding); material warnings and active narrowing may legitimately consume space. Opening/closing must produce zero canvas displacement.

Selection indication belongs to actual state: retain Journeys/Traces and seven-view pressed states; expose highlight value and discovered-event checked state. Launcher actions are not persistent destinations and must not receive a false “current page” mark. Their opened dialog title identifies location and Close restores the originating workspace.

### Reasonable alternative: temporary navigation drawer

Use a labeled Tools button to open an overlaid side drawer containing the same groups; never a persistent sidebar. A wider panel improves long labels and future tool growth, but takes more graph area, adds a modal/focus-trap decision on phones and gives occasional actions undue navigation weight. It still needs all dismissal and scope rules. Prefer the smaller anchored disclosure for the current ten controls; reconsider a drawer only if the inventory grows materially.

## Risks / Trade-offs

- More trigger width shortens the phone view strip → retain horizontal view navigation, readable Outils label and measured no-wrap toolbar geometry.
- Grouped panel grows vertically → bounded scrolling and persistent Close; no nested navigation or graph resizing.
- “Tools” is broad → explicit group names and selected-app title; “Saved views” gives References a familiar task name without changing saved data.
- Reusing Disclosure blindly misses focus exit/app-change coordination → verify those behaviors explicitly; do not broaden a shared component refactor unless necessary.

## Migration Plan

After Alex approves, change presentation/localization only, preserve app/view state and routes, run targeted acceptance below, and revert presentation changes if regressions occur. No data migration or flag is needed. Preserve both October31 archives and their outstanding native zoom/screen-reader checks. This proposal adds its own manual acceptance and does not mark previous checks complete.
