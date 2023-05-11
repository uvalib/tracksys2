<template>
   <Panel header="Related Master Files" v-if="component.relatedMasterFiles.length > 0">
      <DataTable :value="component.relatedMasterFiles" ref="componentMasterFilesTable" dataKey="id"
         showGridlines stripedRows responsiveLayout="scroll" class="p-datatable-sm"
         :lazy="false" :paginator="true" :alwaysShowPaginator="true" :rows="15"
         :rowsPerPageOptions="[15,30,50,100]" paginatorPosition="top"
         paginatorTemplate="FirstPageLink PrevPageLink CurrentPageReport NextPageLink LastPageLink RowsPerPageDropdown"
         currentPageReportTemplate="{first} - {last} of {totalRecords}"
         v-model:selection="selectedMasterFiles" :selectAll="selectAll" @select-all-change="onSelectAllChange" @row-select="onRowSelect" @row-unselect="onRowUnselect"
      >
         <template #paginatorstart>
            <div class="master-file-acts" id="sticky-toolbar">
               <DPGButton label="PDF of Selected" @click="pdfClicked()" class="p-button-secondary" :disabled="!filesSelected" />
               <DPGButton label="Download Selected" @click="downloadClicked()" class="p-button-secondary" :disabled="!filesSelected" />
               <template  v-if="userStore.isAdmin || userStore.isSupervisor">
                  <LookupDialog :disabled="!filesSelected" label="Assign Metadata" @selected="assignMetadata" target="metadata" :create="true"/>
               </template>
            </div>
         </template>
         <Column selectionMode="multiple" headerStyle="width: 3em"></Column>
         <Column field="id" header="ID">
            <template #body="slotProps">
               <router-link :to="`/masterfiles/${slotProps.data.id}`">{{slotProps.data.id}}</router-link>
            </template>
         </Column>
         <Column field="pid" header="PID"/>
         <Column field="filename" header="File Name"/>
         <Column field="title" header="Title"/>
         <Column field="description" header="Description"/>
         <Column field="metadata.tiitle" header="Metadata">
            <template #body="slotProps">
               <router-link :to="`/metadata/${slotProps.data.metadata.id}`">{{slotProps.data.metadata.title}}</router-link>
            </template>
         </Column>
         <Column field="thumbnailURL" header="Thumb" class="thumb">
            <template #body="slotProps">
               <a :href="slotProps.data.viewerURL" target="_blank">
                  <img :src="slotProps.data.thumbnailURL" :class="{exemplar: slotProps.data.exemplar}"/>
               </a>
            </template>
         </Column>
      </DataTable>
   </Panel>
   <Panel header="Related Master Files" v-else>
      <p>No master files are associated with this component.</p>
   </Panel>
</template>

<script setup>
import { useComponentsStore } from '@/stores/components'
import { useUserStore } from '@/stores/user'
import { useSystemStore } from '@/stores/system'
import Panel from 'primevue/panel'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import { ref, onMounted, onUnmounted, computed } from 'vue'
import LookupDialog from '@/components/LookupDialog.vue'

const component = useComponentsStore()
const userStore = useUserStore()
const systemStore = useSystemStore()

const toolbarTop = ref(0)
const toolbarHeight = ref(0)
const toolbarWidth = ref(0)
const toolbar = ref(null)
const selectedMasterFiles = ref([])
const selectAll = ref(false)

const filesSelected = computed(() => {
   return selectedMasterFiles.value.length > 0
})
const selectedFileNames = computed(() => {
   let filenames = []
   selectedMasterFiles.value.forEach( s => {
      filenames.push(s.filename)
   })
   return filenames
})
const selectedIDs = computed(() => {
   let ids = []
   selectedMasterFiles.value.forEach( s => {
      ids.push(s.id)
   })
   return ids
})

onMounted(() => {
   setTimeout( () => {
      let tb = null
      let tbs = document.getElementsByClassName("p-paginator-top")
      if ( tbs ) {
         tb = tbs[0]
      }
      if ( tb) {
         toolbar.value = tb
         toolbarHeight.value = tb.offsetHeight
         toolbarWidth.value = tb.offsetWidth
         toolbarTop.value = 0

         // walk the parents of the toolbar and add each top value
         // to find the top of the toolbar relative to document top
         let ele = tb
         if (ele.offsetParent) {
            do {
               toolbarTop.value += ele.offsetTop
               ele = ele.offsetParent
            } while (ele)
         }
      }
   }, 1000)
   window.addEventListener("scroll", scrollHandler)
})

onUnmounted(() => {
   window.removeEventListener("scroll", scrollHandler)
})

function scrollHandler( ) {
   if ( toolbar.value) {
      if ( window.scrollY <= toolbarTop.value ) {
         if ( toolbar.value.classList.contains("sticky") ) {
            toolbar.value.classList.remove("sticky")
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].classList.add("sticky")
               dts[0].style.top = `0px`
            }
         }
      } else {
         if ( toolbar.value.classList.contains("sticky") == false ) {
            let dts = document.getElementsByClassName("p-datatable-wrapper")
            if ( dts ) {
               dts[0].classList.add("sticky")
               dts[0].style.top = `${toolbarHeight.value}px`
            }
            toolbar.value.classList.add("sticky")
            toolbar.value.style.width = `${toolbarWidth.value}px`
         }
      }
   }
}

function downloadClicked() {
   component.downloadFromArchive(userStore.computeID, selectedMasterFiles.value[0].unitID, selectedFileNames.value)
   clearSelections()
}
function pdfClicked() {
   let ids = []
   selectedMasterFiles.value.forEach( s => {
      ids.push(s.id)
   })
   let token = new Date().getTime()
   let unitID = selectedMasterFiles.value[0].unitID
   let mdPID = selectedMasterFiles.value[0].metadata.pid
   let url = `${systemStore.pdfURL}/${mdPID}?unit=${unitID}&token=${token}&pages=${ids.join(',')}`
   window.open(url)
}
async function assignMetadata( metadataID ) {
   await component.assignMetadata(metadataID, selectedMasterFiles.value[0].unitID, selectedIDs.value)
   clearSelections()
}
function onRowSelect() {
   selectAll.value = selectedMasterFiles.value < component.relatedMasterFiles.length
}
function onRowUnselect() {
   selectAll.value  = false
}
function onSelectAllChange(event) {
   selectAll.value = event.checked
   if (selectAll.value) {
      selectedMasterFiles.value = component.relatedMasterFiles
   }
   else {
      selectedMasterFiles.value = []
   }
}
function clearSelections() {
   selectAll.value = false
   selectedMasterFiles.value = []
}

</script>

<style scoped lang="scss">
:deep(.p-datatable-wrapper.sticky) {
   position: relative;
}
:deep(.p-paginator-top)  {
   div.p-paginator.p-component {
      padding: 0 0 15px 0;
   }
}
:deep(.p-paginator-top.sticky) {
   position: fixed;
   z-index: 1000;
   top: 0;
   div.p-paginator.p-component {
      padding: 0 0 10px 0;
      border-bottom: 2px solid var(--uvalib-grey-lightest);
      border-radius: 0;
   }
}

.p-datatable-sm {
   font-size: 0.9em;
}
.master-file-acts {
   font-size: 0.85em;
   text-align: right;
   button.p-button {
      margin-left: 5px;
   }
}
:deep(td.thumb) {
   width: 160px !important;
   text-align: center !important;
}
img.exemplar {
   border: 5px solid var(--uvalib-teal);
}

</style>