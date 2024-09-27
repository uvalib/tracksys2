import { createRouter, createWebHistory } from 'vue-router'

import { useCookies } from '@vueuse/integrations/useCookies'
import { useUserStore } from '@/stores/user'

const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL),
   routes: [
      {
         path: '/',
         name: 'home',
         component: () => import('../views/HomeView.vue')
      },
      {
         path: '/aptrust',
         name: 'aptrust',
         component: () => import('../views/APTrustSubmissions.vue')
      },
      {
         path: '/archivesspace',
         name: 'archivesspace',
         component: () => import('../views/ArchivesSpaceReviews.vue')
      },
      {
         path: '/collections',
         name: 'collections',
         component: () => import('../views/CollectionsList.vue')
      },
      {
         path: '/published/:type',
         name: 'published',
         component: () => import('../views/PublishedList.vue')
      },
      {
         path: '/hathitrust',
         name: 'hathitrust',
         component: () => import('../views/HathiTrustSubmissions.vue')
      },
      {
         path: '/jobs',
         name: 'jobs',
         component: () => import('../views/JobStatusList.vue')
      },
      {
         path: '/jobs/:id',
         name: 'jobdetail',
         component: () => import('../views/JobStatusDetails.vue')
      },
      {
         path: '/masterfiles/:id',
         name: 'masterfile',
         component: () => import('../views/MasterFileDetails.vue')
      },
      {
         path: '/masterfiles/:id/edit',
         name: 'masterfileedit',
         component: () => import('../views/MasterFileEdit.vue')
      },
      {
         path: '/metadata/:id',
         name: 'metadata',
         component: () => import('../views/MetadataDetails.vue')
      },
      {
         path: '/metadata/:id/edit',
         name: 'metadataedit',
         component: () => import('../views/MetadataEdit.vue')
      },
      {
         path: '/orders',
         name: 'order',
         component: () => import('../views/OrdersList.vue')
      },
      {
         path: '/orders/new',
         name: 'neworder',
         component: () => import('../views/OrderEdit.vue')
      },
      {
         path: '/orders/:id',
         name: 'orderdetails',
         component: () => import('../views/OrderDetails.vue')
      },
      {
         path: '/orders/:id/edit',
         name: 'orderedit',
         component: () => import('../views/OrderEdit.vue')
      },
      {
         path: '/units/:id',
         name: 'unit',
         component: () => import('../views/UnitDetails.vue')
      },
      {
         path: '/units/:id/edit',
         name: 'unitedit',
         component: () => import('../views/UnitEdit.vue')
      },
      {
         path: '/components/:id',
         name: 'component',
         component: () => import('../views/ComponentDetails.vue')
      },
      {
         path: '/equipment',
         name: 'equipment',
         component: () => import('../views/DigitizationEquipment.vue')
      },
      {
         path: '/staff',
         name: 'staff',
         component: () => import('../views/StaffMembers.vue')
      },
      {
         path: '/customers',
         name: 'customers',
         component: () => import('../views/CustomerList.vue')
      },
      {
         path: '/signedout',
         name: 'signedout',
         component: () => import('../views/SignedOut.vue')
      },
      {
         path: '/forbidden',
         name: 'forbidden',
         component: () => import('../views/ForbiddenView.vue')
      },
      {
         path: '/:pathMatch(.*)*',
         name: "not_found",
         component: () => import('../views/NotFound.vue')
      }
   ]
})

router.beforeEach( async (to) => {
   console.log("BEFORE ROUTE "+to.path)
   const userStore = useUserStore()
   const cookies = useCookies()
   const noAuthRoutes = ["not_found", "forbidden", "signedout"]

   if (to.path === '/granted') {
      const jwtStr = cookies.get("ts2_jwt")
      userStore.setJWT(jwtStr)
      if ( userStore.isSignedIn  ) {
         console.log(`GRANTED [${jwtStr}]`)
         let priorURL = localStorage.getItem('tsPriorURL')
         localStorage.removeItem("tsPriorURL")
         if ( priorURL && priorURL != "/granted" && priorURL != "/") {
            console.log("RESTORE "+priorURL)
            return priorURL
         }
         return "/"
      }
      return {name: "forbidden"}
   }

   // for all other routes, pull the existing jwt from storage from storage and set in the user store.
   // depending upon the page requested, this token may or may not be used.
   const jwtStr = localStorage.getItem('ts2_jwt')
   userStore.setJWT(jwtStr)

   if ( noAuthRoutes.includes(to.name)) {
      console.log("NOT A PROTECTED PAGE")
   } else {
      if (userStore.isSignedIn == false) {
         console.log("AUTHENTICATE")
         localStorage.setItem("tsPriorURL", to.fullPath)
         window.location.href = "/authenticate"
         return false   // cancel the original navigation
      } else {
         console.log(`REQUEST AUTHENTICATED PAGE WITH JWT`)
      }
   }
})

export default router
