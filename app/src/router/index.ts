import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '../views/AuthView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'auth',
      component: AuthView
    },
    {
      path: '/census',
      name: 'census',
      component: () => import('../views/CensusView.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  if (to.name === 'auth' && authStore.isLoggedIn) {
    next({ name: 'census' })
  } else if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    next({ name: 'auth' })
  } else {
    next()
  }
})

export default router
