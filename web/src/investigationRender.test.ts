import {test,before,after} from 'node:test';import {strict as assert} from 'node:assert';import {createServer} from 'vite';import type {Journey,Model} from './types.ts';import {defaultPreferences} from './investigation.ts';
import {defaultInvestigationState} from './investigationNavigation.ts';
let server:Awaited<ReturnType<typeof createServer>>,component:any,render:any;
before(async()=>{server=await createServer({server:{middlewareMode:true},appType:'custom'});component=(await server.ssrLoadModule('/src/InvestigationView.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render});after(async()=>{await server?.close()});
const journey=(path:string[])=>({id:'one',scope:'one',tab:'tab',start:'1793390000000000000',lastMapped:'1793390002000000000',end:'0',status:'active',path,steps:path.map((step,i)=>({step,event:step,evidenceId:i,time:String(1793390000000000000n+BigInt(i)*1000000000n),source:'browser'})),evidence:[],warnings:[],partial:false} as Journey);
const model=(j:Journey):Model=>({journeys:[j],nodes:[],edges:[],diagnostics:[],unassigned:[],stats:{},policy:'',warnings:[]});
function html(view:string,path:string[],extra={}){return render(component,{props:{view,model:model(journey(path)),language:'en',onedge(){},onjourney(){},label:'Live observation',population:'live',...extra}}).body.replace(/<!--.*?-->/gs,'') as string}
test('heatmap declares one metric, population, scale, keyboard model and matching evidence',()=>{const out=html('matrix',Array(10).fill('Vote'));assert.match(out,/Transition occurrences/);assert.match(out,/Linear scale/);assert.match(out,/tabindex="0"/);assert.match(out,/Vote → Vote; 9 Transition occurrences; 9 Observed outgoing occurrences/);assert.match(out,/Matching journeys/);assert.match(out,/Outgoing shares include all retained/)});
test('repeats distinguish prevalence, extra occurrences and intensity without causal diagnosis',()=>{const out=html('repeats',Array(10).fill('Vote'));assert.match(out,/9 · 10 total occurrences/);assert.match(out,/1\/1 · 100.0%/);assert.match(out,/9.00 extras per repeating journey/);assert.match(out,/No qualifying non-self returns observed/)});
test('vertical timeline preserves occurrences, relative clock basis and unavailable diagnostics in French',()=>{const out=html('journey-timeline',['A','B','A'],{chosen:'one',language:'fr'});assert.match(out,/vertical-timeline/);assert.match(out,/\+2.000 s/);assert.doesNotMatch(out,/1970/);assert.match(out,/Preuves indisponibles|Preuve indisponible|Éléments indisponibles|Evidence unavailable/);assert.match(out,/limites|Limite|enregistr|conserv/)});
test('outgoing share with zero outgoing denominator remains unavailable',()=>{const out=html('matrix',['A'],{preferences:{...defaultPreferences(),metric:'share'}});assert.match(out,/A → A; — Outgoing occurrence share/);assert.doesNotMatch(out,/A → A; 0.0%/)});
test('stage search preserves native search semantics after transient state restoration',()=>{assert.match(html('matrix',['A']),/type="search"/)});

test('repeat rankings expose affected, prevalence and extra-occurrence order with visible denominators',()=>{
 const make=(id:string,path:string[])=>({...journey(path),id});
 const js=[...Array.from({length:100},(_,i)=>make('a'+i,i<20?['A','A']:['A'])),...Array.from({length:10},(_,i)=>make('b'+i,i<8?['B','B']:['B'])),make('c',Array(30).fill('C'))];
 const view=(sort:string)=>html('repeats',[],{model:{...model(js[0]),journeys:js},transient:{...defaultInvestigationState(),repeatSort:sort}});
 const affected=view('journeys'),prevalence=view('prevalence'),extras=view('extras');
 assert.match(affected,/20\/100 · 20.0%/);assert.match(affected,/8\/10 · 80.0%/);
 assert.ok(affected.indexOf('<th>A</th>')<affected.indexOf('<th>B</th>'));
 assert.ok(prevalence.indexOf('<th>B</th>')<prevalence.indexOf('<th>A</th>'));
 assert.ok(extras.indexOf('<th>C</th>')<extras.indexOf('<th>A</th>'));
 assert.match(affected,/partial beginnings limit interpretation/);
});

test('timeline grouping starts off and preserves every original ordinal when enabled',()=>{
 const path=[...Array(10).fill('Vote'),'A','Vote'];
 const off=html('journey-timeline',path,{chosen:'one'}),on=html('journey-timeline',path,{chosen:'one',preferences:{...defaultPreferences(),groupRepeats:true}});
 assert.equal((off.match(/<li>/g)??[]).length,12);
 assert.equal((on.match(/<li>/g)??[]).length,3);assert.match(on,/×10 Occurrences/);
 for(let ordinal=1;ordinal<=10;ordinal++)assert.match(on,new RegExp('<h4>'+ordinal+'\\. Vote'));
 assert.match(on,/<h4>12\. Vote/);assert.match(on,/Evidence details · 1–10/);
});

test('empty and one-journey entry retain explicit recovery and direct inspection in both languages',()=>{
 for(const language of ['en','fr']){
  const empty=html('journey-timeline',[],{language,model:{...model(journey([])),journeys:[]}});
  assert.match(empty,/<select/);assert.doesNotMatch(empty,/class="journey-story"/);
  const one=html('matrix',['A'],{language});assert.match(one,/id="one-journey-live"/);
 }
 const empty=html('journey-timeline',[],{model:{...model(journey([])),journeys:[]}});
 assert.match(empty,/No retained journey detail/);assert.match(empty,/Summaries cannot reconstruct ordered detail/);assert.match(empty,/Browse matching journeys/);
});

test('investigation task, metric, evidence and recovery labels retain French parity',async()=>{
 const {readFile}=await import('node:fs/promises');const {translate}=await import('./i18n.ts');
 const keys=new Set(['Choose an investigation task','How do observed stages connect?','Where do journeys go next?','Which activities repeat?','Which ordered routes are shared?','What happened in one journey?','Select a stage, transition or matching journey to inspect supporting evidence.','Unit: journeys sharing an ordered prefix','Unit: ordered stage occurrences','Unit: observed occurrences, journeys and visits','Live observation','Current selection','App-wide baseline','All app traffic','Selected visit','Unavailable','Definition','Observed partial evidence','Refresh']);
 for(const file of ['InvestigationView.svelte','MatchingExamples.svelte','JourneyStory.svelte']){const source=await readFile(new URL(file,import.meta.url),'utf8');for(const match of source.matchAll(/\bt\('([^']+)'\)/g))keys.add(match[1]);}
 for(const key of keys)assert.ok(['Occurrences','Transition'].includes(key)||translate(key,'fr')!==key,`Missing French label: ${key}`);
});
