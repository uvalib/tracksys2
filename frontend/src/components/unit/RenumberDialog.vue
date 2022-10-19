<template>
   <DPGButton @click="show" class="p-button-secondary renumber" label="Renumber Selected" :disabled="props.disabled"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Renumber Master Files" :style="{width: '400px'}">
      <FormKit type="form" id="renumber" :actions="false" @submit="renumberPages">
         <FormKit label="New starting page number" type="text" v-model="startPage" required autofocus/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
            <FormKit type="submit" label="Renumber" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'

const props = defineProps({
   disabled: {
      type: Boolean,
      default: false
   },
   filenames: {
      type: Array,
      required: true
   }
})

const unitsStore = useUnitsStore()

const isOpen = ref(false)
const startPage = ref()

async function renumberPages() {
   await unitsStore.renumberPages(startPage.value, props.filenames)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   startPage.value = ""
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
button.p-button-secondary.renumber {
   margin-left: 5px;
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
