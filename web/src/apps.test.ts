import test from 'node:test';
import assert from 'node:assert/strict';
import {appURL,initialMode,WorkspaceLifetime} from './apps.ts';
test('Journeys is the default; explicit trace and reference links take precedence',()=>{
 assert.equal(initialMode(''),'journey');assert.equal(initialMode('?app=b'),'journey');assert.equal(initialMode('?app=b&trace=abc'),'trace');assert.equal(initialMode('?app=b&reference=1'),'trace');
 assert.equal(appURL('B'),'/?app=B');assert.equal(appURL('B',{trace:'abc',reference:'1'}),'/?app=B&trace=abc&reference=1');
});
test('leaving A aborts requests and rejects late results even if transport ignores abort',async()=>{
 const lifetime=new WorkspaceLifetime();let resolve!:(value:string)=>void;let signal!:AbortSignal;
 const late=lifetime.run(s=>{signal=s;return new Promise<string>(r=>resolve=r)});lifetime.close();assert.equal(signal.aborted,true);resolve('A');await assert.rejects(late,{name:'AbortError'});
 const b=new WorkspaceLifetime();assert.equal(await b.run(async()=> 'B'),'B');
});
test('sign-out invalidates the old request epoch and permits a fresh login',async()=>{
 const lifetime=new WorkspaceLifetime();let finish!:(v:string)=>void;
 const stale=lifetime.run(()=>new Promise<string>(r=>finish=r));lifetime.invalidate();finish('private');await assert.rejects(stale,{name:'AbortError'});assert.equal(await lifetime.run(async()=> 'fresh'),'fresh');
});
test('reference links carry the registration generation as well as app identity',()=>{
 assert.equal(appURL('A',{generation:'g',trace:'t',reference:'1'}),'/?app=A&generation=g&trace=t&reference=1');
});
