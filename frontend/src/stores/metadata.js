import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useMetadataStore = defineStore('metadata', {
	state: () => ({
      detail: {
         id: 0,
         type: "",
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
         supplementalSystem: "",
         externalURI: "",
         externalSystem: null
      },
      dl: {
         pid: "",
         inDL: false,
         inDPLA: false,
         useRight: null,
         useRightRationale: "",
         creatorDeathDate: "",
         availabilityPolicy: null,
         collectionID: "",
         dateDLIngest: null,
         dateDLUpdate: null,
      },
      archivesSpace: {
         id: "",
         createdBy: "",
         createDate: null,
         level: "",
         URL: "",
         repo: "",
         collectionTitle: "",
         language: "",
         dates: ""
      },
      jstor: {
         id: "",
         ssid: "",
         desc: "",
         creator: "",
         date: "",
         collectionID: "",
         collection: "",
         width: 0,
         height: 0,
      },
      apollo: {
         pid: "",
         type: "",
         collectionPID: "",
         collectionTitle: "",
         collectionBarcode: "",
         collectionCatalogKey: "",
         itemURL: "",
         collectionURL: "",
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
      },
      searchHits: [],
      totalSearchHits: 0,
      sirsiMatch: {
         catalogKey: "",
         barcode: "",
         callNumber: "",
         title: "",
         creatorName: "",
         creatorType: "",
         year: "",
         publicationPlace: "",
         location: "",
         collectionID: "",
         error: "",
         searching: false,
      },
      asMatch: {
         error: "",
         searching: false,
         validatedURL: "",
         title: "",
         id: "",
      }
   }),
	getters: {
	},
	actions: {
      resetSearch() {
         this.searchHits = []
         this.totalSearchHits = 0
         this.sirsiMatch.catalogKey = ""
         this.sirsiMatch.barcode = ""
         this.sirsiMatch.callNumber = ""
         this.sirsiMatch.title = ""
         this.sirsiMatch.creatorName = ""
         this.sirsiMatch.collectionID = ""
         this.sirsiMatch.error = ""
         this.sirsiMatch.searching = false
         this.asMatch.error = ""
         this.asMatch.searching = false
         this.asMatch.title = ""
         this.asMatch.id = ""
      },
      async validateArchivesSpaceURI( uri ) {
         this.resetSearch()
         this.asMatch.searching = true
         return axios.get(`/api/metadata/archivesspace?uri=${encodeURIComponent(uri)}`).then(response => {
            this.asMatch.validatedURL = response.data.uri
            this.asMatch.title = response.data.detail.title
            this.asMatch.id = response.data.detail.id
            this.asMatch.searching = false
         }).catch( e => {
            this.asMatch.searching = false
            this.asMatch.error = e
            if (e.response) {
               this.asMatch.error = e.response.data
            }
         })
      },
      async sirsiLookup( barcode, catKey ) {
         this.resetSearch()
         this.sirsiMatch.searching = true
         let url = "/api/metadata/sirsi"
         if ( catKey != "") {
            url += `?ckey=${catKey}`
         } else {
            url += `?barcode=${barcode}`
         }
         return axios.get(url).then(response => {
            this.sirsiMatch.catalogKey = response.data.catalogKey
            this.sirsiMatch.barcode =  response.data.barcode
            this.sirsiMatch.callNumber =  response.data.callNumber
            this.sirsiMatch.title =  response.data.title
            this.sirsiMatch.creatorName =  response.data.creatorName
            this.sirsiMatch.creatorType =  response.data.creatorType
            this.sirsiMatch.year =  response.data.year
            this.sirsiMatch.publicationPlace =  response.data.publicationPlace
            this.sirsiMatch.location =  response.data.location
            this.sirsiMatch.collectionID =  response.data.collectionID
            this.sirsiMatch.searching = false
         }).catch( e => {
            this.sirsiMatch.searching = false
            this.sirsiMatch.error = e
            if (e.response) {
               this.sirsiMatch.error = e.response.data
            }
         })
      },
      async lookup( query ) {
         const system = useSystemStore()
         let url = `/api/search?scope=metadata&q=${encodeURIComponent(query)}&start=0&limit=30`
         return axios.get(url).then(response => {
            this.searchHits = response.data.metadata.hits
            this.totalSearchHits = response.data.metadata.total
         }).catch( e => {
            system.setError(e)
         })
      },
      async publish( ) {
         const system = useSystemStore()
         system.working = true
         return axios.post( `${system.jobsURL}/metadata/${this.detail.id}/publish` ).then( () => {
            this.getDetails(this.detail.id)
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      async create( data ) {
         const system = useSystemStore()
         system.working = true
         return axios.post("/api/metadata", data).then(response => {
            this.setMetadataDetails(response.data)
            this.searchHits = [ {
               id: this.detail.id,
               pid: this.dl.pid,
               type: this.detail.type,
               barcode: this.detail.barcode,
               callNumber: this.detail.callNumber,
               catalogKey: this.detail.catalogKey,
               creatorName: this.detail.creatorName,
               title: this.detail.title,
               virgoURL: "",
            }]
            this.totalSearchHits = 1
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
      async uploadXML( fileData ) {
         const system = useSystemStore()
         var formData = new FormData()
         formData.append("xml", fileData)
         let url = `/api/metadata/${this.detail.id}/xml`
         return axios.post(url, formData, {
            headers: {
              'Content-Type': 'multipart/form-data'
            }
         }).then( resp => {
            this.detail.xmlMetadata = resp.data
            system.toastMessage("XML Uploaded", "XML metadata has successfully been uploaded.")
         }).catch( e => {
            system.setError(`Upload XML meadtata file '${fileData.name}' failed: ${e}`)
         })
      },

      setMetadataDetails( details ) {
         // general info
         this.detail.id = details.metadata.id
         this.detail.type = details.metadata.type
         this.detail.barcode = details.metadata.barcode
         this.detail.catalogKey = details.metadata.catalogKey
         this.detail.callNumber = details.metadata.callNumber
         this.detail.title = details.metadata.title
         this.detail.creatorName = details.metadata.creatorName
         if (this.detail.type == "XmlMetadata" || this.detail.type == "SirsiMetadata") {
            if ( details.details.title && details.details.title != "") {
               this.detail.title = details.details.title
            }
            if ( details.details.creatorName && details.details.creatorName != "") {
               this.detail.creatorName = details.details.creatorName
            }
            this.detail.creatorNameType = details.details.creatorType
            this.detail.year = details.details.year
            this.detail.publicationPlace = details.details.publicationPlace
            this.detail.location = details.details.location
            this.thumbURL = details.details.previewURL
            this.viewerURL = details.details.objectURL
            this.virgoURL = details.virgoURL
            this.detail.xmlMetadata = details.metadata.descMetadata
         } else  {
            this.detail.externalSystem = details.metadata.externalSystem
            this.detail.externalURI = details.metadata.externalURI
            if (details.metadata.externalSystem.name == "ArchivesSpace" ) {
               this.archivesSpace.error = ""
               if ( details.asDetails) {
                  this.archivesSpace.id = details.asDetails.id
                  this.archivesSpace.createdBy = details.asDetails.created_by
                  this.archivesSpace.createDate = details.asDetails.create_time.split("T")[0]
                  this.archivesSpace.level = details.asDetails.level
                  this.archivesSpace.URL = details.asDetails.url
                  this.archivesSpace.repo = details.asDetails.repo
                  this.archivesSpace.collectionTitle = details.asDetails.collection_title
                  this.archivesSpace.language = details.asDetails.language
                  this.archivesSpace.dates = details.asDetails.dates
               } else {
                  this.archivesSpace.error = "Unable to get details for "+details.metadata.externalURI
               }
            } else if (details.metadata.externalSystem.name == "JSTOR Forum") {
               this.jstor = details.jstorDetails
            } else if (details.metadata.externalSystem.name == "Apollo") {
               this.apollo = details.apolloDetails
            }
         }
         if (details.metadata.supplementalURI) {
            this.detail.supplementalURL = `${details.metadata.supplementalSystem.publicURL}/${details.metadata.supplementalURI}`
            this.detail.supplementalSystem = details.metadata.supplementalSystem.name
         }

         // DL info
         this.dl.pid = details.metadata.pid
         this.dl.inDL = (details.metadata.dateDLIngest != null)
         this.dl.inDPLA = details.metadata.dpla
         this.dl.useRight = details.metadata.useRight
         this.dl.useRightRationale = details.metadata.useRightRationale
         if ( details.metadata.creatorDeathDate > 0) {
            this.dl.creatorDeathDate = `${details.metadata.creatorDeathDate}`
         }
         this.dl.availabilityPolicy = details.metadata.availabilityPolicy
         this.dl.collectionID = details.metadata.collectionID
         this.dl.dateDLIngest = details.metadata.dateDLIngest
         this.dl.dateDLUpdate = details.metadata.dateDLUpdate

         // admin / other info
         this.other.parentID = details.metadata.parentID
         this.other.isManuscript = details.metadata.isManuscript
         this.other.isPersonalItem = details.metadata.isPersonalItem
         this.other.ocrHint = details.metadata.ocrHint
         this.other.ocrLanguageHint = details.metadata.ocrLanguageHint
         this.other.preservationTier = details.metadata.preservationTier

         this.setRelatedItems(details.units)
      },
      async getDetails( metadataID ) {
         if (this.detail.id == metadataID) return

         const system = useSystemStore()
         system.working = true
         return axios.get( `/api/metadata/${metadataID}` ).then(response => {
            this.setMetadataDetails(response.data)
            system.working = false
         }).catch( e => {
            if (e.response && e.response.status == 404) {
               system.working = false
               this.router.push("/not_found")
            } else {
               system.setError(e)
            }
         })
      },

      async deleteMetadata() {
         const system = useSystemStore()
         system.working = true
         return axios.delete( `/api/metadata/${this.detail.id}` ).then( () => {
            system.working = false
            this.router.push("/")
         }).catch( e => {
            system.setError(e)
         })
      },

      submitEdit( update ) {
         const system = useSystemStore()
         system.working = true
         axios.post( `/api/metadata/${this.detail.id}`, update ).then(response => {
            this.setMetadataDetails(response.data)
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },

      setRelatedItems( units ) {
         this.related.units = []
         this.related.orders = []
         let orderIDs = []
         units.forEach( r => {
            let u = {
               id: r.id,
               reorder: r.reorder,
               inDL: r.includeInDL,
               dateArchived: r.dateArchived,
               dateDLDeliverablesReady: r.dateDLDeliverablesReady,
               datePatronDeliverablesReady: r.datePatronDeliverablesReady,
               masterFilesCount: r.masterFilesCount,
               intendedUse: r.intendedUse,
               metadata: r.metadata
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
      }
   }
})