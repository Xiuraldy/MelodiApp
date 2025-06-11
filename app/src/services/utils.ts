const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

interface Options {
  method?: string
  data?: any
  headers?: any
}
type ApiCallOptions = Options | undefined

export async function apiCall(
  path: string,
  { method = 'GET', data, headers }: ApiCallOptions = {}
) {
  const sessionToken = sessionStorage.getItem('token')
  const response = await fetch(BASE_URL + path, {
    method,
    mode: 'cors',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${sessionToken}`,
      ...headers
    },
    body: JSON.stringify(data)
  })

  if (method === 'DELETE') {
    return response.ok
  }

  const jsonObj = await response.json()
  return jsonObj
}
