import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {createServer} from 'vite';
import {mkdtempSync,rmSync} from 'node:fs';
import {tmpdir} from 'node:os';
import {join} from 'node:path';
const cacheDir=mkdtempSync(join(tmpdir(),'sentrame-connection-render-'));
let server:Awaited<ReturnType<typeof createServer>>,component:any,render:any;
const originalWindow=Object.getOwnPropertyDescriptor(globalThis,'window');
before(async()=>{
 Object.defineProperty(globalThis,'window',{configurable:true,value:{location:{origin:'http://127.0.0.1:8096'}}});
 server=await createServer({cacheDir,server:{middlewareMode:true,hmr:false},appType:'custom'});
 component=(await server.ssrLoadModule('/src/Apps.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render;
});
after(async()=>{await server?.close();rmSync(cacheDir,{recursive:true,force:true});if(originalWindow)Object.defineProperty(globalThis,'window',originalWindow);else Reflect.deleteProperty(globalThis,'window')});
function html(language='en',apps=[{id:'a',generation:'g',name:'Client app',state:'active'}],selected='a'){
 return render(component,{props:{apps,selected,language,request:()=>{throw new Error('Keys must only load after an explicit owner action')},onselect:()=>{},onclose:()=>{}}}).body.replace(/<!--.*?-->/gs,'') as string;
}
test('app connection offers one masked write key and standard signal endpoints',()=>{
 const out=html();assert.equal((out.match(/<fieldset/g)??[]).length,1);
 for(const label of ['Telemetry key','Reveal key','Rotate key','Copy Authorization header','••••••••••••••••','/v1/logs','/v1/traces','/v1/metrics'])assert.ok(out.includes(label),label);
 assert.doesNotMatch(out,/<legend>Browser|<legend>Server|server secret|using this role/);
 assert.ok(out.includes('permits telemetry writes only'));assert.ok(out.includes('does not prove authentic activity or saved changes'));
});
test('French connection presents equivalent one-key actions and limits',()=>{
 const out=html('fr');for(const label of ['Clé de télémétrie','Afficher la clé','Renouveler la clé','uniquement l’envoi de télémétrie'])assert.ok(out.includes(label),label);
 assert.equal((out.match(/<fieldset/g)??[]).length,1);assert.doesNotMatch(out,/<legend>Browser|<legend>Server/);
});
test('unregistered workspace does not offer or load credentials',()=>{
 const out=html('en',[],'');assert.ok(out.includes('Register app'));assert.doesNotMatch(out,/Reveal key|Rotate key|<fieldset/);
});
