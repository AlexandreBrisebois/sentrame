# Spec Delta

## ADDED Requirements

### Requirement: Consistent live and baseline identity
All journey views, legends, inspectors, example lists and diagnostics SHALL consistently distinguish live and baseline using distinct hues, visible population labels and non-color cues. Population identity SHALL remain distinguishable from volume, recency, error, selection and unmapped-state encodings, including when a relationship belongs to both populations.

#### Scenario: Shared relationship
- **WHEN** the same directed transition exists in baseline and live evidence
- **THEN** both contributions remain independently identifiable and inspectable, with separately labeled counts and no blended total
- **AND** isolating either population does not restart observation or change captured evidence

#### Scenario: Error in both populations
- **WHEN** baseline and live both contain error evidence on a shared connection
- **THEN** error markings preserve population labels and non-color distinction rather than making both populations indistinguishably red
- **AND** error evidence is not represented as proof that either whole journey failed

#### Scenario: Color-independent interpretation
- **WHEN** the developer uses monochrome, reduced motion or assistive technology
- **THEN** baseline and live remain distinguishable without hue, flashing, reduced opacity alone or a hover-only legend

### Requirement: Visible observation and baseline contract
Live workspace context SHALL identify observation start, live bounds and scope, the fixed seven-day baseline bounds and app-wide scope, applicable version/definition filters and each population's coverage. Toggling visibility SHALL preserve these bounds. Baseline recomputation SHALL disclose corrections without silently shifting the comparison window.

#### Scenario: One live visit against app history
- **WHEN** one selected visit is observed over an authorized app-wide baseline
- **THEN** the interface explicitly labels that scope difference beside both populations, including when a population has zero retained journeys
- **AND** it does not imply a matched-visit comparison

#### Scenario: Baseline correction after late intake
- **WHEN** newly retained evidence or lifecycle changes revise counts within the baseline's fixed bounds
- **THEN** the baseline remains on the same bounds and a correction or refreshed-context indication distinguishes that change from new live activity
- **AND** late evidence is classified under the existing occurrence-time policy rather than by receipt color

### Requirement: Population-safe comparison metrics
Comparison SHALL expose counts, units, denominators, coverage, scope, window length and definition compatibility for both populations. Raw live-versus-baseline percentage change SHALL NOT compare unequal windows or scopes. Only supported within-population shares with compatible meanings SHALL offer percentage-point differences; unsupported comparisons SHALL explain why they are unavailable.

#### Scenario: Minutes against seven days
- **WHEN** five minutes of live traffic is displayed against seven days of baseline traffic
- **THEN** raw counts remain separately labeled and the interface does not claim a traffic decrease from their raw difference
- **AND** no implied statistical significance or release causality is assigned to an observed difference

#### Scenario: Scoped or incomplete population
- **WHEN** one live visit is compared with app-wide history, or either source population is incomplete or definition-incompatible
- **THEN** numeric change badges are unavailable with a plain explanation and separately labeled observed counts or shares remain inspectable

#### Scenario: Supported share difference
- **WHEN** both populations are app-wide, have compatible stage definitions and complete supported transition denominators
- **THEN** an optional outgoing-occurrence share difference names both denominators and uses percentage points, while disclosing unequal windows and available sample support
- **AND** a zero baseline denominator is unavailable rather than infinite change

### Requirement: Comparable heatmap and flow scales
Side-by-side population views SHALL align stage or prefix identity and declare their metric scale. Comparable heatmaps SHALL use the same numeric domain with distinct population legends. Shared graph volume encodings SHALL retain a common unit and scale. Live rescaling SHALL be explicit and SHALL NOT silently change the meaning of unchanged values.

#### Scenario: Equal counts in two heatmaps
- **WHEN** a live and baseline cell both contain ten occurrences
- **THEN** each uses the equivalent position in the shared numeric scale with its population hue and visible count, even when other cells have very different maxima
- **AND** independently normalized views are not presented as a magnitude comparison

#### Scenario: New count exceeds a locked domain
- **WHEN** live counts exceed the current heatmap domain during inspection
- **THEN** exact values remain visible, capped intensity is disclosed, and a rescale action is available without silently recoloring every unchanged cell
- **AND** displayed normalized shares are not substituted for occurrence-volume widths in the existing Sankey

### Requirement: Calm stable live investigation
Incoming activity SHALL update supported counts without automatically moving the selected journey, focused control, stage order, scroll position or graph arrangement. Bounded recent-update cues SHALL indicate changed live evidence without persistent flashing. Following new activity SHALL be explicit and reversible; background baseline corrections SHALL NOT impersonate live arrivals.

