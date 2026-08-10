/**
 * 插件中心 mock：已装插件列表
 */
import { mockRoute, mockDelay } from './registry'

mockRoute('GET', '/api/plugin/list', async () => {
  await mockDelay()
  const list = [
    { name: 'example', version: '0.1.0' },
    { name: 'resource', version: '0.1.0' },
  ]
  return { code: 0, msg: 'ok', data: { list, total: list.length } }
})
