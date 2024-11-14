<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useUsers } from '@/services/users'
import type { User } from '@/types';
import { useRouter } from 'vue-router';

const { loadUsers, deleteUser, createUser, editUser, logout, users } = useUsers()

const router = useRouter()
const newUser = reactive({
  username: '',
  email: '',
  password: ''
})

const editMode = ref(false)
const currentUser = reactive({
  ID: null as number | null,
  username: '',
  email: '',
  password: ''
})

function setEditUser(user: User) {
  currentUser.ID = user.ID
  currentUser.username = user.username
  currentUser.email = user.email
  currentUser.password = ''
  editMode.value = true
}

async function saveEditUser() {
  if (currentUser.ID) {
    await editUser(currentUser.ID, {
      username: currentUser.username,
      email: currentUser.email,
      password: currentUser.password || undefined 
    })
    editMode.value = false
    loadUsers()
  }
}

onMounted(loadUsers)

</script>
<template>
  <section>
    <h1>Users</h1>
    <ul>
      <li v-for="user in users" :key="user.ID">
        {{ user.ID }}. {{ user.username }} - {{ user.firstname }} {{ user.lastname }}
        <button @click="() => deleteUser(user.ID)">Borrar</button>
        <button @click="() => setEditUser(user)">Editar</button>
      </li>
    </ul>
    
    <h2 v-if="!editMode">Add User</h2>
    <form v-if="!editMode" @submit.prevent="() => createUser(newUser)">
      <input v-model="newUser.username" type="text" placeholder="Username" />
      <input v-model="newUser.email" type="email" placeholder="Email" />
      <input v-model="newUser.password" type="password" placeholder="Password" />
      <button type="submit">Add User</button>
    </form>

    <h2 v-else>Edit User</h2>
    <form v-if="editMode" @submit.prevent="saveEditUser">
      <input v-model="currentUser.username" type="text" placeholder="Username" />
      <input v-model="currentUser.email" type="email" placeholder="Email" />
      <input v-model="currentUser.password" type="password" placeholder="New Password" />
      <button type="submit">Save Changes</button>
      <button type="button" @click="editMode = false">Cancel</button>
    </form>
  </section>

  <button @click="async() => {
    await logout()
  }">Cerrar sesión</button>
  
</template>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
