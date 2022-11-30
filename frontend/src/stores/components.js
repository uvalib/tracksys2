import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

function dataForComponentNode(node) {
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
      dateDLIngest: node.dateDLIngest,
      dateDLUpdate: node.dateDLUpdate,
   }
   return data
}

function getLabel( data )  {
   let label = data.title
   if (label == "") {
      label = data.label
      if (label == "") {
         label = data.descripton
      }
   }
   return label
}

function getNodeChildren( parentNode, children ) {
   children.forEach( c => {
      let newChild = {
         key: c.pid,
         label: getLabel(c),
         data: dataForComponentNode(c)
      }
      parentNode.children.push( newChild )
   })
}

export const useComponentsStore = defineStore('components', {
	state: () => ({
      nodes: []
	}),
	getters: {
	},
	actions: {
      getComponentTree( id ) {
         const system = useSystemStore()
         system.working = true
         axios.get( `/api/components/${id}` ).then(response => {
            let root = {
               key: response.data.pid,
               label: getLabel(response.data),
               data: dataForComponentNode(response.data),
               children: []
            }
            getNodeChildren(root, response.data.children)

            this.nodes = [root]
            system.working = false
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})