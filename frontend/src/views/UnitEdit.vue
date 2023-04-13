<template>
   <h2>Unit {{route.params.id}}</h2>
   <div class="edit-form">
      <FormKit type="form" id="customer-detail" :actions="false" @submit="submitChanges">
         <div class="split">
            <Panel header="General Information">
               <FormKit label="Status" type="select" v-model="edited.status" :options="unitStatuses" required outer-class="first" />
               <FormKit label="Source URL" type="text" v-model="edited.patronSourceURL"/>
               <FormKit label="Special Instructions" type="textarea" rows="4" v-model="edited.specialInstructions"/>
               <FormKit label="Staff Notes" type="textarea" rows="3" v-model="edited.staffNotes"/>
            </Panel>
            <div class="sep"></div>
            <div class="column">
               <Panel header="Related Information">
                  <div class="split">
                     <div class="related">
                        <label>Order ID</label>
                        <div class="item">
                           <span>{{displayOrderID}}</span>
                           <LookupDialog target="orders" @selected="orderSelected" class="small-button" />
                        </div>
                     </div>
                     <div class="related">
                        <label>Metadata ID</label>
                        <div class="item">
                           <span>{{displayMetadataID}}</span>
                           <LookupDialog target="metadata" @selected="metadataSelected" :create="true"  class="small-button" />
                        </div>
                     </div>
                  </div>
               </Panel>
               <div class="sep"></div>
               <Panel header="Digitization Information">
                  <div class="dd">
                     <FormKit label="Intended Use" type="select" v-model="edited.intendedUseID" outer-class="first" :options="intendedUses"
                        placeholder="Select an intended use" required />
                  </div>

                  <div class="checkbox" v-if="unitsStore.canPublishToVirgo">
                     <input type="checkbox" v-model="edited.includeInDL"/>
                     <span class="label">Include in Virgo</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="edited.ocrMasterFiles"/>
                     <span class="label">OCR Master Files</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="edited.removeWatermark"/>
                     <span class="label">Remove Watermark</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="edited.completeScan"/>
                     <span class="label">Complete Scan</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="edited.throwAway"/>
                     <span class="label">
                        Throw Away
                        <span class="note">
                          (Throw away scans will not be sent to preservation. They are one-time scans made for a single patron.)
                        </span>
                     </span>
                  </div>
               </Panel>
            </div>
         </div>
         <div class="acts">
            <DPGButton label="Cancel" class="p-button-secondary" @click="cancelEdit()"/>
            <FormKit type="submit" label="Save" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import Panel from 'primevue/panel'
import LookupDialog from '@/components/LookupDialog.vue'

const route = useRoute()
const router = useRouter()
const unitsStore = useUnitsStore()
const systemStore = useSystemStore()

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
const displayMetadataID = computed( () => {
   if (edited.value.metadataID && edited.value.metadataID) {
      return edited.value.metadataID
   }
   return "None"
})
const displayOrderID = computed( () => {
   if (edited.value.orderID && edited.value.orderID) {
      return edited.value.orderID
   }
   return "None"
})

const edited = ref({
   status: "",
   patronSourceURL: "",
   specialInstructions: "",
   staffNotes: "",
   completeScan: false,
   throwAway: false,
   orderID: "",
   metadataID: "",
   intendedUseID: 0,
   ocrMasterFiles: false,
   removeWatermark: false,
   includeInDL: false,
})

onMounted( async () =>{
   let unitID = route.params.id
   await unitsStore.getDetails(unitID)
   document.title = `Edit | Unit ${unitID}`

   edited.value.status = unitsStore.detail.status
   edited.value.patronSourceURL = unitsStore.detail.patronSourceURL
   edited.value.staffNotes = unitsStore.detail.staffNotes
   edited.value.specialInstructions = unitsStore.detail.specialInstructions
   edited.value.completeScan = unitsStore.detail.completeScan
   edited.value.throwAway = unitsStore.detail.throwAway
   edited.value.metadataID = unitsStore.detail.metadataID
   edited.value.orderID = unitsStore.detail.orderID
   if (unitsStore.detail.intendedUse) {
      edited.value.intendedUseID =  unitsStore.detail.intendedUse.id
   }
   edited.value.ocrMasterFiles = unitsStore.detail.ocrMasterFiles
   edited.value.removeWatermark = unitsStore.detail.removeWatermark
   edited.value.includeInDL = unitsStore.detail.includeInDL
})

function orderSelected( o ) {
   edited.value.orderID = o
}
function metadataSelected( o ) {
   edited.value.metadataID = o
}
function cancelEdit() {
   router.push(`/units/${route.params.id}`)
}

async function submitChanges() {
   await unitsStore.submitEdit( edited.value )
   if (systemStore.showError == false) {
      router.push(`/units/${unitsStore.detail.id}`)
   }
}
</script>


<style lang="scss" scoped>
.edit-form {
   width: 80%;
   margin: 30px auto 0 auto;

   .note {
      text-align: left;
      font-size: 0.8em;
      font-style: italic;
      padding:0;
      margin: 0 0 0 5px;
      display: inline-block
   }
   .dd {
     margin-bottom: 20px;
   }
   .checkbox {
      padding-bottom: 10px;
   }

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      :deep(.p-panel), .column {
         flex-grow: 1;
      }
      .sep {
         display: inline-block;
         width: 20px;
      }
      .related {
         label {
            display: block;
            text-align: left;
         }
         .item {
            text-align: left;
            display: flex;
            flex-flow: row nowrap;
            justify-content: flex-start;
            align-items: center;
            margin-top: 10px;
            :deep(button.small-button) {
               padding: 3px 15px;
               font-size: 0.85em;
               margin-left: 10px;
            }
         }
      }
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 25px 0;
   button {
      margin-right: 10px;
   }
}
</style>
