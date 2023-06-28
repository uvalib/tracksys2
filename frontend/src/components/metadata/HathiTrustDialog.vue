<template>
   <Dialog v-model:visible="showDialog" :style="{width: '550px'}" header="HathiTrust Status" :modal="true" position="top" @hide="emit('closed')">
      <FormKit type="form" id="hathitrust-detail" :actions="false" @submit="submitChanges">
         <div class="split">
            <FormKit label="Date Requested" type="date" v-model="status.requestedAt" readonly outer-class="first" />
            <span class="sep"/>
            <FormKit label="Date Finished" type="date" v-model="status.finishedAt" outer-class="first"/>
         </div>
         <div class="group">
            <p class="group-name">Metadata</p>
            <div class="split">
               <FormKit label="Date Submitted" type="date" v-model="status.metadataSubmittedAt" outer-class="first" />
               <span class="sep"/>
               <FormKit label="Status" type="text" v-model="status.metadataStatus" outer-class="first" />
            </div>
         </div>
         <div class="group">
            <p class="group-name">Package</p>
            <div class="split">
               <FormKit label="Date Created" type="date" v-model="status.packageCreatedAt" />
               <span class="sep"/>
               <FormKit label="Date Submitted" type="date" v-model="status.packageSubmittedAt" />
               <span class="sep"/>
               <FormKit label="Status" type="text" v-model="status.packageStatus" />
            </div>
         </div>
         <FormKit label="Notes" type="textarea" v-model="status.notes" />
         <p class="error" v-if="metadataStore.error ">{{ metadataStore.error }}</p>
         <div class="acts">
            <DPGButton @click="closeDialog" label="Cancel" class="p-button-secondary"/>
            <FormKit type="submit" label="Update" :wrapper-class="submitClass" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import Dialog from 'primevue/dialog'

const metadataStore = useMetadataStore()

const emit = defineEmits( ['closed' ])

const showDialog = ref(true)

const status = ref({
   requestedAt: null,
   packageCreatedAt: null,
   packageSubmittedAt: null,
   packageStatus: "",
   metadataSubmittedAt: null,
   metadataStatus: "",
   finishedAt: null,
   notes: ""}
)

const submitClass = computed(() => {
   let btnClass = "submit-button"
   if ( metadataStore.working) {
      btnClass += " disabled"
   }
   return btnClass
})

onMounted( () => {
   metadataStore.error = ""
   status.value.requestedAt = cleanDate(metadataStore.hathiTrustStatus.requestedAt)
   status.value.packageCreatedAt = cleanDate(metadataStore.hathiTrustStatus.packageCreatedAt)
   status.value.packageSubmittedAt = cleanDate(metadataStore.hathiTrustStatus.packageSubmittedAt)
   status.value.packageStatus = metadataStore.hathiTrustStatus.packageStatus
   status.value.metadataSubmittedAt = cleanDate(metadataStore.hathiTrustStatus.metadataSubmittedAt)
   status.value.metadataStatus = metadataStore.hathiTrustStatus.metadataStatus
   status.value.finishedAt = cleanDate(metadataStore.hathiTrustStatus.finishedAt)
   status.value.notes = metadataStore.hathiTrustStatus.notes
   showDialog.value = true
})

const cleanDate = ( (dateStr) => {
   if (dateStr) {
      return dateStr.split("T")[0]
   }
   return null
})

const closeDialog = (() => {
   showDialog.value = false
   emit("closed")
})

const submitChanges = ( async () => {
   await metadataStore.updateHathiTrustStatus( status.value )
   if ( metadataStore.error == "" ) {
      closeDialog()
   }
})

</script>

<style scoped lang="scss">
.group {
   border: 1px solid var(--uvalib-grey-lightest);
   padding: 10px;
   margin-top: 10px;
   border-radius: 5px;

   p.group-name {
      margin: 0 0 5px 0;
      font-weight: bold;
   }
}
.split {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: flex-end;
   :deep(.formkit-outer) {
      flex-grow: 1;
   }
   .sep {
      display: inline-block;
      width: 10px;
   }
}
p.error {
   text-align: center;
   font-style: italic;
   color: var(--uvalib-red-emergency);
   margin: 10px 0 0 0;
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 20px 0 10px 0;
   button {
      margin-right: 10px;
   }
}
</style>