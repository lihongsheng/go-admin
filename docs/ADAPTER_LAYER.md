# 中间件/适配层：API 封装与统一 Mock 方案（以资源管理为示例）

对应需求文档「中间件/适配层」交付。本文档**以资源管理为完整示例**，介绍：

1. 统一 Mock 方案如何实现（拦截器 + 注册表）
2. 需求文档要求的四个标准方法（`listResources / getResourceDetail / createResource / updateResource`）如何接入 mock
3. 如何运行时切换 Mock API 与真实 API

## 定位与职责边界

```
┌──────────┐   调 api 封装方法     ┌──────────────────┐   统一传输封装     ┌──────────┐
│  前端页面  │ ──────────────────▶ │  api/resource.js  │ ───────────────▶ │  后端 API  │
│ (UI 状态)  │ ◀────────────────── │  （API 封装层）    │ ◀─────────────── │ (业务逻辑) │
└──────────┘     返回业务数据      └──────────────────┘     标准响应结构    └──────────┘
  loading/空/失败                  Token注入/错误码转换                    建表/查询
  页面自己处理                      Mock↔真实 动态切换（统一方案）            分页/筛选
```

| 层 | 只管什么 | 不管什么 |
|---|---|---|
| 前端页面 | UI 状态（loading / 空数据 / 请求失败）、表单交互 | 不拼接 URL、不感知后端错误码 |
| API 封装层 | 接口方法（列表/详情/创建/更新/批量删除） | 不写业务逻辑、不感知 Mock 细节 |
| 拦截器（传输层） | Token 注入、错误码转友好提示、Mock/真实分发 | 不关心页面 |
| 后端 | 分页/筛选/建表/事务 | 不关心前端框架 |

## 要求.md 四接口一览

需求文档要求的四个标准方法与本项目实现的对应关系（页面唯一数据入口 `web/src/api/resource.js`）：

| 要求.md 标准方法 | API 封装 | 真实 API（后端） | Mock 注册（`api/mock/resource.js`） |
|---|---|---|---|
| `listResources` | `resourceList(params)` | `GET /api/plugin/resource/v1/list` | `mockRoute('GET', '/api/plugin/resource/v1/list', ...)` |
| `getResourceDetail` | `resourceDetail(id)` | `GET /api/plugin/resource/v1/:id` | `mockRoute('GET', '/api/plugin/resource/v1/:id', ...)` |
| `createResource` | `resourceCreate(data)` | `POST /api/plugin/resource/v1` | `mockRoute('POST', '/api/plugin/resource/v1', ...)` |
| `updateResource` | `resourceUpdate(id, data)` | `PUT /api/plugin/resource/v1/:id` | `mockRoute('PUT', '/api/plugin/resource/v1/:id', ...)` |
| （扩展）批量删除 | `resourceBatchDelete(ids)` | `DELETE /api/plugin/resource/v1/batch` | `mockRoute('DELETE', '/api/plugin/resource/v1/batch', ...)` |

**核心机制**：四个方法全部走 `request`（真实实现）。Mock 开启时，axios 拦截器按 (method, url) 命中注册表，把请求"劫持"为 mock handler 的模拟响应；关闭时请求原样放行到真实后端。**方法与 mock 一一对应，无需在方法内部写任何 mock 逻辑。**

---

## 统一 Mock 方案实现详解

### 总体调用链路（以 `createResource` 为例）

```
页面 submitAdd() → resourceCreate(data) → request.post('/api/plugin/resource/v1', data)
  → axios 请求拦截器（utils/request.js）
      ├─ ① Token 注入（Authorization: Bearer xxx）
      ├─ ② getMockEnabled() = true？
      │    ├─ mockResolve('POST', '/api/plugin/resource/v1') 命中注册表？
      │    │    ├─ 是 → 执行 mock handler（写入 mock 内存数组）→ 挂自定义 adapter 返回模拟响应
      │    │    └─ 否 → console.warn 提示 + 放行真实请求
      └─ ③ getMockEnabled() = false → 放行真实请求
  → 响应拦截器：code===0 放行 / 非 0 弹 msg / 401 跳登录 / 9001 跳安装
```

### 1. 注册表核心（`web/src/api/mock/registry.js`）

