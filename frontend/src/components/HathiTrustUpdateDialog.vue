<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Batch Update HathiTrust Status">
      <UButton @click="showUpdateClicked" label="Batch Update HathiTrust Status" color="secondary" :disabled="isEnabled==false"/>
      <template  #body>
         <div class="columns">
            <div class="item">
               <label>Field</label>
               <USelect v-model="hathiTrustField" :items="fields" placeholder="Select a field" 
                   @change="fieldChanged" class="w-full"/>
            </div>
            <div>
               <label>Value</label>
               <USelect v-if="hathiTrustField == 'metadata_status' || hathiTrustField == 'package_status'" 
                  v-model="hathiTrustValue" :items="statuses" placeholder="Select a status" class="w-full"
               />
               <UInputDate v-else v-model="hathiTrustValue" class="w-full"  />
            </div>
         </div>
      </template>
      <template #footer="{close}">
         <UButton label="Cancel" color="secondary" @click="close"/>
         <UButton label="Update" @click="updateHathiTrustStatuses"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref,computed } from 'vue'
import { useHathiTrustStore } from '@/stores/hathitrust'

const props = defineProps({
   orderID: {
      type: Number,
      default: -1
   },
   ids: {
      type: Array,
      default: null
   },
})

const fields = [
   {label: "Date Metadata Submitted", value: "metadata_submitted_at"},
   {label: "Metadata Status", value: "metadata_status"},
   {label: "Date Package Submitted", value: "package_submitted_at"},
   {label: "Package Status", value: "package_status"},
   {label: "Date Finished", value: "finished_at"},
]

const statuses = [
   {label: "Pending", value: "pending"},
   {label: "Submitted", value: "submitted"},
   {label: "Accepted", value: "accepted"},
   {label: "Failed", value: "failed"},
]

const hathiTrust = useHathiTrustStore()
const isOpen = ref(false)
const hathiTrustField = ref("metadata_status")
const hathiTrustValue = ref()

const fieldChanged = (() => {
   console.log("CHANGE FIELD")
   hathiTrustValue.value = null
})
const isEnabled = computed(()  => {
   if ( props.orderID > -1) return true
   return props.ids.length > 0
})

const showUpdateClicked = (() => {
   isOpen.value = true
   hathiTrustField.value = "metadata_status"
   hathiTrustValue.value = ""
})
const updateHathiTrustStatuses = ( () => {
   if (props.orderID > -1) {
      hathiTrust.batchUpdateOrder(props.orderID,  hathiTrustField.value, hathiTrustValue.value.toString())
   } else {
      hathiTrust.batchUpdate(props.ids,  hathiTrustField.value, hathiTrustValue.value.toString())
   }
   isOpen.value = false
})
</script>

<style lang="scss" scoped>
div.columns {
   display: flex;
   flex-flow: column;
   justify-content: flex-start;
   align-items: stretch;
   gap: 20px;
   label {
      display: block;
   }

   .item {
      display: flex;
      flex-flow: column;
      justify-content: flex-start;
      align-items: stretch;
      gap: 5px;
      .p-datepicker {
         width: 100%;
      }
   }
}
</style>