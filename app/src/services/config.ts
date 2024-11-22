import type { Config } from '@/types'
import { apiCall } from '@/services/utils'
import { ref } from 'vue'

export function useConfig() {
  const userConfig = ref<Config>({} as Config)

  async function loadConfig() {
    const baseURL = '/user-config'
    const data = await apiCall(baseURL)
    userConfig.value = data
  }

  async function updateConfig(filters: object) {
    const baseURL = '/user-config'
    const data = await apiCall(baseURL, {
      method: 'PUT',
      data: filters
    })
    userConfig.value = data
  }

  return { loadConfig, userConfig, updateConfig }
}
