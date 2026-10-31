import {test} from 'node:test';
import {strict as assert} from 'node:assert';
import {contextKey,matchesCriterion,prefixTree} from './investigation.ts';
import {defaultInvestigationState,returnAllowed,reconcileOrder} from './investigationNavigation.ts';
import type {Journey} from './types.ts';
const journey=(id:string,path:string[],gap=-1)=>({id,path,start:'1',steps:path.map((step,i)=>({step,evidenceId:i,gapBefore:i===gap})),evidence:[{trace:'trace-a'}]} as Journey);
const context={app:'app',generation:'1',query:'scope=one',population:'baseline' as const};
test('grouped prefix examples equal displayed support, including final run boundary and source gaps',()=>{
 const js=[journey('ten',Array(10).fill('Vote')),journey('eleven',Array(11).fill('Vote')),journey('next',[...Array(10).fill('Vote'),'Done']),journey('gap',Array(11).fill('Vote'),10)];
 for(const node of prefixTree(js,new Set(),true).nodes){
  assert.ok(node.criterion,'displayed prefix carries its exact matching rule');
  assert.deepEqual(js.filter(j=>matchesCriterion(j,node.criterion)).map(j=>j.id),node.journeyIds);
 }
 const ten=prefixTree(js,new Set(),true).nodes.find(n=>n.labels[0]==='Vote ×10')!;
 assert.deepEqual(ten.journeyIds,['ten','next','gap']);
 assert.deepEqual(ten.path,Array(10).fill('Vote'));
});
test('transient state survives serializable remount without storing raw evidence; population state is independent',()=>{
 const live=defaultInvestigationState(),baseline=defaultInvestigationState();
 live.examples.search='Vote';live.examples.page=2;live.examples.chosen='journey';live.examples.story.expanded=['evidence:4'];live.repeatSort='prevalence';live.maxDepth=8;
 assert.deepEqual(structuredClone(live),live);assert.equal(baseline.examples.chosen,'');assert.equal(baseline.examples.story.expanded.length,0);
});
test('trace return rejects another app, generation, query, population, expired target and changed criterion',()=>{
 const j=journey('one',['A','B']);
 const frame={context:contextKey(context),view:'matrix',criterion:{kind:'transition' as const,from:'A',to:'B'},journey:'one',trace:'trace-a',focus:'trace-button',scroll:20};
 assert.ok(returnAllowed(frame,context,[j]));
 for(const altered of [{app:'other'},{generation:'2'},{query:'scope=two'},{population:'live' as const}])assert.equal(returnAllowed(frame,{...context,...altered},[j]),false);
 assert.equal(returnAllowed(frame,context,[]),false);
 assert.equal(returnAllowed(frame,context,[{...j,evidence:[]}]),false);
 assert.equal(returnAllowed(frame,context,[journey('one',['C','B'])]),false);
});
test('prefix order pins existing siblings, appends arrivals and changes only on explicit reorder',()=>{
 const first=reconcileOrder([],['A','B']);
 assert.deepEqual(reconcileOrder(first,['B','C','A']),['A','B','C']);
 assert.deepEqual(reconcileOrder(first,['B','C']),['B','C']);
 assert.deepEqual(reconcileOrder([],['B','A']),['B','A']);
});
test('same-name stages and colliding journey IDs cannot restore across app, generation, query or population',()=>{
 const j=journey('colliding',['A','B','A']);
 const criteria=[{kind:'stage',stage:'A'},{kind:'transition',from:'A',to:'B'},{kind:'repeat',stage:'A'},{kind:'return',from:'B',to:'A'},{kind:'prefix',path:['A','B']},{kind:'exact',path:['A','B','A']}] as const;
 for(const criterion of criteria){
  const frame={context:contextKey(context),view:'sequence',criterion:{...criterion,...('path' in criterion?{path:[...criterion.path]}:{})},journey:j.id,trace:'trace-a',focus:'',scroll:0};
  assert.ok(returnAllowed(frame,context,[j]),criterion.kind);
  for(const altered of [{app:'other'},{generation:'replacement'},{query:'scope=other'},{population:'current' as const},{population:'live' as const}]){
   assert.equal(returnAllowed(frame,{...context,...altered},[j]),false,criterion.kind+' must retain its complete context');
  }
 }
});
test('trace-only returns require retained unassigned authority and cannot borrow a journey or another population',()=>{
 const frame={context:contextKey(context),view:'matrix',criterion:null,journey:'',trace:'trace-a',focus:'',scroll:0};
 const unassigned=[{trace:'trace-a'}] as Parameters<typeof returnAllowed>[3];
 assert.ok(returnAllowed(frame,context,[],unassigned));
 assert.equal(returnAllowed(frame,context,[journey('one',['A'])],[]),false);
 assert.equal(returnAllowed(frame,{...context,population:'live'},[],unassigned),false);
 assert.equal(returnAllowed({...frame,criterion:{kind:'stage',stage:'A'}},context,[],unassigned),false);
 assert.equal(returnAllowed({...frame,trace:''},context,[],unassigned),false);
});
