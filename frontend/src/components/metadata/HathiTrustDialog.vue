<template>
   <Dialog v-model:visible="showDialog" :style="{width: '550px'}" header="HathiTrust Status" :modal="true" position="top" @hide="emit('closed')">
      <form id="hathitrust-detail" @submit="submitChanges">
         <div class="split">
            <FormField id="drequest" label="Date Requested">
               <InputText id="drequest" v-model="requestedAt" type="text" :readonly="true"/>   
            </FormField>
            <FormField id="dfinish" label="Date Finished">
               <InputText id="dfinish" v-model="finishedAt" type="date"/>   
            </FormField>
         </div>
        <div class="group">
            <p class="group-name">Metadata</p>
            <div class="split">
               <FormField id="dsubmit" label="Date Submitted">
                  <InputText id="dsubmit" v-model="metadataSubmittedAt" type="date"/>   
               </FormField>
               <FormField id="mdstatus" label="Status">
                  <InputText id="mdstatus" v-model="metadataStatus" type="text"/>   
               </FormField>
            </div>
         </div>
         <div class="group">
            <p class="group-name">Package</p>
            <div class="split">
               <FormField id="pcreate" label="Date Created">
                  <InputText id="pcreate" v-model="packageCreatedAt" type="date"/>   
               </FormField>
               <FormField id="psubmit" label="Date Submitted">
                  <InputText id="psubmit" v-model="packageSubmittedAt" type="date"/>   
               </FormField>
               <FormField id="pstatus" label="Status">
                  <InputText id="pstatus" v-model="packageStatus" type="text"/>   
               </FormField>
            </div>
         </div>
         <FormField id="htnotes" label="Notes">
            <Textarea id="htnotes" v-model="notes" autoResize rows="3" /> 
         </FormField>
         <Message v-if="metadataStore.error" severity="error" size="small" variant="simple">{{metadataStore.error}}</Message>
         <div class="acts">
            <DPGButton @click="closeDialog" label="Cancel" severity="secondary"/>
            <DPGButton label="Update" type="submit" />
         </div>
      </form>
   </Dialog>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import { useForm } from 'vee-validate'
import FormField from '@/components/FormField.vue'

const { resetForm, handleSubmit, defineField } = useForm({})

const [requestedAt] = defineField('requestedAt')
const [packageCreatedAt] = defineField('packageCreatedAt')
const [packageSubmittedAt] = defineField('packageSubmittedAt')
const [packageStatus] = defineField('packageStatus')
const [metadataSubmittedAt] = defineField('metadataSubmittedAt')
const [metadataStatus] = defineField('metadataStatus')
const [finishedAt] = defineField('finishedAt')
const [notes] = defineField('notes')

const metadataStore = useMetadataStore()

const emit = defineEmits( ['closed' ])

const showDialog = ref(true)

onMounted( () => {
   metadataStore.error = ""
   resetForm( {
      values: {
         requestedAt: cleanDate(metadataStore.hathiTrustStatus.requestedAt),
         packageCreatedAt: cleanDate(metadataStore.hathiTrustStatus.packageCreatedAt),
         packageSubmittedAt: cleanDate(metadataStore.hathiTrustStatus.packageSubmittedAt),
         packageStatus: metadataStore.hathiTrustStatus.packageStatus,
         metadataSubmittedAt: cleanDate(metadataStore.hathiTrustStatus.metadataSubmittedAt),
         metadataStatus: metadataStore.hathiTrustStatus.metadataStatus,
         finishedAt: cleanDate(metadataStore.hathiTrustStatus.finishedAt),
         notes: metadataStore.hathiTrustStatus.notes,
      }
   })
  
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

const submitChanges = handleSubmit( async (values) => {
   await metadataStore.updateHathiTrustStatus( values )
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