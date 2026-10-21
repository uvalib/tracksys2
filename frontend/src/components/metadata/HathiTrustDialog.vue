<template>
   <UModal v-model:open="showDialog" :modal="true" :dismissible="false" :close="false" title="HathiTrustStatus" >
      <UButton @click="showClicked" label="HathiTrust Status" color="secondary" />
      <template #body>
         <UForm :state="state" @submit="submitChanges" class="column gap-4 w-full">
            <div class="row-left gap-4">
               <UFormField name="requestedAt" label="Date Requested" class="grow">  
                  <UInput v-model="state.requestedAt" :readonly="true" variant="outline" class="w-full"/>   
               </UFormField>
               <UFormField name="finishedAt" label="Date Finished" class="grow">  
                  <DatePicker v-model="state.finishedAt" class="w-full"/>   
               </UFormField>
            </div>
            <UCard title="Metadata" :ui="{body: 'pt-4!'}">
               <div class="row-left gap-4">
                  <UFormField name="metadataSubmittedAt" label="Date Submitted" class="grow">  
                     <DatePicker v-model="state.metadataSubmittedAt" class="w-full"/>   
                  </UFormField>
                  <UFormField name="metadataStatus" label="Status" class="grow">  
                     <USelect v-model="state.metadataStatus" :items='statuses' class="w-full"/>   
                  </UFormField>
               </div>
            </UCard>
            <UCard title="Package" :ui="{body: 'pt-4!'}">
               <div class="flex flex-row flex-wrap justify-start gap-4">
                  <UFormField name="packageCreatedAt" label="Date Created" class="grow">  
                     <DatePicker v-model="state.packageCreatedAt" class="w-full"/>   
                  </UFormField>
                  <UFormField name="packageSubmittedAt" label="Date Submitted" class="grow">  
                     <DatePicker v-model="state.packageSubmittedAt" class="w-full"/>   
                  </UFormField>
                  <UFormField name="packageStatus" label="Status" class="grow">  
                     <USelect v-model="state.packageStatus" :items='statuses' class="w-full"/>   
                  </UFormField>
               </div>
            </UCard>
            <div class="row-right gap-2">
               <UButton @click="closeDialog" label="Cancel" color="secondary"/>
               <UButton label="Update" type="submit" />  
            </div>
         </UForm>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import DatePicker from '@/components/DatePicker.vue'
import { toCalendarDate, parseAbsolute } from '@internationalized/date'

const state = ref({
   requestedAt: null,
   packageCreatedAt: null,
   packageSubmittedAt: null,
   packageStatus: "",
   metadataSubmittedAt: null,
   metadataStatus: "",
   finishedAt: null,
   notes: ""
})

const statuses = ['pending', 'submitted', 'accepted', 'failed']

const metadataStore = useMetadataStore()

const emit = defineEmits( ['closed' ])

const showDialog = ref(false)

const showClicked = (() => {
   metadataStore.error = ""
   state.value = {
      packageStatus: metadataStore.hathiTrustStatus.packageStatus,
      metadataStatus: metadataStore.hathiTrustStatus.metadataStatus,
      notes: metadataStore.hathiTrustStatus.notes,
   }
   if ( metadataStore.hathiTrustStatus.requestedAt != null) {
      state.value.requestedAt = toCalendarDate( parseAbsolute(metadataStore.hathiTrustStatus.requestedAt) )
   }
   if ( metadataStore.hathiTrustStatus.metadataSubmittedAt != null) {
      state.value.metadataSubmittedAt = toCalendarDate( parseAbsolute(metadataStore.hathiTrustStatus.metadataSubmittedAt) )
   }  
   if ( metadataStore.hathiTrustStatus.packageCreatedAt != null) {
      state.value.packageCreatedAt = toCalendarDate( parseAbsolute(metadataStore.hathiTrustStatus.packageCreatedAt) )
   } 
   if ( metadataStore.hathiTrustStatus.packageSubmittedAt != null) {
      state.value.packageSubmittedAt = toCalendarDate( parseAbsolute(metadataStore.hathiTrustStatus.packageSubmittedAt) )
   } 
   if ( metadataStore.hathiTrustStatus.finishedAt != null) {
      state.value.finishedAt = toCalendarDate( parseAbsolute(metadataStore.hathiTrustStatus.finishedAt) )
   }  
   showDialog.value = true
})

const closeDialog = (() => {
   showDialog.value = false
   emit("closed")
})

const submitChanges = ( async () => {
   let update = {
      requestedAt: state.requestedAt.toString(),
      packageCreatedAt: state.packageCreatedAt.toString(),
      packageSubmittedAt: state.packageSubmittedAt.toString(),
      packageStatus: state.packageStatus,
      metadataSubmittedAt: state.metadataSubmittedAt.toString(),
      metadataStatus: state.metadataStatus,
      finishedAt: state.finishedAt.toString(),
      notes: state.notes
   }
   await metadataStore.updateHathiTrustStatus( update )
   if ( metadataStore.error == "" ) {
      closeDialog()
   }
})

</script>

<style scoped lang="scss">
.group {
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 0.3rem;
   padding: 15px;
   display: flex;
   flex-direction: column;
   gap: 10px;

   p.group-name {
      font-weight: bold;
      margin:0;
      padding: 0;
   }
}
.split {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: flex-end;
   gap: 15px;
}

</style>