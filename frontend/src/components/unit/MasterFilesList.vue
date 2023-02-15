<template>
   <div class="details">
      <Panel header="Master Files" v-if="unitsStore.masterFiles.length > 0">
         <div class="master-file-acts">
            <template v-if="detail.reorder==false && userStore.isAdmin">
               <DPGButton label="Add" @click="addClicked()" class="p-button-secondary" :loading="unitsStore.updateInProgress" />
               <DPGButton label="Replace" @click="replaceClicked()" class="p-button-secondary" :loading="unitsStore.updateInProgress" />
            </template>
            <template v-if="userStore.isAdmin && (detail.dateArchived==null || detail.reorder || detail.dateDLDeliverablesReady == null)">
               <DPGButton label="Delete Selected" @click="deleteClicked()" class="p-button-secondary" :disabled="!filesSelected" />
            </template>
            <RenumberDialog v-if="userStore.isAdmin || userStore.isSupervisor" :disabled="!filesSelected" :filenames="selectedFileNames" />
            <template v-if="detail.dateArchived != null || (detail.reorder && detail.datePatronDeliverablesReady)">
               <DPGButton label="Download Selected" @click="downloadClicked()" class="p-button-secondary" :disabled="!filesSelected" />
            </template>
            <template v-if="detail.metadata && detail.dateArchived != null && detail.reorder == false">
               <DPGButton label="PDF of Selected" @click="pdfClicked()" class="p-button-secondary" :disabled="!filesSelected" />
            </template>
            <template  v-if="userStore.isAdmin || userStore.isSupervisor">
               <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selected="assignMetadata" target="metadata" :create="true"/>
               <LookupDialog :disabled="!filesSelected" label="Assign Componment" @selected="assignComponent" target="component" />
            </template>
         </div>
         <DataTable :value="unitsStore.masterFiles" ref="unitMasterFilesTable" dataKey="id"
            showGridlines stripedRows responsiveLayout="scroll" class="p-datatable-sm"
            :lazy="false" :paginator="unitsStore.masterFiles.length > 15" :rows="15" :rowsPerPageOptions="[15,30,50,100]"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            v-model:selection="selectedMasterFiles" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
         >
            <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
            <Column field="metadata.title" header="Metadata">
               <template #body="slotProps">
                  <router-link :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.title}}</router-link>
               </template>
            </Column>
            <Column field="filename" header="File Name"/>
            <Column field="title" header="Title"/>
            <Column field="description" header="Description"/>
            <Column field="thumbnailURL" header="Thumb" class="thumb">
               <template #body="slotProps">
                  <a :href="slotProps.data.viewerURL" target="_blank">
                     <img :src="slotProps.data.thumbnailURL" :class="{exemplar: slotProps.data.exemplar}"/>
                  </a>
               </template>
            </Column>
            <Column header="" class="row-acts">
               <template #body="slotProps">
                  <DPGButton label="View" class="p-button-secondary first" @click="viewClicked(slotProps.data)" />
                  <DPGButton label="Download" class="p-button-secondary" @click="downloadFile(slotProps.data)"/>
                  <DPGButton v-if="slotProps.data.exemplar==false && (detail.intendedUse && detail.intendedUse.id == 110 || detail.includeInDL)"
                     label="Exemplar" class="p-button-secondary" @click="exemplarClicked(slotProps.data)"/>
               </template>
            </Column>
         </DataTable>
      </Panel>
      <Panel header="Master Files" v-else>
         <template v-if="cloneMasterFiles == false">
            <p>No master files are associated with this unit.</p>
            <DPGButton label="Clone Existing Master Files" class="p-button-secondary" @click="cloneClicked()" />
         </template>
         <CloneMasterFiles v-else @canceled="cloneCanceled()" @cloned="cloneCompleted()" />
      </Panel>
   </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useConfirm } from "primevue/useconfirm"
import RenumberDialog from './RenumberDialog.vue'
import CloneMasterFiles from './CloneMasterFiles.vue'
import LookupDialog from '@/components/LookupDialog.vue'

