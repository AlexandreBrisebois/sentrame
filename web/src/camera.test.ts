import test from 'node:test';
import assert from 'node:assert/strict';
import {zoomAt} from './camera.ts';
test('pinch zoom preserves the graph point beneath the pointer',()=>{
 const before={x:40,y:-20,k:.8},pointer={x:430,y:270},after=zoomAt(before,1.5,pointer);
 assert.ok(Math.abs(after.k-1.2)<1e-9);
 assert.ok(Math.abs((pointer.x-before.x)/before.k-(pointer.x-after.x)/after.k)<1e-9);
 assert.ok(Math.abs((pointer.y-before.y)/before.k-(pointer.y-after.y)/after.k)<1e-9);
});
test('pinch clamps extreme gestures without losing its focal point',()=>{
 const camera={x:0,y:0,k:1},p={x:100,y:100};
 assert.equal(zoomAt(camera,100000,p).k,3);assert.equal(zoomAt(camera,.000001,p).k,.08);
 assert.deepEqual(zoomAt(camera,NaN,p),camera);
});
