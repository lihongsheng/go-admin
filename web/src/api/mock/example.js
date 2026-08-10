/**
 * 示例插件（笔记管理）mock：列表 / 新增 / 删除
 * 覆盖 views/plugin/example/index.vue 直接调用的三个接口
 */
import { mockRoute, mockDelay, mockNow } from './registry'

let nextId = 10
const notes = [
  { id: 1, title: 'Hello', content: '示例插件初始化笔记' },
  { id: 2, title: 'Mock 演示', content: '整体 Mock 模式下笔记数据' },
]

mockRoute('GET', '/api/plugin/example/v1/note/list', async () => {
  await mockDelay()
  return { code: 0, msg: 'ok', data: { list: [...notes], total: notes.length } }
})

mockRoute('POST', '/api/plugin/example/v1/note', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const row = { id: nextId++, title: body.title || '', content: body.content || '', created_at: mockNow() }
  notes.unshift(row)
  return { code: 0, msg: 'ok', data: row }
})

mockRoute('DELETE', '/api/plugin/example/v1/note/:id', async ({ params }) => {
  await mockDelay()
  const i = notes.findIndex(n => n.id === Number(params.id))
  if (i === -1) return { code: 1, msg: '笔记不存在' }
  notes.splice(i, 1)
  return { code: 0, msg: 'ok' }
})
