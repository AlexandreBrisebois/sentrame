export type RegisteredApp={id:string;generation:string;name:string;state:string};
export function appURL(id:string,target?:{trace?:string;reference?:string;generation?:string}) {
 const q=new URLSearchParams({app:id});
 if(target?.generation)q.set('generation',target.generation);
 if(target?.trace)q.set('trace',target.trace);
 if(target?.reference)q.set('reference',target.reference);
 return '/?'+q;
}
export function initialMode(search:string) {const q=new URLSearchParams(search);return q.has('trace')||q.has('reference')?'trace':'journey'}
// Each mounted workspace owns one request lifetime. Leaving it aborts all reads
// and mutations; even a transport ignoring abort cannot publish a late response.
export class WorkspaceLifetime {
 controller=new AbortController();epoch=0;
 close(){this.epoch++;this.controller.abort()}
 invalidate(){this.close();this.controller=new AbortController()}
 async run<T>(request:(signal:AbortSignal)=>Promise<T>):Promise<T>{const epoch=this.epoch;const result=await request(this.controller.signal);if(epoch!==this.epoch||this.controller.signal.aborted)throw new DOMException('Workspace closed','AbortError');return result}
}
