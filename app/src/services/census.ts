import type { Person } from '@/types'
import { apiCall } from '@/services/utils'
import { ref } from 'vue'

export function useCensus() {
  const people = ref<Person[]>([])
  const totalRecords = ref(0)

  async function loadCensus(filters: Record<string, string>) {
    const query = new URLSearchParams(filters).toString()
    const baseURL = query ? `/census?${query}` : '/census'
    const data = await apiCall(baseURL)
    people.value = data.current
    totalRecords.value = data.total
  }

  return { loadCensus, people, totalRecords }
}
