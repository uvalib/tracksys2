import { defineStore } from 'pinia'
import axios from 'axios'

export const useSystemStore = defineStore('system', {
	state: () => ({
      working: false,
		version: "unknown",
		error: "",
      reportsURL: "",
      projectsURL: ""
	}),
	getters: {
	},
	actions: {
      setError( e ) {
         this.error = e
         this.working = false
      },
      getConfig() {
         this.working = true
         axios.get("/api/config").then(response => {
            this.version = response.data.version
            this.reportsURL = response.data.reportsURL
            this.projectsURL = response.data.projectsURL
            this.working = false
         }).catch( e => {
            this.error =  e
            this.working = false
         })
      },
	}
})