<template>
   <h2>Master File {{route.params.id}}</h2>
   <div class="edit-form">
      <FormKit type="form" id="masterfile-detail" :actions="false" @submit="submitChanges">
         <div class="split">
            <Panel header="General Information">
               <FormKit label="Title" type="text" v-model="edited.title" required outer-class="first" />
               <FormKit label="Description" type="textarea" :rows="4" v-model="edited.description"/>
               <FormKit label="Orientation" type="select" v-model="edited.orientation" :options="orientations" />
            </Panel>
            <div class="sep"></div>
            <div class="column">
               <Panel header="Related Information">
                  <div class="metadata-lookup">
                     <label>Metadata ID</label>
                     <div class="item">
                        <span>{{displayMetadataID}}</span>
                        <LookupDialog target="metadata" @selected="metadataSelected" class="small-button"/>
                     </div>
                  </div>
                  <template v-if="edited.updateLocation">
                     <FormKit label="Container Type" type="select" v-model="edited.containerTypeID" :options="containerTypes" outer-class="first" />
                     <FormKit label="Container ID" type="text" v-model="edited.containerID"  />
                     <FormKit label="Folder" type="text" v-model="edited.folderID" />
                     <FormKit label="Notes" type="textarea" :rows="3" v-model="edited.notes" />
                  </template>
                  <template v-else>
                     <DPGButton label="Add Location" class="p-button-secondary" @click="edited.updateLocation = true"/>
                  </template>
               </Panel>
            </div>
         </div>
         <div class="acts">
            <DPGButton label="Cancel" class="p-button-secondary" @click="cancelEdit()"/>
            <FormKit type="submit" label="Save" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useSystemStore } from '@/stores/system'
import Panel from 'primevue/panel'
import LookupDialog from '../components/unit/LookupDialog.vue'

const route = useRoute()
const router = useRouter()
const masterFiles = useMasterFilesStore()
const systemStore = useSystemStore()

const edited = ref({
   orientation: 0,
   title: "",
   description: "",
   metadataID: 0,
   updateLocation: false,
   containerTypeID: 0,
   containerID: "",
   folderID: "",
   notes: ""
})

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
   if (edited.value.metadataID && edited.value.metadataID) {
      return edited.value.metadataID
   }
   return "None"
})

onMounted( async () =>{
   let mfID = route.params.id
   await masterFiles.getDetails(mfID)
   document.title = `Edit | Master File ${mfID}`

   edited.value.metadataID = masterFiles.details.metadataID
   edited.value.title = masterFiles.details.title
   edited.value.description = masterFiles.details.description
   edited.value.orientation = masterFiles.details.orientation
   edited.value.updateLocation = false
   if (masterFiles.details.locations.length > 0) {
      edited.value.containerTypeID = masterFiles.details.locations[0].containerType.id
      edited.value.containerID = masterFiles.details.locations[0].containerID
      edited.value.folderID = masterFiles.details.locations[0].folderID
      edited.value.notes = masterFiles.details.locations[0].notes
      edited.value.updateLocation = true
   }
})

function metadataSelected( o ) {
   edited.value.metadataID = o
}
function cancelEdit() {
   router.push(`/units/${route.params.id}`)
}

async function submitChanges() {
   await masterFiles.submitEdit( edited.value )
   if (systemStore.showError == false) {
      router.push(`/masterfiles/${masterFiles.details.id}`)
   }
}
</script>


<style lang="scss" scoped>
.edit-form {
   width: 80%;
   margin: 30px auto 0 auto;

   .metadata-lookup {
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      margin-bottom: 15px;

      label {
         display: block;
         text-align: left;
      }
      .item {
         margin-left: 15px;
         display: flex;
         flex-flow: row nowrap;
         justify-content: flex-start;
         align-items: center;
         margin-top: 10px;
         :deep(button.small-button) {
            padding: 3px 15px;
            font-size: 0.85em;
            margin-left: 10px;
         }
      }
   }

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      :deep(.p-panel), .column {
         flex-grow: 1;
      }
      .sep {
         display: inline-block;
         width: 20px;
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
   padding: 25px 0;
   button {
      margin-right: 10px;
   }
}
</style>
