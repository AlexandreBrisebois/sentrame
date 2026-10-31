<script lang="ts">
 import type {Journey,Model} from './types';
 import {loopProjection} from './journeyViews';
 import {dateTime} from './helpers';
 import {edgeKey} from './graph';
 import {translate,type Language} from './i18n';
 let {view,model,language,selectedStep,selectedEdge,onedge,onjourney,partial=false,label,chosen=$bindable('')}:{view:string;model:Model;language:Language;selectedStep:string;selectedEdge:string;onedge:(a:string,b:string)=>void;onjourney:(id:string)=>void;partial?:boolean;label:string;chosen?:string}=$props();
 const t=(s:string)=>translate(s,language);
 let limited=$derived(partial||model.journeys.some(j=>j.partial||j.rightCensored||(j.steps??[]).some(s=>s.gapBefore)));
 let projection=$derived(loopProjection(model.journeys));
 let allStages=$derived([...new Set(model.journeys.flatMap(j=>(j.steps??[]).map(s=>s.step)))].sort());
 let stages=$derived(allStages.slice(0,100));
 let edgeMap=$derived(new Map(projection.edges.map(e=>[edgeKey(e.from,e.to),e])));
 let selected=$derived(model.journeys.find(j=>j.id===chosen));
 function evidence(j:Journey,id:number){return j.evidence?.find(e=>e.id===id)}
 function prefixCount(path:string[]){return projection.prefixCounts.get(JSON.stringify(path))??0}
