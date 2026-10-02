<template>
   <h2>Master File Audit Report</h2>
   <div  v-if="auditStore.loading" class="wait-wrap">
      <WaitSpinner/>
   </div>
   <div v-else class="report">
      <div class="row-right items-center">
         <label>Year:</label>
         <USelect v-model="auditStore.targetYear" :items="auditStore.auditYears"/>
         <UButton color="secondary" @click="auditStore.getAuditReport()" label="Generate Report"/>
      </div>

      <BarChart :data="auditStore" :options="options" style="max-height:800px;" />

      <div class="total">
         <label>Total Audited:</label><span class="total">{{auditStore.totalAudited}}</span>
      </div>
      <p class="error" v-if="auditStore.error">{{auditStore.error}}</p>
   </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import WaitSpinner from "@/components/WaitSpinner.vue"
import { useAuditStore } from '@/stores/audit'
import BarChart from "./BarChart.vue"

const auditStore = useAuditStore()

const options = ref({
   title: {
      display: false,
   },
   legend: {
      display: false
   },
   plugins: {
      legend: {
         display: false,
      },
      colors: {
         enabled: false
      }
   },
})

onMounted( () => {
   auditStore.getAuditReport()
})
</script>

<style scoped lang="scss">
.wait-wrap {
   text-align: center;
   margin-top: 10%;
}
.report {
   margin: 10px 50px;
   .total {
      text-align: center;
      margin: 20px 0;
   }
}
</style>