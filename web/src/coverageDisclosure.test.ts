import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {readFile} from 'node:fs/promises';
import {createServer} from 'vite';
import type {Journey,Model,Snapshot} from './types.ts';
import {comparisonReason,transitionMetrics,type PopulationCoverage} from './investigation.ts';
import {translate} from './i18n.ts';
let server:Awaited<ReturnType<typeof createServer>>,component:any,graph:any,investigation:any,render:any;
before(async()=>{
 server=await createServer({server:{middlewareMode:true,hmr:false},appType:'custom'});
 component=(await server.ssrLoadModule('/src/CoverageDisclosure.svelte')).default;
 graph=(await server.ssrLoadModule('/src/Graph.svelte')).default;
 investigation=(await server.ssrLoadModule('/src/InvestigationView.svelte')).default;
 render=(await server.ssrLoadModule('svelte/server')).render;
});
after(async()=>{await server?.close()});
const empty=():Model=>({journeys:[],nodes:[],edges:[],diagnostics:[],unassigned:[],stats:{},policy:'',warnings:[]});
const cov=(extra:Partial<PopulationCoverage>={}):PopulationCoverage=>({from:'1793390000000000000',to:'1793393600000000000',scope:'',selectionPolicy:'evidence-window',journeys:10,sourcePartial:false,partialJourneys:0,censoredJourneys:0,definitions:[1],graphNodesOmitted:0,graphEdgesOmitted:0,...extra});
const snap=(extra:Partial<Snapshot>={}):Snapshot=>({revision:'1',now:1793393600000,evidence:[],definitions:[],activeDefinition:0,views:[],model:empty(),coverage:{source:{truncated:false,recordLimit:10000,collectionCompleteness:'unknown; exporter sampling/loss not observed'},current:cov(),evidence:{limit:2000,omitted:0}},...extra});
const body=(c:any,props:any)=>render(c,{props}).body.replace(/<!--.*?-->/gs,'');
const text=(s:string)=>s.replace(/<[^>]+>/g,' ').replace(/\s+/g,' ').trim();
const html=(snapshot:Snapshot|null,language:'en'|'fr'='en',live=false,baseline=false)=>body(component,{snapshot,language,live,baseline});

test('missing independent coverage stays unavailable rather than inferred complete from legacy limits',()=>{
 const en=text(html(snap({coverage:undefined,limits:{truncated:false}})));
 assert.match(en,/unavailable|unknown/i);
 assert.doesNotMatch(en,/source complete|complete source|0 omitted records/i);
 for(const language of ['en','fr'] as const){assert.ok(text(html(null,language)).length>0);assert.doesNotMatch(html(null,language),/undefined|NaN/);}
});

test('graph and received-list caps remain independent from uncapped source and population',()=>{
 const s=snap();s.coverage!.current=cov({graphNodesOmitted:200,graphEdgesOmitted:202});
 s.coverage!.evidence={limit:2000,omitted:8000};s.evidence=Array(2000).fill({}) as Snapshot['evidence'];
 s.model.nodes=Array.from({length:100},(_,i)=>({id:'n'+i,label:'n'+i,occurrences:1,journeys:1,visits:1,stops:0,pending:0,timing:[]}));
 s.model.edges=Array.from({length:100},(_,i)=>({from:'n'+i,to:'n'+i,occurrences:1,journeys:1,visits:1}));
 s.limits={truncated:true};const out=text(html(s));
 for(const n of ['10000','100','200','202','2000','8000'])assert.ok(out.includes(n),'Missing independent count '+n);
 assert.match(out,/unknown/i);
 assert.doesNotMatch(out,/Query work limit reached\. Counts are partial/);
 assert.equal(comparisonReason(s.coverage!.current,cov()),null);
});

