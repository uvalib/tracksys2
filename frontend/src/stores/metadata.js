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
         xmlMetadata: "",
         supplementalURL: "",
         supplementalSystem: ""
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
      archivesSpace: {
         id: "",
         title: "",
         createdBy: "",
         createDate: null,
         level: "",
         URL: "",
         repo: "",
         collectionTitle: "",
         language: "",
         dates: ""
      },
      other: {
         parentID: 0,
         isManuscript: false,
         isPersonalItem: false,
         ocrHint: null,
         ocrLanguageHint: "",
         preservationTier: null,
      },
      thumbURL: "",
      viewerURL: "",
      virgoURL: "",
      related: {
         units: [],
         orders: []
      }
   }),
	getters: {
	},
	actions: {
      async publish( ) {
         const system = useSystemStore()
         system.working = true
         return axios.post( `/api/metadata/${this.detail.id}/publish` ).then( () => {
            this.getDetails(this.detail.id)
         }).catch( e => {
            system.setError(e)
         })
      },
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
               this.detail.xmlMetadata = response.data.metadata.descMetadata
            } else  {
               this.detail.externalSystem = response.data.metadata.externalSystem.name
               this.detail.externalURL = `${response.data.metadata.externalSystem.publicURL}${response.data.metadata.externalURI}`
               if (response.data.metadata.externalSystem.name == "ArchivesSpace") {
                  this.archivesSpace.id = response.data.asDetails.id
                  this.archivesSpace.title = response.data.asDetails.title
                  this.archivesSpace.createdBy = response.data.asDetails.created_by
                  this.archivesSpace.createDate = response.data.asDetails.create_time.split("T")[0]
                  this.archivesSpace.level = response.data.asDetails.level
                  this.archivesSpace.URL = response.data.asDetails.url
                  this.archivesSpace.repo = response.data.asDetails.repo
                  this.archivesSpace.collectionTitle = response.data.asDetails.collection_title
                  this.archivesSpace.language = response.data.asDetails.language
                  this.archivesSpace.dates = response.data.asDetails.dates
               }
            }
            if (response.data.metadata.supplementalURI) {
               this.detail.supplementalURL = `${response.data.metadata.supplementalSystem.publicURL}/${response.data.metadata.supplementalURI}`
               this.detail.supplementalSystem = response.data.metadata.supplementalSystem.name
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
            this.other.isPersonalItem = response.data.metadata.isPersonalItem
            this.other.ocrHint = response.data.metadata.ocrHint
            this.other.ocrLanguageHint = response.data.metadata.ocrLanguageHint
            this.other.preservationTier = response.data.metadata.preservationTier

            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      getRelatedItems( metadataID ) {
         const system = useSystemStore()
         this.related.units = []
         this.related.orders = []
         axios.get( `/api/metadata/${metadataID}/related` ).then(response => {
            let orderIDs = []
            response.data.forEach( r => {
               let u = {
                  id: r.id,
                  reorder: r.reorder,
                  inDL: r.includeInDL,
                  dateArchived: r.dateArchived,
                  dateDLDeliverablesReady: r.dateDLDeliverablesReady,
                  datePatronDeliverablesReady: r.datePatronDeliverablesReady,
                  masterFilesCount: r.masterFilesCount
               }
               if ( r.intendedUse) {
                  u.intendedUse = r.intendedUse.name
               }
               this.related.units.push(u)
               if (orderIDs.includes(r.order.id) == false ) {
                  orderIDs.push(r.order.id)
                  this.related.orders.push({
                     id: r.order.id,
                     title: r.order.title,
                     customer: r.order.customer,
                     agency: r.order.agency,
                     staffNotes: r.order.staffNotes,
                     specialInstructions: r.order.specialInstructions,
                  })
               }
            })
         }).catch( e => {
            system.setError(e)
         })

      }

   }
})