import type { Entry } from '@/types'
import { apiCall } from '@/services/utils'
import { ref } from 'vue'

export function useEntries() {
  const entries = ref<Entry[]>([])

  async function loadEntries() {
    const data = await apiCall('/meEntries')
    entries.value = data as Entry[]
  }

  async function createEntry(newEntry: { title: string; content: string }) {
    const createdEntry = await apiCall('/entries', {
      method: 'POST',
      data: newEntry
    })
    entries.value = [...entries.value, createdEntry]
  }

  async function editEntry(id: number, updatedEntry: { title: string; content: string }) {
    const editedUser = await apiCall(`/entries/${id}`, {
      method: 'PUT',
      data: updatedEntry
    })
    entries.value = entries.value.map((entry) => (entry.ID === id ? editedUser : entry))
  }

  async function deleteEntry(id: number) {
    const deleted = await apiCall(`/entry/${id}`, {
      method: 'DELETE'
    })
    if (deleted) {
      entries.value = entries.value.filter((entry) => entry.ID !== id)
    }
  }

  return { loadEntries, createEntry, deleteEntry, editEntry, entries }
}
