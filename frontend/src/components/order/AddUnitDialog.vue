<template>
   <DPGButton @click="show" label="Add Unit" class="p-button-secondary add"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Unit" :style="{width: '750px'}">
      <FormKit type="form" id="customer-detail" :actions="false" @submit="createUnit">
         <Panel header="Unit Metadata" class="margin-bottom">
            <div class="lookup">
               <input type="text" v-model="metadataSearch"  @keydown.stop.prevent.enter="lookupMetadata"/>
               <DPGButton @click="lookupMetadata" label="Lookup" class="p-button-secondary"/>
            </div>
            <template v-if="metadataStore.searched">
               <div class="no-results" v-if="metadataStore.totalSearchHits == 0">
                  No matching metadata records found.
               </div>
               <div v-else class="hits">
                  <div class="scroller">
                     <DataTable :value="metadataStore.searchHits" ref="metadataHitsTable" dataKey="id"
                        stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
                        v-model:selection="selectedMetadata" selectionMode="single"
                        :lazy="false" :paginator="false" :rows="30" removableSort>
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
               Find a metadata record for the new unit.
            </div>
         </Panel>
         <Panel header="Digitization Information">
            <FormKit label="Intended Use" type="select" v-model="unitInfo.intendedUseID" outer-class="first" :options="intendedUses"/>
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
                  <span class="label">Include in DL</span>
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

const systemStore = useSystemStore()
const metadataStore = useMetadataStore()
const ordersStore = useOrdersStore()

const isOpen = ref(false)
const error = ref("")
const metadataSearch = ref("")
const selectedMetadata = ref()
const unitInfo = ref({
   intendedUseID: null,
   sourceURL: "",
   specialInstructions: "",
   staffNotes: "",
   completeScan: false,
   throwAway: false,
   includeInDL: false,
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

async function lookupMetadata() {
   await metadataStore.lookup( metadataSearch.value )
}

async function createUnit() {
   error.value = ""
   console.log(selectedMetadata.value)
   if ( !selectedMetadata.value) {

      error.value = "Please select a metadata record for the new unit."
   } else {
      await ordersStore.addUnit(selectedMetadata.value.id, unitInfo.value)
      hide()
   }
}

function hide() {
   isOpen.value=false
}
function show() {
   isOpen.value = true
   error.value = ""
   metadataSearch.value = ""
   selectedMetadata.value = null
   metadataStore.resetSearch()
}
</script>

<style lang="scss" scoped>
button.p-button.add {
   font-size: 0.8em;
   padding: 5px 20px;
}
div.p-panel {
   font-size: 0.85em;
   div.lookup {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: flex-start;
      input[type-text] {
         flex-grow: 1;
         margin-right: 5px;
      }
      button {
         font-size: 0.9em;
         display: inline-block;
         margin-left: 5px;
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
:deep(div.formkit-outer.first) {
   .dpg-form-label {
      margin-top: 5px;
   }
}
div.opts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   margin: 20px 0 10px 0;
   div.checkbox {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: center;
      margin: 0;
      input[type=checkbox] {
         width: 18px;
         height: 18px;
         margin-right: 10px;
         display: inline-block;
      }
      span {
         display: inline-block;
      }
   }
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
