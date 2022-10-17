<template>
   <h2>Unit {{route.params.id}}</h2>
   <DPGButton label="Edit" class="p-button-secondary edit" @click="editUnit()"/>
   <div v-if="detail.lastError" class="last-error">
      <span>Recent Error:</span>
      <router-link :to="`/jobs/${detail.lastError.jobID}`">{{detail.lastError.error}}</router-link>
   </div>
   <div class="details" v-if="systemStore.working==false">
      <Panel header="General Information">
         <dl>
            <DataDisplay label="Status" :value="detail.status">
               <span :class="`status ${detail.status}`">{{displayStatus(detail.status)}}</span>
            </DataDisplay>
            <DataDisplay label="Patron Source URL" :value="detail.patronSourceURL">
               <a :href="detail.patronSourceURL" target="_blank">{{detail.patronSourceURL}}<i class="icon fas fa-external-link"></i></a>
            </DataDisplay>
            <DataDisplay label="Special Instructions" :value="detail.specialInstructions"/>
            <DataDisplay label="Staff Notes" :value="detail.staffNotes"/>
            <DataDisplay label="Complete Scan" :value="detail.completeScan.toString()">
               <span :class="`flag ${detail.completeScan.toString()}`">{{displayFlag(detail.completeScan)}}</span>
            </DataDisplay>
            <DataDisplay label="Throwaway" :value="detail.throwAway.toString()">
               <span :class="`flag ${detail.throwAway.toString()}`">{{displayFlag(detail.throwAway)}}</span>
            </DataDisplay>
         </dl>
      </Panel>
      <Panel header="Digitization Information">
         <dl>
            <template v-if="detail.intendedUse">
               <DataDisplay label="Intended Use" :value="detail.intendedUse.name"/>
               <DataDisplay label="Deliverable Format" :value="detail.intendedUse.deliverableFormat"/>
               <DataDisplay label="Deliverable Resolution" :value="detail.intendedUse.deliverableResolution"/>
            </template>
            <template v-else>
               <DataDisplay label="Intended Use" value=""/>
               <DataDisplay label="Deliverable Format" value=""/>
               <DataDisplay label="Deliverable Resolution" value=""/>
            </template>
            <DataDisplay label="OCR Master Files" :value="detail.ocrMasterFiles.toString()">
               <span :class="`flag ${detail.ocrMasterFiles.toString()}`">{{displayFlag(detail.ocrMasterFiles)}}</span>
            </DataDisplay>
            <DataDisplay label="Remove Watermark" :value="detail.removeWatermark.toString()">
               <span :class="`flag ${detail.removeWatermark.toString()}`">{{displayFlag(detail.removeWatermark)}}</span>
            </DataDisplay>
            <DataDisplay label="Date Archived" :value="formatDate(detail.dateArchived)" />
            <DataDisplay label="Date Patron Deliverables Ready" :value="formatDate(detail.datePatronDeliverablesReady)" />
            <DataDisplay label="In Digital Library" :value="detail.includeInDL.toString()">
               <span :class="`flag ${detail.includeInDL.toString()}`">{{displayFlag(detail.includeInDL)}}</span>
            </DataDisplay>
            <DataDisplay label="Date DL Deliverables Ready" :value="formatDate(detail.dateDLDeliverablesReady)" />
         </dl>
         <div class="acts" v-if="!detail.projectID">
            <CreateProjectDialog />
         </div>
      </Panel>
   </div>
   <div class="details" v-if="systemStore.working==false">
      <Panel header="Related Information" class="small">
         <dl>
            <DataDisplay v-if="detail.metadata" label="Metadata" :value="detail.metadata.pid">
               <router-link :to="`/metadata/${detail.metadata.id}`">{{detail.metadata.pid}}: {{detail.metadata.title}}</router-link>
            </DataDisplay>
            <DataDisplay v-else label="Metadata" value=""/>
            <DataDisplay label="Order" :value="`${detail.orderID}`">
               <router-link :to="`/orders/${detail.orderID}`">#{{detail.orderID}}</router-link>
            </DataDisplay>
            <DataDisplay v-if="detail.projectID"  label="Project" :value="`${detail.projectID}`">
               <a :href="`${systemStore.projectsURL}/projects/${detail.projectID}`" target="_blank">
                  #{{detail.projectID}}<i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
         </dl>
      </Panel>
      <Panel header="Attachments">
         <div v-if="detail.attachments.length > 0">
            <DataTable :value="detail.attachments" ref="attachmentsTable" dataKey="id"
               stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm attachments-data"
            >
               <Column field="filename" header="File Name" />
               <Column field="description" header="Description" />
               <Column header="" class="row-acts">
                  <template #body="slotProps">
                     <DPGButton label="Download" class="p-button-secondary" @click="downloadAttachment(slotProps.data.id)"/>
                     <DPGButton label="Delete" class="p-button-secondary" @click="deleteAttachment(slotProps.data.id)"/>
                  </template>
               </Column>
            </DataTable>
         </div>
         <div v-else>
            No attachments are associated with this unit.
         </div>
         <p v-if='detail.status != "approved" &&  detail.status != "done"'>
            Attachments cannot be aded to unapproved units.
         </p>
         <div v-else class="toolbar">
            <DPGButton label="Add Attachment" class="p-button-secondary" @click="addAttachmentClicked()"/>
         </div>
      </Panel>
   </div>
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
                     <img :src="slotProps.data.thumbnailURL" />
                  </a>
               </template>
            </Column>
            <Column header="" class="row-acts nowrap">
               <template #body="slotProps">
                  <router-link :to="`/masterfiles/${slotProps.data.id}`">Details</router-link>
               </template>
            </Column>
         </DataTable>
      </Panel>
   </div>
