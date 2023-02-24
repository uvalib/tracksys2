import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

function getNodeData(node) {
   let data = {
      id: node.id,
      pid: node.pid,
      title: node.title.replace(/\s+/g, ' ').trim(),
      label: node.label.replace(/\s+/g, ' ').trim(),
      description: node.description.replace(/\s+/g, ' ').trim(),
      date: node.date.replace(/\s+/g, ' ').trim(),
      level: node.level,
      barcode: node.barcode,
      eadID: node.eadID,
      componentType: node.componentType.name,
      masterFileCount: node.masterFileCount
   }
   return data
}

function getLabel( data )  {
   // label is concatenated data for all to facilitate filtering. the label is not displayed.
   return `${data.title} ${data.label} ${data.description} ${data.pid} ${data.date} ${data.eadID}`
}

function getNodeChildren( parentNode, children ) {
   children.forEach( c => {
      let newChild = {
         key: `${c.id}`,
         label: getLabel(c),
         data: getNodeData(c),
         children: [],
         selectable: true
      }
      parentNode.children.push( newChild )
      if (c.children && c.children.length > 0) {
         getNodeChildren(newChild, c.children)
      }
   })
}

export const useComponentsStore = defineStore('components', {
	state: () => ({
      selectedComponent: 0,
      nodes: [],
      relatedMasterFiles: [],
      loadingMasterFiles: false,
      searchHits: [],
      totalSearchHits: 0,
	}),
	getters: {
      title: state => {
         if (state.nodes[0] && state.nodes[0].data) {
            return `${state.nodes[0].label}`
         }
         return "Component Details"
      },
	},
	actions: {
      async getComponentTree( id ) {
         const system = useSystemStore()
         this.nodes = []
         system.working = true
         this.selectedComponent = parseInt(id,10)
         return axios.get( `/api/components/${id}` ).then(response => {
            let component = response.data.component
            let root = {
               key: `${component.id}`,
               label: getLabel(component),
               data: getNodeData(component),
               children: []
            }
            if ( component.children ) {
               getNodeChildren(root, component.children)
            }

            this.nodes = [root]
            this.relatedMasterFiles = []
            if (response.data.masterFiles) {
               this.relatedMasterFiles = response.data.masterFiles
            }
            system.working = false
         }).catch( e => {
            if (e.response && e.response.status == 404) {
               this.selectedComponent = -1
               this.router.push("/not_found")
               system.working = false
            } else {
               system.setError(e)
            }
         })
      },
      loadRelatedMasterFiles( id ) {
         this.selectedComponent = parseInt(id,10)
         this.loadingMasterFiles = true
         axios.get( `/api/components/${id}/masterfiles` ).then(response => {
            this.relatedMasterFiles = response.data
            this.loadingMasterFiles = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
            this.loadingMasterFiles = false
         })
      },
      async lookup( query ) {
         const system = useSystemStore()
         let url = `/api/search?scope=components&q=${encodeURIComponent(query)}&start=0&limit=30`
         return axios.get(url).then(response => {
            this.searchHits = response.data.components.hits
            this.totalSearchHits = response.data.components.total
         }).catch( e => {
            system.setError(e)
         })
      },
      downloadFromArchive( computeID, unitID, files ) {
         const system = useSystemStore()
         let payload = {computeID: computeID, files: files}
         let url = `${system.jobsURL}/units/${unitID}/copy`
         axios.post(url, payload).then( () => {
            system.toastMessage("Archive Download", `Master files are being downloaded from the archive.`)
         }).catch( e => {
            system.setError(e)
         })
      },
      assignMetadata( metadataID, unitID, masterFileIDs) {
         const system = useSystemStore()
         let data = {ids: masterFileIDs, metadataID:  parseInt(metadataID,10) }
         axios.post(`${system.jobsURL}/units/${unitID}/masterfiles/metadata`, data).then( () => {
            this.loadRelatedMasterFiles( this.selectedComponent )
            system.toastMessage("Assign Metadata Success", 'The selected master files have been assigned new metadata.')
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})