无依赖的独立模块，维护 (method, url) → handler 映射，URL 支持 `:id` 通配：

```js
const routes = []

function patternToRegExp(pattern) {
  const segs = pattern.split('/').map(s =>
    s.startsWith(':') ? '([^/]+)' : s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  )
  return new RegExp('^' + segs.join('/') + '$')
}

// 注册 mock 路由：handler({ query, body, params }) → { code, msg, data }
// 返回 undefined 表示不处理（放行真实请求）
export function mockRoute(method, pattern, handler) {
  routes.push({ method: method.toUpperCase(), pattern, re: patternToRegExp(pattern), handler })
}

// 供 request 拦截器调用：命中返回 { handler, params }，未命中返回 null
export function mockResolve(method, url) {
  const path = String(url || '').split('?')[0]      // 去掉 query
  const m = routes.find(r => r.method === method.toUpperCase() && r.re.test(path))
  if (!m) return null
  // 提取 :id 等路径参数
  const names = m.pattern.split('/').filter(s => s.startsWith(':')).map(s => s.slice(1))
  const values = m.re.exec(path).slice(1)
  const params = {}
  names.forEach((n, i) => { params[n] = values[i] })
  return { handler: m.handler, params }
}
```

### 2. 请求拦截器分发（`web/src/utils/request.js`）

Mock 命中时**挂自定义 adapter**（不能直接返回伪响应对象——axios 会把拦截器返回值当作 config 继续走 dispatchRequest，导致 transformData 抛错）：

```js
request.interceptors.request.use(async cfg => {
  const userStore = useUserStore()
  if (userStore.token) cfg.headers.Authorization = 'Bearer ' + userStore.token

  // 统一 Mock 分发：开关动态读取（运行时可切换）
  if (getMockEnabled()) {
    const hit = mockResolve(cfg.method, cfg.url)
    if (hit) {
      const body = await hit.handler({ query: cfg.params, body: cfg.data, params: hit.params })
      if (body !== undefined) {
        // 自定义 adapter 跳过网络层，返回标准 AxiosResponse
        cfg.adapter = async () => ({
          data: body, status: 200, statusText: 'OK', headers: {}, config: cfg, request: {}
        })
      }
    } else {
      console.warn('[mock] 未注册的接口，放行真实请求:', cfg.method, cfg.url)
    }
  }
  return cfg
})
```

Mock 响应 `{ code, msg, data }` 与真实后端同构，**复用响应拦截器**的统一错误处理（0 码放行、非 0 码弹提示、401 跳登录）。

### 3. handler 约定

```js
mockRoute('METHOD', '/api/path/:id', async ({ query, body, params }) => {
  await mockDelay()                    // 模拟 150-400ms 网络延迟
  // query  = GET 查询参数（axios params）
  // body   = POST/PUT 请求体（FormData 用 body.get('file') 取文件）
  // params = URL 路径参数（:id 等）
  return { code: 0, msg: 'ok', data: { ... } }   // 必须与真实后端同构
  // 返回 undefined = 不处理，放行真实请求
})
```

### 4. 实现要点（踩过的坑）

- **ESM 循环导入**：注册表核心拆到无依赖的 `registry.js`，各 mock 模块从它 import `mockRoute`，`index.js` 只做聚合——避免 `Cannot access 'routes' before initialization`
- **拦截器返回伪响应**：axios 会把 request 拦截器返回值当作 config 继续走 dispatchRequest——必须用自定义 adapter 返回 AxiosResponse
- **数据结构一致性**：mock 返回结构必须与真实后端一致（`{ code, msg, data }`），才能复用响应拦截器

---

## 四个接口的 Mock 接入（要求.md 四方法逐一示例）

以下为 `web/src/api/mock/resource.js` 真实代码，四个 handler 与要求.md 四方法一一对应。

### `listResources`（列表）的 mock

