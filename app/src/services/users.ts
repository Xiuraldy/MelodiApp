import type { User } from '@/types'
import { apiCall } from '@/services/utils'
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

export function useUsers() {
  const users = ref<User[]>([])
  const router = useRouter()
  const authStore = useAuthStore()

  async function loadUsers() {
    console.log('authStore.isLoggedIn', authStore.isLoggedIn)
    const data = await apiCall('/users')
    users.value = data as User[]
  }

  async function createUser(newUser: { username: string; email: string; password: string }) {
    const createdUser = await apiCall('/users', {
      method: 'POST',
      data: newUser
    })
    users.value = [...users.value, createdUser]
  }

  async function editUser(
    id: number,
    updatedUser: { username: string; email: string; password?: string }
  ) {
    const editedUser = await apiCall(`/users/${id}`, {
      method: 'PUT',
      data: updatedUser
    })
    users.value = users.value.map((user) => (user.ID === id ? editedUser : user))
  }

  async function deleteUser(id: number) {
    const deleted = await apiCall(`/users/${id}`, {
      method: 'DELETE'
    })
    if (deleted) {
      users.value = users.value.filter((user) => user.ID !== id)
    }
  }

  async function logout() {
    try {
      await apiCall('/auth/logout', {
        method: 'DELETE'
      })
      authStore.clearSession()
    } catch (error) {
      console.error('Error during logout:', error)
    }
  }

  return { loadUsers, createUser, deleteUser, editUser, logout, users }
}
