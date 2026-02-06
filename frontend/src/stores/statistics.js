import { defineStore } from 'pinia'
import axios from 'axios'
import dayjs from 'dayjs'

export const useStatsStore = defineStore('stats', {
	state: () => ({
		dateRangeType: "before",
		startDate: new Date(),
		endDate: null,
		imageStats: {
			total: 0,
			DL: 0,
			DPLA: 0,
			error: "",
			loading: false,
		},
		storageStats: {
			total: 0,
			DL: 0,
			loading: false,
		},
		metadataStats: {
			total: 0,
			sirsi: 0,
			xml: 0,
			totalDL: 0,
			sirsiDL: 0,
			xmlDL: 0,
			totalDPLA: 0,
			sirsiDPLA: 0,
			xmlDPLA: 0,
			error: "",
			loading: false,
		},
		archiveStats: {
			bound: 0,
			manuscript: 0,
			photo: 0,
			error: "",
			loading: false,
		},
		publishedStats: {
			loading: false,
			virgo: [],
			archivesSpace: [],
			error: ""
		}
	}),
	getters: {
	},
	actions: {
		getAllStats(force) {
			if (this.storageStats.total == 0 || force == true) {
				this.getImageSats()
				this.getStorageSats()
				this.getMetadataSats()
				this.getArchiveSats()
				this.getPublishedSats()
			}
		},

		getImageSats() {
			let dateParam = getDateParam(this.dateRangeType, this.startDate, this.endDate)
			let url = "/api/stats/images"
			if (dateParam != "") {
				url += "?date="+encodeURIComponent(dateParam)
			}
			this.imageStats.rangeText = dateParam
			this.imageStats.loading = true
			axios.get(url).then(response => {
				this.imageStats.total = response.data.total
				this.imageStats.DL = response.data.dl
				this.imageStats.DPLA = response.data.dpla
				this.imageStats.loading = false
				this.imageStats.error = ""
			}).catch(e => {
            this.imageStats.error = e
				this.imageStats.loading = false
         })
		},

		getStorageSats() {
			let url = "/api/stats/storage"
			this.storageStats.loading = true
			axios.get(url).then(response => {
				this.storageStats.total = response.data.total
				this.storageStats.DL = response.data.dl
				this.storageStats.loading = false
				this.storageStats.error = ""
			}).catch(e => {
            this.storageStats.error = e
				this.storageStats.loading = false
         })
		},

		getPublishedSats() {
			let url = "/api/stats/published"
			this.publishedStats.loading = true
			axios.get(url).then(response => {
				this.publishedStats.virgo = response.data.virgo
				this.publishedStats.archivesSpace = response.data.archivesSpace
				this.publishedStats.loading = false
				this.publishedStats.error = ""
			}).catch(e => {
            this.publishedStats.error = e
				this.publishedStats.loading = false
         })
		},

		getMetadataSats() {
			let dateParam = getDateParam(this.dateRangeType, this.startDate, this.endDate)
			let url = "/api/stats/metadata"
			if (dateParam != "") {
				url += "?date="+encodeURIComponent(dateParam)
			}
			this.metadataStats.rangeText = dateParam
			this.metadataStats.loading = true
			axios.get(url).then(response => {
				this.metadataStats.total = response.data.all.total
				this.metadataStats.sirsi = response.data.all.sirsi
				this.metadataStats.xml = response.data.all.xml
				this.metadataStats.totalDL = response.data.DL.total
				this.metadataStats.sirsiDL = response.data.DL.sirsi
				this.metadataStats.xmlDL = response.data.DL.xml
				this.metadataStats.totalDPLA = response.data.DPLA.total
				this.metadataStats.sirsiDPLA = response.data.DPLA.sirsi
				this.metadataStats.xmlDPLA = response.data.DPLA.xml
				this.metadataStats.error = ""
				this.metadataStats.loading = false
			}).catch(e => {
            this.metadataStats.error = e
				this.metadataStats.loading = false
         })
		},

		getArchiveSats() {
			let dateParam = getDateParam(this.dateRangeType, this.startDate, this.endDate)
			let url = "/api/stats/archive"
			if (dateParam != "") {
				url += "?date="+encodeURIComponent(dateParam)
			}
			this.archiveStats.rangeText = dateParam
			this.archiveStats.loading = true
			axios.get(url).then(response => {
				this.archiveStats.bound = response.data.bound
				this.archiveStats.manuscript = response.data.manuscript
				this.archiveStats.photo = response.data.photo
				this.archiveStats.loading = false
				this.archiveStats.error = ""
			}).catch(e => {
            this.archiveStats.error = e
				this.archiveStats.loading = false
         })
		},
	}
})

function getDateParam(rangeType, startDate, endDate) {
	let dateParam = ""
	if (rangeType == "before") {
		dateParam = `BEFORE ${dayjs(startDate).format("YYYY-MM-DD")}`
	} else if (rangeType == "after") {
		dateParam = `AFTER ${dayjs(startDate).format("YYYY-MM-DD")}`
	} else {
		dateParam = `${dayjs(startDate).format("YYYY-MM-DD")} TO ${dayjs(endDate).format("YYYY-MM-DD")}`
	}
	return dateParam
}
