<template>
   <UModal v-model:open="open" :modal="true" :dismissible="false" :close="false" title="Upload XML">
      <UButton @click="show()" label="Upload XML" />
      <template #body>
         <UFileUpload accept="text/xml" label="Drop your XML file here" 
            v-model="uploadFile" layout="list" :fileImage="false" position="inside"/>
      </template>
      <template #footer="{ close }">
         <UButton label="Cancel" color="secondary" @click="close" />
         <UButton label="Upload" @click="uploadClicked()" />
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'  

const emit = defineEmits( ['upload' ])
const open = ref(false)
const uploadFile = ref()

const show = ( () => {
   uploadFile.value = null
   open.value = true
})

const uploadClicked = ( async () => {
   emit('upload', uploadFile.value)
   open.value = false
})
</script>

<style scoped lang="scss">
</style>