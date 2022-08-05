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
