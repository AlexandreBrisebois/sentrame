import type {TraceOperation,TraceRelationship} from './types';

export function traceRows(operations:TraceOperation[],relationships:TraceRelationship[]) {
 const parents=new Map<number,number>();
 for(const r of relationships)if(r.type==='parent'&&r.state==='resolved'&&r.to!==undefined)parents.set(r.from,r.to);
 const present=new Set(operations.map(o=>o.evidence.id));
 const children=new Map<number,TraceOperation[]>();
 for(const operation of operations){const parent=parents.get(operation.evidence.id);if(parent!==undefined&&present.has(parent)){const group=children.get(parent)??[];group.push(operation);children.set(parent,group)}}
 const rows:{operation:TraceOperation,depth:number}[]=[],seen=new Set<number>();
 const visit=(operation:TraceOperation,depth:number)=>{
  const id=operation.evidence.id;if(seen.has(id))return;seen.add(id);rows.push({operation,depth});
  for(const child of children.get(id)??[])visit(child,depth+1);
 };
 // Root order follows the API; parent structure, never clock proximity, groups children.
 for(const operation of operations){const parent=parents.get(operation.evidence.id);if(parent===undefined||!present.has(parent))visit(operation,0)}
 // Malformed cycles remain bounded and visible without inventing a missing root.
 for(const operation of operations)if(!seen.has(operation.evidence.id))visit(operation,0);
 return rows;
}
export function operationBar(operation:TraceOperation,operations:TraceOperation[]) {
 const starts=operations.filter(o=>!o.evidence.fallback&&BigInt(o.evidence.time)>0n).map(o=>BigInt(o.evidence.time));
 if(!starts.length||operation.durationNs===null)return null;
 const origin=starts.reduce((a,b)=>a<b?a:b),ends=operations.filter(o=>o.durationNs!==null).map(o=>BigInt(o.evidence.time)+BigInt(o.durationNs!));
 const end=ends.reduce((a,b)=>a>b?a:b,origin),range=end-origin;if(range<=0n)return null;
 return {left:Number((BigInt(operation.evidence.time)-origin)*10000n/range)/100,width:Number(BigInt(operation.durationNs)*10000n/range)/100};
}
export function contextualMapping(evidence:import('./types').Evidence,suggestion?:string) {
 return {event:evidence.name,step:suggestion===evidence.attrs?.['stage.name']?suggestion:evidence.name,source:evidence.source,kind:'step',key:'',value:'',service:evidence.service,carrier:evidence.signal==='traces'?'span':evidence.signal==='span-events'?'span-event':'log'};
}
