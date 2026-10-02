<template>
   <h2>Patron Deliveries</h2>
   <div class="p-4">
      <WaitSpinner v-if="statsStore.deliveries.loading"/>
      <LineChart v-else :data="statsStore.deliveries" :options="options"/>
      <p class="error" v-if="statsStore.deliveries.error">{{statsStore.deliveries.error}}</p>
      <div class="row-right items-center p-4">
         <label>Year:</label>
         <UInput v-model="tgtYear" />
         <UButton color="secondary" @click="loadStats" label="Generate" />
      </div>
   </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {useStatsStore} from '@/stores/statistics'
import LineChart from "./LineChart.vue"
import WaitSpinner from "@/components/WaitSpinner.vue"

const tgtYear = ref( new Date().getFullYear() )
const statsStore = useStatsStore()

const options = ref({
   responsive: true,
})

const loadStats = ( () => {
   statsStore.getPatronDeliveries(tgtYear.value)
})

onMounted( () => {
   loadStats()
})
</script>

<style scoped lang="scss">

h3 {
   margin: 10px 0 5px 10px;
   padding-bottom: 5px;
   text-align: left;
   border-bottom: 1px solid var(--uvalib-grey-light);
}
</style>