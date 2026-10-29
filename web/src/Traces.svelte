<script lang="ts">
 import {onMount} from 'svelte';
 import type {TraceDetail,Evidence,Definition,SavedView} from './types';
 import {traceRows,operationBar} from './trace';
 import {translate,type Language} from './i18n';
 let {request,revision,language='en',readOnly=false,onconfigure}:{request:(path:string,body?:unknown)=>Promise<any>;revision:string|number;language?:Language;readOnly?:boolean;onconfigure?:(e:Evidence,suggestion?:string)=>void}=$props();
 const t=(s:string)=>translate(s,language);
 let list=$state<{trace:string;spans:number;records:number;services:string[]}[]>([]), selected=$state(''),detail=$state<TraceDetail|null>(null),error=$state(''),loading=$state(false),view=$state('tree'),truncated=$state(false);
 let definitions=$state<Definition[]>([]),previewDefinition=$state('');
 let references=$state<SavedView[]>([]),selectedReference=$state(''),referenceName=$state('');
 let generation=0;
 let autoRefresh=$state(true);
 onMount(()=>{const stream=new EventSource('/api/events');stream.addEventListener('reset',()=>{if(autoRefresh)void load(selected)});stream.onerror=()=>{if(autoRefresh)void load(selected)};return()=>{generation++;stream.close()}});
 async function load(id:string){const current=++generation;loading=true;error='';try{const response=await request('traces');if(current!==generation)return;list=response.traces;truncated=response.truncated;const defs=await request('definitions');if(current!==generation)return;definitions=defs;const refs=await request('views');if(current!==generation)return;references=refs.filter((v:SavedView)=>v.mode==='trace');if(id){const d=await request('traces/'+encodeURIComponent(id)+(selectedReference?'?reference='+selectedReference:previewDefinition?'?previewDefinition='+previewDefinition:''));if(current===generation)detail=d}else detail=null}catch(e){if(current===generation){detail=null;error=e instanceof Error?e.message:t('Trace unavailable')}}finally{if(current===generation)loading=false}}
 $effect(()=>{revision;previewDefinition;selectedReference;void load(selected);return()=>{generation++}});
 function openReference(id:string){const ref=references.find(v=>String(v.id)===id);selectedReference=id;if(ref?.trace){selected=ref.trace;previewDefinition=String(ref.definition||'');if(ref.layout&&['tree','timeline','table'].includes(ref.layout))view=ref.layout}}
 async function saveReference(){if(readOnly||!selected||!referenceName.trim())return;try{const ref=await request('views',{name:referenceName.trim(),mode:'trace',trace:selected,layout:view,definition:Number(previewDefinition||0)});referenceName='';await load(selected);openReference(String(ref.id))}catch(e){error=e instanceof Error?e.message:t('Trace unavailable')}}
 const rows=$derived(detail?traceRows(detail.operations,detail.relationships):[]);
 async function copy(){if(detail){try{await navigator.clipboard.writeText(JSON.stringify(detail,null,2))}catch{error=t('Copy unavailable. Select the displayed sanitized text.')}}}
