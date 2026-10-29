export function acceptsSnapshot(request: number, current: number): boolean { return request === current; }
// Area, not radius, grows linearly with the selected duration. Unknown time is a neutral outlined marker.
export function nodeRadius(milliseconds: number | null | undefined): number {
  if (milliseconds == null || !Number.isFinite(milliseconds) || milliseconds < 0) return 22;
  return Math.sqrt(22 * 22 + Math.min(milliseconds, 1800000) / 1800000 * (56 * 56 - 22 * 22));
}
export function duration(milliseconds: number | null | undefined, unavailable = 'Unavailable'): string {
  if (milliseconds == null || !Number.isFinite(milliseconds)) return unavailable;
  if (milliseconds < 1000) return `${Math.round(milliseconds)} ms`;
  if (milliseconds < 60000) return `${(milliseconds / 1000).toFixed(1)} s`;
  return `${(milliseconds / 60000).toFixed(1)} min`;
}
export function dateTime(nanos: number | string | undefined): string {
  if (!nanos) return '—';
  const date = typeof nanos==='string'&&!/^\d+$/.test(nanos) ? new Date(nanos) : new Date(timestampMilliseconds(nanos));
  return Number.isFinite(date.getTime()) ? date.toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'medium', timeZone: 'UTC' })+' UTC' : String(nanos);
}
export function isoLocal(value: string): string { return value ? new Date(value + (value.endsWith('Z') ? '' : 'Z')).toISOString() : ''; }
export function histogramQuantile(histogram: Record<string,number>, quantile:number): string | null {
  const entries = Object.entries(histogram).map(([bucket,count])=>[Number(bucket),count] as const).filter(([bucket,count])=>Number.isFinite(bucket)&&Number.isFinite(count)&&count>0).sort((a,b)=>a[0]-b[0]);
  const total=entries.reduce((n,entry)=>n+entry[1],0);
  if(!total||quantile<0||quantile>1)return null;
  const rank=Math.max(1,Math.ceil(total*quantile));let n=0;
  for(const [bucket,count] of entries){n+=count;if(n<rank)continue;if(bucket===-1)return '<1 ms';if(bucket===999)return '>24 h';const lo=Math.pow(1.1,bucket),hi=Math.pow(1.1,bucket+1);if(hi<1000)return `${lo.toFixed(1)}–${hi.toFixed(1)} ms`;if(hi<60000)return `${(lo/1000).toFixed(1)}–${(hi/1000).toFixed(1)} s`;return `${(lo/60000).toFixed(1)}–${(hi/60000).toFixed(1)} min`;}
  return null;
}
export type SummaryGroupView = {key:string;days:number;records:number;journeys:number;completed:number;inactive:number;exits:number;histograms:Record<string,Record<string,number>>;excluded:Record<string,number>;exclusionsAvailable:boolean};
export function mergeSummaryGroups(days:{groups?:Record<string,any>}[]):SummaryGroupView[]{
  const groups=new Map<string,SummaryGroupView>();
  for(const day of days){for(const [key,value] of Object.entries(day.groups??{})){
    const g=groups.get(key)??{key,days:0,records:0,journeys:0,completed:0,inactive:0,exits:0,histograms:{},excluded:{},exclusionsAvailable:true};
    g.days++;g.records+=value.records??0;g.journeys+=value.journeys??0;g.completed+=value.completed??0;g.inactive+=value.inactive??0;g.exits+=value.exits??0;
    for(const [stage,hist] of Object.entries(value.histograms??{})){const h=g.histograms[stage]??{};for(const [bucket,count] of Object.entries(hist as Record<string,number>)){h[bucket]=(h[bucket]??0)+count}g.histograms[stage]=h}
    if(value.excluded&&typeof value.excluded==='object'){for(const [stage,count] of Object.entries(value.excluded as Record<string,number>))g.excluded[stage]=(g.excluded[stage]??0)+count}else g.exclusionsAvailable=false;
    groups.set(key,g);
  }}return [...groups.values()];
}
export type Discovery = {name:string;count:number;sources:string[];assigned:number};
export function discoveryChoices(discovered:Discovery[]|undefined,evidence:{name:string;source:string;scope:string;tab:string}[]):[string,{count:number;sources:Set<string>;assigned:number}][]{
  if(discovered!==undefined)return discovered.map(d=>[d.name,{count:d.count,sources:new Set(d.sources??[]),assigned:d.assigned}] as [string,{count:number;sources:Set<string>;assigned:number}]).sort((a,b)=>a[0].localeCompare(b[0]));
  const names=new Map<string,{count:number;sources:Set<string>;assigned:number}>();for(const e of evidence){const v=names.get(e.name)??{count:0,sources:new Set<string>(),assigned:0};v.count++;v.sources.add(e.source);if(e.scope&&(e.source!=='browser'||e.tab))v.assigned++;names.set(e.name,v)}return [...names.entries()].sort((a,b)=>a[0].localeCompare(b[0]));
}
export function timestampMilliseconds(value:number|string):number {if(typeof value==='string'&&/^\d{16,}$/.test(value))return Number(BigInt(value)/1000000n);const n=Number(value);return n>1e15?n/1e6:n;}
export function evidenceCopy(value:object):string {return JSON.stringify(value,null,2);}
export function graphGeometry(containerWidth:number,nodeCount:number){const width=Math.max(150,Math.floor(containerWidth));const columns=Math.max(1,Math.floor(width/150));const rows=Math.max(1,Math.ceil(nodeCount/columns));return {width,columns,cellWidth:width/columns,rows,height:rows*195+50};}
export function receiptLagMilliseconds(receipt:number|string,occurrence:number|string):number|null{if(typeof receipt==='string'&&typeof occurrence==='string'&&/^\d+$/.test(receipt)&&/^\d+$/.test(occurrence)){const difference=BigInt(receipt)-BigInt(occurrence);if(difference<0n)return null;const milliseconds=Number(difference)/1e6;return Number.isFinite(milliseconds)?milliseconds:null;}const milliseconds=(Number(receipt)-Number(occurrence))/1e6;return Number.isFinite(milliseconds)&&milliseconds>=0?milliseconds:null;}
