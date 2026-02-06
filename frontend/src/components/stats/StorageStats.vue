<template>
   <div class="stats-card">
      <h3>Storage Statistics</h3>
      <div  v-if="statsStore.storageStats.loading" class="wait-wrap">
         <WaitSpinner/>
      </div>
      <div v-else class="stats">
         <dl>
            <dt>Total Image Size (GB):</dt>
            <dd>{{numberWithCommas(statsStore.storageStats.total)}}</dd>
            <dt>DL Total Image Size (GB):</dt>
            <dd>{{numberWithCommas(statsStore.storageStats.DL)}}</dd>
         </dl>
      </div>
      <p class="error" v-if="statsStore.storageStats.error">{{statsStore.storageStats.error}}</p>
   </div>
</template>

<script setup>
import {useStatsStore} from '@/stores/statistics'
import WaitSpinner from "@/components/WaitSpinner.vue"

const statsStore = useStatsStore()

function numberWithCommas(num) {
    return Math.round(num).toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}
</script>

<style lang="scss" scoped>
.stats-card {
   margin: 10px;
   text-align: left;
   border: 1px solid var(--uvalib-grey-light);
   box-shadow: var(--box-shadow-light);
   h3 {
      background: var(--uvalib-grey-lightest);
      font-size: 1em;
      text-align: left;
      margin:0;
      padding: 5px 10px;
      border-bottom: 1px solid var(--uvalib-grey-light);
      font-weight: 500;;
   }
   .wait-wrap {
      text-align: center;
      padding: 30px 0 ;
   }
   .stats {
      padding: 10px;

      dl {
         margin: 10px 30px 0 30px;
         display: inline-grid;
         grid-template-columns: max-content 2fr;
         grid-column-gap: 10px;
         font-size: 0.9em;
         text-align: left;
         box-sizing: border-box;

         dt {
            font-weight: bold;
            text-align: right;
         }
         dd {
            margin: 0 0 10px 0;
            word-break: break-word;
            -webkit-hyphens: auto;
            -moz-hyphens: auto;
            hyphens: auto;
         }
      }
   }
}
</style>

