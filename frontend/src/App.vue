<template>

<ConfirmDialog position="top" :closable="false"/> <!-- FIXME use composable nuxt -->

   <UApp :toaster="toaster">

      <UHeader mode="slideover" id="uva-header" title="UVA Library" to="https://library.virginia.edu">
         <template #title>
            <div class="library-link">
               <UvaLibraryLogo />
            </div>
         </template>

         <!-- this is the main menu. shows up in the center of the header if size allows -->
         <UNavigationMenu v-if="userStore.isSignedIn" highlight content-orientation="vertical" :items="menuItems" />

          <template #right>
            <div class="site-link">
               <RouterLink to="/">Tracksys</RouterLink>
               <p class="version">{{ systemStore.version }}</p>
            </div>
         </template>
      </UHeader>

      <UMain>
         <div class="content" v-if="configuring==false">
            <router-view />
         </div>
      </UMain>

      <WaitSpinner v-if="systemStore.working" :overlay="true" message="Please wait..." />

      <UModal v-model:open="systemStore.showError" :modal="true" :dismissible="false" title="System Error">
         <template #body>
            <div style="text-align: left" v-html="systemStore.error"></div>
         </template>
      </UModal>
   </UApp>
</template>

<script setup>
import UvaLibraryLogo from "@/components/UvaLibraryLogo.vue"
import WaitSpinner from "@/components/WaitSpinner.vue"
import { useSystemStore } from "@/stores/system"
import { useUserStore } from "@/stores/user"
import { useSearchStore } from "@/stores/search"
import { onBeforeMount, watch, ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const systemStore = useSystemStore()
const userStore = useUserStore()
const searchStore = useSearchStore()
const toast = useToast()
const router = useRouter()

const toaster = { duration: 5000, position: "top-center" }

const configuring = ref(true)
const menuItems = computed(()=> {
   return [ 
      {label: "Home", onSelect: () => homeClicked()}, 
      {label: "Orders", to: "/orders"},  
      {label: "Collections", to: "/collections"}, 
      {label: "Published", children: [
         {label: "Virgo", to: "/published/virgo"},    
         {label: "ArchivesSpace", to: "/published/archivesspace"},
         {label: "DPLA", to: "/published/depla"}
      ]},
      {label: "Job Statuses", to: "/jobs"},
      {label: "Digitization", children: [
         {label: "Equipment", to: `${systemStore.projectsURL}/equipment`, target: "_blank"},
         {label: "Projects", to: `${systemStore.projectsURL}`, target: "_blank"},
         {label: "Reports", to: `${systemStore.projectsURL}/reports`, target: "_blank"},
         {label: "Statistics", to: "/statistics"},
         {label: "Patron Deliveries", to: "/deliveries"},
      ]},
      {label: "Miscellaneous", children: [
         {label: "ArchivesSpace Reviews", to: "/archivesspace"},
         {label: "HathiTrust Submissions", to: "/hathitrust"},
         {label: "Master File Audit", to: "/audit-report"},
         {label: "Customers", to: "/customers"},
         {label: "Staff Members", to: "/staff"},
      ]},
      {label: userStore.signedInUser, children: [
         {label: "Sign Out", icon: 'i-lucide-log-out',  onSelect: () => signout()}    
      ]},
   ]
})

watch(() => systemStore.toast.show, (newShow) => {
   if ( newShow == true) {
      let color = "primary"
      if ( systemStore.toast.error) {
        color = "error"
      }
      toast.add({
         title: systemStore.toast.summary,
         description: systemStore.toast.message,
         color: color
      })
      systemStore.clearToastMessage()
   }
})

const homeClicked = (() => {
   searchStore.resetSearch()
   router.push("/")
})
const signOut = (() => {
   userStore.signout()
   router.push("/signedout")
})

onBeforeMount( async () => {
   document.title = `Tracksys`
   await systemStore.getConfig()
   configuring.value = false
})

</script>

<style scoped lang="scss">
div.library-link {
   width: 220px;
}
div.site-link {
   font-size: 1.5em;
   a {
      color: white !important;
      padding: 3px 6px;
      border-radius: 0.3rem;
      &:hover {
         background: var(--uvalib-blue-alt);
         text-decoration: none !important;
      }
   }
   p.version {
      margin: 0;
      font-size: 0.5em;
      text-align: right;
      padding: 0;
   }
}

</style>
