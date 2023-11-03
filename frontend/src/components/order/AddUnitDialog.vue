<template>
   <DPGButton @click="show" :label="props.label" class="p-button-secondary add"  :class="{ small: size=='small'}"/>
   <Dialog v-model:visible="isOpen" :modal="true" :header="title" :style="{width: '750px'}" position="top">
      <FormKit v-if="mode=='unit'" type="form" id="customer-detail" :actions="false" @submit="createUnit">
         <Panel header="Unit Metadata" class="margin-bottom">
            <div class="lookup">
               <input type="text" v-model="metadataSearch"  @keydown.stop.prevent.enter="lookupMetadata"/>
               <DPGButton @click="lookupMetadata" label="Lookup" class="p-button-secondary"/>
               <DPGButton @click="createMetadata" label="Create" class="p-button-secondary"/>
            </div>
            <template v-if="searched">
               <div class="hits">
                  <div class="scroller">
                     <DataTable :value="metadataStore.searchHits" ref="metadataHitsTable" dataKey="id"
                        stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
                        v-model:selection="selectedMetadata" selectionMode="single"
                        :lazy="false" :paginator="false" :rows="30" removableSort>
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
            <FormKit label="Intended Use" type="select" v-model="unitInfo.intendedUseID" outer-class="first" :options="intendedUses" placeholder="Select an intended use" required/>
            <FormKit label="Source URL" type="text" v-model="unitInfo.sourceURL"/>
            <FormKit label="Special Instructions" type="textarea" rows="4" v-model="unitInfo.specialInstructions"/>
            <FormKit label="Staff Notes" type="textarea" rows="2" v-model="unitInfo.staffNotes"/>
            <div class="opts">
               <div class="checkbox">
                  <input type="checkbox" v-model="unitInfo.completeScan"/>
                  <span class="label">Complete Scan</span>
               </div>
               <div class="checkbox">
                  <input type="checkbox" v-model="unitInfo.throwAway"/>
                  <span class="label">Throw Away</span>
               </div>
               <div class="checkbox">
                  <input type="checkbox" v-model="unitInfo.includeInDL"/>
                  <span class="label">Include in Virgo</span>
               </div>
            </div>
         </Panel>
         <p class="error">{{error}}</p>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
            <FormKit type="submit" label="Add Unit" wrapper-class="submit-button" />
         </div>
      </FormKit>
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
import NewMetadataPanel from '@/components/NewMetadataPanel.vue'

const props = defineProps({
   label: {
      type: String,
      default: "Add Unit",
   },
   size: {
      type: String,
      default: "normal",
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
const error = ref("")
const metadataSearch = ref("")
const selectedMetadata = ref()
const unitInfo = ref({
   intendedUseID: 0,
   sourceURL: "",
   specialInstructions: "",
   staffNotes: "",
   completeScan: false,
   throwAway: false,
   includeInDL: false,
})

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

function createMetadata() {
   mode.value = "metadata"
}
function metadataCreateCanceled() {
   mode.value = "unit"
}
function metadataCreated() {
   metadataSearch.value = metadataStore.detail.pid
   searched.value = true
   selectedMetadata.value = metadataStore.searchHits[0]
   mode.value = "unit"
}

async function lookupMetadata() {
   await metadataStore.lookup( metadataSearch.value )
   searched.value = true
}

async function createUnit() {
   error.value = ""
   if ( !selectedMetadata.value) {
      error.value = "Please select a metadata record for the new unit."
   } else {
      if (props.item != null ) {
         await ordersStore.addUnit(selectedMetadata.value.id, unitInfo.value, props.item.id)
      } else {
         await ordersStore.addUnit(selectedMetadata.value.id, unitInfo.value)
      }

      if (systemStore.error == "") {
         hide()
      }
   }
}

function hide() {
   isOpen.value=false
}
function show() {
   mode.value = "unit"
   isOpen.value = true
   error.value = ""
   metadataSearch.value = ""
   selectedMetadata.value = null
   metadataStore.resetSearch()
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
      unitInfo.value.specialInstructions = si
      unitInfo.value.sourceURL = props.item.sourceURL
      lookupMetadata()
   }
}
</script>

<style lang="scss" scoped>
button.p-button.add.small {
   font-size: 0.7em;
   padding: 5px 10px;
}
div.p-panel {
   font-size: 0.85em;
   div.lookup {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-start;

      button {
         font-size: 0.9em;
         display: inline-block;
         margin-left: 5px;
         padding: 0.5em 1em;
         width: 75px;
      }
   }
   div.hint {
      margin-top: 15px;
   }
   div.no-results {
      margin: 15px 0 0 0;
      font-size: 1.2em;
      font-weight: 600;
      text-align: center;
   }
   :deep(div.hits) {
      margin-top: 15px;
      table.p-datatable-table {
         font-size: 0.75em !important;
      }
      .scroller {
         max-height: 250px;
         overflow-y: scroll;
      }
   }
}
div.margin-bottom {
   margin-bottom: 25px;
}
div.opts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   margin: 20px 0 10px 0;
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

.error {
   padding: 0;
   margin: 15px 0 0 0;
   text-align: center;
   color: var(--uvalib-red-emergency);
}
</style>
