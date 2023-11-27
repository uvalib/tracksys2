<template>
   <Panel  v-if="apTrustPreservation" header="APTrust Information">
      <template v-if="metadataStore.apTrustStatus">
         <dl v-if="metadataStore.detail.isCollection">
            <DataDisplay label="Requested" :value="formatDate(metadataStore.apTrustStatus.requestedAt)"/>
            <DataDisplay label="Finished" :value="formatDate(metadataStore.apTrustStatus.finishedAt)"/>
            <DataDisplay label="Status" :value="metadataStore.apTrustStatus.status"/>
            <DataDisplay label="Note" value="Initial submission of the collection is complete. Resubmissions are handed at the item level."/>
         </dl>
         <dl v-else>
            <DataDisplay label="Bag" :value="metadataStore.apTrustStatus.bag"/>
            <DataDisplay label="Requested" :value="formatDate(metadataStore.apTrustStatus.requestedAt)"/>
            <DataDisplay label="Submitted" :value="formatDate(metadataStore.apTrustStatus.submittedAt)"/>
            <DataDisplay label="Finished" :value="formatDate(metadataStore.apTrustStatus.finishedAt)"/>
            <DataDisplay label="ID" :value="metadataStore.apTrustStatus.id"/>
            <DataDisplay label="eTag" :value="metadataStore.apTrustStatus.etag"/>
            <DataDisplay label="Object ID" :value="metadataStore.apTrustStatus.objectIdentifier">
               <a class="supplemental" :href="`${systemStore.apTrustURL}/objects?identifier=${metadataStore.apTrustStatus.objectIdentifier}`" target="_blank">
                  {{metadataStore.apTrustStatus.objectIdentifier}}
                  <i class="icon fas fa-external-link"></i>
               </a>
            </DataDisplay>
            <DataDisplay label="Group Identifier" :value="metadataStore.apTrustStatus.groupIdentifier"/>
            <DataDisplay label="Storage" :value="metadataStore.apTrustStatus.storage"/>
            <DataDisplay label="Status" :value="metadataStore.apTrustStatus.status"/>
            <DataDisplay v-if="metadataStore.apTrustStatus.status != 'Success'" label="Note" :value="metadataStore.apTrustStatus.note"/>
         </dl>
      </template>
      <div v-else>
         <div>Preservation has been requested but the item has not been submitted</div>
      </div>
      <div class="apt-acts">
         <DPGButton v-if="canSubmitAPTrust" label="Submit to APTrust" class="p-button-secondary apt-submit" @click="apTrustSubmitClicked" />
         <DPGButton v-else-if="canReubmitAPTrust" label="Resubmit to APTrust" class="p-button-secondary apt-submit" @click="apTrustResubmitClicked" />
         <APTrustReportDialog v-if="canGetAPTrustReport"/>
      </div>
   </Panel>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import DataDisplay from '@/components/DataDisplay.vue'
import APTrustReportDialog from '@/components/metadata/APTrustReportDialog.vue'
import dayjs from 'dayjs'
import { useConfirm } from "primevue/useconfirm"

const confirm = useConfirm()
const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const userStore = useUserStore()

const aptSubmitted = ref(false)

const canGetAPTrustReport = computed(() => {
   if (metadataStore.detail.isCollection == false) return false
   if (metadataStore.apTrustStatus == null) return false
   let yearStr =  metadataStore.apTrustStatus.requestedAt.split("-")[0]
   let year = parseInt(yearStr, 10)
   // group status only works for item submitted from 2023 onwards
   return year >= 2023
})

const canReubmitAPTrust = computed (() => {
   if ( userStore.isAdmin == false || aptSubmitted.value == true) return false
   if ( metadataStore.apTrustStatus == null) return false

   return metadataStore.detail.isCollection == false && (metadataStore.apTrustStatus.status == "Success")
})
const canSubmitAPTrust = computed (() => {
   if ( userStore.isAdmin == false || aptSubmitted.value == true) return false
   if ( metadataStore.apTrustStatus == null) {
      return true
   }
   return metadataStore.detail.isCollection == false && (metadataStore.apTrustStatus.status == "Failed" || metadataStore.apTrustStatus.status == "Canceled")
})
const apTrustPreservation = computed( () => {
   if ( metadataStore.detail.preservationTier && metadataStore.detail.preservationTier.id > 1 ) return true
   return false
})

const apTrustResubmitClicked = ( () => {
   confirm.require({
      message: "Resubmission will replace the current version with a newly created one. This cannot be reversed. Are you sure?",
      header: 'Confirm APTrust Resubmission',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
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
         rejectClass: 'p-button-secondary',
         accept: () => {
            doApTrustSubmission(false)
         },
      })
   } else {
      doApTrustSubmission(false)
   }
})

const doApTrustSubmission = ( async (resubmit) => {
   aptSubmitted.value = true
   await metadataStore.sendToAPTRust(resubmit)
   if (systemStore.error == "") {
      systemStore.toastMessage('Submitted', 'This item has begun the APTrust submission process; check the job status page for updates')
   }
})

const formatDate = (( date ) => {
   if (date) {
      return dayjs(date).format("YYYY-MM-DD HH:mm")
   }
   return ""
})

</script>

<style scoped lang="scss">
.apt-acts {
   margin-top: 10px;
   text-align: right;
}
</style>