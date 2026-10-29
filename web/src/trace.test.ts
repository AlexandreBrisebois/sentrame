import {test} from 'node:test';
import assert from 'node:assert/strict';
import {traceRows} from './trace.ts';

test('trace tree rows use explicit resolved parents and preserve exact native time',()=>{
 const ops=[{evidence:{id:1,span:'p',time:'1793276003272967001',end:'1793276004272967002'},durationNs:'1000000001',timing:'complete'},{evidence:{id:2,span:'c'},durationNs:null,timing:'right-censored'},{evidence:{id:3,span:'s'},durationNs:null,timing:'unavailable'}];
 const rows=traceRows(ops as any,[{type:'parent',from:2,to:1,state:'resolved'}] as any);
 assert.equal(rows[1].depth,1);assert.equal(rows[2].depth,0);assert.equal(rows[0].operation.durationNs,'1000000001');
});
test('ambiguous and cyclic parents never manufacture a tree chain',()=>{
 const ops=[{evidence:{id:1}},{evidence:{id:2}}];
 assert.deepEqual(traceRows(ops as any,[{type:'parent',from:1,to:2,state:'ambiguous'}] as any).map(x=>x.depth),[0,0]);
 assert.ok(traceRows(ops as any,[{type:'parent',from:1,to:2,state:'resolved'},{type:'parent',from:2,to:1,state:'resolved'}] as any).every(x=>x.depth<=2));
});

test('business mapping starts with retained service, carrier and trusted producer, without custom context',async()=>{
 const {contextualMapping}=await import('./trace.ts');
 const mapping=contextualMapping({name:'vote.persist',service:'voting-api',source:'server',signal:'traces',attrs:{}} as any);
 assert.equal(mapping.service,'voting-api');assert.equal(mapping.carrier,'span');assert.equal(mapping.source,'server');assert.equal(mapping.key,'');
});

test('stage hint stays evidence until an explicit suggestion is accepted into a literal rule',async()=>{
 const {contextualMapping}=await import('./trace.ts');const e={name:'vote.committed',service:'voting-api',source:'server',signal:'logs',attrs:{'stage.name':'voting.persist'}} as any;
 assert.equal(contextualMapping(e).step,'vote.committed');
 assert.equal(contextualMapping(e,'voting.persist').step,'voting.persist');
 assert.equal(contextualMapping(e,'voting.persist').key,'');
});

test('tree groups explicit children beneath their parent despite skewed start order',()=>{
 const ops=[{evidence:{id:2,time:'1'}},{evidence:{id:3,time:'2'}},{evidence:{id:1,time:'3'}}];
 const rows=traceRows(ops as any,[{type:'parent',from:2,to:1,state:'resolved'}] as any);
 assert.deepEqual(rows.map(r=>[r.operation.evidence.id,r.depth]),[[3,0],[1,0],[2,1]]);
 assert.equal(rows[2].operation.evidence.time,'1');
});
