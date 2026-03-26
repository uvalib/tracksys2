<template>
   <h2>Master File {{route.params.id}}</h2>
   <div class="edit-form">
      <form id="masterfile-detail" @submit="submitChanges">
         <div class="split">
            <Panel header="General Information">
               <FormField id="mftitle" label="Title">
                  <InputText id="mftitle" v-model="title" type="text" />   
               </FormField>
               <FormField id="mfdesc" label="Desctiption">
                  <Textarea id="mfdesc" v-model="description" autoResize rows="4" /> 
               </FormField>
               <FormField id="orient" label="Orientation">
                  <Select id="orient" v-model="orientation"  :options="orientations" optionLabel="label" optionValue="id" />   
               </FormField>
            </Panel>
            <div class="column">
               <Panel header="Related Information">
                  <div class="metadata-lookup">
                     <label>Metadata ID</label>
                     <div class="item">
                        <span>{{displayMetadataID}}</span>
                        <LookupDialog target="metadata" @selected="metadataSelected"/>
                     </div>
                  </div>
                  <template v-if="updateLocation">
                     <FormField id="containertype" label="Container Type">
                        <Select id="containertype" v-model="containerTypeID"  :options="containerTypes" 
                           optionLabel="label" optionValue="id" placeholder="Select a type"
                        />   
                     </FormField>
                     <FormField id="container" label="Container ID">
                        <InputText id="container" v-model="containerID" type="text" />   
                     </FormField>
                     <FormField id="folder" label="Folder">
                        <InputText id="folder" v-model="folderID" type="text" />   
                     </FormField>
                     <FormField id="locnotes" label="Notes">
                        <Textarea id="locnotes" v-model="notes" autoResize rows="3" /> 
                     </FormField>
                  </template>
                  <template v-else>
                     <DPGButton label="Add Location" severity="secondary" @click="updateLocation = true"/>
                  </template>
               </Panel>
            </div>
         </div>
         <div class="acts">
            <DPGButton label="Cancel" severity="secondary" @click="cancelEdit()"/>
            <DPGButton label="Save" type="submit" />
         </div>
      </form>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, computed } from 'vue'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useSystemStore } from '@/stores/system'
import Panel from 'primevue/panel'
import LookupDialog from '@/components/LookupDialog.vue'

import { useForm } from 'vee-validate'
import FormField from '@/components/FormField.vue'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'

const { values, setValues, resetForm, handleSubmit, defineField } = useForm({})

const [orientation] = defineField('orientation') 
const [title] = defineField('title')
const [description] = defineField('description')
const [updateLocation] = defineField('updateLocation')
const [containerTypeID] = defineField('containerTypeID')
const [containerID] = defineField('containerID')
const [folderID] = defineField('folderID')
const [notes] = defineField('notes')

const route = useRoute()
const router = useRouter()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()

const containerTypes = computed(() => {
   let out = []
   systemStore.containerTypes.forEach( c => {
      out.push( {label: c.name, id: c.id} )
   })
   return out
})

const orientations = computed(() => {
   let out = []
   // { normal: 0, flip_y_axis: 1, rotate90: 2, rotate180: 3, rotate270: 4 }
   out.push( {label: "Normal", id: 0} )
   out.push( {label: "Flip Y Axis", id: 1} )
   out.push( {label: "Rotate 90&deg;", id: 2} )
   out.push( {label: "Rotate 180&deg;", id: 3} )
   out.push( {label: "Rotate 270&deg;", id: 4} )
   return out
})

const displayMetadataID = computed( () => {
   if (values.metadataID) {
      return values.metadataID
   }
   return "None"
})

onMounted( async () =>{
   let mfID = route.params.id
   await masterFiles.getDetails(mfID)
   document.title = `Edit | Master File ${mfID}`

   let vals = {
      metadataID: masterFiles.details.metadataID,
      title: masterFiles.details.title,
      description: masterFiles.details.description,
      orientation: masterFiles.details.orientation,
      updateLocation: false,
      containerTypeID: 1,
   }
   if ( !vals.orientation) {
      vals.orientation = 0
   }
   if (masterFiles.details.locations.length > 0) {
      vals.containerTypeID = masterFiles.details.locations[0].containerType.id
      vals.containerID = masterFiles.details.locations[0].containerID
      vals.folderID = masterFiles.details.locations[0].folderID
      vals.notes = masterFiles.details.locations[0].notes
      vals.updateLocation = true
   }
   resetForm({ values: vals})
   console.log(vals)
})

const metadataSelected = (( metadataID ) => {
   setValues({ metadataID: metadataID})
})

const cancelEdit = (() => {
   router.push(`/masterfiles/${masterFiles.details.id}`)
})

const submitChanges = handleSubmit( async (values) => {
   await masterFiles.submitEdit( values )
   if (systemStore.showError == false) {
      router.push(`/masterfiles/${masterFiles.details.id}`)
   }
})
</script>


<style lang="scss" scoped>
.edit-form {
   width: 80%;
   margin: 30px auto 0 auto;
   form {
      display: flex;
      flex-direction: column;
      gap: 15px;
   }
   :deep(.p-panel-content) {
      display: flex;
      flex-direction: column;
      gap: 15px;
   }

   .metadata-lookup {
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      gap: 15px;
      .item {
         display: flex;
         flex-flow: row nowrap;
         justify-content: flex-start;
         align-items: center;
         gap: 10px;
      }
   }

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      gap: 15px;
      text-align: left;
      :deep(.p-panel), .column {
         flex-grow: 1;
      }
      .related {
         label {
            display: block;
            text-align: left;
         }
         .item {
            text-align: left;
            display: flex;
            flex-flow: row nowrap;
            justify-content: flex-start;
            align-items: center;
         }
      }
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   gap:10px;
}
</style>
