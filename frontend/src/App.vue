<template>

<ConfirmDialog position="top" :closable="false"/>

   <UApp :toaster="toaster">

      <div class="header" role="banner" id="uva-header">
         <div class="main-header">
            <div class="library-link">
               <a target="_blank" href="https://library.virginia.edu">
                  <UvaLibraryLogo />
               </a>
            </div>
            <div class="site-link">
               <router-link @click="homeClicked" to="/">Tracksys</router-link>
               <p class="version">v{{ systemStore.version }}</p>
            </div>
         </div>
         <MenuBar v-if="userStore.jwt" />
      </div>

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
import MenuBar from "@/components/MenuBar.vue"
import WaitSpinner from "@/components/WaitSpinner.vue"
import { useSystemStore } from "@/stores/system"
import { useUserStore } from "@/stores/user"
import { useSearchStore } from "@/stores/search"
import { onBeforeMount, watch, ref } from 'vue'

const systemStore = useSystemStore()
const userStore = useUserStore()
const searchStore = useSearchStore()
const toast = useToast()

const toaster = { duration: 5000, position: "top-center" }

const configuring = ref(true)

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
})

const errorClosed = (() => {
   systemStore.clearError()
})

onBeforeMount( async () => {
   document.title = `Tracksys`
   await systemStore.getConfig()
   configuring.value = false
})

</script>

<style scoped lang="scss">
div.header {
   background-color: var(--uvalib-brand-blue);
   color: white;
   text-align: left;
   position: relative;
   box-sizing: border-box;

   .main-header {
      display: flex;
      flex-direction: row;
      flex-wrap: nowrap;
      justify-content: space-between;
      align-content: stretch;
      align-items: center;
      padding: 1vw 20px 5px 10px;
      a {
         color: white !important;
      }
   }
}
// a {
//    color: var(--uvalib-blue-alt-dark);
//    font-weight: 500;
//    text-decoration: none;
//    &:hover {
//       text-decoration: underline;
//    }
// }
p.version {
   margin: 5px 0 0 0;
   font-size: 0.5em;
   text-align: right;
}
div.library-link {
   width: 220px;
   order: 0;
   flex: 0 1 auto;
   align-self: flex-start;
}
div.site-link {
   order: 0;
   font-size: 1.5em;
   a {
      color: white;
      text-decoration: none;
      &:hover {
         text-decoration: underline;
      }
   }
}
div.content {
   position: relative;
}
</style>
