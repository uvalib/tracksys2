<template>
   <h2>Statistics</h2>
   <div class="date-range">
      <div class="row-left items-center">
         <label>From:</label>
         <USelect v-model="rangeType" @change="modeChanged" :items="rangeTypes" />
         <UInputDate v-if="rangeType != 'between'" v-model="dates.start" variant="soft"/>
         <UInputDate v-else range v-model="dates" variant="soft" separator-icon="i-lucide-arrow-right"/>
      </div>
      <UButton @click="getAllClicked" label="Get All Statistics"/>
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
import { onMounted, shallowRef, ref } from 'vue'
import { today, getLocalTimeZone } from '@internationalized/date'
import {useStatsStore} from '@/stores/statistics'
import ImageStats from '@/components/stats/ImageStats.vue'
import StorageStats from '@/components/stats/StorageStats.vue'
import MetadataStats from '@/components/stats/MetadataStats.vue'
import ArchiveStats from '@/components/stats/ArchiveStats.vue'
import WaitSpinner from "@/components/WaitSpinner.vue"

const statsStore = useStatsStore()

const rangeTypes = [
   {label: "BEFORE", value: "before"},
   {label: "AFTER", value: "after"},
   {label: "BETWEEN", value: "between"}
]

const rangeType = ref("before")
const dates = shallowRef({
  start: today( getLocalTimeZone() ),
  end: today( getLocalTimeZone() ).add({months: 3}) 
})

onMounted( () => {
   statsStore.getAllStats(false, rangeType.value, dates.value.start.toString(), dates.value.end.toString())
})

const modeChanged = (() => {
   console.log("MODE CHANGED")
   if ( rangeType.value == "between") {
      dates.value.end = dates.value.start.add({months: 3}) 
   } else {
      dates.value.end = null
   }
})

function getAllClicked() {
   statsStore.getAllStats(true, rangeType.value, dates.value.start.toString(), dates.value.end.toString()) 
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