```js
mockRoute('GET', '/api/plugin/resource/v1/list', async ({ query }) => {
  await mockDelay()
  // query = 页面传入的分页/搜索参数，与真实后端解析规则一致
  const { keyword = '', type = '', status = 0, page = 1, limit = 10 } = query || {}
  const kw = String(keyword).trim().toLowerCase()
  const filtered = mockResources.filter(r =>
    (!kw || r.name.toLowerCase().includes(kw)) &&
    (!type || r.type === type) &&
    (!Number(status) || r.status === Number(status))
  )
  const start = (Number(page) - 1) * Number(limit)
  return { code: 0, msg: 'ok',
    data: { list: filtered.slice(start, start + Number(limit)), total: filtered.length } }
})
```

### `getResourceDetail`（详情）的 mock

```js
mockRoute('GET', '/api/plugin/resource/v1/:id', async ({ params }) => {
  await mockDelay()
  // params.id = URL 中的 :id（由 mockResolve 提取）
  const row = mockResources.find(r => r.id === Number(params.id))
  if (!row) return { code: 1, msg: '资源不存在' }   // 错误结构同真实后端
  return { code: 0, msg: 'ok', data: { ...row } }
})
```

### `createResource`（创建）的 mock

```js
mockRoute('POST', '/api/plugin/resource/v1', async ({ body }) => {
  await mockDelay()
  body = body || {}
  // 构造与真实后端相同的资源对象，写入内存数组（前端刷新后重置）
  const row = {
    id: nextId++, name: body.name, type: body.type, status: body.status ?? 1,
    status_at: mockNow(), created_at: mockNow(), updated_at: mockNow()
  }
  mockResources.unshift(row)
  return { code: 0, msg: 'ok', data: { ...row } }
})
```

### `updateResource`（更新）的 mock

```js
mockRoute('PUT', '/api/plugin/resource/v1/:id', async ({ params, body }) => {
  await mockDelay()
  body = body || {}
  const row = mockResources.find(r => r.id === Number(params.id))
  if (!row) return { code: 1, msg: '资源不存在' }
  Object.assign(row, { name: body.name ?? row.name, type: body.type ?? row.type, updated_at: mockNow() })
  // 与真实后端逻辑一致：仅状态变更时刷新状态更新时间
  if (body.status && row.status !== Number(body.status)) {
    row.status = Number(body.status)
    row.status_at = mockNow()
  }
  return { code: 0, msg: 'ok', data: { ...row } }
})
```

> 批量删除为扩展方法（`mockRoute('DELETE', '/api/plugin/resource/v1/batch', ...)`），不在四方法契约内。

**注册**：在 `web/src/api/mock/index.js` 底部副作用导入即可生效：

```js
import './resource'
```

---

## 如何切换 Mock API 与真实 API

切换开关只维护在 `web/src/api/mock/index.js`（仅此一层），**页面与 API 封装零改动**。

### 开关实现

```js
// 默认值：环境变量 VITE_USE_MOCK（编译期），可被运行时覆盖
const DEFAULT_MOCK = import.meta.env.VITE_USE_MOCK === 'true'
const STORAGE_KEY = 'go-admin-mock-enabled'

export function getMockEnabled() {
  const v = localStorage.getItem(STORAGE_KEY)
  if (v === null) return DEFAULT_MOCK
  return v === '1'
}

export function setMockEnabled(enabled) {
  localStorage.setItem(STORAGE_KEY, enabled ? '1' : '0')
}
```

优先级：**localStorage > 环境变量**。

### 方式一：布局右上角 Mock 开关（运行时 UI）

`web/src/layout/index.vue` 头部右侧：

```js
const mockEnabled = ref(getMockEnabled())
function onMockToggle() {
  setMockEnabled(mockEnabled.value)   // 持久化到 localStorage
  window.location.reload()            // 刷新应用新模式
}
```

| 操作 | 效果 |
|---|---|
| 打开 Mock（开关亮起） | 刷新后全系统接口（含资源管理四接口）返回模拟数据，免安装免登录 |
| 关闭 Mock | 刷新后走真实后端（无后端时进入安装向导 / 登录页） |

### 方式二：代码动态切换

```js
import { setMockEnabled } from '@/api/mock'

setMockEnabled(true)    // 切到 Mock API（后续请求立即生效，无需刷新）
setMockEnabled(false)   // 切到真实 API
```

### 方式三：环境变量（默认值）

```bash
VITE_USE_MOCK=true npm run dev     # 默认开启 Mock（仍可被运行时开关覆盖）
npm run dev                        # 默认走真实 API
```

