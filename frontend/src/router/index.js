import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/DashboardView.vue'

const routes = [
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: DashboardView,
    meta: { title: 'DASHBOARD' }
  },
  {
    path: '/agents',
    name: 'Agents',
    component: () => import('@/views/AgentsView.vue'),
    meta: { title: 'AGENTS' }
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: () => import('@/views/TasksView.vue'),
    meta: { title: 'TASKS' }
  },
  {
    path: '/telegram',
    name: 'Telegram',
    component: () => import('@/views/TelegramView.vue'),
    meta: { title: 'TELEGRAM' }
  },
  {
    path: '/memory',
    name: 'Memory',
    component: () => import('@/views/MemoryView.vue'),
    meta: { title: 'MEMORY' }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/SettingsView.vue'),
    meta: { title: 'SETTINGS' }
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.afterEach((to) => {
  document.title = `${to.meta.title || 'AI Empire'} — AI Empire`
})

export default router
