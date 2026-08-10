/**
 * 用户管理 mock：分页列表 / 创建 / 更新 / 删除 + 头像上传
 * 数据源 mockUsers 同时被 auth.js 复用（登录用户信息保持一致）
 */
import { mockRoute, mockDelay, mockNow } from './registry'
import { mockRoles } from './role'

let nextId = 10
export const mockUsers = [
  {
    id: 1, username: 'admin', nickname: '管理员', email: 'admin@go-admin.dev', phone: '13800000001',
    avatar: '', status: 1, system_type: 0, mch_id: 0,
    roles: [{ id: 1, name: '超级管理员', default_router: '' }],
    created_at: '2026-07-01 10:00:00', updated_at: '2026-07-01 10:00:00'
  },
  {
    id: 2, username: 'ops', nickname: '运维小张', email: 'ops@go-admin.dev', phone: '13800000002',
    avatar: '', status: 1, system_type: 0, mch_id: 0,
    roles: [{ id: 2, name: '平台管理员', default_router: '/system/user' }],
    created_at: '2026-07-05 09:30:00', updated_at: '2026-07-05 09:30:00'
  },
  {
    id: 3, username: 'mch01', nickname: '商户一号', email: 'mch01@go-admin.dev', phone: '13800000003',
    avatar: '', status: 1, system_type: 1, mch_id: 1,
    roles: [{ id: 3, name: '商户管理员', default_router: '' }],
    created_at: '2026-07-10 14:00:00', updated_at: '2026-07-10 14:00:00'
  },
  {
    id: 4, username: 'guest', nickname: '访客', email: '', phone: '',
    avatar: '', status: 2, system_type: 0, mch_id: 0,
    roles: [],
    created_at: '2026-07-20 11:00:00', updated_at: '2026-07-20 11:00:00'
  },
]

export const findMockUser = (id) => mockUsers.find(u => u.id === Number(id))

// 角色 id → 角色对象（用于回显 role_ids → roles 数组）
function rolesOf(roleIds) {
  return (roleIds || []).map(rid => {
    const r = mockRoles.find(x => x.id === Number(rid))
    return r ? { id: r.id, name: r.name, default_router: r.default_router || '' } : { id: rid, name: '角色' + rid, default_router: '' }
  })
}

mockRoute('GET', '/api/v1/system/user/list', async ({ query }) => {
  await mockDelay()
  const { keyword = '', page = 1, size = 10 } = query || {}
  const kw = String(keyword).trim().toLowerCase()
  const filtered = mockUsers.filter(u =>
    !kw || u.username.toLowerCase().includes(kw) || u.nickname.toLowerCase().includes(kw)
  )
  const start = (Number(page) - 1) * Number(size)
  return {
    code: 0, msg: 'ok',
    data: { list: filtered.slice(start, start + Number(size)), total: filtered.length }
  }
})

mockRoute('POST', '/api/v1/system/user', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const row = {
    id: nextId++, username: body.username, nickname: body.nickname || '', email: body.email || '',
    phone: body.phone || '', avatar: body.avatar || '', status: body.status ?? 1,
    system_type: body.system_type ?? 0, mch_id: body.mch_id ?? null,
    roles: rolesOf(body.role_ids),
    created_at: mockNow(), updated_at: mockNow()
  }
  mockUsers.push(row)
  return { code: 0, msg: 'ok', data: row }
})

mockRoute('PUT', '/api/v1/system/user', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const u = findMockUser(body.id)
  if (!u) return { code: 1, msg: '用户不存在' }
  Object.assign(u, {
    nickname: body.nickname ?? u.nickname,
    email: body.email ?? u.email,
    phone: body.phone ?? u.phone,
    avatar: body.avatar ?? u.avatar,
    status: body.status ?? u.status,
    system_type: body.system_type ?? u.system_type,
    mch_id: body.mch_id ?? u.mch_id,
    roles: body.role_ids ? rolesOf(body.role_ids) : u.roles,
    updated_at: mockNow()
  })
  return { code: 0, msg: 'ok', data: { ...u } }
})

mockRoute('DELETE', '/api/v1/system/user/:id', async ({ params }) => {
  await mockDelay()
  const i = mockUsers.findIndex(u => u.id === Number(params.id))
  if (i === -1) return { code: 1, msg: '用户不存在' }
  mockUsers.splice(i, 1)
  return { code: 0, msg: 'ok' }
})

// 头像/文件上传：不真正存文件，返回假 URL
mockRoute('POST', '/api/v1/base/upload', async ({ body }) => {
  await mockDelay()
  let name = 'mock-file.png'
  if (body && typeof body.get === 'function') {
    const f = body.get('file')
    if (f && f.name) name = f.name
  }
  return { code: 0, msg: 'ok', data: { url: '/uploads/mock/' + encodeURIComponent(name), name } }
})
