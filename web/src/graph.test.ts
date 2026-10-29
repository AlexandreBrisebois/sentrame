import test from 'node:test';
import assert from 'node:assert/strict';
import { graphLayout, edgeKey, lensEvidence, transitions } from './graph.ts';
import type { Model, Journey, MapNode } from './types.ts';
const node=(id:string,count=10)=>({id,label:id,occurrences:count,journeys:count,visits:count,stops:0,pending:0,timing:[]}) as MapNode;
const journey=(id:string,path:string[],extra={})=>({id,path,steps:path.map(step=>({step})),definition:1,version:'1',status:'completed',evidence:[],warnings:[],partial:false,...extra}) as Journey;
const model=(journeys:Journey[])=>({journeys,nodes:[],edges:[],stats:{},unassigned:[],diagnostics:[],warnings:[],policy:''}) as Model;
test('horizontal hierarchy orders directed progression and retains cycles',()=>{
 const r=graphLayout([node('A'),node('B'),node('C')],[{from:'A',to:'B',occurrences:10},{from:'B',to:'C',occurrences:10},{from:'B',to:'B',occurrences:2}],'hierarchy');
 assert.ok(r.nodes[0].x<r.nodes[1].x);assert.ok(r.nodes[1].x<r.nodes[2].x);assert.equal(r.edges.length,3);assert.ok(r.edges[2].returning);assert.ok(r.edges.every(e=>!e.path.includes('NaN')));
});
test('force positions are deterministic and do not move with live count changes',()=>{
 const links=[{from:'A',to:'B',occurrences:1}];
 const a=graphLayout([node('A'),node('B')],links,'force');
 const b=graphLayout([node('A',900),node('B',120)],links.map(e=>({...e,occurrences:50})),'force');
 assert.deepEqual(a.nodes.map(n=>[n.x,n.y]),b.nodes.map(n=>[n.x,n.y]));assert.notEqual(a.nodes[0].x,a.nodes[1].x);
});
test('Sankey bands and bars share one linear volume scale, including returns',()=>{
 const r=graphLayout([node('A',10),node('B',5),node('C',5)],[{from:'A',to:'B',occurrences:5},{from:'A',to:'C',occurrences:5},{from:'C',to:'A',occurrences:2}],'sankey');
 const a=r.nodes.find(n=>n.id==='A')!,b=r.nodes.find(n=>n.id==='B')!;
 assert.equal(a.height/b.height,2);assert.equal(r.edges[0].width+r.edges[1].width,a.height);assert.ok(r.edges.every(e=>e.width>0&&Number.isFinite(e.width)));assert.equal(r.edges.length,3);
});
test('window gaps never fabricate adjacency in graph highlights',()=>{
 const j=journey('gap',['A','B'],{steps:[{step:'A'},{step:'B',gapBefore:true}],recovered:true});
 assert.deepEqual(transitions(j),[]);assert.equal(lensEvidence(model([j]),'recovery').edges.size,0);
});
test('rarity is a full comparable sequence, not an uncommon link or low sample guess',()=>{
 const js=Array.from({length:19},(_,i)=>journey(String(i),['A','B','C']));js.push(journey('rare',['A','B','B','C']));
 const r=lensEvidence(model(js),'rare');assert.deepEqual([...r.journeyIds],['rare']);assert.ok(r.edges.has(edgeKey('B','B')));assert.ok(r.edges.has(edgeKey('A','B')));
 assert.equal(lensEvidence(model(js.slice(1)),'rare').journeyIds.size,0);
 assert.equal(lensEvidence(model([...js.slice(0,19),journey('other',['A','B','B','C'],{version:'2'})]),'rare').journeyIds.size,0);
});
test('exceptions and recovery are evidence lenses, without changing journey outcomes',()=>{
 const j=journey('recovered',['A','B'],{status:'completed',recovered:true,evidence:[{severity:17,attrs:{}}]});
 assert.ok(lensEvidence(model([j]),'exceptions').journeyIds.has(j.id));assert.ok(lensEvidence(model([j]),'recovery').journeyIds.has(j.id));assert.equal(j.status,'completed');
});
test('Sankey comparison keeps each period volume and ports separate on a common scale',async()=>{
 const {sankeyLayer}=await import('./graph.ts');
 const ns=[node('A',10),node('B',10),node('C',10)];
 const baseline=[{from:'A',to:'B',occurrences:10}],live=[{from:'A',to:'C',occurrences:10}];
 const envelope=graphLayout(ns,[...baseline,...live],'sankey');
 const b=sankeyLayer(envelope.nodes,[node('A',10),node('B',10)],baseline,envelope.scale);
 const l=sankeyLayer(envelope.nodes,[node('A',10),node('C',10)],live,envelope.scale,8);
 assert.equal(b.nodes.find(n=>n.id==='A')!.height,10*envelope.scale);
 assert.equal(l.nodes.find(n=>n.id==='A')!.height,10*envelope.scale);
 assert.equal(b.edges[0].width,10*envelope.scale);
 assert.equal(l.edges[0].width,10*envelope.scale);
 assert.equal(b.edges[0].from,'A');assert.equal(b.edges[0].to,'B');assert.equal(l.edges[0].to,'C');
 assert.ok(b.nodes.find(n=>n.id==='A')!.height<envelope.nodes.find(n=>n.id==='A')!.height);
 assert.equal(b.edges.length,1);assert.equal(l.edges.length,1);
});
test('skipped-column connections reserve external channels clear of intermediate nodes',()=>{
 const nodes=[node('A'),node('B'),node('C')],edges=[{from:'A',to:'B',occurrences:5},{from:'B',to:'C',occurrences:5},{from:'A',to:'C',occurrences:5}];
 for(const layout of ['hierarchy','sankey'] as const){
  const r=graphLayout(nodes,edges,layout),middle=r.nodes.find(n=>n.id==='B')!,skip=r.edges.find(e=>e.from==='A'&&e.to==='C')!;
  assert.ok(skip.channelY!==undefined);assert.ok(skip.channelY!+skip.width/2<middle.y-middle.height/2);
  assert.ok(skip.channelY!-skip.width/2>=0,'channel remains inside bounds');
  assert.ok(r.height>Math.max(...r.nodes.map(n=>n.y+n.height/2)));
 }
});
function pathTangents(path:string){
 let cursor=[0,0];const result:{start:number[];end:number[]}[]=[];
 for(const match of path.matchAll(/([MLQC])([^MLQC]+)/g)){
  const command=match[1],v=match[2].trim().split(/[,\s]+/).map(Number),to=v.slice(-2);
  if(command==='M'){cursor=to;continue;}
  const first=command==='L'?to:v.slice(0,2),last=command==='L'?cursor:command==='Q'?v.slice(0,2):v.slice(2,4);
  const start=[first[0]-cursor[0],first[1]-cursor[1]],end=[to[0]-last[0],to[1]-last[1]];
  if(Math.hypot(...start)>1e-7&&Math.hypot(...end)>1e-7)result.push({start,end});cursor=to;
 }
 return result;
}
test('external lanes and return loops have continuous tangents without square elbows',()=>{
 const ns=[node('A',100),node('B',100),node('C',100)],es=[{from:'A',to:'B',occurrences:40},{from:'B',to:'C',occurrences:40},{from:'A',to:'C',occurrences:60},{from:'B',to:'B',occurrences:60}];
 for(const layout of ['hierarchy','sankey'] as const){
  const result=graphLayout(ns,es,layout);
  for(const edge of result.edges.filter(e=>e.channelY!==undefined)){
   const pieces=pathTangents(edge.path);assert.ok(pieces.length>1);
   for(let i=1;i<pieces.length;i++){
    const a=pieces[i-1].end,b=pieces[i].start,normal=Math.hypot(...a)*Math.hypot(...b);
    assert.ok(Math.abs(a[0]*b[1]-a[1]*b[0])/normal<1e-7,'connection tangent stays continuous across curve/line joins');
    assert.ok(a[0]*b[0]+a[1]*b[1]>0,'joins do not double back into elbows');
   }
  }
 }
});

