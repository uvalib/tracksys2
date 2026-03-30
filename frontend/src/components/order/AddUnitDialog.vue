<template>
   <DPGButton @click="show" :label="props.label" severity="secondary" />
   <Dialog v-model:visible="isOpen" :modal="true" :header="title" :style="{width: '750px'}" position="top" >
      <form v-if="mode=='unit'" @submit="createUnit">
         <Panel header="Unit Metadata">
            <div class="lookup">
               <InputText v-model="metadataSearch"  @keydown.stop.prevent.enter="lookupMetadata" fluid/>
               <DPGButton @click="lookupMetadata" label="Lookup" severity="secondary" :loading="metadataStore.working" loadingIcon="pi pi-spin pi-spinner"/>
               <DPGButton @click="createMetadata" label="Create" severity="secondary"/>
            </div>
            <Message v-if="errors.metadataID" severity="error" size="small" variant="simple">{{ errors.metadataID }}</Message>
            <template v-if="searched">
               <div class="hits">
                  <div class="scroller">
                     <DataTable :value="metadataStore.searchHits" ref="metadataHitsTable" dataKey="id"
                        stripedRows showGridlines size="small"
                        @update:selection="metadataSelected" v-model:selection="selectedMetadata" selectionMode="single"
                        :lazy="false" :paginator="false" :rows="30" removableSort
                     >
                        <template #empty>No matching metadata records</template>
                        <Column field="pid" header="PID" :sortable="true"/>
                        <Column field="type" header="Type" :sortable="true"/>
                        <Column field="title" header="title" :sortable="true"/>
                        <Column field="callNumber" header="Call Number" :sortable="true"/>
                        <Column field="barcode" header="Barcode" :sortable="true"/>
                     </DataTable>
                  </div>
               </div>
            </template>
            <div v-else class="hint">
               Find or create a metadata record for the new unit.
            </div>
         </Panel>
         <Panel header="Digitization Information">
            <FormField id="intendeduse" label="Intended Use" :error="errors.intendedUseID" :required="true">
               <Select id="intendeduse" v-model="intendedUseID"  :options="intendedUses" optionLabel="label" optionValue="value" placeholder="Select an intended use" />   
            </FormField>
            <FormField id="srcurl" label="Source URL">
               <InputText id="srcurl" v-model="sourceURL" type="text" />   
            </FormField>
            <FormField id="specialinst" label="Special Instructions">
               <Textarea id="specialinst" v-model="specialInstructions" rows="3"/>   
            </FormField>
            <FormField id="staffnotes" label="Staff Notes">
               <Textarea id="staffnotes" v-model="staffNotes" rows="3"/>   
            </FormField>
            <div class="opts">
               <div class="checkbox">
                  <input type="checkbox" v-model="completeScan"/>
                  <span class="label">Complete Scan</span>
               </div>
               <div class="checkbox">
                  <input type="checkbox" v-model="throwAway"/>
                  <span class="label">Throw Away</span>
               </div>
               <div class="checkbox">
                  <input type="checkbox" v-model="includeInDL"/>
                  <span class="label">Include in Virgo</span>
               </div>
            </div>
         </Panel>
         <div class="acts">
            <DPGButton @click="isOpen=false" label="Cancel" severity="secondary"/>
            <DPGButton label="Add Unit" type="submit" />
         </div>
      </form>
      <NewMetadataPanel v-else @canceled="metadataCreateCanceled" @created="metadataCreated" />
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import Panel from 'primevue/panel'
import { useSystemStore } from '@/stores/system'
import { useMetadataStore } from '@/stores/metadata'
import { useOrdersStore } from '@/stores/orders'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import NewMetadataPanel from '@/components/metadata/NewMetadataPanel.vue'

import { useForm } from 'vee-validate'
import * as yup from 'yup'
import FormField from '@/components/FormField.vue'
import InputText from 'primevue/inputtext'
import Select from 'primevue/select'
import Textarea from 'primevue/textarea'
import Message from 'primevue/message'

