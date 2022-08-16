<template>
   <Menubar :model="items">
   <template #end>
      <span class="signout" tabindex="0" @click="signOut">Sign Out</span>
   </template>
   </Menubar>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import Menubar from 'primevue/menubar'
import { useSystemStore } from "@/stores/system"
import { useUserStore } from "@/stores/user"
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const systemStore = useSystemStore()
const router = useRouter()

const items = ref([])

function signOut() {
   userStore.signout()
   router.push("signedout")
}

onMounted(() => {
   setTimeout( () => {
      items.value = [
         {label: "Dashboard", to: "/"},
         {label: "Job Statuses", to: "/jobs"},
         {label: "Digitization", items: [
            {label: "Projects", url: systemStore.projectsURL, target: "_blank"},
            {label: "Reports", url: systemStore.reportsURL, target: "_blank"},
         ]}
      ]
   }, 50)
})
</script>

<style scoped lang="scss">
.p-menubar {
   padding: 0;
   border-radius: 0;
}
.signout {
   cursor: pointer;
   padding: 0.5rem;
   margin: 0 10px 0 0;
   &:hover {
      text-decoration: underline;
      background: #e9ecef;
   }
}
</style>