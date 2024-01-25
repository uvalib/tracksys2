<template>
   <Menubar :model="items">
      <template #item="{ label, item, props }">
         <router-link v-if="item.route" :to="item.route">
            {{ label }}
         </router-link>
         <a v-else :href="item.url" :target="item.target" v-bind="props.action">
            <span v-bind="props.label">{{ label }}</span>
            <span v-if="item.items" class="pi pi-fw pi-angle-down" v-bind="props.submenuicon" />
         </a>
      </template>
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
         {label: "Home", route: "/"},
         {label: "Orders", route: "/orders"},
         {label: "Collections", route: "/collections"},
         {label: "Job Statuses", route: "/jobs"},
         {label: "Digitization", items: [
            {label: "Equipment", route: '/equipment'},
            {label: "Projects", url: systemStore.projectsURL, target: "_blank"},
            {label: "Reports", url: `${systemStore.reportsURL}/reports`, target: "_blank"},
            {label: "Statistics", url: systemStore.reportsURL, target: "_blank"},
         ]},
         {label: "Miscellaneous", items: [
            {label: "APTrust Submissions", route: "/aptrust"},
            {label: "ArchivesSpace Reviews", route: "/archivesspace"},
            {label: "Hathirust Submissions", route: "/hathitrust"},
            {label: "Customers", route: "/customers"},
            {label: "Staff Members", route: "/staff"},
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
   li.p-menuitem {
      div.p-menuitem-content a {
         color: #495057 !important;
         padding: 0.75rem 1rem !important;
         display: block;
         border-radius: 0;
         white-space: nowrap;
         &:hover {
            text-decoration: none !important;
         }
      }
   }
}
.global-search {
   margin: 5px;
   width: 300px;
}
</style>