<template>
   <Panel header="Related Master Files">
      <DataTable :value="component.relatedMasterFiles" ref="componentMasterFilesTable" dataKey="id"
         showGridlines stripedRows responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="true" :alwaysShowPaginator="true" :rows="15"
         :rowsPerPageOptions="[15,30,50,100]" paginatorPosition="top"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:selection="selectedMasterFiles" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
      >
         <template #empty><p class="none">No master files are associated with this component</p></template>
         <template #paginatorstart>
            <div class="master-file-acts">
               <DPGButton label="PDF" @click="pdfClicked()" severity="secondary" :disabled="!filesSelected" />
               <DPGButton label="Download" @click="downloadClicked()" severity="secondary" :disabled="!filesSelected" />
               <template  v-if="userStore.isAdmin || userStore.isSupervisor">
                  <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selected="assignMetadata" target="metadata" :create="true"/>
               </template>
            </div>
         </template>
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID"/>
         <Column field="filename" header="File Name"/>
         <Column field="title" header="Title"/>
         <Column field="description" header="Description"/>
         <Column field="metadata.tiitle" header="Metadata">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.title}}</router-link>
            </template>
         </Column>
         <Column field="thumbnailURL" header="Thumb" class="thumb">
            <template #body="slotProps">
               <a :href="slotProps.data.viewerURL" target="_blank">
                  <img :src="slotProps.data.thumbnailURL" :class="{exemplar: slotProps.data.exemplar}"/>
               </a>
            </template>
         </Column>
      </DataTable>
   </Panel>
   <Dialog v-model:visible="pdfStore.downloading" :modal="true" header="Generating PDF" :style="{width: '350px'}">
      <div class="download">
         <p>PDF generation in progress...</p>
         <ProgressBar :value="pdfStore.percent"/>
      </div>
   </Dialog>
</template>

<script setup>
import { useComponentsStore } from '@/stores/components'
import { useUserStore } from '@/stores/user'
import { usePDFStore } from '@/stores/pdf'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { ref, computed, watch } from 'vue'
import LookupDialog from '@/components/LookupDialog.vue'
import ProgressBar from 'primevue/progressbar'
import Dialog from 'primevue/dialog'
import { usePinnable } from '@/composables/pin'

usePinnable("p-datatable-paginator-top")

const component = useComponentsStore()
const userStore = useUserStore()
const pdfStore = usePDFStore()

const selectedMasterFiles = ref([])
const selectAll = ref(false)

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

watch(() => component.loadingMasterFiles, (newVal) => {
   if ( newVal == false) {
      selectedMasterFiles.value = []
      selectAll.value = false
   }
})

const downloadClicked = (() => {
   component.downloadFromArchive(userStore.computeID, selectedMasterFiles.value[0].unitID, selectedFileNames.value)
   clearSelections()
})

const pdfClicked = (() => {
   let ids = []
   let unitID = 0
   selectedMasterFiles.value.forEach( s => {
      ids.push(s.id)
      if (unitID == 0) {
         unitID = s.unitID
      }
   })
   pdfStore.requestPDF( unitID, ids )
})

const assignMetadata = ( async ( metadataID ) => {
   await component.assignMetadata(metadataID, selectedMasterFiles.value[0].unitID, selectedIDs.value)
   clearSelections()
})

const onRowSelect = (() => {
   selectAll.value = selectedMasterFiles.value < component.relatedMasterFiles.length
})

const onRowUnselect = (() => {
   selectAll.value  = false
})

const onSelectAllChange = ((event) => {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedMasterFiles.value = component.relatedMasterFiles
   }
   else {
      selectedMasterFiles.value = []
   }
})

const clearSelections = (() => {
   selectAll.value = false
   selectedMasterFiles.value = []
})

</script>

<style scoped lang="scss">
:deep(div.p-panel-content) {
   padding-top: 0;
}
p.none {
   text-align:center;
}
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

</style>