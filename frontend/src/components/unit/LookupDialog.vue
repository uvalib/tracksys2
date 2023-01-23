<template>
   <DPGButton @click="show" class="p-button-secondary lookup-button" :label="props.label" :disabled="props.disabled" :class="props.class" />
   <Dialog v-model:visible="isOpen" :modal="true" :header="dialogTitle">
      <div class="lookup" v-if="mode=='lookup'">
         <input type="text" v-model="query"  @keydown.stop.prevent.enter="lookupRecords" autofocus :placeholder="searchPlaceholder"/>
         <DPGButton @click="lookupRecords" label="Lookup" class="p-button-secondary"/>
         <DPGButton @click="createMetadata" label="Create" class="p-button-secondary" v-if="props.create"/>
      </div>
      <NewMetadataPanel v-else @canceled="metadataCreateCanceled" @created="metadataCreated" />
      <template v-if="searched && mode=='lookup'">
         <div class="no-results"
            v-if="(target=='metadata' && metadataStore.totalSearchHits == 0) || (target=='orders' && ordersStore.totalLookupHits == 0) || (target=='component' && componentsStore.totalLookupHits == 0)">
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
                  <Column field="title" header="Title" :sortable="true" >
                     <template #body="slotProps">{{truncateTitle(slotProps.data.title)}}</template>
                  </Column>
                  <Column field="callNumber" header="Call Number" :sortable="true"/>
                  <Column field="barcode" header="Barcode" :sortable="true"/>
               </DataTable>

               <DataTable v-if="target=='component'" :value="componentsStore.searchHits" ref="componentHitsTable" dataKey="id"
                  stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
                  v-model:selection="selectedHit" selectionMode="single"
                  :lazy="false" :paginator="false" :rows="30" removableSort
               >
                  <Column field="id" header="ID" :sortable="true"/>
                  <Column field="pid" header="PID" :sortable="true"/>
                  <Column field="title" header="Title" :sortable="true" >
                     <template #body="slotProps">{{truncateTitle(slotProps.data.title)}}</template>
                  </Column>
                  <Column field="label" header="Label" :sortable="true" >
                     <template #body="slotProps">{{truncateTitle(slotProps.data.label)}}</template>
                  </Column>
                  <Column field="date" header="Date" :sortable="false"/>
                  <Column field="eadID" header="EAD ID" :sortable="false"/>
                  <Column field="masterFileCount" header="Master Files"/>
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
            <DPGButton @click="okClicked" label="Select" :disabled="!selectedHit" />
         </div>
      </template>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import { useMetadataStore } from '@/stores/metadata'
import { useOrdersStore } from '@/stores/orders'
import { useComponentsStore } from '@/stores/components'
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
   },
   label: {
      type: String,
      default: "Lookup"
   },
   disabled: {
      type: Boolean,
      default: false
   },
   class: {
      type: String,
      default: "wombat"
   }
})

const metadataStore = useMetadataStore()
const ordersStore = useOrdersStore()
const componentsStore = useComponentsStore()

const isOpen = ref(false)
const mode = ref("lookup")
const searched = ref(false)
const query = ref("")
const selectedHit = ref()

function createMetadata() {
   mode.value = "create"
   searched.value = false
   query.value = ""
   selectedHit.value = null
}
function metadataCreateCanceled() {
   mode.value = "lookup"
}
function metadataCreated() {
   selectedHit.value = metadataStore.searchHits[0]
   emit("selected", selectedHit.value.id)
   hide()
}

const searchPlaceholder = computed(() => {
   if (props.target == "metadata") {
      return "Lookup metadata..."
   }
   if (props.target == "component") {
      return "Lookup component..."
   }
   return "Lookup order..."
})

const dialogTitle = computed(() => {
   if (props.target == "metadata") {
      if (mode.value == "create") return "Create Metadata"
      return "Metadata Lookup"
   }
   if (props.target == "component") {
      return "Component Lookup"
   }
   return "Order Lookup"
})

async function lookupRecords() {
   selectedHit.value = null
   searched.value = false
   if (props.target == "metadata") {
      await metadataStore.lookup( query.value )
   } else if  (props.target == "component") {
      await componentsStore.lookup( query.value )
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
button.lookup-button {
   margin-left: 5px;
}
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
