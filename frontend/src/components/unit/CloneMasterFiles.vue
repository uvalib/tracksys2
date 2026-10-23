<template>
   <WaitSpinner v-if="cloneStore.status == 'cloning'" :overlay="true" message="<div>Please wait...</div><p>Masterfile clone is in progress</p>" />
   <div v-if="cloneStore.initialized && cloneStore.sourceUnits.length == 0">
      <div>There are no suittable units to clone master files from.</div>
      <UButton label="OK" color="secondary" @click="cancelClicked()" />
   </div>
   <div v-else class="column gap-4">
      <label>Source unit:
         <USelect v-model="selectedUnitID" @change="unitPicked" :items="cloneStore.sourceUnits" 
            valueKey="id" labelKey="id" class="w-fit min-w-30" placeholder="Unit.." />
      </label>
      <div class="row-left gap-4">
         <UCard title="Source Master Files" class="flex-1">
            <div v-if="cloneStore.status == 'loading'" class="text-center">
               <WaitSpinner :overlay="false" message="Loading files..." />
            </div>
            <UScrollArea v-else :items="srcMasterfiles" v-slot="{ item }" class="h-150" virtualize>
               <div role="button" tabindex="0" @click="item.selected = !item.selected"
                  class="image row-left gap-4 border border-brand-grey-100 p-4 w-full my-1!"
                  :class="{selected: item.selected}"
               >
                  <dl class="w-full text-sm">
                     <DataDisplay label="Unit ID" :value="item.unitID" />
                     <DataDisplay label="Filename" :value="item.filename" />
                     <DataDisplay label="Title" :value="item.title" />
                  </dl> 
                  <img :src="item.thumbnailURL" class="max-h-25">
               </div>
            </UScrollArea>
         </UCard>

         <div class="flex items-center justify-center gap-2 md:flex-col">
            <UButton icon="i-lucide-chevrons-right"color="secondary" @click="moveAllTo('clone')"/>
            <UButton icon="i-lucide-chevron-right" color="secondary" :disabled="srcSelected == false" @click="moveTo('clone')"/>
            <UButton icon="i-lucide-chevron-left" color="secondary" :disabled="cloneSelected == false" @click="moveTo('src')" />
            <UButton icon="i-lucide-chevrons-left"color="secondary" @click="moveAllTo('src')"/>
         </div>
         
         <UCard title="Clone Master Files" class="flex-1">
            <UScrollArea :items="cloneMasterfiles" v-slot="{ item }" class="h-150" virtualize>
               <div role="button" tabindex="0" @click="item.selected = !item.selected"
                  class="image row-left gap-4 border border-brand-grey-100 p-4 w-full my-1!"
                  :class="{selected: item.selected}"
               >
                  <dl class="w-full text-sm">
                     <DataDisplay label="Unit ID" :value="item.unitID" />
                     <DataDisplay label="Filename" :value="item.filename" />
                     <DataDisplay label="Title" :value="item.title" />
                  </dl> 
                  <img :src="item.thumbnailURL" class="max-h-25">
               </div>
            </UScrollArea>   
         </UCard>

         <div class="flex items-center justify-center gap-2 md:flex-col">
            <UButton icon="i-lucide-chevrons-up" color="secondary" :disabled="cloneSelectCount != 1" @click="reorder('first')"/>
            <UButton icon="i-lucide-chevron-up"color="secondary" :disabled="cloneSelectCount != 1" @click="reorder('up')"/>
            <UButton icon="i-lucide-chevron-down" color="secondary":disabled="cloneSelectCount != 1" @click="reorder('down')"/>
            <UButton icon="i-lucide-chevrons-down"color="secondary" :disabled="cloneSelectCount != 1" @click="reorder('last')"/>
         </div>
      </div>
      <div class="row-right gap-4">
         <UButton label="Cancel" color="secondary" @click="cancelClicked()" />
         <UButton label="Clone Master Files" @click="cloneClicked()" :disabled="cloneMasterfiles.length == 0 || cloneStore.inProgress" />
      </div>
   </div>
</template>

<script setup>
import { onMounted, ref, watch, computed } from 'vue'
import { useUnitsStore } from '@/stores/units'
import { useCloneStore } from '@/stores/clone'
import DataDisplay from '../DataDisplay.vue'
import { useConfirm } from "@/composables/useConfirm"
import WaitSpinner from "@/components/WaitSpinner.vue"

const unitsStore = useUnitsStore()
const cloneStore = useCloneStore()

const emit = defineEmits( ['canceled', 'cloned' ])

const selectedUnitID = ref(null)
const srcMasterfiles = ref([])
const cloneMasterfiles = ref([])

