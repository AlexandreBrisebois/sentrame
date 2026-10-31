<script lang="ts">
 import type {Model} from './types';
 import {investigationObservation,observationCriterion,type Criterion,type Metric,type Population,type PopulationCoverage,type InvestigationObservation} from './investigation';
 import type {Language} from './i18n';
 let {model,criterion,population,coverage,metric='occurrences',language,selectedStep='',selectedEdge='',onexamples,oncopy}:{model:Model;criterion?:Criterion|null;population:Population;coverage?:PopulationCoverage;metric?:Metric;language:Language;selectedStep?:string;selectedEdge?:string;onexamples?:(c:Criterion)=>void;oncopy?:(c:Criterion,p:Population,observation:InvestigationObservation)=>void}=$props();
 let remembered=$state<Criterion|null>(null),selectionKey=$state('');
 $effect(()=>{const key=JSON.stringify([selectedStep,selectedEdge]);if(key!==selectionKey){selectionKey=key;remembered=null;}if(criterion)remembered=criterion});
 const selected=$derived(observationCriterion(criterion??remembered,selectedStep,selectedEdge));
 const observation=$derived(investigationObservation(model,selected,metric,coverage,language));
 const fr=$derived(language==='fr');
</script>
<section class="observation-summary" aria-label={fr?'Observation étayée':'Supported observation'}>
 <h3>{fr?'Observation étayée':'Supported observation'}</h3>
 <p data-observation="description">{observation.description}</p>
 <p data-observation="support">{fr?'Support':'Support'}: {observation.supportJourneys} {fr?'parcours conservés':'retained journeys'}. {#if observation.lowSupport}<strong>{fr?'Support faible : moins de 20 parcours ; observation descriptive uniquement.':'Low support: fewer than 20 journeys; descriptive observation only.'}</strong>{/if}</p>
 {#if observation.partial}<p class="warning">{fr?'Preuves partielles : les occurrences et dénominateurs décrivent uniquement les détails conservés disponibles.':'Partial evidence: occurrences and denominators describe available retained detail only.'}</p>{/if}
 <p>{fr?'La répétition ou le retour ne prouve ni défaut, ni abandon, ni cause. La navigation sans objectif ne crée pas de taux de conversion.':'Repetition or return does not establish a defect, abandonment or cause. Goal-free navigation does not create a conversion rate.'}</p>
 {#if onexamples&&observation.available}<button class="text-button" onclick={()=>onexamples?.(selected)}>{fr?'Examiner les exemples étayés':'Inspect supporting examples'}</button>{/if}
 {#if oncopy}<button class="secondary" onclick={()=>oncopy?.(selected,population,observation)}>{fr?'Copier le résumé de l’investigation':'Copy investigation summary'}</button>{/if}
</section>
<style>.observation-summary{margin-block:12px;padding:12px;border:1px solid var(--line);border-radius:8px;overflow-wrap:anywhere}.observation-summary h3{margin:0 0 6px}.observation-summary p{margin:6px 0}.observation-summary button{margin:4px 8px 0 0}</style>
