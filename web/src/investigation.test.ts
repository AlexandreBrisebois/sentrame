import {test} from 'node:test';
import {strict as assert} from 'node:assert';
import {matchesCriterion,groupOccurrences,transitionMetrics,repeatMetrics,prefixTree,comparisonReason,contextKey,matchingExamples} from './investigation.ts';
import type {Journey} from './types.ts';
const j=(id:string,path:string[],gaps:number[]=[])=>({id,scope:id,tab:'tab',start:'100',status:'active',steps:path.map((step,i)=>({step,event:step,evidenceId:i,time:String(i*1000000000),source:'browser',gapBefore:gaps.includes(i)})),path} as Journey);
test('criterion membership is exact and respects direct gaps, repeats, returns and prefixes',()=>{
 const a=j('a',['A','B','A']);assert.ok(matchesCriterion(a,{kind:'return',from:'B',to:'A'}));assert.ok(matchesCriterion(a,{kind:'repeat',stage:'A'}));assert.ok(!matchesCriterion(a,{kind:'repeat',stage:'B'}));
 assert.ok(matchesCriterion(a,{kind:'prefix',path:['A','B']}));assert.ok(!matchesCriterion(a,{kind:'exact',path:['A','B']}));assert.ok(!matchesCriterion(j('gap',['A','B'],[1]),{kind:'transition',from:'A',to:'B'}));
 assert.notEqual(contextKey({app:'a',generation:'1',query:'x',population:'live'}),contextKey({app:'a',generation:'1',query:'x',population:'baseline'}));
});
test('metrics retain hidden destinations in denominators and preserve nine/one/one',()=>{
 const rows=transitionMetrics([j('loop',Array(10).fill('Vote'))]);assert.deepEqual(rows.map(e=>[e.occurrences,e.journeys,e.visits,e.share]),[[9,1,1,1]]);
 const paths=[...Array.from({length:6},(_,i)=>j('b'+i,['A','B'])),...Array.from({length:4},(_,i)=>j('c'+i,['A','C']))];const e=transitionMetrics(paths);assert.equal(e.find(e=>e.to==='B')!.share,.6);assert.equal(e.find(e=>e.to==='C')!.denominator,10);
});
test('prevalence uses reaching journeys; intensity uses affected journeys',()=>{
 const js=[...Array.from({length:100},(_,i)=>j('a'+i,i<20?['A','A']:['A'])),...Array.from({length:10},(_,i)=>j('b'+i,i<8?['B','B']:['B']))];
 const r=repeatMetrics(js);assert.equal(r.find(r=>r.stage==='A')!.prevalence,.2);assert.equal(r.find(r=>r.stage==='B')!.prevalence,.8);
 const vote=repeatMetrics([j('v',Array(10).fill('Vote'))])[0];assert.equal(vote.extras,9);assert.equal(vote.intensity,9);
});
test('grouping retains ordinal evidence and never crosses gaps or nonconsecutive returns',()=>{
 const g=groupOccurrences(j('v',[...Array(10).fill('Vote'),'A','Vote']).steps,true);assert.equal(g[0].occurrences.length,10);assert.equal(g.length,3);assert.equal(groupOccurrences(j('gap',['A','A'],[1]).steps,true).length,2);
});
test('prefix tree shares support and distinguishes closure, gaps and caps',()=>{
 const tree=prefixTree([j('a',['A','B','C']),j('b',['A','B','D']),j('c',['A','E'])],new Set(['["A"]']),false);assert.equal(tree.nodes[0].journeyIds.length,3);assert.equal(tree.nodes.find(n=>n.path.join('/')==='A/B')!.journeyIds.length,2);
 const gap=prefixTree([j('g',['A','B','C'],[2])],new Set(),false);assert.ok(!gap.nodes.some(n=>n.path.includes('C')));
 const dense=prefixTree(Array.from({length:300},(_,i)=>j(String(i),['A',String(i)])),new Set(['["A"]']),false);assert.ok(dense.nodes.length<=200);assert.ok(dense.omitted>0);
});
test('comparison rejects scope, coverage, mixed definitions and zero support',()=>{
 const cov={sourcePartial:false,partialJourneys:0,definitions:[1],scope:''};assert.equal(comparisonReason(cov,cov),null);assert.ok(comparisonReason({...cov,scope:'one'},cov));assert.ok(comparisonReason({...cov,sourcePartial:true},cov));assert.ok(comparisonReason({...cov,definitions:[1,2]},cov));
});
test('examples keep criterion, newest ordering, stable tie break and bounded pages',()=>{
 const js=Array.from({length:45},(_,i)=>({...j(String(i).padStart(2,'0'),['A','A']),start:String(i)}));const p=matchingExamples(js,{kind:'repeat',stage:'A'},'', '', '',0);assert.equal(p.total,45);assert.equal(p.rows.length,20);assert.equal(p.rows[0].id,'44');assert.equal(matchingExamples(js,{kind:'repeat',stage:'B'},'','','',0).total,0);
});