#### Scenario: Read while traffic arrives
- **WHEN** a developer inspects a matrix cell, repeat example or expanded timeline while a burst arrives
- **THEN** the same target stays selected and in position, changed evidence is indicated calmly, and new items can be revealed intentionally

#### Scenario: Follow mode and manual inspection
- **WHEN** following activity is enabled and the developer scrolls away or selects an earlier occurrence
- **THEN** automatic following stops with a clear resume-follow action while collection and authorized updates continue

#### Scenario: Fresh cue expires
- **WHEN** a contribution has remained unchanged beyond the documented recent-update interval
- **THEN** its freshness cue clears while its Live identity remains
- **AND** reduced motion conveys the same update state without animation

### Requirement: Controllable live presentation and freshness
The workspace SHALL distinguish connected, waiting, paused, reconnecting, stale, refreshing and failed presentation states, show last successful snapshot freshness, and explain that collection continues when viewing pauses. Resume SHALL reconcile authorized evidence without duplicate counts. Queued event counts SHALL be shown only when actually known.

#### Scenario: Pause during a burst
- **WHEN** updates arrive while viewing is paused
- **THEN** visible evidence and positions remain stable, an updates-available indication appears, and resume refreshes the authorized current state
- **AND** update notifications are not misrepresented as counts of unique records or journeys

#### Scenario: Connection without fresh data
- **WHEN** the transport is connected but no recent successful snapshot is available
- **THEN** connection state and data freshness are shown separately, with a retry action where appropriate
- **AND** the interface does not claim no application activity merely because updates are absent

#### Scenario: Reset or replay gap
- **WHEN** reconnection requires a reset rather than continuous replay
- **THEN** a continuity warning remains until reconciliation, counts are not duplicated, and still-valid investigation state is preserved
- **AND** previously displayed evidence that has expired is removed or invalidated under lifecycle policy even if presentation was paused

### Requirement: Provenance-preserving live drilldown
Every live or baseline selection and drilldown SHALL retain its population identity, window, scope and revision context. A shared stage or transition SHALL offer independently labeled supporting evidence for both populations. Missing live evidence SHALL NOT silently fall back to baseline evidence, nor baseline to live.

#### Scenario: Shared transition inspector
- **WHEN** a developer selects a connection present in both populations
- **THEN** the inspector shows independently labeled counts and example actions for Live and Baseline, and its timeline and diagnostic navigation retain the chosen population

#### Scenario: One side unavailable
- **WHEN** live detail is absent but a baseline example exists for the same stage names
- **THEN** the live selection explains its missing detail and baseline evidence requires an explicit labeled action

### Requirement: Bounded live updates and announcements
Live visualization work and accessibility announcements SHALL be coalesced and bounded without blocking intake or user input. New app, scope or filter contexts SHALL discard superseded updates. Revocation, deletion and expiry SHALL override retained presentation state; paused inspection SHALL NOT preserve unauthorized or expired raw evidence.

#### Scenario: Sustained updates
- **WHEN** incoming notifications exceed the rendering cadence
- **THEN** the view reconciles to a current authorized snapshot, keeps controls responsive and announces summarized changes instead of interrupting for every event

#### Scenario: Context changes during a refresh
- **WHEN** an earlier live refresh completes after switching apps or filters
- **THEN** its counts, freshness cues and selected examples cannot enter the new context
- **AND** loss of access clears the affected evidence and stops its update source

### Requirement: Accurate live availability language
Live views SHALL distinguish observed zero from no retained detail, insufficient coverage, waiting for activity, authorization limits and failed reads. Newly observed paths SHALL be described relative to the selected baseline window and coverage rather than as first-ever behavior.

#### Scenario: Path absent from baseline
- **WHEN** a path appears live and is not present in retained baseline evidence
- **THEN** it may be labeled not observed in this baseline, with bounds and coverage available
- **AND** it is not called new to the application or an anomaly solely because it is absent there

#### Scenario: Baseline not authorized
- **WHEN** a viewer lacks access to app history
- **THEN** baseline controls explain their unavailability without querying, displaying or copying restricted baseline evidence

### Requirement: Existing reset-snapshot transport semantics
The existing periodic five-second reset stream SHALL trigger bounded snapshot reconciliation. Routine periodic resets SHALL not be described as new continuity gaps; reconnect or an unexpected context SHALL remain marked until reconciliation. Real intake/SSE qualification SHALL be separate from a clearly identified high-rate client notification stress harness.

#### Scenario: Routine reset without new evidence
- **WHEN** a periodic reset reports the same revision on a continuous connection
- **THEN** freshness and lifecycle checks continue without a new gap warning or invented event count
