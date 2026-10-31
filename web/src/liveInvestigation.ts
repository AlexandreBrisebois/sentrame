// Notifications are wake-ups, not record deltas. One read and one pending wake-up.
export class CoalescedRefresh {
 private running=false;private pending=false;private timer:ReturnType<typeof setTimeout>|undefined;private closed=false;
 private read:()=>Promise<void>;private interval:number;
 constructor(read:()=>Promise<void>,interval=500){this.read=read;this.interval=interval}
 request(){if(this.closed)return;this.pending=true;if(!this.running&&!this.timer)this.timer=setTimeout(()=>{this.timer=undefined;void this.flush()},this.interval)}
 private async flush(){if(this.closed||!this.pending)return;this.pending=false;this.running=true;try{await this.read()}finally{this.running=false;if(this.pending)this.request()}}
 invalidate(){this.pending=false;if(this.timer)clearTimeout(this.timer);this.timer=undefined}
 close(){this.invalidate();this.closed=true}
}
export function snapshotFreshness(last:number,now:number,bound=5000){return !last?'No successful snapshot':now-last>bound?'Stale snapshot':'Current snapshot'}
export function resetNeedsWarning(hadSnapshot:boolean,reconnected:boolean,unexpectedContext:boolean){return hadSnapshot&&(reconnected||unexpectedContext)}

export type RefreshTarget='both'|'current'|'baseline';
// Manual actions and scheduled wake-ups share one read lane. Pending requests
// merge, with a visible reconciliation taking priority over validation-only work.
export class SerializedRefresh {
 private active=false;
 private pending:{validationOnly:boolean;target:RefreshTarget;promise:Promise<void>;resolve:()=>void;reject:(reason:unknown)=>void}|null=null;
 private read:(validationOnly:boolean,target:RefreshTarget)=>Promise<void>;
 constructor(read:(validationOnly:boolean,target:RefreshTarget)=>Promise<void>){this.read=read}
 async request(validationOnly=false,target:RefreshTarget='both'):Promise<void>{
  if(this.active){
   if(!this.pending){let resolve=()=>{},reject=(_:unknown)=>{};const promise=new Promise<void>((r,j)=>{resolve=r;reject=j});this.pending={validationOnly,target,promise,resolve,reject}}
   else{this.pending.validationOnly=this.pending.validationOnly&&validationOnly;if(this.pending.target!==target)this.pending.target='both'}
   return this.pending.promise;
  }
  this.active=true;
  try{await this.read(validationOnly,target)}finally{
   this.active=false;const next=this.pending;this.pending=null;
   if(next)void this.request(next.validationOnly,next.target).then(next.resolve,next.reject);
  }
 }
}
