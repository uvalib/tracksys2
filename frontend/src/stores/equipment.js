import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import axios from 'axios'

export const useEquipmentStore = defineStore('equipment', {
	state: () => ({
      workstations: [],
      equipment: [],
      pendingEquipment: {
         workstationID: 0,
         changed: false,
         equipment: []
      }
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
      workstationSelected( wsID ) {
         this.pendingEquipment.workstationID = wsID
         this.pendingEquipment.changed = false
         let ws = this.workstations.find( ws => ws.id == wsID )
         this.pendingEquipment.equipment = ws.equipment.slice()
      },
      deactivateWorkstation( wsID ) {
         axios.post( `/api/workstation/${wsID}/update?status=1` ).then(() => {
            let tgtWS = this.workstations.find(ws => ws.id == wsID)
            if (tgtWS) {
               tgtWS.status = 1
            }
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         })
      },
      activateWorkstation( wsID ) {
         axios.post( `/api/workstation/${wsID}/update?status=0` ).then(() => {
            let tgtWS = this.workstations.find(ws => ws.id == wsID)
            if (tgtWS) {
               tgtWS.status = 0
            }
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         })
      },
      retireWorkstation( wsID ) {
         axios.post( `/api/workstation/${wsID}/update?status=2` ).then(() => {
            let wsIdx = this.workstations.findIndex(ws => ws.id == wsID)
            this.workstations.splice(wsIdx, 1)
         }).catch( e => {
            const system = useSystemStore()
            system.setError(e)
         })
      },
	}
})