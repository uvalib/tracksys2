<template>
   <DPGButton @click="show" class="p-button-secondary metadata" label="Assign Component" :disabled="props.disabled"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Assign Component" :style="{width: '400px'}">
      <FormKit type="form" id="assign-metadata" :actions="false" @submit="assignComponent">
         <div class="lookup">
            <FormKit label="Component ID" type="text" v-model="componentID" required autofocus outer-class="first" />
            <LookupDialog target="component" @selected="componentSelected" />
         </div>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
            <FormKit type="submit" label="Assign" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'
import LookupDialog from './LookupDialog.vue';

const props = defineProps({
   disabled: {
      type: Boolean,
      default: false
   },
   ids: {
      type: Array,
      required: true
   }
})

const unitsStore = useUnitsStore()

const isOpen = ref(false)
const componentID = ref()

function componentSelected( md ) {
   componentID.value = md
}
async function assignComponent() {
   await unitsStore.assignComponent(componentID.value, props.ids)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   componentID.value = ""
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
button.p-button-secondary.metadata {
   margin-left: 5px;
}
.lookup {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: end;
   padding-bottom: 5px;
   :deep(.formkit-outer) {
      flex-grow: 1;
   }
}
.spacer {
   display: inline-block;
   margin: 0 5px;
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 15px 0 10px 0;
   margin: 0;
   button {
      margin-right: 10px;
   }
}
</style>
