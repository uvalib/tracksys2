<template>
   <UModal v-model:open="isOpen" :modal="true" :dismissible="false" :close="false" title="Manage Tags">
      <UButton @click="show()" size="sm" color="secondary" label="Manage Tags"/>
      <template #body>
         <UListbox :items="tags" virtualize valueKey="id" labelKey="tag" @update:modelValue="addTag"
            :filter="{placeholder: 'Search...', icon: 'i-lucide-search'}"
         />
          <div class="row-left gap-2 pt-4">
            <UInput v-model="newTag" placeholder="New Tag" class="grow"/>
            <UButton label="Add" color="secondary" @click="createTag()"/>
         </div>
         <div class="pt-4">
            <label>Current Tags:</label>
            <div v-if=" masterFiles.details.tags.length == 0" class="p-4">None</div>
            <div v-else class="flex flex-row flex-wrap gap-2 pl-4 pt-2">
               <UButton v-for="t in masterFiles.details.tags" size="sm"  color="secondary" trailing-icon="i-lucide-circle-x"
                  class="rounded-full" :label="t.tag"  :key="t.id" @click="removeTag(t)"/>
            </div>
         </div>
      <p class="error" v-if="tagStore.error">{{tagStore.error}}</p>
      </template>
      <template #footer="{close}">
         <UButton @click="close" color="secondary" label="Done"/>  
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useMasterFilesStore } from '@/stores/masterfiles'
import { useTagsStore } from '@/stores/tags'

const masterFiles = useMasterFilesStore()
const tagStore = useTagsStore()

const isOpen = ref(false)
const newTag = ref("")

const tags = computed(() => {
   return tagStore.tags.filter( t => isUsed(t) == false)
})

const createTag = ( async () => {
   tagStore.setError("")
   let cleanVal = newTag.value.trim()
   if (cleanVal == "") {
      tagStore.setError("Please enter a new tag name")
      return
   }

   let idx = tagStore.tags.findIndex( t => t.tag.trim() == cleanVal)
   if (idx > -1) {
      tagStore.setError(`Tag '${cleanVal}' already exists. Please create a unique tag.`)
   }
   await tagStore.createTag(cleanVal)
   newTag.value = ""

   let tag = tagStore.tags.find( t => t.tag.trim() == cleanVal)
   masterFiles.addTag(tag)
})

const removeTag = (( tag ) => {
   masterFiles.removeTag(tag)
})
const addTag =( (tagID) => {
   masterFiles.addTag(tagStore.tags.find( t => t.id == tagID))
})
const isUsed = ( ( tgt ) => {
   let idx = masterFiles.details.tags.findIndex( t => t.tag == tgt.tag )
   return idx > -1
})
const show = (() => {
   tagStore.getTags()
   isOpen.value = true
})
</script>

<style lang="scss" scoped>
p.error {
   text-align: center;
   color: var(--uvalib-red-emergency);
   padding: 0;
   margin: 10px 0 0 0;
}

</style>
