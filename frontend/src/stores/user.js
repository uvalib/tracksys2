import axios from 'axios'
import { defineStore } from 'pinia'
import { useSystemStore } from './system'
import { useJwt } from '@vueuse/integrations/useJwt'

export const useUserStore = defineStore('user', {
   state: () => ({
      jwt: "",
      firstName: "",
      lastName: "",
      role: "",
      computeID: "",
      ID: 0
   }),
   getters: {
      isAdmin: state => {
         return state.role == "admin"
      },
      isSupervisor: state => {
         return state.role == "supervisor"
      },
      isStudent: state => {
         return state.role == "student"
      },
      isViewer: state => {
         return state.role == "viewer"
      },
      signedInUser: state => {
         return `${state.firstName} ${state.lastName} (${state.computeID})`
      },
      isSignedIn: state => {
         return state.jwt != "" && state.computeID != ""
      }
   },
   actions: {
      signout() {
         localStorage.removeItem("ts2_jwt")
         localStorage.removeItem("tsPriorURL")
         this.jwt = ""
         this.firstName = ""
         this.lastName = ""
         this.role = ""
         this.computeID = ""
         this.ID = 0
      },
      setJWT(jwt) {
         if (jwt == this.jwt || jwt == "" || jwt == null || jwt == "null")  return

         this.jwt = jwt
         localStorage.setItem("ts2_jwt", jwt)

         const { payload } = useJwt(jwt)
         this.ID = payload.value.userID
         this.computeID = payload.value.computeID
         this.firstName = payload.value.firstName
         this.lastName = payload.value.lastName
         this.role = payload.value.role

         // add interceptor to put bearer token in header
         const system = useSystemStore()
         axios.interceptors.request.use(config => {
            let url = config.url
            if ( url.match(system.iiifManifestURL) || url.match(system.jobsURL) ) {
               console.log("SKIP AUTH HEADER")
               return config
            }
            console.log("ADD AUTH HEADER")
            config.headers['Authorization'] = 'Bearer ' + jwt
            return config
         }, error => {
            return Promise.reject(error)
         })

         // Catch 401 errors and redirect to an expired auth page
         axios.interceptors.response.use(
            res => res,
            err => {
               console.log("failed response for "+err.config.url)
               console.log(err)
               if (err.config.url.match(/\/authenticate/)) {
                  this.router.push("/forbidden")
               } else {
                  if (err.response && err.response.status == 401) {
                     localStorage.removeItem("ts2_jwt")
                     this.jwt = ""
                     this.firstName = ""
                     this.lastName = ""
                     this.role = ""
                     this.computeID = ""
                     this.ID = 0
                     system.working = false
                     this.router.push("/signedout?expired=1")
                     return new Promise(() => { })
                  }
               }
               return Promise.reject(err)
            }
         )
      },
   }
})