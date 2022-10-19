<template>
   <div class="details" v-if="systemStore.working==false && unitsStore.masterFiles.length > 0">
      <Panel header="Master Files">
         <DataTable :value="unitsStore.masterFiles" ref="unitMasterFilesTable" dataKey="id"
            showGridlines stripedRows responsiveLayout="scroll" class="p-datatable-sm"
            :lazy="false" :paginator="unitsStore.masterFiles.length > 15" :rows="15" :rowsPerPageOptions="[15,30,50,100]"
            paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
            currentPageReportTemplate="{first} - {last} of {totalRecords}"
            v-model:selection="selectedMasterFiles" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
         >
            <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
            <Column field="pid" header="PID"/>
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
                  <DPGButton v-if="slotProps.data.exemplar==false && (detail.intendedUse.id == 110 || detail.includeInDL)"
                     label="Exemplar" class="p-button-secondary" @click="exemplarClicked(slotProps.data)"/>
               </template>
            </Column>
         </DataTable>
      </Panel>
   </div>
</template>

<script setup>
import { ref } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'

const systemStore = useSystemStore()
const unitsStore = useUnitsStore()
const userStore = useUserStore()
const router = useRouter()

const selectedMasterFiles = ref([])
const selectAll = ref(false)

const { detail } = storeToRefs(unitsStore)

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

</script>

<style scoped lang="scss">
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
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