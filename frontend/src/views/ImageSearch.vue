<template>
   <h2>Search Images</h2>
   <div class="image-search">
      <p class="hint">Set a similarity threshold, select a search image then click 'Search Images'</p>
      <div class="labels">
         <span>More Similar</span>
         <span>Less Similar</span>
      </div>
      <USlider class="w-full" :min="5" :max="20" v-model="searchStore.distance" @change="slideChanged"/>
      <UFileUpload accept="image/*" label="Drop your image here" v-model="lookupImage" @change="uploadImageChanged"/>
      <UButton label="Search Images" @click="startImageUpload" :disabled="!lookupImage"/>
   </div>
   <div class="search">
      <SimilarImages v-if="system.working == false && searchStore.similarSearch" />
   </div>
</template>

<script setup>
import { ref } from 'vue'
import { useSearchStore } from '../stores/search'
import { useSystemStore } from '../stores/system'
import SimilarImages from '@/components/results/SimilarImages.vue'

const searchStore = useSearchStore()
const system = useSystemStore()

const lookupImage = ref()

const slideChanged = ( () => {
   if (searchStore.similarSearch == true && searchStore.searchPHash !== 0) {
      searchStore.imageSearch()
   }
})

const uploadImageChanged = (() => {
   searchStore.resetImageSearch()
})

const startImageUpload = ( async () => {
   searchStore.uploadSearchImage( lookupImage.value )
})

</script>

<style lang="scss" scoped>
.image-search {
   margin: 25px auto 0 auto;
   width: 600px;
   display: flex;
   flex-direction: column;
   gap: 25px;
   align-items: center;
   label {
      font-weight: bold;
   }
   .labels {
      width: 100%;
      font-size: 0.85em;
      display: flex;
      flex-flow: row nowrap;
      justify-content: space-between;
   }
}
</style>