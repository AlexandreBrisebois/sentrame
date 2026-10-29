// Invitation credentials belong only in the fragment. Return the same local
// path without its fragment before the application makes authenticated reads.
export function invitationLocation(location:{hash:string;pathname:string;search:string}){
 const fragment=new URLSearchParams(location.hash.replace(/^#/,''));
 if(!fragment.has('invite'))return null;
 return {token:fragment.get('invite')||'invalid',cleanURL:location.pathname+location.search};
}
export function invitationIsLoopback(url:string){
 try{return ['127.0.0.1','localhost','[::1]'].includes(new URL(url).hostname.toLowerCase())}catch{return false}
}
