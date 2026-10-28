<template>
   <h2>Unit {{route.params.id}}</h2>
   <UForm :schema="schema" :state="state" class="flex flex-col gap-4 w-full max-w-300 mx-auto! mt-10!" @submit="submitChanges">
      <div class="row-between gap-4">
         <UCard title="General Information" class="flex-1">
            <div class="column gap-4">
               <UFormField name="status" label="Status" class="grow">   
                  <USelect v-model="state.status" :items="unitStatuses" class="w-full"/>   
               </UFormField>
               <UFormField name="patronSourceURL" label="Source URL" class="grow">   
                  <UInput v-model="state.patronSourceURL" class="w-full"/>   
               </UFormField>
               <UFormField name="specialInstructions" label="Special Instructions" class="grow">   
                  <UTextarea v-model="state.specialInstructions" class="w-full" :rows="4"/>   
               </UFormField>
               <UFormField name="staffNotes" label="Staff Notes" class="grow">   
                  <UTextarea v-model="state.staffNotes" class="w-full" :rows="4"/>   
               </UFormField>
            </div>
         </UCard>
         <div class="column gap-4 flex-1">
            <UCard title="Related Information">
               <div class="row-between">
                  <div class="row-left gap-2 items-end">
                     <UFormField name="orderID" label="Order ID" :required="true">   
                        <UInput v-model="state.orderID" :readonly="true"/>   
                     </UFormField>
                     <LookupDialog target="orders" @selected="orderSelected" />
                  </div>
                  <div class="row-left gap-2 items-end">
                     <UFormField name="metadataID" label="Metadata ID" :required="true">   
                        <UInput v-model="state.metadataID" :readonly="true"/>   
                     </UFormField>  
                     <LookupDialog target="metadata" @selected="metadataSelected" :create="true"  /> 
                  </div>
               </div>
            </UCard>
            <UCard title="Digitization Information">
               <div class="column gap-4">
                  <UFormField name="intendedUseID" label="Intended Use" class="grow" :required="true">   
                     <USelect v-model="state.intendedUseID" :items="intendedUses" class="w-full"/>   
                  </UFormField>  
                  <UCheckbox label="Include in Virgo" v-model="state.includeInDL" class="w-full" v-if="unitsStore.canPublishToVirgo"/>  
                  <UCheckbox label="OCR Master Files" v-model="state.ocrMasterFiles" class="w-full"/>  
                  <UCheckbox label="Remove Watermark" v-model="state.removeWatermark" class="w-full"/>  
                  <UCheckbox label="Complete Scan" v-model="state.completeScan" class="w-full"/>  
                  <UCheckbox label="Throw Away" v-model="state.throwAway" class="w-full"/> 
                  <div class="pl-4 text-sm italic">Throw away scans will not be sent to preservation. They are one-time scans made for a single patron.</div>
               </div> 
            </UCard>
         </div>
      </div>
      <div class="row-right gap-2">
         <UButton label="Cancel" color="secondary" @click="cancelEdit()"/>
         <UButton label="Save" type="submit"/>
      </div>
   </UForm>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import LookupDialog from '@/components/LookupDialog.vue'
import { useConfirm } from "@/composables/useConfirm"
import * as yup from 'yup'

const route = useRoute()
const router = useRouter()
const unitsStore = useUnitsStore()
const systemStore = useSystemStore()

const schema = yup.object({
   orderID: yup.number().min(1, 'OrderID is required'),
   metadataID: yup.number().min(1, 'MetadataID is required'),
   intendedUseID: yup.number().min(1,'Intended use is required'),
})

const state = ref({
   metadataID: null,
   orderID: null,
   status: "",
   patronSourceURL: "",
   specialInstructions: "",
   staffNotes: "",
   intendedUseID: null,
   ocrMasterFiles: false,
   removeWatermark: false,
   includeInDL: false,
   throwAway: false,
   completeScan: false
})

const unitStatuses = computed(() => {
   let out = []
   out.push( {label: "Approved", value: "approved"} )
   out.push( {label: "Canceled", value: "canceled"} )
   out.push( {label: "Done", value: "done"} )
   out.push( {label: "Error", value: "error"} )
   out.push( {label: "Unapproved", value: "unapproved"} )
   return out
})

const intendedUses = computed(() => {
   let out = []
   systemStore.intendedUses.forEach( a => {
      if (a.name == "Digital Collection Building") {
         out.push( {label: `${a.name}: Highest Possible resolution TIFF`, value: a.id} )
      } else if (a.deliverableFormat == 'pdf') {
         out.push( {label: `${a.name}: PDF`, value: a.id} )
      } else {
         let dpi = "DPI"
         if (a.deliverableResolution == "Highest Possible") {
            dpi = "resolution"
         }
         out.push( {label: `${a.name}: ${a.deliverableResolution} ${dpi} ${a.deliverableFormat.toUpperCase()}`, value: a.id} )
      }
   })
   return out
})

onMounted( async () =>{
   let unitID = route.params.id
   await unitsStore.getDetails(unitID)
   document.title = `Edit | Unit ${unitID}`

   state.value = {
      status: unitsStore.detail.status,
      patronSourceURL: unitsStore.detail.patronSourceURL,
      staffNotes: unitsStore.detail.staffNotes,
      specialInstructions: unitsStore.detail.specialInstructions,
      completeScan: unitsStore.detail.completeScan,
      throwAway: unitsStore.detail.throwAway,
      metadataID: unitsStore.detail.metadataID,
      orderID: unitsStore.detail.orderID,
      ocrMasterFiles: unitsStore.detail.ocrMasterFiles,
      removeWatermark: unitsStore.detail.removeWatermark,
      includeInDL: unitsStore.detail.includeInDL,
      intendedUseID: null,
   }
   if (unitsStore.detail.intendedUse) {
     state.value.intendedUseID =  unitsStore.detail.intendedUse.id
   }
   if (unitsStore.canPublishToVirgo == false ) {
      state.value.includeInDL = false
   }
})

const orderSelected = (( newOrderID ) => {
   state.value.orderID = newOrderID
})

const metadataSelected = (( newMetadataID ) => {
   state.value.metadataID = newMetadataID
})

const cancelEdit = (() => {
   router.push(`/units/${route.params.id}`)
})

const submitChanges = ( async () => {
   if ( state.value.metadataID != unitsStore.detail.metadataID) {
      const msg = "You have changed the metadata record for this unit. All master files will also be updated to use this metadata. Are you sure?",
      resp = await useConfirm("Confirm Metadata Change", msg, "Update")
      if (resp) {
         await unitsStore.submitEdit( state.value )
         if (systemStore.showError == false) {
            router.push(`/units/${unitsStore.detail.id}`)
         }
      }
   } else {
      await unitsStore.submitEdit( state.value )
      if (systemStore.showError == false) {
         router.push(`/units/${unitsStore.detail.id}`)
      }
   }
})
</script>

<style lang="scss" scoped>
</style>
