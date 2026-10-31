import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {createServer} from 'vite';
import type {Model} from './types.ts';
let server:Awaited<ReturnType<typeof createServer>>,component:any,render:any;
before(async()=>{server=await createServer({server:{middlewareMode:true},appType:'custom'});component=(await server.ssrLoadModule('/src/Graph.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render});
after(async()=>{await server?.close()});
const empty:Model={journeys:[],nodes:[],edges:[],stats:{},diagnostics:[],unassigned:[],policy:'',warnings:[]};
const populated:Model={...empty,nodes:['A','B'].map(id=>({id,label:id,occurrences:1,journeys:1,visits:1,stops:0,pending:0,timing:[]})),edges:[{from:'A',to:'B',occurrences:1,journeys:1,visits:1}]};
test('isolating baseline omits hidden Live counts instead of representing hidden contributions as zero',()=>{
 for(const layout of ['hierarchy','force','sankey']){
  const html=render(component,{props:{model:empty,baseline:populated,liveMode:true,liveVisible:false,layout}}).body;
  assert.match(html,/aria-label="A; B 1 baseline"/);
  assert.match(html,/aria-label="A → B; 1 baseline transitions"/);
  assert.doesNotMatch(html,/L 0 live|A → B; 0 transitions/);
 }
});
test('shared and live-only presentations retain independent supported counts',()=>{
 const shared=render(component,{props:{model:populated,baseline:populated,liveMode:true}}).body;
 assert.match(shared,/aria-label="A; B 1 baseline · L 1 live"/);
 assert.match(shared,/aria-label="A → B; 1 transitions; 1 baseline transitions"/);
 const live=render(component,{props:{model:populated,liveMode:true}}).body;
 assert.match(live,/aria-label="A; L 1 live"/);
 assert.doesNotMatch(live,/baseline transitions/);
});
