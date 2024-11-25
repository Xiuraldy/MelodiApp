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

  const response = await fetch('http://54.207.125.158:8080/auth/login', {
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
  email: '',
  password: ''
})

async function register() {

  const response = await fetch('http://54.207.125.158:8080/auth/register', {
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
    <div class="content-form">
      <form v-if="activeLogin" @submit.prevent="signIn">
        <h1>Login</h1>
        <div class="content-inputs-button">
          <input v-model="loginInputs.email" type="text" placeholder="Email" />
          <input v-model="loginInputs.password" type="password" placeholder="Password" />
          <div class="error" v-if="loginError">⚠︎ Error: {{ loginError }}</div>
          <button>Login</button>
        </div>
        <p>You don't have a user?<a @click="activeLogin = false">Register</a></p>
      </form>
  
      <form v-if="!activeLogin" @submit.prevent="register">
        <h1>Register</h1>
        <div class="content-inputs-button">
          <input v-model="registerInputs.username" type="text" placeholder="Username" />
          <input v-model="registerInputs.email" type="text" placeholder="Email" />
          <input v-model="registerInputs.password" type="password" placeholder="Password" />
          <div class="error" v-if="registerError">⚠︎ Error: {{ registerError }}</div>
          <button>Register</button>
        </div>
        <p>Do you already have a user?<a @click="activeLogin = true">Login</a></p>
      </form>
    </div>
  </main>
</template>

<style scoped>

h1 {
  width: -webkit-fill-available;
  background-color: #fff;
  display: flex;
  justify-content: center;
  color: var(--color-secundary);
  font-weight: 200;
  border-top-left-radius: 40px;
  border-top-right-radius: 40px;
  border: 2px solid var(--color-secundary);

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
  justify-content: space-between;
  align-items: center;
  background-color: var(--color-secundary);
  height: 70vh;
  width: 50vh;
  border-radius: 50px;
}

template {
  font-family: poppins;
}

.content-inputs-button {
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: -webkit-fill-available;
  padding: 30px;
}

input, a {
  font-family: poppins;
  font-size: medium;
  outline: none;
  border-radius: 5px;
  border: 2px solid #fff;
  background-color: transparent;
  color: #fff;
  cursor: pointer;
  margin: 5px;
}

input {
  cursor: text;
  padding-left: 10px;
}

input::placeholder {
  color: #fff;
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
  background-color: var(--color-primary);
  width: -webkit-fill-available;
  color: #fff;
  border-bottom-left-radius: 50px;
  border-bottom-right-radius: 50px;
}

a {
  font-size: initial;
  border: 2px solid #fff;
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
}
</style>