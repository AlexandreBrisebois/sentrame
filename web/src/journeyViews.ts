import type {Journey} from './types';

export type RepeatCount = {stage:string; occurrences:number; journeyIds:string[]};
export type ReturnCount = {from:string; to:string; occurrences:number; journeyIds:string[]};
export type Transition = ReturnCount & {visits:number};
export type Prefix = {path:string[]; journeyIds:string[]; continued:number; partial:number};

// The input is the server's captured, filtered journey projection, not a page of evidence.
export function loopProjection(journeys:Journey[], maxSteps=8, maxBranches=20) {
 const repeats=new Map<string,RepeatCount>(), returns=new Map<string,ReturnCount>();
 const edges=new Map<string,Transition & {visitKeys:Set<string>}>(), prefixes=new Map<string,Prefix>(), prefixCounts=new Map<string,number>();
 for(const j of journeys){
  const seen=new Set<string>(), repeated=new Set<string>(), returned=new Set<string>(), edgeSeen=new Set<string>();
  const steps=j.steps??[];
  for(let i=0;i<steps.length;i++){
   const s=steps[i], previous=steps[i-1];
   if(previous&&!s.gapBefore){
    const key=JSON.stringify([previous.step,s.step]);
    const e=edges.get(key)??{from:previous.step,to:s.step,occurrences:0,journeyIds:[],visits:0,visitKeys:new Set<string>()};
    e.occurrences++;if(!edgeSeen.has(key)){e.journeyIds.push(j.id);edgeSeen.add(key)}e.visitKeys.add(JSON.stringify([j.scope,j.tab]));e.visits=e.visitKeys.size;edges.set(key,e);
   }
   if(seen.has(s.step)){
    const r=repeats.get(s.step)??{stage:s.step,occurrences:0,journeyIds:[]};r.occurrences++;
    if(!repeated.has(s.step)){r.journeyIds.push(j.id);repeated.add(s.step)} repeats.set(s.step,r);
    if(previous&&!s.gapBefore&&previous.step!==s.step){const key=JSON.stringify([previous.step,s.step]);const v=returns.get(key)??{from:previous.step,to:s.step,occurrences:0,journeyIds:[]};v.occurrences++;if(!returned.has(key)){v.journeyIds.push(j.id);returned.add(key)}returns.set(key,v)}
   }
   seen.add(s.step);
  }
  const gap=steps.findIndex((s,i)=>i>0&&s.gapBefore);
  const contiguous=gap<0?steps:steps.slice(0,gap);
  const path=contiguous.slice(0,maxSteps).map(s=>s.step); if(!path.length)continue;
  for(let n=1;n<=path.length;n++){const prefix=JSON.stringify(path.slice(0,n));prefixCounts.set(prefix,(prefixCounts.get(prefix)??0)+1)}
  const key=JSON.stringify(path), p=prefixes.get(key)??{path,journeyIds:[],continued:0,partial:0};
  p.journeyIds.push(j.id);if(steps.length>path.length||j.rightCensored||j.status==='active')p.continued++;if(j.partial||steps[0]?.gapBefore)p.partial++;prefixes.set(key,p);
 }
 const byCount=<T extends {occurrences:number}>(a:T,b:T)=>b.occurrences-a.occurrences||JSON.stringify(a).localeCompare(JSON.stringify(b));
 const all=[...prefixes.values()].sort((a,b)=>b.journeyIds.length-a.journeyIds.length||JSON.stringify(a.path).localeCompare(JSON.stringify(b.path)));
 const omitted=all.slice(maxBranches);
 return {prefixCounts,edges:[...edges.values()].map(({visitKeys:_,...edge})=>edge).sort(byCount),repeats:[...repeats.values()].sort(byCount),returns:[...returns.values()].sort(byCount),prefixes:all.slice(0,maxBranches),omittedJourneys:omitted.reduce((n,p)=>n+p.journeyIds.length,0),omittedBranches:omitted.length,maxSteps,maxBranches};
}

// The cutoff and range boundaries use nanoseconds; absent detail is not zero traffic.
export function coveragePartial(truncated:boolean, fromNs:string|number|undefined, detailFrom:string|number|undefined):boolean {
 if(truncated)return true;
 if(fromNs==null||detailFrom==null)return false;
 try{return BigInt(fromNs)<BigInt(detailFrom)}catch{return true}
}
