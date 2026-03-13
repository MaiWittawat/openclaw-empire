import axios from 'axios'

// Base API instance - connects to your FastAPI backend
const api = axios.create({
  baseURL: import.meta.env.VIT_SERVER_URL || 'http://localhost:8212/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor - attach auth token if available
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor - handle errors globally
api.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error.message)
  }
)

// ── AGENTS ──
export const agentsApi = {
  getAll: () => api.get('/agents'),
  getById: (id) => api.get(`/agents/${id}`),
  updateStatus: (id, status) => api.patch(`/agents/${id}/status`, { status }),
}

// ── TASKS ──
export const tasksApi = {
  getAll: (params) => api.get('/tasks', { params }),
  getById: (id) => api.get(`/tasks/${id}`),
  create: (data) => api.post('/tasks', data),
  cancel: (id) => api.delete(`/tasks/${id}`),
  send: (agentId, task) => api.post('/tasks/send', { agent_id: agentId, task }),
}

// ── STATS ──
export const statsApi = {
  getOverview: () => api.get('/stats/overview'),
  getTokenUsage: (days = 7) => api.get('/stats/tokens', { params: { days } }),
  getSystemInfo: () => api.get('/stats/system'),
}

// ── TELEGRAM ──
export const telegramApi = {
  getFeed: (limit = 50) => api.get('/telegram/feed', { params: { limit } }),
  sendMessage: (agentId, message) =>
    api.post('/telegram/send', { agent_id: agentId, message }),
}

// ── MEMORY ──
export const memoryApi = {
  getAll: () => api.get('/memory'),
  search: (query) => api.get('/memory/search', { params: { q: query } }),
}

export default api
