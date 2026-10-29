export type TimelineData={from:string;to:string;bins:{from:string;to:string;count:number}[];total:number;truncated:boolean;unit:string};
export function clampBrush(start:number,end:number,length:number){
 const size=Math.max(0,Math.floor(length));if(!size)return {start:0,end:0};
 const first=Math.max(0,Math.min(size-1,Math.round(Number.isFinite(start)?start:0)));
 const last=Math.max(first+1,Math.min(size,Math.round(Number.isFinite(end)?end:size)));
 return {start:first,end:last};
}
// Datetime-local fields in this workspace explicitly represent UTC, independent
// of the browser's locale or this Mac's timezone.
export function timelineInstant(value:string){
 const raw=/^\d{4}-\d\d-\d\dT\d\d:\d\d(?::\d\d(?:\.\d+)?)?$/.test(value)?value+'Z':value;
 return Date.parse(raw);
}
export function rangeFromBrush(data:TimelineData,start:number,end:number){
 const brush=clampBrush(start,end,data.bins.length);
 return {from:new Date(timelineInstant(data.bins[brush.start]?.from??data.from)).toISOString(),to:new Date(timelineInstant(data.bins[brush.end-1]?.to??data.to)).toISOString()};
}
export function brushFromBounds(data:TimelineData,from:string,to:string){
 const low=timelineInstant(data.from),high=timelineInstant(data.to),parsedFrom=timelineInstant(from),parsedTo=timelineInstant(to);
 const begin=Number.isFinite(parsedFrom)?parsedFrom:low,finish=Number.isFinite(parsedTo)?parsedTo:high;
 const count=data.bins.length;
 if(!count)return {start:0,end:0,clamped:begin<low||finish>high,outside:finish<=low||begin>=high};
 let start=data.bins.findIndex(bin=>timelineInstant(bin.to)>begin);if(start<0)start=count-1;
 let end=data.bins.findIndex(bin=>timelineInstant(bin.to)>=finish);end=end<0?count:end+1;
 return {...clampBrush(start,end,count),clamped:begin<low||finish>high||finish<=begin,outside:finish<=low||begin>=high};
}
export function recordsInBrush(data:TimelineData,start:number,end:number){
 const brush=clampBrush(start,end,data.bins.length);
 return data.bins.slice(brush.start,brush.end).reduce((sum,bin)=>sum+(Number.isFinite(bin.count)?Math.max(0,bin.count):0),0);
}
