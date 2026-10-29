import test from 'node:test';
import assert from 'node:assert/strict';
import {errorEvidence,evidenceEdgeKey} from './graphEvidence.ts';
import type {Model,Journey} from './types.ts';
const journey=(steps:any[],evidence:any[])=>({id:'j',status:'completed',steps,path:steps.map(s=>s.step),evidence}) as Journey;
const model=(j:Journey)=>({journeys:[j]}) as Model;
test('error color stays local to observed error-related transitions, without failing recovery',()=>{
 const j=journey([{step:'Start',evidenceId:1},{step:'Error',evidenceId:2},{step:'Recovered',evidenceId:3},{step:'Done',evidenceId:4}],[{id:1,severity:0},{id:2,severity:17},{id:3,severity:0},{id:4,severity:0}]);
 const errors=errorEvidence(model(j));
 assert.ok(errors.edges.has(evidenceEdgeKey('Start','Error')));assert.ok(errors.edges.has(evidenceEdgeKey('Error','Recovered')));
 assert.equal(errors.edges.has(evidenceEdgeKey('Recovered','Done')),false);assert.equal(j.status,'completed');
});
test('unmapped error context marks its observed span but never bridges a window gap',()=>{
 const steps=[{step:'A',evidenceId:1},{step:'B',evidenceId:3}];const evidence=[{id:1},{id:2,severity:17},{id:3}];
 assert.ok(errorEvidence(model(journey(steps,evidence))).edges.has(evidenceEdgeKey('A','B')));
 assert.equal(errorEvidence(model(journey([steps[0],{...steps[1],gapBefore:true}],evidence))).edges.size,0);
 assert.equal(errorEvidence(model(journey(steps,[{id:1},{id:2,severity:13},{id:3}]))).edges.size,0);
});
test('standard message-only exceptions qualify while metric attributes do not become error paths',()=>{
 const steps=[{step:'A',evidenceId:1},{step:'B',evidenceId:3}];
 const evidence=[{id:1},{id:2,signal:'span-events',attrs:{'exception.message':'Controlled diagnostic'}},{id:3}];
 assert.ok(errorEvidence(model(journey(steps,evidence))).edges.has(evidenceEdgeKey('A','B')));
 assert.equal(errorEvidence(model(journey(steps,[{id:1},{id:2,signal:'metrics',severity:17,attrs:{'exception.type':'MetricLabel'}},{id:3}]))).edges.size,0);
});