### 切换后的请求链路对比（以 `createResource` 为例）

```
Mock API（开关开）：
  resourceCreate(data) → request.post('/api/plugin/resource/v1', data)
    → 拦截器命中 mockRoute('POST', '/api/plugin/resource/v1')
    → 写入 mock 内存数组，模拟延迟
    → 返回 { code: 0, msg: "ok", data: {...} }          ← 不发 HTTP 请求

真实 API（开关关）：
  resourceCreate(data) → request.post('/api/plugin/resource/v1', data)
    → 拦截器放行 → 注入 Bearer Token
    → Gin: POST /api/plugin/resource/v1
    → Service 层校验 + 业务 → Repo 层 GORM 写入 plugin_resource_ecs
    → 返回 { code: 0, msg: "ok", data: {...} }
```

两次返回结构完全一致，页面无需感知差异。

### 验证切换是否生效

| 方法 | 现象 |
|---|---|
| 浏览器 Network 面板 | Mock 开：请求不出现（被拦截）；Mock 关：出现真实请求（响应 200） |
| 控制台 | Mock 开：无网络请求日志；Mock 关：Network 可见 |
| 数据变化 | Mock 开：列表 5 条 mock 数据，新增后 6 条（内存）；Mock 关：真实库数据 |
| localStorage | `go-admin-mock-enabled` = `'1'`（Mock）/ `'0'`（真实） |

---

## Token 注入与错误码转换

由 `web/src/utils/request.js` axios 拦截器统一完成（API 封装层之下的传输层）：

| 处理 | 行为 |
|---|---|
| Token 注入 | 请求拦截器自动加 `Authorization: Bearer <token>` |
| 业务错误 | `code !== 0` → `ElMessage.error(msg)` 并 reject |
| 未安装 | `code === 9001` → 跳转 `/install` |
| 未授权 | `code === 401` → 清 token → 跳转 `/login` |
| Mock 分发 | Mock 开关开启时按注册表命中 → 自定义 adapter 返回模拟响应 |

## 三层全链路联调示例

### 1. 页面调用（只依赖 api 封装）

```js
// web/src/views/plugin/resource/index.vue
import { resourceList, resourceCreate } from '@/api/resource'

async function load() {
  loading.value = true
  try {
    const { data } = await resourceList({ page: 1, limit: 10, keyword: 'web' })
    list.value = data.list
    total.value = data.total
  } finally { loading.value = false }
}

async function submitAdd() {
  await resourceCreate({ name: 'web-server-02', type: '通用型 g6', status: 1 })
}
```

### 2. 页面 → 后端（真实模式）

```
resourceCreate(data)
  → request.post('/api/plugin/resource/v1', data)
  → (拦截器) 注入 Bearer Token
  → Gin: POST /api/plugin/resource/v1
  → Service 层校验 + 业务
  → Repo 层 GORM 写入 plugin_resource_ecs
  → 返回 { code: 0, msg: "ok", data: { id, name, ... } }
```

### 3. 页面 → Mock（Mock 模式，统一方案）

```
resourceCreate(data)
  → request.post('/api/plugin/resource/v1', data)
  → 拦截器命中 mock 注册表（mock/resource.js，开关已开启）
  → 内存写入 mockResources，模拟延迟 150-400ms
  → 返回与真实后端完全一致的 { code: 0, msg: "ok", data: {...} }
```

### 4. 联调验证脚本（Playwright）

```bash
VITE_USE_MOCK=true npm run dev   # 起 Mock 前端
# 浏览器自动进入系统（免安装免登录）
# 资源管理：列表 5 条 → 新增 6 条 → 勾选批量删除 4 条
# 右上角关 Mock → 刷新 → 走真实后端（本地无后端时进入安装向导/登录）
```

## 与其他模块的关系

- **统一 Mock 层**（`web/src/api/mock/`）：拦截器级中间件，覆盖包括资源管理（`mock/resource.js` 注册 `/api/plugin/resource/v1/*`）在内的全部接口；API 封装层只写真实实现，mock 由本层统一分发
- **免安装免登录**：Mock 开关开启时，路由守卫自动注入演示 token，跳过安装向导与登录页
