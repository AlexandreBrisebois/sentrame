import type {Journey,StepOccurrence,Model} from './types';

export type Population='current'|'live'|'baseline';
export type Criterion={kind:'all'}|{kind:'stage'|'repeat';stage:string}|{kind:'transition'|'return';from:string;to:string}|{kind:'prefix'|'exact';path:string[];grouped?:boolean};
export type InvestigationContext={app:string;generation:string;query:string;population:Population};
export type ReturnFrame={context:string;view:string;criterion:Criterion|null;journey:string;focus:string;scroll:number;pageScroll?:number;trace?:string};
export type Metric='occurrences'|'journeys'|'visits'|'share';
export type Preferences={version:1;metric:Metric;stages:string[];sort:string;groupRepeats:boolean};
export type PopulationCoverage={from:string;to:string;scope:string;selectionPolicy:string;journeys:number;sourcePartial:boolean;partialJourneys:number;censoredJourneys:number;definitions:number[];graphNodesOmitted:number;graphEdgesOmitted:number};
export const defaultPreferences=():Preferences=>({version:1,metric:'occurrences',stages:[],sort:'traffic',groupRepeats:false});
export const contextKey=(c:InvestigationContext)=>JSON.stringify([c.app,c.generation,c.query,c.population]);
export function validReturn(frame:ReturnFrame,context:InvestigationContext,journeys:Journey[]){return frame.context===contextKey(context)&&(!frame.journey||journeys.some(j=>j.id===frame.journey))}
export function matchesCriterion(j:Journey,c:Criterion):boolean {
 const steps=j.steps??[];
 if(c.kind==='all')return true;
 if(c.kind==='stage')return steps.some(s=>s.step===c.stage);
 if(c.kind==='repeat')return steps.filter(s=>s.step===c.stage).length>=2;
 if(c.kind==='transition'||c.kind==='return')return steps.some((s,i)=>i>0&&!s.gapBefore&&s.step===c.to&&steps[i-1].step===c.from&&(c.kind==='transition'||c.from!==c.to&&steps.slice(0,i).some(s=>s.step===c.to)));
 if(c.kind!=='prefix'&&c.kind!=='exact')return false;
 if(c.grouped){const expected=groupOccurrences(c.path.map(step=>({step} as StepOccurrence)),true),actual=groupOccurrences(steps,true);if(!expected.every((g,i)=>actual[i]?.stage===g.stage&&actual[i].occurrences.length===g.occurrences.length&&(i===0||!actual[i].gap)))return false;}
 return c.path.every((s,i)=>steps[i]?.step===s&&(i===0||!steps[i].gapBefore))&&(c.kind!=='exact'||steps.length===c.path.length);
}
export type OccurrenceGroup={stage:string;start:number;end:number;occurrences:StepOccurrence[];gap:boolean};
export function groupOccurrences(steps:StepOccurrence[],group:boolean):OccurrenceGroup[]{
 const result:OccurrenceGroup[]=[];
 steps.forEach((s,i)=>{const last=result.at(-1);if(group&&last?.stage===s.step&&!s.gapBefore){last.occurrences.push(s);last.end=i}else result.push({stage:s.step,start:i,end:i,occurrences:[s],gap:!!s.gapBefore})});return result;
}
export function transitionMetrics(journeys:Journey[]){
 const edges=new Map<string,{from:string;to:string;occurrences:number;ids:Set<string>;visits:Set<string>}>();
 for(const j of journeys){const steps=j.steps??[];for(let i=1;i<steps.length;i++){if(steps[i].gapBefore)continue;const from=steps[i-1].step,to=steps[i].step,key=JSON.stringify([from,to]);const e=edges.get(key)??{from,to,occurrences:0,ids:new Set<string>(),visits:new Set<string>()};e.occurrences++;e.ids.add(j.id);e.visits.add(JSON.stringify([j.scope,j.tab]));edges.set(key,e)}}
 const denominators=new Map<string,number>();for(const e of edges.values())denominators.set(e.from,(denominators.get(e.from)??0)+e.occurrences);
 return [...edges.values()].map(e=>({...e,journeyIds:[...e.ids],journeys:e.ids.size,visits:e.visits.size,denominator:denominators.get(e.from)??0,share:denominators.get(e.from)?e.occurrences/denominators.get(e.from)!:null})).sort((a,b)=>b.occurrences-a.occurrences||a.from.localeCompare(b.from)||a.to.localeCompare(b.to));
}
export function repeatMetrics(journeys:Journey[]){
 const rows=new Map<string,{stage:string;total:number;extras:number;reaching:number;journeyIds:string[]}>();
 for(const j of journeys){const counts=new Map<string,number>();for(const s of j.steps??[])counts.set(s.step,(counts.get(s.step)??0)+1);for(const [stage,n]of counts){const r=rows.get(stage)??{stage,total:0,extras:0,reaching:0,journeyIds:[]};r.total+=n;r.reaching++;if(n>1){r.extras+=n-1;r.journeyIds.push(j.id)}rows.set(stage,r)}}
 return [...rows.values()].map(r=>({...r,affected:r.journeyIds.length,prevalence:r.reaching?r.journeyIds.length/r.reaching:null,intensity:r.journeyIds.length?r.extras/r.journeyIds.length:null})).filter(r=>r.affected>0);
}
export type PrefixNode={key:string;path:string[];criterion:Criterion;labels:string[];journeyIds:string[];depth:number;children:number;closed:number;active:number;censored:number;inactive:number;gap:number;capped:number};
export function prefixTree(journeys:Journey[],expanded:Set<string>,group=false,depthLimit=3,branchLimit=5,order:string[]=[]){
 const map=new Map<string,PrefixNode>();
 for(const j of journeys){const groups=groupOccurrences(j.steps??[],group);const gap=groups.findIndex((g,i)=>i>0&&g.gap);const contiguous=gap<0?groups:groups.slice(0,gap);const path:string[]=[];
  for(let i=0;i<Math.min(32,contiguous.length);i++){const g=contiguous[i];path.push(g.stage+(group?` ×${g.occurrences.length}`:''));const key=JSON.stringify(path);const n=map.get(key)??{key,path:(j.steps??[]).slice(0,g.end+1).map(s=>s.step),criterion:{kind:'prefix',path:(j.steps??[]).slice(0,g.end+1).map(s=>s.step),grouped:group},labels:[...path],journeyIds:[],depth:i+1,children:0,closed:0,active:0,censored:0,inactive:0,gap:0,capped:0};n.journeyIds.push(j.id);
   if(i===contiguous.length-1||i===31){if(contiguous.length>32)n.capped++;else if(gap>=0)n.gap++;else if(j.rightCensored)n.censored++;else if(j.status==='active')n.active++;else if(['completed','exited','conflict'].includes(j.capturedStatus||j.status))n.closed++;else n.inactive++}map.set(key,n);
  }
 }
 const children=new Map<string,PrefixNode[]>();for(const n of map.values()){const parent=JSON.stringify(n.labels.slice(0,-1));const list=children.get(parent)??[];list.push(n);children.set(parent,list)}
 const ranks=new Map(order.map((key,i)=>[key,i]));
 for(const [key,list]of children){list.sort((a,b)=>(ranks.get(a.key)??Infinity)-(ranks.get(b.key)??Infinity)||b.journeyIds.length-a.journeyIds.length||a.key.localeCompare(b.key));if(map.has(key))map.get(key)!.children=list.length}
 const visible:PrefixNode[]=[];let omitted=0;
 function visit(key:string){const list=children.get(key)??[];const shown=list.slice(0,expanded.has(key+'/*')?200:branchLimit);omitted+=list.slice(shown.length).reduce((n,x)=>n+x.journeyIds.length,0);for(const n of shown){if(visible.length>=200){omitted+=n.journeyIds.length;continue}visible.push(n);if(n.depth<depthLimit||expanded.has(n.key))visit(n.key)}}visit('[]');
 return {nodes:visible,omitted,rootBranches:children.get('[]')?.length??0,total:journeys.length,order:[...children.values()].flatMap(list=>list.map(n=>n.key))};
}
// Display groups use the same bounded uninterrupted prefix as the flow. Exact
// captured-path support and prefix support remain separate, overlapping counts.
export function rankedPathGroups(journeys:Journey[],limit=20){
 const groups=new Map<string,{path:string[];members:number;closed:number;active:number;censored:number;inactive:number;gap:number;capped:number;partial:number}>();
 for(const j of journeys){const steps=j.steps??[],gap=steps.findIndex((s,i)=>i>0&&s.gapBefore),contiguous=gap<0?steps:steps.slice(0,gap),path=contiguous.slice(0,32).map(s=>s.step);if(!path.length)continue;
  const key=JSON.stringify(path),g=groups.get(key)??{path,members:0,closed:0,active:0,censored:0,inactive:0,gap:0,capped:0,partial:0};g.members++;if(j.partial||steps[0]?.gapBefore)g.partial++;
  if(contiguous.length>32)g.capped++;else if(gap>=0)g.gap++;else if(j.rightCensored)g.censored++;else if(j.status==='active')g.active++;else if(['completed','exited','conflict'].includes(j.capturedStatus||j.status))g.closed++;else g.inactive++;groups.set(key,g);
 }
 const all=[...groups.values()].sort((a,b)=>b.members-a.members||JSON.stringify(a.path).localeCompare(JSON.stringify(b.path)));
 return {groups:all.slice(0,limit).map(g=>({...g,prefixSupport:journeys.filter(j=>matchesCriterion(j,{kind:'prefix',path:g.path})).length,exactSupport:journeys.filter(j=>matchesCriterion(j,{kind:'exact',path:g.path})).length})),omittedGroups:Math.max(0,all.length-limit),omittedMembers:all.slice(limit).reduce((sum,g)=>sum+g.members,0),total:journeys.length};
}
export function matchingExamples(journeys:Journey[],criterion:Criterion,text:string,status:string,from:string,page:number){
 let fromNs:bigint|null=null;try{if(from)fromNs=BigInt(Date.parse(from))*1000000n}catch{}
 const matched=journeys.filter(j=>matchesCriterion(j,criterion)&&(!status||j.status===status)&&(!text||(j.path??[]).join(' → ').toLocaleLowerCase().includes(text.toLocaleLowerCase()))&&(fromNs===null||BigInt(j.start)>=fromNs));
 matched.sort((a,b)=>BigInt(a.start)===BigInt(b.start)?a.id.localeCompare(b.id):BigInt(a.start)>BigInt(b.start)?-1:1);
 const start=Math.min(Math.max(0,page)*20,Math.max(0,Math.ceil(matched.length/20)-1)*20);return {rows:matched.slice(start,start+20),total:matched.length,start};
}
export function comparisonReason(a:Partial<PopulationCoverage>|undefined,b:Partial<PopulationCoverage>|undefined):string|null {
 if(!a||!b)return 'Both populations must be available.';
 if(a.scope||b.scope)return 'A selected visit and app-wide history have different scopes.';
 if(a.sourcePartial||b.sourcePartial||a.partialJourneys||b.partialJourneys)return 'Partial source coverage does not support a difference.';
 if(a.definitions?.length!==1||b.definitions?.length!==1||a.definitions[0]!==b.definitions[0])return 'Stage definitions are missing, mixed or incompatible.';
 return null;
}
export function criterionText(c:Criterion){return c.kind==='all'?'All retained journeys':c.kind==='stage'||c.kind==='repeat'?`${c.kind}: ${c.stage}`:c.kind==='transition'||c.kind==='return'?`${c.kind}: ${c.from} → ${c.to}`:`${c.kind}: ${"path" in c?c.path.join(' → '):""}`}
export function lostEvidence(previous:Model|undefined,next:Model|undefined){if(!previous)return false;const ids=new Set((next?.journeys??[]).flatMap(j=>(j.evidence??[]).map(e=>e.id)).concat((next?.unassigned??[]).map(e=>e.id)));return [...previous.journeys.flatMap(j=>(j.evidence??[]).map(e=>e.id)),...(previous.unassigned??[]).map(e=>e.id)].some(id=>!ids.has(id))}
// Summaries project declared analysis fields only. Raw evidence, attributes, messages,
// credentials, URLs and arbitrary snapshot keys never enter the copy payload.
export type InvestigationObservation={kind:Criterion['kind'];metric:Metric;available:boolean;description:string;unit:string;count:number|null;denominator:number|null;supportJourneys:number;lowSupport:boolean;partial:boolean;reachingJourneys?:number;extraOccurrences?:number;totalOccurrences?:number;intensity?:number|null;references:{state:'available'|'unavailable';total:number;journeys:string[];traces:string[];omitted:number}};
export function sanitizeSummaryText(value:unknown,max=160):string {
 return (typeof value==='string'?value:'').replace(/[\u0000-\u001f\u007f]/g,' ').replace(/(?:bearer\s+)[a-z0-9._~+\/-]+/gi,'[redacted]').replace(/\b(?:password|passwd|secret|token|api[_-]?key|authorization|cookie|csrf)\s*[:=]\s*(?:"[^"]*"|'[^']*'|[^\s,;]+)/gi,'[redacted]').replace(/https?:\/\/[^\s]+/gi,'[URL omitted]').trim().slice(0,max);
}
const safeCount=(n:unknown)=>typeof n==='number'&&Number.isFinite(n)&&n>=0?n:null;
export function observationCriterion(criterion:Criterion|null|undefined,selectedStep='',selectedEdge=''):Criterion {
 if(criterion)return criterion;
 try{const pair=JSON.parse(selectedEdge);if(Array.isArray(pair)&&pair.length===2&&pair.every(x=>typeof x==='string'))return {kind:'transition',from:pair[0],to:pair[1]}}catch{}
 return selectedStep?{kind:'stage',stage:selectedStep}:{kind:'all'};
}
export function investigationObservation(model:Model,criterion:Criterion,metric:Metric='occurrences',coverage?:PopulationCoverage,language='en'):InvestigationObservation {
 if(criterion.kind!=='transition')metric='occurrences';
 const fr=language==='fr',js=model.journeys??[],matched=js.filter(j=>matchesCriterion(j,criterion));
 const partial=!!coverage?.sourcePartial||!!coverage?.partialJourneys||js.some(j=>j.partial||(j.steps??[]).some(s=>s.gapBefore));
 let count:number|null=matched.length,denominator:number|null=js.length,unit=fr?'parcours conservés':'retained journeys',description='',extras:Pick<InvestigationObservation,'reachingJourneys'|'extraOccurrences'|'totalOccurrences'|'intensity'>={};
 const stage=(s:string)=>sanitizeSummaryText(s),path=('path' in criterion?criterion.path:[]).map(stage).join(' → ');
 if(criterion.kind==='transition'||criterion.kind==='return'){
  let edges=transitionMetrics(js);if(criterion.kind==='return')edges=transitionMetrics(matched);
  const edge=edges.find(e=>e.from===criterion.from&&e.to===criterion.to),occurrences=criterion.kind==='return'?matched.reduce((sum,j)=>sum+(j.steps??[]).filter((s,i,steps)=>i>0&&!s.gapBefore&&steps[i-1].step===criterion.from&&s.step===criterion.to&&criterion.from!==criterion.to&&steps.slice(0,i).some(x=>x.step===criterion.to)).length,0):edge?.occurrences??0;
  denominator=criterion.kind==='return'?null:transitionMetrics(js).filter(e=>e.from===criterion.from).reduce((sum,e)=>sum+e.occurrences,0);
  count=metric==='share'?(denominator?occurrences/denominator:null):metric==='journeys'?matched.length:metric==='visits'?new Set(matched.map(j=>JSON.stringify([j.scope,j.tab]))).size:occurrences;
  unit=metric==='share'?(fr?'part des occurrences sortantes':'outgoing occurrence share'):metric==='journeys'?(fr?'parcours concernés':'affected journeys'):metric==='visits'?(fr?'visites concernées':'affected visits'):(fr?'occurrences de transition':'transition occurrences');
  description=`${stage(criterion.from)} → ${stage(criterion.to)}: ${count===null?'—':metric==='share'?(100*count).toFixed(1)+'%':count} ${unit}`;
  if(denominator!==null)description+=fr?` ; ${denominator} occurrences sortantes observées depuis ${stage(criterion.from)}.`:`; ${denominator} observed outgoing occurrences from ${stage(criterion.from)}.`;
 }else if(criterion.kind==='repeat'){
  const row=repeatMetrics(js).find(r=>r.stage===criterion.stage);count=row?.extras??0;denominator=row?.reaching??js.filter(j=>j.steps.some(s=>s.step===criterion.stage)).length;unit=fr?'occurrences supplémentaires':'extra occurrences';
  extras={reachingJourneys:denominator,extraOccurrences:count,totalOccurrences:row?.total??js.reduce((sum,j)=>sum+j.steps.filter(s=>s.step===criterion.stage).length,0),intensity:row?.intensity??null};
  description=fr?`${stage(criterion.stage)} : ${count} occurrences supplémentaires dans ${matched.length}/${denominator} parcours atteignant cette étape ; ${extras.totalOccurrences} occurrences totales. Intensité : ${extras.intensity==null?'—':extras.intensity.toFixed(2)} supplémentaires par parcours répétant.`:`${stage(criterion.stage)}: ${count} extra occurrences in ${matched.length}/${denominator} journeys reaching this stage; ${extras.totalOccurrences} total occurrences. Intensity: ${extras.intensity==null?'—':extras.intensity.toFixed(2)} extras per repeating journey.`;
 }else if(criterion.kind==='stage'){
  count=js.reduce((sum,j)=>sum+j.steps.filter(s=>s.step===criterion.stage).length,0);unit=fr?'occurrences d’étape':'stage occurrences';
  description=fr?`${stage(criterion.stage)} : ${count} occurrences dans ${matched.length}/${js.length} parcours conservés.`:`${stage(criterion.stage)}: ${count} occurrences in ${matched.length}/${js.length} retained journeys.`;
 }else if(criterion.kind==='prefix'||criterion.kind==='exact'){
  description=fr?`${matched.length}/${js.length} parcours ${criterion.kind==='prefix'?'partagent le préfixe enregistré':'correspondent au groupe exact'} : ${path}.`:`${matched.length}/${js.length} journeys ${criterion.kind==='prefix'?'share the captured prefix':'match the exact group'}: ${path}.`;
 }else description=fr?`${js.length} parcours conservés dans cette sélection. Les limites et résultats explicitement enregistrés restent distincts des observations de navigation.`:`${js.length} retained journeys in this selection. Captured boundaries and explicit outcomes remain separate from navigation observations.`;
 const available=criterion.kind==='all'?js.length>0:matched.length>0;
 if(!available)description=fr?'Observation sélectionnée indisponible dans les détails actuellement conservés. Aucun exemple antérieur n’est rétabli.':'Selected observation is unavailable in currently retained detail. Earlier examples are not restored.';
 const safeID=(id:string)=>/^[a-z0-9_.:-]{1,160}$/i.test(id)&&!/(token|password|secret|csrf|cookie)/i.test(id);
 const journeys=matched.map(j=>j.id).filter(safeID).slice(0,20),traces=[...new Set(matched.flatMap(j=>j.evidence??[]).map(e=>e.trace).filter(id=>/^[a-f0-9]{32}$/i.test(id)))].slice(0,20);
 return {kind:criterion.kind,metric,available,description,unit,count,denominator,supportJourneys:matched.length,lowSupport:matched.length<20,partial,...extras,references:{state:available?'available':'unavailable',total:matched.length,journeys,traces,omitted:Math.max(0,matched.length-journeys.length)}};
}
export function safeInvestigationSummary(context:{app:string;population:Population;bounds:unknown;scope:string;version:string;definition:string;coverage:unknown;sourceCoverage?:unknown;evidenceCoverage?:unknown;freshness?:unknown},criterion:Criterion,observation:InvestigationObservation){
 const cov=context.coverage&&typeof context.coverage==='object'?context.coverage as Record<string,unknown>:null,bounds=context.bounds&&typeof context.bounds==='object'?context.bounds as Record<string,unknown>:null;
 const date=(v:unknown)=>typeof v==='string'&&Number.isFinite(Date.parse(v))?new Date(v).toISOString():null;
 const source=context.sourceCoverage&&typeof context.sourceCoverage==='object'?context.sourceCoverage as Record<string,unknown>:null,evidence=context.evidenceCoverage&&typeof context.evidenceCoverage==='object'?context.evidenceCoverage as Record<string,unknown>:null,freshness=context.freshness&&typeof context.freshness==='object'?context.freshness as Record<string,unknown>:null;
 const readState=freshness&&['current','stale','unavailable','paused'].includes(String(freshness.state))?String(freshness.state):'unavailable';
 const independentCoverage={source:source?{recordLimit:safeCount(source.recordLimit),truncated:source.truncated===true,exporterCompleteness:'unknown'}:{state:'unavailable',exporterCompleteness:'unknown'},receivedRecordList:evidence?{population:'current-query',returned:safeCount(evidence.returned),limit:safeCount(evidence.limit),omitted:safeCount(evidence.omitted)}:{state:'unavailable'},freshness:{state:readState,lastSuccessfulDisplay:date(freshness?.lastSuccessfulDisplay)}};
 const coverage=cov?{journeys:safeCount(cov.journeys),sourcePartial:cov.sourcePartial===true,partialJourneys:safeCount(cov.partialJourneys),censoredJourneys:safeCount(cov.censoredJourneys),graphNodesOmitted:safeCount(cov.graphNodesOmitted),graphEdgesOmitted:safeCount(cov.graphEdgesOmitted),exporterCompleteness:'unknown'}:{state:'unavailable',exporterCompleteness:'unknown'};
 const c=criterion.kind==='all'?{kind:criterion.kind}:criterion.kind==='stage'||criterion.kind==='repeat'?{kind:criterion.kind,stage:sanitizeSummaryText(criterion.stage)}:criterion.kind==='transition'||criterion.kind==='return'?{kind:criterion.kind,from:sanitizeSummaryText(criterion.from),to:sanitizeSummaryText(criterion.to)}:{kind:criterion.kind,path:('path' in criterion?criterion.path:[]).slice(0,32).map(x=>sanitizeSummaryText(x)),grouped:'grouped' in criterion&&criterion.grouped===true};
 const obs={kind:observation.kind,metric:observation.metric,available:observation.available===true,description:sanitizeSummaryText(observation.description,2000),unit:sanitizeSummaryText(observation.unit),count:safeCount(observation.count),denominator:safeCount(observation.denominator),supportJourneys:safeCount(observation.supportJourneys),lowSupport:observation.lowSupport===true,partial:observation.partial===true,...(observation.reachingJourneys!==undefined?{reachingJourneys:safeCount(observation.reachingJourneys),extraOccurrences:safeCount(observation.extraOccurrences),totalOccurrences:safeCount(observation.totalOccurrences),intensity:safeCount(observation.intensity)}:{})};
 const ids=(values:string[])=>values.filter(v=>/^[a-z0-9_.:-]{1,160}$/i.test(v)&&!/(token|password|secret|csrf|cookie)/i.test(v)).slice(0,20);
 const references=observation.available?{state:'available',total:safeCount(observation.references.total),journeys:ids(observation.references.journeys),traces:observation.references.traces.filter(v=>/^[a-f0-9]{32}$/i.test(v)).slice(0,20),omitted:safeCount(observation.references.omitted)}:{state:'unavailable',journeys:[],traces:[]};
 return JSON.stringify({format:'sentrame-investigation-v1',app:sanitizeSummaryText(context.app),population:['current','live','baseline'].includes(context.population)?context.population:'current',bounds:{from:date(bounds?.from),to:date(bounds?.to)},scope:sanitizeSummaryText(context.scope),version:sanitizeSummaryText(context.version),definition:sanitizeSummaryText(context.definition),criterion:c,observation:obs,coverage,independentCoverage,references,limits:{rawEvidence:'omitted',exporterCompleteness:'unknown',referenceValidity:'current authorized retained detail; recheck on opening'}},null,2);
}
