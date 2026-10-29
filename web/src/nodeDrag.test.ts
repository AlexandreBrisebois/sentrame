import test from 'node:test';
import assert from 'node:assert/strict';
import {dragPosition} from './nodeDrag.ts';
test('node drag uses screen-pixel threshold and graph coordinates at any zoom',()=>{
 const start={pointer:{x:100,y:100},node:{x:300,y:200},scale:.5};
 assert.equal(dragPosition(start,{x:103,y:101}),null);
 assert.deepEqual(dragPosition(start,{x:110,y:120}),{x:320,y:240});
 assert.deepEqual(dragPosition(start,{x:100,y:100},true),{x:300,y:200});
});
