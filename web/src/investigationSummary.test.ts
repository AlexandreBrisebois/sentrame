import {test,before,after} from 'node:test';
import {strict as assert} from 'node:assert';
import {createServer} from 'vite';
import {investigationObservation,safeInvestigationSummary,sanitizeSummaryText,observationCriterion} from './investigation.ts';
import type {Journey,Model} from './types.ts';
const journey=(id:string,path:string[],extra={})=>({id,scope:id,tab:'tab',path,steps:path.map((step,i)=>({step,event:step,evidenceId:i,time:String(i+1),source:'browser'})),evidence:[],partial:false,navigationOnly:true,...extra} as Journey);
const model=(journeys:Journey[]):Model=>({journeys,nodes:[],edges:[],diagnostics:[],unassigned:[],stats:{},policy:'',warnings:[]});
const coverage={from:'2026-10-23T00:00:00Z',to:'2026-10-30T00:00:00Z',scope:'',selectionPolicy:'evidence-window',journeys:1,sourcePartial:false,partialJourneys:0,censoredJourneys:0,definitions:[1],graphNodesOmitted:0,graphEdgesOmitted:0};
test('repeat observation preserves extras, total, reaching support, prevalence denominator and intensity',()=>{
 const observation=investigationObservation(model([journey('one',Array(10).fill('Vote')),journey('two',['Vote'])]),{kind:'repeat',stage:'Vote'});
 assert.deepEqual([observation.extraOccurrences,observation.totalOccurrences,observation.supportJourneys,observation.reachingJourneys,observation.intensity],[9,11,1,2,9]);
 assert.match(observation.description,/9 extra occurrences in 1\/2 journeys/);assert.ok(observation.lowSupport);assert.doesNotMatch(observation.description,/conversion|caus|defect/);
});
test('transition observation preserves nine one one and hidden destination outgoing denominator',()=>{
 const js=[...Array.from({length:6},(_,i)=>journey('b'+i,['A','B'])),...Array.from({length:4},(_,i)=>journey('c'+i,['A','C']))];
 const c={kind:'transition' as const,from:'A',to:'B'};
 const share=investigationObservation(model(js),c,'share');assert.equal(share.count,.6);assert.equal(share.denominator,10);assert.match(share.description,/60.0% outgoing occurrence share; 10 observed/);
 const loop=investigationObservation(model([journey('v',Array(10).fill('Vote'))]),{kind:'transition',from:'Vote',to:'Vote'});assert.equal(loop.count,9);assert.equal(loop.supportJourneys,1);
 assert.equal(investigationObservation(model(js),c,'visits').count,6);
});
test('prefix, exact, return and partial observations preserve captured criteria without inferred gap edges',()=>{
 const js=[journey('a',['A','B','A']),journey('b',['A','B','C'],{partial:true})];
 assert.equal(investigationObservation(model(js),{kind:'prefix',path:['A','B']}).count,2);
 assert.equal(investigationObservation(model(js),{kind:'exact',path:['A','B']}).available,false);
 const returned=investigationObservation(model(js),{kind:'return',from:'B',to:'A'});assert.equal(returned.count,1);assert.equal(returned.denominator,null);assert.ok(returned.partial);
 const gap=journey('gap',['A','B']);gap.steps[1].gapBefore=true;assert.equal(investigationObservation(model([gap]),{kind:'transition',from:'A',to:'B'}).available,false);
});
test('copy projection is bounded allowlisted context and selected observation without credentials, payloads or stale references',()=>{
 const m=model([journey('one',['A','B'],{evidence:[{trace:'0123456789abcdef0123456789abcdef',attrs:{password:'RAW-PAYLOAD-SENTINEL'},body:'RAW-BODY-SENTINEL'}]})]);
 const criterion={kind:'transition' as const,from:'A',to:'B'},observation=investigationObservation(m,criterion);
 const context={app:'Synthetic token=APP-SECRET',population:'live' as const,bounds:{from:coverage.from,to:coverage.to,password:'BOUND-SECRET'},scope:'scope-1',version:'v1 cookie=COOKIE-SECRET',definition:'1',coverage:{...coverage,password:'COVERAGE-SECRET',payload:m}};
 const copied=JSON.parse(safeInvestigationSummary(context,criterion,observation));
 assert.deepEqual(copied.bounds,{from:'2026-10-23T00:00:00.000Z',to:'2026-10-30T00:00:00.000Z'});assert.deepEqual(copied.references.journeys,['one']);assert.equal(copied.references.traces.length,1);assert.equal(copied.observation.denominator,1);assert.equal(copied.coverage.exporterCompleteness,'unknown');
 assert.doesNotMatch(JSON.stringify(copied),/APP-SECRET|COOKIE-SECRET|BOUND-SECRET|COVERAGE-SECRET|RAW-PAYLOAD-SENTINEL|RAW-BODY-SENTINEL/);
 const expired=JSON.parse(safeInvestigationSummary(context,criterion,investigationObservation(model([]),criterion)));assert.equal(expired.references.state,'unavailable');assert.deepEqual(expired.references.journeys,[]);assert.deepEqual(expired.references.traces,[]);assert.equal(expired.observation.available,false);
 assert.equal(sanitizeSummaryText('Bearer abc.def password=foo https://example.test/secret\nhello'),'[redacted] [redacted] [URL omitted] hello');
});
test('copy keeps independent source and list limits separate from denominators and labels stale/paused/unavailable reads',()=>{
 const m=model([journey('one',['A','B'])]),criterion={kind:'transition' as const,from:'A',to:'B'};
 for(const state of ['current','stale','unavailable','paused']){
  const copied=JSON.parse(safeInvestigationSummary({app:'app',population:'baseline',bounds:coverage,scope:'',version:'',definition:'1',coverage,sourceCoverage:{recordLimit:10000,truncated:true,secret:'SOURCE-SECRET'},evidenceCoverage:{returned:2000,limit:2000,omitted:8000,body:'LIST-RAW'},freshness:{state,lastSuccessfulDisplay:'2026-10-30T01:02:03Z',credential:'READ-SECRET'}},criterion,investigationObservation(m,criterion)));
  assert.equal(copied.observation.denominator,1);assert.deepEqual(copied.independentCoverage.source,{recordLimit:10000,truncated:true,exporterCompleteness:'unknown'});assert.deepEqual(copied.independentCoverage.receivedRecordList,{population:'current-query',returned:2000,limit:2000,omitted:8000});assert.equal(copied.independentCoverage.freshness.state,state);assert.equal(copied.independentCoverage.freshness.lastSuccessfulDisplay,'2026-10-30T01:02:03.000Z');assert.doesNotMatch(JSON.stringify(copied),/SOURCE-SECRET|LIST-RAW|READ-SECRET/);
 }
 const unknown=JSON.parse(safeInvestigationSummary({app:'app',population:'current',bounds:null,scope:'',version:'',definition:'',coverage:null},criterion,investigationObservation(model([]),criterion)));assert.equal(unknown.independentCoverage.source.state,'unavailable');assert.equal(unknown.independentCoverage.receivedRecordList.state,'unavailable');assert.equal(unknown.independentCoverage.freshness.state,'unavailable');
});
test('summary selection follows explicit criterion, retained edge, then stage without another population fallback',()=>{
 assert.deepEqual(observationCriterion(null,'A','["B","C"]'),{kind:'transition',from:'B',to:'C'});assert.deepEqual(observationCriterion(null,'A','invalid'),{kind:'stage',stage:'A'});assert.deepEqual(observationCriterion({kind:'repeat',stage:'D'},'A','["B","C"]'),{kind:'repeat',stage:'D'});
});
let server:Awaited<ReturnType<typeof createServer>>,component:any,render:any;
before(async()=>{server=await createServer({server:{middlewareMode:true,hmr:false},appType:'custom'});component=(await server.ssrLoadModule('/src/InvestigationSummary.svelte')).default;render=(await server.ssrLoadModule('svelte/server')).render});after(async()=>{await server?.close()});
test('EN FR summary UI states selected units, small support, partial coverage and neutral definitions',()=>{
 for(const language of ['en','fr']){
  const out=render(component,{props:{model:model([journey('one',Array(10).fill('Vote'),{partial:true})]),criterion:{kind:'repeat',stage:'Vote'},population:'live',coverage,language,oncopy(){},onexamples(){}}}).body;
  assert.match(out,/9/);assert.match(out,/1\/1/);assert.match(out,language==='en'?/Low support/:/Support faible/);assert.match(out,language==='en'?/Partial evidence/:/Preuves partielles/);assert.match(out,language==='en'?/Goal-free navigation/:/navigation sans objectif/);assert.match(out,language==='en'?/Copy investigation summary/:/Copier le résumé/);
 }
});