watch(() => cloneStore.status, (newStatus) => {
   if ( newStatus == 'success') {
      unitsStore.getMasterFiles( unitsStore.detail.id )
      emit("cloned")
   }
})

onMounted(() => {
   cloneStore.getSourceUnits( unitsStore.detail.id)
   srcMasterfiles.value = []
   cloneMasterfiles.value = []
})

const isInCloneList = (( mf ) => {
   let cloneIdx =  cloneMasterfiles.value.findIndex( cmf => cmf.id == mf.id)
   return cloneIdx == -1
})

const srcSelected = computed(() => {
   return srcMasterfiles.value.findIndex( mf => mf.selected) > -1
})
const cloneSelected = computed(() => {
   return cloneMasterfiles.value.findIndex( mf => mf.selected) > -1
})
const cloneSelectCount = computed(() => {
   let cnt = 0 
   cloneMasterfiles.value.forEach( mf => {
      if (mf.selected) cnt++
   })
   return cnt
})
const reorder = ((place) => {
   // reorder os only active when one clone file is selected
   const tgtIdx = cloneMasterfiles.value.findIndex(mf => mf.selected) 
   if ( tgtIdx == -1 ) return 
   const tgt = cloneMasterfiles.value[tgtIdx]
   if ( place == 'first') {
      if (tgtIdx > 0) {
         cloneMasterfiles.value.splice(tgtIdx, 1)
         cloneMasterfiles.value.unshift(tgt)
      }
   } else if ( place == 'last') {
      if (tgtIdx < (cloneMasterfiles.value.length-1)) {
         cloneMasterfiles.value.splice(tgtIdx, 1)
         cloneMasterfiles.value.push(tgt)
      }
   } else if ( place == 'up' ) {
      if (tgtIdx > 0) {
         cloneMasterfiles.value.splice(tgtIdx, 1)
         cloneMasterfiles.value.splice(tgtIdx-1,0, tgt)   
      }
   } else if ( place == 'down') {
      if (tgtIdx < (cloneMasterfiles.value.length-1)) {
         cloneMasterfiles.value.splice(tgtIdx, 1)
         cloneMasterfiles.value.splice(tgtIdx+1,0, tgt)   
      }
   }
})

const unitPicked = ( async () => {
   await cloneStore.getMasterFiles(selectedUnitID.value)
   srcMasterfiles.value = cloneStore.masterFiles.filter( isInCloneList )
})

const moveTo = ( (tgt) => {
   if ( tgt == 'clone' ) {
      const selected = srcMasterfiles.value.filter( mf => mf.selected )
      srcMasterfiles.value =  srcMasterfiles.value.filter( mf => mf.selected == false)
      cloneMasterfiles.value.push(...selected)
      cloneMasterfiles.value.forEach( mf => mf.selected = false)
   } else {
      // when moving back to src, put files mack in the original order
      const selected = cloneMasterfiles.value.filter( mf => mf.selected )
      cloneMasterfiles.value =  cloneMasterfiles.value.filter( mf => mf.selected == false)
      selected.forEach( mf => {
         mf.selected = false
         srcMasterfiles.value.splice(mf.index, 0, mf)
      })
   }
})
const moveAllTo = ( (tgt) => {
   if ( tgt == 'clone' ) {
      cloneMasterfiles.value.push(...srcMasterfiles.value)
      srcMasterfiles.value = []
   } else {
      // put all items back into src list, tne ensure original order is restored
      srcMasterfiles.value.push(...cloneMasterfiles.value)
      cloneMasterfiles.value = []  
      srcMasterfiles.value.sort( (a,b) => {
         if ( a.index < b.index) return -1
         if ( a.index > b.index) return 1
         return 0
      }) 
   }
})

const cancelClicked = (() => {
   emit("canceled")
})

const cloneClicked = (async () => {
   const msg = 'Clone selected master files into this unit?'
   const resp = await useConfirm("Confirm Clone", msg, "Clone")
   if (resp) {
      cloneStore.cloneMasterFiles( unitsStore.detail.id, cloneMasterfiles.value)
   }
})
</script>

<style scoped lang="scss">
.image {
   cursor: pointer;
   &:hover {
      background-color:  var(--color-brand-blue-alt-300) !important;
   }
   &:focus {
      outline: 2px dashed var(  --color-brand-blue-300 );
      outline-offset: -5px;
   }
}
.image.selected {
   background-color: var(--color-brand-blue-alt-400);
   border-color: var(--color-brand-blue-alt);
}
dl {
   margin: 0 !important;
}
</style>