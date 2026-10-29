<script lang="ts">
 import {onMount} from 'svelte';
 import {dateTime} from './helpers';
 import {invitationIsLoopback} from './invitation';
 import {translate,type Language} from './i18n';
 type Invite={id:number;label:string;created:string;expires:string;redeemed:boolean;revoked:boolean};
 let {language='en',request}:{language?:Language;request:(path:string,body?:unknown)=>Promise<any>}=$props();
 const t=(key:string)=>translate(key,language);
 let label=$state(''),hours=$state(24),invites=$state<Invite[]>([]),url=$state(''),busy=$state(false),error=$state(''),copied=$state(false),now=$state(Date.now());
 let alive=true;
 async function refresh(){try{const result=await request('invitations');if(alive)invites=result}catch(e){if(alive)error=e instanceof Error?e.message:t('Request failed. Try again.')}}
 async function create(){busy=true;error='';url='';copied=false;try{const result=await request('invitations',{label:label.trim(),hours:Number(hours)});if(!alive)return;url=result.url;label='';await refresh()}catch(e){if(alive)error=e instanceof Error?e.message:t('Request failed. Try again.')}finally{if(alive)busy=false}}
 async function revoke(id:number){busy=true;error='';try{await request('invitations/revoke',{id});if(alive){url='';await refresh()}}catch(e){if(alive)error=e instanceof Error?e.message:t('Request failed. Try again.')}finally{if(alive)busy=false}}
 async function copy(){try{await navigator.clipboard.writeText(url);copied=true}catch{error=t('Copy unavailable. Select the displayed link.')}}
 function status(invite:Invite){return invite.revoked?'Revoked':Date.parse(invite.expires)<=now?'Expired':invite.redeemed?'Redeemed':'Waiting for viewer'}
 onMount(()=>{alive=true;void refresh();const poll=setInterval(()=>{now=Date.now();void refresh()},5000);return()=>{alive=false;url='';clearInterval(poll)}});
</script>
<section class="invitations" aria-label={t('Temporary viewers')}>
 <p class="section-description">{t('Invite someone to browse the whole application in read-only mode. They can view retained historical and live journeys, references, errors and sanitized diagnostic details. They cannot configure, save, delete or invite anyone.')}</p>
 <form class="panel" onsubmit={e=>{e.preventDefault();void create()}}>
  <div class="form-grid"><label>{t('Viewer label (optional)')}<input bind:value={label} maxlength="100" placeholder={t('For example: Product review')}/></label><label>{t('Access expires after')}<select bind:value={hours}><option value={1}>{t('1 hour')}</option><option value={24}>{t('24 hours')}</option><option value={168}>{t('7 days')}</option></select></label></div>
  <p class="inspector-note">{t('Expiry starts when the link is created. Anyone with the link can redeem it once; send it privately. Revocation blocks new requests immediately and closes live streams within five seconds. Already copied data cannot be recalled.')}</p>
  <button disabled={busy}>{t('Create invitation link')}</button>
 </form>
 {#if error}<p class="error" role="alert">{error}</p>{/if}
 {#if url}<div class="panel created-link"><h3>{t('Invitation link')}</h3><label>{t('Copy this link now')}<input type="text" readonly value={url} autocomplete="off" spellcheck="false"/></label><p class="inspector-note">{t('This link is shown only here and is not saved in browser storage.')}{#if invitationIsLoopback(url)} {t('This loopback address works only on this Mac; this does not expose Sentrame to another computer.')}{/if}</p><button class="secondary" onclick={copy}>{t(copied?'Copied':'Copy invitation link')}</button></div>{/if}
 <h3>{t('Invitation access')}</h3>
 {#each invites as invite (invite.id)}<article class="panel invite-entry"><div><strong>{invite.label||t('Temporary viewer')}</strong><p class="muted">{t(status(invite))} · {t('Expires')}: {dateTime(invite.expires)}</p></div>{#if !invite.revoked&&Date.parse(invite.expires)>now}<button class="secondary" disabled={busy} aria-label={t('Revoke access')+': '+(invite.label||t('Temporary viewer'))} onclick={()=>revoke(invite.id)}>{t('Revoke access')}</button>{/if}</article>{:else}<p class="empty-state">{t('No invitations yet')}</p>{/each}
 <p class="inspector-note">{t('Viewer sessions end on expiry, revocation or service restart. A redeemed link cannot be reused; create a new invitation if access is needed again.')}</p>
</section>
<style>
.created-link input{font-family:ui-monospace,SFMono-Regular,Consolas,monospace;font-size:11px}.created-link label{margin-bottom:12px}.invite-entry{display:flex;gap:15px;align-items:center;justify-content:space-between}.invite-entry p{margin:5px 0 0;font-size:11px}.invite-entry button{flex-shrink:0}.invitations>h3{margin-top:24px}@media(max-width:600px){.invite-entry{align-items:flex-start;flex-direction:column}}
</style>
