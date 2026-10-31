<script lang="ts">
 import type {Snapshot} from './types';
 import {translate,type Language} from './i18n';
 let {snapshot,language,live=false,baseline=false}:{snapshot:Snapshot|null;language:Language;live?:boolean;baseline?:boolean}=$props();
 const t=(s:string)=>translate(s,language);
 const coverage=$derived(snapshot?.coverage);
</script>
<section class="coverage-disclosure" aria-label={t('Coverage and display limits')}>
 <h3>{t('Coverage and display limits')}</h3>
 <p data-coverage="source"><strong>{t('Source-record scan')}</strong>: {#if coverage?.source}{t('Record limit')} {coverage.source.recordLimit} · {t(coverage.source.truncated?'Source scan capped. Counts describe scanned retained evidence; narrow filters.':'Source scan not capped. Retention and selection still apply.')}{:else}{t('Coverage metadata unavailable.')}{/if} {t('Exporter completeness is unknown; sampling and delivery loss are not measured.')}</p>
 {#each [{id:'current',label:live?'Live observation':'Current selection',cov:coverage?.current,model:snapshot?.model},...(live&&baseline?[{id:'baseline',label:'App-wide baseline',cov:coverage?.baseline,model:snapshot?.baselineModel}]:[])] as pop}
  <p data-coverage={'population-'+pop.id}><strong>{t(pop.label)}</strong>: {#if pop.cov&&pop.model}{pop.cov.journeys} {t('Retained journeys in this selection')} · {t(pop.cov.sourcePartial?'Observed partial evidence':'Retained source')} · {pop.cov.partialJourneys} {t('Partial beginnings')} · {pop.cov.censoredJourneys} {t('Censored continuations')}. <span>{t('Graph projection')}: {pop.model.nodes.length} {t('returned nodes')} · {pop.cov.graphNodesOmitted} {t('omitted nodes')} · {pop.model.edges.length} {t('returned connections')} · {pop.cov.graphEdgesOmitted} {t('omitted connections')}.</span>{:else}{t('Coverage metadata unavailable.')}{/if}</p>
 {/each}
 <p data-coverage="evidence"><strong>{t('Received-record list')}</strong>: {#if coverage?.evidence}{snapshot?.evidence.length??0} {t('returned records')} · {t('Record limit')} {coverage.evidence.limit} · {coverage.evidence.omitted} {t('omitted records')}.{:else}{t('Coverage metadata unavailable.')}{/if}</p>
 <p>{t('Graph and received-record list omissions do not reduce ordered-journey denominators. Narrow filters or inspect matching journeys for retained evidence beyond the display.')}</p>
</section>
<style>
 .coverage-disclosure{padding:8px 24px;font-size:11px;border-top:1px solid var(--line);overflow-wrap:anywhere}
 h3{font:inherit;font-weight:700;margin:0 0 4px}p{margin:3px 0}
 @media(max-width:600px){.coverage-disclosure{padding-inline:14px}}
</style>
