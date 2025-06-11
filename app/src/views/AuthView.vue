<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const loginInputs = reactive({
  email: '',
  password: ''
})

const loginError = ref('') 
const registerError = ref('') 


async function signIn() {

  const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(loginInputs),
  })

  const data = await response.json()

  if (data.error) {
    loginError.value = data.error 
    return
  }

  authStore.setSession(data.token)
}

const registerInputs = reactive({
  username: '',
  lastname: '',
  email: '',
  password: '',
  celphone: '',
  role: ''
})

async function register() {

  const response = await fetch('http://localhost:8080/api/v1/auth/register', {
    method: 'POST',
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(registerInputs),
  })

  const data = await response.json()

  if (data.error) {
    registerError.value = data.error 
    return
  }

  authStore.setSession(data.token)
}

const activeLogin = ref(true);

</script>

<template>
  <main>
    <img v-if="!activeLogin" src="/auth/group 2.jpg" alt="">
    <div class="content-form">
      <form v-if="activeLogin" @submit.prevent="signIn">
        <h1>Inicio de Sesión</h1>
        <div class="content-inputs-button">
          <input v-model="loginInputs.email" type="text" placeholder="Correo Electrónico" />
          <input v-model="loginInputs.password" type="password" placeholder="Contraseña" />
          <button>Inicia Sesión</button>
          <div class="error" v-if="loginError">⚠︎ Error: {{ loginError }}</div>
        </div>
        <p>¿No tienes usuario?<a @click="activeLogin = false">Registrate</a></p>
      </form>
  
      <form v-if="!activeLogin" @submit.prevent="register">
        <h1>Registro</h1>
        <div class="content-inputs-button">
          <input v-model="registerInputs.username" type="text" placeholder="Nombre" />
          <input v-model="registerInputs.lastname" type="text" placeholder="Apellido" />
          <input v-model="registerInputs.email" type="text" placeholder="Correo Electrónico" />
          <input v-model="registerInputs.celphone" type="text" placeholder="Número de Celular" />
          <select v-model="registerInputs.role" type="text">
            <option value="" disabled selected hidden>Soy...</option>
            <option value="singer">Cantante</option>
            <option value="guitarist">Guitarrista</option>
            <option value="electricGuitarist">Guitarrista Eléctrico</option>
            <option value="pianist">Pianista</option>
            <option value="saxophonist">Saxofonista</option>
            <option value="drummer">Baterista</option>
            <option value="bassist">Bajista</option>
          </select>
          <input v-model="registerInputs.password" type="password" placeholder="Contraseña" />
          <button>Registrate</button>
          <div class="error" v-if="registerError">⚠︎ Error: {{ registerError }}</div>
        </div>
        <p>¿Tienes usuario?<a @click="activeLogin = true">Inicia Sesión</a></p>
      </form>
    </div>
    <img v-if="activeLogin" src="/auth/group.jpg" alt="">
  </main>
</template>

<style scoped>

h1 {
  width: -webkit-fill-available;
  display: flex;
  justify-content: flex-start;
  color: var(--black);
  font-weight: 600;
}

main {
  display: flex;
  justify-content: space-evenly;
  flex-wrap: wrap;
}

img {
  height: 70vh;
  border-radius: 50px;
}

.content-form {
  display: flex;
  align-items: center;
  flex-direction: column;
  justify-content: center;
}

form {
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 70vh;
  width: 50vh;
  border-radius: 50px;
  justify-content: space-evenly;
}

template {
  font-family: poppins;
}

.content-inputs-button {
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: -webkit-fill-available;
}

input, a, select {
  appearance: none;
  font-family: poppins;
  font-size: medium;
  outline: none;
  border-radius: 5px;
  border: 2px solid var(--black);
  color: var(--black);
  background-color: transparent;
  cursor: pointer;
  margin: 5px;
}

input {
  cursor: text;
  padding-left: 10px;
}

input::placeholder {
  color: var(--black);
}

select {
  padding-left: 10px;
}


button {
  border-radius: 5px;
  border: none;
  background-color: var(--color-tertiary);
  padding: 5px;
  color: #fff;
  cursor: pointer;
  margin: 5px;
}

button:active {
  position: relative;
  bottom: -3px;
}

p {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: -webkit-fill-available;
  color: var(--black);
}

a {
  font-size: initial;
  border: 2px solid var(--black);
  padding-right: 6px;
  padding-left: 6px;
  margin-top: -1px;
}

.error {
  display: flex;
  transition: all 0.3s ease;
  margin: 5px;
  justify-content: space-around;
  background-color: #fff;
  color: var(--color-error);
  margin-bottom: -5px;
  border-radius: 5px 0px;
}

@media (max-width: 600px) {
  form {
    width: -webkit-fill-available;
    margin: 30px;
  }

  img {
    display: none;
  }
}
</style>