import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'
import dayjs from 'dayjs'

export const useMetadataStore = defineStore('metadata', {
	state: () => ({
      working: false,
      error: "",
      detail: {
         id: 0,
         pid: "",
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
         externalSystem: null,
         parentID: 0,
         isCollection: false,
         isManuscript: false,
         isPersonalItem: false,
         ocrHint: null,
         ocrLanguageHint: "",
         preservationTier: null,
         inDL: false,
         inDPLA: false,
         inHathiTrust: false,
         useRightName: "",
         useRightURI: "",
         useRightStatement: "",
         creatorDeathDate: "",
         availabilityPolicy: null,
         collectionID: "",
         collectionFacet: "",
         dateDLIngest: null,
         dateDLUpdate: null,
         thumbURL: "",
         viewerURL: "",
         virgoURL: "",
      },
      hathiTrustStatus: {
         requestedAt: null,
         packageCreatedAt: null,
         packageSubmittedAt: null,
         packageStatus: "",
         metadataSubmittedAt: null,
         metadataStatus: "",
         finishedAt: null,
         notes: "",
      },
      apTrustStatus: {
         etag: "",
			objectID: "",
			status: "",
			note: "",
			submittedAt: "",
			finishedAt: ""
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
         dates: "",
         publishedAt: ""
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
      related: {
         units: [],
         orders: [],
         masterFiles: [],
         collection: null
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
      hasMasterFiles: state => {
         if ( state.related.masterFiles.length > 0 ) return true
         let hasFiles = false
         state.related.units.some( u => {
            if (u.masterFilesCount > 0) {
               hasFiles = true
            }
            return hasFiles == true
         })
         return hasFiles
      },
      canPublishToVirgo: state => {
         if ( state.detail.type == 'ExternalMetadata' ) return false

         let canPublish =  false
         state.related.units.forEach( u => {
            if ( u.inDL ) {
               canPublish = true
            }
         })
         return canPublish
      }
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
         this.sirsiMatch.error = ""
         this.sirsiMatch.searching = false
         this.asMatch.error = ""
         this.asMatch.searching = false
         this.asMatch.title = ""
         this.asMatch.id = ""
      },
      resetArchivesSpaceErrors() {
         this.asMatch.error = ""
         this.asMatch.searching = false
         this.asMatch.title = ""
         this.asMatch.id = ""
         this.archivesSpace.error = ""
      },
      async validateArchivesSpaceURI( uri ) {
         this.resetSearch()
         this.asMatch.searching = true
         return axios.get(`/api/metadata/archivesspace?uri=${encodeURIComponent(uri)}`).then(response => {
            this.asMatch.validatedURL = response.data.uri
            this.asMatch.title = response.data.detail.title
            if (response.data.detail.dates) {
               this.asMatch.title = `${this.asMatch.title}, ${response.data.detail.dates}`
            }
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
            this.sirsiMatch.searching = false
         }).catch( e => {
            this.sirsiMatch.searching = false
            this.sirsiMatch.error = e
            if (e.response) {
               this.sirsiMatch.error = e.response.data
            }
         })
      },
      async lookup( query, collectionID ) {
         const system = useSystemStore()
         let url = `/api/search?scope=metadata&q=${encodeURIComponent(query)}&start=0&limit=30`
         if ( collectionID) {
            url += "&collection=true"
         }
         return axios.get(url).then(response => {
            this.searchHits = response.data.metadata.hits
            this.totalSearchHits = response.data.metadata.total
         }).catch( e => {
            system.setError(e)
         })
      },
      async publish( ) {
         const system = useSystemStore()
         return axios.post( `${system.jobsURL}/metadata/${this.detail.id}/publish` ).then( () => {
            var now = dayjs().format("YYYY-MM-DD hh:mm A")
            if ( this.detail.dateDLIngest ) {
               this.detail.dateDLUpdate = now
            } else {
               this.detail.dateDLIngest = now
            }
         }).catch( e => {
            system.setError(e)
         })
      },
      async publishToArchivesSpace( userID, immediate ) {
         const system = useSystemStore()
         let payload = {userID: `${userID}`, metadataID: `${this.detail.id}`}
         let url = `${system.jobsURL}/archivesspace/publish?immediate=${immediate}`
         return axios.post( url, payload ).then( () => {
            var now = dayjs().format("YYYY-MM-DD hh:mm A")
            this.archivesSpace.publishedAt = now
         }).catch( e => {
            system.setError(e)
         })
      },
      async create( data ) {
         const system = useSystemStore()
         return axios.post("/api/metadata", data).then(response => {
            this.setMetadataDetails(response.data)
            this.searchHits = [ {
               id: this.detail.id,
               pid: this.detail.pid,
               type: this.detail.type,
               barcode: this.detail.barcode,
               callNumber: this.detail.callNumber,
               catalogKey: this.detail.catalogKey,
               creatorName: this.detail.creatorName,
               title: this.detail.title,
            }]
            this.totalSearchHits = 1
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
            this.detail.xmlMetadata = resp.data.metadata
            this.detail.title = resp.data.title
            if ( this.detail.dateDLIngest || this.detail.dateDLUpdate ) {
               this.publish()
            }
            system.toastMessage("XML Uploaded", "XML metadata has successfully been uploaded.")
         }).catch( e => {
            system.setError(`Upload XML meadtata file '${fileData.name}' failed: ${e.response.data}`)
         })
      },

      setMetadataDetails( details ) {
         // general info
         this.$reset()
         this.detail.id = details.metadata.id
         this.detail.type = details.metadata.type
         this.detail.barcode = details.metadata.barcode
         this.detail.catalogKey = details.metadata.catalogKey
         this.detail.callNumber = details.metadata.callNumber
         this.detail.title = details.metadata.title
         this.detail.creatorName = details.metadata.creatorName
         this.detail.parentID = details.metadata.parentID
         this.detail.isCollection = details.metadata.isCollection
         this.detail.isManuscript = details.metadata.isManuscript
         this.detail.isPersonalItem = details.metadata.isPersonalItem
         this.detail.ocrHint = details.metadata.ocrHint
         this.detail.ocrLanguageHint = details.metadata.ocrLanguageHint
         this.detail.preservationTier = details.metadata.preservationTier
         this.detail.thumbURL = details.thumbURL
         this.detail.viewerURL = details.viewerURL
         this.detail.virgoURL = details.virgoURL

         if ( details.metadata.apTrustStatus) {
            this.apTrustStatus.etag = details.metadata.apTrustStatus.etag
            this.apTrustStatus.objectID = details.metadata.apTrustStatus.objectID
            this.apTrustStatus.status = details.metadata.apTrustStatus.status
            this.apTrustStatus.note = details.metadata.apTrustStatus.note
            this.apTrustStatus.submittedAt = details.metadata.apTrustStatus.submittedAt
            this.apTrustStatus.finishedAt = details.metadata.apTrustStatus.finishedAt
         }

         if (this.detail.type == "SirsiMetadata") {
            if ( details.sirsiDetails.title && details.sirsiDetails.title != "") {
               this.detail.title = details.sirsiDetails.title
            }
            if ( details.sirsiDetails.creatorName && details.sirsiDetails.creatorName != "") {
               this.detail.creatorName = details.sirsiDetails.creatorName
            }
            this.detail.creatorNameType = details.sirsiDetails.creatorType
            this.detail.year = details.sirsiDetails.year
            this.detail.publicationPlace = details.sirsiDetails.publicationPlace
            this.detail.location = details.sirsiDetails.location
            this.detail.useRightName =  details.sirsiDetails.useRightName
            this.detail.useRightURI =  details.sirsiDetails.useRightURI
            this.detail.useRightStatement =  details.sirsiDetails.useRightStatement
         }

         if ( this.detail.type == "XmlMetadata") {
            this.detail.xmlMetadata = details.metadata.descMetadata
         }

         if (this.detail.type == "ExternalMetadata") {
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
                  if (details.asDetails.published_at) {
                     this.archivesSpace.publishedAt = details.asDetails.published_at.split("T")[0]
                  }
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
         this.detail.pid = details.metadata.pid
         this.detail.inDL = (details.metadata.dateDLIngest != null)
         this.detail.inDPLA = details.metadata.dpla
         this.detail.inHathiTrust = details.metadata.hathiTrust
         if ( this.detail.inHathiTrust ) {
            this.hathiTrustStatus = details.metadata.hathiTrustStatus
         }
         if ( details.metadata.creatorDeathDate > 0) {
            this.detail.creatorDeathDate = `${details.metadata.creatorDeathDate}`
         }
         this.detail.availabilityPolicy = details.metadata.availabilityPolicy
         this.detail.collectionID = details.metadata.collectionID
         this.detail.collectionFacet = details.metadata.collectionFacet
         this.detail.dateDLIngest = details.metadata.dateDLIngest
         this.detail.dateDLUpdate = details.metadata.dateDLUpdate

         this.setRelatedItems(details.units, details.masterFiles, details.collectionRecord)
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

      async updateHathiTrustStatus( status ) {
         this.working = true
         this.error = ""
         await axios.post( `/api/metadata/${this.detail.id}/hathitrust`, status ).then(response => {
            this.hathiTrustStatus = response.data
            this.working = false
         }).catch( e => {
            this.working = false
            this.error = e
         })
      },

      setRelatedItems( units, masterFiles, collection ) {
         this.related.units = []
         this.related.orders = []
         this.related.masterFiles = []
         this.related.collection = null

         if ( collection ) {
            this.related.collection = collection
         }
         if ( masterFiles ) {
            this.related.masterFiles = masterFiles
         }
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