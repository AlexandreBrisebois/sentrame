<script lang="ts">
  import { onMount } from 'svelte';
  import { zoomAt, type Camera } from './camera';
  import { dragPosition, type DragStart } from './nodeDrag';
  import { errorEvidence } from './graphEvidence';
  import type { Model } from './types';
  import type { Language } from './i18n';
  import { translate } from './i18n';
  import { graphLayout, repositionGraph, sankeyLayer, edgeKey, type NodePositions, lensEvidence, type Layout, type Lens, type GraphNode, type GraphEdge } from './graph';
  let {model,baseline=null,language='en',selected='',onselect=(_:string)=>{},layout='hierarchy',lens='none',discovery=[],selectedEdge='',onedge=(_:string,__:string)=>{},goal=''}:{model:Model;baseline?:Model|null;language?:Language;selected?:string;onselect?:(id:string)=>void;layout?:Layout;lens?:Lens;discovery?:{id:string;label:string;occurrences:number}[];selectedEdge?:string;onedge?:(from:string,to:string)=>void;goal?:string}=$props();
  const t=(s:string)=>translate(s,language);
  const marker=$props.id();
  let frame:HTMLDivElement;
  let viewport=$state({width:900,height:600});
  let camera=$state({x:0,y:0,k:1});
  let userCamera=$state(false);
  let measured=$state(false);
  let initialized=false;
  let pan:{x:number;y:number;cx:number;cy:number}|null=null;
  let placements=$state<Record<Layout,NodePositions>>({hierarchy:{},force:{},sankey:{}});
  let nodeDrag:{id:string;pointerId:number;layout:Layout;start:DragStart;moved:boolean;previous?:{x:number;y:number}}|null=null;
  let suppressClick=false;
  const nodeLimit=100;
  let allNodes=$derived.by(()=>{
    const result=new Map<string,GraphNode>();
    for(const n of baseline?.nodes??[])result.set(n.id,{...n});
    for(const n of model.nodes??[]){const old=result.get(n.id);result.set(n.id,{...n,occurrences:Math.max(n.occurrences,old?.occurrences??0)});}
    for(const n of discovery)if(!result.has(n.id))result.set(n.id,{...n,discovered:true});
    return [...result.values()];
  });
  let allEdges=$derived.by(()=>{const result=new Map<string,GraphEdge>();for(const e of [...(baseline?.edges??[]),...(model.edges??[])]){const key=edgeKey(e.from,e.to),old=result.get(key);result.set(key,{...e,occurrences:Math.max(e.occurrences,old?.occurrences??0)});}return [...result.values()];});
  let automaticGeometry=$derived(graphLayout(allNodes.slice(0,nodeLimit),allEdges,layout));
  let geometry=$derived(repositionGraph(automaticGeometry,placements[layout],layout));
  let hasPlacements=$derived(Object.keys(placements[layout]).length>0);
  let baselineGeometry=$derived(layout==='sankey'?sankeyLayer(geometry.nodes,baseline?.nodes??[],baseline?.edges??[],geometry.scale):geometry);
  let liveGeometry=$derived(layout==='sankey'?sankeyLayer(geometry.nodes,[...model.nodes,...discovery.map(n=>({...n,discovered:true}))],model.edges,geometry.scale,baseline?8:18):geometry);
  let errors=$derived(errorEvidence(model));
  let baselineErrors=$derived(baseline?errorEvidence(baseline):{edges:new Set<string>(),nodes:new Set<string>()});
  let highlights=$derived(lensEvidence(model,lens));
  let endingStatuses=$derived.by(()=>{const result=new Map<string,Map<string,number>>();for(const j of model.journeys){const id=j.steps?.at(-1)?.step??j.path?.at(-1);if(!id)continue;const statuses=result.get(id)??new Map<string,number>();const status=j.status==='active'?'Pending':j.status==='completed'?'Completed':j.status==='inactive'?'Inactive':j.status==='exited'?'Exited':'Unknown';statuses.set(status,(statuses.get(status)??0)+1);result.set(id,statuses);}return result;});
  function endingLabel(id:string){const entries=[...(endingStatuses.get(id)??[])];if(entries.length===1)return `${t(entries[0][0])} · ${entries[0][1]}`;return `${t('Mixed endings')} · ${entries.reduce((n,[_,count])=>n+count,0)}`;}
  function endingDetail(id:string){return [...(endingStatuses.get(id)??[])].map(([status,count])=>`${t(status)}: ${count}`).join(' · ');}
  let neighbors=$derived.by(()=>{const ids=new Set([selected]);for(const e of allEdges)if(e.from===selected||e.to===selected){ids.add(e.from);ids.add(e.to);}return ids;});
  let visibleCamera=$derived(userCamera?camera:fitCamera());
  function fitCamera(){const available=Math.max(150,viewport.height-125),k=Math.max(.08,Math.min(1.2,(viewport.width-90)/geometry.width,available/geometry.height));return {x:(viewport.width-geometry.width*k)/2,y:15+(available-geometry.height*k)/2,k};}
  function fit(){camera=fitCamera();userCamera=true;}
  function zoom(factor:number){camera=zoomAt(visibleCamera,factor,{x:viewport.width/2,y:viewport.height/2});userCamera=true;}
  let gesture:{camera:Camera;point:{x:number;y:number}}|null=null;
  function gesturePoint(e:{clientX?:number;clientY?:number}){const rect=frame.getBoundingClientRect();return {x:e.clientX?e.clientX-rect.left:viewport.width/2,y:e.clientY?e.clientY-rect.top:viewport.height/2};}
  function wheel(e:WheelEvent){
    e.preventDefault();if(gesture)return;
    if(e.ctrlKey){camera=zoomAt(visibleCamera,Math.exp(-Math.max(-200,Math.min(200,e.deltaY))*.012),gesturePoint(e));}
    else{const unit=e.deltaMode===1?16:e.deltaMode===2?viewport.height:1;const c=visibleCamera;camera={...c,x:c.x-e.deltaX*unit,y:c.y-e.deltaY*unit};}
    userCamera=true;
  }
  function gestureStart(e:Event){e.preventDefault();gesture={camera:{...visibleCamera},point:gesturePoint(e as Event&{clientX:number;clientY:number})};}
  function gestureChange(e:Event){e.preventDefault();if(!gesture)return;const scale=(e as Event&{scale:number}).scale;camera=zoomAt(gesture.camera,scale,gesture.point);userCamera=true;}
  function gestureEnd(e:Event){e.preventDefault();gesture=null;}
  function place(id:string,position:{x:number;y:number},view=layout){placements={...placements,[view]:{...placements[view],[id]:position}};}
  function resetLayout(){placements={...placements,[layout]:{}};fit();}
  function pointerDown(e:PointerEvent){
    if(e.button!==0||nodeDrag||pan)return;
    suppressClick=false;
    const target=e.target as Element,nodeElement=target.closest('[data-node]');
    if(nodeElement){
      const id=nodeElement.getAttribute('data-node')!,node=geometry.nodes.find(n=>n.id===id);if(!node)return;
      const c=visibleCamera;camera={...c};userCamera=true;
      nodeDrag={id,pointerId:e.pointerId,layout,start:{pointer:{x:e.clientX,y:e.clientY},node:{x:node.x,y:node.y},scale:c.k},moved:false,previous:placements[layout][id]};
      (nodeElement as SVGElement).setPointerCapture(e.pointerId);return;
    }
    if(target.closest('[data-item]'))return;
    const c=visibleCamera;pan={x:e.clientX,y:e.clientY,cx:c.x,cy:c.y};camera={...c};userCamera=true;(e.currentTarget as SVGSVGElement).setPointerCapture(e.pointerId);
  }
  function pointerMove(e:PointerEvent){
    if(nodeDrag&&nodeDrag.pointerId===e.pointerId){
      const position=dragPosition(nodeDrag.start,{x:e.clientX,y:e.clientY},nodeDrag.moved);
      if(position){nodeDrag.moved=true;suppressClick=true;place(nodeDrag.id,position,nodeDrag.layout);}return;
    }
    if(pan)camera={...camera,x:pan.cx+e.clientX-pan.x,y:pan.cy+e.clientY-pan.y};
  }
  function pointerUp(e:PointerEvent){
    if(nodeDrag&&nodeDrag.pointerId===e.pointerId){
      if(e.type==='pointercancel'){
        const restored={...placements[nodeDrag.layout]};
        if(nodeDrag.previous)restored[nodeDrag.id]=nodeDrag.previous;else delete restored[nodeDrag.id];
        placements={...placements,[nodeDrag.layout]:restored};
      }
      nodeDrag=null;
    }
    pan=null;
  }
  function nodeClick(e:MouseEvent,id:string){if(suppressClick){suppressClick=false;return;}onselect(id);}
  function nodeKey(e:KeyboardEvent,id:string){
    const moves:Record<string,[number,number]>={ArrowLeft:[-20,0],ArrowRight:[20,0],ArrowUp:[0,-20],ArrowDown:[0,20]};
    if(e.shiftKey&&moves[e.key]){
      e.preventDefault();e.stopPropagation();const node=geometry.nodes.find(n=>n.id===id)!;
      const [dx,dy]=moves[e.key];place(id,{x:node.x+dx,y:node.y+dy});return;
    }
    itemKey(e,()=>onselect(id));
  }
  function canvasKey(e:KeyboardEvent){if(e.target!==e.currentTarget)return;const keys:Record<string,[number,number]>={ArrowLeft:[40,0],ArrowRight:[-40,0],ArrowUp:[0,40],ArrowDown:[0,-40]};if(keys[e.key]){e.preventDefault();const c=visibleCamera;camera={...c,x:c.x+keys[e.key][0],y:c.y+keys[e.key][1]};userCamera=true;}else if(e.key==='+'||e.key==='='){e.preventDefault();zoom(1.2);}else if(e.key==='-'){e.preventDefault();zoom(1/1.2);}else if(e.key==='0'){e.preventDefault();fit();}}
  function itemKey(e:KeyboardEvent,action:()=>void){if(e.key==='Enter'||e.key===' '){e.preventDefault();e.stopPropagation();action();}}
  function edgeActive(from:string,to:string){return selectedEdge===edgeKey(from,to)||!!selected&&(from===selected||to===selected)||lens!=='none'&&highlights.edges.has(edgeKey(from,to));}
  function edgeOpacity(from:string,to:string){if(selectedEdge)return selectedEdge===edgeKey(from,to)?1:.18;if(selected)return from===selected||to===selected?1:.16;if(lens!=='none')return highlights.edges.has(edgeKey(from,to))?.95:.18;if(errors.edges.has(edgeKey(from,to)))return .9;return layout==='sankey'?.4:.58;}
  function nodeOpacity(id:string){if(selectedEdge){const edge=geometry.edges.find(e=>edgeKey(e.from,e.to)===selectedEdge);return edge&&(id===edge.from||id===edge.to)?1:.3;}if(selected)return neighbors.has(id)?1:.3;if(lens!=='none')return highlights.nodes.has(id)?1:.35;return 1;}
  function short(label:string){return label.length>23?label.slice(0,22)+'…':label;}
  const counts=(id:string)=>{const n=model.nodes.find(n=>n.id===id),b=baseline?.nodes?.find(n=>n.id===id);return b?`${b.occurrences} ${t('baseline')} · ${n?.occurrences??0} ${t('live')}`:`${n?.occurrences??discovery.find(n=>n.id===id)?.occurrences??0} ${t('occurrences')}`;};
  $effect(()=>{if(measured&&geometry.nodes.length&&!initialized){camera=fitCamera();userCamera=true;initialized=true;}});
  onMount(()=>{const observer=new ResizeObserver(entries=>{const r=entries[0]?.contentRect;if(r){viewport={width:Math.max(300,r.width),height:Math.max(260,r.height)};measured=true;}});observer.observe(frame);frame.addEventListener('wheel',wheel,{passive:false});frame.addEventListener('gesturestart',gestureStart,{passive:false});frame.addEventListener('gesturechange',gestureChange,{passive:false});frame.addEventListener('gestureend',gestureEnd,{passive:false});return()=>{observer.disconnect();frame.removeEventListener('wheel',wheel);frame.removeEventListener('gesturestart',gestureStart);frame.removeEventListener('gesturechange',gestureChange);frame.removeEventListener('gestureend',gestureEnd);};});
