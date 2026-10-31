<script lang="ts">
 import {onMount,tick} from 'svelte';
 import type {Model} from './types';
 import type {Criterion} from './investigation';
 import {matchingExamples} from './investigation';
 import {defaultExampleState,defaultStoryState,type ExampleState} from './investigationNavigation';
 import {dateTime} from './helpers';
 import {translate,type Language} from './i18n';
 import JourneyStory from './JourneyStory.svelte';
 let {model,criterion,language,population,group=false,ontrace,transient=$bindable(defaultExampleState())}:{model:Model;criterion:Criterion;language:Language;population:string;group?:boolean;ontrace?:(trace:string,journey:string)=>void;transient?:ExampleState}=$props();
 const t=(s:string)=>translate(s,language);
 let root:HTMLDivElement;
 const result=$derived(matchingExamples(model.journeys,criterion,transient.search,transient.status,transient.from,transient.page));
 const selected=$derived(model.journeys.find(j=>j.id===transient.chosen));
 const id=(suffix:string)=>'examples-'+encodeURIComponent(population)+'-'+suffix;
 const scroller=()=>root?.closest('dialog');
 function focus(target:string){(document.getElementById(target)??document.getElementById(id('search')))?.focus({preventScroll:true})}
 async function choose(journey:string,invoker:string){
  transient.listScroll=scroller()?.scrollTop??0;transient.returnFocus=invoker;transient.chosen=journey;transient.story=defaultStoryState();transient.scroll=0;
  await tick();focus(id('back'));scroller()?.scrollTo({top:0});
 }
 async function back(){transient.chosen='';await tick();focus(transient.returnFocus);scroller()?.scrollTo({top:transient.listScroll});transient.scroll=transient.listScroll}
 onMount(()=>{
  const dialog=scroller();
  const remember=()=>{transient.scroll=dialog?.scrollTop??0};
  dialog?.addEventListener('scroll',remember);
  void tick().then(()=>{if(!root.isConnected)return;if(transient.focus)focus(transient.focus);dialog?.scrollTo({top:transient.scroll})});
  return()=>dialog?.removeEventListener('scroll',remember);
 });
</script>
<div bind:this={root} onfocusin={e=>{if(e.target instanceof HTMLElement&&e.target.id)transient.focus=e.target.id}}>
<p><strong>{t(population)}</strong> · {t(criterion.kind==='all'?'All retained journeys':criterion.kind)}{criterion.kind==='all'?'':': '+('stage' in criterion?criterion.stage:'from' in criterion?criterion.from+' → '+criterion.to:criterion.path.join(' → '))}</p>
{#if transient.chosen}
 <button id={id('back')} class="secondary" onclick={back}>{t('Back to matching journeys')}</button>
 {#if selected}<JourneyStory journey={selected} {model} {language} {group} {ontrace} bind:transient={transient.story}/>{:else}<p role="status">{t('Evidence unavailable')}</p>{/if}
{:else}
 <p>{t('Newest captured start first; journey ID breaks ties. Examples are not representative samples.')}</p>
 <div class="investigation-controls">
  <label>{t('Path search')}<input id={id('search')} type="search" bind:value={transient.search} oninput={()=>transient.page=0}/></label>
  <label>{t('Status')}<select id={id('status')} bind:value={transient.status} onchange={()=>transient.page=0}><option value="">{t('All')}</option>{#each ['active','completed','exited','inactive'] as s}<option value={s}>{t(s==='active'?'Active':s==='completed'?'Completed':s==='exited'?'Exited':'Inactive')}</option>{/each}</select></label>
  <label>{t('Started after')}<input id={id('from')} type="datetime-local" bind:value={transient.from} onchange={()=>transient.page=0}/></label>
 </div>
 <p>{result.total} {t('Matching journeys')} · {Math.min(result.start+20,result.total)} / {result.total}</p>
 <ol class="matching-examples">{#each result.rows as j (j.id)}<li>
  <button id={id(j.id)} class="text-button" onclick={e=>choose(j.id,e.currentTarget.id)}>{dateTime(j.start)} · {t(j.status==='active'?'Active':j.status==='completed'?'Completed':j.status==='inactive'?'Inactive':'Exited')} · {j.path.slice(0,8).join(' → ')}{j.path.length>8?' …':''}</button>
  {#if j.partial||j.rightCensored}<small>{t('Observed partial evidence')}</small>{/if}<details><summary>{t('Journey ID')}</summary><code>{j.id}</code></details>
 </li>{:else}<li>{t('No matching retained journeys. Narrowing and unavailable detail can limit examples.')}</li>{/each}</ol>
 <div class="button-row"><button id={id('previous')} class="secondary" disabled={result.start===0} onclick={()=>transient.page--}>{t('Previous')}</button><button id={id('next')} class="secondary" disabled={result.start+20>=result.total} onclick={()=>transient.page++}>{t('Next')}</button></div>
{/if}
</div>
