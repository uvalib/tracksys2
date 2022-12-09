<script setup>
import { onBeforeMount } from 'vue'
import { useEquipmentStore } from '@/stores/equipment'

const equipmentStore = useEquipmentStore()

onBeforeMount( async () => {
   document.title = "Equipment"
   equipmentStore.getEquipment()
})
</script>

<template>
   <h2>Equipment</h2>
   <div class="equipment">

   </div>
</template>

<style scoped lang="scss">
   .equipment {
      min-height: 600px;
   }
</style>