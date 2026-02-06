<template>
   <div class="stats-card">
      <h3>Metadata Statistics</h3>
      <div  v-if="statsStore.metadataStats.loading" class="wait-wrap">
         <WaitSpinner/>
      </div>
      <div v-else class="stats">
         <dl>
            <dt>Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.total)}}</dd>
            <dt>SIRSI Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.sirsi)}}</dd>
            <dt>XML Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.xml)}}</dd>
            <dt>DL Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.totalDL)}}</dd>
            <dt>DL SIRSI Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.sirsiDL)}}</dd>
            <dt>DL XML Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.xmlDL)}}</dd>
             <dt>DPLA Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.totalDPLA)}}</dd>
            <dt>DPLA SIRSI Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.sirsiDPLA)}}</dd>
            <dt>DPLA XML Metadata Count:</dt>
            <dd>{{numberWithCommas(statsStore.metadataStats.xmlDPLA)}}</dd>
         </dl>
      </div>
      <p class="error" v-if="statsStore.metadataStats.error">{{statsStore.metadataStats.error}}</p>
   </div>
</template>

<script setup>
import {useStatsStore} from '@/stores/statistics'
import WaitSpinner from "@/components/WaitSpinner.vue"

const statsStore = useStatsStore()

function numberWithCommas(num) {
    return num.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
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

