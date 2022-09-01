import { defineStore } from 'pinia'
import axios from 'axios'
import { useSearchStore } from './search'

export const useSystemStore = defineStore('system', {
	state: () => ({
      working: false,
		version: "unknown",
		error: "",
      showError: false,
      reportsURL: "",
      projectsURL: "",
      academicStatuses: [],
      agencies: []
	}),
	getters: {
	},
	actions: {
      setError( e ) {
         this.error = e
         this.showError = true
         this.working = false
      },
      getConfig() {
         this.working = true
         axios.get("/api/config").then(response => {
            this.version = response.data.version
            this.reportsURL = response.data.reportsURL
            this.projectsURL = response.data.projectsURL
            this.academicStatuses = response.data.controlledVocabularies.academicStatuses
            this.agencies = response.data.controlledVocabularies.agencies
            const searchStore = useSearchStore()
            searchStore.setGlobalSearchFields(response.data.searchFields)
            this.working = false
         }).catch( e => {
            this.error =  e
            this.working = false
         })
      },
	}
})