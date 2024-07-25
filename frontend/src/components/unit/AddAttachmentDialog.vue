<template>
   <DPGButton @click="show" severity="secondary" label="Add Attachment"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Attachment" :style="{width: '400px'}" :closable="false">
      <FormKit type="form" id="add-attachment" :actions="false" @submit="addAttachment">
         <FormKit label="" type="file" v-model="info.attachment" required outer-class="first" />
         <FormKit label="Brief Description" type="textarea" rows="4" v-model="info.description"/>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" severity="secondary"/>
            <FormKit type="submit" label="Attach File" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'

const unitsStore = useUnitsStore()

const isOpen = ref(false)
const info = ref({
   attachment: null,
   description: ""
})

async function addAttachment() {
   await unitsStore.attachFile(info.value)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   info.value.description = ""
   info.value.attachment = null
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 15px 0 10px 0;
   margin: 0;
   gap: 10px;
}
</style>