</script>
<section class="trace-workspace" aria-label={t('Trace exploration')}>
 <header><h2>{t('Trace exploration')}</h2><label class="toggle"><input type="checkbox" bind:checked={autoRefresh}/>{t('Automatic refresh')}</label><button onclick={()=>load(selected)} disabled={loading}>{t('Refresh')}</button></header>
 <p>{t('Observed operations and relationships. A trace is not a business journey; inclusive durations are not page dwell or unique total time.')}</p>
 <label>{t('Observed trace')}<select bind:value={selected} onchange={()=>selectedReference=''}><option value="">{t('Choose a trace')}</option>{#each list as item}<option value={item.trace}>{item.services.join(', ')} · {item.trace.slice(0,12)} · {item.spans} {t('spans')}</option>{/each}</select></label>
 <label>{t('Trace references')}<select value={selectedReference} onchange={e=>openReference(e.currentTarget.value)}><option value="">{t('No reference')}</option>{#each references as ref}<option value={String(ref.id)}>{ref.name}</option>{/each}</select></label>
 {#if loading}<p role="status">{t('Loading…')}</p>{/if}{#if error}<p role="alert">{error}</p>{/if}
 {#if truncated}<p>{t('Bounded results: narrow collection or inspect a recent trace. Some evidence is omitted.')}</p>{/if}
 {#if !loading&&!list.length}<p>{t('No retained traces. Export ordinary spans with valid trace/span IDs; no stage or session attributes are required.')}</p>{/if}
 {#if detail}
 <label>{t('Business interpretation')}<select bind:value={previewDefinition} onchange={()=>selectedReference=''}><option value="">{t('Capture-time definitions')}</option>{#each definitions as def}<option value={String(def.version)}>{t('Non-persisting preview')} v{def.version}</option>{/each}</select></label>
 <p>{t(detail.interpretation)}</p>
 {#if !readOnly}<form onsubmit={e=>{e.preventDefault();void saveReference()}}><label>{t('Reference name')}<input bind:value={referenceName} maxlength="160" required/></label><button disabled={!referenceName.trim()||loading}>{t('Save trace reference')}</button></form>{/if}
 <p class="mono">{detail.trace} · {t('Revision')} {detail.revision}</p>
 <p>{t('Completeness unknown: sampling, missing exports, retention and bounds can omit work. A sampled flag does not establish a sampling rate.')}</p>
 {#if detail.truncated}<p>{t('Trace detail truncated at 2,000 records.')}</p>{/if}
 <div class="button-row" aria-label={t('Trace view')}>{#each ['tree','timeline','table'] as option}<button aria-pressed={view===option} onclick={()=>view=option}>{t(option)}</button>{/each}<button onclick={copy}>{t('Copy sanitized trace')}</button></div>
 <div class="table-wrap"><table><caption>{t('Native inclusive span timing')} · {detail.computation}</caption><thead><tr><th>{t('Operation')}</th><th>{t('Producer / service')}</th><th>{t('Start / end (exact ns)')}</th><th>{t('Duration (ns) / quality')}</th><th>{t('Partial evidence')}</th></tr></thead><tbody>
 {#each rows as {operation:op,depth}}
 {@const e=op.evidence}{@const bar=operationBar(op,detail.operations)}
 <tr><td><span style:padding-left={view==='tree'?`${Math.min(depth,16)}rem`:'0'}>{view==='tree'&&depth?'↳ ':''}{e.name}</span>{#each detail.business.filter(b=>b.evidenceId===e.id) as b}<p>{b.stage} · v{b.definition} · {b.basis}{b.goal?' · '+t('Configured goal at span end'):''}</p>{/each}{#if !readOnly&&onconfigure}<button class="text-button" onclick={()=>onconfigure?.(e)}>{t('Give business meaning')}</button>{#each detail.suggestions.filter(s=>s.evidenceId===e.id) as suggestion}<button class="text-button" onclick={()=>onconfigure?.(e,suggestion.label)}>{t('Use emitted suggestion')}: {suggestion.label}</button>{/each}{/if}{#if e.attrs?.['stage.name']}<small>{t('Emitted hint')}: {e.attrs['stage.name']}</small>{/if}<small class="mono">{e.span}</small>{#if view==='timeline'&&bar}<div class="trace-track"><span style:left={`${bar.left}%`} style:width={`${Math.max(bar.width,0.2)}%`}></span></div>{/if}</td><td>{e.source} / {e.service||t('unknown')}</td><td class="mono">{e.time}<br/>{e.end==='0'?t('Unavailable'):e.end}</td><td class="mono">{op.durationNs??'—'}<br/>{t(op.timing)}</td><td>{t('SDK dropped attributes / events / links')}: {e.droppedAttributes??0}/{e.droppedEvents??0}/{e.droppedLinks??0}<br/>{t('Projection omitted links')}: {e.omittedLinks??0}<br/>{t('Projection omitted attributes')}: {e.omittedAttributes??0}<br/>{t('Name projection')}: {e.nameState??'unknown'}<br/>{t('Flags')}: {e.flags??0}</td></tr>
 {/each}
 </tbody></table></div>
 <details><summary>{t('Typed relationships')}</summary><div class="table-wrap"><table><thead><tr><th>{t('Relationship')}</th><th>{t('From record')}</th><th>{t('Target trace / span')}</th><th>{t('State')}</th></tr></thead><tbody>{#each detail.relationships as r}<tr><td>{r.type}</td><td>{r.from}</td><td class="mono">{r.targetTrace}<br/>{r.targetSpan||'—'}</td><td>{t(r.state)}</td></tr>{/each}</tbody></table></div></details>
 <details><summary>{t('Span events in export order')}</summary>{#each detail.events as e}<p><span class="mono">{e.span} #{e.eventIndex} · {e.time}</span> · {e.name} · {e.source}{#each detail.business.filter(b=>b.evidenceId===e.id) as b} · {b.stage} · v{b.definition}{b.goal?' · '+t('Configured event goal'):''}{/each}{#if !readOnly&&onconfigure}<button class="text-button" onclick={()=>onconfigure?.(e)}>{t('Give business meaning')}</button>{#each detail.suggestions.filter(s=>s.evidenceId===e.id) as suggestion}<button class="text-button" onclick={()=>onconfigure?.(e,suggestion.label)}>{t('Use emitted suggestion')}: {suggestion.label}</button>{/each}{/if}</p>{/each}</details>
 <details><summary>{t('Correlated logs')}</summary>{#each detail.logs as e}<p><span class="mono">{e.time} · {e.span||t('trace-only')}</span> · {e.name} · {e.source}{#each detail.business.filter(b=>b.evidenceId===e.id) as b} · {b.stage} · v{b.definition}{b.goal?' · '+t('Configured event goal'):''}{/each}{#if !readOnly&&onconfigure}<button class="text-button" onclick={()=>onconfigure?.(e)}>{t('Give business meaning')}</button>{#each detail.suggestions.filter(s=>s.evidenceId===e.id) as suggestion}<button class="text-button" onclick={()=>onconfigure?.(e,suggestion.label)}>{t('Use emitted suggestion')}: {suggestion.label}</button>{/each}{/if}</p>{/each}</details>
 {/if}
</section>
<style>
 .trace-workspace{padding:24px 26px;overflow:auto;min-height:0;display:flex;flex-direction:column;gap:18px}
 .trace-workspace>*{flex-shrink:0}
 .trace-workspace>p{margin:0;max-width:90ch;overflow-wrap:anywhere;font-size:13px}
 .trace-workspace>label{max-width:720px;font-size:12px}
 .trace-workspace header{display:flex;flex-wrap:wrap;gap:12px 24px;align-items:center}
 .trace-workspace header h2{margin-right:auto}
 .trace-workspace header .toggle{min-height:42px;font-size:12px}
 .trace-workspace form{display:flex;align-items:end;flex-wrap:wrap;gap:12px;max-width:720px}
 .trace-workspace form label{flex:1 1 240px;min-width:0;font-size:12px}
 .trace-workspace .button-row{margin:0;gap:8px 12px}
 .trace-workspace details{border-top:1px solid var(--line);padding-top:8px}
 .trace-workspace details p{margin:12px 0;overflow-wrap:anywhere;font-size:12px}
 td{vertical-align:top}td small{display:block}td p{margin:8px 0}
 td :global(.text-button),details :global(.text-button){display:flex;width:fit-content;max-width:100%;margin-top:8px;min-height:36px}
 .trace-track{position:relative;height:1rem;width:16rem;background:var(--surface,#eee);margin-top:8px}.trace-track span{position:absolute;height:100%;background:#547345}
 table{min-width:44rem}th,td{padding:12px;text-align:left}
 @media(max-width:760px){.trace-workspace{padding:18px;gap:18px}.trace-workspace header h2{flex-basis:100%;font-size:25px}.trace-workspace form button{width:100%}}
</style>
