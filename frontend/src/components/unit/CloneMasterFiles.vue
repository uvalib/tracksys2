<template>
   <div v-if="cloneStore.initialized && cloneStore.sourceUnits.length == 0">
      <p>There are no suittable units to clone master files from.</p>
      <div class="acts">
         <DPGButton label="OK" @click="cancelClicked()" />
      </div>
   </div>
   <div v-else class="clone-masterfiles">
      <div class="units">
         <label>Source unit:</label>
         <select v-model="selectedUnitID" @change="unitPicked">
            <option disabled value="0">Select a unit</option>
            <option v-for="opt in cloneStore.sourceUnits" :key="`as-${opt.id}`" :value="opt.id">{{opt.id}}</option>
         </select>
      </div>
      <PickList v-if="masterFiles" v-model="masterFiles" dataKey="id"
         listStyle="height:500px" :showSourceControls="false"
      >
         <template #sourceheader>
            Source Master Files
         </template>
         <template #targetheader>
            Master Files
         </template>
         <template #item="slotProps">
            <div class="master-file">
               <div>
                  <dl>
                     <DataDisplay label="Filename" :value="slotProps.item.filename" />
                     <DataDisplay label="Title" :value="slotProps.item.title" />
                  </dl>
               </div>
               <img :src="slotProps.item.thumbnailURL">
            </div>
         </template>
      </PickList>
      <div class="acts">
         <DPGButton label="Cancel" class="p-button-secondary" @click="cancelClicked()" />
         <DPGButton label="Clone Master Files" @click="cloneClicked()" :disabled="masterFiles[1].length == 0" />
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useCloneStore } from '@/stores/clone'
import PickList from 'primevue/picklist'
import DataDisplay from '../DataDisplay.vue';

const unitsStore = useUnitsStore()
const cloneStore = useCloneStore()

const emit = defineEmits( ['canceled', 'cloned' ])

const selectedUnitID = ref(0)
const masterFiles = ref([[],[]])

onMounted(() => {
   cloneStore.getSourceUnits( unitsStore.detail.id)
})

async function unitPicked() {
   await cloneStore.getMasterFiles(selectedUnitID.value)
   masterFiles.value=[cloneStore.masterFiles, []]
}

function cancelClicked() {
   emit("canceled")
}
function cloneClicked() {
   console.log(masterFiles.value[1])
   emit("cloned")
}
</script>

<style scoped lang="scss">
.units {
   padding-bottom: 15px;
   label {
      font-weight: 600;
   }
   select {
      width: auto;
      display: inline-block;
      margin: 10px;
   }
}
.master-file {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   align-items: flex-start;
   border-bottom: 1px solid var(--uvalib-grey-light);
   padding-bottom: 15px;
}
.acts {
   text-align: center;
   padding: 25px 0 0 0;
   button.p-button {
      margin-left: 10px;
   }
}
</style>