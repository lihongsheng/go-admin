# 整体 Mock 模式

全系统接口返回模拟数据的演示模式：**免后端、免安装、免登录**，打开即用。mock 数据统一维护在 `web/src/api/mock/`，页面代码零改动。

## 原理

mock 分为两层：

```
请求 → axios 请求拦截器
        ├─ getMockEnabled() = true
        │    ├─ 命中注册表？ → 自定义 adapter 直接返回模拟响应 { code, msg, data }
        │    └─ 未命中     → 放行真实请求（console.warn 提示）
        └─ getMockEnabled() = false → 放行真实后端
```

### 1. 拦截器中间件（整体 mock）

`web/src/api/mock/registry.js` 维护 (method, url) → handler 注册表（URL 支持 `:id` 参数），`web/src/utils/request.js` 拦截器按 `getMockEnabled()` 动态分发。

```js
// 注册 mock 路由（示例：mock/auth.js）
mockRoute('POST', '/api/v1/base/login', async ({ body }) => {
  // ...校验 + 返回 { code, msg, data: { token } }
})
```

覆盖模块：

| 文件 | 覆盖接口 |
|---|---|
| auth.js | 验证码（SVG 图片显示数字）/ 登录 / 用户信息 / 登出 / 安装状态 |
| user.js | 用户分页 CRUD + 头像上传（数据源与登录信息一致） |
| role.js | 角色 CRUD + 角色授权 + 默认首页 |
| menu.js | 完整菜单树（驱动动态路由与按钮权限码）+ 菜单 CRUD |
| mch.js / plugin.js / example.js | 商户 / 插件中心 / 示例插件 |
| resource.js | 资源管理（覆盖 api/resource.js 四个方法） |

### 2. 统一方案覆盖资源管理

资源管理同样走注册表：`mock/resource.js` 注册了 `/api/plugin/resource/v1/*` 全部路由（list / :id / POST / PUT / batch），**覆盖 api/resource.js 的四个方法**（resourceList / resourceDetail / resourceCreate / resourceUpdate）及批量删除。mock 开启时，页面经 API 封装发起的请求被拦截器分发到 mock handler；关闭时直达真实后端——**页面代码零改动**。

## 免安装 / 免登录

`web/src/permission.js` 路由守卫：

```js
// Mock 开启且无 token → 自动注入演示 token，跳过登录页
if (getMockEnabled() && !userStore.token) {
  userStore.setToken('mock-token-1')
}
```

- `/install/status`、`/install/check-db` mock 返回 `installed: true` → 跳过安装向导
- `fetchInfo` 返回 mock 管理员（admin / 超级管理员）→ 正常进入系统

## 如何开启（3 种方式）

| 方式 | 操作 | 生效时机 | 持久化 |
|---|---|---|---|
| ① 环境变量（默认值） | `VITE_USE_MOCK=true npm run dev` | 启动时 | 无（每次启动生效） |
| ② 运行时 UI | 布局右上角 **Mock 开关** | 点击立即生效 | localStorage |
| ③ 代码 | `setMockEnabled(true / false)` | 调用后生效 | localStorage |

优先级：**localStorage > env 默认值**（`getMockEnabled()` 先读 localStorage，无则用 env）。

验证是否生效：

```js
localStorage.getItem('go-admin-mock-enabled')  // '1' = Mock 开，'0' = 关
// 控制台应看到 [plugin] example loaded / [plugin] resource loaded
// 右上角 Mock 开关亮起
```

## 动态切换

布局右上角 **Mock 开关**（`web/src/layout/index.vue`）：

```js
const mockEnabled = ref(getMockEnabled())
function onMockToggle() {
  setMockEnabled(mockEnabled.value)   // localStorage 持久化
  window.location.reload()            // 刷新应用新模式
}
```

| 操作 | 效果 |
|---|---|
| 打开 Mock | 刷新后免安装免登录进入系统，全接口返回模拟数据 |
| 关闭 Mock | 刷新后走真实后端（无后端时进入安装向导 / 登录页） |

> 登录页没有布局，无法用开关切回——重新打开方式：直接访问系统后恢复（或 `localStorage.setItem('go-admin-mock-enabled', '1')` 后刷新）。

## 启动方式

