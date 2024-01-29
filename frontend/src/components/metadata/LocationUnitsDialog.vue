<template>
   <Dialog v-model:visible="showDialog" :style="{width: '1000px'}" :header="`Units in Folder ${props.folder}`" :modal="true" position="top" @hide="emit('closed')">
      <div  v-if="metadataStore.locationUnits">
         <p>Multiple units contain digitized materials from this folder. Pick one to view:</p>
         <div class="units">
            <DataTable :value="metadataStore.locationUnits" ref="folderUnitsTable" dataKey="id"
               stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm" :lazy="false"
            >
               <Column field="id" header="">
                  <template #body="slotProps">
                     <DPGButton @click="viewUnitClicked(slotProps.data.id)" label="View" class="p-button-secondary"/>
                  </template>
               </Column>
               <Column field="orderID" header="Order ID" class="nowrap" />
               <Column field="orderTitle" header="Order Title" />
               <Column field="id" header="Unit ID"  class="nowrap"/>
               <Column field="intendedUseID" header="Intended Use"  class="nowrap">
                  <template #body="slotProps">{{ intendedUse(slotProps.data.intendedUseID) }}</template>
               </Column>
               <Column field="completeScan" header="Complete Scan" class="nowrap"/>
               <Column field="dateArchived" header="Date Ingested"  class="nowrap">
                  <template #body="slotProps">
                     <span v-if="slotProps.data.dateArchived">
                        {{ $formatDate(slotProps.data.dateArchived) }}
                     </span>
                     <span v-else class="na">N/A</span>
                  </template>
               </Column>
               <Column field="staffNotes" header="Staff Notes" />

            </DataTable>
         </div>
      </div>
      <div v-else class="no-units">
         No units exist for this folder.
      </div>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import { useMetadataStore } from '@/stores/metadata'
import { useSystemStore } from '@/stores/system'
import Dialog from 'primevue/dialog'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const metadataStore = useMetadataStore()
const systemStore = useSystemStore()

const emit = defineEmits( ['closed', 'unit' ])
const props = defineProps({
   folder: {
      type: String,
      required: true
   },
})

const showDialog = ref(true)

const intendedUse = ( (id) => {
   let iu = systemStore.intendedUses.find( u => u.id == id)
   if ( iu ) {
      return iu.name
   }
   return "Unknown"
})

const viewUnitClicked = ((unitID) => {
   emit("unit", unitID)
   showDialog.value = false
})
</script>

<style scoped lang="scss">
:deep(th.nowrap), :deep(td.nowrap) {
   white-space: nowrap !important;
}
.no-units {
   text-align: center;
   font-size: 1.2em;
   padding: 30px;
}
p {
   text-align: center;
   padding: 0;
   margin: 10px 0 25px 0;
   font-weight: 500;
}
.units .p-datatable {
   font-size: 0.85em;
   margin: 0 10px 15px 10px;

   .na {
      color: #ccc;
   }
}
</style>