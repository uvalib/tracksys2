import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useEquipmentStore = defineStore('equipment', {
	state: () => ({
      workstations: [],
      equipment: []
	}),
	getters: {
      scanners: state => {
         return state.equipment.filter( e => e.type == "Scanner")
      },
      lenses: state => {
         return state.equipment.filter( e => e.type == "Lens")
      },
      cameraBodies: state => {
         return state.equipment.filter( e => e.type == "CameraBody")
      },
      digitalBacks: state => {
         return state.equipment.filter( e => e.type == "DigitalBack")
      },
	},
	actions: {
      getEquipment( ) {
         const system = useSystemStore()
         axios.get( `/api/equipment` ).then(response => {
            this.workstations = response.data.workstations
            this.equipment = response.data.equipment
         }).catch( e => {
            system.setError(e)
         })
      },
	}
})