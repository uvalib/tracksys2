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
      availabilityPolicies: [],
      categories: [],
      containerTypes: [],
      intendedUses: [],
      ocrHints: [],
      ocrLanguageHints: [],
      preservationTiers: [],
      useRights: [],
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
         if (e.response && e.response.data) {
            this.error = e.response.data
         }
         this.showError = true
         this.working = false
      },
      async getConfig() {
         this.working = true
         return axios.get("/config").then(response => {
            this.version = response.data.version
            this.reportsURL = response.data.reportsURL
            this.projectsURL = response.data.projectsURL
            this.jobsURL = response.data.jobsURL
            this.pdfURL = response.data.pdfURL
            this.iiifManifestURL = response.data.iiifManifestURL
            this.academicStatuses = response.data.controlledVocabularies.academicStatuses
            this.agencies = response.data.controlledVocabularies.agencies
            this.availabilityPolicies = response.data.controlledVocabularies.availabilityPolicies
            this.categories = response.data.controlledVocabularies.categories
            this.containerTypes = response.data.controlledVocabularies.containerTypes
            this.intendedUses = response.data.controlledVocabularies.intendedUses
            this.ocrHints = response.data.controlledVocabularies.ocrHints
            this.ocrLanguageHints = response.data.controlledVocabularies.ocrLanguageHints
            this.preservationTiers = response.data.controlledVocabularies.preservationTiers
            this.useRights = response.data.controlledVocabularies.useRights
            this.workflows = response.data.controlledVocabularies.workflows
            const searchStore = useSearchStore()
            searchStore.setGlobalSearchFields(response.data.searchFields)
            this.working = false
         }).catch( e => {
            this.setError(  e )
         })
      },
	}
})