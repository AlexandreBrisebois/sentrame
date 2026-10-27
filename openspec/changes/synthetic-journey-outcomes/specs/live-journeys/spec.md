# Live Journeys

## Purpose

Show newly accepted journey activity within a selected application or visit while preserving historical context, viewer scope and honest latency. This is a proposed contract.

## ADDED Requirements

### Requirement: Switchable live scope
The system SHALL stream selected-app traffic or one isolated anonymous visit/session through SSE. Authorized viewers SHALL be able to switch freely while preserving applicable app/version/date/step filters. The current scope SHALL remain explicit and no event outside the authorized app/visit scope SHALL be disclosed.

#### Scenario: App to isolated visit and back
- **WHEN** a viewer selects one visit from an app stream and later returns to all selected-app traffic
- **THEN** the selected app and applicable filters persist and the interface states which traffic is now included

#### Scenario: Restricted invited viewer
- **WHEN** an invite is restricted to one live visit
- **THEN** changing a URL, filter or SSE request cannot widen that viewer's app, visit, history or diagnostics access

### Requirement: Clean live window and fixed baseline
The system SHALL offer a toggle between a clean live window and the selected app's accumulated traffic over the last seven days, with the chosen live visit highlighted over it. App-baseline access SHALL be authorized separately from visit access. Compatible app/version/definition filters SHALL remain explicit. Baseline bounds SHALL freeze at live start; counts SHALL remain distinct without overlay double counting.

#### Scenario: Baseline toggle
- **WHEN** a viewer enables the baseline after starting live observation
- **THEN** the display shows selected-app accumulated traffic for seven days ending at observation start, with the chosen live visit highlighted and separately labeled counts since start

#### Scenario: Filter change
- **WHEN** an applicable filter changes
- **THEN** each layer shows its declared authorized scope, compatible filters and definition, preserves its window and exposes recomputation or inactive filters

#### Scenario: App baseline behind one live visit
- **WHEN** a viewer authorized for app history selects one live visit and enables the baseline toggle
- **THEN** the historical layer remains labeled app-wide while only that visit is highlighted as live; neither count is represented as the other scope

#### Scenario: Restricted viewer enables baseline
- **WHEN** a visit-only viewer lacks app-history permission
- **THEN** a toggle, URL or API request cannot disclose the app baseline; the invitation policy governs whether access is granted, and any separately permitted earlier visit detail does not change baseline semantics

### Requirement: Resumable bounded streaming
The system SHALL bound per-viewer buffers and connections and support replay or an explicit gap/reset after SSE disconnection. Reconnection SHALL NOT silently duplicate displayed contributions. Slow viewers SHALL NOT block ingestion. Authorization revocation SHALL stop affected streams and prevent resumed access.

#### Scenario: Connection interruption
- **WHEN** a client reconnects with a prior cursor inside the replay window
- **THEN** it receives the missing authorized updates without inflating counts

#### Scenario: Expired cursor or slow viewer
- **WHEN** replay is unavailable or a viewer exceeds its buffer
- **THEN** the UI receives an explicit gap/reset and refreshes an authorized snapshot rather than pretending continuity

### Requirement: Consistent snapshots and lifecycle revisions
The system SHALL pair snapshots and cursors from one consistent committed revision. Timeout, expiry and deletion SHALL update or invalidate open views within a documented bound without new telemetry. Restart or notification loss SHALL recover by replay or explicit reset. Updates from superseded filter subscriptions SHALL NOT contaminate the current view.

#### Scenario: Silence crosses timeout or expiry
- **WHEN** no new record arrives as a journey becomes inactive or its detail expires
- **THEN** the open view updates or is invalidated within the declared freshness bound and does not indefinitely retain active or expired counts

#### Scenario: Commit during snapshot handoff
- **WHEN** an event commits while a snapshot/subscription is established, or the service restarts before publishing its notification
- **THEN** the consistent cursor and replay/reset policy prevents silent omission or count inflation

#### Scenario: Rapid filter switching
- **WHEN** old subscription updates arrive after a viewer changes filters
- **THEN** the UI discards or isolates them and shows only the current authorized subscription revision

### Requirement: Visible event time and observation latency
The system SHALL distinguish occurrence time, receipt time and display freshness. Exporter buffering, offline clients, retries and clock uncertainty SHALL be visible where known. The interface SHALL NOT promise instantaneous action visibility or treat paused/disconnected viewing as absence of application activity.

#### Scenario: Late batch during a clean window
- **WHEN** evidence received after live start occurred before the live window
- **THEN** the active event-time policy places it consistently, marks it late and does not silently move historical traffic into live action counts

#### Scenario: Pause viewing
- **WHEN** a viewer pauses live animation or the connection becomes stale
- **THEN** collection continues, the display marks its state and resumption reconciles counts within authorized replay limits

### Requirement: Qualified streaming through the shared tunnel
The SSE response SHALL use `text/event-stream` with explicit authenticated non-cacheable handling. Qualification SHALL cover flushing, buffering, heartbeat/idle behavior, disconnect/reconnect and revocation through the actual shared Tunnel route. HTTPS reachability or a direct-origin stream SHALL NOT establish external streaming compatibility.

#### Scenario: Proxy buffering or disconnect
- **WHEN** a controlled future qualification streams through the shared Tunnel and the proxy buffers or disconnects
- **THEN** the test distinguishes source delay from route buffering, verifies replay/reset and authorization, and records the route unqualified until required behavior is demonstrated

### Requirement: Provisional ongoing stage time
The live view SHALL distinguish an ongoing stage's provisional elapsed/observed-active time from completed duration statistics. It SHALL expose source/receipt latency and clock/activity uncertainty. Closing, late correction, replay or tab switching SHALL reconcile the same occurrence without double counting; ongoing time SHALL NOT silently enter completed-duration percentiles.

#### Scenario: Live stage closes after a delayed batch
- **WHEN** an entry is visible live and its matching end arrives later
- **THEN** the provisional display reconciles to one qualified occurrence under its timing policy, preserving exporter/tunnel delay context and adding at most one supported contribution to completed summaries
