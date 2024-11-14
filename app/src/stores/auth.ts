import type { JWTPayload } from '@/types'
import { jwtDecode } from 'jwt-decode'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', () => {
  const router = useRouter()
  const session = ref('')

  const isLoggedIn = computed(() => !!session.value)
  //session.value -> string = "12345..." -> truthy
  // -> con el primer bang pasa a false
  // -> con el segundo bang pasa a true
  //session.value -> string = "" -> falsey
  // -> con el primer bang pasa a true
  // -> con el segundo bang pasa a false

  async function init() {
    const tokenStr = sessionStorage.getItem('token')
    if (tokenStr) {
      setSession(tokenStr)
    }
  }

  function setSession(tokenStr: string) {
    const payload = jwtDecode(tokenStr) as JWTPayload
    const now = new Date()
    const diff = payload.MapClaims.eat * 1000 - now.getTime()

    // If we want to ask for a session refresh
    // setTimeout(
    //   () => {
    //     confirm('Session almost finished')
    //   },
    //   diff - 30 * 1000
    // )
    sessionStorage.setItem('token', tokenStr)
    session.value = payload.session
    setTimeout(() => {
      clearSession()
    }, diff)
    router.push('/users')
  }

  function clearSession() {
    session.value = ''
    sessionStorage.removeItem('token')
    router.push('/')
  }

  return { isLoggedIn, init, clearSession, setSession }
})
