import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: { 'Content-Type': 'application/json' }
})

// Request interceptor
api.interceptors.request.use(
  config => {
    const token = localStorage.getItem('token')
    if (token) config.headers.Authorization = `Bearer ${token}`
    return config
  },
  error => Promise.reject(error)
)

// Response interceptor
api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error.response?.data || error.message)
    return Promise.reject(error)
  }
)

// Agents
export const agentApi = {
  getAll: () => api.get('/agents'),
  getById: (id) => api.get(`/agents/${id}`),
  sendTask: (agentId, task) => api.post(`/agents/${agentId}/task`, { task }),
}

// Tasks
export const taskApi = {
  getAll: (params) => api.get('/tasks', { params }),
  getById: (id) => api.get(`/tasks/${id}`),
  create: (data) => api.post('/tasks', data),
  getHistory: () => api.get('/tasks/history'),
}

// Stats / Overview
export const statsApi = {
  getOverview: () => api.get('/stats/overview'),
  getTokenUsage: (days = 7) => api.get(`/stats/tokens?days=${days}`),
}

// Telegram
export const telegramApi = {
  getFeed: (limit = 20) => api.get(`/telegram/feed?limit=${limit}`),
  sendMessage: (agentId, message) => api.post('/telegram/send', { agent_id: agentId, message }),
}

export default api
