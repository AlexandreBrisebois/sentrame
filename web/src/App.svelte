<script lang="ts">
 import Workspace from './Workspace.svelte';
 import {appURL} from './apps';
 let appID=$state(new URLSearchParams(window.location.search).get('app')??'');
 let settings=$state(false),revision=$state(0);
 function select(id:string,connection=false){settings=connection;appID=id;history.pushState(null,'',appURL(id));revision++}
 function navigate(){appID=new URLSearchParams(window.location.search).get('app')??'';settings=false;revision++}
</script>
<svelte:window onpopstate={navigate}/>
{#key revision}<Workspace initialAppID={appID} initialSettings={settings} onselect={select}/>{/key}
