import type {DiscoveryContext,Mapping} from './types';

export function receivedMappingContexts(event:string,contexts:DiscoveryContext[]):DiscoveryContext[] {
 const groups=new Map<string,DiscoveryContext>();
 for(const context of contexts){
  if(context.name!==event||!['span','span-event','log'].includes(context.carrier))continue;
  const key=JSON.stringify([context.service,context.source,context.signal,context.carrier]);
  const old=groups.get(key);
  groups.set(key,{...context,count:(old?.count??0)+context.count});
 }
 return [...groups.values()];
}
export function applyReceivedContext(mapping:Mapping,context?:DiscoveryContext):Mapping {
 return {...mapping,service:context?.service??'',source:context?.source??'',signal:context?.signal??'',carrier:context?.carrier??'',key:'',value:''};
}
export function initialReceivedMapping(event:string,contexts:DiscoveryContext[]):Mapping {
 const options=receivedMappingContexts(event,contexts);
 return applyReceivedContext({event,step:event,source:'',kind:'step',key:'',value:''},options.length===1?options[0]:undefined);
}
