<script lang="ts">
 import Disclosure from './Disclosure.svelte';
 import {translate,type Language} from './i18n';
 let {layout,language,live=false}:{layout:string;language:Language;live?:boolean}=$props();
 const t=(s:string)=>translate(s,language);
</script>
<Disclosure label="Help" {language}>
 <p><strong>{t(layout==='matrix'?'Where do journeys go next?':layout==='repeats'?'Which activities repeat?':layout==='sequence'?'Which ordered routes are shared?':layout==='journey-timeline'?'What happened in one journey?':'How do observed stages connect?')}</strong></p>
 <p>{t(layout==='sequence'?'Unit: journeys sharing an ordered prefix':layout==='journey-timeline'?'Unit: ordered stage occurrences':'Unit: observed occurrences, journeys and visits')}</p>
 <p>{t('Select a stage, transition or matching journey to inspect supporting evidence.')}</p>
 {#if live}<p>{t('Live uses L and solid contributions; Baseline uses B and patterned contributions. Error evidence is separate from population.')}</p>{/if}
</Disclosure>
