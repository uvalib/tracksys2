<template>
   <DPGButton @click="show" icon="pi pi-search" class="p-button-rounded p-button-text" />
   <Dialog v-model:visible="isOpen" :modal="true" :header="dialogTitle">
      <div class="lookup" v-if="mode=='lookup'">
         <input type="text" v-model="query"  @keydown.stop.prevent.enter="lookupMetadata" autofocus/>
         <DPGButton @click="lookupMetadata" label="Lookup" class="p-button-secondary"/>
         <DPGButton @click="createMetadata" label="Create" class="p-button-secondary" v-if="props.create"/>
      </div>
      <NewMetadataPanel v-else @canceled="metadataCreateCanceled" @created="metadataCreated" />
      <template v-if="searched">
         <div class="no-results" v-if="(target=='metadata' && metadataStore.totalSearchHits == 0) || (target=='orders' && ordersStore.totalLookupHits == 0)">
            No matching records found.
         </div>
         <div v-else class="hits">
            <div class="scroller">
               <DataTable v-if="target=='metadata'" :value="metadataStore.searchHits" ref="metadataHitsTable" dataKey="id"
                  stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
                  v-model:selection="selectedHit" selectionMode="single"
                  :lazy="false" :paginator="false" :rows="30" removableSort
               >
                  <Column field="id" header="ID" :sortable="true"/>
                  <Column field="pid" header="PID" :sortable="true"/>
                  <Column field="type" header="Type" :sortable="true"/>
                  <Column field="title" header="title" :sortable="true" >
                     <template #body="slotProps">{{truncateTitle(slotProps.data.title)}}</template>
                  </Column>
                  <Column field="callNumber" header="Call Number" :sortable="true"/>
                  <Column field="barcode" header="Barcode" :sortable="true"/>
               </DataTable>
               <DataTable v-if="target=='orders'" :value="ordersStore.lookupHits" ref="orderHitsTable" dataKey="id"
                  stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
                  v-model:selection="selectedHit" selectionMode="single"
                  :lazy="false" :paginator="false" :rows="30" removableSort
               >
                  <Column field="id" header="ID" :sortable="true"/>
                  <Column field="dateDue" header="Date Due" :sortable="true" />
                  <Column field="title" header="Title" :sortable="true"/>
                  <Column field="specialInstructions" header="Special Instructions" :sortable="true"/>
                  <Column field="lastName" header="Customer" class="nowrap" >
                  <template #body="slotProps">
                     {{slotProps.data.customer.lastName}}, {{slotProps.data.customer.firstName}}
                  </template>
               </Column>
               <Column field="agency.name" header="Agency" />
               </DataTable>
            </div>
         </div>
      </template>
      <template #footer>
         <div class="acts" v-if="mode=='lookup'">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
            <DPGButton @click="okClicked" label="Select" :disabled="selectedHit == null" />
         </div>
      </template>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import { useMetadataStore } from '@/stores/metadata'
import { useOrdersStore } from '@/stores/orders'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import NewMetadataPanel from '../order/NewMetadataPanel.vue'

const emit = defineEmits( ['selected' ])
const props = defineProps({
   target: {
      type: String,
      required: true
   },
   create: {
      type: Boolean,
      default: false
   }
})

const metadataStore = useMetadataStore()
const ordersStore = useOrdersStore()

const isOpen = ref(false)
const mode = ref("lookup")
const searched = ref(false)
const query = ref("")
const selectedHit = ref()

function createMetadata() {
   mode.value = "create"
}
function metadataCreateCanceled() {
   mode.value = "lookup"
   hide()
}
function metadataCreated() {
   selectedHit.value = metadataStore.searchHits[0]
   emit("selected", selectedHit.value.id)
   hide()
}

const dialogTitle = computed(() => {
   if (props.target == "metadata") {
      if (mode.value == "create") return "Create Metadata"
      return "Metadata Lookup"
   }
   return "Order Lookup"
})

async function lookupMetadata() {
   if (props.target == "metadata") {
      await metadataStore.lookup( query.value )
   } else {
      await ordersStore.lookup( query.value )
   }
   searched.value = true
}

function truncateTitle(t) {
   if (t.length < 75) return t
   return t.slice(0,75)+"..."
}
function okClicked() {
   emit("selected", selectedHit.value.id)
   hide()
}
function hide() {
   isOpen.value=false
}
function show() {
   isOpen.value = true
   selectedHit.value = null
   searched.value = false
   query.value = ""
   mode.value = "lookup"
}
</script>

<style lang="scss" scoped>
div.lookup {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-start;
   align-items: stretch;
   input[type=text] {
      flex-grow: 1;
      margin: 0;
   }
   button.p-button-secondary {
      font-size: 0.8em;
      display: inline-block;
      margin-left: 5px;
      overflow: unset;
   }
}
div.no-results {
   margin: 15px 0 0 0;
   font-size: 1.2em;
   font-weight: 600;
   text-align: center;
}
div.hits {
   margin-top: 15px;
   :deep(table.p-datatable-table) {
      font-size: 0.75em !important;
   }
   .scroller {
      max-height: 250px;
      overflow-y: scroll;
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 0 0 10px 0;
   button {
      margin-right: 10px;
   }
}
</style>
