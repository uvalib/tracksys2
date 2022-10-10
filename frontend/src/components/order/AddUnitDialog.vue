<template>
   <DPGButton @click="show" label="Add Unit" class="p-button-secondary add"/>
   <Dialog v-model:visible="isOpen" :modal="true" header="Add Unit" :style="{width: '650px'}">
      <FormKit type="form" id="customer-detail" :actions="false" @submit="createUnit">
         <Panel header="Unit Metadata" class="margin-bottom">
         </Panel>
         <Panel header="Digitization Information">
            <FormKit label="Intended Use" type="select" v-model="unitInfo.intendedUseID" outer-class="first" :options="intendedUses"/>
            <FormKit label="Source URL" type="text" v-model="unitInfo.sourceURL"/>
            <FormKit label="Special Instructions" type="textarea" rows="4" v-model="unitInfo.specialInstructions"/>
            <FormKit label="Staff Notes" type="textarea" rows="2" v-model="unitInfo.staffNotes"/>
            <div class="opts">
               <div class="checkbox">
                  <input type="checkbox" v-model="unitInfo.completeScan"/>
                  <span class="label">Complete Scan</span>
               </div>
               <div class="checkbox">
                  <input type="checkbox" v-model="unitInfo.throwAway"/>
                  <span class="label">Throw Away</span>
               </div>
               <div class="checkbox">
                  <input type="checkbox" v-model="unitInfo.includeInDL"/>
                  <span class="label">Include in DL</span>
               </div>
            </div>
         </Panel>
         <p class="error">{{error}}</p>
         <div class="acts">
            <DPGButton @click="hide" label="Cancel" class="p-button-secondary"/>
            <span class="spacer"></span>
            <FormKit type="submit" label="Save" wrapper-class="submit-button" />
         </div>
      </FormKit>
   </Dialog>
</template>

<script setup>
import { ref, computed } from 'vue'
import Dialog from 'primevue/dialog'
import Panel from 'primevue/panel'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()

const isOpen = ref(false)
const error = ref("")
const unitInfo = ref({
   intendedUseID: null,
   sourceURL: "",
   specialInstructions: "",
   staffNotes: "",
   completeScan: false,
   throwAway: false,
   includeInDL: false,
})

const intendedUses = computed(() => {
   let out = []
   systemStore.intendedUses.forEach( a => {
      if (a.name == "Digital Collection Building") {
         out.push( {label: `${a.name}: Highest Possible resolution TIFF`, value: a.id} )
      } else {
         let dpi = "DPI"
         if (a.deliverableResolution == "Highest Possible") {
            dpi = "resolution"
         }
         out.push( {label: `${a.name}: ${a.deliverableResolution} ${dpi} ${a.deliverableFormat.toUpperCase()}`, value: a.id} )
      }
   })
   return out
})

function createUnit() {
   hide()
}

function hide() {
   isOpen.value=false
}
function show() {
   isOpen.value = true
   error.value = ""
}
</script>

<style lang="scss" scoped>
button.p-button.add {
   font-size: 0.8em;
   padding: 5px 20px;
}
div.p-panel {
   font-size: 0.85em;
}
div.margin-bottom {
   margin-bottom: 25px;
}
:deep(div.formkit-outer.first) {
   .dpg-form-label {
      margin-top: 5px;
   }
}
div.opts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: space-between;
   margin: 20px 0 10px 0;
   div.checkbox {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
      align-items: center;
      margin: 0;
      input[type=checkbox] {
         width: 18px;
         height: 18px;
         margin-right: 10px;
         display: inline-block;
      }
      span {
         display: inline-block;
      }
   }
}
.acts {
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   padding: 20px 0 10px 0;
   button {
      margin-right: 10px;
   }
}

.error {
   padding: 0;
   margin: 0;
   text-align: center;
   color: var(--uvalib-red-emergency);
}
</style>