```bash
# 开发模式
cd web
VITE_USE_MOCK=true npm run dev        # 打开 http://localhost:5173

# Docker（构建时注入开关）
docker compose up -d web-mock         # http://localhost:8080（无需后端/数据库）
```

`web/Dockerfile` 支持构建参数：

```dockerfile
ARG VITE_USE_MOCK=false
RUN VITE_USE_MOCK=$VITE_USE_MOCK npm run build:prod
```

`docker-compose.yml` 中 `web-mock` 服务：

```yaml
web-mock:
  build:
    context: ./web
    args:
      VITE_USE_MOCK: "true"
  ports:
    - "81:80"
```

## 未注册接口行为

Mock 模式下未注册的接口会**放行真实请求**并在控制台提示：

```
[mock] 未注册的接口，放行真实请求: get /api/xxx
```

## 新增页面如何接入 Mock

新页面（如订单管理）有两条路径，按需选择：

### 路径 A：API 封装模式（推荐）

```js
// web/src/api/order.js —— 接口方法封装（走 request）
import request from '@/utils/request'

export const orderList = (params) => request.get('/api/plugin/order/v1/list', { params })
export const orderCreate = (data) => request.post('/api/plugin/order/v1', data)
// ...
```

页面只 import `api/order.js` 的方法，**mock 由统一方案覆盖**：在 `api/mock/` 注册对应路由即可（参考 [ADAPTER_LAYER.md](ADAPTER_LAYER.md)）。

### 路径 B：拦截器注册表（已有 api 封装的模块快速 mock）

```js
// api/order.js 照常写（走 request）
export const orderList = (params) => request.get('/api/plugin/order/v1/list', { params })

// mock/order.js 注册 mock —— 页面代码零改动
mockRoute('GET', '/api/plugin/order/v1/list', async ({ query }) => {
  await mockDelay()
  return { code: 0, msg: 'ok', data: { list: [], total: 0 } }
})

// mock/index.js 底部导入
import './order'
```

## 给新增接口添加 Mock

三步即可让新接口被 mock 覆盖（页面代码零改动）：

**1. 在对应模块文件注册路由**（如 `mock/auth.js`）：

```js
import { mockRoute, mockDelay } from './registry'

mockRoute('GET', '/api/v1/system/user/:id', async ({ params }) => {
  await mockDelay()                       // 模拟 150-400ms 网络延迟
  const u = mockUsers.find(x => x.id === Number(params.id))
  if (!u) return { code: 1, msg: '用户不存在' }
  return { code: 0, msg: 'ok', data: { ...u } }
})
```

**2. handler 约定**：

| 入参 | 说明 |
|---|---|
| `{ query }` | GET 查询参数（axios `params`） |
| `{ body }` | 请求体（POST/PUT JSON；FormData 可用 `body.get('file')`） |
| `{ params }` | URL 路径参数（`:id` 等） |

返回值必须是统一结构 `{ code, msg, data }`——与真实后端一致，响应拦截器统一处理（0 码放行、非 0 码弹提示、401 跳登录）。返回 `undefined` 表示不处理、放行真实请求。

**3. 新建模块文件时**，在 `mock/index.js` 底部副作用导入：

```js
import './your-module'
```

### 示例：为 `/api/v1/system/mch/no/:mchNo` 添加 mock

```js
mockRoute('GET', '/api/v1/system/mch/no/:mchNo', async ({ params }) => {
  await mockDelay()
  const m = mockMchs.find(x => x.mch_no === params.mchNo)
  if (!m) return { code: 1, msg: '商户不存在' }
  return { code: 0, msg: 'ok', data: { ...m } }
})
```

> 统一方案：所有模块（含资源管理）都在注册表维护 mock，无需在适配层重复实现。

## 实现要点备忘

- **ESM 循环导入**：注册表核心拆到无依赖的 `registry.js`，各 mock 模块从它 import `mockRoute`，`index.js` 只做聚合，避免 `Cannot access 'routes' before initialization`
- **拦截器返回伪响应**：axios 会把 request 拦截器返回值当作 config 继续走 dispatchRequest，导致 transformData 抛错——正确做法是命中时挂**自定义 adapter** 返回 AxiosResponse
- **数据结构一致性**：mock 返回结构必须与真实后端一致（`{ code, msg, data }`），复用响应拦截器的错误处理（401 跳登录 / 非 0 码弹提示）
