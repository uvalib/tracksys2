<template>
   <DPGButton @click="show" severity="secondary" label="Add to Collection"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add to Collection" :style="{width: '780px'}">
      <VirtualScroller :items="collectionStore.collections" :itemSize="30" showLoader class="collections" :showLoader="system.working" >
         <template v-slot:item="{ item }">
            <div class="collection" >
               <RadioButton v-model="collectionID" :inputId="item.pid" name="collection" :value="item.id" />
               <label :for="item.pid" class="ml-2">{{ item.pid }} - {{ item.title }}</label>
            </div>
         </template>
      </VirtualScroller>
      <div class="acts">
         <DPGButton @click="hide" severity="secondary" label="Cancel"/>
         <DPGButton @click="addToCollection" label="Add" :disabled="collectionID==0"/>
      </div>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useCollectionsStore } from '@/stores/collections'
import { useSystemStore } from '@/stores/system'
import VirtualScroller from 'primevue/virtualscroller'
import RadioButton from 'primevue/radiobutton'

const props = defineProps({
   metadataID: {
      type: Number,
      required: true,
   },
})

const collectionStore = useCollectionsStore()
const system = useSystemStore()

const isOpen = ref(false)
const collectionID = ref(0)

const addToCollection = (async () => {
   await collectionStore.addToCollection( collectionID.value, props.metadataID )
   isOpen.value=false
})

const hide = (() => {
   isOpen.value=false
})

const show = (() => {
   collectionID.value = 0
   collectionStore.getCollections()
   isOpen.value = true

})
</script>

<style lang="scss" scoped>
.collections {
   min-height: 350px;
   border: 1px solid var(--uvalib-grey-light);
   margin: 10px 5px 20px 5px;
   .collection {
      padding:5px;
      display: flex;
      flex-flow: row nowrap;
      gap: 10px;
      justify-self: flex-start;
      align-items: center;
      label {
         white-space: nowrap;
      };
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   gap: 10px;
}
</style>
