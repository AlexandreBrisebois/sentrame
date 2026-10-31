import type { Definition, Mapping } from './types';

export function isCompletionGoal(active:Definition|undefined,mapping:Mapping):boolean {
 if(!active?.goal||active.goal!==mapping.event)return false;
 const selector=active.goalSelector;
 return !selector||(['event','source','service','signal','carrier','key','value'] as const).every(key=>(selector[key]??'')===(mapping[key]??''));
}
export function completionStage(active:Definition|undefined):string {
 const stages=new Set(active?.mappings.filter(mapping=>isCompletionGoal(active,mapping)).map(mapping=>mapping.step));
 return stages.size===1?[...stages][0]:'';
}
// A node edit is a copy of the complete active definition, never an isolated
// replacement that accidentally drops other stages or changes the goal.
export function editEventDefinition(active:Definition, mapping:Mapping, replace:number, makeGoal=false):Definition {
 const next:Definition={...active,mappings:active.mappings.map(m=>({...m}))};
 if(replace>=0&&replace<next.mappings.length)next.mappings[replace]={...mapping};
 else next.mappings.push({...mapping});
 if(makeGoal){next.goal=mapping.event;next.goalSelector={...mapping};}else if(replace>=0&&active.goalSelector&&JSON.stringify(active.goalSelector)===JSON.stringify(active.mappings[replace]))next.goalSelector={...mapping};
 return next;
}
export function captureReferenceBounds(mode:string,from:string,to:string,liveStart:string,now:string){
 return mode==='live'?{from:liveStart,to:now}:{from,to:to||now};
}
export function referenceBound(value:string,captured:string){if(value&&value===captured.replace(/Z$/,''))return captured;return value?new Date(value.endsWith('Z')?value:value+'Z').toISOString():'';}
export function removeEventMapping(active:Definition,index:number):Definition {
 const mapping=active.mappings[index];
 if(!mapping||mapping.event===active.goal||active.mappings.length<=1)throw new Error('Choose another operational goal before removing this mapping.');
 return {...active,mappings:active.mappings.filter((_,i)=>i!==index).map(m=>({...m}))};
}

// Missing journey correlation affects interpretation, not permission to save a
// valid mapping. Keep the disabled action and its visible explanation together.
export function mappingActivationBlocker(stage:string,goal:string,preview:{conflict?:string;missingContext?:number}|null,goalKind?:string):string {
 if(!stage.trim())return 'Enter a stage name before activating.';
 if(goalKind==='exit')return 'A completion goal cannot also be an exit. Choose Step for a successful milestone.';
 if(!preview)return 'Checking this configuration…';
 if(preview.conflict)return 'Resolve the preview conflict before activating.';
 return '';
}

export function currentMappingPreview<T extends {definitionFingerprint?:string}>(definition:Definition,preview:T|null):T|null {
 return preview?.definitionFingerprint===JSON.stringify(definition)?preview:null;
}
