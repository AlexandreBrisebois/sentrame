import test from 'node:test';
import assert from 'node:assert/strict';
import { editEventDefinition, captureReferenceBounds } from './workspace.ts';
const mapping=(event:string,step=event)=>({event,step,source:'server',kind:'step',key:'',value:''});
test('contextual event edit preserves every other mapping and the operational goal',()=>{
 const active={version:3,goal:'done',fallbackTiming:false,mappings:[mapping('start'),mapping('vote'),mapping('done')]};
 const result=editEventDefinition(active,mapping('vote','Choose'),1);
 assert.equal(result.mappings[1].step,'Choose');assert.equal(result.goal,'done');assert.equal(result.version,3);
 assert.equal(active.mappings[1].step,'vote');assert.deepEqual(result.mappings[0],active.mappings[0]);assert.deepEqual(result.mappings[2],active.mappings[2]);
});
test('only an explicit goal action changes goal; adding mapping keeps old mappings',()=>{
 const active={version:1,goal:'done',fallbackTiming:true,mappings:[mapping('done')]};
 const result=editEventDefinition(active,mapping('new'),-1,true);
 assert.equal(result.goal,'new');assert.equal(result.mappings.length,2);assert.equal(active.mappings.length,1);
});
test('saving a live reference freezes the live window instead of old historical controls',()=>{
 assert.deepEqual(captureReferenceBounds('live','2025-01-01','','2026-10-28T12:00:00Z','2026-10-28T13:00:00Z'),{from:'2026-10-28T12:00:00Z',to:'2026-10-28T13:00:00Z'});
 assert.deepEqual(captureReferenceBounds('historical','2026-10-27T00:00:00Z','','','2026-10-28T13:00:00Z'),{from:'2026-10-27T00:00:00Z',to:'2026-10-28T13:00:00Z'});
});
test('contextual edits accept reactive proxy definitions without cloning the proxy',()=>{
 const active=new Proxy({version:1,goal:'done',fallbackTiming:false,mappings:[mapping('done')]},{});
 assert.equal(editEventDefinition(active,mapping('done','Finished'),0).mappings[0].step,'Finished');
});
test('opening an existing reference preserves its captured nanoseconds until explicitly edited',async()=>{
 const {referenceBound}=await import('./workspace.ts');
 assert.equal(referenceBound('2026-10-28T12:34:56.123456789','2026-10-28T12:34:56.123456789Z'),'2026-10-28T12:34:56.123456789Z');
 assert.equal(referenceBound('2026-10-28T12:35','2026-10-28T12:34:56.123456789Z'),'2026-10-28T12:35:00.000Z');
 assert.equal(referenceBound('','2026-10-28T12:34:56.123456789Z'),'');
});
test('unmapping preserves other rules and prevents removing the active goal',async()=>{
 const {removeEventMapping}=await import('./workspace.ts');
 const active={version:3,goal:'done',fallbackTiming:false,mappings:[mapping('start'),mapping('done')]};
 assert.equal(removeEventMapping(active,0).mappings.length,1);
 assert.equal(active.mappings.length,2);
 assert.throws(()=>removeEventMapping(active,1),/another operational goal/);
});
