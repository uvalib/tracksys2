<template>
   <h2>Master File Audit Report</h2>
   <div  v-if="auditStore.loading" class="wait-wrap">
      <WaitSpinner/>
   </div>
   <div v-else class="report">
      <div class="control-bar">
         <label>Year:</label>
         <select v-model="auditStore.targetYear">
            <option v-for="w in auditStore.auditYears" :value="w.value" :key="`wf${w.value}`">{{w.label}}</option>
         </select>
         <DPGButton severity="secondary" @click="auditStore.getAuditReport()" label="Generate Report"/>
      </div>


      <Chart type="bar" :data="auditStore" :options="options" style="max-height:800px;      "/>

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
import Chart from 'primevue/chart'

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
   .control-bar {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-end;
      align-items: anchor-center;
      gap: 10px;
      select {
         width: 100px;
      }
   }
   .total {
      text-align: center;
      margin: 20px 0;
   }
}
</style>