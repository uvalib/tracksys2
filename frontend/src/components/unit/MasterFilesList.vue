<template>
   <div class="details">
      <Panel header="Master Files" v-if="unitsStore.masterFiles.length > 0">
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
               <div class="master-file-acts" id="sticky-toolbar">
                  <template v-if="detail.reorder==false && userStore.isAdmin">
                     <DPGButton label="Add" @click="addClicked()" class="p-button-secondary" :loading="unitsStore.updateInProgress" />
                     <DPGButton label="Replace" @click="replaceClicked()" class="p-button-secondary" :loading="unitsStore.updateInProgress" />
                  </template>
                  <template v-if="userStore.isAdmin && (detail.dateArchived==null || detail.reorder || detail.dateDLDeliverablesReady == null)">
                     <DPGButton label="Delete Selected" @click="deleteClicked()" class="p-button-secondary" :disabled="!filesSelected" />
                  </template>
                  <RenumberDialog v-if="userStore.isAdmin || userStore.isSupervisor" :disabled="!filesSelected" :filenames="selectedFileNames" />
                  <template v-if="unitsStore.canDownload">
                     <DPGButton label="Download Selected" @click="downloadClicked()" class="p-button-secondary" :disabled="!filesSelected" />
                  </template>
                  <template v-if="unitsStore.canPDF">
                     <DPGButton label="PDF of Selected" @click="pdfClicked()" class="p-button-secondary" :disabled="filesSelected == false" />
                  </template>
                  <template  v-if="userStore.isAdmin || userStore.isSupervisor">
                     <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selected="assignMetadata" target="metadata" :create="true"/>
                     <LookupDialog :disabled="!filesSelected" label="Assign Componment" @selected="assignComponent" target="component" />
                  </template>
               </div>
            </template>
            <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
            <Column field="metadata.title" header="Metadata" filterField="metadata.title" :showFilterMatchModes="false" >
               <template #filter="{filterModel}">
                  <InputText type="text" v-model="filterModel.value" placeholder="Metadata title"/>
               </template>
               <template #body="slotProps">
                  <router-link :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.title}}</router-link>
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
            <Column header="" class="row-acts">
               <template #body="slotProps">
                  <DPGButton label="View" class="p-button-secondary first" @click="viewClicked(slotProps.data)" />
                  <DPGButton label="Download Image" class="p-button-secondary" @click="downloadFile(slotProps.data)" v-if="unitsStore.canDownload"/>
                  <DPGButton label="Download PDF" class="p-button-secondary" @click="downloadPDF(slotProps.data)" v-if="unitsStore.canPDF"/>
                  <DPGButton v-if="slotProps.data.exemplar==false && (detail.intendedUse && detail.intendedUse.id == 110 || detail.includeInDL)"
                     label="Set Exemplar" class="p-button-secondary" @click="exemplarClicked(slotProps.data)"/>
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
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useSystemStore } from '@/stores/system'
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
import CloneMasterFiles from './CloneMasterFiles.vue'
import LookupDialog from '@/components/LookupDialog.vue'
import InputText from 'primevue/inputtext'
import { FilterMatchMode } from 'primevue/api'

const confirm = useConfirm()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()
const router = useRouter()

const selectedMasterFiles = ref([])
const selectAll = ref(false)
const cloneMasterFiles = ref(false)
const toolbarTop = ref(0)
const toolbarHeight = ref(0)
const toolbarWidth = ref(0)
const toolbar = ref(null)

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

onMounted(() => {
   setTimeout( () => {
      let tb = null
      let tbs = document.getElementsByClassName("p-paginator-top")
      if ( tbs ) {
         tb = tbs[0]
      }
      if ( tb) {
         toolbar.value = tb
         toolbarHeight.value = tb.offsetHeight
         toolbarWidth.value = tb.offsetWidth
         toolbarTop.value = 0

         // walk the parents of the toolbar and add each top value
         // to find the top of the toolbar relative to document top
         let ele = tb
         if (ele.offsetParent) {
            do {
               toolbarTop.value += ele.offsetTop
               ele = ele.offsetParent
            } while (ele)
         }
      }
   }, 1000)
   window.addEventListener("scroll", scrollHandler)
})

onUnmounted(() => {
   window.removeEventListener("scroll", scrollHandler)
})

function scrollHandler( ) {
   if ( toolbar.value) {
      if ( window.scrollY <= toolbarTop.value ) {
         if ( toolbar.value.classList.contains("sticky") ) {
            toolbar.value.classList.remove("sticky")
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].classList.add("sticky")
               dts[0].style.top = `0px`
            }
         }
      } else {
         if ( toolbar.value.classList.contains("sticky") == false ) {
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].classList.add("sticky")
               dts[0].style.top = `${toolbarHeight.value}px`
            }
            toolbar.value.classList.add("sticky")
            toolbar.value.style.width = `${toolbarWidth.value}px`
         }
      }
   }
}

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
   requestPDF( ids )
}
const requestPDF = (( masterFileIDs ) => {
   if (unitsStore.hasText == false ) {
      pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs )
   } else {
      confirm.require({
         message: `This unit has transcription or OCR text. Include it with the PDF?`,
         header: 'Include Text',
         icon: 'pi pi-question-circle',
         rejectClass: 'p-button-secondary',
         accept: () => {
            pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs, true )
         },
         reject: () => {
            pdfStore.requestPDF( unitsStore.detail.id, masterFileIDs )
         }
      })
   }
})
function downloadPDF(info) {
   requestPDF( [info.id] )
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
         clearSelections()
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
         clearSelections()
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
   clearSelections()
}

async function assignComponent( componentID ) {
   await unitsStore.assignComponent(componentID, selectedIDs.value)
   clearSelections()
}

function clearSelections() {
   selectAll.value = false
   selectedMasterFiles.value = []
}

</script>

<style scoped lang="scss">
:deep(.p-datatable-wrapper.sticky) {
   position: relative;
}
:deep(.p-paginator-top)  {
   div.p-paginator.p-component {
      padding: 0 0 15px 0;
   }
}
:deep(.p-paginator-top.sticky) {
   position: fixed;
   z-index: 1000;
   top: 0;
   div.p-paginator.p-component {
      padding: 0 0 10px 0;
      border-bottom: 2px solid var(--uvalib-grey-lightest);
      border-radius: 0;
   }
}
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