const confirm = useConfirm()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()
const userStore = useUserStore()
const router = useRouter()

const selectedMasterFiles = ref([])
const selectAll = ref(false)
const cloneMasterFiles = ref(false)

const { detail } = storeToRefs(unitsStore)

const filesSelected = computed(() => {
   return selectedMasterFiles.value.length > 0
})
const selectedFileNames = computed(() => {
   let filenames = []
   selectedMasterFiles.value.forEach( s => {
      filenames.push(s.filename)
   })
   return filenames
})
const selectedIDs = computed(() => {
   let ids = []
   selectedMasterFiles.value.forEach( s => {
      ids.push(s.id)
   })
   return ids
})

function cloneClicked() {
   cloneMasterFiles.value = true
}
function cloneCanceled() {
   cloneMasterFiles.value = false
}
function cloneCompleted() {
   cloneMasterFiles.value = false
   systemStore.toastMessage("Clone Success", 'All master files have been cloned.')
}
function downloadClicked() {
   unitsStore.downloadFromArchive(userStore.computeID, selectedFileNames.value)
}
function pdfClicked() {
   let ids = []
   selectedMasterFiles.value.forEach( s => {
      ids.push(s.id)
   })
   let token = new Date().getTime()
   let url = `${systemStore.pdfURL}/${unitsStore.detail.metadata.pid}?unit=${unitsStore.detail.id}&token=${token}&pages=${ids.join(',')}`
   window.open(url)
}
function replaceClicked() {
   let unitDir = `${unitsStore.detail.id}`.padStart(9, '0')
   confirm.require({
      message: `Replace master files with .tif files from ./finalization/unit_update/${unitDir}?`,
      header: 'Confirm Replace Master Files',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         unitsStore.replaceMasterFiles()
      }
   })
}

function addClicked() {
   let unitDir = `${unitsStore.detail.id}`.padStart(9, '0')
   confirm.require({
      message: `Add all .tif files from ./finalization/unit_update/${unitDir} to this unit?`,
      header: 'Confirm Add Master Files',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         unitsStore.addMasterFiles()
      }
   })
}

function deleteClicked() {
   confirm.require({
      message: 'Are you sure you want delete the selected master files? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Master Files',
      icon: 'pi pi-exclamation-triangle',
      rejectClass: 'p-button-secondary',
      accept: async () => {
         unitsStore.deleteMasterFiles(selectedFileNames.value)
      }
   })
}

function downloadFile( info) {
   unitsStore.downloadFromArchive(userStore.computeID, info.filename )
}
function viewClicked(data) {
   router.push(`/masterfiles/${data.id}`)
}
function exemplarClicked( info) {
   unitsStore.setExemplar( info.id )
}
function onRowSelect() {
   selectAll.value = selectedMasterFiles.value < unitsStore.masterFiles.length
}
function onRowUnselect() {
   selectAll.value  = false
}
function onSelectAllChange(event) {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedMasterFiles.value = unitsStore.masterFiles
   }
   else {
      selectedMasterFiles.value = []
   }
}

async function assignMetadata( metadataID ) {
   await unitsStore.assignMetadata(metadataID, selectedIDs.value)
}

async function assignComponent( componentID ) {
   await unitsStore.assignComponent(componentID, selectedIDs.value)
}

</script>

<style scoped lang="scss">
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;

   .p-datatable-sm {
      font-size: 0.9em;
   }

   .master-file-acts {
      font-size: 0.85em;
      text-align: right;
      padding: 0 0 15px 0;

      button.p-button {
         margin-left: 5px;
      }
   }

   :deep(td.thumb) {
      width: 160px !important;
      text-align: center !important;
   }
   img.exemplar {
      border: 5px solid var(--uvalib-teal);
   }
}
:deep(td.row-acts) {
   vertical-align: top;

   button.p-button.first {
      margin: 0;
   }
   button.p-button {
      font-size: 0.75em;
      padding: 3px 6px;
      display: block;
      width: 100%;
      margin-top: 5px;
   }
}
</style>