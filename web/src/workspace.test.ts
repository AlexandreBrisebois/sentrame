import test from 'node:test';
import assert from 'node:assert/strict';
import { editEventDefinition, captureReferenceBounds, isCompletionGoal, completionStage, mappingActivationBlocker } from './workspace.ts';
const mapping=(event:string,step=event)=>({event,step,source:'server',kind:'step',key:'',value:''});
test('graph goal stage respects conditioned selector rather than the first same-name mapping',()=>{
 const home={...mapping('pwa.screen.presented','Home'),source:'browser',key:'stage.name',value:'home'};
 const search={...home,step:'Search',value:'search'};
 const active={version:2,goal:search.event,goalSelector:search,mappings:[home,search],fallbackTiming:false};
 assert.equal(completionStage(active),'Search');
 assert.equal(completionStage({...active,goalSelector:{...search,service:'other'}}),'');
 assert.equal(completionStage({...active,goalSelector:undefined}), '');
 assert.equal(completionStage({...active,goalSelector:undefined,mappings:[search]}),'Search');
 assert.equal(completionStage({...active,goal:''}),'');
 assert.equal(completionStage(undefined),'');
});
test('completion goal reflects the exact configured selector, not every same-name event',()=>{
 const goal={...mapping('saved'),service:'api',key:'stage.name',value:'planner'};
 const active={version:2,goal:'saved',goalSelector:goal,mappings:[goal],fallbackTiming:false};
 assert.equal(isCompletionGoal(active,goal),true);
 assert.equal(isCompletionGoal(active,{...goal,value:'search'}),false);
 assert.equal(isCompletionGoal(active,{...goal,source:'browser'}),false);
 assert.equal(isCompletionGoal({...active,goal:''},goal),false);
 assert.equal(isCompletionGoal({...active,goalSelector:undefined},goal),true);
});
test('a received unmapped event can become a goal without replacing other mappings',()=>{
 const active={version:1,goal:'',mappings:[mapping('opened')],fallbackTiming:false};
 const next=editEventDefinition(active,mapping('saved','Meal added to planner'),-1,true);
 assert.equal(next.goal,'saved');
 assert.equal(isCompletionGoal(next,next.mappings[1]),true);
 assert.equal(next.mappings.length,2);
 assert.equal(active.goal,'');
 assert.equal(mappingActivationBlocker('Saved','saved',{conflict:''},'step'),'');
 assert.equal(mappingActivationBlocker('Saved','saved',{conflict:''},'exit'),'A completion goal cannot also be an exit. Choose Step for a successful milestone.');
});
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

test('event activation waits for automatic validation without requiring or inventing a goal',async()=>{
 const {mappingActivationBlocker}=await import('./workspace.ts');
 const empty={version:1,goal:'',fallbackTiming:false,mappings:[]};
 const next=editEventDefinition(empty,mapping('pwa.app.opened'),-1);
 assert.equal(next.goal,'');
 assert.equal(mappingActivationBlocker(next.mappings[0].step,next.goal,null),'Checking this configuration…');
 assert.equal(mappingActivationBlocker(next.mappings[0].step,next.goal,{conflict:''}),'');
 const withGoal=editEventDefinition(empty,mapping('pwa.app.opened'),-1,true);
 assert.equal(mappingActivationBlocker(withGoal.mappings[0].step,withGoal.goal,null),'Checking this configuration…');
 assert.equal(mappingActivationBlocker(withGoal.mappings[0].step,withGoal.goal,{conflict:'',missingContext:4}), '');
});
test('activation blocks missing stage and server conflicts while preserving an existing goal',async()=>{
 const {mappingActivationBlocker}=await import('./workspace.ts');
 assert.equal(mappingActivationBlocker('  ','done',{conflict:''}),'Enter a stage name before activating.');
 assert.equal(mappingActivationBlocker('Open','done',{conflict:'Conflicting event mappings'}),'Resolve the preview conflict before activating.');
 const active={version:2,goal:'done',fallbackTiming:false,mappings:[mapping('done')]};
 const next=editEventDefinition(active,mapping('pwa.app.opened'),-1);
 assert.equal(next.goal,'done');
 assert.equal(mappingActivationBlocker('Open',next.goal,{conflict:''}),'');
});

test('automatic validation cannot authorize a changed definition or goal',async()=>{
 const {currentMappingPreview}=await import('./workspace.ts');
 const definition={version:1,goal:'',mappings:[mapping('screen')],fallbackTiming:false};
 const preview={definitionFingerprint:JSON.stringify(definition),conflict:''};
 assert.equal(currentMappingPreview(definition,preview),preview);
 assert.equal(currentMappingPreview({...definition,goal:'screen'},preview),null);
 assert.equal(currentMappingPreview({...definition,mappings:[mapping('screen','Changed stage')]},preview),null);
 assert.equal(currentMappingPreview({...definition,mappings:[...definition.mappings,mapping('other')]},preview),null);
 assert.equal(currentMappingPreview(definition,null),null);
});
