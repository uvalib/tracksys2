<script setup>
import { onBeforeMount, ref } from 'vue'
import { useEquipmentStore } from '@/stores/equipment'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import EquipmentPanel from '../components/equipment/EquipmentPanel.vue'

const equipmentStore = useEquipmentStore()

const selectedWorkstation = ref()

onBeforeMount( async () => {
   document.title = "Equipment"
   equipmentStore.getEquipment()
})

function clearSetup() {
   console.log("clear")
}

function statusClass(statusID) {
   if (statusID == 1) {
      return "ws-status inactive"
   }
   return "ws-status active"
}
</script>

<template>
   <h2>Equipment</h2>
   <div class="equipment">
      <div class="columns">
         <Panel header="Workstations">
            <DataTable :value="equipmentStore.workstations" ref="workstationTable" dataKey="id"
               stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
               :lazy="false" :paginator="false"
               selectionMode="single" v-model:selection="selectedWorkstation"
               :rows="equipmentStore.workstations.length" :totalRecords="equipmentStore.workstations.length"
            >
               <Column field="status" header="Status" class="status-col">
                  <template #body="slotProps">
                     <span :class="statusClass(slotProps.data.status)"></span>
                  </template>
               </Column>
               <Column field="name" header="Name" class="wide"/>
               <Column field="projectCount" header="Projects"/>
               <Column header="Actions" class="row-acts">
                  <template #body="slotProps">
                     <DPGButton v-if="slotProps.data.status==0" label="Deactivate"  class="p-button-secondary first" @click="deactivateWorkstation(slotProps.data.id)"/>
                     <DPGButton v-else label="Activate"  class="p-button-secondary first" @click="activateWorkstation(slotProps.data.id)"/>
                     <DPGButton label="Retire"  class="p-button-secondary" @click="retireWorkstation(slotProps.data.id)"/>
                  </template>
               </Column>
            </DataTable>
         </Panel>
         <Panel header="Workstation Setup">
            <h3 v-if="selectedWorkstation == null">Select a workstation to view details</h3>
            <template v-else>
               <DataTable :value="selectedWorkstation.equipment" ref="setupTable" dataKey="id"
                  stripedRows showGridlines class="p-datatable-sm"
               >
                  <Column field="type" header="Type"/>
                  <Column field="name" header="Name"/>
                  <Column field="serialNumber" header="Serial Number"/>

               </DataTable>
               <div class="setup-acts">
                  <DPGButton label="Clear Setup" class="p-button-secondary" @click="clearSetup"/>
                  <DPGButton label="Save Setup Changes" class="p-button-secondary" @click="clearSetup"/>
               </div>
            </template>
         </Panel>
      </div>
      <Panel header="Equipment">
         <TabView class="results">
            <TabPanel header="Camera Bodies">
               <EquipmentPanel :equipment="equipmentStore.cameraBodies" />
            </TabPanel>
            <TabPanel header="Lenses">
               <EquipmentPanel :equipment="equipmentStore.lenses" />
            </TabPanel>
            <TabPanel header="Digital Backs">
               <EquipmentPanel :equipment="equipmentStore.digitalBacks" />
            </TabPanel>
            <TabPanel header="Scanners">
               <EquipmentPanel :equipment="equipmentStore.scanners" />
            </TabPanel>
         </TabView>
      </Panel>
   </div>
</template>

<style scoped lang="scss">
.equipment {
   min-height: 600px;
   padding: 15px 25px;
   h3 {
      text-align: center;
      color: var(--uvalib-text);
      font-weight: 500;
      font-size: 1em;
   }

   .setup-acts {
      padding: 15px 0 0 0;
      text-align: right;
      button.p-button {
         margin-left: 10px;
         font-size: 0.85em;
      }
   }

   .p-component.p-datatable {
      font-size: 0.9em;
   }

   :deep(.wide) {
      width: 100%;
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

   .columns {
      padding: 0 25px 10px 25px;
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-start;

      div.p-panel {
         margin: 10px;
         flex: 45%;
         text-align: left;
      }
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
</style>