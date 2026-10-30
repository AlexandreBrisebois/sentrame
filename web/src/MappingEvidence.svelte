<script lang="ts">
 import {translate,type Language} from './i18n';
 import {dateTime} from './helpers';
 import type {MappingEvidencePreview} from './types';
 let {preview,loading,error,language,hasCondition,oncondition}:{preview:MappingEvidencePreview|null;loading:boolean;error:string;language:Language;hasCondition:boolean;oncondition:(key:string,value:string)=>void}=$props();
 const t=(s:string)=>translate(s,language);
 const examples=$derived(preview?(preview.matches?preview.matchingExamples:preview.examples):[]);
</script>
<section class="mapping-evidence" aria-label={t('Received records for this mapping')} aria-busy={loading}>
 <h3>{t('Received records for this mapping')}</h3>
 <p>{t('Examples use retained app-wide history for the selected event, service, carrier and producer, independent of the graph time range.')}</p>
 {#if loading}<p role="status">{t('Loading received records…')}</p>
 {:else if error}<p class="error" role="alert">{error}</p>
 {:else if preview}
 {#if preview.truncated}<p class="warning">{t('Results are partial. Some retained records were omitted.')}</p>{/if}
 <p class="match-count" role="status"><strong>{preview.matches} {t('Matching records')}</strong> · {preview.candidates} {t('before the optional condition')}</p>
 {#if preview.matches&&preview.missingContext}<p class="warning">{preview.missingContext} {t('matching records are missing a session ID or, for browser events, a tab ID. You can save this mapping, but these records cannot be connected into journeys.')}</p>{:else if preview.matches}<p>{t('All matching records have the IDs needed to connect them into journeys.')}</p>{/if}
 {#if !preview.candidates}<p class="warning">{t('No retained records match this event context. Choose a received context above or check the exact service, carrier and producer.')}</p>
 {:else}
  {#if !preview.matches}<p class="warning">{t('This condition matches no records. The examples and available values below show this event context before applying the condition.')}</p>{/if}
  <h4>{t('Pick a condition from received data (optional)')}</h4>
  <p>{t('Choose a field and value below, or leave both condition fields empty to include every record in this event context.')}</p>
  <div class="condition-options">{#each preview.conditions as condition}<button type="button" class="secondary" onclick={()=>oncondition(condition.key,condition.value)}><span><code>{condition.key}</code><strong>{condition.value}</strong></span><small>{condition.count} {t('records')}</small></button>{:else}<p>{t('No selectable attributes were retained for these records. Leave the condition empty.')}</p>{/each}</div>
  {#if hasCondition}<button type="button" class="text-button" onclick={()=>oncondition('','')}>{t('Clear condition')}</button>{/if}
  {#if preview.conditionsTruncated}<p>{t('Available values are bounded. Some values are omitted.')}</p>{/if}
  <details open><summary>{t(preview.matches?'Matching examples':'Examples before the condition')} · {examples.length}/{preview.matches||preview.candidates}</summary>
   {#each examples as record}<article><p><strong>{record.service||t('unknown')}</strong> · {record.source} / {record.signal}<br/>{dateTime(record.time)} · {record.version||t('unknown')}</p>
    <dl>{#each Object.entries(record.attrs??{}) as [key,value]}<div><dt>{key}</dt><dd>{value}</dd></div>{:else}<div><dt>{t('Attributes')}</dt><dd>{t('No retained attributes')}</dd></div>{/each}</dl>
   </article>{/each}
  </details>
 {/if}
 {:else}<p>{t('No received examples are available.')}</p>{/if}
</section>
<style>
 .mapping-evidence{border:1px solid var(--line);border-radius:7px;background:#f4f7ef;padding:15px;min-width:0}.mapping-evidence h3{font-size:13px;margin:0 0 8px}.mapping-evidence h4{font-size:12px;margin:16px 0 8px}.mapping-evidence p{font-size:11px;line-height:1.7;color:var(--muted);margin:0 0 10px}.mapping-evidence .match-count{color:var(--ink)}.condition-options{display:grid;gap:6px;max-height:220px;overflow:auto;padding:3px;margin:-3px}.condition-options button{text-align:left;justify-content:space-between;padding:9px 11px;min-width:0;font-size:11px;background:var(--surface)}.condition-options button>span{min-width:0}.condition-options code{font-size:10px;display:block;overflow-wrap:anywhere}.condition-options strong{display:block;font-weight:600;overflow-wrap:anywhere}.condition-options small{flex:none;color:var(--muted);font-size:10px}.mapping-evidence>.text-button{font-size:11px;margin-top:8px}.mapping-evidence details{margin-top:13px;border-top:1px solid var(--line)}.mapping-evidence summary{font-size:12px}article{border-top:1px solid var(--line);padding:12px 0 0;margin-top:6px}article p{overflow-wrap:anywhere}dl{font-size:11px;margin:0 0 12px}dl>div{display:grid;grid-template-columns:minmax(0,1fr) minmax(0,1.3fr);gap:10px;padding:4px 0}dt{color:var(--muted);overflow-wrap:anywhere}dd{margin:0;overflow-wrap:anywhere}
</style>
