<template>
   <h2>Metadata {{route.params.id}}</h2>
   <div class="details" v-if="systemStore.working==false">
      <Panel header="General Information">

      </Panel>
   </div>
</template>

<script setup>
import { onBeforeMount } from 'vue'
import { useRoute, onBeforeRouteUpdate } from 'vue-router'
import { useSystemStore } from '@/stores/system'
import Panel from 'primevue/panel'

const route = useRoute()
const systemStore = useSystemStore()

onBeforeRouteUpdate(async (to) => {
   let mdID = to.params.id
   console.log("before metadata " + mdID)
})

onBeforeMount(() => {
   let mdID = route.params.id
   console.log("before metadata " + mdID)
})

</script>

<style scoped lang="scss">
.details {
   padding: 0 25px 10px 25px;
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;

   :deep(div.p-panel) {
      margin: 10px;
      flex-grow: 1;
      text-align: left;
   }
}
</style>