<template>
   <h2>Unit {{route.params.id}}</h2>
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
            <DataDisplay label="Intended Use" :value="detail.intendedUse.name"/>
            <DataDisplay label="Deliverable Format" :value="detail.intendedUse.deliverableFormat"/>
            <DataDisplay label="Deliverable Resolution" :value="detail.intendedUse.deliverableResolution"/>
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
   <div class="details" v-if="systemStore.working==false">
      <Panel header="Master Files">
      </Panel>
   </div>
</template>

<script setup>
import { onBeforeMount } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useUnitsStore } from '@/stores/units'
import Panel from 'primevue/panel'
import { storeToRefs } from "pinia"
import DataDisplay from '../components/DataDisplay.vue'
import dayjs from 'dayjs'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const route = useRoute()
const systemStore = useSystemStore()
const unitsStore = useUnitsStore()

const { detail } = storeToRefs(unitsStore)

onBeforeRouteUpdate(async (to) => {
   let uID = to.params.id
   unitsStore.getDetails( uID )
})

onBeforeMount(() => {
   let uID = route.params.id
   unitsStore.getDetails( uID )
})

function downloadAttachment(id) {
   alert(id)
}

function deleteAttachment(id) {
   alert(id)
}

function addAttachmentClicked() {
   alert("not yet implemenetd")
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
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
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