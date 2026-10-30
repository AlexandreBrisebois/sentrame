# Using Sentrame

Sentrame turns received application events into an explorable journey map. Use it to see what happened before and after a stage, find unusual routes, inspect errors, and compare live activity with recent traffic. A graph shows observed evidence; a missing connection, quiet visit or error does not by itself prove failure or its cause.

This guide uses the English interface labels. **EN / FR** changes the interface language. Screenshots below show synthetic review data, not a connected production application.

## Open the workspace

For the existing local review workspace, run this from the Sentrame repository and leave the terminal open:

```sh
task start:workspace
```

Open [http://127.0.0.1:8096](http://127.0.0.1:8096). Use the exact configured address: `127.0.0.1` and `localhost` are different origins for access checks. The owner signs in with the configured passphrase. The local setup stores it in the ignored `.sentrame/local.env` file as `SENTRAME_PASSPHRASE`; keep that file private. There is no shared default password.

A restart requires a new sign-in. For first-time installation, synthetic demonstrations and REST Client requests, use the [local setup and review guide](LOCAL-REVIEW.md). For container configuration, storage and network access, use the [NAS launch guide](NAS.md). Preparing a container does not establish that a particular NAS or external network route has been qualified.

## Choose your application

Sign in, choose **Add app**, enter a unique name and register. Connection settings open with generated browser/server credentials, masked until reveal or copy. Configure your existing standard exporter; no Sentrame SDK or app-ID field is needed. Server credentials stay server-side. New apps have empty mappings.

Use the app selector for separate products and **App settings** to rename, rotate a role or remove an app after typing its name. Each app has its own history, mappings, references, lifecycle and invitations. App URLs and tabs keep their own selection; switching closes old editors/streams and clears old app data. Journeys is first and opens by default for owners and viewers, including after switching from Traces. Explicit trace links and saved trace references open Traces. A journey empty state stays in Journeys with connection/mapping/correlation guidance.

Temporary invitations name the selected app and its read-only whole-app scope and lifetime. They do not grant app administration, other apps or credentials. See [storage, recovery and limits](MULTI-APP.md). This release is greenfield: register apps again in a fresh directory; old single-app data and credentials are not imported.

## Start with the graph

The workspace keeps the graph in the foreground. Select a node or connection for details; use **☰** for less frequent tools.

| View | When to use it | What its shapes mean |
| --- | --- | --- |
| **Hierarchy** | Follow progression from left to right. | Cards are stages; arrows show observed direction; thicker connections mean more transition occurrences. |
| **Force** | Explore connected groups, loops and less linear paths. | Circles are stages arranged into a settled network. It is not a continually moving simulation. |
| **Sankey** | Compare where traffic goes and how it splits. | Bar height and ribbon width represent occurrence volume on a shared scale. Returns and repeats remain visible. |

Traffic is not a count of people. One journey may visit a stage several times, so node occurrences, transition occurrences, journeys and visits can differ. A repeat such as Vote → Vote can be healthy repeated activity.

Click a node to open its inspector. **Incoming** shows observed connections into it; **Outgoing** shows connections from it. Select a neighboring row to inspect that connection. The inspector also shows stage counts, **Time at stage**, mapped event names and **Show journeys**. A connection inspector shows that transition’s counts and access to its ordered evidence. Selecting a node emphasizes its neighborhood; it does not delete other evidence.

![Hierarchy with a rare-path highlight in synthetic review data](verification/workspace/hierarchy-rare.png)

### Move around without losing your place

- Drag empty space to pan. Two-finger scrolling also pans; pinch handlers support zooming at the pointer. Use **−**, **+** and **Fit** if gestures are unavailable.
- Drag a node to make a route easier to read. **Force** allows both axes; **Hierarchy** and **Sankey** keep their columns and allow vertical movement within routing bounds. Connections follow the node.
- A click opens the inspector; a drag does not. Each view keeps its own moved positions while that graph remains mounted. Reloading the page clears manual placement.
- **Reset layout** restores the current view’s automatic positions. **Fit** changes the camera to show the graph; it does not reset node placement.

For keyboard use, Tab to nodes and connections, then press Enter or Space to inspect them. With a node focused, Shift plus arrow keys moves it, subject to the view’s movement rules. With the canvas focused, arrow keys pan, `+`/`−` zoom and `0` fits. Keyboard focus remains visibly marked.

**Table** provides numeric map and timing alternatives. **Journeys & evidence** shows ordered paths and individual retained records. These views are also useful when many connections overlap. Physical trackpad, touch and assistive-technology behavior should be checked on the device you intend to use.

## Events, steps and activity

An **event** is an incoming telemetry record, such as `vote.committed`. A **mapping** gives that event a journey meaning. A **step** is one mapped occurrence in an ordered journey; a graph stage groups occurrences with the same stage label.

For example, ten separate `vote.committed` events mapped to **Vote** can produce ten Vote occurrences and nine Vote → Vote transitions. They do not become ten different stage names.

**Activity** is supporting evidence that does not add a step to the graph. An event saying a tab became hidden can help explain observed active time without cluttering the business path with another node. It remains evidence, not a measurement of a person’s attention.

Standard telemetry can be accepted without optional journey identifiers. Received records that cannot support a journey remain **Unassigned records**. Sentrame does not invent a visit or connection to make them fit.

If the historical timeline has bars but the journey graph is empty, those bars count **received records**, not journeys. When the selected range contains unassigned records, the empty graph confirms receipt and offers **View received records**. **Journeys & evidence → Unassigned records** opens automatically when there are no displayed journeys; it shows the event name, producer/signal, service/version, occurrence time and any emitted stage hint. **☰ → Unmapped events** also lists received event names across the retained app history.

A standalone log such as `pwa.app.opened` without a trace ID cannot appear in **Traces**. Without supported journey correlation, it also cannot become a journey merely by assigning a stage label. An empty trace list offers **View received records** to return to the journey workspace. For a live investigation, widen the historical range if the live observation began after the record arrived.

## Discover and configure events

Use **☰ → Unmapped events** to see received event names without mappings. **Show discovered events** puts them on the canvas as dashed, unconnected nodes. A disconnected discovery node means no journey relationship has been inferred; it does not mean the receiver failed.

Select an event node, then **Configure event…**. For an already mapped stage, find the event under **Events in this stage** and select it there. Each configuration field has an **ⓘ** button with guidance; click it or focus it and press Enter/Space.

| Field | What to enter |
| --- | --- |
| **Stage** | A short graph label, such as `Voting` or `Vote saved`. Several events may share a stage name. For a timing start/end pair, use the same stage name. Activity still requires a label, but does not add a graph node. |
| **Kind** | How this event affects the path or timing; see below. |
| **Trusted producer** | **Browser**, **Server** or **Any source**. Choose Server when the event must prove a server-side saved change. This filters by the receiver’s configured producer assignment; it does not make a browser claim trustworthy. |
| **Condition key** | Optional retained attribute name to test, such as `reason`. |
| **Equals** | The exact value required for that key, such as `exhausted`. It is case-sensitive literal text. |

### Choose the event kind

| Kind | Effect |
| --- | --- |
| **step** | Adds an occurrence to the journey path. Use for an action or milestone such as Vote saved. |
| **explicit exit** | Ends the journey with its emitted reason. Exhausting available recipes, for example, is not automatically a failure. |
| **stage start** | Adds an entry step and starts explicit stage timing. |
| **stage end** | Supplies the matching timing end without adding another graph step. |
| **activity** | Keeps supporting evidence without adding a graph step. Recognized `activity.state` values—`active`, `visible`, `idle`, `hidden`—can support observed active-time calculations. |

A native span uses its own start/end and native identity. Explicit paired timing is an alternative for a broader interval and needs compatible start/end evidence with the same stage label, emitted `stage.occurrence.id`, producer and service. Changing the Kind dropdown alone cannot supply missing instrumentation. All mapped kinds refresh the inactivity timeout. An ordinary mapped step may instead have **Next-event estimated** timing when that policy is enabled; that is an estimate, not an observed end.

### Use conditions only when necessary

Leave **both Condition key and Equals empty** to match this event by name and selected producer alone.

Suppose **Received examples** shows this attribute on `journey.exited`:

```json
{"reason": "exhausted"}
```

Enter:

| Condition key | Equals |
| --- | --- |
| `reason` | `exhausted` |

That matches `journey.exited` records with `reason` equal to `exhausted`. A record with `reason: agreed`, or no `reason`, does not match this example. Type `exhausted` without quotes; do not enter the whole JSON object or `reason=exhausted` in Equals. `Exhausted` is a different value.

Another example is **Condition key:** `activity.state`, **Equals:** `hidden`. Use keys and values actually present in **Received examples**. Attributes removed by the privacy policy cannot be used in conditions.

Choose **Preview matches** before **Activate new definition**. The preview checks retained examples and reports matching records and missing safe journey context. Activation creates a prospective definition; it does not rewrite existing captured journeys under the new rules.

Use **☰ → Mappings & goal** for the active configuration and **Set operational goal**. The operational goal is the event that completes future journeys. Merely investigating a stage or error does not make it a goal. A captured journey starts at its first mapped occurrence, completes at its configured goal, or closes as inactive after 30 minutes without mapped activity. The next genuine mapped occurrence starts a new journey. Inactive means the outcome is unknown, not failed. Ongoing captured journeys retain their original definition.

## Select a time range or follow one visit

Select **Historical** or **Live** to open **Time & filters**. In Historical mode, histogram bars show retained record counts by occurrence time. They do not count journeys, people, or metric values.

1. Choose a chart context: **1h**, **24h**, **7d** or **30d**.
2. Drag across bars, or move the **Start** and **End** handles.
3. Release to apply that interval to the graph. Green marks only the included span.

The chart context is independent of the graph selection, so you can see surrounding traffic while choosing a narrower slice. If existing dates extend beyond the chart, the handles show their visible portion; the graph keeps its range until you choose a new one. **Timeline data table** gives the counts numerically. **Quick range** offers recent intervals, and **Advanced: exact times or correlation ID** exposes precise bounds and known-ID search.

Dates are **UTC**, including the advanced fields, regardless of your computer’s timezone. The range includes its start and excludes its end: `10:00 → 11:00` includes an event at 10:00 but not one at exactly 11:00. A graph slice preserves window gaps instead of joining stages across omitted evidence.

![Shared timeline selection in synthetic review data](verification/workspace/timeline-slicer.png)

You do not need to know a session ID to follow a visit. Select a node, choose **Show journeys**, identify a journey by its path, time and outcome, and use **Follow this visit** or **Follow this operation**. The selected scope appears above the graph with an **X**. Clear it to return to all traffic under the other filters. A correlation ID groups related events; it is not a person’s identity, and it should not be reused to join separate visits.

## Highlight a pattern

The **Highlight** label describes the dropdown’s purpose: emphasize matching evidence while other parts of the graph fade. Highlighting does not remove data or change journey outcomes.

| Option | What it emphasizes |
| --- | --- |
| **None** | Normal graph emphasis. |
| **Rare paths** | Full sequences with at most 10% share among at least 20 comparable closed journeys under the same definition and version. Partial, ongoing and right-censored paths are excluded from this rarity cohort. |
| **Endings** | Stages where journeys are pending, inactive, explicitly exited or completed. Different outcomes remain named. |
| **Warnings & exceptions** | Journeys containing warning or error evidence. |
| **Recovery** | Journeys with observed error/recovery evidence that later completed. |

Rare does not mean defective. A highlighted link may also belong to a common path; rarity is based on the whole sequence. With fewer than 20 eligible comparable journeys, Sentrame does not guess which paths are rare.

Red connections identify local error-related transitions, with text/symbol equivalents and access to evidence. An error in one part of a journey does not make every subsequent connection red or undo a later completion.

## Explore and save references

Start from the question you want to investigate:

- Select a stage, then **Paths into here…** to look backward toward it.
- Select **Paths from here…** to inspect later recorded paths.
- Choose **Use as beginning…**, then select an end stage to explore a bounded segment.
- Use **References** to save the current investigation, or choose **Current workspace** for the whole current view.

Exploring these projections is read-only and does not require saving. A reference is a reusable query, not a new journey definition. Saving retains the time bounds, visit/version/definition filters, selected layout and highlight. Beginning/end modes allow the first or last repeated endpoint. **Open view →** applies the saved context. **Clear reference** returns to the ordinary filtered workspace.

References preserve captured journey boundaries. They cannot merge separate captured journeys, fill missing starts or resurrect expired detail. A live reference saves a fixed observation window; it does not save a combined live-plus-baseline dataset.

Owners can choose **Delete reference** on a saved card, then confirm or cancel. This removes only that saved query. It does not delete collected events, journeys or other references. Temporary viewers may open references and investigate paths, but cannot save or delete them.

## Watch live activity

Switch to **Live** to begin an observation window. New retained evidence arrives through server-sent updates; reconnect/reset refreshes the snapshot rather than adding the same counts again.

- With **Seven-day baseline** off, view the clean live window.
- Turn it on to compare against the selected app’s fixed seven days before live start. Baseline and live volumes are counted separately; a selected visit does not turn the baseline into a single-visit baseline.
- **Pause updates** pauses viewing, not ingestion. **Resume updates** reconciles the view with current evidence.

Buffering, sampling, retries and clock differences can delay or omit observations. A quiet live screen does not establish that nobody is using the application. An ongoing stage’s provisional time is separate from completed-duration statistics.

## Inspect errors and time honestly

**Errors**, or **☰ → Errors & exceptions**, opens application-wide diagnostics. You do not need to be the person using the application, and an error does not need a visit ID to be discoverable.

Choose **Whole application** or **Current graph filters**, search by message, exception, service or trace ID, and expand a record. Inspect its message, available stack, producer, occurrence/receipt times and technical IDs. **Explore related context** follows available correlation; proximity alone is not proof that an error caused a particular result. **Copy sanitized context** copies the retained projection.

Messages, exception stacks and custom error-log bodies are bounded and scrubbed. Some details may have been omitted or truncated; earlier records can lack fields removed by an older policy. Stacks are shown as emitted, without applying source maps. Scrubbing is imperfect: controlled technical messages are appropriate, customer content and secrets are not. Synthetic .NET-like and Next.js-like examples do not qualify actual .NET or Next.js exporters.

![Scrubbed diagnostic inspection in synthetic review data](verification/workspace/diagnostic-inspector.png)

In **Time at stage**, compare values only when their basis and definition are compatible:

| Timing basis | How to interpret it |
| --- | --- |
| **Measured elapsed** | Explicit entry-to-exit wall-clock duration. It may include idle or background time. |
| **Observed active** | Duration supported by explicit activity intervals, with bounded heartbeat assumptions. It does not measure attention. |
| **Next-event estimated** | An enabled fallback using adjacent compatible evidence. It is not an explicit stage end. |

**Median** describes the middle completed sample. **Slower cases (p90)** requires at least 20 completed compatible samples; ongoing, missing or excluded observations are not zero. Review sample count, exclusions and completeness alongside any duration. Missing clocks, mismatched occurrences or incomplete boundaries can make timing unavailable. Occurrence time, receipt time and **Receipt lag** answer different questions; clocks may differ.

## Retention, collection health and limits

Open **☰ → Data lifecycle** for retention and summary information.

- **Ordered detail:** 30 days by occurrence time, with a disclosed receipt-time fallback where supported. Expired records no longer appear in detailed reads.
- **Retained summaries:** supported app-wide whole-UTC-day counts and timing distributions for 12 calendar months. They cannot reconstruct old individual paths, visits or arbitrary new filters. Summary timing quantiles are approximate histogram bins; daily medians are not averaged.
- **Collection health:** reported completeness, privacy, freshness and verified-transport policies. It is not a live guarantee that every exporter event arrived.
- **Limits:** work caps for one view, including a 10,000-record model scan and 2,000 returned evidence records. The graph and error workspace have their own bounded displays. A truncated result is partial; narrow the time range or filters. These caps are different from retention.

The correction horizon ends at captured journey start plus 30 days. After detail expires, Sentrame cannot subtract one individual’s contribution from finalized summaries. Owner deletion tools distinguish a retained visit, a whole UTC summary day, and **Delete app data and disable intake**. Confirm the intended scope before proceeding. Whole-app deletion disables that intake generation. Client copies, filesystem snapshots and backups are outside server erasure guarantees; visit/journey correlation does not provide unlimited replay protection.

## Invite a temporary viewer

Owners open **☰ → Temporary viewers**:

1. Add an optional **Viewer label**, such as Product review.
2. Choose **Access expires after**: **1 hour**, **24 hours** or **7 days**. The default is 24 hours, and the clock starts when the link is created.
3. Choose **Create invitation link**, then **Copy invitation link**. Copy it now; Sentrame does not show the original link again from the saved invitation list.
4. Send it privately to the intended viewer. Anyone holding the link can redeem it once. To end access early, find the entry under **Invitation access** and select **Revoke access**.

The viewer opens the link and explicitly chooses **Open read-only workspace**. The invitation credential is removed from the address bar before normal app requests. If already signed in as the owner, sign out first or open the invitation in a separate browser profile; an invitation does not downgrade the owner’s current access.

A temporary viewer can browse the whole selected application’s retained history, live maps, baseline, references, errors and sanitized diagnostic details. They may change their own view filters, layout and highlights, follow a visit, and copy displayed sanitized evidence. They cannot change mappings or goals, save/delete references, delete collected data, or create/revoke invitations. The server enforces these restrictions independently of hidden buttons.

The viewer sees **Read-only access** and an expiry time. Expiry or revocation blocks subsequent protected requests; active live streams are checked within five seconds. The open app rechecks access even while historical viewing or live updates are paused. Once access ends, it clears the workspace and asks for a new invitation. Previously copied information cannot be recalled.

Viewer sessions also end on service restart. A redeemed link stays used and cannot sign in again; create a new invitation if access is needed afterward. An unused, unexpired, unrevoked invitation can still be redeemed after a restart.

A `127.0.0.1`, `localhost` or loopback IPv6 link works only on the machine opening it. Creating an invitation does not expose the service to another computer. An external viewer needs an already configured, reachable address and the correct origin; see the [NAS guide](NAS.md). Actual external-route and Home Server qualification remain separate.

## When the screen does not show what you expected

| Symptom | What to check |
| --- | --- |
| Empty graph | Clear a selected visit/reference, widen the time range, check the version/definition filters, then inspect **Unmapped events** and **Unassigned records**. Successful receipt without safe journey context is different from ingestion failure. |
| Event appears but has no connections | It may be unmapped or lack supported correlation. Configure its meaning and inspect received context; do not infer a link from nearby timestamps. |
| No rare-path highlight | At least 20 closed, comparable, non-partial/uncensored journeys are needed. A low count is not evidence that all paths are common. |
| Timing is unavailable | Check native span start/end or the explicitly selected pair, compatible source/clocks and completed samples. For p90, check the 20-sample minimum. |
| Counts differ between histogram and graph | The histogram counts retained records; graph nodes/edges count mapped occurrences/transitions. Journeys and visits are separate units. |
| Results say partial or truncated | Narrow the period or filters. Displayed counts are not complete population totals. |
| New UI is not visible | Refresh to load built browser assets. A backend change needs a service restart; restarting requires sign-in again. |
| Viewer link no longer works | It may be expired, revoked or already used, or its redeemed session may have ended on restart. Ask the owner for a new link. |
| Another device cannot open a local invite | Loopback points at that other device, not this service. Review the configured reachable address; an invite cannot create network access. |
| Real exporter sends nothing visible | Check the protected receiver’s advertised signal/HTTP combination, producer credentials and emitted timestamps/context. Use the setup/integration guides; this UI cannot establish delivery completeness. |

For repeatable sample traffic and actual OTLP request files, return to the [local review guide](LOCAL-REVIEW.md). Qualification notes distinguish what was source-tested, exercised locally or not run; do not treat a demonstration as proof of every SDK, route or device.


## Trace exploration

After sign-in, **Traces** lists observed execution contexts, with service and observed-span counts. Select one, then choose **tree**, **timeline** or **table**. Parent edges describe native context; parallel children never become chronological transitions. Missing and ambiguous targets stay labeled. Open **Typed relationships**, **Span events in export order** and **Correlated logs** for retained context. Logs with only a trace ID remain trace-level evidence. **Copy sanitized trace** copies the same protected projection.

Exact start/end and duration nanoseconds remain decimal strings. Durations are inclusive operation time, not exclusive time, page dwell, attention or a sum of unique time. Missing ends are right-censored; invalid or receipt-fallback clocks are unavailable. Sampling, SDK drop counts, projection link omissions and read truncation disclose partial evidence. A sampling flag supplies no sampling rate. Turn off **Automatic refresh** to pause updates; Refresh reconciles a full bounded read. Revocation clears access through the shared session boundary.

Switch to **Journeys** for configured chronological business paths. Trace-only telemetry does not create a journey, person count or conversion funnel. The two modes preserve separate counting units. Service/carrier contextual mappings and explicit definition previews are available at the qualified local seam; optional stage hint suggestions are available for explicit owner acceptance. See [qualification](verification/standard-otel.md) for what actually ran; wire tests do not qualify a browser exporter.


## Give ordinary telemetry business meaning

Received logs also offer **Configure event** directly under **Unassigned records**. The editor fills in the received service, producer and carrier when there is one observed context. When an event exists in several contexts, choose **Received event context** to select the actual service and producer; exact service suggestions are also available. **Received records for this mapping** loads automatically from retained app-wide history, independently of the graph time range. It shows the matching count and up to five recent examples with their retained attributes. **Pick a condition from received data (optional)** fills both Condition key and Equals from an observed field/value, with its occurrence count. Leave both fields empty to match all records in the selected context. If a condition matches nothing, the editor shows examples and alternative values from that context before the condition. Choices are limited to 64 retained, allowlisted field/value pairs, with omissions disclosed. A condition or selector change requires another **Preview matches** before activation. On an app with no goal yet, the first definition requires an explicitly chosen operational goal. Select **End a journey when this event occurs (operational goal)** only for the intended ending event, then choose **Preview matches**. Any edit requires another preview. The activation button explains an unmet stage, goal, preview or conflict requirement. Missing correlation is reported separately and does not block saving a valid definition; those records remain unassigned.

On a trace operation, span event or log, choose **Give business meaning**. The editor starts with the retained exact service/carrier/source/name. Enter a literal Stage label, choose Kind and an optional retained key/value condition, then **Preview matches** before activation. Select **Use this exact mapping as the operational goal** only for the actual intended boundary. No stage tag is needed. A trusted `vote.committed` span event can map to **Vote saved** while a similarly named log remains separate.

Native-span timing measures inclusive operation elapsed time using native identity. Point observations have no duration by label. Explicit paired boundaries and next-event estimates are separate policies; a broader page stage needs real supporting observations. Do not sum nested operations into unique page time. An end-boundary span goal requires a valid completed span; a named event goal uses its event timestamp. A generic successful HTTP request and browser intent do not prove persistence.

Default trace labels use capture-time definitions. Use the interpretation menu for an explicit non-persisting preview under another definition. Activation changes future captures, without rewriting retained evidence, previous labels or captured journeys. Two matching explicit rules block activation; there is no hidden rule order. Read-only viewers may inspect previews but cannot configure or activate.


## Optional emitted label suggestions

An unmapped record may offer **Use emitted suggestion: voting.persist**. Choosing it fills a literal Stage in the editor; Preview and activation remain required. There is no automatically activated mapping. Different emitted labels remain separate suggestions, with at most32 per contextual selector and omissions disclosed. Missing/invalid tags offer none and do not prevent ordinary-name mapping.

One explicit owner mapping always wins. If the emitter says `voting` while the owner says **Vote saved**, the stage remains **Vote saved**, and the emitted hint remains inspectable. Two overlapping explicit rules must be disambiguated. Child records do not inherit a span's label. Changed emitted hints do not rename prior captures or establish a goal, Kind, session or authority.


## Scoped native timing and references

Explicit scope/tab can relate multiple request traces into one configured journey. The graph counts chronological stage occurrences; opening a trace still shows its native operation structure. Native span timing now appears as **Native inclusive operation**, separate from **Measured elapsed**, **Observed active** and **Next-event estimated**. Each distribution carries its definition/computation, sample count and exclusions; p90 needs20 compatible completed samples. Missing ends/invalid clocks do not become zero samples, and competing timing bases for the same occurrence are unavailable until the owner chooses one.

In Traces, choose or save a **Trace reference** with a name, current view and definition interpretation. It stays distinct from journey goal/segment references. An explicit preview pins that definition; the default uses recorded capture context. Expired/deleted targets remain unavailable even when the saved query survives. A late earlier mapped prefix within the supported inactivity gap retains the known capture definition instead of adopting the latest labels. No personal identity index is created.

Observed-active time requires actual browser visibility/idle observations within a complete native browser interval or an optional paired page interval in the same capture. It measures capped positive observed intervals; gaps remain unknown, and server operation time cannot borrow browser presence.
