<template>
   <h2>{{componentsStore.title}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <Panel>
         <template #header>
            <div class="cmp-header">
               <span class="title">Component Hierarchy</span>
               <span class="hint">Scroll down to view related master files</span>
            </div>
         </template>
         <Tree :value="componentsStore.nodes" :expandedKeys="expandedKeys"
            :filter="true" filterMode="lenient" scrollHeight="450px"
            selectionMode="single" v-model:selectionKeys="selectedKey" @node-select="componentSelected"
         >
            <template #default="slotProps">
               <dl :id="`node-${slotProps.node.data.id}`">
                  <DataDisplay label="ID" :value="slotProps.node.data.id" />
                  <DataDisplay label="PID" :value="slotProps.node.data.pid" />
                  <DataDisplay label="Type" :value="slotProps.node.data.componentType" />
                  <DataDisplay v-if="slotProps.node.data.title"  label="Title" :value="slotProps.node.data.title.trim()" />
                  <DataDisplay v-if="slotProps.node.data.label"  label="Label" :value="slotProps.node.data.label.trim()" />
                  <DataDisplay v-if="slotProps.node.data.description"  label="Description" :value="slotProps.node.data.description.trim()" />
                  <DataDisplay v-if="slotProps.node.data.date" label="Date" :value="slotProps.node.data.date" />
                  <DataDisplay v-if="slotProps.node.data.level" label="Level" :value="slotProps.node.data.level" />
                  <DataDisplay v-if="slotProps.node.data.eadID" label="EAD ID" :value="slotProps.node.data.eadID" />
                  <DataDisplay v-if="slotProps.node.data.barcode" label="Barcode" :value="slotProps.node.data.barcode" />
                  <DataDisplay label="Master Files" :value="slotProps.node.data.masterFileCount" blankValue="0" />
               </dl>
            </template>
         </Tree>
      </Panel>
      <div class="master-files" id="master-files">
         <MasterFiles />
      </div>
   </div>
</template>

<script setup>
import { onBeforeMount, ref, nextTick } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useComponentsStore } from '@/stores/components'
import Tree from 'primevue/tree'
import Panel from 'primevue/panel'
import DataDisplay from '../components/DataDisplay.vue'
import MasterFiles from '../components/component/MasterFiles.vue'

const route = useRoute()
const systemStore = useSystemStore()
const componentsStore = useComponentsStore()

const expandedKeys = ref({})
const selectedKey = ref({})

onBeforeRouteUpdate(async (to) => {
   let cID = to.params.id
   await componentsStore.getComponentTree(cID)
   expandSelectedComponent( cID )
})

onBeforeMount( async () => {
   let cID = route.params.id
   document.title = componentsStore.title
   await componentsStore.getComponentTree(cID)
   expandSelectedComponent( cID )
})

function componentSelected( tgtComponent ) {
   let cID = parseInt(tgtComponent.key, 10)
   componentsStore.loadRelatedMasterFiles(cID)
   let mfEle = document.getElementById("master-files")
   var headerOffset = 40
   var elementPosition = mfEle.getBoundingClientRect().top
   var offsetPosition = elementPosition - headerOffset
   window.scrollBy({
      top: offsetPosition,
      behavior: "smooth"
   })
}

function expandSelectedComponent( cID ) {
   let nodes = componentsStore.nodes
   if ( cID != nodes[0].key) {
      expandedKeys.value[nodes[0].key] = true
      findNode(nodes[0], cID)
   }

   nextTick( () =>{
      let eleID = `node-${cID}`
      let componentEle = document.getElementById(eleID)
      var w = document.getElementsByClassName("p-tree-wrapper")[0]
      var elementPosition = componentEle.offsetTop - w.offsetTop - 40

      w.scrollBy({
         top: elementPosition,
         behavior: "smooth"
      })
      selectedKey.value[cID] = true
   })
}

function findNode( currNode, tgtKey) {
   if ( currNode.key == tgtKey) return true

   let found = false
   currNode.children.some( n => {
      if (n.key == tgtKey) {
         found = true
         expandedKeys.value[currNode.key] = true
      } else {
         if (findNode(n, tgtKey)) {
            found = true
            expandedKeys.value[n.key] = true
         }
      }
      return found == true
   })

   return found
}

</script>

<style scoped lang="scss">
.details {
   padding: 10px 50px;

   .cmp-header {
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      width: 100%;
      .title {
         font-weight: 600;
      }
      .hint {
         font-size: 0.9em;
         color: #aaa;
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