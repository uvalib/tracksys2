<template>
   <UModal v-model:open="showDialog" :modal="true" :dismissible="false" :close="false" title="Submit HathiTrust Metadata">
      <UButton @click="showDialog = true" label="Submit HathiTrust Metadata" />
      <template #body>
         <div class="hathi-panel">
            <p>Submit all metadata for candidate units in this order</p>
            <div>
               <label>Submission Mode</label>
               <USelect v-model="submitMode" placeholder="Select a submission mode" :items="modes" />
               <p class="hint">Development mode will log the metadata to the job log</p>
            </div>
            <div>
               <label>Submission Name</label>
               <input type="text" v-model="submitName" />
               <p class="hint">This is an identifier that will appended to the submission file name to help identify it later. For example: batch12</p>
            </div>
         </div>
      </template>
      <template #footer="{ close }">
         <UButton label="Cancel" color="secondary" @click="close"/>
         <UButton label="Submit"  @click="submitClicked" :disabled="submitDisabled"/>
      </template>
   </UModal>
</template>

<script setup>
import { ref, computed } from 'vue'

const emit = defineEmits( ['submit' ])
const props = defineProps({
   order: {
      type: Number,
      required: true
   }
})

const showDialog = ref(false)
const submitMode = ref("")
const submitName = ref(`order${props.order}`)

const modes = [
   {value: "dev", label: "Development (no submission)"},
   {value: "prod", label: "Production"}
]

const submitDisabled = computed( () => {
   return ( submitMode.value == "" || submitName.value == "")
})

const submitClicked = (() => {
   emit("submit", {mode: submitMode.value, name: submitName.value})
   showDialog.value = false
})

</script>

<style lang="scss" scoped>
.hathi-panel {
   display: flex;
   flex-direction: column;
   gap: 15px;;

   p {
      margin:0;
   }
   p.hint {
      font-size: 0.85em;
      color: var(--uvalib-grey);
      margin: 0;
   }
   label {
      display: block;
      font-weight: bold;
      margin-bottom: 5px;
   }
   .buttons {
      text-align: right;
      margin: 0 0 5px 0;
      button {
         margin-left: 10px;
      }
   }
}
</style>