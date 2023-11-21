<template>
   <DPGButton v-if="metadataStore.apTrustStatus" label="Get APTrust Status Report"
      class="p-button-secondary apt-submit" @click="apTrustStatusClicked" :loading="collectionStore.working" />
   <Dialog v-model:visible="showReport" header="APTrust Collection Status Report" :modal="true" position="top" style="width: 80%;" >
      <div class="error" v-if="collectionStore.apTrustStatus.errorMessage">
         {{ collectionStore.apTrustStatus.errorMessage }}
      </div>
      <div v-else class="resport">
         <div class="report-summary">
            <span>
               <b>Total submitted</b>: {{ collectionStore.apTrustStatus.totalSubmitted }}
            </span>
            <span>
               <b>Success count:</b> {{ collectionStore.apTrustStatus.successCount }}
            </span>
            <span>
               <b>Error count</b>: {{ collectionStore.apTrustStatus.failures.length }}
            </span>
         </div>
         <div v-if="collectionStore.apTrustStatus.failures.length > 0" class="error-details">
            <label>Errors</label>
            <div class="error-scroller">
               <DataTable :value="collectionStore.apTrustStatus.failures" ref="apTrustStatusTable" dataKey="id"
                  stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm" :lazy="false"
                  v-model:selection="selectedErrors" :resizableColumns="true" columnResizeMode="fit"
               >
                  <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>
                  <Column field="id" header="Metadata">
                     <template #body="slotProps">
                        <router-link :to="`/metadata/${slotProps.data.id}`">
                           {{ slotProps.data.pid }} - {{ truncateTitle ( slotProps.data.title )}}
                        </router-link>
                     </template>
                  </Column>
                  <Column field="error" header="Error Message" />
               </DataTable>
            </div>
         </div>
         <div v-else class="success">
            All items successfully submitted to APTrust
         </div>
      </div>
      <div class="acts">
         <DPGButton @click="resubmitClicked" label="Resubmit Selected" class="p-button-secondary" :disabled="selectedErrors.length == 0"/>
         <DPGButton @click="closeDialog" label="Close" class="p-button-secondary"/>
      </div>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import { useCollectionsStore } from '@/stores/collections'
import { useMetadataStore } from '@/stores/metadata'
import { useSystemStore } from '@/stores/system'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useConfirm } from "primevue/useconfirm"

const confirm = useConfirm()
const collectionStore = useCollectionsStore()
const metadataStore = useMetadataStore()
const systemStore = useSystemStore()

const showReport = ref(false)
const selectedErrors = ref([])

const truncateTitle = ((t) => {
   if (t.length < 100) return t
   return t.slice(0,100)+"..."
})
const closeDialog = (() => {
   showReport.value = false
})
const apTrustStatusClicked = (async () => {
   await collectionStore.getAPTrustStatus()
   showReport.value = true
})
const resubmitClicked = (() => {
   confirm.require({
      message: `Resubmit ${selectedErrors.value.length} selected metadata records. This process may take several days to complete. Are you sure?`,
      header: 'Confirm APTrust Resubmission',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: ( async () => {
         await collectionStore.apTrustResubmit( selectedErrors.value.map( s => s.id) )
      }),
   })
})
</script>

<style lang="scss" scoped>
.error-details  {
   margin-top: 15px;
   label {
      font-weight: bold;
      margin-bottom: 10px;
      display: block;
   }
   .error-scroller {
      max-height: 300px;
      overflow-y: scroll;
      a {
         color: var(--uvalib-brand-blue-light);
         font-weight: 600;
         text-decoration: none;

         &:hover {
            text-decoration: underline;
         }
      }
   }
}
.success {
   margin: 25px 10px 5px 10px;
   text-align: center;
}
.report-summary {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
}
.error {
   font-size: 1.2em;
   text-align: center;
   margin-top: 10px;
   color: var(--uvalib-red-darker);
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 10px 0 10px 0;
   button {
      margin-left: 10px;
   }
}
</style>