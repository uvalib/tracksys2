<template>
   <div v-if="apTrust.working" class="loading"><WaitSpinner :overlay="false" message="Processing APTRust request..." /></div>
   <template v-else-if="apTrust.hasItemStatus">
      <dl v-if="metadataStore.detail.isCollection">
         <DataDisplay label="Requested" :value="$formatDateTime(apTrust.itemStatus.requestedAt)"/>
         <DataDisplay label="Finished" :value="$formatDateTime(apTrust.itemStatus.finishedAt)"/>
         <DataDisplay label="Status" :value="apTrust.itemStatus.status"/>
         <DataDisplay label="Note" value="Initial submission of the collection is complete. Resubmissions are handed at the item level."/>
      </dl>
      <dl v-else>
         <DataDisplay label="Bag" :value="apTrust.itemStatus.bag"/>
         <DataDisplay label="Requested" :value="$formatDateTime(apTrust.itemStatus.requestedAt)"/>
         <DataDisplay label="Submitted" :value="$formatDateTime(apTrust.itemStatus.submittedAt)"/>
         <DataDisplay label="Finished" :value="$formatDateTime(apTrust.itemStatus.finishedAt)"/>
         <DataDisplay label="ID" :value="apTrust.itemStatus.id"/>
         <DataDisplay label="eTag" :value="apTrust.itemStatus.etag"/>
         <DataDisplay label="Object ID" :value="apTrust.itemStatus.objectIdentifier">
            <a class="supplemental" :href="`${systemStore.apTrustURL}/objects?identifier=${apTrust.itemStatus.objectIdentifier}`" target="_blank">
               {{apTrust.itemStatus.objectIdentifier}}
               <i class="icon fas fa-external-link"></i>
            </a>
         </DataDisplay>
         <DataDisplay label="Group Identifier" :value="apTrust.itemStatus.groupIdentifier"/>
         <DataDisplay label="Storage" :value="apTrust.itemStatus.storage"/>
         <DataDisplay label="Status" :value="apTrust.itemStatus.status"/>
         <DataDisplay v-if="apTrust.itemStatus.status != 'Success'" label="Note" :value="apTrust.itemStatus.note"/>
      </dl>
      <p class="error" v-if="apTrust.itemStatus.errorMessage">apTrust.itemStatus.errorMessage</p>
   </template>
   <div v-else>
      <div>Preservation has been requested but the item has not been submitted</div>
   </div>
   <div class="apt-acts"  v-if="apTrust.working == false && props.readonly==false">
      <DPGButton v-if="canSubmitAPTrust" label="Submit to APTrust" severity="secondary" @click="apTrustSubmitClicked" :loading="apTrust.working" />
      <DPGButton v-else-if="canReubmitAPTrust" label="Resubmit to APTrust" severity="secondary" @click="apTrustResubmitClicked" :loading="apTrust.working" />
      <APTrustReportDialog v-if="canGetAPTrustReport"/>
   </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useUserStore } from '@/stores/user'
import { useAPTrustStore } from '@/stores/aptrust'
import WaitSpinner from "@/components/WaitSpinner.vue"
import DataDisplay from '@/components/DataDisplay.vue'
import APTrustReportDialog from '@/components/aptrust/APTrustReportDialog.vue'
import { useConfirm } from "primevue/useconfirm"


const props = defineProps({
   readonly: {
      type: Boolean,
      default: false
   },
})

const confirm = useConfirm()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const userStore = useUserStore()
const apTrust = useAPTrustStore()

const aptSubmitted = ref(false)

const canGetAPTrustReport = computed(() => {
   if (metadataStore.detail.isCollection == false) return false
   if (apTrust.hasItemStatus == false) return false
   let yearStr =  apTrust.itemStatus.requestedAt.split("-")[0]
   let year = parseInt(yearStr, 10)
   // group status only works for item submitted from 2023 onwards
   return year >= 2023
})

const canReubmitAPTrust = computed (() => {
   if ( userStore.isAdmin == false || aptSubmitted.value == true || metadataStore.detail.isCollection) return false
   if ( apTrust.hasItemStatus == false) return false
   return  apTrust.itemStatus.status == "Success"
})

const canSubmitAPTrust = computed (() => {
   if ( userStore.isAdmin == false || aptSubmitted.value == true ) return false
   if ( apTrust.hasItemStatus == false) return true
   return metadataStore.detail.isCollection == false && (apTrust.itemStatus.status == "Failed" || apTrust.itemStatus.status == "Canceled")
})

const apTrustResubmitClicked = ( () => {
   confirm.require({
      message: "Resubmission will replace the current version with a newly created one. This cannot be reversed. Are you sure?",
      header: 'Confirm APTrust Resubmission',
      icon: 'pi pi-question-circle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Resubmit'
      },
      accept: () => {
         doApTrustSubmission(true)
      },
   })
})

const apTrustSubmitClicked = ( () => {
   if (metadataStore.detail.isCollection) {
      confirm.require({
         message: "Submitting a collection record to APTrust will also submit all collection items. Are you sure?",
         header: 'Confirm APTrust Submission',
         icon: 'pi pi-question-circle',
         rejectProps: {
            label: 'Cancel',
            severity: 'secondary'
         },
         acceptProps: {
            label: 'Submit'
         },
         accept: () => {
            doApTrustSubmission(false)
         },
      })
   } else {
      doApTrustSubmission(false)
   }
})

const doApTrustSubmission = ( async (resubmit) => {
   await apTrust.submitItem( metadataStore.detail.id, resubmit)
   if (systemStore.error == "") {
      systemStore.toastMessage('Submitted', 'This item has begun the APTrust submission process; check the job status page for updates')
      aptSubmitted.value = true
   } else {
      aptSubmitted.value = false
   }
})

</script>

<style scoped lang="scss">
.loading{
   text-align: center;
   font-size: 0.7em;
}
.apt-acts {
   margin-top: 10px;
   text-align: right;
}
dl {
   margin: 10px 30px 0 30px;
   display: inline-grid;
   grid-template-columns: max-content 2fr;
   grid-column-gap: 10px;
   font-size: 0.9em;
   text-align: left;
   box-sizing: border-box;
}
</style>