test('manual node positions reroute all views without changing evidence or traffic',async()=>{
 const {repositionGraph,sankeyLayer}=await import('./graph.ts');
 const ns=[node('A'),node('B'),node('C')],es=[{from:'A',to:'B',occurrences:5},{from:'B',to:'B',occurrences:5},{from:'B',to:'C',occurrences:5}];
 for(const layout of ['hierarchy','force','sankey'] as const){
  const original=graphLayout(ns,es,layout),b=original.nodes.find(n=>n.id==='B')!;
  const moved=repositionGraph(original,{B:{x:b.x+80,y:b.y+100}},layout),after=moved.nodes.find(n=>n.id==='B')!;
  assert.equal(after.x,layout==='force'?b.x+80:b.x);assert.equal(after.y,b.y+100);
  assert.equal(moved.edges.length,es.length);assert.deepEqual(moved.edges.map(e=>e.occurrences),original.edges.map(e=>e.occurrences));
  assert.ok(moved.edges.every((e,i)=>e.path!==original.edges[i].path));
  assert.ok(moved.edges.every(e=>!e.path.includes('NaN')));assert.deepEqual(original.nodes.find(n=>n.id==='B'),b);
  assert.deepEqual(repositionGraph(original,{},layout),original,'reset restores automatic layout');
  const bounded=repositionGraph(original,{B:{x:-10000,y:-10000}},layout);
  assert.ok(bounded.nodes.find(n=>n.id==='B')!.x>0);assert.ok(bounded.nodes.find(n=>n.id==='B')!.y>0);
  assert.ok(bounded.edges.every(e=>e.channelY===undefined||e.channelY-e.width/2>=0));
  if(layout==='sankey'){
   const layer=sankeyLayer(moved.nodes,ns,es,moved.scale,8);
   assert.equal(layer.nodes.find(n=>n.id==='B')!.y,after.y);assert.equal(layer.edges[0].width,es[0].occurrences*moved.scale);
  }
 }
});
