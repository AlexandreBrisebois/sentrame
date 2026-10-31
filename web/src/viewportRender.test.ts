import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {createServer} from 'vite';
import type {Snapshot} from './types';
let server:Awaited<ReturnType<typeof createServer>>,coverage:any,help:any,render:any;
before(async()=>{server=await createServer({server:{middlewareMode:true,hmr:false},appType:'custom'});coverage=(await server.ssrLoadModule('/src/CoverageDisclosure.svelte')).default;help=(await server.ssrLoadModule('/src/InvestigationHelp.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render});
after(async()=>{await server?.close()});
const model={journeys:[],nodes:[],edges:[],stats:{},diagnostics:[],unassigned:[],policy:'',warnings:[]};
const population={journeys:4,sourcePartial:false,partialJourneys:0,censoredJourneys:0,graphNodesOmitted:0,graphEdgesOmitted:0,definitions:[1]};
function snapshot(){return {model,baselineModel:model,evidence:[],coverage:{source:{recordLimit:10000,truncated:false},current:{...population},baseline:{...population},evidence:{limit:2000,omitted:0}}} as unknown as Snapshot}
function html(s:Snapshot|null,extra={}){return render(coverage,{props:{snapshot:s,language:'en',...extra}}).body.replace(/<!--.*?-->/gs,'') as string}
const visible=(out:string)=>out.replace(/<div\b[^>]*popover="auto"[^>]*>[\s\S]*?<\/div>/g,'');
test('compact coverage retains details without a complete-traffic claim or zero-omission warning',()=>{const out=html(snapshot());assert.match(out,/<button[^>]*aria-label="Coverage · collection unknown"/);assert.match(out,/popover="auto"/);assert.match(out,/collection unknown/);assert.match(out,/Graph and received-record list omissions do not reduce ordered-journey denominators/);assert.doesNotMatch(visible(out),/class="coverage-warning"/)});
test('material source, population and display limits remain visible with coverage closed',()=>{const s=snapshot();s.coverage!.source.truncated=true;s.coverage!.current!.sourcePartial=true;s.coverage!.current!.partialJourneys=2;s.coverage!.baseline!.censoredJourneys=3;s.coverage!.current!.graphNodesOmitted=8;s.coverage!.evidence.omitted=40;const out=visible(html(s,{live:true,baseline:true}));for(const value of ['Source scan capped','Live observation','2 Partial beginnings','App-wide baseline','3 Censored continuations','8 omitted nodes','40 omitted records'])assert.ok(out.includes(value),value);assert.match(out,/Received-record list/)});
test('unknown coverage remains visible and English/French disclose equivalent scopes',()=>{assert.match(visible(html(null)),/Coverage metadata unavailable/);const s=snapshot();s.coverage!.current!.graphEdgesOmitted=7;const out=visible(html(s,{language:'fr'}));assert.match(out,/7 connexions omises/);assert.match(out,/sélection actuelle/i)});
test('help exposes task meaning with localized native disclosure and population legend only when relevant',()=>{const out=render(help,{props:{layout:'matrix',language:'fr',live:true}}).body;assert.match(out,/popover="auto"/);assert.match(out,/aria-expanded="false"/);assert.match(out,/Où vont ensuite les parcours/);assert.match(out,/Unité/);assert.match(out,/Fermer/);assert.match(out,/référence/i);const historical=render(help,{props:{layout:'hierarchy',language:'en'}}).body;assert.doesNotMatch(historical,/Baseline uses B/)});
test('each coverage limitation independently produces a visible scope-specific caveat',()=>{
 const cases:[string,(s:Snapshot)=>void,string][]=[
  ['source',s=>s.coverage!.source.truncated=true,'Source scan capped'],
  ['partial source',s=>s.coverage!.current!.sourcePartial=true,'Observed partial evidence'],
  ['partial start',s=>s.coverage!.current!.partialJourneys=2,'2 Partial beginnings'],
  ['censoring',s=>s.coverage!.current!.censoredJourneys=3,'3 Censored continuations'],
  ['nodes',s=>s.coverage!.current!.graphNodesOmitted=8,'8 omitted nodes'],
  ['connections',s=>s.coverage!.current!.graphEdgesOmitted=7,'7 omitted connections'],
  ['record list',s=>s.coverage!.evidence.omitted=40,'40 omitted records'],
  ['baseline',s=>s.coverage!.baseline!.sourcePartial=true,'App-wide baseline: Observed partial evidence']
 ];
 for(const [name,mutate,expected] of cases){const s=snapshot();mutate(s);const out=visible(html(s,{live:true,baseline:true}));assert.ok(out.includes(expected),name);assert.equal(s.coverage!.current!.journeys,4,name+' retained denominator');}
});
