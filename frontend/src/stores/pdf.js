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
                  this.percent = 100
                  this.downloadPDF()
                  clearInterval(this.intervalID)
               } else if ( resp.data == "FAILED") {
                  this.downloading = false
                  this.percent = 0
                  system.setError("PDF generation failed")
                  clearInterval(this.intervalID)
               } else {
                  this.percent = parseInt(resp.data.replace("%", ""), 10)
               }
            }).catch( e => {
               system.setError(e)
               this.downloading = false
               this.percent = 0
            })
         }, 1000)
      },

      downloadPDF() {
         let downloadURL = `/api/units/${this.unitID}/pdf/download?token=${this.token}`
         if (this.includeText) {
            let pages = "all"
            if ( this.targetMasterFiles.length > 0) {
               pages = this.targetMasterFiles.join(",")
            }
            downloadURL += `&text=1&pages=${pages}`
         }

         const system = useSystemStore()
         axios.get(downloadURL, {responseType: "blob"}).then((response) => {
            let dataType = "application/zip"
            let ext = "zip"
            if ( response.headers['content-type'] == "application/pdf") {
               dataType = "application/pdf"
               ext = "pdf"
            }
            var fileURL = window.URL.createObjectURL(response.data, { type: dataType})
            var fileLink = document.createElement('a')
            fileLink.href = fileURL
            fileLink.setAttribute('download', `unit-${this.unitID}.${ext}`)
            fileLink.click()
            fileLink.remove()
            window.URL.revokeObjectURL(fileURL)
            this.downloading = false
         }).catch( e => {
            system.setError(e)
            this.downloading = false
            this.percent = 0
         })
      }
   }
})