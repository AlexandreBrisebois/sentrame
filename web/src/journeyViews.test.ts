import {test} from 'node:test';
import {strict as assert} from 'node:assert';
import {loopProjection,coveragePartial} from './journeyViews.ts';
import {transitionMetrics} from './investigation.ts';
import type {Journey} from './types.ts';
const journey=(id:string,path:string[],gaps:number[]=[])=>({id,steps:path.map((step,i)=>({step,event:step,evidenceId:i,time:i,source:'browser',gapBefore:gaps.includes(i)}))} as Journey);
test('counts repeats and directed returns once per journey',()=>{
 const p=loopProjection([journey('1',['A','B','A','B','B'])]);
 assert.deepEqual(p.repeats.map(x=>[x.stage,x.occurrences,x.journeyIds.length]),[['B',2,1],['A',1,1]]);
 assert.deepEqual(p.returns.map(x=>[x.from,x.to,x.occurrences]),[['A','B',1],['B','A',1]]);
});
test('gaps do not create returns or invented prefixes',()=>{
 const p=loopProjection([journey('1',['A','B','A'],[2])]);
 assert.equal(p.returns.length,0);assert.equal(p.repeats[0].occurrences,1);assert.deepEqual(p.edges.map(e=>[e.from,e.to]),[['A','B']]);assert.deepEqual(p.prefixes[0].path,['A','B']);assert.equal(p.prefixes[0].continued,1);
});
test('prefixes preserve order and disclose limits',()=>{
 const p=loopProjection([journey('1',['A','B','A']),journey('2',['B','A','B']),journey('3',['A','B','C'])],2,1);
 assert.equal(p.prefixes[0].journeyIds.length,2);assert.equal(p.prefixes[0].continued,2);assert.equal(p.omittedJourneys,1);
});

test('transition counts distinguish occurrences, journeys and visits',()=>{
 const a=journey('1',['A','B','A','B']);a.scope='visit';a.tab='tab';
 const b=journey('2',['A','B']);b.scope='visit';b.tab='tab';
 const edge=loopProjection([a,b]).edges.find(e=>e.from==='A'&&e.to==='B')!;
 assert.equal(edge.occurrences,3);assert.equal(edge.journeyIds.length,2);assert.equal(edge.visits,1);
});

test('partial starts and clipped continuations remain explicit',()=>{
 const j=journey('partial',['B','A']);j.partial=true;j.rightCensored=true;
 const p=loopProjection([j]);assert.equal(p.prefixes[0].partial,1);assert.equal(p.prefixes[0].continued,1);assert.equal(p.returns.length,0);
});
test('identical edges do not collapse different ordered paths',()=>{
 const a=journey('a',['A','B','A','C','A']);const b=journey('b',['A','C','A','B','A']);
 const p=loopProjection([a,b]);assert.equal(p.prefixes.length,2);assert.notDeepEqual(p.prefixes[0].path,p.prefixes[1].path);
});

test('coverage reports truncation and expired detail independently for each range',()=>{
 assert.equal(coveragePartial(true,'200','100'),true);
 assert.equal(coveragePartial(false,'99','100'),true);
 assert.equal(coveragePartial(false,'100','100'),false);
 assert.equal(coveragePartial(false,'invalid','100'),true);
});

test('prefix widths count journeys even when branches diverge or are omitted',()=>{
 const p=loopProjection([journey('1',['A','B','A']),journey('2',['A','B','C']),journey('3',['A','C'])],8,1);
 assert.equal(p.prefixCounts.get(JSON.stringify(['A'])),3);
 assert.equal(p.prefixCounts.get(JSON.stringify(['A','B'])),2);
 assert.equal(p.prefixCounts.get(JSON.stringify(['A','B','A'])),1);
 assert.equal(p.omittedJourneys,2);
});
test('active journeys are not shown as ended',()=>{
 const j=journey('1',['A']);j.status='active';assert.equal(loopProjection([j]).prefixes[0].continued,1);
});

test('selected transition metrics preserve self/reverse direction, hidden destinations and partial evidence',()=>{
 const vote=journey('vote',Array(10).fill('Vote'));vote.scope='visit';vote.tab='tab';
 const self=transitionMetrics([vote])[0];assert.deepEqual([self.occurrences,self.journeys,self.visits,self.denominator,self.share],[9,1,1,9,1]);
 const paths=[...Array.from({length:6},(_,i)=>journey('b'+i,['A','B'])),...Array.from({length:4},(_,i)=>journey('c'+i,['A','C']))];
 const complete=transitionMetrics(paths),b=complete.find(e=>e.to==='B')!,c=complete.find(e=>e.to==='C')!;
 assert.deepEqual([b.share,c.share,b.denominator,c.denominator],[.6,.4,10,10]);
 assert.equal(complete.filter(e=>e.to!=='C')[0].denominator,10,'display filtering cannot change source denominator');
 const partial=journey('partial',['A','B','A','C'],[3]);partial.partial=true;
 const edges=transitionMetrics([partial]);assert.deepEqual(edges.map(e=>[e.from,e.to,e.occurrences]).sort(),[['A','B',1],['B','A',1]]);
 assert.equal(edges.find(e=>e.from==='A')!.denominator,1,'missing adjacency is never added to an observed denominator');
 assert.deepEqual(transitionMetrics([journey('single',['A'])]),[],'no outgoing denominator does not fabricate a zero-share edge');
 assert.deepEqual(transitionMetrics([]),[]);
});
