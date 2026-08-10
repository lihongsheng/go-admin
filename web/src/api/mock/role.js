/**
 * 角色管理 mock：角色 CRUD / 角色授权 / 默认首页
 * 授权结果按角色 ID 存内存（menu_ids 全量，前端会自行过滤父节点勾选叶子）
 */
import { mockRoute, mockDelay, mockNow } from './registry'
import { collectAllIds } from './menu'

let nextId = 10
export const mockRoles = [
  { id: 1, name: '超级管理员', default_router: '', remark: '拥有全部权限', status: 1, system_type: 0, mch_id: 0, created_at: '2026-07-01 10:00:00', updated_at: '2026-07-01 10:00:00' },
  { id: 2, name: '平台管理员', default_router: '/system/user', remark: '管理平台用户与角色', status: 1, system_type: 0, mch_id: 0, created_at: '2026-07-02 10:00:00', updated_at: '2026-07-02 10:00:00' },
  { id: 3, name: '商户管理员', default_router: '', remark: '商户侧管理员', status: 1, system_type: 1, mch_id: 1, created_at: '2026-07-03 10:00:00', updated_at: '2026-07-03 10:00:00' },
  { id: 4, name: '代理管理员', default_router: '', remark: '代理侧管理员', status: 1, system_type: 2, mch_id: 0, created_at: '2026-07-04 10:00:00', updated_at: '2026-07-04 10:00:00' },
]

// 角色授权结果：roleId → { menu_ids, default_router }
const roleAuthStore = new Map()

mockRoute('GET', '/api/v1/system/role/list', async ({ query }) => {
  await mockDelay()
  const sysType = Number(query && query.system_type)
  const list = sysType ? mockRoles.filter(r => r.system_type === sysType) : mockRoles
  return { code: 0, msg: 'ok', data: { list, total: list.length } }
})

mockRoute('POST', '/api/v1/system/role', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const row = {
    id: nextId++, name: body.name, default_router: body.default_router || '',
    remark: body.remark || '', status: body.status ?? 1,
    system_type: body.system_type ?? 0, mch_id: body.mch_id ?? null,
    created_at: mockNow(), updated_at: mockNow()
  }
  mockRoles.push(row)
  return { code: 0, msg: 'ok', data: row }
})

mockRoute('PUT', '/api/v1/system/role', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const r = mockRoles.find(x => x.id === Number(body.id))
  if (!r) return { code: 1, msg: '角色不存在' }
  Object.assign(r, {
    name: body.name ?? r.name,
    default_router: body.default_router ?? r.default_router,
    remark: body.remark ?? r.remark,
    status: body.status ?? r.status,
    system_type: body.system_type ?? r.system_type,
    mch_id: body.mch_id ?? r.mch_id,
    updated_at: mockNow()
  })
  return { code: 0, msg: 'ok', data: { ...r } }
})

mockRoute('DELETE', '/api/v1/system/role/:id', async ({ params }) => {
  await mockDelay()
  const i = mockRoles.findIndex(r => r.id === Number(params.id))
  if (i === -1) return { code: 1, msg: '角色不存在' }
  mockRoles.splice(i, 1)
  roleAuthStore.delete(Number(params.id))
  return { code: 0, msg: 'ok' }
})

mockRoute('GET', '/api/v1/system/role/auth/:id', async ({ params }) => {
  await mockDelay()
  const rid = Number(params.id)
  const role = mockRoles.find(r => r.id === rid)
  if (!role) return { code: 1, msg: '角色不存在' }
  const saved = roleAuthStore.get(rid) || { menu_ids: collectAllIds(), default_router: role.default_router || '' }
  return {
    code: 0, msg: 'ok',
    data: { menu_ids: saved.menu_ids, default_router: saved.default_router, system_type: role.system_type || 0 }
  }
})

mockRoute('POST', '/api/v1/system/role/auth', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const role = mockRoles.find(r => r.id === Number(body.role_id))
  if (!role) return { code: 1, msg: '角色不存在' }
  roleAuthStore.set(Number(body.role_id), { menu_ids: body.menu_ids || [], default_router: role.default_router || '' })
  return { code: 0, msg: 'ok' }
})

mockRoute('PUT', '/api/v1/system/role/:id/default-router', async ({ params, body }) => {
  await mockDelay()
  const r = mockRoles.find(x => x.id === Number(params.id))
  if (!r) return { code: 1, msg: '角色不存在' }
  r.default_router = (body && body.default_router) || ''
  return { code: 0, msg: 'ok' }
})
