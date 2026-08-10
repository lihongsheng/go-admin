/**
 * 菜单管理 mock：完整菜单树（驱动动态路由与权限码）+ 菜单 CRUD
 * 菜单树结构与真实种子数据保持一致，保证 mock 模式下前端路由/按钮权限照常工作
 */
import { mockRoute, mockDelay, mockNow } from './registry'

export const menuTreeData = [
  {
    id: 1, type: 'menu', path: '/dashboard', name: 'Dashboard', component: 'dashboard/index',
    title: '仪表盘', icon: 'dashboard', sort: 1, hidden: false, keep_alive: false
  },
  {
    id: 2, type: 'catalog', path: '/system', name: 'System', component: 'Layout',
    title: '系统管理', icon: 'setting', sort: 10,
    children: [
      {
        id: 3, type: 'menu', path: 'user', name: 'SysUser', component: 'system/user/index',
        title: '用户管理', icon: 'user', sort: 1, hidden: false, keep_alive: false,
        children: [
          { id: 31, type: 'button', name: '新增用户', permission: 'user:add' },
          { id: 32, type: 'button', name: '编辑用户', permission: 'user:edit' },
          { id: 33, type: 'button', name: '删除用户', permission: 'user:del' },
        ]
      },
      {
        id: 4, type: 'menu', path: 'role', name: 'SysRole', component: 'system/role/index',
        title: '角色管理', icon: 'peoples', sort: 2, hidden: false, keep_alive: false,
        children: [
          { id: 41, type: 'button', name: '新增角色', permission: 'role:add' },
          { id: 42, type: 'button', name: '编辑角色', permission: 'role:edit' },
          { id: 43, type: 'button', name: '删除角色', permission: 'role:del' },
          { id: 44, type: 'button', name: '角色授权', permission: 'role:auth' },
        ]
      },
      {
        id: 5, type: 'menu', path: 'menu', name: 'SysMenu', component: 'system/menu/index',
        title: '菜单管理', icon: 'tree-table', sort: 3, hidden: false, keep_alive: false,
        children: [
          { id: 51, type: 'button', name: '新增菜单', permission: 'menu:add' },
          { id: 52, type: 'button', name: '编辑菜单', permission: 'menu:edit' },
          { id: 53, type: 'button', name: '删除菜单', permission: 'menu:del' },
        ]
      },
    ]
  },
  {
    id: 6, type: 'catalog', path: '/plugin', name: 'PluginCenter', component: 'Layout',
    title: '插件中心', icon: 'component', sort: 90,
    children: [
      {
        id: 7, type: 'menu', path: 'list', name: 'PluginList', component: 'plugin/list/index',
        title: '已装插件', icon: 'list', sort: 1, hidden: false, keep_alive: false
      },
    ]
  },
  {
    id: 8, type: 'menu', path: '/plugin/resource', name: 'PluginResource', component: 'plugin/resource/index',
    title: '资源管理', icon: 'monitor', sort: 92, hidden: false, keep_alive: false,
    children: [
      { id: 81, type: 'button', name: '新增资源', permission: 'resource:add' },
      { id: 82, type: 'button', name: '编辑资源', permission: 'resource:edit' },
      { id: 83, type: 'button', name: '批量删除资源', permission: 'resource:del' },
    ]
  },
]

let nextMenuId = 100

// 递归收集所有节点 id（角色授权默认全选）
export function collectAllIds(nodes = menuTreeData) {
  const ids = []
  const walk = (list) => {
    for (const n of list) {
      ids.push(n.id)
      if (n.children) walk(n.children)
    }
  }
  walk(nodes)
  return ids
}

function findNode(nodes, id) {
  for (const n of nodes) {
    if (n.id === Number(id)) return n
    if (n.children) {
      const hit = findNode(n.children, id)
      if (hit) return hit
    }
  }
  return null
}

function removeNode(nodes, id) {
  const i = nodes.findIndex(n => n.id === Number(id))
  if (i > -1) { nodes.splice(i, 1); return true }
  for (const n of nodes) {
    if (n.children && removeNode(n.children, id)) return true
  }
  return false
}

mockRoute('GET', '/api/v1/system/menu/tree', async () => {
  await mockDelay()
  // 深拷贝返回，避免页面直接改到 mock 数据
  return { code: 0, msg: 'ok', data: { list: JSON.parse(JSON.stringify(menuTreeData)) } }
})

// 登录后的菜单树（permission store 用它生成动态路由与按钮权限码）
mockRoute('GET', '/api/v1/base/menu', async () => {
  await mockDelay()
  return { code: 0, msg: 'ok', data: { menus: JSON.parse(JSON.stringify(menuTreeData)) } }
})

mockRoute('POST', '/api/v1/system/menu', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const row = {
    id: nextMenuId++,
    parent_id: body.parent_id || 0,
    type: body.type || 'menu',
    title: body.title || '', name: body.name || '',
    path: body.path || '', component: body.component || '',
    permission: body.permission || '', redirect: body.redirect || '',
    icon: body.icon || '', sort: body.sort ?? 0,
    hidden: !!body.hidden, keep_alive: !!body.keep_alive,
    api_rules: body.api_rules || '',
    created_at: mockNow(), updated_at: mockNow()
  }
  if (row.parent_id) {
    const parent = findNode(menuTreeData, row.parent_id)
    if (parent) {
      if (!parent.children) parent.children = []
      parent.children.push(row)
    } else {
      return { code: 1, msg: '父菜单不存在' }
    }
  } else {
    menuTreeData.push(row)
  }
  return { code: 0, msg: 'ok', data: row }
})

mockRoute('PUT', '/api/v1/system/menu', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const m = findNode(menuTreeData, body.id)
  if (!m) return { code: 1, msg: '菜单不存在' }
  Object.assign(m, {
    parent_id: body.parent_id ?? m.parent_id,
    type: body.type ?? m.type,
    title: body.title ?? m.title,
    name: body.name ?? m.name,
    path: body.path ?? m.path,
    component: body.component ?? m.component,
    permission: body.permission ?? m.permission,
    redirect: body.redirect ?? m.redirect,
    icon: body.icon ?? m.icon,
    sort: body.sort ?? m.sort,
    hidden: !!body.hidden, keep_alive: !!body.keep_alive,
    api_rules: body.api_rules ?? m.api_rules,
    updated_at: mockNow()
  })
  return { code: 0, msg: 'ok', data: { ...m } }
})

mockRoute('DELETE', '/api/v1/system/menu/:id', async ({ params }) => {
  await mockDelay()
  if (!removeNode(menuTreeData, params.id)) return { code: 1, msg: '菜单不存在' }
  return { code: 0, msg: 'ok' }
})
