<template>
   <h2>{{componentsStore.title}}</h2>
   <div class="column gap-6 p-6" v-if="systemStore.working==false">
      <div class="row-left gap-6 items-start">
         <UCard class="basis-2/3 text-sm">
            <template #header>
               <div class="row-between items-center">
                  <span>Component Hierarchy</span>
                  <UButton label="Export JSON" size="sm" color="secondary" @click="exportJSON()"/>
               </div>
            </template>
            <UTree :items="componentsStore.nodes" class="h-75" :get-key="i => i.key" virtualize
               v-model="selectedItem" v-model:expanded="expandedKeys" @update:modelValue="nodeSelected" />
         </UCard>
         <UCard title="Details" class="basis-1/3 text-sm" >
            <dl v-if="selectedItem">
               <DataDisplay label="ID" :value="selectedItem.data.id" />
               <DataDisplay label="PID" :value="selectedItem.data.pid" />
               <DataDisplay label="Type" :value="selectedItem.data.componentType" />
               <DataDisplay v-if="selectedItem.data.title"  label="Title" :value="selectedItem.data.title.trim()" />
               <DataDisplay v-if="selectedItem.data.label"  label="Label" :value="selectedItem.data.label.trim()" />
               <DataDisplay v-if="selectedItem.data.description"  label="Description" :value="selectedItem.data.description.trim()" />
               <DataDisplay v-if="selectedItem.data.date" label="Date" :value="selectedItem.data.date" />
               <DataDisplay v-if="selectedItem.data.level" label="Level" :value="selectedItem.data.level" />
               <DataDisplay v-if="selectedItem.data.eadID" label="EAD ID" :value="selectedItem.data.eadID" />
               <DataDisplay v-if="selectedItem.data.barcode" label="Barcode" :value="selectedItem.data.barcode" />
               <DataDisplay label="Master Files" :value="selectedItem.data.masterFileCount" blankValue="0" />
            </dl>
         </UCard>
      </div>
      <div class="master-files" id="master-files" v-if="selectedItem">
         <MasterFiles />
      </div>
   </div>
</template>

<script setup>
import { onBeforeMount, ref, nextTick } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useComponentsStore } from '@/stores/components'
import MasterFiles from '../components/component/MasterFiles.vue'

const route = useRoute()
const systemStore = useSystemStore()
const componentsStore = useComponentsStore()

const selectedItem = ref()
const expandedKeys = ref([])

onBeforeRouteUpdate(async (to) => {
   let cID = to.params.id
   await componentsStore.getComponentTree(cID)
   expandSelectedComponent( cID )
})

onBeforeMount( async () => {
   let cID = route.params.id
   document.title = componentsStore.title
   await componentsStore.getComponentTree(cID)
   if (componentsStore.selectedComponent > -1) {
      expandSelectedComponent( cID )
   }
})

const exportJSON = (() => {
   const fileURL = window.URL.createObjectURL(new Blob([JSON.stringify(componentsStore.nodes, null, 3)], { type: 'application/json' }))
   const fileLink = document.createElement('a')
   fileLink.href =  fileURL
   fileLink.setAttribute('download', `comopnent.json`)
   document.body.appendChild(fileLink)
   fileLink.click()
   window.URL.revokeObjectURL(fileURL)
})

const nodeSelected = ((node) => {
   if ( node ) {
      componentsStore.loadRelatedMasterFiles(node.key)
   }
})

const expandSelectedComponent = (( cID ) => {
   expandedKeys.value = []
   let nodes = componentsStore.nodes
   if ( cID != nodes[0].key) {
      expandedKeys.value.push( nodes[0].key )
      selectedItem.value = findNode(nodes[0], cID)
   }

   
   
   nextTick( () => {
     scrollToTarget()
   })
})

const scrollToTarget = (async () => {
   let found = false 
   while ( found == false ) {
      let el = document.querySelector('[aria-selected="true"]')
      if (el) {
         await el.scrollIntoView({ block: 'nearest', behavior: 'instant' })
         found = true
      } else {
         let eles = document.querySelectorAll("li:nth-last-of-type(1)")
         await eles[1]?.scrollIntoView({ block: 'nearest', behavior: 'instant' }) 
      }
   }
})

const addExpandKey = ((key) => {
   if ( expandedKeys.value.includes(key) == false ) {
      expandedKeys.value.push(key)
   }
})

const findNode = (( currNode, tgtKey) => {
   if ( currNode.key == tgtKey) {
      return currNode 
   }

   let found = null 
   currNode.children.some( n => {
      if (n.key == tgtKey) {
         found = n
         addExpandKey(currNode.key)
      } else {
         found = findNode(n, tgtKey)
         if (found) {
            addExpandKey(n.key)
         }
      }
      return found != null 
   })

   return found
})

</script>

<style scoped lang="scss">
.details {
   padding: 10px 50px;

   .cmp-header {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      align-items: center;
      width: 100%;
      .title {
         font-weight: 600;
      }
      .hint {
         font-size: 0.9em;
         color: #aaa;
      }
      button {
         font-size: 0.85em;
         padding: 3px 15px;
         margin: 0;
      }
   }
   .master-files {
      margin-top: 15px;
   }
   .p-tree.p-component {
      border:0;
      padding: 0;
   }
   :deep(dd) {
      margin-bottom: 4px !important;
   }
   :deep(span.p-treenode-label) {
      text-align: left;
      width: 100%;
      padding: 10px 0;
   }
   :deep(.p-treenode-content.p-treenode-selectable) {
      border: 1px solid var(--uvalib-grey-light);
   }
   :deep(.p-treenode.p-treenode-leaf) {
      padding-right: 0;
      margin: 2px 0;
   }

   :deep(.p-tree .p-treenode-children) {
      padding: 0 0 0 3rem;
      font-size: 0.95em;
   }

}
</style>