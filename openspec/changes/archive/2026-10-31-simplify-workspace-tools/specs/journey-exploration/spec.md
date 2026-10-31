# Spec Delta

## ADDED Requirements

### Requirement: Discoverable scoped workspace tools
The workspace SHALL offer a visible localized Tools disclosure in the existing compact toolbar. It SHALL group graph presentation, investigation and app-management actions, identify the selected app and preserve every authorized destination. App selection SHALL remain distinct. Existing direct view controls and graph space SHALL be preserved without adding header rows.

#### Scenario: Find a tool in either language
- **WHEN** an owner opens Tools in English or Outils in French
- **THEN** View, Investigate and Manage this app groups expose highlight, discovered-event visibility, saved views, save-current-view, unmapped events, mappings/goal, temporary viewers, errors, data lifecycle and evidence explanation with localized labels
- **AND** the current highlight and discovered-event visibility expose actual selected state, while launcher actions do not claim to be a current page

#### Scenario: Viewer or app transition
- **WHEN** a temporary viewer opens Tools or the selected app/access changes
- **THEN** only authorized actions are offered, existing read-only destinations remain reachable, and app/access changes close obsolete tool panels without leaking or submitting old-app state
- **AND** app switching remains in the separate app selector and owner App settings remains reachable from Account

#### Scenario: Preserve compact graph navigation
- **WHEN** Tools is closed, opened or dismissed at desktop or390 CSS-pixel width in English or French
- **THEN** all seven direct views remain reachable, phone view navigation retains horizontal scrolling and focus reveal, and opening or closing Tools causes no graph displacement or added header row
- **AND** the normal historical reference fixture retains a canvas start no lower than165px desktop and201px phone, allowing existing material warnings and intentional narrowing to consume necessary space

### Requirement: Predictable workspace tool disclosure interaction
Tools SHALL provide a labeled, viewport-bounded nonmodal panel with Close, Escape and outside dismissal, visible keyboard focus and at least44px interactive targets. Trigger state and panel relationship SHALL be accessible. Only one workspace/account disclosure SHALL remain open. Action dialogs SHALL preserve app context and recover focus to Tools when their invoking action is absent.

#### Scenario: Keyboard open and dismiss
- **WHEN** a person opens Tools with Enter or Space and traverses its controls with Tab or Shift+Tab
- **THEN** controls follow visible reading order, Escape or Close restores trigger focus, and tabbing out dismisses without trapping or stealing focus

#### Scenario: Pointer or competing disclosure
- **WHEN** a person clicks outside Tools or opens Account, Help, Coverage or Context
- **THEN** Tools closes, focus follows the intended outside target, and no overlapping tool panel remains
- **AND** opening Tools closes any preceding disclosure

#### Scenario: Open an action and return
- **WHEN** a person chooses a Tools action that opens an existing dialog and then closes it
- **THEN** Tools closes before that dialog opens, its localized title identifies the chosen destination, authorized app context is preserved and focus returns to the Tools trigger

#### Scenario: Narrow screen and zoom
- **WHEN** a person operates Tools at390 CSS pixels or200% native browser zoom in either language
- **THEN** labels wrap readably, every action and Close remains reachable through bounded panel scrolling, focus stays visible and the page has no horizontal overflow caused by Tools
