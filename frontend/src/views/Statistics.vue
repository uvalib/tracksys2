<template>
   <h2>Statistics</h2>
   <div class="date-range">
      <div class="date-picker">
         <label>From:</label>
         <select v-model="statsStore.dateRangeType" @change="modeChanged">
            <option value="before">BEFORE</option>
            <option value="after">AFTER</option>
            <option value="between">BETWEEN</option>
         </select>
         <DatePicker v-model="statsStore.startDate" showIcon :showOnFocus="false" dateFormat="yy-mm-dd"/>
         <DatePicker v-if="statsStore.dateRangeType == 'between'" v-model="statsStore.endDate" showIcon :showOnFocus="false" dateFormat="yy-mm-dd"/>
      </div>
      <DPGButton @click="getAllClicked" label="Get All Statistics"/>
   </div>
   <div class="stats">
      <div class="column">
         <StorageStats />
         <ImageStats />
         <ArchiveStats />
      </div>
      <div class="column">
            <MetadataStats />
      </div>
   </div>
   <div class="stats">
      <div class="column">
         <h3>Recent Virgo Publications</h3>
         <div  v-if="statsStore.publishedStats.loading" class="wait-wrap">
            <WaitSpinner/>
         </div>
         <div class="ext-system">
            <table>
               <tbody>
                  <tr>
                     <th></th><th>Title</th><th>Thumbnail</th><th>Details</th>
                  </tr>
                  <tr v-for="(rec,idx) in statsStore.publishedStats.virgo" :key="`as${rec.id}`">
                     <td class="num">{{idx+1}}.</td>
                     <td class="title">{{rec.title}}</td>
                     <td><img :src="rec.thumbURL"/></td>
                     <td><router-link :to="`/metadata/${rec.id}`">Details</router-link></td>
                  </tr>
               </tbody>
            </table>
         </div>
      </div>
      <div class="column">
         <h3>Recent ArchivesSpace Publications</h3>
         <div  v-if="statsStore.publishedStats.loading" class="wait-wrap">
            <WaitSpinner/>
         </div>
         <div v-else class="ext-system">
            <table>
               <tbody>
                  <tr>
                     <th></th><th>Title</th><th>Details</th><th>Link</th>
                  </tr>
                  <tr v-for="(rec,idx) in statsStore.publishedStats.archivesSpace" :key="`as${rec.id}`">
                     <td class="num">{{idx+1}}.</td>
                     <td class="title">{{rec.title}}</td>
                     <td><router-link :to="`/metadata/${rec.id}`">Details</router-link></td>
                     <td><a :href="rec.externalURL" target="_blank">ArchivesSpace</a></td>
                  </tr>
               </tbody>
            </table>
         </div>
      </div>
   </div>
</template>

<script setup>
import { onMounted } from 'vue'
import DatePicker from 'primevue/datepicker'
import {useStatsStore} from '@/stores/statistics'
import ImageStats from '@/components/stats/ImageStats.vue'
import StorageStats from '@/components/stats/StorageStats.vue'
import MetadataStats from '@/components/stats/MetadataStats.vue'
import ArchiveStats from '@/components/stats/ArchiveStats.vue'
import WaitSpinner from "@/components/WaitSpinner.vue"

const statsStore = useStatsStore()

onMounted( () => {
   statsStore.getAllStats(false)
})

const modeChanged = (() => {
   if ( statsStore.dateRangeType == "between") {
      statsStore.endDate = new Date()
      statsStore.endDate =  statsStore.endDate.setMonth(statsStore.startDate.getMonth() + 3)
   }
})

function getAllClicked() {
   statsStore.getAllStats(true)
}
</script>

<style scoped lang="scss">
.date-range {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   padding: 10px 15px;
   border-bottom: 1px solid var(--uvalib-grey-light);
   border-top: 1px solid var(--uvalib-grey-light);
   .date-picker {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: anchor-center;
      gap: 5px;
      .p-datepicker {
         width: 300px;
      }
   }
}
.stats {
   margin: 10px;
   display: flex;
   flex-flow: row wrap;
   text-align: left;
   gap: 15px;
   h3 {
      margin: 10px 0 5px 10px;
      padding-bottom: 5px;
      text-align: left;
      border-bottom: 1px solid var(--uvalib-grey-light);
   }
   .wait-wrap {
      padding: 20px 10px;
   }
   .column {
      width: 48%;
   }
   table {
      margin: 10px 0 0 10px;
      border-collapse: collapse;
      border: 1px solid #dedede;
      box-shadow: var(--box-shadow-light);
      th {
         background-color: #efefef;
         text-align: left;
         padding: 4px 10px 4px 5px;
         border-bottom: 1px solid #ccc;
      }
      td {
         vertical-align: middle;
         background: white;
         border-bottom: 1px solid #dedede;
         padding: 4px 10px 4px 5px;
      }
   }
}
</style>
