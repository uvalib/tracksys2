<template>
   <div v-if="equipment.length == 0">
      <h3>No equipment found</h3>
   </div>
   <DataTable v-else :value="props.equipment" ref="equipmentTable" dataKey="id"
      stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
      :lazy="false" :paginator="false" :rows="props.equipment.length"
   >
      <Column header="" headerStyle="width: 3em" v-if="equipmentStore.pendingEquipment.workstationID > 0">
         <template #body="slotProps">
            <Checkbox :value="slotProps.data" v-model="equipmentStore.pendingEquipment.equipment" :disabled="isItemDisabled(slotProps.data.id)" @click="equipmentClicked"/>
         </template>
      </Column>
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
            <DPGButton label="Retire"  class="p-button-secondary" @click="retireWorkstation(slotProps.data.id)" :disabled="workstation(slotProps.data.id) != 'N/A'"/>
         </template>
      </Column>
   </DataTable>
</template>

<script setup>
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { useEquipmentStore } from '@/stores/equipment'
import Checkbox from 'primevue/checkbox'

const props = defineProps({
   equipment: {
      type: Array,
      required: true
   }
})

const equipmentStore = useEquipmentStore()

function equipmentClicked() {
   equipmentStore.pendingEquipment.changed = true
}

function isItemDisabled(equipID) {
   let equipWS = workstation(equipID)
   if  ( equipWS == "N/A") return false
   let tgtWS = equipmentStore.workstations.find( ws => ws.id == equipmentStore.pendingEquipment.workstationID)
   if (tgtWS) {
      return tgtWS.name != equipWS
   }
   return true
}

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