import { createRouter, createWebHistory } from 'vue-router'
import AuthView from '../views/AuthView.vue'
import ProgramationsView from '../views/ProgramationsView.vue'
import { useAuthStore } from '@/stores/auth'
import SongsView from '@/views/SongsView.vue'
import UsersView from '@/views/UsersView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'auth',
      component: AuthView
    },
    {
      path: '/programations',
      name: 'programations',
      component: ProgramationsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/songs',
      name: 'songs',
      component: SongsView,
      meta: { requiresAuth: true }
    },
    {
      path: '/users',
      name: 'users',
      component: UsersView,
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
