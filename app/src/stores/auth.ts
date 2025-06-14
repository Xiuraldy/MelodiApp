import type { JWTPayload } from '@/types'
import { jwtDecode } from 'jwt-decode'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  const session = ref('')

  const isLoggedIn = computed(() => !!session.value)

  async function init() {
    const tokenStr = sessionStorage.getItem('token')
    if (tokenStr) {
      setSession(tokenStr)
    }
  }

  function setSession(tokenStr: string) {
    console.log('entra')
    const payload = jwtDecode(tokenStr) as JWTPayload
    console.log('payload', payload)
    const now = new Date()
    const diff = payload.exp * 1000 - now.getTime()
    sessionStorage.setItem('token', tokenStr)
    session.value = payload.sub
    setTimeout(() => {
      clearSession()
    }, diff)
    console.log('router.push')
    router.push('/programations')
  }

  function clearSession() {
    session.value = ''
    sessionStorage.removeItem('token')
    router.push('/')
  }

  return { isLoggedIn, init, clearSession, setSession, session }
})
