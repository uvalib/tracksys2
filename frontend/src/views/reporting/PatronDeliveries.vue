<template>
   <h2>Patron Deliveries</h2>
   <div class="report">
      <Chart type="line" :data="statsStore.deliveries" :options="options"/>
      <p class="error" v-if="statsStore.deliveries.error">{{statsStore.deliveries.error}}</p>
      <div class="controls">
         <span class="year-picker">
            <label>Year:<input v-model="tgtYear"></label>
         </span>
         <button @click="loadStats">Generate</button>
      </div>
   </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {useStatsStore} from '@/stores/statistics'
import WaitSpinner from '@/components/WaitSpinner.vue'
import Chart from 'primevue/chart'

const tgtYear = ref( new Date().getFullYear() )
const statsStore = useStatsStore()

const options = ref({
   responsive: true,
   plugins: {
      legend: {
         position: 'top',
      },
   },
})

const loadStats = (() => {
   statsStore.getPatronDeliveries(tgtYear.value)
})

onMounted( () => {
   statsStore.getPatronDeliveries(tgtYear.value)
})
</script>

<style scoped lang="scss">

h3 {
   margin: 10px 0 5px 10px;
   padding-bottom: 5px;
   text-align: left;
   border-bottom: 1px solid var(--uvalib-grey-light);
}
.wait-wrap {
   padding: 20px 10px;
}

.report {
      padding: 10px;
      .controls {
         border-top: 1px solid var(--uvalib-grey-lightest);
         display: flex;
         flex-flow: row wrap;
         justify-content: flex-end;
         padding-top: 15px;
         margin-top: 5px;
         input {
            margin: 0 10px;
            width: 85px;
            color: var(--uvalib-text);
            text-align: center;
         }
      }
   }
</style>