import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuditStore = defineStore('audit', {
	state: () => ({
		loading: false,
		labels: [],
		datasets: [{data: [], backgroundColor: []}],
		totalAudited: 0,
		error: "",
		auditYears: [],
		targetYear: "all"
	}),
	getters: {
	},
	actions: {
		getAuditReport( ) {
			this.loading = true
			this.labels = []
			this.datasets = [],
			this.totalAudited = 0
			if (this.auditYears.length == 0 ) {
				this.auditYears.push({label: "All", value: "all"})
				let startYear = 2009
				let endYear = new Date().getFullYear()
				for ( let y=startYear; y<=endYear;y++) {
					this.auditYears.push({label: `${y}`, value: `${y}`})
				}
			}
			axios.get(`/api/reports/audit?year=${this.targetYear}`).then(response => {
				let result = {data: [],
					backgroundColor: ["#44aacc","#cc4444","#cc4444","#cc4444", "#cc4444"],
				}
				response.data.results.forEach( r => {
					this.labels.push( `${r.label} (${r.total})` )
					result.data.push( r.total )
				})

				this.datasets.push(result)
				this.totalAudited = response.data.totalAudited
				this.loading = false
				this.error = ""
			}).catch(e => {
            this.error = e
				this.loading = false
         })
		},
	}
})