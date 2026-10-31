import {test} from 'node:test';
import {strict as assert} from 'node:assert';
import {prefixTree,rankedPathGroups,matchesCriterion,matchingExamples} from './investigation.ts';
import type {Journey} from './types.ts';
const j=(id:string,path:string[],extra:Partial<Journey>={},gaps:number[]=[])=>({id,scope:id,tab:'tab',start:'100',status:'active',path,steps:path.map((step,i)=>({step,event:step,evidenceId:i,time:String(i),source:'browser',gapBefore:gaps.includes(i)})),...extra} as Journey);
test('equal aggregate edges retain different captured order and shared versus exact support',()=>{
 const js=[j('abac',['A','B','A','C']),j('acab',['A','C','A','B']),j('ab',['A','B'],{status:'completed'}),j('abc',['A','B','C'])];
 const rows=rankedPathGroups(js).groups;assert.equal(rows.length,4);const ab=rows.find(g=>g.path.join('/')==='A/B')!;assert.deepEqual([ab.members,ab.prefixSupport,ab.exactSupport,ab.closed],[1,3,1,1]);assert.deepEqual(matchingExamples(js,{kind:'exact',path:ab.path},'','','',0).rows.map(j=>j.id),['ab']);
 const tree=prefixTree(js,new Set());assert.equal(tree.nodes.find(n=>n.path.join('/')==='A/B')!.journeyIds.length,3);assert.ok(tree.nodes.some(n=>n.path.join('/')==='A/C'));
});
test('three-step five-child default permits root and per-parent expansion without exceeding 200 nodes',()=>{
 const js=Array.from({length:9},(_,a)=>Array.from({length:9},(_,b)=>j(`${a}-${b}`,[`R${a}`,`C${b}`,'Third','Fourth']))).flat();
 const initial=prefixTree(js,new Set());assert.equal(initial.rootBranches,9);assert.equal(initial.nodes.filter(n=>n.depth===1).length,5);assert.equal(initial.nodes.filter(n=>n.depth===2&&n.path[0]==='R0').length,5);assert.equal(Math.max(...initial.nodes.map(n=>n.depth)),3);assert.ok(initial.omitted>0);
 const expanded=prefixTree(js,new Set(['[]/*','["R0"]/*','["R0","C0","Third"]']));assert.equal(expanded.nodes.filter(n=>n.depth===1).length,9);assert.equal(expanded.nodes.filter(n=>n.depth===2&&n.path[0]==='R0').length,9);assert.ok(expanded.nodes.some(n=>n.path.join('/')==='R0/C0/Third/Fourth'));assert.ok(expanded.nodes.length<=200);
});
test('32-step cap and global 200-node cap preserve exact long-path drilldown membership',()=>{
 const path=Array.from({length:40},(_,i)=>`Step${i}`),js=[j('long',path),...Array.from({length:299},(_,i)=>j('root'+i,['Root'+i]))];
 const tree=prefixTree(js,new Set(['[]/*']),false,32);assert.equal(tree.nodes.length,200);assert.ok(tree.omitted>0);assert.ok(tree.nodes.every(n=>n.depth<=32));
 const long=rankedPathGroups([js[0]]).groups[0];assert.deepEqual([long.path.length,long.members,long.prefixSupport,long.exactSupport,long.capped],[32,1,1,0,1]);const examples=matchingExamples(js,{kind:'prefix',path:long.path},'','','',0);assert.deepEqual(examples.rows.map(j=>j.id),['long']);assert.equal(examples.rows[0].steps.length,40);
});
test('gap, source partial, display cap, active/censored and explicit closure are separately counted',()=>{
 const long=Array.from({length:40},(_,i)=>`S${i}`);const js=[j('closed',['A'],{status:'completed'}),j('inactive',['B'],{status:'inactive'}),j('active',['C']),j('censored',['D'],{rightCensored:true}),j('gap',['E','X','Y'],{partial:true},[1]),j('late-gap',long,{},[36])];
 const rows=rankedPathGroups(js).groups;assert.equal(rows.find(g=>g.path[0]==='A')!.closed,1);assert.equal(rows.find(g=>g.path[0]==='B')!.inactive,1);assert.equal(rows.find(g=>g.path[0]==='B')!.closed,0);assert.equal(rows.find(g=>g.path[0]==='C')!.active,1);assert.equal(rows.find(g=>g.path[0]==='D')!.censored,1);const gap=rows.find(g=>g.path[0]==='E')!;assert.deepEqual([gap.gap,gap.partial,gap.exactSupport],[1,1,0]);assert.equal(rows.find(g=>g.path[0]==='S0')!.capped,1);const tree=prefixTree([js[5]],new Set(),false,32);assert.equal(tree.nodes.at(-1)!.capped,1);assert.equal(tree.nodes.at(-1)!.gap,0);
});
test('ranked list omission counts retain all source members and browse-all paging reaches uncommon paths',()=>{
 const js=Array.from({length:25},(_,i)=>j(String(i).padStart(2,'0'),['Path'+i]));const ranked=rankedPathGroups(js);assert.deepEqual([ranked.groups.length,ranked.omittedGroups,ranked.omittedMembers,ranked.total],[20,5,5,25]);assert.equal(ranked.groups.reduce((sum,g)=>sum+g.members,0)+ranked.omittedMembers,25);const last=matchingExamples(js,{kind:'all'},'Path24','','',0);assert.deepEqual(last.rows.map(j=>j.id),['24']);
});
test('repeat grouping preserves support for exact group length without merging longer repeats or source gaps',()=>{
 const js=[j('ten',Array(10).fill('Vote')),j('eleven',Array(11).fill('Vote')),j('aba',['A','B','A']),j('gap',['A','A'],{},[1])];const tree=prefixTree(js,new Set(),true);const ten=tree.nodes.find(n=>n.labels[0]==='Vote ×10')!;assert.deepEqual(ten.journeyIds,['ten']);assert.equal(matchesCriterion(js[1],ten.criterion),false);assert.deepEqual(matchingExamples(js,ten.criterion,'','','',0).rows.map(j=>j.id),['ten']);assert.ok(tree.nodes.some(n=>n.labels.join('/')==='A ×1/B ×1/A ×1'));assert.equal(tree.nodes.find(n=>n.labels[0]==='A ×1')!.gap,1);
});
