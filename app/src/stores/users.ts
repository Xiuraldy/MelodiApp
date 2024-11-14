import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'

export const useUsersStore = defineStore('users', () => {
  const users = ref<User[]>([])

  function setUsers(input: User[]) {
    users.value = input
  }

  return { setUsers, users }
})
