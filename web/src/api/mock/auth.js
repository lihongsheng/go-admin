/**
 * 认证 mock：验证码 / 登录 / 用户信息 / 登出 / 安装状态
 * 登录用户信息与 mock/user.js 的 mockUsers 保持一致
 * 演示约定：用户名必须存在（默认 admin），密码非空即可，验证码以图片显示
 */
import { mockRoute, mockDelay } from './registry'
import { mockUsers } from './user'

// 验证码：id → code，SVG 图片直接画出 code，用户照着输入
const captchaStore = new Map()

function randCode() {
  return String(Math.floor(1000 + Math.random() * 9000))
}

function svgDataUri(code) {
  const colors = ['#409EFF', '#E6A23C', '#67C23A', '#F56C6C']
  const chars = code.split('').map((ch, i) =>
    `<text x="${16 + i * 20}" y="27" font-size="22" font-weight="bold" fill="${colors[i % 4]}" ` +
    `transform="rotate(${(Math.random() * 24 - 12).toFixed(0)} ${16 + i * 20} 27)">${ch}</text>`
  ).join('')
  const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="104" height="38">` +
    `<rect width="104" height="38" fill="#f5f7fa" rx="4"/>` +
    `<line x1="0" y1="${8 + Math.random() * 22}" x2="104" y2="${8 + Math.random() * 22}" stroke="#e4e7ed" stroke-width="1"/>` +
    chars + `</svg>`
  return 'data:image/svg+xml;base64,' + btoa(svg)
}

mockRoute('GET', '/api/v1/base/captcha', async () => {
  await mockDelay()
  const id = 'mock-' + Date.now() + '-' + Math.floor(Math.random() * 1e6)
  const code = randCode()
  captchaStore.set(id, code)
  return { code: 0, msg: 'ok', data: { captcha_id: id, captcha_b64: svgDataUri(code) } }
})

mockRoute('POST', '/api/v1/base/login', async ({ body }) => {
  await mockDelay()
  body = body || {}
  const { username, password, captcha_id, captcha_code } = body
  if (captcha_id && captchaStore.get(captcha_id) !== captcha_code) {
    return { code: 1, msg: '验证码错误' }
  }
  const user = mockUsers.find(u => u.username === username)
  if (!user) return { code: 1, msg: '用户不存在' }
  if (user.status === 2) return { code: 1, msg: '用户已被禁用' }
  if (!password) return { code: 1, msg: '密码不能为空' }
  // Mock 约定：密码任意非空即可登录
  return { code: 0, msg: 'ok', data: { token: 'mock-token-' + user.id } }
})

mockRoute('GET', '/api/v1/base/info', async () => {
  await mockDelay()
  const admin = mockUsers.find(u => u.id === 1) || mockUsers[0]
  return { code: 0, msg: 'ok', data: { ...admin } }
})

mockRoute('POST', '/api/v1/base/logout', async () => {
  await mockDelay()
  return { code: 0, msg: 'ok' }
})

// 安装状态：mock 模式下视为已安装，跳过安装向导
mockRoute('GET', '/install/status', async () => {
  await mockDelay()
  return { code: 0, msg: 'ok', data: { installed: true } }
})

// 安装向导的 DB 探测：mock 模式下直接返回已安装，避免进入安装页
mockRoute('GET', '/install/check-db', async () => {
  await mockDelay()
  return {
    code: 0, msg: 'ok',
    data: { configured: true, connected: true, installed: true, driver: 'mysql', reason: '' }
  }
})
