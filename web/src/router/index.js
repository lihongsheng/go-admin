import { createRouter, createWebHistory } from 'vue-router'

// 常驻路由（不走权限）
export const constantRoutes = [
  { path: '/install', name: 'Install', component: () => import('@/views/install/index.vue'), meta: { title: '安装向导' } },
  { path: '/login',   name: 'Login',   component: () => import('@/views/login/index.vue'),   meta: { title: '登录' } },
  { path: '/404',     name: '404',     component: () => import('@/views/error/404.vue') },
  {
    path: '/', name: 'Layout', component: () => import('@/layout/index.vue'),
    children: [
      { path: 'dashboard', name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'dashboard' } }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes
})

export default router