const props = defineProps({
   label: {
      type: String,
      default: "Add Unit",
   },
   item: {
      type: Object,
      default: null,
   },
})

const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const ordersStore = useOrdersStore()

const isOpen = ref(false)
const mode = ref("unit")
const searched = ref(false)
const metadataSearch = ref("")
const selectedMetadata = ref()

const { errors, handleSubmit, defineField, resetForm, setValues} = useForm({
  validationSchema: yup.object().shape({
      metadataID: yup.number().min(1, "Metadata ID is required"),
      intendedUseID: yup.number().min(1,'Intended use is required'),
   })
})

const [intendedUseID] = defineField('intendedUseID')
const [sourceURL] = defineField('sourceURL')
const [specialInstructions] = defineField('specialInstructions')
const [staffNotes] = defineField('staffNotes')
const [completeScan] = defineField('completeScan')
const [throwAway] = defineField('throwAway')
const [includeInDL] = defineField('includeInDL')

const title = computed(() => {
   if (mode.value == "unit") return "Add Unit"
   return "Create Metadata"
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

const metadataSelected = ( (e) => {
   if (e) {
      setValues({metadataID: e.id})
   }
})

const createMetadata = (() => {
   mode.value = "metadata"
})

const metadataCreateCanceled = (() => {
   mode.value = "unit"
})

const metadataCreated = (() => {
   console.log(metadataStore.searchHits[0])
   metadataSearch.value = metadataStore.detail.pid
   searched.value = true
   selectedMetadata.value = metadataStore.searchHits[0]
   setValues({metadataID: metadataStore.searchHits[0].id})
   mode.value = "unit"
})

const lookupMetadata = ( async () => {
   await metadataStore.lookup( metadataSearch.value )
   searched.value = true
})

const createUnit = handleSubmit( async (values) => {
   await ordersStore.addUnit(values)
   if (systemStore.error == "") {
      isOpen.value = false
   }
})

const show = (() => {
   mode.value = "unit"
   isOpen.value = true
   metadataSearch.value = ""
   selectedMetadata.value = null
   metadataStore.resetSearch()

   let val = {
      itemID: 0,
      metadataID: 0,
      intendedUseID: 0,
      sourceURL: "",
      specialInstructions: "",
      staffNotes: "",
      completeScan: false,
      throwAway: false,
      includeInDL: false,
   }

   if (props.item) {
      metadataSearch.value = props.item.title
      var si = `Title: ${props.item.title}`
      si += `\nPages to Digitize: ${props.item.pages}`
      if (props.item.callNumber && props.item.callNumber != "") {
         si += `\nCall Number: ${props.item.callNumber}`
         metadataSearch.value = props.item.callNumber
      }
      if (props.item.author && props.item.author != "") {
         si += `\nAuthor: ${props.item.author}`
      }
      if (props.item.year && props.item.year != "") {
         si += `\nYear: ${props.item.year}`
      }
      if (props.item.location && props.item.location != "") {
         si += `\nLocation: ${props.item.location}`
      }
      if (props.item.description && props.item.description != "") {
         si += `\nDescription: ${props.item.description}`
      }
      val.itemID = props.item.id
      val.specialInstructions = si
      val.sourceURL = props.item.sourceURL
      val.intendedUseID = props.item.intendedUse.id
      lookupMetadata()
   }
    resetForm({values: val })
})

</script>

<style lang="scss" scoped>
div.p-panel {
   div.lookup {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: stretch;
      gap: 10px;
      margin-top: 20px;
   }
   div.no-results {
      margin: 15px 0 0 0;
      font-size: 1.2em;
      font-weight: 600;
      text-align: center;
   }
   div.hits {
      .scroller {
         max-height: 250px;
         overflow-y: scroll;
      }
   }
}
:deep(.p-panel-content) {
   display: flex;
   flex-direction: column;
   gap: 15px;
}
</style>
