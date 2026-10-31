# Tasks

Alex approved the compact Tools design and implementation. Automated implementation checks are recorded in IMPLEMENTATION-HANDOFF.md; native accessibility/device acceptance remains separate.

## 1. Present grouped tools

- [x] 1.1 Replace the hamburger with localized Tools/Outils in the same toolbar slot; group and retain the complete destination inventory in design.md. Verify owner/viewer inventories, highlight availability, selected states, and unchanged destinations with focused render/browser checks.
- [x] 1.2 Add the bounded, wrapping panel, selected-app title, Close and44px targets in the preserved visual theme. Verify1280×900 and390×844 EN/FR screenshots, long app names, no page overflow, internal scroll to every action, seven-view reachability and normal-fixture canvas starts165px/approximately200px with zero opening displacement.
- [x] 1.3 Update USER-GUIDE.md and LOCAL-REVIEW.md for Tools, groups and Saved views wording; verify each documented route against the rendered UI and preserve Account/App settings instructions.

## 2. Make interactions predictable

- [x] 2.1 Implement labeled nonmodal disclosure semantics and Enter/Space, Tab/Shift+Tab exit, Close, Escape and outside dismissal. Verify real browser key events and pointer/touch checks preserve intended focus and never trap it; opening another disclosure replaces the first in both directions.
- [x] 2.2 Close Tools before existing action dialogs and on app/access transitions; preserve permissions and dialog return to Tools. Verify owner/viewer actions, empty app, switch during slow reads, access expiry, destination return, and current investigation state without mutating a shared receiver.

## 3. Integrate and accept

- [x] 3.1 Run frontend check/build and existing relevant navigation, invitation and viewport tests; record current source/asset identity, commands and outcomes in a concise handoff. Confirm no API/data migration and no changes to archived records.
- [ ] 3.2 Observe actual200% browser zoom and native keyboard/screen-reader operation in EN/FR plus phone touch use; record label, expanded state, grouping, selected state, dismissal and focus observations. Keep any unavailable check explicitly pending; do not use headless width as native zoom/screen-reader acceptance or close earlier archives' pending checks.
