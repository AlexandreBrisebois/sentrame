# Synthetic client event plan

The fixture generator translates readable JSON into official standard OTLP and executes version-controlled REST requests through the actual receiver/store/protected-read seam. It uses disposable selected loopback state, one registered demonstration app and its single telemetry key. This is not authorization to touch a real application or retained deployment.

| Observation | Useful question | Evidence and limit |
| --- | --- | --- |
| Navigation/screen entry | Which path and completion screen were observed? | Explicit random client visit/tab and stable event identity; a screen observation is not proof of shared agreement or persistence |
| Intent/result observation | What did the client request and subsequently report? | Two records remain distinct; qualified action identity is separate from record replay; healthy repetition is preserved |
| Error/retry/recovery | What friction and later continuation were observed? | Supported explicit correlation, bounded scrubbed diagnostics; severity/time proximity does not establish cause or termination |
| Stage/page/presence | Which duration basis is actually supported? | Native, paired elapsed, next-event estimate and actual observed-active intervals remain separate; missing ends are censored, never attention |
| Missing context/claimed origin | Was valid telemetry received without a reconstructable journey? | Missing visit/tab remains unassigned; payload source/runtime/kind cannot grant stronger assurance |

Examples use `vote.result.observed` and optional `completion.entered`, not trusted commitments. Standard carrier/service/name mappings and optional stage hints add prospective meaning without a role selector or SDK replacement. Native span identity supplies operation timing; additional custom context is requested only for the client journey question.

Generator/parser tests are separate from actual wire/SDK/browser/container/deployment qualification. Record passed/failed/blocked/not-run checks for the current source version. Older reports and [historical server integration evidence](../verification/historical-server-integration.md) preserve previous results without advertising current server telemetry features.
