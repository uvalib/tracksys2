import { defineStore } from 'pinia'
import axios from 'axios'
import { useSearchStore } from './search'


export const useSystemStore = defineStore('system', {
	state: () => ({
      working: false,
		version: "unknown",
		error: "",
      showError: false,
      jobsURL: "",
      reportsURL: "",
      pdfURL: "",
      projectsURL: "",
      iiifManifestURL: "",
      academicStatuses: [],
      agencies: [],
      categories: [],
      containerTypes: [],
      intendedUses: [],
      workflows: [],
      toast: {
         summary: "",
         message: "",
         show: false
      }
	}),
	getters: {
	},
	actions: {
      toastMessage( summary, message ) {
         this.toast.summary = summary
         this.toast.message = message
         this.toast.show = true
      },
      clearToastMessage() {
         this.toast.summary = ""
         this.toast.message = ""
         this.toast.show = false
      },
      setError( e ) {
         this.error = e
         this.showError = true
         this.working = false
      },
      async getConfig() {
         return axios.get("/api/config").then(response => {
            this.version = response.data.version
            this.reportsURL = response.data.reportsURL
            this.projectsURL = response.data.projectsURL
            this.jobsURL = response.data.jobsURL
            this.pdfURL = response.data.pdfURL
            this.iiifManifestURL = response.data.iiifManifestURL
            this.academicStatuses = response.data.controlledVocabularies.academicStatuses
            this.agencies = response.data.controlledVocabularies.agencies
            this.categories = response.data.controlledVocabularies.categories
            this.containerTypes = response.data.controlledVocabularies.containerTypes
            this.intendedUses = response.data.controlledVocabularies.intendedUses
            this.workflows = response.data.controlledVocabularies.workflows
            const searchStore = useSearchStore()
            searchStore.setGlobalSearchFields(response.data.searchFields)
         }).catch( e => {
            this.error =  e
            this.working = false
         })
      },
	}
})