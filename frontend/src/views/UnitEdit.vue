<template>
   <h2>Unit {{route.params.id}}</h2>
   <div class="edit-form">
      <form @submit="submitChanges">
         <div class="split">
            <Panel header="General Information">
               <FormField id="status" label="Status">
                  <Select id="status" fluid v-model="status"  :options="unitStatuses" optionLabel="label" optionValue="value" placeholder="Select a status" />   
               </FormField>
               <FormField id="srcurl" label="Source URL">
                  <InputText id="srcurl" v-model="patronSourceURL" type="text" />   
               </FormField>
               <FormField id="specialinst" label="Special Instructions">
                  <Textarea id="specialinst" v-model="specialInstructions" rows="3" fluid/>   
               </FormField>
               <FormField id="staffnotes" label="Staff Notes">
                  <Textarea id="staffnotes" v-model="staffNotes" rows="3" fluid/>   
               </FormField>
            </Panel>
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
                  <Message v-if="errors.metadataID" severity="error" size="small" variant="simple">{{ errors.metadataID}}</Message>
                  <Message v-if="errors.orderID" severity="error" size="small" variant="simple">{{ errors.orderID }}</Message>
               </Panel>
               <Panel header="Digitization Information">
                  <FormField id="intendeduse" label="Intended Use" :error="errors.intendedUseID" :required="true">
                     <Select id="intendeduse" v-model="intendedUseID"  :options="intendedUses" optionLabel="label" optionValue="value" placeholder="Select an intended use" />   
                  </FormField>
                  <div class="checkbox" v-if="unitsStore.canPublishToVirgo">
                     <input type="checkbox" v-model="includeInDL"/>
                     <span class="label">Include in Virgo</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="ocrMasterFiles"/>
                     <span class="label">OCR Master Files</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="removeWatermark"/>
                     <span class="label">Remove Watermark</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="completeScan"/>
                     <span class="label">Complete Scan</span>
                  </div>
                  <div class="checkbox">
                     <input type="checkbox" v-model="throwAway"/>
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
            <DPGButton label="Cancel" severity="secondary" @click="cancelEdit()"/>
            <DPGButton label="Save" type="submit"/>
         </div>
      </form>
   </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { onMounted, ref, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useSystemStore } from '@/stores/system'
import Panel from 'primevue/panel'
import Select from 'primevue/select'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'
import LookupDialog from '@/components/LookupDialog.vue'
import { useConfirm } from "primevue/useconfirm"
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'

const confirm = useConfirm()
const route = useRoute()
const router = useRouter()
const unitsStore = useUnitsStore()
const systemStore = useSystemStore()

const { values, errors, resetForm, handleSubmit, defineField, setValues } = useForm({
  validationSchema: yup.object().shape({
      orderID: yup.number().min(1, 'OrderID is required'),
      metadataID: yup.number().min(1, 'MetadataID is required'),
      intendedUseID: yup.number().min(1,'Intended use is required'),
   })
})

const [status] = defineField('status')
const [patronSourceURL] = defineField('patronSourceURL')
const [specialInstructions] = defineField('specialInstructions')
const [staffNotes] = defineField('staffNotes')
const [intendedUseID] = defineField('intendedUseID')
const [ocrMasterFiles] = defineField('ocrMasterFiles')
const [removeWatermark] = defineField('removeWatermark')
const [includeInDL] = defineField('includeInDL')
const [throwAway] = defineField('throwAway')
const [completeScan] = defineField('completeScan')

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
   if (values.metadataID) {
      return values.metadataID
   }
   return "None"
})

const displayOrderID = computed( () => {
   if (values.orderID) {
      return values.orderID
   }
   return "None"
})

onMounted( async () =>{
   let unitID = route.params.id
   await unitsStore.getDetails(unitID)
   document.title = `Edit | Unit ${unitID}`

   let vals = {
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
      intendedUseID: 0,
   }
   if (unitsStore.detail.intendedUse) {
     vals.intendedUseID =  unitsStore.detail.intendedUse.id
   }
   resetForm({values: vals})
})

const orderSelected = (( newOrderID ) => {
   setValues({ orderID: newOrderID})
})

const metadataSelected = (( newMetadataID ) => {
   setValues({ metadataID: newMetadataID} )
})

const cancelEdit = (() => {
   router.push(`/units/${route.params.id}`)
})

const submitChanges = handleSubmit( async (values) => {
   if ( values.metadataID != unitsStore.detail.metadataID) {
      confirm.require({
         message: "You have changed the metadata record for this unit. All master files will also be updated to use this metadata. Are you sure?",
         header: 'Confirm Metadata Change',
         icon: 'pi pi-question-circle',
         rejectProps: {
            label: 'Cancel',
            severity: 'secondary'
         },
         acceptProps: {
            label: 'Update'
         },
         accept: async () => {
            await unitsStore.submitEdit( values )
            if (systemStore.showError == false) {
               router.push(`/units/${unitsStore.detail.id}`)
            }
         },
      })
   } else {
      await unitsStore.submitEdit( values )
      if (systemStore.showError == false) {
         router.push(`/units/${unitsStore.detail.id}`)
      }
   }
})
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

   .column {
      display: flex;
      flex-direction: column;
      gap: 15px;
   }

   .split {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      gap: 20px;
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
            gap: 10px;
         }
      }
   }
}
:deep(.p-panel) {
   flex-grow: 1;
   .p-panel-content {
      display: flex;
      flex-direction: column;
      gap: 15px;
      text-align: left;
   }
}
</style>
