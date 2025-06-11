<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from './stores/auth'
import { apiCall } from '@/services/utils'
import { onBeforeMount } from 'vue'

const authStore = useAuthStore()

onBeforeMount(authStore.init)

async function logout() {
  try {
    
    await apiCall('/auth/logout', {
      method: 'POST', 
    })
    authStore.clearSession()
  } catch (error) {
    console.error('Error during logout:', error)

    authStore.clearSession()
  }
}
</script>

<template>
  <div class="img-nav">
    <img src="/logo.png" alt="logo" />
      <nav>
        <RouterLink v-if="authStore.isLoggedIn" to="/programations">Programaciones</RouterLink>
        <RouterLink v-if="authStore.isLoggedIn" to="/songs">Canciones</RouterLink>
        <RouterLink v-if="authStore.isLoggedIn" to="/users">Usuarios</RouterLink>
        <button v-if="authStore.isLoggedIn" @click="async() => {
          await logout()
        }">Cerrar Sesión</button>
      </nav>
  </div>
  <RouterView />
  <footer>
    <h3>© Proyecto SENA - MelodiApp</h3> 
  </footer>
</template>

<style scoped>

h1 {
  font-size: 22px;
}

img {
  width: 100px;
}

nav {
  width: 100%;
  font-size: 12px;
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 10vh;
  background-color: transparent;
  color: #fff;
  padding: 30px;
  transition: all 0.3s ease;
}

nav a.router-link-exact-active {
  color: #fff;
  background-color: var(--color-secundary);
  padding: 0px 20px 0px 20px;
  font-size: initial;
  border-radius: 5px;
}

nav a.router-link-exact-active:hover {
  background-color: transparent;
}

nav a {
  display: inline-block;
  padding: 0 1rem;
  border-left: 1px solid var(--color-border);
}

nav a:first-of-type {
  border: 0;
}

.img-nav {
  display: flex;
  align-items: center;
}

footer {
  height: 6vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--color-tertiary);
}

h3 {
  margin-right: 10px;
}

button:hover {
  background-color: var(--color-secundary);
  color: var(--white);
  border: 2px solid var(--color-secundary);
}

button:active {
  position: relative;
  bottom: -3px;
}

@media (max-width: 600px) {
    h3 {
      margin-right: 5px;
      font-size: 15px;
    }

    nav {
      padding: 20px;
      display: flex;
      flex-direction: column;
    }
}

</style>