test('current and baseline partial/censored and graph counts are independently disclosed',()=>{
 const s=snap({baselineModel:empty()});
 s.coverage!.current=cov({journeys:11,partialJourneys:2,censoredJourneys:3,graphNodesOmitted:17,graphEdgesOmitted:19});
 s.coverage!.baseline=cov({journeys:29,sourcePartial:true,partialJourneys:5,censoredJourneys:7,graphNodesOmitted:23,graphEdgesOmitted:31});
 const out=text(html(s,'en',true,true));
 for(const n of [11,2,3,17,19,29,5,7,23,31])assert.match(out,new RegExp('(?:^|\\D)'+n+'(?:\\D|$)'),'Missing population count '+n);
 assert.match(out,/Live|live/);assert.match(out,/baseline/i);
 assert.match(out,/partial/i);assert.match(out,/censored|continuation/i);
 assert.ok(comparisonReason(s.coverage!.current,s.coverage!.baseline));
 const hidden=text(html(s,'en',true,false));assert.doesNotMatch(hidden,/31/);
});

test('retained source work cap and unknown exporter completeness are separate disclosures',()=>{
 const s=snap();s.coverage!.source.truncated=true;s.coverage!.current=cov({sourcePartial:true});
 const out=text(html(s));assert.match(out,/10000/);assert.match(out,/limit|truncat|partial/i);assert.match(out,/unknown/i);
 const uncapped=text(html(snap()));assert.match(uncapped,/unknown/i);
});

test('coverage explanation has English and French parity',async()=>{
 const source=await readFile(new URL('./CoverageDisclosure.svelte',import.meta.url),'utf8');
 const keys=[...source.matchAll(/\bt\('([^']+)'\)/g)].map(m=>m[1]);
 assert.ok(keys.length>5,'Coverage needs explicit labels and explanations');
 for(const key of keys)assert.notEqual(translate(key,'fr'),key,'Missing French coverage label '+key);
 assert.notEqual(text(html(snap(),'en')),text(html(snap(),'fr')));
});

test('ordered metrics preserve full outgoing denominator under graph and evidence presentation caps',()=>{
 const js=Array.from({length:10},(_,i)=>({id:'journey-'+i,scope:'scope-'+i,tab:'tab',definition:1,version:'v1',start:'1793390000000000000',end:'0',lastMapped:'1793390001000000000',status:'active',path:['A',i<6?'B':'C'],steps:['A',i<6?'B':'C'].map((step,k)=>({step,event:step,evidenceId:i*2+k,time:String(1793390000000000000n+BigInt(k)*1000000000n),source:'browser'})),evidence:[],timings:[],warnings:[],partial:false,recovered:false} as Journey));
 const m={...empty(),journeys:js};const capped=cov({graphNodesOmitted:200,graphEdgesOmitted:202});
 const rows=transitionMetrics(m.journeys);assert.deepEqual(rows.map(x=>[x.to,x.occurrences,x.denominator,x.share]),[['B',6,10,.6],['C',4,10,.4]]);
 const out=body(investigation,{view:'matrix',model:m,language:'en',onedge(){},onjourney(){},coverage:capped,label:'Current selection',population:'current'});
 assert.match(out,/A → B; 6 Transition occurrences; 10 Observed outgoing occurrences; Retained source/);
 assert.doesNotMatch(out,/Coverage is partial/);assert.equal(comparisonReason(capped,cov()),null);
});

test('graph disjoint union display cap discloses additional hidden nodes and edges in each layout',()=>{
 const population=(prefix:string)=>{
  const m=empty();m.nodes=Array.from({length:100},(_,i)=>({id:prefix+i,label:prefix+i,occurrences:1,journeys:1,visits:1,stops:0,pending:0,timing:[]}));
  m.edges=Array.from({length:99},(_,i)=>({from:prefix+i,to:prefix+(i+1),occurrences:1,journeys:1,visits:1}));return m;
 };
 for(const layout of ['hierarchy','force','sankey']){
  const out=body(graph,{model:population('live-'),baseline:population('base-'),layout,language:'en',liveMode:true});
  assert.equal((out.match(/data-node=/g)??[]).length,100);
  const legend=text(out.match(/<span data-coverage="graph-display"[^>]*>(.*?)<\/span>/s)?.[1]??'');
  assert.match(legend,/100\s*\/\s*200/,'100 of 200 returned nodes');
  assert.match(legend,/99\s*\/\s*198/,'99 of 198 returned edges');
  assert.match(legend,/omitt|hidden|display/i);
 }
});
