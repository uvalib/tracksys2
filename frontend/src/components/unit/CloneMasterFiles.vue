<template>
   <WaitSpinner v-if="cloneStore.status == 'cloning'" :overlay="true" message="<div>Please wait...</div><p>Masterfile clone is in progress</p>" />
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
         @move-to-source=" removeCloneItem()"
         @move-all-to-source=" removeCloneItem()"
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
                     <DataDisplay label="Unit ID" :value="slotProps.item.unitID" />
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
         <DPGButton label="Clone Master Files" @click="cloneClicked()" :disabled="masterFiles[1].length == 0 || cloneStore.inProgress" />
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref, watch, nextTick } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useCloneStore } from '@/stores/clone'
import PickList from 'primevue/picklist'
import DataDisplay from '../DataDisplay.vue'
import { useConfirm } from "primevue/useconfirm"
import WaitSpinner from "@/components/WaitSpinner.vue"

const confirm = useConfirm()
const unitsStore = useUnitsStore()
const cloneStore = useCloneStore()

const emit = defineEmits( ['canceled', 'cloned' ])

const selectedUnitID = ref(0)
const masterFiles = ref([[],[]])

watch(() => cloneStore.status, (newStatus) => {
   if ( newStatus == 'success') {
      unitsStore.getMasterFiles( unitsStore.detail.id )
      emit("cloned")
   }
})

onMounted(() => {
   cloneStore.getSourceUnits( unitsStore.detail.id)
})

const removeCloneItem = (() => {
   nextTick( () => {
      masterFiles.value[0] = cloneStore.masterFiles.filter( isInCloneList )
   })
})

const isInCloneList = (( mf ) => {
   let cloneIdx =  masterFiles.value[1].findIndex( cmf => cmf.id == mf.id)
   return cloneIdx == -1
})

const unitPicked = ( async () => {
   await cloneStore.getMasterFiles(selectedUnitID.value)
   masterFiles.value[0] = cloneStore.masterFiles.filter( isInCloneList )
})

const cancelClicked = (() => {
   emit("canceled")
})

const cloneClicked = (() => {
   confirm.require({
      message: 'Clone selected master files into this unit?',
      header: 'Confirm Clone',
      icon: 'pi pi-question-circle',
      rejectClass: 'p-button-secondary',
      accept: () => {
         cloneStore.cloneMasterFiles( unitsStore.detail.id, masterFiles.value[1])
      }
   })
})
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