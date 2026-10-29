export type DragStart={pointer:{x:number;y:number};node:{x:number;y:number};scale:number};
// Threshold is in screen pixels; movement is in graph units. Once dragging,
// returning to the starting point must still suppress the subsequent click.
export function dragPosition(start:DragStart,pointer:{x:number;y:number},dragging=false){
 const dx=pointer.x-start.pointer.x,dy=pointer.y-start.pointer.y;
 if(!dragging&&Math.hypot(dx,dy)<5)return null;
 return {x:start.node.x+dx/start.scale,y:start.node.y+dy/start.scale};
}
