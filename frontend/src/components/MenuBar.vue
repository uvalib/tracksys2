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
         <IconField iconPosition="left" v-if="showSearch" class="global-search">
            <InputIcon class="pi pi-search"> </InputIcon>
            <InputText v-model="newQuery" placeholder="Search TrackSys" @keyup.enter="searchEntered" />
         </IconField>
      </template>
   </Menubar>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Menubar from 'primevue/menubar'
import { useSystemStore } from "@/stores/system"
import { useSearchStore } from '@/stores/search'
import { useUserStore } from "@/stores/user"
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
import InputText from 'primevue/inputtext'
import { useRoute, useRouter } from 'vue-router'

const systemStore = useSystemStore()
const searchStore = useSearchStore()
const userStore = useUserStore()
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
         {label: "Home", command: ()=>menuLinkClicked("/")},
         {label: "Orders", command: ()=>menuLinkClicked("/orders")},
         {label: "Collections", command: ()=>menuLinkClicked("/collections")},
         {label: "Published", items: [
            {label: "Virgo", command: ()=>menuLinkClicked("/published/virgo")},
            {label: "ArchivesSpace", command: ()=>menuLinkClicked("/published/archivesspace")},
            {label: "DPLA", command: ()=>menuLinkClicked("/published/dpla")},
         ]},
         {label: "Job Statuses", command: ()=>menuLinkClicked("/jobs")},
         {label: "Digitization", items: [
            {label: "Equipment", command: ()=>menuLinkClicked("/equipment")},
            {label: "Projects", url: systemStore.projectsURL, target: "_blank"},
            {label: "Reports", url: `${systemStore.reportsURL}/reports`, target: "_blank"},
            {label: "Statistics", url: systemStore.reportsURL, target: "_blank"},
         ]},
         {label: "Miscellaneous", items: [
            {label: "APTrust Submissions", command: ()=>menuLinkClicked("/aptrust")},
            {label: "ArchivesSpace Reviews", command: ()=>menuLinkClicked("/archivesspace")},
            {label: "HathiTrust Submissions", command: ()=>menuLinkClicked("/hathitrust")},
            {label: "Customers", command: ()=>menuLinkClicked("/customers")},
            {label: "Staff Members", command: ()=>menuLinkClicked("/staff")},
         ]},
         {label: userStore.signedInUser, items: [
            {label: "Sign Out", command: ()=>signOut() },
         ]}
      ]
   }, 500)
})

const menuLinkClicked = ( (tgtRoute) => {
   router.push(tgtRoute)
})

const signOut = (() => {
   userStore.signout()
   router.push("/signedout")
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
.global-search {
   margin:0;
   input[type=text] {
      margin-bottom: 0 !important;
   }
}
</style>