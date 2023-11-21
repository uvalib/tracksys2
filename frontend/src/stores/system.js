import { defineStore } from 'pinia'
import axios from 'axios'
import { useSearchStore } from './search'


export const useSystemStore = defineStore('system', {
	state: () => ({
      working: false,
		version: "unknown",
		error: "",
      showError: false,
      apTrustURL: "",
      jobsURL: "",
      reportsURL: "",
      projectsURL: "",
      iiifManifestURL: "",
      academicStatuses: [],
      agencies: [],
      availabilityPolicies: [],
      collectionFacets: [],
      categories: [],
      containerTypes: [],
      externalSystems: [],
      intendedUses: [],
      ocrHints: [],
      ocrLanguageHints: [],
      preservationTiers: [],
      useRights: [],
      workflows: [],
      toast: {
         error: false,
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
         this.toast.error = false
      },
      toastError( summary, message, success=true ) {
         this.toast.summary = summary
         this.toast.message = message
         this.toast.show = true
         this.toast.error = true
      },
      clearToastMessage() {
         this.toast.summary = ""
         this.toast.message = ""
         this.toast.show = false
         this.toast.error = false
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
            this.apTrustURL = response.data.apTrustURL
            this.reportsURL = response.data.reportsURL
            this.projectsURL = response.data.projectsURL
            this.jobsURL = response.data.jobsURL
            this.iiifManifestURL = response.data.iiifManifestURL
            this.academicStatuses = response.data.controlledVocabularies.academicStatuses
            this.agencies = response.data.controlledVocabularies.agencies
            this.availabilityPolicies = response.data.controlledVocabularies.availabilityPolicies
            this.categories = response.data.controlledVocabularies.categories
            this.collectionFacets = response.data.controlledVocabularies.collectionFacets
            this.containerTypes = response.data.controlledVocabularies.containerTypes
            this.externalSystems = response.data.controlledVocabularies.externalSystems
            this.intendedUses = response.data.controlledVocabularies.intendedUses
            this.ocrHints = response.data.controlledVocabularies.ocrHints
            this.ocrLanguageHints = response.data.controlledVocabularies.ocrLanguageHints
            this.preservationTiers = response.data.controlledVocabularies.preservationTiers
            this.useRights = response.data.controlledVocabularies.useRights
            this.workflows = response.data.controlledVocabularies.workflows
            const searchStore = useSearchStore()
            searchStore.setGlobalSearchFields(response.data.searchFields)
            this.working = false
         }).catch( err => {
            if (err.response && err.response.status == 401) {
               console.log("Session expired in getConfig")
            } else {
               this.setError(  err )
            }
         })
      },
      async createCollectionFacet( newFacet ) {
         this.working = true
         let data = {facet: newFacet}
         return axios.post("/api/collection-facet", data).then(response => {
            this.collectionFacets = response.data
            this.working = false
         }).catch( err => {
            this.setError(  err )
         })
      }
	}
})