import {test} from 'node:test';import {strict as assert} from 'node:assert';
import {CoalescedRefresh,snapshotFreshness,resetNeedsWarning} from './liveInvestigation.ts';
test('bursts retain one in-flight and one pending refresh; close rejects later requests',async()=>{
 let reads=0,active=0,max=0,release:()=>void=()=>{};
 const refresh=new CoalescedRefresh(async()=>{reads++;active++;max=Math.max(max,active);await new Promise<void>(r=>release=r);active--},1);
 for(let i=0;i<1200;i++)refresh.request();await new Promise(r=>setTimeout(r,10));assert.equal(reads,1);for(let i=0;i<1200;i++)refresh.request();assert.equal(reads,1);release();await new Promise(r=>setTimeout(r,10));assert.equal(reads,2);assert.equal(max,1);refresh.close();release();refresh.request();await new Promise(r=>setTimeout(r,5));assert.equal(reads,2);
});
test('transport connection does not imply data freshness and routine reset is not a gap',()=>{assert.equal(snapshotFreshness(0,100),'No successful snapshot');assert.equal(snapshotFreshness(1,6002),'Stale snapshot');assert.ok(!resetNeedsWarning(true,false,false));assert.ok(resetNeedsWarning(true,true,false))});
test('manual and scheduled reads serialize, merge one pending request and retain explicit reconciliation',async()=>{
 const {SerializedRefresh}=await import('./liveInvestigation.ts');
 let active=0,max=0;const calls:{validationOnly:boolean;target:string}[]=[],releases:(()=>void)[]=[];
 const lane=new SerializedRefresh(async(validationOnly,target)=>{active++;max=Math.max(max,active);calls.push({validationOnly,target});await new Promise<void>(r=>releases.push(r));active--});
 const first=lane.request(true,'current');
 const queued=[];for(let i=0;i<1200;i++)queued.push(lane.request(true,'baseline'));queued.push(lane.request(false,'current'));
 assert.equal(calls.length,1);releases.shift()!();await first;await new Promise(r=>setTimeout(r,0));
 assert.deepEqual(calls,[{validationOnly:true,target:'current'},{validationOnly:false,target:'both'}]);assert.equal(max,1);
 releases.shift()!();await Promise.all(queued);assert.equal(calls.length,2);assert.equal(active,0);
});
test('a failed read releases the serial lane and the merged pending action still runs',async()=>{
 const {SerializedRefresh}=await import('./liveInvestigation.ts');let fail:()=>void=()=>{},calls=0;
 const lane=new SerializedRefresh(async()=>{calls++;if(calls===1)await new Promise<void>((_,reject)=>fail=()=>reject(new Error('read failed')))});
 const first=lane.request().catch(e=>e.message),next=lane.request(false,'baseline');fail();
 assert.equal(await first,'read failed');await next;assert.equal(calls,2);await lane.request();assert.equal(calls,3);
});
