import { createRouter, createWebHistory } from 'vue-router'

import HomeView from '../views/HomeView.vue'

import APTrustSubmissions from '../views/APTrustSubmissions.vue'
import ArchivesSpaceReviews from '../views/ArchivesSpaceReviews.vue'
import CollectionsList from '../views/CollectionsList.vue'
import HathiTrustSubmissions from '../views/HathiTrustSubmissions.vue'
import JobStatusList from '../views/JobStatusList.vue'
import JobStatusDetails from '../views/JobStatusDetails.vue'
import MasterFileDetails from '../views/MasterFileDetails.vue'
import MasterFileEdit from '../views/MasterFileEdit.vue'
import MetadataDetails from '../views/MetadataDetails.vue'
import MetadataEdit from '../views/MetadataEdit.vue'
import OrdersList from '../views/OrdersList.vue'
import OrderDetails from '../views/OrderDetails.vue'
import OrderEdit from '../views/OrderEdit.vue'
import UnitDetails from '../views/UnitDetails.vue'
import UnitEdit from '../views/UnitEdit.vue'

import ComponentDetails from '../views/ComponentDetails.vue'
import DigitizationEquipment from '../views/DigitizationEquipment.vue'

import StaffMembers from '../views/StaffMembers.vue'
import CustomerList from '../views/CustomerList.vue'

import SignedOut from '../views/SignedOut.vue'
import ForbiddenView from '../views/ForbiddenView.vue'
import NotFound from '../views/NotFound.vue'

import VueCookies from 'vue-cookies'
import { useUserStore } from '@/stores/user'

const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL),
   routes: [
      {
         path: '/',
         name: 'home',
         component: HomeView
      },
      {
         path: '/aptrust',
         name: 'aptrust',
         component: APTrustSubmissions
      },
      {
         path: '/archivesspace',
         name: 'archivesspace',
         component: ArchivesSpaceReviews
      },
      {
         path: '/collections',
         name: 'collections',
         component: CollectionsList
      },
      {
         path: '/hathitrust',
         name: 'hathitrust',
         component: HathiTrustSubmissions
      },
      {
         path: '/jobs',
         name: 'jobs',
         component: JobStatusList
      },
      {
         path: '/jobs/:id',
         name: 'jobdetail',
         component: JobStatusDetails
      },
      {
         path: '/masterfiles/:id',
         name: 'masterfile',
         component: MasterFileDetails
      },
      {
         path: '/masterfiles/:id/edit',
         name: 'masterfileedit',
         component: MasterFileEdit
      },
      {
         path: '/metadata/:id',
         name: 'metadata',
         component: MetadataDetails
      },
      {
         path: '/metadata/:id/edit',
         name: 'metadataedit',
         component: MetadataEdit
      },
      {
         path: '/orders',
         name: 'order',
         component: OrdersList
      },
      {
         path: '/orders/new',
         name: 'neworder',
         component: OrderEdit
      },
      {
         path: '/orders/:id',
         name: 'orderdetails',
         component: OrderDetails
      },
      {
         path: '/orders/:id/edit',
         name: 'orderedit',
         component: OrderEdit
      },
      {
         path: '/units/:id',
         name: 'unit',
         component: UnitDetails
      },
      {
         path: '/units/:id/edit',
         name: 'unitedit',
         component: UnitEdit
      },
      {
         path: '/components/:id',
         name: 'component',
         component: ComponentDetails
      },
      {
         path: '/equipment',
         name: 'equipment',
         component: DigitizationEquipment
      },
      {
         path: '/staff',
         name: 'staff',
         component: StaffMembers
      },
      {
         path: '/customers',
         name: 'customers',
         component: CustomerList
      },
      {
         path: '/signedout',
         name: 'signedout',
         component: SignedOut
      },
      {
         path: '/forbidden',
         name: 'forbidden',
         component: ForbiddenView
      },
      {
         path: '/:pathMatch(.*)*',
         name: "not_found",
         component: NotFound
      }
   ]
})

router.beforeEach((to, _from, next) => {
   const userStore = useUserStore()
   if (to.path === '/granted') {
      let jwtStr = VueCookies.get("ts2_jwt")
      userStore.setJWT(jwtStr)
      let priorURL = localStorage.getItem('tsPriorURL')
      localStorage.removeItem("tsPriorURL")
      if ( priorURL && priorURL != "/granted" && priorURL != "/") {
         console.log("RESTORE "+priorURL)
         next(priorURL)
      } else {
         next("/")
      }
   } else if (to.name !== 'not_found' && to.name !== 'forbidden' && to.name !== "signedout") {
      localStorage.setItem("tsPriorURL", to.fullPath)
      let jwtStr = localStorage.getItem('ts2_jwt')
      if (jwtStr) {
         userStore.setJWT(jwtStr)
         next()
      } else {
         window.location.href = "/authenticate"
      }
   } else {
      next()
   }
})

export default router
