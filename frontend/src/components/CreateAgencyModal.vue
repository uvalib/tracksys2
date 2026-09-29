<template>
   <UModal v-model:open="open" :modal="true" :dismissible="false" :close="false" title="Create Agency">
      <UButton @click="show()" label="Create Agency" />
      <template #body>
         <div class="agency">
            <div class="row">
               <label>Name</label>
               <UInput v-model="newAgencyName" autofocus/>
            </div>
            <div class="row">
               <label>Description</label>
               <UTextarea :rows="4" v-model="newAgencyDesc"/>
            </div>
         </div>
      </template>
      <template #footer="{ close }">
         <UButton label="Cancel" color="secondary" @click="close" />
         <UButton label="Create" @click="createAgency()" />
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useSystemStore } from '../stores/system'

const system = useSystemStore()

const open = ref(false)
const newAgencyName = ref("")
const newAgencyDesc = ref("")

const show = ( () => {
   newAgencyDesc.value = ""
   newAgencyName.value = ""
   open.value = true
})

const createAgency = ( async () => {
   await system.createAgency(newAgencyName.value, newAgencyDesc.value)
   open.value = false
})
</script>

<style scoped lang="scss">
.agency {
   display: flex;
   flex-direction: column;
   gap: 10px;
   .row {
      display: flex;
      flex-direction: column;
      gap: 5px;  
   }
}
</style>