</script>
<div class="workspace-graph" class:sankey-view={layout==='sankey'} bind:this={frame}>
  <!-- svelte-ignore a11y_no_noninteractive_tabindex, a11y_no_noninteractive_element_interactions (Interactive SVG application implements keyboard pan/zoom and focusable nodes/edges; equivalent table lives in workspace.) -->
  <svg class="graph-svg" viewBox={`0 0 ${viewport.width} ${viewport.height}`} role="application" aria-label={`${t('Journey map')} · ${layout}`} aria-describedby={`${marker}-help`} tabindex="0" onkeydown={canvasKey} onpointerdown={pointerDown} onpointermove={pointerMove} onpointerup={pointerUp} onpointercancel={pointerUp}>
    <defs><marker id={`${marker}-arrow`} viewBox="0 0 7 7" refX="6.5" refY="3.5" markerWidth={6/visibleCamera.k} markerHeight={6/visibleCamera.k} markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 7 3.5 L 0 7 z" fill="#54796a"/></marker><marker id={`${marker}-baseline-arrow`} viewBox="0 0 7 7" refX="6.5" refY="3.5" markerWidth={6/visibleCamera.k} markerHeight={6/visibleCamera.k} markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 7 3.5 L 0 7 z" fill="#9dad9e"/></marker><marker id={`${marker}-error-arrow`} viewBox="0 0 7 7" refX="6.5" refY="3.5" markerWidth={6/visibleCamera.k} markerHeight={6/visibleCamera.k} markerUnits="userSpaceOnUse" orient="auto"><path d="M 0 0 L 7 3.5 L 0 7 z" fill="#a3493d"/></marker></defs>
    <g transform={`translate(${visibleCamera.x} ${visibleCamera.y}) scale(${visibleCamera.k})`}>
      {#each geometry.edges as edge (edgeKey(edge.from,edge.to))}
        {@const key=edgeKey(edge.from,edge.to)}
        {@const live=model.edges.find(e=>edgeKey(e.from,e.to)===key)}
        {@const base=baseline?.edges?.find(e=>edgeKey(e.from,e.to)===key)}
        {@const active=edgeActive(edge.from,edge.to)}
        {@const baseRoute=baselineGeometry.edges.find(e=>edgeKey(e.from,e.to)===key)}
        {@const liveRoute=liveGeometry.edges.find(e=>edgeKey(e.from,e.to)===key)}
        {#if base&&baseRoute}<path d={baseRoute.path} fill="none" stroke={baselineErrors.edges.has(key)?'#a3493d':'#aebcaf'} stroke-width={layout==='sankey'?baseRoute.width:edge.width+2} opacity={edgeOpacity(edge.from,edge.to)*.7} marker-end={layout!=='sankey'||edge.returning?`url(#${marker}-${baselineErrors.edges.has(key)?'error-arrow':'baseline-arrow'})`:undefined} class="flow-line"/>{/if}
        {#if live&&liveRoute}<path d={liveRoute.path} fill="none" stroke={errors.edges.has(key)?'#a3493d':active?'#285e4d':'#789386'} stroke-width={layout==='sankey'?liveRoute.width:edge.width} opacity={edgeOpacity(edge.from,edge.to)} marker-end={layout!=='sankey'||edge.returning?`url(#${marker}-${errors.edges.has(key)?'error-arrow':'arrow'})`:undefined} class="flow-line"/>{:else if !base}<path d={edge.path} fill="none" stroke={errors.edges.has(key)?'#a3493d':'#789386'} stroke-width={edge.width} opacity={edgeOpacity(edge.from,edge.to)} marker-end={layout!=='sankey'?`url(#${marker}-${errors.edges.has(key)?'error-arrow':'arrow'})`:undefined} class="flow-line"/>{/if}
        <path data-item="edge" class="edge-hit" d={liveRoute?.path??baseRoute?.path??edge.path} fill="none" stroke="transparent" stroke-width={Math.max(24/visibleCamera.k,edge.width)} role="button" tabindex="0" aria-label={`${edge.from} → ${edge.to}; ${live?.occurrences??0} ${t('transitions')}${base?`; ${base.occurrences} ${t('baseline transitions')}`:''}${errors.edges.has(key)||baselineErrors.edges.has(key)?'; '+t('Error evidence on this connection'):''}`} aria-pressed={selectedEdge===key} onclick={()=>onedge(edge.from,edge.to)} onkeydown={e=>itemKey(e,()=>onedge(edge.from,edge.to))}><title>{edge.from} → {edge.to} · {live?.occurrences??0} {t('transitions')}{base?` · ${base.occurrences} ${t('baseline transitions')}`:""}{edge.returning?' · '+t('return / repeat path'):''}{errors.edges.has(key)||baselineErrors.edges.has(key)?' · '+t('Error evidence on this connection'):''}</title></path>
        {#if base&&baseRoute&&liveRoute&&layout==='sankey'}<path data-item="edge" d={baseRoute.path} fill="none" stroke="transparent" stroke-width={Math.max(24/visibleCamera.k,baseRoute.width)} class="edge-hit" onclick={()=>onedge(edge.from,edge.to)} aria-hidden="true"/>{/if}
      {/each}
      {#each geometry.nodes as node (node.id)}
        {@const chosen=selected===node.id}
        {@const highlighted=lens!=='none'&&highlights.nodes.has(node.id)}
        {@const isGoal=goal===node.id}
        {@const liveBar=liveGeometry.nodes.find(n=>n.id===node.id)}
        {@const baseBar=baseline?baselineGeometry.nodes.find(n=>n.id===node.id):undefined}
        <g data-item="node" data-node={node.id} class="graph-node" class:chosen class:highlighted class:hierarchy-node={layout==='hierarchy'} role="button" tabindex="0" aria-label={`${node.label}; ${counts(node.id)}${node.discovered?'; '+t('Unmapped event'):''}${isGoal?'; '+t('Operational goal'):''}`} aria-pressed={chosen} onclick={e=>nodeClick(e,node.id)} onkeydown={e=>nodeKey(e,node.id)} opacity={nodeOpacity(node.id)}>
          <title>{node.label} · {counts(node.id)}{node.discovered?' · '+t('Unmapped event — no journey relationship inferred'):''}{lens==='endings'?' · '+endingDetail(node.id):''}</title>
          {#if layout==='force'}
            <circle class="node-shape" cx={node.x} cy={node.y} r="23" fill={chosen?'#153f39':'#fffefa'} stroke={chosen?'#153f39':'#789181'} stroke-width={chosen?2.5:1.4} stroke-dasharray={node.discovered?'4 4':undefined}/>
            {#if isGoal}<text class="goal-mark" x={node.x} y={node.y+5} text-anchor="middle" fill={chosen?'#fffefa':'#376b54'}>✓</text>{/if}
            <text class="node-label" x={node.x} y={node.y+46} text-anchor="middle">{short(node.label)}</text>
            <text class="node-count" x={node.x} y={node.y+64} text-anchor="middle">{node.discovered?t('Unmapped')+' · ':''}{counts(node.id)}</text>
          {:else if layout==='sankey'}
            {#if baseBar}<rect class="node-shape sankey-bar" x={baseBar.x-baseBar.width/2} y={baseBar.y-baseBar.height/2} width={baseBar.width} height={baseBar.height} rx="1" fill="#b5c2b5" stroke={chosen?'#426956':'#b5c2b5'} stroke-width={chosen?1.5:0}/>{/if}
            {#if liveBar}<rect class="node-shape sankey-bar" x={liveBar.x-liveBar.width/2} y={liveBar.y-liveBar.height/2} width={liveBar.width} height={liveBar.height} rx="1" fill={chosen?'#153f39':'#537c66'} stroke={chosen?'#153f39':'#537c66'} stroke-width={chosen?1.5:0} stroke-dasharray={node.discovered?'4 4':undefined}/>{/if}
            <rect x={node.x-17} y={node.y-Math.max(22,node.height/2)} width="34" height={Math.max(44,node.height)} fill="transparent"/>
            <text class="node-label" x={node.x+18} y={node.y-4}>{isGoal?'✓ ':''}{short(node.label)}</text><text class="node-count" x={node.x+18} y={node.y+15}>{node.discovered?t('Unmapped')+' · ':''}{counts(node.id)}</text>
          {:else}
            <rect class="node-shape" x={node.x-node.width/2} y={node.y-node.height/2} width={node.width} height={node.height} rx="9" fill={chosen?'#153f39':'#fffefa'} stroke={chosen?'#153f39':'#9aaa9a'} stroke-width={chosen?2.5:1.2} stroke-dasharray={node.discovered?'4 4':undefined}/>
            <text class="node-label" x={node.x} y={node.y-3} text-anchor="middle" fill={chosen?'#fffefa':undefined}>{isGoal?'✓ ':''}{short(node.label)}</text>
            <text class="node-count" x={node.x} y={node.y+16} text-anchor="middle" fill={chosen?'#d8e4da':undefined}>{node.discovered?t('Unmapped')+' · ':''}{counts(node.id)}</text>
          {/if}
          {#if errors.nodes.has(node.id)||baselineErrors.nodes.has(node.id)}<text x={node.x-node.width/2-9} y={node.y+4} class="error-mark" aria-label={t('Error evidence attached')}>!</text>{/if}
          {#if highlighted}<circle cx={node.x-node.width/2-9} cy={node.y-node.height/2-7} r="5" fill="#fffefa" stroke="#366650" stroke-width="2"/><text x={node.x-node.width/2-9} y={node.y-node.height/2-17} text-anchor="middle" class="lens-mark">{lens==='exceptions'?'!':lens==='recovery'?'↪':lens==='endings'?endingLabel(node.id):t('Rare')}</text>{/if}
        </g>
      {/each}
    </g>
  </svg>
  {#if !allNodes.length}<div class="graph-empty"><span aria-hidden="true">↔</span><h2>{t('No journeys in this view')}</h2><p>{t('Receive events or adjust the time and visit filters.')}</p></div>{/if}
  <div class="canvas-navigation" aria-label={t('Graph navigation')}><button class="secondary" aria-label={t('Zoom out')} onclick={()=>zoom(1/1.2)}>−</button><span>{Math.round(visibleCamera.k*100)}%</span><button class="secondary" aria-label={t('Zoom in')} onclick={()=>zoom(1.2)}>+</button><button class="secondary" onclick={fit}>{t('Fit')}</button><button class="secondary" disabled={!hasPlacements} onclick={resetLayout}>{t('Reset layout')}</button></div>
  <div class="canvas-legend">{#if errors.edges.size||baselineErrors.edges.size}<span class="error-key">! {t('Red connections: error evidence')}</span>{/if}<span>{t(layout==='sankey'?'Bar height and ribbon width: occurrences':'Arrows: observed direction · line width: traffic')}</span>{#if baseline}<span>{t('Muted: seven-day baseline · Green: live · counts stay separate')}</span>{/if}{#if discovery.length}<span>{t('Dashed outline: unmapped event')}</span>{/if}{#if lens==='rare'}<span>{t('Rare: ≤10% of ≥20 closed, comparable journeys; partial / ongoing paths excluded. Highlighted links may also occur in common paths.')}</span>{/if}{#if allNodes.length>nodeLimit}<strong>{t('Showing')} {nodeLimit}/{allNodes.length} {t('nodes. Narrow filters; omitted relationships are not joined.')}</strong>{/if}</div>
  <p id={`${marker}-help`} class="sr-only">{t('Tab to select stages and connections. Enter opens the inspector. On the canvas use arrow keys to pan, plus and minus to zoom, zero to fit. Drag nodes to rearrange; Shift plus arrow keys also moves nodes. Hierarchy and Sankey keep columns fixed. Drag empty space or two-finger scroll to pan; pinch to zoom. Timing and ordered evidence are in the inspector and table.')}</p>
</div>
<style>
.workspace-graph{height:100%;min-height:300px;width:100%;position:relative;overflow:hidden;background:var(--paper,#f5f4ee)}
.graph-svg{display:block;width:100%;height:100%;touch-action:none;cursor:grab}.graph-svg:active{cursor:grabbing}.graph-svg:focus,.graph-svg:focus-visible{outline:none}.graph-svg{-webkit-tap-highlight-color:transparent}.graph-svg:focus-visible~.canvas-navigation{box-shadow:0 0 0 2px var(--focus,#bd782a)}.flow-line{pointer-events:none;stroke-linecap:round;stroke-linejoin:round}.sankey-view .flow-line{stroke-linecap:butt}.edge-hit{cursor:pointer;outline:none;-webkit-tap-highlight-color:transparent}.edge-hit:focus{outline:none}.edge-hit:focus-visible{outline:none;stroke:#bd782a;stroke-opacity:.35}.graph-node{user-select:none;cursor:grab;outline:none;-webkit-tap-highlight-color:transparent}.graph-node:active{cursor:grabbing}.graph-node:focus-visible .node-shape{stroke:#bd782a;stroke-width:3}.node-label{font-family:Inter,ui-sans-serif,system-ui,sans-serif;font-size:12px;font-weight:600;fill:#253b36;paint-order:stroke;stroke:#f5f4ee;stroke-width:3px;stroke-linejoin:round}.node-count{font-family:Inter,ui-sans-serif,system-ui,sans-serif;font-size:10px;fill:#626e66;paint-order:stroke;stroke:#f5f4ee;stroke-width:2px}.graph-node.chosen.hierarchy-node .node-label{stroke:none;fill:#fffefa}.graph-node.chosen.hierarchy-node .node-count{stroke:none;fill:#d8e4da}.goal-mark{font-size:16px}.error-mark{fill:#a3493d;font:700 14px system-ui;paint-order:stroke;stroke:var(--paper,#f5f4ee);stroke-width:3px}.error-key{color:#a3493d}.lens-mark{fill:#366650;font:600 10px system-ui;paint-order:stroke;stroke:#f5f4ee;stroke-width:3px}.canvas-navigation{position:absolute;right:22px;bottom:18px;display:flex;align-items:center;gap:2px;padding:3px;background:#fffefaeF;border:1px solid var(--line,#d9dfd5);border-radius:8px;box-shadow:0 2px 12px #153f3908}.canvas-navigation button{min-width:38px;min-height:36px;padding:6px 10px;border:0;background:transparent}.canvas-navigation span{font-size:10px;color:var(--muted,#626e66);min-width:40px;text-align:center}.canvas-legend{position:absolute;bottom:15px;left:24px;right:330px;display:flex;gap:5px 20px;flex-wrap:wrap;pointer-events:none;font-size:10px;color:var(--muted,#626e66);max-width:750px}.canvas-legend span,.canvas-legend strong{background:#f5f4eee8;padding:2px 0}.graph-empty{position:absolute;top:35%;width:100%;text-align:center;pointer-events:none;color:var(--muted,#626e66)}.graph-empty span{font-size:38px}.graph-empty h2{font-size:25px}.graph-empty p{font-size:12px;margin-top:15px}@media(max-width:600px){.canvas-navigation{right:12px;left:12px;bottom:12px;flex-wrap:wrap;justify-content:center;gap:4px}.canvas-navigation button{min-height:44px}.canvas-legend{left:14px;right:14px;bottom:78px;font-size:9px;gap:4px 12px}}
</style>
