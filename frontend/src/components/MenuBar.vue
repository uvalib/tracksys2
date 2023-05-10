<template>
   <Menubar :model="items">
      <template #end>
         <span class="global-search p-input-icon-right" v-if="showSearch">
            <i class="pi pi-search" />
            <InputText v-model="newQuery" @keyup.enter="searchEntered" placeholder="Search Tracksys..."/>
         </span>
      </template>
   </Menubar>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Menubar from 'primevue/menubar'
import { useSystemStore } from "@/stores/system"
import { useSearchStore } from '@/stores/search'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'

const systemStore = useSystemStore()
const searchStore = useSearchStore()
const route = useRoute()
const router = useRouter()

const items = ref([])
const newQuery = ref("")

const showSearch = computed(() => {
   return route.name != "home"
})

onMounted(() => {
   setTimeout( () => {
      items.value = [
         {label: "Home", to: "/"},
         {label: "Orders", to: "/orders"},
         {label: "Collections", to: "/collections"},
         {label: "Job Statuses", to: "/jobs"},
         {label: "Digitization", items: [
            {label: "Equipment", to: '/equipment'},
            {label: "Projects", url: systemStore.projectsURL, target: "_blank"},
            {label: "Reports", url: `${systemStore.reportsURL}/reports`, target: "_blank"},
            {label: "Statistics", url: systemStore.reportsURL, target: "_blank"},
         ]},
         {label: "Miscellaneous", items: [
            {label: "Staff Members", to: "/staff"},
            {label: "Customers", to: "/customers"},
         ]}
      ]
   }, 500)
})

const searchEntered = (() => {
   if (newQuery.value.length > 0) {
      searchStore.resetSearch()
      let query = Object.assign({}, route.query)
      query.q = newQuery.value
      router.push({path: "/", query: query})
      newQuery.value = ""
   }
})
</script>

<style scoped lang="scss">
.p-menubar {
   padding: 0;
   border-radius: 0;
   min-height: 48px;
}
.global-search {
   margin: 5px;
   width: 300px;
}
</style>