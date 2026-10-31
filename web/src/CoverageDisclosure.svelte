<script lang="ts">
 import Disclosure from './Disclosure.svelte';
 import type {Snapshot} from './types';
 import {translate,type Language} from './i18n';
 let {snapshot,language,live=false,baseline=false,warningsOnly=false,hideWarnings=false}:{snapshot:Snapshot|null;language:Language;live?:boolean;baseline?:boolean;warningsOnly?:boolean;hideWarnings?:boolean}=$props();
 const t=(s:string)=>translate(s,language);
 const coverage=$derived(snapshot?.coverage);
 const populations=$derived([{id:'current',label:live?'Live observation':'Current selection',cov:coverage?.current,model:snapshot?.model},...(live&&baseline?[{id:'baseline',label:'App-wide baseline',cov:coverage?.baseline,model:snapshot?.baselineModel}]:[])]);
</script>
{#snippet warnings()}
 <div class="coverage-warnings">
 {#if !coverage?.source}<p class="coverage-warning">{t('Source-record scan')}: {t('Coverage metadata unavailable.')}</p>{:else if coverage.source.truncated}<p class="coverage-warning">{t('Source-record scan')}: {t('Source scan capped.')}</p>{/if}
 {#each populations as pop}
  {#if !pop.cov||!pop.model}<p class="coverage-warning">{t(pop.label)}: {t('Coverage metadata unavailable.')}</p>
  {:else if pop.cov.sourcePartial||pop.cov.partialJourneys||pop.cov.censoredJourneys||pop.cov.graphNodesOmitted||pop.cov.graphEdgesOmitted}
   <p class="coverage-warning">{t(pop.label)}: {#if pop.cov.sourcePartial}{t('Observed partial evidence')} · {/if}{#if pop.cov.partialJourneys}{pop.cov.partialJourneys} {t('Partial beginnings')} · {/if}{#if pop.cov.censoredJourneys}{pop.cov.censoredJourneys} {t('Censored continuations')} · {/if}{#if pop.cov.graphNodesOmitted||pop.cov.graphEdgesOmitted}{t('Graph projection')}: {pop.cov.graphNodesOmitted} {t('omitted nodes')} · {pop.cov.graphEdgesOmitted} {t('omitted connections')}{/if}</p>
  {/if}
 {/each}
 {#if !coverage?.evidence}<p class="coverage-warning">{t('Received-record list')}: {t('Coverage metadata unavailable.')}</p>{:else if coverage.evidence.omitted}<p class="coverage-warning">{t('Received-record list')}: {coverage.evidence.omitted} {t('omitted records')} · {t('Display only; ordered-journey denominators unchanged.')}</p>{/if}
 </div>
{/snippet}
{#if warningsOnly}{@render warnings()}{:else}
<section class="coverage-disclosure" aria-label={t('Coverage and display limits')}>
 <Disclosure label="Coverage" note="collection unknown" {language}>
 <h3>{t('Coverage and display limits')}</h3>
 <p data-coverage="source"><strong>{t('Source-record scan')}</strong>: {#if coverage?.source}{t('Record limit')} {coverage.source.recordLimit} · {t(coverage.source.truncated?'Source scan capped. Counts describe scanned retained evidence; narrow filters.':'Source scan not capped. Retention and selection still apply.')}{:else}{t('Coverage metadata unavailable.')}{/if} {t('Exporter completeness is unknown; sampling and delivery loss are not measured.')}</p>
 {#each populations as pop}
  <p data-coverage={'population-'+pop.id}><strong>{t(pop.label)}</strong>: {#if pop.cov&&pop.model}{pop.cov.journeys} {t('Retained journeys in this selection')} · {t(pop.cov.sourcePartial?'Observed partial evidence':'Retained source')} · {pop.cov.partialJourneys} {t('Partial beginnings')} · {pop.cov.censoredJourneys} {t('Censored continuations')}. <span>{t('Graph projection')}: {pop.model.nodes.length} {t('returned nodes')} · {pop.cov.graphNodesOmitted} {t('omitted nodes')} · {pop.model.edges.length} {t('returned connections')} · {pop.cov.graphEdgesOmitted} {t('omitted connections')}.</span>{:else}{t('Coverage metadata unavailable.')}{/if}</p>
 {/each}
 <p data-coverage="evidence"><strong>{t('Received-record list')}</strong>: {#if coverage?.evidence}{snapshot?.evidence.length??0} {t('returned records')} · {t('Record limit')} {coverage.evidence.limit} · {coverage.evidence.omitted} {t('omitted records')}.{:else}{t('Coverage metadata unavailable.')}{/if}</p>
 <p>{t('Graph and received-record list omissions do not reduce ordered-journey denominators. Narrow filters or inspect matching journeys for retained evidence beyond the display.')}</p>
</Disclosure>
 {#if !hideWarnings}{@render warnings()}{/if}
</section>
{/if}
<style>
 .coverage-disclosure{min-width:0;font-size:12px;overflow-wrap:anywhere}.coverage-warning{color:var(--ink);border-left:2px solid #8c642c;padding-left:8px}.coverage-warnings:empty{display:none}
 h3{font:inherit;font-weight:700;margin:0 0 4px}p{margin:3px 0}

</style>
