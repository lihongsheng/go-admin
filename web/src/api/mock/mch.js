/**
 * 商户管理 mock：分页列表 / 创建 / 更新 / 详情 / 状态切换
 * 仅平台用户在用户/角色表单中搜索商户时使用；商户管理页面也走此 mock
 */
import { mockRoute, mockDelay, mockNow } from './registry'

let nextId = 10
const mockMchs = [
  { id: 1, mch_no: 'M10001', mch_name: '示例科技（上海）', linker: '张三', phone: '13800000011', email: 'a@demo.com', status: 1, address: '上海市浦东新区', reason: '首个演示商户', created_at: '2026-07-01 10:00:00', updated_at: '2026-07-01 10:00:00' },
  { id: 2, mch_no: 'M10002', mch_name: '示例网络（杭州）', linker: '李四', phone: '13800000022', email: 'b@demo.com', status: 1, address: '杭州市西湖区', reason: '', created_at: '2026-07-08 10:00:00', updated_at: '2026-07-08 10:00:00' },
  { id: 3, mch_no: 'M10003', mch_name: '示例传媒（北京）', linker: '王五', phone: '13800000033', email: 'c@demo.com', status: 2, address: '北京市朝阳区', reason: '已停用', created_at: '2026-07-15 10:00:00', updated_at: '2026-07-20 10:00:00' },
]

mockRoute('GET', '/api/v1/system/mch/list', async ({ query }) => {
  await mockDelay()
  const { mch_no = '', mch_name = '', status = 0, page = 1, limit = 10 } = query || {}
  const filtered = mockMchs.filter(m =>
    (!mch_no || m.mch_no.includes(mch_no)) &&
    (!mch_name || m.mch_name.includes(mch_name)) &&
    (!Number(status) || m.status === Number(status))
  )
  const start = (Number(page) - 1) * Number(limit)
  return { code: 0, msg: 'ok', data: { list: filtered.slice(start, start + Number(limit)), total: filtered.length } }
})

mockRoute('POST', '/api/v1/system/mch', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const row = {
    id: nextId++, mch_no: 'M' + (10000 + nextId), mch_name: body.mch_name || '',
    linker: body.linker || '', phone: body.phone || '', email: body.email || '',
    status: body.status ?? 1, address: body.address || '', reason: body.reason || '',
    created_at: mockNow(), updated_at: mockNow()
  }
  mockMchs.push(row)
  return { code: 0, msg: 'ok', data: row }
})

mockRoute('PUT', '/api/v1/system/mch', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const m = mockMchs.find(x => x.id === Number(body.id))
  if (!m) return { code: 1, msg: '商户不存在' }
  Object.assign(m, {
    mch_name: body.mch_name ?? m.mch_name,
    linker: body.linker ?? m.linker,
    phone: body.phone ?? m.phone,
    email: body.email ?? m.email,
    status: body.status ?? m.status,
    address: body.address ?? m.address,
    reason: body.reason ?? m.reason,
    updated_at: mockNow()
  })
  return { code: 0, msg: 'ok', data: { ...m } }
})

mockRoute('GET', '/api/v1/system/mch/:id', async ({ params }) => {
  await mockDelay()
  const m = mockMchs.find(x => x.id === Number(params.id))
  if (!m) return { code: 1, msg: '商户不存在' }
  return { code: 0, msg: 'ok', data: { ...m } }
})

mockRoute('GET', '/api/v1/system/mch/no/:mchNo', async ({ params }) => {
  await mockDelay()
  const m = mockMchs.find(x => x.mch_no === params.mchNo)
  if (!m) return { code: 1, msg: '商户不存在' }
  return { code: 0, msg: 'ok', data: { ...m } }
})

mockRoute('PUT', '/api/v1/system/mch/status', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const m = mockMchs.find(x => x.id === Number(body.id))
  if (!m) return { code: 1, msg: '商户不存在' }
  m.status = body.status ?? m.status
  m.updated_at = mockNow()
  return { code: 0, msg: 'ok', data: { ...m } }
})
