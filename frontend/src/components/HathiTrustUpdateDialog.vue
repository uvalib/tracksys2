<template>
   <DPGButton label="Batch Update HathiTrust Status" class="batch p-button-secondary" @click="showUpdateClicked" :disabled="isEnabled==false"/>
   <Dialog v-model:visible="showHathiTrustUpdate" :modal="true" header="Batch Update HathiTrust Status" position="top" >
      <div class="hathi-panel">
         <div class="columns">
            <div>
               <label>Field</label>
               <select v-model="hathiTrustField">
                  <option value="" disabled selected>Select a field</option>
                  <option value="metadata_submitted_at">Date Metadata Submitted</option>
                  <option value="metadata_status">Metadata Status</option>
                  <option value="package_submitted_at">Date Package Submitted</option>
                  <option value="package_status">Package Status</option>
                  <option value="finished_at">Date Finished</option>
               </select>
            </div>
            <div>
               <label>Value</label>
               <select v-if="hathiTrustField == 'metadata_status' || hathiTrustField == 'package_status'" v-model="hathiTrustValue">
                  <option value="" disabled selected>Select a status</option>
                  <option value="pending">Pending</option>
                  <option value="submitted">Submitted</option>
                  <option value="accepted">Accepted</option>
                  <option value="failed">Failed</option>
               </select>
               <Calendar v-else v-model="hathiTrustValue"  dateFormat="yy-mm-dd" showButtonBar/>
            </div>
         </div>
         <div class="buttons">
            <DPGButton label="Cancel" class="p-button-secondary" @click="showHathiTrustUpdate = false"/>
            <DPGButton label="Update" class="p-button-primary" @click="updateHathiTrustStatuses"/>
         </div>
      </div>
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import Calendar from 'primevue/calendar'
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

const hathiTrust = useHathiTrustStore()
const showHathiTrustUpdate = ref(false)
const hathiTrustField = ref("metadata_status")
const hathiTrustValue = ref("")

const isEnabled = computed(()  => {
   if ( props.orderID > -1) return true
   return props.ids.length > 0
})

const showUpdateClicked = (() => {
   showHathiTrustUpdate.value = true
   hathiTrustField.value = "metadata_status"
   hathiTrustValue.value = ""
})
const updateHathiTrustStatuses = ( () => {
   if (props.orderID > -1) {
      hathiTrust.batchUpdateOrder(props.orderID,  hathiTrustField.value, hathiTrustValue.value)
   } else {
      hathiTrust.batchUpdate(props.ids,  hathiTrustField.value, hathiTrustValue.value)
   }
   showHathiTrustUpdate.value = false
})
</script>

<style lang="scss" scoped>
button.batch {
   margin-right: 10px;
}
.hathi-panel {
   p {
      padding: 0;
      margin: 0 0 15px 0;
   }
   label {
      font-weight: 500;
      display: block;
      margin-bottom: 5px
   }
   select, input[type=text] {
      padding: 9px;
   }
   div.columns {
      display: flex;
      flex-flow: row nowrap;
      .p-calendar.p-component.p-inputwrapper {
         width: 100%;
      }
      div {
         flex-grow: 1;
      }
      div:last-of-type {
         margin-left: 15px;
      }
   }
   .buttons {
      text-align: right;
      margin-top: 15px;
      button {
         margin-left: 10px;
      }
   }
}
</style>