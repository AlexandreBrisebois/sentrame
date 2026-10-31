import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {createServer} from 'vite';
import type {Journey,Model,Evidence} from './types.ts';
let server:Awaited<ReturnType<typeof createServer>>,component:any,render:any;
before(async()=>{server=await createServer({server:{middlewareMode:true},appType:'custom'});component=(await server.ssrLoadModule('/src/JourneyStory.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render});
after(async()=>{await server?.close()});
const origin=1793390000000000000n;
const journey=(path=['A','B','A']):Journey=>({id:'story-one',scope:'visit',tab:'tab',definition:1,version:'test',mixedVersion:false,start:String(origin-1000000000n),end:'0',lastMapped:String(origin+2000000000n),status:'active',reason:'',path,steps:path.map((step,i)=>({step,event:step,evidenceId:i+1,time:String(origin+BigInt(i)*1000000000n),source:'browser'})),evidence:[],timings:[],warnings:[],partial:false,recovered:false});
const evidence=(id:number,extra={}):Evidence=>({id,name:'synthetic.'+id,source:'browser',signal:'logs',time:String(origin),receipt:String(origin+1000000n),observed:String(origin),service:'synthetic-story',version:'test',scope:'visit',tab:'tab',record:'record-'+id,action:'',occurrence:'',trace:'',span:'',sequence:id,severity:9,attrs:{},fallback:false,definition:1,...extra});
const model=(j:Journey,diagnostics:any[]=[]):Model=>({journeys:[j],nodes:[],edges:[],diagnostics,unassigned:[],stats:{},policy:'',warnings:[]});
const html=(j:Journey,extra:Record<string,unknown>={})=>render(component,{props:{journey:j,model:model(j),language:'en',...extra}}).body.replace(/<!--.*?-->/gs,'') as string;
test('current story uses first retained occurrence for relative time and preserves captured boundary',()=>{
 const j=journey();j.partial=true;j.rightCensored=true;j.steps[2].gapBefore=true;j.windowFrom=String(origin);j.windowTo=String(origin+3000000000n);
 const out=html(j);assert.match(out,/Observed span: 2.000 s/);assert.match(out,/\+0.000 s/);assert.match(out,/\+2.000 s\s*· Δ 1.000 s/);assert.match(out,/Captured boundary/);assert.match(out,/Evidence window/);assert.match(out,/not the captured beginning/);assert.match(out,/Missing intermediate steps/);assert.match(out,/Continuation is outside this view/);assert.equal((out.match(/<strong>A<\/strong>/g)??[]).length,2);
});
test('absent and zero clock bounds stay unavailable instead of epoch or fabricated zero duration',()=>{
 for(const value of ['0','',undefined]){const j=journey([]);j.start=value as string;j.end=value as string;j.lastMapped=value as string;j.windowFrom=value;j.windowTo=value;const out=html(j);assert.match(out,/Observed span: Unavailable/);assert.match(out,/Captured boundary: Unavailable → Unavailable/);assert.match(out,/No retained journey detail/);assert.doesNotMatch(out,/1970|Evidence window|Observed span: 0.000 s/)}
});
test('receipt fallback and inconsistent occurrence clocks cannot become measured relative elapsed',()=>{
 const j=journey();j.evidence=[evidence(1,{fallback:true})];let out=html(j);assert.match(out,/Observed span: Unavailable/);assert.match(out,/Relative timing unavailable/);assert.doesNotMatch(out,/Observed span: 2.000 s/);
 j.evidence=[];j.steps[2].time=String(origin-1000000000n);out=html(j);assert.match(out,/Observed span: Unavailable/);assert.match(out,/Δ Unavailable/);assert.doesNotMatch(out,/-1.000 s|-2.000 s/);
 j.steps[1].time='0';out=html(j);assert.match(out,/Relative timing unavailable/);assert.doesNotMatch(out,/1970/);
});
test('placed diagnostic markers require exact evidence membership, unplaced remain separate',()=>{
 const j=journey();j.evidence=[evidence(2,{trace:'0123456789abcdef0123456789abcdef'}),evidence(99)];const m=model(j,[{journeyId:j.id,evidenceId:2,type:'PlacedError',severity:17},{journeyId:j.id,evidenceId:99,type:'UnplacedError',severity:17},{journeyId:'other',evidenceId:2,type:'OtherJourneyError',severity:17}]);
 const out=html(j,{model:m,ontrace(){}}),placed=out.indexOf('PlacedError'),listEnd=out.indexOf('</ol>'),unplaced=out.indexOf('UnplacedError');assert.ok(placed<listEnd);assert.ok(unplaced>listEnd);assert.doesNotMatch(out,/OtherJourneyError/);assert.match(out,/exact stage attachment or causality/);assert.match(out,/Open retained trace/);assert.match(out,/Occurrence time/);assert.match(out,/Observed time/);assert.match(out,/Receipt time/);assert.match(out,new RegExp(String(origin)));
});
test('expired raw detail remains unavailable while retained stage order stays inspectable',()=>{
 const j=journey();const m=model(j,[{journeyId:j.id,evidenceId:2,type:'RetainedMarker',severity:17}]);const out=html(j,{model:m,ontrace(){}});assert.match(out,/Evidence unavailable/);assert.match(out,/RetainedMarker/);assert.equal((out.match(/<li>/g)??[]).length,3);assert.doesNotMatch(out,/Open retained trace/);
});
test('long current story preserves every occurrence and grouping keeps original ordinals',()=>{
 const j=journey([...Array(180).fill('Vote'),...Array.from({length:20},(_,i)=>'Tail'+i)]);let out=html(j);assert.equal((out.match(/<li>/g)??[]).length,200);assert.match(out,/<h4>200\. Tail19/);out=html(j,{group:true});assert.equal((out.match(/<li>/g)??[]).length,21);assert.match(out,/×180 Occurrences/);assert.match(out,/Evidence details · 1–180/);assert.match(out,/<h4>180\. Vote/);assert.match(out,/<h4>200\. Tail19/);
});
test('recovery and goal-free navigation remain descriptive in both languages',()=>{
 const j=journey();j.navigationOnly=true;j.recovered=true;j.warnings=['Display uses occurrence time and stable receipt tie-break; causal order is unproven.'];for(const language of ['en','fr']){const out=html(j,{language});assert.match(out,language==='fr'?/Rétabli/:/Recovered/);assert.match(out,language==='fr'?/Navigation sans objectif/:/Goal-free navigation/);assert.doesNotMatch(out,/conversion rate|attention measured|failed journey/i);assert.match(out,language==='fr'?/causal|causalité/:/causal order is unproven/)}
});
