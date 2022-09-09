import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useMetadataStore = defineStore('metadata', {
	state: () => ({
      detail: {
         id: 0,
         barcode: "",
         catalogKey: "",
         callNumber: "",
         title: "",
         creatorName: "",
         creatorNameType: "",
         year: "",
         publicationPlace: "",
         location: "",
      },
      dl: {
         pid: "",
         inDL: false,
         inDPLA: false,
         useRight: null,
         useRightRationale: "",
         creatorDeathDate: "",
         availabilityPolicy: "",
         collectionFacet: "",
         dateDLIngest: null,
         dateDLUpdate: null,
      },
      other: {
         parentID: 0,
         isManuscript: false,
         ocrHint: null,
         ocrLanguageHint: "",
         preservationTier: null,
      },
      thumbURL: "",
      viewerURL: "",
      virgoURL: ""
   }),
	getters: {
	},
	actions: {
      getDetails( metadataID ) {
         const system = useSystemStore()
         system.working = true
         axios.get( `/api/metadata/${metadataID}` ).then(response => {
            // general info
            this.detail.id = response.data.metadata.id
            this.detail.type = response.data.metadata.type
            this.detail.barcode = response.data.metadata.barcode
            this.detail.catalogKey = response.data.metadata.catalogKey
            this.detail.callNumber = response.data.metadata.callNumber
            this.detail.title = response.data.metadata.title
            this.detail.creatorName = response.data.metadata.creatorName
            if (this.detail.type == "XmlMetadata" || this.detail.type == "SirsiMetadata") {
               if ( response.data.details.title && response.data.details.title != "") {
                  this.detail.title = response.data.details.title
               }
               if ( response.data.details.creatorName && response.data.details.creatorName != "") {
                  this.detail.creatorName = response.data.details.creatorName
               }
               this.detail.creatorNameType = response.data.details.creatorType
               this.detail.year = response.data.details.year
               this.detail.publicationPlace = response.data.details.publicationPlace
               this.detail.location = response.data.details.location
               this.thumbURL = response.data.details.previewURL
               this.viewerURL = response.data.details.objectURL
               this.virgoURL = response.data.virgoURL
            }

            // DL info
            this.dl.pid = response.data.metadata.pid
            this.dl.inDL = (response.data.metadata.dateDLIngest != null)
            this.dl.inDPLA = response.data.metadata.dpla
            this.dl.useRight = response.data.metadata.useRight
            this.dl.useRightRationale = response.data.metadata.useRightRationale
            if ( response.data.metadata.creatorDeathDate > 0) {
               this.dl.creatorDeathDate = `${response.data.metadata.creatorDeathDate}`
            }
            this.dl.availability = response.data.metadata.availability
            this.dl.collectionFacet = response.data.metadata.collectionFacet
            this.dl.dateDLIngest = response.data.metadata.dateDLIngest
            this.dl.dateDLUpdate = response.data.metadata.dateDLUpdate

            // admin / other info
            this.other.parentID = response.data.metadata.parentID
            this.other.isManuscript = response.data.metadata.isManuscript
            this.other.ocrHint = response.data.metadata.ocrHint
            this.other.ocrLanguageHint = response.data.metadata.ocrLanguageHint
            this.other.preservationTier = response.data.metadata.preservationTier

            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      }
   }
})