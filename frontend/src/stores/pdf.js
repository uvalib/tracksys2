import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const usePDFStore = defineStore('pdf', {
	state: () => ({
      unitID: 0,
      downloading: false,
      percent: 0,
      intervalID: -1,
      token: ""
   }),
	getters: {
	},
	actions: {
      requestPDF( unitID, masterFileIDs ) {
         if (this.downloading) return

         const system = useSystemStore()
         this.downloading = true
         this.percent = 0
         this.token = ""
         this.unitID = unitID
         let url = `/api/units/${this.unitID}/pdf`
         if ( masterFileIDs ) {
            url += `?pages=${masterFileIDs.join(",")}`
         }
         axios.get(url).then( resp => {
            if ( resp.data.status == "READY") {
               this.percent = 100
               this.token = resp.data.token
               this.downloadPDF()
            } else if ( resp.data.status == "FAILED") {
               this.downloading = false
               this.percent = 0
               system.setError("Unable to generate PDF")
            } else {
               this.token = resp.data.token
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
            let statusURL = `/api/units/${this.unitID}/pdf/status`
            if (this.token != "") {
               statusURL += `?token=${this.token}`
            }
            axios.get(statusURL).then( resp => {
               console.log("PDF STATUS: "+resp.data)
               if ( resp.data == "READY") {
                  console.log("DONE")
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
         let downloadURL = `/api/units/${this.unitID}/pdf/download`
         if (this.token != "") {
            downloadURL += `?token=${this.token}`
         }

         const system = useSystemStore()
         axios.get(downloadURL, {responseType: "blob"}).then((response) => {
            var fileURL = window.URL.createObjectURL(response.data, { type: 'application/pdf'})
            var fileLink = document.createElement('a')
            fileLink.href = fileURL
            fileLink.setAttribute('download', `unit-${this.unitID}.pdf`)
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