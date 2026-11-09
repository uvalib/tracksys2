<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Customer Email">
      <UButton @click="isOpen=true" label="View Customer Email" color="secondary"/>
      <template  #body>
         <div v-html="props.message" class="email"></div>   
      </template>
      <template #footer="{close}">
         <UButton @click="close" label="Close" color="secondary"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'

const isOpen = ref(false)

const props = defineProps({
   message: {
      type: String,
      required: true
   }
})
</script>

<style lang="scss" scoped>
:deep(a) {
   color: var(--uvalib-blue-alt);
   border-bottom: 1px dashed var(--uvalib-blue-alt);
   &:hover {
      border-bottom: 1px solid var(--uvalib-brand-blue-lighter);
      color: var(--uvalib-brand-blue-lighter);
   }
   &:focus, &:focus-visible {
      outline-offset: 3px;
      outline-color: var( --uvalib-brand-blue-lightest );
   }
}
:deep(p) {
   margin: 10px 0;
}
:deep(h2) {
   font-size: 1.2em;
   font-weight: bold;
}
:deep(h3) {
   margin: 10px 0;
   font-weight: bold;
}
div.email {
   padding: 10px;
   font-size: 0.85em;
}</style>