<script lang="ts">
 import Disclosure from './Disclosure.svelte';
 import {translate,type Language} from './i18n';
 import type {Lens} from './graph';
 let {language,appName,readOnly,graphLensAvailable,hasEvidence,lens=$bindable('none'),showDiscovery=$bindable(false),onopen,onreference}:{language:Language;appName:string;readOnly:boolean;graphLensAvailable:boolean;hasEvidence:boolean;lens?:Lens;showDiscovery?:boolean;onopen:(tool:string)=>void;onreference:()=>void}=$props();
 const t=(key:string)=>translate(key,language);
</script>
<div class="workspace-tools-panel" data-workspace-tools>
 <Disclosure label="Tools" {language} triggerId="workspace-menu-button" heading={`${t('Tools for')} ${appName}`} compact closeOnExit>
  <section aria-label={t('View')}>
   <h3>{t('View')}</h3>
   {#if graphLensAvailable}<label class="highlight"><span>{t('Highlight')}</span><select aria-label={t('Graph highlight')} bind:value={lens} disabled={!hasEvidence}><option value="none">{t('None')}</option><option value="rare">{t('Rare paths')}</option><option value="endings">{t('Endings')}</option><option value="exceptions">{t('Warnings & exceptions')}</option><option value="recovery">{t('Recovery')}</option></select></label>{/if}
   <label class="discovered"><input type="checkbox" bind:checked={showDiscovery}/><span>{t('Show discovered events')}</span></label>
  </section>
  <section aria-label={t('Investigate')}>
   <h3>{t('Investigate')}</h3>
   <button onclick={onreference}>{t('Saved views')}</button>
   {#if !readOnly}<button onclick={onreference}>{t('Save current view…')}</button>{/if}
   <button onclick={()=>onopen('errors')}>{t('Errors & exceptions')}</button>
   <button onclick={()=>onopen('about')}>{t('About this evidence')}</button>
  </section>
  <section aria-label={t('Manage this app')}>
   <h3>{t('Manage this app')}</h3>
   <button onclick={()=>onopen('discovery')}>{t('Unmapped events')}</button>
   {#if !readOnly}<button onclick={()=>onopen('configuration')}>{t('Mappings & goal')}</button><button onclick={()=>onopen('invitations')}>{t('Temporary viewers')}</button>{/if}
   <button onclick={()=>onopen('lifecycle')}>{t('Data lifecycle')}</button>
  </section>
 </Disclosure>
</div>
<style>
 .workspace-tools-panel{display:inline-flex;flex:none}section+section{border-top:1px solid var(--line);margin-top:8px;padding-top:8px}h3{font:600 12px/1.5 system-ui,sans-serif;color:var(--muted);margin:4px 8px}button{display:block;width:100%;min-width:44px;min-height:44px;padding:10px 8px;text-align:left;border:0;background:transparent;color:var(--ink);font-size:13px;font-weight:500;line-height:1.5;white-space:normal;overflow-wrap:anywhere}button:hover{background:#edf1e7}label{margin:0;font-size:13px;line-height:1.5}.highlight{display:grid;gap:4px;padding:4px 8px}.highlight select{min-height:44px;width:100%;font-size:13px}.discovered{display:flex;align-items:center;gap:10px;min-height:44px;padding:6px 8px;cursor:pointer}.discovered input{width:18px;height:18px;flex:none;margin:0}.discovered span{min-width:0}
</style>
