<script lang="ts">
 import {defaultStoryState,type StoryState} from './investigationNavigation';
 import type {Journey,Model,Evidence} from './types';
 import {groupOccurrences} from './investigation';
 import {dateTime} from './helpers';
 import {translate,type Language} from './i18n';
 let {journey,model,language,group=false,ontrace,transient=$bindable(defaultStoryState())}:{journey:Journey;model:Model;language:Language;group?:boolean;transient?:StoryState;ontrace?:(trace:string,journey:string)=>void}=$props();
 function toggle(key:string,open:boolean){if(open&&!transient.expanded.includes(key))transient.expanded=[...transient.expanded,key];else if(!open&&transient.expanded.includes(key))transient.expanded=transient.expanded.filter(k=>k!==key)}
 const t=(s:string)=>translate(s,language);
 const groups=$derived(groupOccurrences(journey.steps??[],group));
 const diagnostics=$derived(model.diagnostics.filter(d=>d.journeyId===journey.id));
 const evidence=(id:number)=>journey.evidence?.find(e=>e.id===id);
 function clock(time:number|string|undefined){try{if(time===undefined||time==='')return null;const ns=BigInt(time);return ns>0n?ns:null}catch{return null}}
 function exact(time:number|string|undefined){return clock(time)===null?t('Unavailable'):dateTime(time)}
 function relative(fromIndex:number,toIndex:number){const first=journey.steps?.[fromIndex],last=journey.steps?.[toIndex];if(!first||!last||evidence(first.evidenceId)?.fallback||evidence(last.evidenceId)?.fallback)return null;const start=clock(first.time),end=clock(last.time);if(start===null||end===null||end<start)return null;const seconds=Number(end-start)/1e9;return Number.isFinite(seconds)?seconds.toFixed(3)+' s':null}
 function elapsed(index:number){return relative(0,index)??t('Unavailable')}
 function delta(i:number){return relative(i-1,i)??t('Unavailable')}
 const uncertainClock=$derived((journey.steps??[]).some((_,i)=>relative(0,i)===null||(i>0&&relative(i-1,i)===null)));
</script>
{#snippet detail(e:Evidence|undefined,location:string)}
 {#if e}<p>{e.name} · {e.source} / {e.signal} · {e.service||t('unknown')} · {e.version||t('unknown')}</p><dl class="evidence-clocks"><dt>{t('Occurrence time')}</dt><dd>{exact(e.time)} <code>{e.time}</code></dd><dt>{t('Observed time')}</dt><dd>{exact(e.observed)} <code>{e.observed}</code></dd><dt>{t('Receipt time')}</dt><dd>{exact(e.receipt)} <code>{e.receipt}</code></dd><dt>{t('Source sequence')}</dt><dd>{e.sequence||t('Unavailable')}</dd></dl><pre>{JSON.stringify(e.attrs??{},null,2)}</pre>{#if e.trace&&ontrace}<button id={'trace-'+journey.id+'-'+e.id+'-'+location} class="secondary" onclick={()=>ontrace?.(e.trace,journey.id)}>{t('Open retained trace')}</button>{/if}{:else}<p>{t('Evidence unavailable')}</p>{/if}
{/snippet}
<article class="journey-story">
 <h3>{t('Journey')} · {exact(journey.start)} · {t(journey.status==='completed'?'Completed':journey.status==='active'?'Active':journey.status==='inactive'?'Inactive':'Exited')}</h3>
 <p class="mono">{journey.id}</p>
 <p>{journey.navigationOnly?t('Goal-free navigation'):t('Captured outcome')} · {journey.steps?.length??0} {t('Stage occurrences')} · {t('Observed span')}: {elapsed((journey.steps?.length??0)-1)}</p>
 <p>{t('Captured boundary')}: {exact(journey.start)} → {exact(clock(journey.end)!==null?journey.end:journey.lastMapped)}</p>
 {#if clock(journey.windowFrom)!==null||clock(journey.windowTo)!==null}<p>{t('Evidence window')}: {exact(journey.windowFrom)} → {exact(journey.windowTo)}</p>{/if}
 <p>{t(journey.partial?'Elapsed origin is the first available occurrence, not the captured beginning.':'Elapsed since first available occurrence; not page attention or stage duration.')}</p>
 {#if journey.recovered}<p>{t('Recovered')}</p>{/if}
 {#if uncertainClock}<p class="warning">{t('Relative timing unavailable: missing, fallback or inconsistent occurrence clocks.')}</p>{/if}
 {#each journey.warnings??[] as warning}<p class="warning">{t(warning)}</p>{/each}
 {#if journey.rightCensored}<p class="warning">{t('Continuation is outside this view.')}</p>{/if}
 <ol class="journey-timeline vertical-timeline">
 {#each groups as g (g.start)}<li>{#if g.gap}<strong>{t('Missing intermediate steps')}</strong>{/if}<strong>{g.stage}</strong>{#if g.occurrences.length>1} ×{g.occurrences.length} {t('Occurrences')}{/if}<p>+{elapsed(g.start)}{#if g.start>0} · Δ {delta(g.start)}{/if} · {exact(g.occurrences[0].time)}</p>
  {#each diagnostics.filter(d=>g.occurrences.some(s=>s.evidenceId===d.evidenceId)) as d}<details open={transient.expanded.includes('diagnostic:'+d.evidenceId)} ontoggle={e=>toggle('diagnostic:'+d.evidenceId,e.currentTarget.open)}><summary>! {d.type} · {d.severity}</summary>{@render detail(evidence(d.evidenceId),'diagnostic')}</details>{/each}
  <details open={transient.expanded.includes('evidence:'+g.start)} ontoggle={e=>toggle('evidence:'+g.start,e.currentTarget.open)}><summary>{t('Evidence details')} · {g.start+1}{g.end>g.start?'–'+(g.end+1):''}</summary>{#each g.occurrences as step,i}<section><h4>{g.start+i+1}. {step.step} · {exact(step.time)}</h4>{@render detail(evidence(step.evidenceId),'occurrence')}</section>{/each}</details>
 </li>{/each}</ol>
 {#if !groups.length}<p>{t('No retained journey detail in this selection.')}</p>{/if}
 <h3>{t('Linked diagnostic evidence')}</h3><p>{t('Unplaced diagnostics: ordered evidence does not establish an exact stage attachment or causality.')}</p>
 {#each diagnostics.filter(d=>!journey.steps.some(s=>s.evidenceId===d.evidenceId)) as d}<details open={transient.expanded.includes('diagnostic:'+d.evidenceId)} ontoggle={e=>toggle('diagnostic:'+d.evidenceId,e.currentTarget.open)}><summary>! {d.type} · {d.severity} · {t('Evidence')}: {d.evidenceId}</summary>{@render detail(evidence(d.evidenceId),'diagnostic')}</details>{:else}<p>{t('No unplaced diagnostic evidence')}</p>{/each}
</article>
