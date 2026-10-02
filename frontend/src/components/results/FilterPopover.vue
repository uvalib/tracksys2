<template>
<UPopover v-model:open="open">
   <UButton variant="link" class="text-black" trailing-icon="i-lucide-funnel" :label="props.label" />
   <template #content>
      <div class="p-2 column-close">
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
import { ref } from 'vue'

const open = ref(false)

const emit = defineEmits( ['clear', 'apply'])
const props = defineProps({
   label: {
      type: String,
      required: true
   }
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