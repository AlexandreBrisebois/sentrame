<script lang="ts">
 import {translate,type Language} from './i18n';
 import type {RegisteredApp} from './apps';
 let {apps,selected,language,request,onselect,onclose}:{apps:RegisteredApp[];selected:string;language:Language;request:(path:string,body?:unknown)=>Promise<any>;onselect:(id:string,settings?:boolean)=>void;onclose:()=>void}=$props();
 const t=(s:string)=>translate(s,language);
 let name=$state(''),rename=$state(''),confirmation=$state(''),rotation=$state(''),error=$state(''),notice=$state(''),busy=$state(false);
 let secrets=$state<{browser:string;server:string;endpoint:string}|null>(null),showBrowser=$state(false),showServer=$state(false);
 const current=$derived(apps.find(a=>a.id===selected));
 let loadedID='';$effect(()=>{if(current&&current.id!==loadedID){rename=current.name;loadedID=current.id}});
 async function action(fn:()=>Promise<void>){busy=true;error='';notice='';try{await fn()}catch(e){error=e instanceof Error?t(e.message):t('Request failed. Try again.')}finally{busy=false}}
 async function register(e:SubmitEvent){e.preventDefault();await action(async()=>{const a=await request('apps',{name:name.trim()});onselect(a.id,true)})}
 async function reveal(role:string){await action(async()=>{secrets=await request('apps/'+selected+'/connection');if(role==='browser')showBrowser=true;else showServer=true})}
 async function copy(role:string){await action(async()=>{if(!secrets)secrets=await request('apps/'+selected+'/connection');if(!secrets)return;await navigator.clipboard.writeText('Authorization: Bearer '+(role==='browser'?secrets.browser:secrets.server));notice=t('Copied')})}
 async function rotate(){await action(async()=>{await request('apps/'+selected+'/rotate',{role:rotation});secrets=null;showBrowser=false;showServer=false;rotation='';notice=t('Credential rotated. Update the sender configuration.')})}
 async function remove(){await action(async()=>{await request('apps/'+selected+'/remove',{confirmation});onselect('')})}
</script>
<section aria-label={t('App settings')}>
 {#if current}
 <h2>{current.name}</h2>
 <form onsubmit={e=>{e.preventDefault();void action(async()=>{await request('apps/'+selected+'/rename',{name:rename});onselect(selected)})}}><label>{t('App name')}<input bind:value={rename} required maxlength="80"/></label><button disabled={busy||!rename.trim()}>{t('Rename app')}</button></form>
 <h3>{t('Connection settings')}</h3>
 <p>{t('Configure your existing standard OTLP/HTTP exporter. No Sentrame SDK, app-ID attribute or instrumentation changes are required.')}</p>
 <p>{t('Use the reachable base URL for your sender. Signal URLs end in /v1/logs, /v1/traces or /v1/metrics. A base endpoint exporter adds these paths; a signal endpoint uses the full path.')}</p>
 <code>{window.location.origin}</code><button class="secondary" onclick={()=>action(async()=>{await navigator.clipboard.writeText(window.location.origin);notice=t('Copied')})}>{t('Copy base endpoint')}</button>
 {#each ['logs','traces','metrics'] as signal}<div><code>{window.location.origin}/v1/{signal}</code><button class="text-button" onclick={()=>action(async()=>{await navigator.clipboard.writeText(window.location.origin+'/v1/'+signal);notice=t('Copied')})}>{t('Copy signal endpoint')} · {signal}</button></div>{/each}
 <p>{t('Server credentials belong only in server secret configuration. Browser credentials are untrusted; use the supported app-origin forwarder. Direct cross-origin browser export is not qualified.')}</p>
 {#each ['browser','server'] as role}
 <fieldset><legend>{t(role==='browser'?'Browser':'Server')}</legend><code class="credential">{role==='browser'&&showBrowser?secrets?.browser:role==='server'&&showServer?secrets?.server:'••••••••••••••••'}</code>
 <div class="button-row"><button disabled={busy} onclick={()=>reveal(role)}>{t('Reveal credential')}</button><button disabled={busy} onclick={()=>copy(role)}>{t('Copy Authorization header')}</button><button class="secondary" disabled={busy} onclick={()=>rotation=role}>{t('Rotate credential')}</button></div></fieldset>
 {/each}
 {#if rotation}<div role="group" aria-label={t('Confirm rotation')}><p>{t('Rotation immediately rejects the old credential. Update every sender using this role afterward.')} {t(rotation==='browser'?'Browser':'Server')}</p><button disabled={busy} onclick={rotate}>{t('Confirm rotation')}</button><button class="secondary" onclick={()=>rotation=''}>{t('Cancel')}</button></div>{/if}
 <details><summary>{t('Remove app')}</summary><p>{t('This permanently removes this app and its telemetry, mappings, references and viewer access. Type the exact app name to confirm.')}</p><label>{t('Confirmation')}<input bind:value={confirmation} autocomplete="off"/></label><button disabled={busy||confirmation!==current.name} onclick={remove}>{t('Remove app')}</button></details>
 {:else}<h2>{t('Add app')}</h2><p>{t('Register your app, then connect its existing OpenTelemetry exporter using the generated settings.')}</p>{/if}
 {#each apps.filter(a=>a.state==='deleting') as pending}<div role="status"><strong>{pending.name}</strong><p>{t('Removal pending. Retry removal or restart to resume.')}</p><button disabled={busy} onclick={()=>action(async()=>{await request('apps/'+pending.id+'/remove',{confirmation:pending.name});onselect('')})}>{t('Retry removal')}</button></div>{/each}
 <form onsubmit={register}><h3>{t('Add app')}</h3><label>{t('App name')}<input bind:value={name} required maxlength="80" autocomplete="off"/></label><p>{t('Use a unique name of 1–80 characters. Up to 32 apps per installation.')}</p><button disabled={busy||!name.trim()}>{t('Register app')}</button></form>
 {#if error}<p class="error" role="alert">{error}</p>{/if}{#if notice}<p role="status">{notice}</p>{/if}
 <button class="secondary" onclick={()=>{secrets=null;onclose()}}>{t('Close')}</button>
</section>
