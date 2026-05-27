import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 60000,
  headers: { 'Content-Type': 'application/json' }
})

api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    const msg = error.response?.data?.message || error.message || '请求失败'
    console.error('API Error:', msg)
    if (typeof window !== 'undefined') {
      const event = new CustomEvent('app:toast', { detail: { type: 'error', message: msg } })
      window.dispatchEvent(event)
    }
    return Promise.reject(error)
  }
)

export default api
