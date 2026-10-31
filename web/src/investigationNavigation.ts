import type {Journey,Evidence} from './types';
import {contextKey,matchesCriterion,type InvestigationContext,type ReturnFrame} from './investigation.ts';

// In-memory presentation state only. Never save raw examples or this state in references.
export type StoryState={expanded:string[]};
export type ExampleState={search:string;status:string;from:string;page:number;chosen:string;returnFocus:string;invoker:string;listScroll:number;scroll:number;focus:string;story:StoryState};
export type InvestigationState={search:string;focusCell:string;maxDepth:number;repeatSort:string;prefixOrder:string[];examples:ExampleState;story:StoryState};
export const defaultStoryState=():StoryState=>({expanded:[]});
export const defaultExampleState=():ExampleState=>({search:'',status:'',from:'',page:0,chosen:'',returnFocus:'',invoker:'',listScroll:0,scroll:0,focus:'',story:defaultStoryState()});
export const defaultInvestigationState=():InvestigationState=>({search:'',focusCell:'0:0',maxDepth:3,repeatSort:'journeys',prefixOrder:[],examples:defaultExampleState(),story:defaultStoryState()});
export function returnAllowed(frame:ReturnFrame,context:InvestigationContext,journeys:Journey[],unassigned:Evidence[]=[]){
 if(frame.context!==contextKey(context))return false;
 if(!frame.journey)return !frame.criterion&&!!frame.trace&&unassigned.some(e=>e.trace===frame.trace);
 const journey=journeys.find(j=>j.id===frame.journey);
 return !!journey&&(!frame.criterion||matchesCriterion(journey,frame.criterion))&&(!frame.trace||journey.evidence?.some(e=>e.trace===frame.trace));
}
export function reconcileOrder(previous:string[],ranked:string[]){
 const available=new Set(ranked),known=new Set(previous);
 return [...previous.filter(key=>available.has(key)),...ranked.filter(key=>!known.has(key))];
}
