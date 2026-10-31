<script lang="ts">
 import type {Snippet} from 'svelte';
 import {translate,type Language} from './i18n';
 let {label,note='',icon='',language,children}:{label:string;note?:string;icon?:string;language:Language;children:Snippet}=$props();
 const id=$props.id();
 let open=$state(false),trigger:HTMLButtonElement,container:HTMLDivElement;
 let position=$state({left:12,top:60,height:400});
 function place(){if(!trigger)return;const r=trigger.getBoundingClientRect(),width=Math.min(440,window.innerWidth-24);const top=Math.min(r.bottom+6,window.innerHeight-100);position={left:Math.max(12,Math.min(r.left,window.innerWidth-width-12)),top:Math.max(12,top),height:Math.max(88,window.innerHeight-top-12)}}
 function toggle(){if(open){close();return}place();container.showPopover()}
 function close(){container?.hidePopover();trigger?.focus({preventScroll:true})}
</script>
<svelte:window onresize={()=>{if(open)place()}} onscroll={()=>{if(open)place()}}/>
<div class="flyout-control">
 <button type="button" class="flyout-trigger" bind:this={trigger} aria-label={translate(label,language)+(note?' · '+translate(note,language):'')} aria-expanded={open} aria-controls={id} aria-haspopup="dialog" title={translate(label,language)} onclick={toggle}>{icon||translate(label,language)}</button>
 <div id={id} bind:this={container} popover="auto" role="dialog" aria-label={translate(label,language)} class="disclosure-content" style:left={`${position.left}px`} style:top={`${position.top}px`} style:max-height={`${position.height}px`} ontoggle={e=>open=(e as ToggleEvent).newState==='open'}>
  <header><strong>{translate(label,language)}</strong><button type="button" class="text-button" onclick={close}>{translate('Close',language)}</button></header>
  {@render children()}
 </div>
</div>
<style>
 .flyout-control{display:inline-flex;min-width:0}.flyout-trigger{min-width:44px;min-height:44px;padding:6px 10px;background:transparent;border:0;color:var(--ink);font-size:13px;font-weight:600;white-space:nowrap}
 .flyout-trigger[aria-expanded=true]{background:#e7ece3}
 .disclosure-content{position:fixed;inset:auto;margin:0;width:min(440px,calc(100vw - 24px));box-sizing:border-box;overflow:auto;overscroll-behavior:contain;border:1px solid var(--line);border-radius:8px;background:var(--surface);padding:10px 14px;font-size:13px;line-height:1.55;overflow-wrap:anywhere;color:var(--ink);box-shadow:0 12px 35px #153f3933}
 header{position:sticky;top:-10px;display:flex;align-items:center;justify-content:space-between;gap:12px;background:var(--surface);z-index:1;border-bottom:1px solid var(--line);margin-bottom:8px}header button{min-height:44px;padding:6px 10px}
 .disclosure-content :global(p){margin:6px 0 10px}.disclosure-content :global(label){font-size:13px}.disclosure-content :global(button){min-height:44px}
</style>
