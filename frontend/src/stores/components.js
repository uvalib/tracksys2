import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

function getNodeData(node) {
   let data = {
      id: node.id,
      pid: node.pid,
      title: node.title,
      label: node.label,
      description: node.description,
      date: node.date,
      level: node.level,
      barcode: node.barcode,
      eadID: node.eadID,
      componentType: node.componentType.name,
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
      selectedComponent: "",
      nodes: [],
      relatedMasterFiles: [],
      loadingMasterFiles: false
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
         this.selectedComponent = id
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
            system.setError(e)
         })
      },
      loadRelatedMasterFiles( id ) {
         this.selectedComponent = id
         this.loadingMasterFiles = true
         this.relatedMasterFiles = []
         axios.get( `/api/components/${id}/masterfiles` ).then(response => {
            this.relatedMasterFiles = response.data
            this.loadingMasterFiles = false
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
            this.loadingMasterFiles = false
         })
      }
	}
})