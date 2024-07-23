<template>
   <DPGButton label="Get APTrust Status Report"
      class="p-button-secondary apt-submit" @click="apTrustStatusClicked" :loading="apTrust.loadingReport" />
   <Dialog v-model:visible="showReport" header="APTrust Collection Status Report" :modal="true" position="top" style="width: 80%;" >
      <div class="error" v-if="apTrust.collectionStatus.errorMessage">
         {{ apTrust.collectionStatus.errorMessage }}
      </div>
      <div v-else class="resport">
         <div class="report-summary">
            <span>
               <b>Total submitted</b>: {{ apTrust.collectionStatus.totalSubmitted }}
            </span>
            <span>
               <b>Success count:</b> {{ apTrust.collectionStatus.successCount }}
            </span>
            <span>
               <b>Error count</b>: {{ apTrust.collectionStatus.failures.length }}
            </span>
         </div>
         <div v-if="apTrust.collectionStatus.failures.length > 0" class="error-details">
            <label>Errors</label>
            <div class="error-scroller">
               <DataTable :value="apTrust.collectionStatus.failures" ref="apTrustStatusTable" dataKey="id"
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
import { useAPTrustStore } from '@/stores/aptrust'
import { useMetadataStore } from '@/stores/metadata'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useConfirm } from "primevue/useconfirm"

const confirm = useConfirm()
const apTrust = useAPTrustStore()
const metadataStore = useMetadataStore()

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
   await apTrust.getCollectionStatusReport( metadataStore.detail.id )
   showReport.value = true
})
const resubmitClicked = (() => {
   let msg = `Resubmit ${selectedErrors.value.length} selected metadata records. `
   msg += "Depending upon the number and size of items selected, this process may take several days to complete. Are you sure?"
   confirm.require({
      message: msg,
      header: 'Confirm APTrust Resubmission',
      icon: 'pi pi-question-circle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Resubmit'
      },
      accept: ( async () => {
         await apTrust.resubmitCollectionItems( metadataStore.detail.id, selectedErrors.value.map( s => s.id) )
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