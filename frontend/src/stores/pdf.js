import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const usePDFStore = defineStore('pdf', {
	state: () => ({
      unitID: 0,
      downloading: false,
      includeText: false,
      percent: 0,
      intervalID: -1,
      token: "",
      targetMasterFiles: []
   }),
	getters: {
	},
	actions: {
      requestPDF( unitID, masterFileIDs = [], includeText = false ) {
         if (this.downloading) return

         const system = useSystemStore()
         this.downloading = true
         this.percent = 0
         this.token = ""
         this.unitID = unitID
         this.includeText = includeText
         this.targetMasterFiles = masterFileIDs

         let url = `/api/units/${this.unitID}/pdf`
         if ( masterFileIDs.length > 0 ) {
            url += `?pages=${masterFileIDs.join(",")}`
         }
         axios.get(url).then( resp => {
            this.token = resp.data.token
            if ( resp.data.status == "READY") {
               this.percent = 100
               this.downloading = false
               this.downloadPDF()
            } else if ( resp.data.status == "FAILED") {
               this.downloading = false
               this.percent = 0
               system.setError("Unable to generate PDF")
            } else {
               this.pollStatus()
            }
         }).catch( e => {
            system.setError(e)
            this.downloading = false
            this.percent = 0
         })
      },

      pollStatus() {
         const system = useSystemStore()
         this.intervalID = setInterval( () => {
            let statusURL = `/api/units/${this.unitID}/pdf/status?token=${this.token}`
            axios.get(statusURL).then( resp => {
               if ( resp.data == "READY") {
                  this.downloadPDF()
                  clearInterval(this.intervalID)
                  this.intervalID = -1
               } else if ( resp.data == "FAILED") {
                  this.downloading = false
                  this.percent = 0
                  system.setError("PDF generation failed")
                  clearInterval(this.intervalID)
                  this.intervalID = -1
               } else {
                  this.percent = parseInt(resp.data.replace("%", ""), 10)
               }
            }).catch( e => {
               system.setError(e)
               this.downloading = false
               this.percent = 0
               clearInterval(this.intervalID)
               this.intervalID = -1
            })
         }, 1000)
      },

      downloadPDF() {
         this.downloading = false
         this.percent = 100

         let tgtURL = new URL(window.location.href)
         let downloadURL = `${tgtURL.origin}/pdf?unit=${this.unitID}&token=${this.token}`
         if (this.includeText) {
            let pages = "all"
            if ( this.targetMasterFiles.length > 0) {
               pages = this.targetMasterFiles.join(",")
            }
            downloadURL += `&text=1&pages=${pages}`
         }

         let newWindow = window.open(downloadURL)
         if ( newWindow == null ) {
            window.location.href = downloadURL
         }
      }
   }
})