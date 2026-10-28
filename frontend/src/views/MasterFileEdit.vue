<template>
   <h2>Master File {{route.params.id}}</h2>
   <UForm :state="state" class="w-full max-w-250 mx-auto! mt-10!" @submit="submitChanges">
      <div class="row-between gap-4 items-start">
         <UCard title="General Information" class="flex-1">
            <div class="column gap-4">
               <UFormField name="title" label="Title" class="grow">   
                  <UInput v-model="state.title" class="w-full"/>   
               </UFormField>
               <UFormField name="description" label="Description" class="grow">   
                  <UTextarea v-model="state.description" class="w-full" :rows="4"/>   
               </UFormField>
               <UFormField name="orientation" label="Orientation" class="grow">   
                  <USelect v-model="state.orientation" :items="orientations" class="w-full"/>   
               </UFormField>
            </div>
         </UCard>

         <UCard title="Related Information" class="flex-1">
            <div class="column gap-4">
               <div class="row-between">
                  <span>Metadata ID: {{displayMetadataID}}</span>
                  <LookupDialog target="metadata" @selected="metadataSelected"/>
               </div>
            
               <template v-if="updateLocation">
                  <UFormField name="containerTypeID" label="Container Type" class="grow">   
                     <USelect v-model="state.containerTypeID" :items="containerTypes" class="w-full"/>   
                  </UFormField>
                  <UFormField name="containerID" label="Container ID" class="grow">   
                     <UInput v-model="state.containerID" class="w-full"/>   
                  </UFormField>
                  <UFormField name="folderID" label="Folder" class="grow">   
                     <UInput v-model="state.folderID" class="w-full"/>   
                  </UFormField>
                  <UFormField name="notes" label="Notes" class="grow">   
                     <UTextarea v-model="state.notes" class="w-full" :rows="4"/>   
                  </UFormField>
               </template>
               <div v-else>
                  <UButton label="Add Location" color="secondary" @click="updateLocation = true"/>
               </div>
            </div>
         </UCard>
      </div>
      <div class="row-right gap-2">
         <UButton label="Cancel" color="secondary" @click="cancelEdit()"/>
         <UButton label="Save" type="submit" />
      </div>
   </UForm>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, computed, ref } from 'vue'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useSystemStore } from '@/stores/system'
import LookupDialog from '@/components/LookupDialog.vue'

const state = ref({
   orientation: 0,
   title: "",
   description: "",
   updateLocation: "",
   containerTypeID: null,
   containerID: null,
   folderID: null,
   notes: "",
})
const updateLocation = ref(false)

const route = useRoute()
const router = useRouter()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()

const containerTypes = computed(() => {
   let out = []
   systemStore.containerTypes.forEach( c => {
      out.push( {label: c.name, value: c.id} )
   })
   return out
})

const orientations = computed(() => {
   let out = []
   // { normal: 0, flip_y_axis: 1, rotate90: 2, rotate180: 3, rotate270: 4 }
   out.push( {label: "Normal", value: 0} )
   out.push( {label: "Flip Y Axis", value: 1} )
   out.push( {label: "Rotate 90&deg;", value: 2} )
   out.push( {label: "Rotate 180&deg;", value: 3} )
   out.push( {label: "Rotate 270&deg;", value: 4} )
   return out
})

const displayMetadataID = computed( () => {
   if (state.value.metadataID) {
      return state.value.metadataID
   }
   return "None"
})

onMounted( async () =>{
   let mfID = route.params.id
   await masterFiles.getDetails(mfID)
   document.title = `Edit | Master File ${mfID}`

   state.value = {
      metadataID: masterFiles.details.metadataID,
      title: masterFiles.details.title,
      description: masterFiles.details.description,
      orientation: masterFiles.details.orientation,
      updateLocation: false,
      containerTypeID: 1,
      containerID: null,
      folderID: null
   }
   if ( !state.value.orientation) {
      state.value.orientation = 0
   }

    updateLocation.value = false
   if (masterFiles.details.locations.length > 0) {
      state.value.containerTypeID = masterFiles.details.locations[0].containerType.id
      state.value.containerID = masterFiles.details.locations[0].containerID
      state.value.folderID = masterFiles.details.locations[0].folderID
      state.value.notes = masterFiles.details.locations[0].notes
      updateLocation.value = true
   }
})

const metadataSelected = (( metadataID ) => {
   setValues({ metadataID: metadataID})
})

const cancelEdit = (() => {
   router.push(`/masterfiles/${masterFiles.details.id}`)
})

const submitChanges = ( async () => {
   await masterFiles.submitEdit( state.value )
   if (systemStore.showError == false) {
      router.push(`/masterfiles/${masterFiles.details.id}`)
   }
})
</script>

<style lang="scss" scoped>
</style>
