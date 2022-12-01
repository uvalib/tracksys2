<template>
   <h2>{{componentsStore.title}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <Tree :value="componentsStore.nodes" :expandedKeys="expandedKeys"
         :filter="true" filterMode="strict"
         selectionMode="single" v-model:selectionKeys="selectedKey">
         <template #default="slotProps">
            <dl :id="`node-${slotProps.node.data.id}`">
               <DataDisplay label="ID" :value="slotProps.node.data.id" />
               <DataDisplay label="PID" :value="slotProps.node.data.pid" />
               <DataDisplay label="Type" :value="slotProps.node.data.componentType" />
               <DataDisplay v-if="slotProps.node.data.title"  label="Title" :value="slotProps.node.data.title" />
               <DataDisplay v-if="slotProps.node.data.label"  label="Label" :value="slotProps.node.data.label" />
               <DataDisplay v-if="slotProps.node.data.description"  label="Description" :value="slotProps.node.data.description" />
               <DataDisplay v-if="slotProps.node.data.date" label="Date" :value="slotProps.node.data.date" />
               <DataDisplay v-if="slotProps.node.data.level" label="Level" :value="slotProps.node.data.level" />
               <DataDisplay v-if="slotProps.node.data.eadID" label="EAD ID" :value="slotProps.node.data.eadID" />
               <DataDisplay v-if="slotProps.node.data.barcode" label="Barcode" :value="slotProps.node.data.barcode" />
            </dl>
         </template>
      </Tree>
   </div>
</template>

<script setup>
import { onBeforeMount, ref, nextTick } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import { useComponentsStore } from '@/stores/components'
import Tree from 'primevue/tree'
import DataDisplay from '../components/DataDisplay.vue'

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

function expandSelectedComponent( cID ) {
   let nodes = componentsStore.nodes
   if ( cID == nodes[0].key) {
      return
   }
   expandedKeys.value[nodes[0].key] = true
   findNode(nodes[0], cID)

   nextTick( () =>{
      let eleID = `node-${cID}`
      let componentEle = document.getElementById(eleID)
      var headerOffset = 40
      var elementPosition = componentEle.getBoundingClientRect().top
      var offsetPosition = elementPosition - headerOffset
      window.scrollBy({
         top: offsetPosition,
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
   padding: 15px 25px;
   :deep(span.p-treenode-label) {
      text-align: left;
      // border-bottom: 2px solid var(--uvalib-grey-light);
      // border-top: 1px solid var(--uvalib-grey-light);
      width: 100%;
      padding: 15px;
   }
   :deep(.p-treenode-content.p-treenode-selectable) {
      border: 1px solid var(--uvalib-grey-light);
   }
   :deep(.p-treenode.p-treenode-leaf) {
      padding-right: 0;
      margin: 10px 0;
   }

   :deep(.p-tree .p-treenode-children) {
      padding: 0 0 0 2rem;
   }

}
</style>