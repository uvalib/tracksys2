<template>
   <DPGButton @click="show" severity="secondary" label="Add Attachment"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Attachment" :style="{width: '400px'}" :closable="false">
      <div class="upload">
         <FileUpload ref="uploader" mode="basic" :customUpload="true" name="file" @uploader="startUpload($event)" />
         <FormField id="attdesc" label="Brief Description">
            <Textarea id="attdesc" v-model="description" rows="4"/>   
         </FormField>
      </div>
      <template #footer>
         <DPGButton @click="hide" label="Cancel" severity="secondary"/>
         <DPGButton @click="uploader.upload()" label="Attach File" />
      </template>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import FileUpload from 'primevue/fileupload'
import Textarea from 'primevue/textarea'
import FormField from '@/components/FormField.vue'
import axios from 'axios'

const unit = useUnitsStore()
const system = useSystemStore()

const isOpen = ref(false)
const uploader = ref()
const description = ref("")

const startUpload = ( (event) => {
   const fileData = event.files[0]
   let formData = new FormData()
   formData.append("description", description.value)
   formData.append("name", fileData.name)
   formData.append("file", fileData)
   const url=`${system.jobsURL}/units/${unit.detail.id}/attach`
   axios.post(url, formData, {
      headers: {
         'Content-Type': 'multipart/form-data',
      }
   }).then(async () => {
      let uID = unit.detail.id
      unit.detail.id = 0
      await unit.getDetails(uID)
      hide()
   }).catch((error) => {
      system.setError(error)
   })
})

const hide = (() => {
   isOpen.value=false
})
const show = (() =>{
   description.value = ""
   isOpen.value = true
})
</script>

<style lang="scss" scoped>
.upload {
   display: flex;
   flex-direction: column;
   gap: 15px;
}
</style>
