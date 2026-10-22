<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Add Attachment">
      <UButton @click="show()" color="secondary" label="Add Attachment" />
      <template #body>
         <div class="column gap-4">
            <UFileUpload label="Drop your attachment here" v-model="uploadFile" 
               layout="list" :fileImage="false" position="inside"
            />
            <UFormField label="Brief Description" >
               <UTextarea v-model="description" class="w-full"/>
            </UFormField>
         </div>
      </template>
      <template #footer={close}>
         <UButton @click="close" label="Cancel" color="secondary"/>
         <UButton @click="startUpload" label="Attach File" />
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import axios from 'axios'

const unit = useUnitsStore()
const system = useSystemStore()

const isOpen = ref(false)
const uploadFile = ref()
const description = ref("")

const startUpload = ( () => {
   let formData = new FormData()
   formData.append("description", description.value)
   formData.append("name", uploadFile.value.name)
   formData.append("file", uploadFile.value)
   const url=`${system.jobsURL}/units/${unit.detail.id}/attach`
   axios.post(url, formData, {
      headers: {
         'Content-Type': 'multipart/form-data',
      }
   }).then(async () => {
      let uID = unit.detail.id
      unit.detail.id = 0
      await unit.getDetails(uID)
      isOpen.value = false
   }).catch((error) => {
      system.setError(error)
   })
})

const show = (() =>{
   description.value = ""
   uploadFile.value = null
   isOpen.value = true
})
</script>

<style lang="scss" scoped>
</style>