test('repeat membership and grouping retain separate journeys, partial starts and observed-only returns',()=>{
 const partial={...j('partial',['A','B','A']),partial:true};
 const rows=repeatMetrics([partial,j('single',['A']),j('other',['B'])]);
 assert.deepEqual(rows.map(r=>[r.stage,r.affected,r.reaching,r.total,r.extras,r.prevalence,r.intensity]),[['A',1,2,3,1,.5,1]]);
 assert.equal(matchesCriterion(partial,{kind:'return',from:'B',to:'A'}),true);
 assert.equal(matchesCriterion(j('self',['A','A']),{kind:'return',from:'A',to:'A'}),false);
 assert.deepEqual(repeatMetrics([]),[]);assert.deepEqual(repeatMetrics([j('empty',[])]),[]);
 const a=j('a',['Vote']),b=j('b',['Vote']);assert.equal(repeatMetrics([a,b]).length,0);
 assert.deepEqual([a,b].map(x=>groupOccurrences(x.steps,true).length),[1,1]);
 assert.equal(groupOccurrences(j('return',['A','B','A']).steps,true).length,3);
 assert.equal(groupOccurrences(j('default',['Vote','Vote']).steps,false).length,2);
 const group=groupOccurrences(j('long',Array(10).fill('Vote')).steps,true)[0];
 assert.deepEqual(group.occurrences.map(x=>x.evidenceId),Array.from({length:10},(_,i)=>i));assert.deepEqual([group.start,group.end],[0,9]);
});

test('matching examples combine exact criterion with text/status/time filters and stable page boundaries',()=>{
 const stamp=(s:string)=>String(BigInt(Date.parse(s))*1000000n);
 const js=Array.from({length:45},(_,i)=>({...j(String(i).padStart(2,'0'),['A','B','A']),start:stamp(i<3?'2026-10-29T00:00:00Z':'2026-10-30T00:00:00Z'),status:i%2?'completed':'active'}));
 const criterion={kind:'return' as const,from:'B',to:'A'};
 const p=matchingExamples(js,criterion,'a → b','active','2026-10-30T00:00:00Z',0);
 assert.equal(p.total,21);assert.deepEqual(p.rows.map(x=>x.id),Array.from({length:20},(_,i)=>String(4+i*2).padStart(2,'0')));
 const last=matchingExamples(js,criterion,'a → b','active','2026-10-30T00:00:00Z',99);assert.equal(last.start,20);assert.deepEqual(last.rows.map(x=>x.id),['44']);
 assert.equal(matchingExamples(js,{kind:'exact',path:['A','B']},'','','',0).total,0);
 assert.equal(matchingExamples(js,criterion,'missing','','',0).total,0);
 assert.equal(matchingExamples(js,criterion,'','inactive','',0).total,0);
});
