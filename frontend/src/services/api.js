import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => {
    const body = response.data
    if (body && typeof body === 'object' && 'error' in body) {
      if (body.error) {
        return Promise.reject(new Error(body.error))
      }
      return body.data
    }
    return body
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    const message =
      error.response?.data?.error ||
      error.response?.data?.message ||
      error.message ||
      'Request failed'
    return Promise.reject(new Error(message))
  }
)

export const agentsApi = {
  getAll: () => api.get('/agents'),
}

export const tasksApi = {
  getAll: () => api.get('/tasks'),
  getEvents: (id) => api.get(`/tasks/${id}/events`),
  send: (agentId, title) => api.post('/tasks', { agent_id: agentId, title }),
}

export const statsApi = {
  getOverview: () => api.get('/stats'),
}

export default api
