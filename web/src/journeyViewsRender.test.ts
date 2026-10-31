import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {readFileSync} from 'node:fs';
import {createServer} from 'vite';
import type {Journey,Model,Evidence} from './types.ts';
import {translate} from './i18n.ts';

let server:Awaited<ReturnType<typeof createServer>>, component:any, render:any;
before(async()=>{
 server=await createServer({server:{middlewareMode:true},appType:'custom'});
 component=(await server.ssrLoadModule('/src/JourneyViews.svelte')).default;
 render=(await server.ssrLoadModule('svelte/server')).render;
});
after(async()=>{await server?.close()});
const journey=(id:string,path:string[])=>({id,scope:id,tab:'tab',start:'1793368800000000000',end:'1793368810000000000',status:'completed',path,steps:path.map((step,i)=>({step,event:step,evidenceId:i+1,time:String(1793368800000000000n+BigInt(i)*1000000000n),source:'browser'})),evidence:[],warnings:[],partial:false} as unknown as Journey);
const model=(journeys:Journey[]):Model=>({journeys,nodes:[],edges:[],diagnostics:[],unassigned:[],stats:{},policy:'',warnings:[]});
function html(view:string,m:Model,extra:Record<string,unknown>={}){return render(component,{props:{view,model:m,language:'en',selectedStep:'',selectedEdge:'',onedge(){},onjourney(){},label:'Current selection',...extra}}).body.replace(/<!--.*?-->/gs,'') as string}

test('timeline retains repeated occurrences, captured bounds, gaps and diagnostic detail',()=>{
 const j=journey('timeline',['A','B','A']);j.partial=true;j.steps[2].gapBefore=true;
 j.windowFrom='0';j.windowTo='0';j.end='0';j.lastMapped=j.start;
 j.warnings=['Display uses occurrence time and stable receipt tie-break; causal order is unproven.'];
 j.evidence=[{id:99,name:'synthetic.exception',time:j.start,source:'server',signal:'logs',service:'synthetic',attrs:{'exception.type':'SyntheticError'}} as Evidence];
 const m=model([j]);m.diagnostics=[{journeyId:j.id,evidenceId:99,type:'SyntheticError',severity:17}];
 const out=html('journey-timeline',m,{chosen:j.id});
 assert.match(out,/<ol class="journey-timeline">/);assert.equal((out.match(/<strong>A<\/strong>/g)??[]).length,2);
 assert.match(out,/Captured boundary/);assert.doesNotMatch(out,/1970|Evidence window/);assert.match(out,/Missing intermediate steps/);assert.match(out,/earlier visits may be missing/);
 assert.match(out,/<details><summary>! SyntheticError/);assert.match(out,/synthetic.exception/);assert.match(out,/Show ordered evidence/);
 j.evidence=[];assert.match(html('journey-timeline',m,{chosen:j.id}),/Evidence unavailable/);
 assert.match(html('journey-timeline',m,{chosen:'expired'}),/Expired detail cannot be reconstructed/);
});
test('matrix exposes directed and diagonal counts with table headers and distinguishes unknown zeroes',()=>{
 const m=model([journey('one',['A','B','A','B','B'])]);const out=html('matrix',m);
 assert.match(out,/<th scope="col">A/);assert.match(out,/<th scope="row">B/);
 for(const label of ['A → B: 2 Occurrences, 1 Journeys, 1 Visits','B → A: 1 Occurrences, 1 Journeys, 1 Visits','B → B: 1 Occurrences, 1 Journeys, 1 Visits'])assert.ok(out.includes(label));
 assert.match(out,/<td>0<\/td>/);assert.doesNotMatch(html('matrix',m,{partial:true}),/<td>0<\/td>/);
 m.journeys[0].steps[2].gapBefore=true;const partial=html('matrix',m);assert.match(partial,/<td>—<\/td>/);assert.ok(!partial.includes('B → A:'));
});
test('repeat ranking offers counts and journey actions without treating loops as defects',()=>{
 const out=html('repeats',model([journey('loop',['A','B','A','B']),journey('self',['A','A'])]));
 assert.match(out,/Repetition is not automatically a problem/);assert.match(out,/Additional occurrences/);assert.match(out,/>A<.*?<td>2<\/td><td>2<\/td>/s);
 assert.match(out,/B → A/);assert.match(out,/A → B/);assert.doesNotMatch(out,/A → A/);assert.match(out,/<button[^>]*>Journey loop<\/button>/);
});
test('dense matrix and long sequence remain bounded with labeled omissions and journey widths',()=>{
 const journeys=Array.from({length:120},(_,i)=>journey(String(i),['A',...Array.from({length:12},(_,n)=>`B${i}-${n}`)]));
 const m=model(journeys),matrix=html('matrix',m),sequence=html('sequence',m);
 assert.equal((matrix.match(/scope="col"/g)??[]).length,100);assert.match(matrix,/Matrix limited to the first 100 stages/);assert.match(matrix,/100\/1441 · 1341 omitted stages/);
 assert.equal((sequence.match(/class="sequence-flow-row"/g)??[]).length,20);assert.match(sequence,/Other path groups.*?100.*?100 Journeys/s);
 assert.match(sequence,/1 continue or are censored/);assert.match(sequence,/120 Journeys/);assert.match(sequence,/width:100%/);assert.doesNotMatch(sequence,/13\. B/);
});
test('all fixed view labels and actions have French translations, including dynamic units',()=>{
 const source=readFileSync(new URL('./JourneyViews.svelte',import.meta.url),'utf8');
 const keys=[...source.matchAll(/t\('([^']+)'\)/g)].map(m=>m[1]);
 for(const key of [...keys,'Current selection','Live observation','App-wide baseline','Completed','Active','Inactive','Exited','Display uses occurrence time and stable receipt tie-break; causal order is unproven.'])assert.ok(['Occurrences','Transition'].includes(key)||translate(key,'fr')!==key,`Missing French: ${key}`);
 const m=model([journey('loop',['A','B','A','B'])]);
 for(const view of ['journey-timeline','matrix','repeats','sequence']){
  const out=html(view,m,{language:'fr',chosen:'loop'});
  assert.match(out,/Sélection actuelle/);assert.doesNotMatch(out,/>Show ordered evidence<|>Journey loop<|>Occurrences \/ journeys \/ visits</);
 }
 assert.match(html('matrix',m,{language:'fr'}),/2 Occurrences, 1 Parcours, 1 Visites/);
});

test('French timeline translates partial-window and navigation warnings',()=>{
 const j=journey('partial',['A','C']);j.partial=true;j.rightCensored=true;j.steps[1].gapBefore=true;
 j.warnings=[
  'Correlated navigation; no completion goal configured. Excluded from goal conversion.',
  'Pending follow-up / right-censored; no stopping cause is known.',
  'Window omits the captured beginning; original boundary remains recorded.',
  'Window omits intermediate steps; no direct transition is invented across the gap.',
  'Window ends before a captured closure; pending/right-censored, not an observed dead end.',
  'Reference begins inside a captured journey; earlier evidence remains outside this projection.',
  'Reference ends before later captured evidence; operational goal and captured boundary are unchanged.',
  'Requested endpoint is unavailable or crosses a captured boundary; showing only this recorded segment, without joining journeys.'
 ];
 const out=html('journey-timeline',model([j]),{language:'fr',chosen:j.id});
 for(const warning of j.warnings){assert.notEqual(translate(warning,'fr'),warning);assert.ok(out.includes(translate(warning,'fr')));assert.ok(!out.includes(warning))}
 assert.match(out,/Étapes intermédiaires manquantes/);assert.match(out,/La suite est hors de cette vue/);
});
