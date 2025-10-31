<template>
   <h2>
      <span>Job {{route.params.id}} Processing Log</span>
      <div class="actions">
         <DPGButton v-if="jobsStore.details.status=='running'"
            label="Watch status" @click="tailLog" :disabled="tailIntervalID > -1" />
         <DPGButton v-else label="Jump to end" @click="jumpToEnd" />
         <DPGButton label="Delete" @click="deleteJobLog" />
      </div>
   </h2>
   <div class="status">
      <template v-if="jobsStore.details.error">
         <b class="error">FAILED:</b><span>{{jobsStore.details.error}}</span>
      </template>
      <b class="finished" v-else-if="jobsStore.details.status=='finished'">FINISHED</b>
       <b class="running" v-else>RUNNING...</b>
       <span>
         <label>Associated Object:</label>
         <router-link v-if="getAssociatedObjectLink(jobsStore.details.associatedObject)" :to="getAssociatedObjectLink(jobsStore.details.associatedObject)">
            {{jobsStore.details.associatedObject}}
         </router-link>
         <span v-else></span>
       </span>
   </div>
   <div class="log">
      <div class="scroller">
         <div class="line" v-for="le in jobsStore.details.events" :key="le.id">
            <span class="date">{{ formatTimeStamp(le.timeStamp) }}</span>
            <span class="sep">:</span>
            <span :class="le.level">{{le.level.toUpperCase()}}</span>
            <span class="sep">:</span>
            <span class="txt" v-html="le.text"></span>
         </div>
      </div>
   </div>
</template>

<script setup>
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { useJobsStore } from '@/stores/jobs'
import { useRoute, useRouter } from 'vue-router'
import { useConfirm } from "primevue/useconfirm"
import { useDateFormat } from '@vueuse/core'

const route = useRoute()
const router = useRouter()
const jobsStore = useJobsStore()
const confirm = useConfirm()

const tailIntervalID = ref(-1)

onMounted(() => {
   jobsStore.getJobDetails(route.params.id)
   document.title = `Job #${route.params.id}`
   tailIntervalID.value = -1
})

onBeforeUnmount(() => {
   if ( tailIntervalID.value > -1 ) {
      clearInterval( tailIntervalID.value )
   }
})

const formatTimeStamp = ( (ts) => {
   return useDateFormat(ts, "YYYY-MM-DD HH:mm:ss").value
})

const getAssociatedObjectLink = (( objName ) => {
   if (objName.split(" ").length != 2) {
      return ""
   }
   let objType = objName.split(" ")[0].toLowerCase().trim()
   let objID =  objName.split(" ")[1].toLowerCase().trim()
   if (objType == "unit") {
      return `/units/${objID}`
   }
   if (objType == "order") {
      return `/orders/${objID}`
   }
   if (objType == "metadata") {
      return `/metadata/${objID}`
   }
   return ""
})

const tailLog = (() => {
   tailIntervalID.value = setInterval( ()=> {
      jobsStore.getJobDetails(route.params.id)
      jumpToEnd()
   }, 2000)
})

const jumpToEnd = (() => {
   let lines = document.getElementsByClassName("line")
   lines[lines.length-1].scrollIntoView()
})

const deleteJobLog = (() => {
   confirm.require({
      message: 'Are you sure you want delete this job log? All data will be lost. This cannot be reversed.',
      header: 'Confirm Delete Job Log',
      icon: 'pi pi-exclamation-triangle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: async () => {
         await jobsStore.deleteJobs([ parseInt(route.params.id,10) ])
         router.push("/jobs")
      }
   })
})
</script>

<style scoped lang="scss">
   .status {
      padding: 15px 25px 10px 25px;
      text-align: left;
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
      label {
         font-weight: bold;
         display: inline-block;
         margin-right: 10px;
      }
      b {
         display: inline-block;
         margin-right: 10px;
      }
      b.error {
         color: firebrick;
      }
      b.finished {
         color: var( --uvalib-green-dark);
      }
      b.running {
         color: #629bff;
      }
   }
   .log {
      min-height: 600px;
      padding-bottom: 25px;
      .scroller {
         border-radius: 5px;;
         margin: 0 25px 25px 25px;
         background: #333;
         text-align: left;
         font-family: "Courier New", Courier, monospace;
         color: #ccc;
         padding: 15px;
         span {
            color: #f5f5f5;
            font-weight: bold;
         }
         .txt {
            font-weight: normal;
         }
         span.sep {
            display: inline-block;
            margin: 0 10px;
         }
         .info {
            color: #629bff;
         }
         .error {
            color: #CB9B43;
         }
         .fatal {
            color: #EE4444;
         }
      }
   }
</style>