<template>
   <UModal :dismissible="false" :close="false" title="Submission Notes" @enter="modalOpened">
      <template #body>
        <div class="column gap-2">
            <div>{{ props.item.metadata.pid }} - {{ props.item.metadata.title }}</div>
            <UTextarea class="notes" v-if="edit" v-model="notes" :rows="10" />
            <div v-else class="note-text">{{ notes }}</div>
        </div>
      </template>
      <template #footer=>
         <template v-if="edit">
            <UButton label="Cancel" color="secondary" @click="emit('close', false)"/>
            <UButton label="Submit" color="primary" class="left-margin" @click="submitNoteEdit"/>
         </template>
         <template v-else>
            <UButton label="Close" color="secondary" class="left-margin" @click="emit('close', false)"/>
            <UButton label="Edit" color="primary" @click="edit=true"/>
         </template>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useArchivesSpaceStore } from '@/stores/archivesspace'

const archivesSpace = useArchivesSpaceStore()

const emit = defineEmits( ['close' ])
const props = defineProps({
   item: {
      type: Object,
      reqtured: true
   }
})

const notes = ref("")
const edit = ref(false)

const modalOpened = (() => {
   notes.value = props.item.notes
   edit.value = true
   if ( notes.value.length > 0 ) {
      edit.value = false
   }
})

const submitNoteEdit = ( async () => {
   await archivesSpace.updateNotes( props.item.metadataID, notes.value )
   emit('close', true)
})
</script>

<style scoped lang="scss">
div.note-text {
   height: 250px;
   overflow-y: scroll;
   margin-top: 15px;
   padding: 5px 10px;
   border: 1px solid var(--uvalib-grey-light);
   border-radius: 5px;
   background-color: white;
}
</style>