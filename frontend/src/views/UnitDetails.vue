<template>
   <h2>Unit {{route.params.id}}</h2>
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
   let uID = to.params.id
   console.log("before unit " + uID)
})

onBeforeMount(() => {
   let uID = route.params.id
   console.log("before unit " + uID)
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