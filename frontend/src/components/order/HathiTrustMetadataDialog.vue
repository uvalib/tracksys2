<template>
   <DPGButton label="Submit HathiTrust Metadata" @click="showDialog = true"/>
   <Dialog v-model:visible="showDialog" :modal="true" header="Submit HathiTrust Metadata" position="top" >
      <div class="hathi-panel">
         <p>Submit all metadata for candidate units in this order</p>
         <div>
            <label>Submission Mode</label>
            <select v-model="submitMode">
               <option value="">Select a submission mode</option>
               <option value="dev">Development (no submission)</option>
               <option value="prod">Production</option>
            </select>
            <p class="hint">Development mode will log the metadata to the job log</p>
         </div>
         <div>
            <label>Submission Name</label>
            <input type="text" v-model="submitName" />
            <p class="hint">This is an identifier that will appended to the submission file name to help identify it later. For example: batch12</p>
         </div>
         <div class="buttons">
            <DPGButton label="Cancel" severity="secondary" @click="showDialog = false"/>
            <DPGButton label="Submit"  @click="submitClicked" :disabled="submitDisabled"/>
         </div>
      </div>
   </Dialog>
</template>

<script setup>
import Dialog from 'primevue/dialog'
import { ref, computed } from 'vue'

const emit = defineEmits( ['submit' ])

const showDialog = ref(false)
const submitMode = ref("")
const submitName = ref("")

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