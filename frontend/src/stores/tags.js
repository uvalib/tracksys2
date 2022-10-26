import { defineStore } from 'pinia'
import axios from 'axios'

export const useTagsStore = defineStore('tags', {
	state: () => ({
      query: "",
      loading: false,
      tags: [],
      error: "",

   }),
	getters: {
	},
	actions: {
      setError( err ) {
         this.error = err
      },
      async getTags(  ) {
         this.loading = true
         this.error = ""
         let url = `/api/tags`
         if (this.query != "") {
            url += `?q=${this.query}`
         }
         return axios.get(url ).then(response => {
            this.tags = response.data
            this.loading = false
         }).catch( e => {
            this.loading = false
            this.error = e
         })
      },
      async createTag( name ) {
         this.loading = true
         return axios.post( "/api/tags", {tag: name} ).then(response => {
            this.tags = response.data
            this.loading = false
         }).catch( e => {
            this.loading = false
            this.error = e
         })
      },
   }
})