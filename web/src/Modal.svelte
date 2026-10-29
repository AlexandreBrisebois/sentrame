<script lang="ts">
 import { onMount, type Snippet } from 'svelte';
 let { title, onclose, children, wide=false }:{title:string;onclose:()=>void;children:Snippet;wide?:boolean}=$props();
 let dialog:HTMLDialogElement;
 onMount(()=>{const previous=document.activeElement as HTMLElement|null;dialog.showModal();return()=>{dialog.close();if(previous?.isConnected)previous.focus()}});
</script>
<dialog bind:this={dialog} class:wide class="workspace-dialog" aria-label={title} oncancel={e=>{e.preventDefault();onclose()}} onclick={e=>{if(e.target===dialog)onclose()}}>
 <div class="dialog-content"><header class="dialog-heading"><h2>{title}</h2><button class="icon-button" aria-label="Close / Fermer" onclick={onclose}>×</button></header>{@render children()}</div>
</dialog>
