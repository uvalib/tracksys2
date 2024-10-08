<template>
   <h2>
      <span>Equipment</span>
      <div class="actions" v-if="(userStore.isAdmin || userStore.isSupervisor)" >
         <AddWorkstationDialog />
         <AddEquipmentDialog />
      </div>
   </h2>
   <div class="equipment">
      <div class="columns">
         <Panel header="Workstations">
            <DataTable :value="equipmentStore.workstations" ref="workstationTable" dataKey="id"
               stripedRows showGridlines responsiveLayout="scroll" class="p-datatable-sm"
               @rowSelect="workstationSelected"
               selectionMode="single" v-model:selection="selectedWorkstation"
               :rows="equipmentStore.workstations.length" :totalRecords="equipmentStore.workstations.length"
            >
               <Column field="status" header="Active" class="status-col">
                  <template #body="slotProps">
                     <span :class="statusClass(slotProps.data.status)"></span>
                  </template>
               </Column>
               <Column field="name" header="Name" class="wide"/>
               <Column field="projectCount" header="Projects"/>
               <Column header="Actions" class="row-acts">
                  <template #body="slotProps">
                     <DPGButton v-if="slotProps.data.status==0" label="Deactivate"  severity="secondary" @click="deactivateWorkstation(slotProps.data.id)"/>
                     <DPGButton v-else label="Activate" severity="secondary" @click="activateWorkstation(slotProps.data.id)"/>
                     <DPGButton label="Retire"  severity="secondary" @click="retireWorkstation(slotProps.data.id)" :disabled="slotProps.data.projectCount > 0"/>
                  </template>
               </Column>
            </DataTable>
         </Panel>
         <Panel :header="setupHeader">
            <h3 v-if="selectedWorkstation == null">Select a workstation to view details</h3>
            <template v-else>
               <DataTable :value="equipmentStore.pendingEquipment.equipment" ref="setupTable" dataKey="id"
                  stripedRows showGridlines class="p-datatable-sm"
               >
                  <Column field="type" header="Type"/>
                  <Column field="name" header="Name"/>
                  <Column field="serialNumber" header="Serial Number"/>

               </DataTable>
               <div class="setup-acts">
                  <DPGButton label="Clear Setup" severity="secondary" @click="clearSetup" :disabled="clearAllDisabled"/>
                  <DPGButton label="Save Setup Changes" severity="secondary" @click="saveSetup"
                     :disabled="!(equipmentStore.pendingEquipment.changed==true && equipmentStore.pendingEquipment.equipment.length > 0)"/>
               </div>
            </template>
         </Panel>
      </div>
      <Panel>
         <template #header>
            <div class="equip-header">
               <span class="title">Equipment</span>
               <span class="ws"><b>Workstation</b>: {{selectedWSName}}</span>
            </div>
         </template>
         <Tabs value="bodies" :lazy="true">
            <TabList>
               <Tab value="bodies">Camera Bodies</Tab>
               <Tab value="lenses">Lenses</Tab>
               <Tab value="backs">Digital Backs</Tab>
               <Tab value="scanners">Scanners</Tab>
            </TabList>
            <TabPanels>
               <TabPanel value="bodies">
                  <EquipmentPanel :equipment="equipmentStore.cameraBodies" />
               </TabPanel>
               <TabPanel value="lenses">
                  <EquipmentPanel :equipment="equipmentStore.lenses" />
               </TabPanel>
               <TabPanel value="backs">
                  <EquipmentPanel :equipment="equipmentStore.digitalBacks" />
               </TabPanel>
               <TabPanel value="scanners">
                  <EquipmentPanel :equipment="equipmentStore.scanners" />
               </TabPanel>
            </TabPanels>
         </Tabs>
      </Panel>
   </div>
</template>

<script setup>
import { onBeforeMount, ref, computed } from 'vue'
import { useEquipmentStore } from '@/stores/equipment'
import { useUserStore } from '@/stores/user'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tabs from 'primevue/tabs'
import TabList from 'primevue/tablist'
import Tab from 'primevue/tab'
import TabPanels from 'primevue/tabpanels'
import TabPanel from 'primevue/tabpanel'
import EquipmentPanel from '@/components/equipment/EquipmentPanel.vue'
import { useConfirm } from "primevue/useconfirm"
import AddWorkstationDialog from '@/components/equipment/AddWorkstationDialog.vue'
import AddEquipmentDialog from '@/components/equipment/AddEquipmentDialog.vue'

const confirm = useConfirm()
const equipmentStore = useEquipmentStore()
const userStore = useUserStore()

const selectedWorkstation = ref()

const setupHeader = computed(() => {
   let name = selectedWSName.value
   if (name == "None") return "Workstation Setup"
   return name+" Setup"
})
const selectedWSName = computed(() => {
   if ( selectedWorkstation.value ) {
      let ws = equipmentStore.workstations.find(ws => ws.id == selectedWorkstation.value.id)
      if (ws) {
         return ws.name
      }
   }
   return "None"
})
const clearAllDisabled = computed( () => {
   let ws = equipmentStore.workstations.find(ws => ws.id == selectedWorkstation.value.id)
   if (ws.projectCount > 0) return true
   return equipmentStore.pendingEquipment.equipment.length == 0
})

onBeforeMount( async () => {
   document.title = "Equipment"
   equipmentStore.getEquipment()
})

const deactivateWorkstation = (( wsID ) => {
   equipmentStore.deactivateWorkstation(wsID)
})

const activateWorkstation = (( wsID ) => {
   equipmentStore.activateWorkstation(wsID)
})

const retireWorkstation = (( wsID ) => {
   let ws = equipmentStore.workstations.find(ws => ws.id == wsID)
   confirm.require({
      message: `Retire workstation '${ws.name}'?`,
      header: 'Confirm Retire',
      icon: 'pi pi-question-circle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Retire'
      },
      accept: () => {
         equipmentStore.retireWorkstation(wsID)
      }
   })
})

const workstationSelected = (() => {
   equipmentStore.workstationSelected( selectedWorkstation.value.id )
})

const clearSetup = (() => {
   equipmentStore.clearSetup()
})

const saveSetup = (() => {
   equipmentStore.saveSetup()
})

const statusClass = ((statusID) => {
   if (statusID == 1) {
      return "ws-status inactive"
   }
   return "ws-status active"
})
</script>

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
   .equip-header {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      width: 100%;
      .title {
         font-weight: 600;
      }
      b {
         font-weight: 600;
      }
   }
   .setup-acts {
      padding: 15px 0 0 0;
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      align-items: flex-start;
      gap: 10px;
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
      gap: 5px;
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