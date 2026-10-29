<script lang="ts">
  let {fieldId,label,help,helpLabel}:{fieldId:string;label:string;help:string[];helpLabel:string}=$props();
  let open=$state(false);
</script>
<div class="field-heading">
  <label for={fieldId}>{label}</label>
  <button type="button" class="field-info" aria-label={helpLabel} aria-expanded={open} aria-controls={`${fieldId}-help`} onclick={()=>open=!open}><span aria-hidden="true">i</span></button>
</div>
<div class="field-explanation" id={`${fieldId}-help`} hidden={!open}>
  {#each help as paragraph}<p>{paragraph}</p>{/each}
</div>
<style>
.field-heading{display:flex;align-items:center;justify-content:space-between;gap:8px;min-height:42px}
.field-info{min-height:42px;min-width:42px;padding:4px;border:0;background:transparent;color:var(--muted,#626e66);border-radius:50%}
.field-info span{display:inline-flex;align-items:center;justify-content:center;border:1px solid currentColor;border-radius:50%;width:16px;height:16px;font-family:Georgia,serif;font-size:11px;font-weight:600;line-height:1}
.field-info:hover,.field-info[aria-expanded=true]{background:#e7ece3;color:var(--green,#153f39)}
.field-explanation{background:#edf0e8;border:1px solid var(--line,#d9dfd5);border-radius:8px;padding:11px 13px;margin:3px 0 9px;font-size:12px;font-weight:400;line-height:1.55;letter-spacing:normal;color:var(--ink,#253b36);overflow-wrap:anywhere}
.field-explanation[hidden]{display:none}.field-explanation p{margin:0 0 8px}.field-explanation p:last-child{margin-bottom:0}
</style>