</template>

<script setup>
import { onBeforeMount, ref } from 'vue'
import { useRoute, onBeforeRouteUpdate, useRouter } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import Panel from 'primevue/panel'
import { storeToRefs } from "pinia"
import DataDisplay from '../components/DataDisplay.vue'
import dayjs from 'dayjs'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import CreateProjectDialog from '../components/unit/CreateProjectDialog.vue'

const route = useRoute()
const router = useRouter()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()

const { detail } = storeToRefs(unitsStore)

const selectedMasterFiles = ref([])
const selectAll = ref(false)

onBeforeRouteUpdate(async (to) => {
   let uID = to.params.id
   unitsStore.getDetails( uID )
   unitsStore.getMasterFiles(uID)
})

onBeforeMount(() => {
   let uID = route.params.id
   unitsStore.getDetails( uID )
   unitsStore.getMasterFiles(uID)
   document.title = `Unit #${uID}`
})

function editUnit() {
   router.push(`/units/${route.params.id}/edit`)
}

function downloadAttachment(id) {
   alert(id)
}

function deleteAttachment(id) {
   alert(id)
}

function addAttachmentClicked() {
   alert("not yet implemenetd")
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

function displayStatus( id) {
   if (id == "await_fee") {
      return "Await Fee"
   }
   return id.charAt(0).toUpperCase() + id.slice(1)
}

function displayFlag( f ) {
   if (f) return "Yes"
   return "No"
}

function formatDate( dateStr ) {
   if (dateStr) {
      let d = dayjs(dateStr)
      return d.format("YYYY-MM-DD HH:mm A")
   }
   return ""
}

</script>

<style scoped lang="scss">
button.p-button-secondary.edit {
   position: absolute;
   right:15px;
   top: 15px;
}
.last-error {
   background: var(--uvalib-red-darker);
   padding: 10px;
   color: white;
   span {
      display: inline-block;
      font-weight: bold;
      margin-right: 10px;
   }
   a {
      color: white !important;
   }
}
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   :deep(td.thumb) {
      width: 160px !important;
      text-align: center !important;
   }
   div.p-panel.p-component.small {
      flex: 25%;
   }
   .attachments-data {
      font-size: 0.8em;
   }
   :deep(div.p-panel) {
      margin: 10px;
      flex: 45%;
      text-align: left;
   }
   .toolbar {
      padding: 15px 0 0 0;
      text-align: right;
      font-size: 0.8em;
   }
   .acts {
      font-size: 0.8em;
      text-align: right;
   }
   .row-acts {
      font-size: 0.9em;
      button.p-button {
         padding: 4px 10px;
         width: 100%;
         margin-bottom: 8px;
      }
   }
}
</style>