<template>
   <Panel header="Master Files" class="masterfiles">
      <DataTable :value="unitsStore.masterFiles" ref="unitMasterFilesTable" dataKey="id"
         showGridlines stripedRows responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="true" :alwaysShowPaginator="true" :rows="15"
         :rowsPerPageOptions="[15,30,50,100]" paginatorPosition="top"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:selection="selectedMasterFiles" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
         v-model:filters="filters" filterDisplay="menu"
      >
         <template #paginatorstart>
            <div class="master-file-acts">
               <template v-if="detail.reorder==false && userStore.isAdmin">
                  <DPGButton label="Add" @click="addClicked()" severity="secondary" :loading="unitsStore.updateInProgress" />
                  <DPGButton label="Replace" @click="replaceClicked()" severity="secondary" :loading="unitsStore.updateInProgress" />
               </template>
               <template v-if="userStore.isAdmin && (detail.dateArchived==null || detail.reorder || detail.dateDLDeliverablesReady == null)">
                  <DPGButton label="Delete" @click="deleteClicked()" severity="secondary" :disabled="!filesSelected" />
               </template>
               <RenumberDialog v-if="userStore.isAdmin || userStore.isSupervisor" :disabled="!filesSelected" :filenames="selectedFileNames" />
               <template v-if="unitsStore.canDownload">
                  <DPGButton label="Download" @click="downloadClicked()" severity="secondary" :disabled="!filesSelected" />
               </template>
               <template v-if="unitsStore.canPDF">
                  <DPGButton label="PDF" @click="pdfClicked()" severity="secondary" :disabled="filesSelected == false" />
               </template>
               <template  v-if="userStore.isAdmin || userStore.isSupervisor">
                  <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selected="assignMetadata" target="metadata" :create="true"/>
                  <LookupDialog :disabled="!filesSelected" label="Assign Componment" @selected="assignComponent" target="component" />
               </template>
            </div>
         </template>
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="metadata.title" header="Metadata" filterField="metadata.title" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Metadata title"/>
            </template>
            <template #body="slotProps">
               <router-link v-if="slotProps.data.metadata" :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.title}}</router-link>
               <span v-else class="empty">N/A</span>
            </template>
         </Column>
         <Column field="filename" header="File Name"/>
         <Column field="title" header="Title" filterField="title" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Title"/>
            </template>
         </Column>
         <Column field="description" header="Description" filterField="description" :showFilterMatchModes="false" >
            <template #filter="{filterModel}">
               <InputText type="text" v-model="filterModel.value" placeholder="Description"/>
            </template>
         </Column>
         <Column field="thumbnailURL" header="Thumb" class="thumb">
            <template #body="slotProps">
               <a :href="slotProps.data.viewerURL" target="_blank">
                  <img :src="slotProps.data.thumbnailURL" :class="{exemplar: slotProps.data.exemplar}"/>
               </a>
            </template>
         </Column>
         <Column header="Actions" class="row-acts">
            <template #body="slotProps">
               <DPGButton label="View" severity="secondary" @click="viewClicked(slotProps.data)" size="small"/>
               <DPGButton label="Download Image" severity="secondary" @click="downloadFile(slotProps.data)" v-if="unitsStore.canDownload" size="small"/>
               <DPGButton label="Download PDF" severity="secondary" @click="downloadPDF(slotProps.data)" v-if="unitsStore.canPDF" size="small"/>
               <DPGButton v-if="slotProps.data.exemplar==false && (detail.intendedUse && detail.intendedUse.id == 110 || detail.includeInDL)"
                  label="Set Exemplar" severity="secondary" @click="exemplarClicked(slotProps.data)" size="small"/>
               <DPGButton label="Republish IIIF" severity="secondary" @click="republishIIIF(slotProps.data.id)" v-if="detail.reorder==false && userStore.isAdmin" size="small"/>
            </template>
         </Column>
      </DataTable>
   </Panel>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useConfirm } from "primevue/useconfirm"
