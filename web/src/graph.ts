import type { Model, Journey } from './types';
import {hasErrorEvidence} from './graphEvidence.ts';
export type Layout = 'hierarchy'|'force'|'sankey';
export type Lens = 'none'|'rare'|'endings'|'exceptions'|'recovery';
export type GraphNode = {id:string;label:string;occurrences:number;discovered?:boolean};
export type GraphEdge = {from:string;to:string;occurrences:number};
export type PlacedNode = GraphNode & {x:number;y:number;width:number;height:number};
export type PlacedEdge = GraphEdge & {path:string;width:number;returning:boolean;channelY?:number};
export const edgeKey=(from:string,to:string)=>JSON.stringify([from,to]);
export function transitions(j:Journey):[string,string][]{
 const steps=j.steps?.length?j.steps:j.path.map(step=>({step,gapBefore:false}));
 return steps.slice(1).flatMap((step,i)=>step.gapBefore?[]:[[steps[i].step,step.step] as [string,string]]);
}
export function lensEvidence(model:Model,lens:Lens){
 const nodes=new Set<string>(),edges=new Set<string>(),journeyIds=new Set<string>();
 let chosen:Journey[]=[];
 if(lens==='rare'){
  const eligible=model.journeys.filter(j=>!j.partial&&!j.rightCensored&&j.status!=='active');
  const cohorts=new Map<string,Journey[]>();
  for(const j of eligible){const key=JSON.stringify([j.definition,j.version]);cohorts.set(key,[...(cohorts.get(key)??[]),j]);}
  for(const cohort of cohorts.values()){
   if(cohort.length<20)continue;
   const groups=new Map<string,Journey[]>();
   for(const j of cohort){const sequence=j.steps?.length?j.steps.map(s=>[s.step,!!s.gapBefore]):j.path;const key=JSON.stringify(sequence);groups.set(key,[...(groups.get(key)??[]),j]);}
   for(const group of groups.values())if(group.length/cohort.length<=.1)chosen.push(...group);
  }
 }else if(lens==='exceptions')chosen=model.journeys.filter(j=>j.evidence?.some(e=>e.severity>=13||hasErrorEvidence(e)));
 else if(lens==='recovery')chosen=model.journeys.filter(j=>j.recovered);
 else if(lens==='endings')chosen=model.journeys;
 for(const j of chosen){
  journeyIds.add(j.id);
  if(lens==='endings'){const last=j.steps?.at(-1)?.step??j.path?.at(-1);if(last)nodes.add(last);}
  else{for(const id of j.path??[])nodes.add(id);for(const [a,b] of transitions(j))edges.add(edgeKey(a,b));}
 }
 return {nodes,edges,journeyIds};
}
// Tarjan components keep cycles intact. Longest-path ranks operate on their DAG.
function ranks(ids:string[],edges:GraphEdge[]){
 const index=new Map<string,number>(),low=new Map<string,number>(),stack:string[]=[],open=new Set<string>(),component=new Map<string,number>();let serial=0,components=0;
 const adjacent=new Map(ids.map(id=>[id,edges.filter(e=>e.from===id).map(e=>e.to)]));
 function visit(id:string){index.set(id,serial);low.set(id,serial++);stack.push(id);open.add(id);for(const next of adjacent.get(id)??[]){if(!index.has(next)){visit(next);low.set(id,Math.min(low.get(id)!,low.get(next)!));}else if(open.has(next))low.set(id,Math.min(low.get(id)!,index.get(next)!));}if(low.get(id)===index.get(id)){let next:string;do{next=stack.pop()!;open.delete(next);component.set(next,components);}while(next!==id);components++;}}
 ids.forEach(id=>{if(!index.has(id))visit(id)});
 const links=edges.filter(e=>component.get(e.from)!==component.get(e.to));const depth=new Map<number,number>();
 function level(c:number):number{if(depth.has(c))return depth.get(c)!;const incoming=links.filter(e=>component.get(e.to)===c);const value=incoming.length?Math.max(...incoming.map(e=>level(component.get(e.from)!)+1)):0;depth.set(c,value);return value;}
 return new Map(ids.map(id=>[id,level(component.get(id)!)]));
}
export function graphLayout(input:GraphNode[],inputEdges:GraphEdge[],layout:Layout){
 const ids=new Set(input.map(n=>n.id));const edges=inputEdges.filter(e=>ids.has(e.from)&&ids.has(e.to)&&e.occurrences>0);const rank=ranks([...ids],edges);
 const groups=new Map<number,GraphNode[]>();for(const n of input){const r=rank.get(n.id)??0;groups.set(r,[...(groups.get(r)??[]),n]);}
 // One common scale preserves proportional Sankey traffic, including cycles.
 const volume=(id:string)=>Math.max(1,input.find(n=>n.id===id)?.occurrences??0,edges.filter(e=>e.from===id).reduce((n,e)=>n+e.occurrences,0),edges.filter(e=>e.to===id).reduce((n,e)=>n+e.occurrences,0));
 const columnVolume=Math.max(1,...[...groups.values()].map(ns=>ns.reduce((n,v)=>n+volume(v.id),0)));
 const scale=Math.min(4,440/columnVolume);
 const thickest=Math.max(1,...edges.map(e=>e.occurrences*scale));
 const columnSpacing=layout==='sankey'?Math.max(270,thickest+120):270;
 let nodes:PlacedNode[]=input.map(n=>{const col=groups.get(rank.get(n.id)??0)!;const i=col.findIndex(v=>v.id===n.id);const height=layout==='sankey'?volume(n.id)*scale:layout==='force'?46:56;const previous=col.slice(0,i).reduce((sum,v)=>sum+(layout==='sankey'?volume(v.id)*scale:56)+70,0);return {...n,x:110+(rank.get(n.id)??0)*columnSpacing,y:155+previous+height/2,width:layout==='sankey'?18:layout==='force'?46:170,height};});
 if(layout==='force'&&nodes.length){
  const width=Math.max(760,Math.sqrt(nodes.length)*230),height=Math.max(500,Math.sqrt(nodes.length)*150);
  const sorted=[...nodes].sort((a,b)=>a.id.localeCompare(b.id));
  sorted.forEach((n,i)=>{const angle=i*Math.PI*2/sorted.length;n.x=width/2+Math.cos(angle)*width*.3;n.y=height/2+Math.sin(angle)*height*.3});
  const byId=new Map(nodes.map(n=>[n.id,n]));
  for(let tick=0;tick<240;tick++){
   const forces=new Map(nodes.map(n=>[n.id,{x:(width/2-n.x)*.007,y:(height/2-n.y)*.007}]));
   for(let i=0;i<nodes.length;i++)for(let k=i+1;k<nodes.length;k++){const a=nodes[i],b=nodes[k];let dx=a.x-b.x,dy=a.y-b.y;const d=Math.max(12,Math.hypot(dx,dy));if(!dx&&!dy)dx=1;const push=Math.min(18,15000/(d*d));forces.get(a.id)!.x+=dx/d*push;forces.get(a.id)!.y+=dy/d*push;forces.get(b.id)!.x-=dx/d*push;forces.get(b.id)!.y-=dy/d*push;if(Math.abs(dx)<180&&Math.abs(dy)<95){const sx=(dx<0?-1:1)*(180-Math.abs(dx))*.045,sy=(dy<0?-1:1)*(95-Math.abs(dy))*.045;if(Math.abs(sx)<Math.abs(sy)){forces.get(a.id)!.x+=sx;forces.get(b.id)!.x-=sx;}else{forces.get(a.id)!.y+=sy;forces.get(b.id)!.y-=sy;}}}
   for(const e of edges){if(e.from===e.to)continue;const a=byId.get(e.from)!,b=byId.get(e.to)!,dx=b.x-a.x,dy=b.y-a.y,d=Math.max(1,Math.hypot(dx,dy)),pull=(d-190)*.018;forces.get(a.id)!.x+=dx/d*pull;forces.get(a.id)!.y+=dy/d*pull;forces.get(b.id)!.x-=dx/d*pull;forces.get(b.id)!.y-=dy/d*pull;}
   const cool=.8*(1-tick/270);for(const n of nodes){const f=forces.get(n.id)!;n.x+=Math.max(-12,Math.min(12,f.x))*cool;n.y+=Math.max(-12,Math.min(12,f.y))*cool;}
  }
 }
 const widthOf=(e:GraphEdge)=>layout==='sankey'?e.occurrences*scale:Math.min(4.5,1+Math.sqrt(e.occurrences)*.28);
 const obstacleEdges=layout==='force'?[]:edges.filter(e=>channelNeeded(nodes.find(n=>n.id===e.from)!,nodes.find(n=>n.id===e.to)!,nodes));
 const channelWidth=Math.max(0,...obstacleEdges.map(widthOf));
 const headroom=40+obstacleEdges.reduce((n,e)=>n+widthOf(e)+24,0);
 const minX=Math.min(0,...nodes.map(n=>n.x-n.width/2-(layout==='force'?45:channelWidth+70)));
 const minY=Math.min(0,...nodes.map(n=>n.y-n.height/2-(layout==='force'?110:headroom)));
 nodes=nodes.map(n=>({...n,x:n.x-minX,y:n.y-minY}));
 return repositionGraph({nodes,edges:edges.map(e=>({...e,path:'',width:widthOf(e),returning:false})),width:600,height:410,scale},{},layout);
}
export type NodePositions=Record<string,{x:number;y:number}>;
export type GraphGeometry={nodes:PlacedNode[];edges:PlacedEdge[];width:number;height:number;scale:number};
// Manual placement is presentation state only. Reuse actual evidence and reroute
// connections without rerunning force simulation on every pointer movement.
export function repositionGraph(base:GraphGeometry,positions:NodePositions,layout:Layout):GraphGeometry{
 const minTop=Math.min(...base.nodes.map(n=>n.y-n.height/2));
 const nodes=base.nodes.map(n=>{
  const p=positions[n.id];if(!p)return {...n};
  return {...n,x:layout==='force'?Math.max(n.width/2+45,p.x):n.x,y:Math.max(minTop+n.height/2,p.y)};
 });
 const byId=new Map(nodes.map(n=>[n.id,n])),edges=base.edges,scale=base.scale;
 const widthOf=(e:GraphEdge)=>Math.min(4.5,1+Math.sqrt(e.occurrences)*.28);
 const placed=layout==='force'?edges.map(e=>{
  const a=byId.get(e.from)!,b=byId.get(e.to)!,returning=e.from===e.to||b.x<=a.x,width=widthOf(e);let path='';
  if(e.from===e.to){const x=a.x,y=a.y-a.height/2;path=`M ${x-13} ${y} C ${x-75} ${y-95}, ${x+75} ${y-95}, ${x+13} ${y-3}`;}
  else{const dx=b.x-a.x,dy=b.y-a.y,d=Math.max(1,Math.hypot(dx,dy)),r=26;const sx=a.x+dx/d*r,sy=a.y+dy/d*r,tx=b.x-dx/d*(r+3),ty=b.y-dy/d*(r+3);const reverse=edges.some(v=>v.from===e.to&&v.to===e.from);path=reverse?`M ${sx} ${sy} Q ${(sx+tx)/2-dy/d*30} ${(sy+ty)/2+dx/d*30}, ${tx} ${ty}`:`M ${sx} ${sy} L ${tx} ${ty}`;}
  return {...e,path,width,returning};
 }):routeLayer(nodes,edges,layout,scale,nodes);
 const channelWidth=Math.max(0,...placed.filter(e=>e.channelY!==undefined).map(e=>e.width));
 return {nodes,edges:placed,width:Math.max(600,...nodes.map(n=>n.x+n.width/2+Math.max(140,channelWidth+70))),height:Math.max(410,...nodes.map(n=>n.y+n.height/2+100)),scale};
}
function channelNeeded(a:PlacedNode,b:PlacedNode,obstacles:PlacedNode[]){return b.x<=a.x||obstacles.some(n=>n.id!==a.id&&n.id!==b.id&&n.x>a.x&&n.x<b.x);}
// Skipped ranks and returns use dedicated lanes above every visible rectangle.
// Source/target side corridors stay outside same-column stacked nodes.
function routeLayer(nodes:PlacedNode[],edges:GraphEdge[],layout:Layout,scale:number,obstacles:PlacedNode[]):PlacedEdge[]{
 const byId=new Map(nodes.map(n=>[n.id,n]));const outgoing=new Map<string,number>(),incoming=new Map<string,number>();
 let nextChannel=Math.min(...obstacles.map(n=>n.y-n.height/2))-12;
 return edges.map(e=>{
  const a=byId.get(e.from)!,b=byId.get(e.to)!,width=layout==='sankey'?e.occurrences*scale:Math.min(4.5,1+Math.sqrt(e.occurrences)*.28);
  const sy=layout==='sankey'?a.y-a.height/2+(outgoing.get(a.id)??0)+width/2:a.y,ty=layout==='sankey'?b.y-b.height/2+(incoming.get(b.id)??0)+width/2:b.y;
  outgoing.set(a.id,(outgoing.get(a.id)??0)+width);incoming.set(b.id,(incoming.get(b.id)??0)+width);
  const sx=a.x+a.width/2+(layout==='sankey'?0:2),tx=b.x-b.width/2-(layout==='sankey'?0:5),returning=e.from===e.to||b.x<=a.x;let path='',channelY:number|undefined;
  if(channelNeeded(a,b,obstacles)){
   nextChannel-=width/2+12;channelY=nextChannel;nextChannel-=width/2+12;
   // Circular cubic quarters join with continuous tangents. Their radius
   // exceeds half the ribbon width, avoiding a pinched inner bend on thick flow.
   const desired=Math.max(18,width/2+10),clear=desired+24,right=sx+clear,left=tx-clear;
   const direction=left>right?1:-1;
   const radius=Math.min(desired,(sy-channelY)/2,(ty-channelY)/2,Math.abs(left-right)/2);
   const handle=radius*0.5522847498307936;
   path=`M ${sx} ${sy} L ${right-radius} ${sy} C ${right-radius+handle} ${sy}, ${right} ${sy-radius+handle}, ${right} ${sy-radius} L ${right} ${channelY+radius} C ${right} ${channelY+radius-handle}, ${right+direction*(radius-handle)} ${channelY}, ${right+direction*radius} ${channelY} L ${left-direction*radius} ${channelY} C ${left-direction*(radius-handle)} ${channelY}, ${left} ${channelY+radius-handle}, ${left} ${channelY+radius} L ${left} ${ty-radius} C ${left} ${ty-radius+handle}, ${left+radius-handle} ${ty}, ${left+radius} ${ty} L ${tx} ${ty}`;
  }else{const mid=(sx+tx)/2;path=`M ${sx} ${sy} C ${mid} ${sy}, ${mid} ${ty}, ${tx} ${ty}`;}
  return {...e,path,width,returning,channelY};
 });
}
// Comparison layers share an envelope only for spacing. Each layer uses its own
// actual volumes and independent port offsets; their counts are never combined.
export function sankeyLayer(envelope:PlacedNode[],input:GraphNode[],inputEdges:GraphEdge[],scale:number,barWidth=18){
 const visible=new Set(envelope.map(n=>n.id));
 const ids=new Set(input.filter(n=>visible.has(n.id)).map(n=>n.id));
 const edges=inputEdges.filter(e=>ids.has(e.from)&&ids.has(e.to)&&e.occurrences>0);
 const nodes=envelope.filter(n=>ids.has(n.id)).map(n=>{
  const actual=input.find(v=>v.id===n.id)!;
  const volume=Math.max(actual.occurrences,edges.filter(e=>e.from===n.id).reduce((sum,e)=>sum+e.occurrences,0),edges.filter(e=>e.to===n.id).reduce((sum,e)=>sum+e.occurrences,0));
  return {...n,...actual,width:barWidth,height:volume*scale};
 });
 return {nodes,edges:routeLayer(nodes,edges,'sankey',scale,envelope)};
}
