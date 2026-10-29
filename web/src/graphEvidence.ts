import type {Model,Evidence} from './types';
export const hasErrorEvidence=(e:Evidence)=>e.signal!=='metrics'&&(e.severity>=17||!!e.attrs?.['exception.type']||!!e.attrs?.['exception.message']||!!e.attrs?.['exception.stacktrace']);
export const evidenceEdgeKey=(from:string,to:string)=>JSON.stringify([from,to]);
export function errorEvidence(model:Model){
 const edges=new Set<string>(),nodes=new Set<string>();
 for(const j of model.journeys??[]){
  const steps=j.steps??[],evidence=j.evidence??[];
  const order=new Map(evidence.map((e,i)=>[e.id,i]));
  const add=(a:number,b:number)=>{if(a>=0&&b===a+1&&b<steps.length&&!steps[b].gapBefore)edges.add(evidenceEdgeKey(steps[a].step,steps[b].step));};
  for(let i=0;i<evidence.length;i++){
   const e=evidence[i];if(!hasErrorEvidence(e))continue;
   const mapped=steps.findIndex(s=>s.evidenceId===e.id);
   if(mapped>=0){nodes.add(steps[mapped].step);add(mapped-1,mapped);add(mapped,mapped+1);continue;}
   let before=-1,after=-1;
   for(let k=0;k<steps.length;k++){const at=order.get(steps[k].evidenceId);if(at===undefined)continue;if(at<i)before=k;else if(at>i&&after<0)after=k;}
   add(before,after);
   if(before>=0)nodes.add(steps[before].step);
  }
 }
 return {edges,nodes};
}
