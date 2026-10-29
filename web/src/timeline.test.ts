import test from 'node:test';
import assert from 'node:assert/strict';
import { clampBrush, brushFromBounds, rangeFromBrush, recordsInBrush } from './timeline.ts';
const data={from:'2026-10-28T10:00:00Z',to:'2026-10-28T11:00:00Z',bins:Array.from({length:4},(_,i)=>({from:new Date(Date.parse('2026-10-28T10:00:00Z')+i*900000).toISOString(),to:new Date(Date.parse('2026-10-28T10:00:00Z')+(i+1)*900000).toISOString(),count:[2,5,0,3][i]})),total:10,truncated:false,unit:'retained records by occurrence time'};
test('brush indices clamp to context and preserve a nonempty interval',()=>{
 assert.deepEqual(clampBrush(-3,99,4),{start:0,end:4});
 assert.deepEqual(clampBrush(3,1,4),{start:3,end:4});
 assert.deepEqual(clampBrush(99,99,4),{start:3,end:4});
 assert.deepEqual(clampBrush(0,0,0),{start:0,end:0});
});
test('brush uses real bin boundaries and UTC ISO values',()=>{
 assert.deepEqual(rangeFromBrush(data,1,3),{from:'2026-10-28T10:15:00.000Z',to:'2026-10-28T10:45:00.000Z'});
 const selection=brushFromBounds(data,'2026-10-28T10:20','2026-10-28T10:41');
 assert.equal(selection.start,1);assert.equal(selection.end,3);assert.equal(selection.clamped,false);
});
test('external bounds are displayed as clamped context without changing their values',()=>{
 const selection=brushFromBounds(data,'2026-10-27T10:00:00Z','2026-10-29T10:00:00Z');
 assert.deepEqual(selection,{start:0,end:4,clamped:true,outside:false});
 const outside=brushFromBounds(data,'2026-10-29T10:00:00Z','2026-10-29T11:00:00Z');
 assert.equal(outside.clamped,true);assert.equal(outside.outside,true);assert.equal(outside.start,3);assert.equal(outside.end,4);
});
test('highlighted counts come from the actual bins rather than interpolated traffic',()=>{
 assert.equal(recordsInBrush(data,0,2),7);assert.equal(recordsInBrush(data,1,3),5);assert.equal(recordsInBrush(data,2,3),0);assert.equal(recordsInBrush(data,0,4),data.total);
});
