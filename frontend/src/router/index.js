import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
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
      path: '/metadata/:id',
      name: 'metadata',
      component: () => import('../views/MetadataDetails.vue')
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
     let jwtStr = VueCookies.get("dpg_jwt")
     userStore.setJWT(jwtStr)
     next( "/" )
  } else if (to.name !== 'not_found' && to.name !== 'forbidden' && to.name !== "signedout") {
     let jwtStr = localStorage.getItem('dpg_jwt')
     if ( jwtStr) {
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
