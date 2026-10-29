export type Camera={x:number;y:number;k:number};
// Keep the same world coordinate under the gesture centroid, including at bounds.
export function zoomAt(camera:Camera,factor:number,point:{x:number;y:number}):Camera{
 if(!Number.isFinite(factor)||factor<=0)return camera;
 const k=Math.max(.08,Math.min(3,camera.k*factor)),ratio=k/camera.k;
 return {x:point.x-(point.x-camera.x)*ratio,y:point.y-(point.y-camera.y)*ratio,k};
}
