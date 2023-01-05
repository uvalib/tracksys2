import axios from 'axios'
import { defineStore } from 'pinia'
import { useSystemStore } from './system'

function parseJwt(token) {
   var base64Url = token.split('.')[1]
   var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
   var jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
   }).join(''))

   return JSON.parse(jsonPayload);
}


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
   },
   actions: {
      signout() {
         localStorage.removeItem("ts2_jwt")
         this.jwt = ""
         this.firstName = ""
         this.lastName = ""
         this.role = ""
         this.computeID = ""
         this.ID = 0
      },
      setJWT(jwt) {
         if (jwt != this.jwt) {
            this.jwt = jwt
            localStorage.setItem("ts2_jwt", jwt)

            let parsed = parseJwt(jwt)
            this.ID = parsed.userID
            this.computeID = parsed.computeID
            this.firstName = parsed.firstName
            this.lastName = parsed.lastName
            this.role = parsed.role

            // add interceptor to put bearer token in header
            axios.interceptors.request.use(config => {
               config.headers['Authorization'] = 'Bearer ' + jwt
               return config
            }, error => {
               return Promise.reject(error)
            })

            // Catch 401 errors and redirect to an expired auth page
            const system = useSystemStore()
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
         }
      },
   }
})