import RenumberDialog from './RenumberDialog.vue'
import LookupDialog from '@/components/LookupDialog.vue'
import InputText from 'primevue/inputtext'
import { FilterMatchMode } from '@primevue/core/api'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const confirm = useConfirm()
const unitsStore = useUnitsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()
const router = useRouter()

const selectedMasterFiles = ref([])
const selectAll = ref(false)

const { detail } = storeToRefs(unitsStore)

const filters = ref( {
   'metadata.title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'title': {value: null, matchMode: FilterMatchMode.CONTAINS},
   'description': {value: null, matchMode: FilterMatchMode.CONTAINS},
})

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

const downloadClicked = (() => {
   unitsStore.downloadFromArchive(userStore.computeID, selectedFileNames.value)
})
const pdfClicked = (() => {
   let ids = []
   selectedMasterFiles.value.forEach( s => {
      ids.push(s.id)
   })
   requestPDF( ids )
})
const requestPDF = (( masterFileIDs ) => {
   if (unitsStore.hasText == false ) {
      pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs )
   } else {
      confirm.require({
         message: `This unit has transcription or OCR text. Include it with the PDF?`,
         header: 'Include Text',
         icon: 'pi pi-question-circle',
         rejectProps: {
            label: 'Cancel',
            severity: 'secondary'
         },
         acceptProps: {
            label: 'Include'
         },
         accept: () => {
            pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs, true )
         },
         reject: () => {
            pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs )
         }
      })
   }
})
const downloadPDF = ((info) => {
   requestPDF( [info.id] )
})

const replaceClicked = (() => {
   let unitDir = `${unitsStore.detail.id}`.padStart(9, '0')
   confirm.require({
      message: `Replace master files with .tif files from ./finalization/unit_update/${unitDir}?`,
      header: 'Confirm Replace Master Files',
      icon: 'pi pi-question-circle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Replace'
      },
      accept: async () => {
         unitsStore.replaceMasterFiles()
         clearSelections()
      }
   })
})

const addClicked = (() => {
   let unitDir = `${unitsStore.detail.id}`.padStart(9, '0')
   confirm.require({
      message: `Add all .tif files from ./finalization/unit_update/${unitDir} to this unit?`,
      header: 'Confirm Add Master Files',
      icon: 'pi pi-question-circle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Add'
      },
      accept: async () => {
         unitsStore.addMasterFiles()
      }
   })
})

const deleteClicked = (() => {
   confirm.require({
      message: 'Are you sure you want delete the selected master files? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Master Files',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: async () => {
         unitsStore.deleteMasterFiles(selectedFileNames.value)
         clearSelections()
      }
   })
})
const downloadFile =(( info) => {
   unitsStore.downloadFromArchive(userStore.computeID, info.filename )
})
const viewClicked = ((data) => {
   router.push(`/masterfiles/${data.id}`)
})
const exemplarClicked = (( info) => {
   unitsStore.setExemplar( info.id )
})
const republishIIIF =((masterFileID) => {
   unitsStore.regenerateIIIF( masterFileID )
})
const onRowSelect = (() => {
   selectAll.value = selectedMasterFiles.value < unitsStore.masterFiles.length
})
const onRowUnselect = (() => {
   selectAll.value  = false
})
const onSelectAllChange = ((event) => {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedMasterFiles.value = unitsStore.masterFiles
   }
   else {
      selectedMasterFiles.value = []
   }
})

const assignMetadata = ( async ( metadataID ) => {
   await unitsStore.assignMetadata(metadataID, selectedIDs.value)
   clearSelections()
})

const assignComponent= ( async ( componentID ) => {
   await unitsStore.assignComponent(componentID, selectedIDs.value)
   clearSelections()
})

const clearSelections = (() => {
   selectAll.value = false
   selectedMasterFiles.value = []
})

</script>

<style scoped lang="scss">
div.masterfiles {
   .master-file-acts {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-start;
      gap: 5px;
      font-size: 0.85em;
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
   display: flex;
   flex-direction: column;
   gap: 5px;
   button {
      font-size: 0.85em;
   }
}
</style>