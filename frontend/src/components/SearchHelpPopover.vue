<template>
<DPGButton severity="secondary" size="small" type="button" icon="pi pi-question-circle" label="Search help" @click="toggle" />
<Popover ref="po">
   <div class="content">
      <ul>
         <li>Select 'All Items' to search all types of items TrackSys manages.</li>
         <li>Select an item type to constrain the search to a single type.</li>
         <li>All search terms are stemmed, meaning that a search for <span class="example">prove</span> will match proved, proven and proving.</li>
         <li>If multiple search terms are specifed, records that contain ALL terms in any order (not necessarily adjacent) will be returned.</li>
         <li>Separate multiple search terms with <span class="example">|</span> to return records that match any of the terms.</li>
         <li>Enclose search terms in quotes to find exact matches.<br/>Example: <span class="example">"julian bond papers"</span></li>
         <li>Add <span class="example">*</span> for a wildcard search.<br/>Example: <span class="example">scatter*</span> will match Scattergood and scattered.</li>
         <li>Add <span class="example">?</span> to match any single character. <br/>Example: <span class="example">f?st</span> will match FAST, but not FEAST.</li>
         <li>Add <span class="example">%</span> to match zero or one of any character. <br/>Example: <span class="example">f%ast</span> will match FAST and FEAST.</li>
         <li>
            Call numbers will be matched regardless of spacing.<br>Example: <span class="example">PS3029 .T3 C5 1835</span> and <span class="example">PS3029.T3C51835</span>
            will match PS3029 .T3 C5 1835.
         </li>
      </ul>
      <div>
         Full documentation on searching can be found
         <a href="https://manual.manticoresearch.com/Searching/Full_text_matching/Operators?client=CONFIG#Full-text-operators" target="_blank">here</a>.</div>
   </div>
</Popover>
</template>

<script setup>
import { ref } from 'vue'
import Popover from 'primevue/popover'

const po = ref()

const toggle = ((event) => {
   po.value.toggle(event);
})
</script>

<style lang="scss" scoped>
.content {
   max-width: 400px;
   display: flex;
   flex-direction: column;
   gap: 10px;
   li {
      padding: 5px 0;
   }
   .example {
      font-family: monospace;
      font-size: 1.15em;
   }
}
</style>
