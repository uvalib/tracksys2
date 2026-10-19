<template>
   <UModal :dismissible="false" :close="false" title="Reject Submission" @enter="modalOpened">
      <template #body>
        <div class="column gap-2!">
            <div>{{ props.item.metadata.pid }} - {{ props.item.metadata.title }}</div>
            <label>Please add some notes about the rejection:</label>
            <UTextarea v-model="reason" :rows="6" class="w-full"/>
        </div>
      </template>
      <template #footer=>
         <UButton label="Cancel" color="secondary" @click="emit('close', false)"/>
         <UButton label="Reject" color="primary" class="left-margin" :disabled="reason.length == 0"
            @click="rejectSubmitted"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref } from 'vue'
import { useArchivesSpaceStore } from '@/stores/archivesspace'

const archivesSpace = useArchivesSpaceStore()

const emit = defineEmits( ['close' ])
const props = defineProps({
   userID: {
      type: Number,
      reqtured: true
   },
   item: {
      type: Object,
      reqtured: true
   }
})

const reason = ref("")

const modalOpened = (() => {
   reason.value = ""
})

const rejectSubmitted = ( async () => {
   await archivesSpace.reject( props.userID, props.item.metadataID, reason.value )
   emit('close', true)
})
</script>

<style scoped lang="scss">
</style>