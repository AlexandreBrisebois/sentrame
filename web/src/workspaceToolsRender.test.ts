import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {createServer} from 'vite';
import {mkdtempSync,rmSync} from 'node:fs';
import {tmpdir} from 'node:os';
import {join} from 'node:path';
const cacheDir=mkdtempSync(join(tmpdir(),'sentrame-tools-render-'));
let server:Awaited<ReturnType<typeof createServer>>,component:any,render:any;
before(async()=>{server=await createServer({cacheDir,server:{middlewareMode:true,hmr:false},appType:'custom'});component=(await server.ssrLoadModule('/src/WorkspaceTools.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render});
after(async()=>{await server?.close();rmSync(cacheDir,{recursive:true,force:true})});
function html(extra={}){return render(component,{props:{language:'en',appName:'Test app',readOnly:false,graphLensAvailable:true,hasEvidence:true,lens:'rare',showDiscovery:true,onopen:()=>{},onreference:()=>{},...extra}}).body.replace(/<!--.*?-->/gs,'') as string}
test('tools names scope and groups, with explicit selected display state',()=>{const out=html();for(const value of ['Tools','Tools for Test app','View','Investigate','Manage this app','Saved views','Save current view','Unmapped events','Mappings &amp; goal','Temporary viewers','Errors &amp; exceptions','Data lifecycle','About this evidence'])assert.ok(out.includes(value),value);assert.match(out,/aria-controls=/);assert.match(out,/popover="auto"/);assert.match(out,/<input[^>]*type="checkbox"[^>]*checked/);assert.match(out,/<option[^>]*value="rare"[^>]*selected/);assert.doesNotMatch(out,/role="menu|aria-current/)});
test('viewer retains inspection but not owner writes, unsupported highlighting is absent',()=>{const out=html({readOnly:true,graphLensAvailable:false});for(const value of ['Saved views','Unmapped events','Data lifecycle','Errors &amp; exceptions'])assert.ok(out.includes(value),value);for(const value of ['Save current view','Mappings &amp; goal','Temporary viewers','Graph highlight'])assert.ok(!out.includes(value),value)});
test('French tools exposes equivalent task groups and empty highlight is disabled',()=>{const out=html({language:'fr',hasEvidence:false});for(const value of ['Outils','Outils pour Test app','Affichage','Explorer','Gérer cette application','Vues enregistrées','Fermer'])assert.ok(out.includes(value),value);assert.match(out,/<select[^>]*disabled/)});
