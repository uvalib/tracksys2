<template>
<UPopover v-model:open="open">
   <UButton variant="link" class="text-black" :trailing-icon="icon" :label="props.label" />
   <template #content>
      <div class="p-4 column-close">
         <slot></slot>
         <div class="row-right">
            <UButton color="secondary" label="Clear" size="xs" @click="clearClicked"/>
            <UButton label="Apply" size="xs" @click="applyClicked"/>
         </div>
      </div>
   </template>
</UPopover>
</template>

<script setup>
import { ref, computed } from 'vue'

const open = ref(false)

const emit = defineEmits( ['clear', 'apply'])
const props = defineProps({
   label: {
      type: String,
      required: true
   },
   applied: {
      type: Boolean,
      default: false
   }
})

const icon = computed(() => {
   if ( props.applied) return  'i-lucide-funnel-plus'
   return "i-lucide-funnel"
})

const applyClicked = (() => {
   open.value=false
   emit('apply')
})

const clearClicked = (() => {
   open.value=false
   emit('clear')
})
</script>