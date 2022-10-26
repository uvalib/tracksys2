<template>
   <DPGButton @click="show" class="p-button-secondary" label="Manage Tags"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Master File Tags" :style="{width: '450px'}">
      <div class="toolbar">
         <span class="p-input-icon-right">
            <i class="pi pi-search" />
            <InputText v-model="tagStore.query" placeholder="Tag Search" @input="tagStore.getTags()" autofocus/>
         </span>
         <DPGButton label="Clear" class="p-button-secondary" @click="clearSearch()"/>
      </div>
      <VirtualScroller :items="tagStore.tags" :itemSize="22" showLoader class="taglist" :loading="tagStore.loading" >
         <template v-slot:item="{ item }">
            <div v-if="isUsed(item)==false" class="tag-list-item" @click="addTag(item)">{{ truncateText(item.tag) }}</div>
         </template>
      </VirtualScroller>
      <div class="add">
         <InputText v-model="newTag" placeholder="New Tag"/>
         <DPGButton label="Add New" class="p-button-secondary" @click="createTag()"/>
      </div>
      <div class="selected">
         <Chip v-for="t in masterFiles.details.tags" :label="t.tag" removable :key="t.id" @remove="removeTag(t)"/>
      </div>
      <p class="error" v-if="tagStore.error">{{tagStore.error}}</p>
      <div class="acts">
         <DPGButton @click="hide" label="Done"/>
      </div>
   </Dialog>
</template>

<script setup>
import { ref } from 'vue'
import Dialog from 'primevue/dialog'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useTagsStore } from '@/stores/tags'
import VirtualScroller from 'primevue/virtualscroller'
import InputText from 'primevue/inputtext'
import Chip from 'primevue/chip'

const masterFiles = useMasterFilesStore()
const tagStore = useTagsStore()

const isOpen = ref(false)
const newTag = ref("")

async function createTag() {
   tagStore.setError("")
   let cleanVal = newTag.value.trim()
   if (cleanVal == "") {
      tagStore.setError("Please enter a new tag name")
   }
   let idx = tagStore.tags.findIndex( t => t.tag.trim() == cleanVal)
   if (idx > -1) {
      tagStore.setError(`Tag '${cleanVal}' already exists. Please create a unique tag.`)
   }
   await tagStore.createTag(cleanVal)
   newTag.value = ""

   let tag = tagStore.tags.find( t => t.tag.trim() == cleanVal)
   masterFiles.addTag(tag)
}
function removeTag( tag ) {
   masterFiles.removeTag(tag)
}
function addTag( tag ) {
   masterFiles.addTag(tag)
}
function clearSearch() {
   tagStore.query = ""
   tagStore.getTags()
}
function truncateText(t) {
   if (t.length < 50) return t
   return t.slice(0,47)+"..."
}
function isUsed( tgt ) {
   let idx = masterFiles.details.tags.findIndex( t => t.tag == tgt.tag )
   return idx > -1
}
function hide() {
   isOpen.value=false
}
function show() {
   tagStore.getTags()
   isOpen.value = true

}
</script>

<style lang="scss" scoped>
.tag-list-item {
   font-size: 0.8em;
   height: 22px;
   cursor: pointer;
   padding: 5px;
   &:hover {
      background: var(--uvalib-grey-lightest);
   }
}
.taglist {
   height: 150px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 3px;

   .tag-list-item.disabled {
      color: var(--uvalib-grey-light);
      cursor: default;
      &:hover {
         cursor: default;
         background: none;
      }
   }
}

.add {
   padding: 10px 0;
   text-align: right;
   display: flex;
   flex-flow: row nowrap;
   input {
      flex-grow: 1;
      font-size: 0.8em;
   }
}


.p-button  {
      margin-left: 5px;
      font-size: 0.8em;
   }

.toolbar {
   padding: 10px 0;
   text-align: right;
   display: flex;
   flex-flow: row nowrap;
   .p-input-icon-right {
      flex-grow: 1;
      input {
         width: 100%;
         font-size: 0.8em;
      }
   }
}

.selected {
   display: flex;
   flex-flow: row wrap;
   justify-content: flex-start;
   padding: 10px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 3px;
   margin-top: 15px;
   .p-chip {
      font-size: 0.8em;
      margin: 2px 4px;
   }
}
p.error {
   text-align: center;
   color: var(--uvalib-red-emergency);
   padding: 0;
   margin: 10px 0 0 0;
}

.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 15px 0 10px 0;
   margin: 0;
}
</style>