</script>
<section class="loop-view" aria-label={t(label)}>
 <h2>{t(label)}</h2>
 {#if limited}<p class="warning">{t('Coverage is partial. Counts describe available retained evidence only.')}</p>{/if}
 {#if model.journeys.some(j=>j.partial||j.rightCensored)}<p class="warning">{t('Some journeys have a partial beginning or continuation; earlier visits may be missing.')}</p>{/if}
 {#if view==='journey-timeline'}
  <label>{t('Choose a journey')}<select bind:value={chosen}><option value="">{t('Choose a journey')}</option>{#each model.journeys as j}<option value={j.id}>{j.id.slice(0,8)} · {dateTime(j.start)} · {j.path.join(' → ').slice(0,100)} · {t(j.status==='completed'?'Completed':j.status==='active'?'Active':j.status==='inactive'?'Inactive':'Exited')}</option>{/each}</select></label>
  {#if selected}
   <p>{t('Journey')}: {selected.id}</p>
   <p>{t('Captured boundary')}: {dateTime(selected.start)} → {dateTime(Number(selected.end)?selected.end:selected.lastMapped)} · {t((selected.capturedStatus||selected.status)==='completed'?'Completed':(selected.capturedStatus||selected.status)==='active'?'Active':(selected.capturedStatus||selected.status)==='inactive'?'Inactive':'Exited')}</p>
   {#if Number(selected.windowFrom)||Number(selected.windowTo)}<p>{t('Evidence window')}: {dateTime(Number(selected.windowFrom)?selected.windowFrom:undefined)} → {dateTime(Number(selected.windowTo)?selected.windowTo:undefined)}</p>{/if}
   {#each selected.warnings??[] as warning}<p class="warning">{t(warning)}</p>{/each}
   {#if selected.rightCensored}<p class="warning">{t('Continuation is outside this view.')}</p>{/if}
   <ol class="journey-timeline">{#each selected.steps??[] as step,i}{@const e=evidence(selected,step.evidenceId)}<li>{#if step.gapBefore}<strong>{t('Missing intermediate steps')}</strong> · {/if}<strong>{step.step}</strong> · {dateTime(step.time)}<details><summary>{t('Evidence details')}</summary>{#if e}<p>{e.name} · {e.source} / {e.signal} · {e.service}</p><pre>{JSON.stringify(e.attrs??{},null,2)}</pre>{:else}<p>{t('Evidence unavailable')}</p>{/if}</details></li>{/each}</ol>
   {#if !selected.steps?.length}<p>{t('No retained journey detail in this selection.')}</p>{/if}
   <h3>{t('Linked diagnostic evidence')}</h3>{#each model.diagnostics.filter(d=>d.journeyId===selected.id) as d}{@const e=evidence(selected,d.evidenceId)}<details><summary>! {d.type} · {d.severity} · {t('Evidence')}: {d.evidenceId}</summary>{#if e}<p>{e.name} · {dateTime(e.time)} · {e.source} / {e.signal} · {e.service}</p><pre>{JSON.stringify(e.attrs??{},null,2)}</pre>{:else}<p>{t('Evidence unavailable')}</p>{/if}</details>{:else}<p>{t('No linked diagnostic evidence')}</p>{/each}
   <button class="secondary" onclick={()=>onjourney(selected!.id)}>{t('Show ordered evidence')}</button>
  {:else}<p>{t('Choose an available journey to see its ordered steps. Expired detail cannot be reconstructed.')}</p>{/if}
 {:else if view==='matrix'}
  <p>{t('Rows are origins; columns are destinations. Counts show direct observed transitions.')}</p>
  {#if allStages.length>100}<p class="warning">{t('Matrix limited to the first 100 stages; narrow filters to see omitted stages.')} ({stages.length}/{allStages.length} · {allStages.length-stages.length} {t('omitted stages')})</p>{/if}
  <p>{t('0 means no observed transition in complete coverage; — means coverage is unavailable or partial.')}</p>
  <div class="table-wrap"><table><caption>{t('Occurrences / journeys / visits')}</caption><thead><tr><th>{t('From / to')}</th>{#each stages as stage}<th scope="col">{stage}</th>{/each}</tr></thead><tbody>{#each stages as from}<tr><th scope="row">{from}</th>{#each stages as to}{@const e=edgeMap.get(edgeKey(from,to))}<td>{#if e}<button class="text-button" aria-pressed={selectedEdge===edgeKey(from,to)} aria-label={`${from} → ${to}: ${e.occurrences} ${t('Occurrences')}, ${e.journeyIds.length} ${t('Journeys')}, ${e.visits} ${t('Visits')}`} onclick={()=>onedge(from,to)}>{e.occurrences} / {e.journeyIds.length} / {e.visits}</button>{:else}{limited?'—':'0'}{/if}</td>{/each}</tr>{/each}</tbody></table></div>
  {#each projection.edges.filter(e=>edgeKey(e.from,e.to)===selectedEdge) as edge}<section><h3>{edge.from} → {edge.to} · {t('Evidence')}</h3>{#each edge.journeyIds as id}<button class="text-button" onclick={()=>onjourney(id)}>{t('Journey')} {id.slice(0,8)}</button> {/each}</section>{/each}
  {#if !stages.length}<p>{t('No retained journey detail in this selection.')}</p>{/if}
 {:else if view==='repeats'}
  <p>{t('Repeats are additional occurrences in a journey. Returns are direct transitions to a different previously observed stage. Repetition is not automatically a problem.')}</p>
  <p>{t('Retained journeys in this selection')}: {model.journeys.length}</p>
  <h3>{t('Repeated stages')}</h3><table><thead><tr><th>{t('Stage')}</th><th>{t('Additional occurrences')}</th><th>{t('Journeys')}</th><th>{t('Evidence')}</th></tr></thead><tbody>{#each projection.repeats as r}<tr><th>{r.stage}{#if selectedStep===r.stage} · ✓{/if}</th><td>{r.occurrences}</td><td>{r.journeyIds.length}</td><td>{#each r.journeyIds as id}<button class="text-button" onclick={()=>onjourney(id)}>{t('Journey')} {id.slice(0,8)}</button> {/each}</td></tr>{/each}</tbody></table>
  <h3>{t('Return transitions')}</h3><table><thead><tr><th>{t('Transition')}</th><th>{t('Occurrences')}</th><th>{t('Journeys')}</th><th>{t('Evidence')}</th></tr></thead><tbody>{#each projection.returns as r}<tr><th><button class="text-button" onclick={()=>onedge(r.from,r.to)}>{r.from} → {r.to}</button></th><td>{r.occurrences}</td><td>{r.journeyIds.length}</td><td>{#each r.journeyIds as id}<button class="text-button" onclick={()=>onjourney(id)}>{t('Journey')} {id.slice(0,8)}</button> {/each}</td></tr>{/each}</tbody></table>
 {:else if view==='sequence'}
  <p>{t('Each row is an observed path prefix. Bars show journeys following each prefix; grouped or continued paths are not endings.')}</p>
  <p>{t('Maximum steps / path groups')}: {projection.maxSteps} / {projection.maxBranches}</p>
  <ol class="sequence-flow">{#each projection.prefixes as p}<li class="sequence-flow-row">{#each p.path as stage,i}{@const count=prefixCount(p.path.slice(0,i+1))}<span class="sequence-flow-step"><strong>{i+1}. {stage}</strong><small>{count} {t('Journeys')}</small><span class="sequence-bar-track" aria-hidden="true"><span class="sequence-bar" style={`width:${100*count/Math.max(1,model.journeys.length)}%`}></span></span></span>{/each}{#if p.continued}<strong>… {p.continued} {t('continue or are censored')}</strong>{/if}{#if p.partial}<span>{p.partial} {t('partial')}</span>{/if}<span>{#each p.journeyIds as id}<button class="text-button" onclick={()=>onjourney(id)}>{t('Journey')} {id.slice(0,8)}</button> {/each}</span></li>{/each}</ol>
  {#if projection.omittedBranches}<p>{t('Other path groups')}: {projection.omittedBranches} · {projection.omittedJourneys} {t('Journeys')}</p>{/if}
 {/if}
</section>
