<template>
   <DPGButton @click="show" severity="secondary" label="Manage Tags"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Master File Tags" :style="{width: '450px'}">
      <IconField iconPosition="left">
         <InputIcon class="pi pi-search" />
         <InputText v-model="tagStore.query" placeholder="Search" @input="tagStore.getTags()" autofocus/>
      </IconField>
      <VirtualScroller :items="tagStore.tags" :itemSize="22" showLoader class="taglist" :showLoader="tagStore.loading" >
         <template v-slot:item="{ item }">
            <div v-if="isUsed(item)==false" class="tag-list-item" @click="addTag(item)">{{ truncateText(item.tag) }}</div>
         </template>
      </VirtualScroller>
      <div class="add">
         <InputText v-model="newTag" placeholder="New Tag"/>
         <DPGButton label="Add" severity="secondary" @click="createTag()"/>
      </div>
      <div class="selected">
         <label>Current Tags:</label>
         <div v-if=" masterFiles.details.tags.length == 0" class="none">None</div>
         <div v-else class="cur-tags">
            <Chip v-for="t in masterFiles.details.tags" :label="t.tag" removable :key="t.id" @remove="removeTag(t)"/>
         </div>
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
import IconField from 'primevue/iconfield'
import InputIcon from 'primevue/inputicon'
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
   margin-top: 10px;

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
      margin-right: 5px;
   }
}

.selected {
   padding: 0;
   margin: 10px 0 0 0;
   .none {
      margin: 10px 0 0 15px;
   }
   .cur-tags {
      display: flex;
      flex-flow: row wrap;
      justify-content: flex-start;
      margin-top: 10px;
      .p-chip {
         font-size: 0.8em;
         margin: 2px;
         background-color: #f1f5f9;
      }
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
