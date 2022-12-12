<template>
   <div v-if="equipment.length == 0">
      <h3>No equipment found</h3>
   </div>
   <DataTable v-else :value="props.equipment" ref="equipmentTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="false" :rows="props.equipment.length"
   >
      <Column field="name" header="Name" class="e-wide"/>
      <Column field="serialNumber" header="Serial Number" class="e-wide"/>
      <Column field="" header="Workstation" class="e-wide">
         <template #body="slotProps">
            <span class="workstation">{{workstation(slotProps.data.id)}}</span>
         </template>
      </Column>
      <Column header="Actions" class="row-acts">
         <template #body="slotProps">
            <DPGButton label="Edit"  class="p-button-secondary first" @click="activateWorkstation(slotProps.data.id)"/>
            <DPGButton label="Retire"  class="p-button-secondary" @click="retireWorkstation(slotProps.data.id)"/>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useEquipmentStore } from '@/stores/equipment'

const props = defineProps({
   equipment: {
      type: Array,
      required: true
   }
})

const equipmentStore = useEquipmentStore()

function workstation( equipID ) {
   let wsName = ""
   equipmentStore.workstations.some( ws => {
      let equip = ws.equipment.find( e => e.id == equipID)
      if (equip ) {
         wsName = ws.name
      }
      return wsName != ""
   })
   if (wsName == "") {
      wsName = "N/A"
   }
   return wsName
}

</script>

<stype scoped lang="scss">
.results {
   margin: 20px;
   font-size: 0.9em;
   td.nowrap, th {
      white-space: nowrap;
   }
   th, td {
      font-size: 0.85em;
   }
   span.ws-status {
      width: 20px;
      height: 20px;
      display: inline-block;
      border-radius: 20px;
      background: var(--uvalib-green);
   }

   span.ws-status.inactive {
      background: var(--uvalib-grey-light);
   }
}
.e-wide {
   width: 30%;
}
:deep(.status-col) {
   width: 25px;
   text-align: center;
   padding: 0;
}
:deep(.row-acts) {
   vertical-align: top;
   display: flex;
   flex-flow: row nowrap;
   width: 175px;

   button.p-button.first {
      margin-right: 10px;
   }

   button.p-button {
      font-size: 0.85em;
      padding: 3px 6px;
      display: block;
      width: 100%;
      margin-top: 5px;
   }
}
</stype>