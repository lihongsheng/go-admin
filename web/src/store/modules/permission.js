import { defineStore } from 'pinia'
import { shallowRef } from 'vue'
import { menu } from '@/api/base'
import { RouterView } from 'vue-router'

// 动态 import 所有 views 下 Vue 文件
const modules = import.meta.glob('../../views/**/*.vue')

// catalog 节点仅作为路由容器，透传 <router-view />，不能返回 Layout
// 否则 Vue Router 嵌套匹配时会渲染双层 sidebar
function loadView(component) {
  if (!component) return RouterView
  if (component === 'Layout') return RouterView
  const path = `../../views/${component}.vue`
  return modules[path] || (() => import('@/views/error/404.vue'))
}

// 将服务端菜单节点转为 Vue Router 路由配置
// type=button 的节点不进入路由
function toRoute(node) {
  const routeChildren = (node.children || [])
    .filter(c => c.type !== 'button')
    .map(toRoute)

  // 去掉开头的 / ，因为要作为子路由挂到 Layout (/ ) 下面
  // 服务端返回的根级路径如 "/system"，子级如 "user"
  const r = {
    path: node.path ? node.path.replace(/^\//, '') : '',
    name: node.name,
    component: loadView(node.component),
    meta: {
      title: node.title,
      icon: node.icon,
      hidden: node.hidden,
      keepAlive: node.keep_alive,
      type: node.type
    },
    children: routeChildren.length > 0 ? routeChildren : undefined
  }

  if (node.type === 'catalog' && node.redirect) {
    r.redirect = node.redirect.replace(/^\//, '')
  } else if (node.type === 'catalog' && routeChildren.length > 0) {
    r.redirect = routeChildren[0].path
  }

  return r
}

// 递归收集所有 type=button 节点的权限码
function collectButtonCodes(nodes) {
  const codes = []
  for (const n of nodes) {
    if (n.type === 'button' && n.permission) {
      codes.push(n.permission)
    }
    if (n.children) {
      codes.push(...collectButtonCodes(n.children))
    }
  }
  return codes
}

export const usePermissionStore = defineStore('permission', () => {
  const routes = shallowRef([])
  const btns = shallowRef([])

  async function generateRoutes() {
    const { data } = await menu()
    const serverMenus = data.menus || []

    // 服务端菜单 → 前端路由（过滤 type=button）
    routes.value = serverMenus.map(toRoute)

    // 收集按钮权限码
    btns.value = collectButtonCodes(serverMenus)

    return routes.value
  }

  return { routes, btns, generateRoutes }
})
