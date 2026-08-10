# go-admin

基于 **Go + Gin + GORM** 与 **Vue 3 + Element Plus** 的插件化后台管理系统示例工程。

以「资源管理」为完整链路示例，覆盖需求文档（`prd/要求.md`）中的**交互视觉 / 前端开发 / 后端开发 / 中间件适配层 / 系统设计**五类交付：

- **后端**：Go 资源管理插件（云 ECS），多条件分页 / 详情 / 新增 / 编辑 / 批量删除
- **前端**：Vue 3 资源管理页面（搜索 / 筛选 / 分页 / 右侧抽屉表单 / 顶部多标签页）
- **中间件/适配层**：资源管理 API 封装（列表/详情/创建/更新），Mock / 真实后端按统一开关**运行时动态切换**
- **整体 Mock 模式**：免后端、免安装、免登录，一键演示全系统

## ✨ 特性

- **在线安装向导**：启动自动检测 DB，未配置 / 未安装时进入 `/install` 三步向导（MySQL / PostgreSQL / SQLite）
- **RBAC 权限四件套**：用户 / 角色 / 菜单 / API，菜单树驱动前端**动态路由**，按钮级 `v-permission`
- **插件机制**：`server/plugin/<name>` 一个目录注入 Model / 菜单 / API / 路由，前端 `web/src/plugin/<name>` 提供视图，Model 自动参与迁移与种子数据
- **资源管理插件**（`server/plugin/resource`）：云 ECS 资源管理，分页搜索 / 详情抽屉 / 新增编辑抽屉 / 批量删除，按钮级权限 `resource:add` / `resource:edit` / `resource:del`
- **中间件/适配层**（`web/src/api/resource.js` + `api/mock/`）：资源列表/详情/创建/更新四接口封装，Mock 由统一方案（拦截器注册表）按开关动态分发
- **整体 Mock 模式**（`web/src/api/mock/`）：axios 拦截器统一分发全系统模拟数据，支持**运行时动态切换**（布局右上角 Mock 开关），免安装、免登录直接演示
- **前端体验增强**：顶部多标签页（首页常驻、打开菜单自动追加、可关闭）、表单统一右侧抽屉弹出
- **Swagger API 文档**：注解驱动自动生成（`make swagger`），`/swagger/index.html` 在线调试
- **可观测性**：OpenTelemetry 分布式追踪 / Prometheus 指标 / Zap 日志

## 🏗️ 架构

```
┌──────────────────────────────────────────────────────────────────┐
│  前端页面（views/）                                                │
│   · UI 状态自管：loading / 空数据 / 请求失败 / 正常                  │
│   · 只调用适配层方法，不直接拼接后端 URL                              │
└───────────────┬──────────────────────────────────────────────────┘
                │  import { resourceList, resourceCreate } from '@/api/resource'
┌───────────────▼──────────────────────────────────────────────────┐
│  中间件/适配层（web/src/api/resource.js + api/mock/）              │
│   · 四接口封装：列表 / 详情 / 创建 / 更新（走 request）              │
│   · Mock 由统一方案（拦截器注册表）覆盖，开关 getMockEnabled()       │
│     运行时动态切换（仅此层维护）                                    │
└───────────────┬──────────────────────────────────────────────────┘
                │  axios（拦截器：Token 注入 / 错误码转友好提示 / 整体 Mock 分发）
┌───────────────▼──────────────────────────────────────────────────┐
│  后端（Go / Gin / GORM）                                          │
│  /api/plugin/resource/v1/* → Service 层（业务） → Repo 层（GORM）   │
└──────────────────────────────────────────────────────────────────┘
```

**职责边界**：

| 层 | 只管什么 | 不管什么 |
|---|---|---|
| 前端页面 | UI 状态（加载/空/失败）、表单交互 | 不拼 URL、不感知后端结构 |
| 中间件/适配层 | 4 方法契约、Token 注入、错误码转换、Mock/真实切换 | 不写业务逻辑 |
| 后端 | 分页/筛选/建表/事务 | 不关心前端框架 |

## 🚀 快速开始

### 方式一：Mock 模式（推荐，无需后端）

```bash
cd web
npm i
VITE_USE_MOCK=true npm run dev
```

打开 `http://localhost:5173` —— **免安装、免登录直接进入系统**（admin 演示账号已注入）。页面右上角有 **Mock 开关**，运行时随时切换真实后端。

### 方式二：真实模式

```bash
# 后端（首次访问自动进安装向导，配置 MySQL 后初始化）
cd server
go mod tidy
go run ./cmd/server -c config/config.yaml

# 前端
cd web
npm i
npm run dev
```

### 方式三：Docker Compose（一键）

```bash
docker compose up -d          # 正式模式：web(:80) + server + mysql
docker compose up -d web-mock # Mock 模式前端：http://localhost:8080（免后端）
```

## 📖 Swagger API 文档

接口注解（`@Summary / @Param / @Router`）写在 Handler 上，`make swagger` 自动生成到 `server/docs/`：

```bash
make swagger        # 生成 swagger.json / swagger.yaml / docs.go（依赖 swag，首次先 make swag）
```

- 访问地址（真实模式）：`http://localhost:8989/swagger/index.html`
- 资源管理插件接口已带注解：`/api/plugin/resource/v1/*`（列表/详情/新增/编辑/批量删除）
- Docker 部署下 nginx 已代理 `/swagger/` → 后端

## 🔌 插件开发

放一个目录即可扩展一个完整功能模块（Model / 菜单 / API / 路由 / 页面自动注册）：

```bash
server/plugin/<name>/        # 后端插件（model → dto → repo → service → api → plugin.go）
web/src/views/plugin/<name>/ # 前端页面（菜单 Component 指向，动态路由自动生效）
```

- 实现 `plugin.Plugin` 接口（`Name/Version/Models/Menus/InitServices/RegisterRoute/SeedTable`）
- 菜单 `Menus()` 声明目录/菜单/按钮节点，`ApiRules` 自动写入 Casbin 按钮级权限
- 完整步骤与代码示例见 [docs/PLUGIN_DEVELOPMENT.md](docs/PLUGIN_DEVELOPMENT.md)

## 🧪 Mock 模式详解

**开启方式**（优先级：localStorage > 环境变量）：

```bash
VITE_USE_MOCK=true npm run dev   # ① 环境变量（默认值）
# ② 布局右上角 Mock 开关（运行时点击，立即生效并持久化）
# ③ setMockEnabled(true | false) 代码控制
```

**新增页面接入**——两条路径：

- **API 封装**（推荐）：`api/xxx.js` 封装接口方法（走 request），页面只 import api 封装，mock 由统一方案按 URL 覆盖
- **拦截器模式**：`api/xxx.js` 照常封装，再在 `api/mock/xxx.js` 注册 `mockRoute(...)`（页面零改动）

**新增接口接入**——三步：

```js
// 1. 在对应模块注册路由（如 mock/user.js）
mockRoute('GET', '/api/v1/system/user/:id', async ({ params }) => {
  await mockDelay()
  // ...返回 { code, msg, data }，结构同真实后端；返回 undefined 则放行真实请求
})

// 2. 新建模块时在 mock/index.js 底部副作用导入
import './your-module'
```

完整讲解（两层架构 / 三条开启路径 / 页面接入示例 / handler 约定）见 [docs/MOCK_MODE.md](docs/MOCK_MODE.md)。

## 📁 目录结构

```
server/   Go 后端
├── cmd/server           启动入口
├── config/              config.yaml（含 config.docker.yaml）
├── core/installer       在线安装向导（建库/建表/种子数据）
├── api/v1/              Handler 层（install/base/system）
├── service/ repo/ dto/ model/ enum/   分层架构
├── middleware/          InstallGuard / JWT / Casbin / CORS / Trace
├── initialize/          启动装配（DB/Casbin/Redis/Router/Plugin/Sync）
├── plugin/              插件：example（示例）/ resource（资源管理）
│   └── resource/        model → dto → repo → service → api → plugin.go
└── utils/               casbin / jwt / upload / response 等

web/      Vue 3 前端
├── src/views/           页面（system / plugin / install / login / dashboard）
├── src/api/             API 封装 + 中间件/适配层
│   ├── resource.js           资源管理 API 封装（列表/详情/创建/更新/批量删除）
│   └── mock/                整体 Mock 层（registry 注册表 + 各模块 mock）
├── src/layout/          布局（侧边栏 + 顶部多标签页 TagsView）
├── src/store/           Pinia（user / permission / tags）
└── src/utils/request.js axios 封装（Token 注入 / 错误码 / Mock 分发）
```

## 📚 模块文档

| 文档 | 内容 |
|---|---|
| [docs/RESOURCE_MANAGEMENT.md](docs/RESOURCE_MANAGEMENT.md) | 资源管理模块：表结构 / 接口清单 / 页面操作 / curl 示例 |
| [docs/ADAPTER_LAYER.md](docs/ADAPTER_LAYER.md) | 中间件/适配层：API 封装 / 统一 Mock 方案 / 动态开关 / 三层联调示例 |
| [docs/MOCK_MODE.md](docs/MOCK_MODE.md) | 整体 Mock 模式：拦截器原理 / 免安装免登录 / 动态切换 / Docker 部署 |
| [docs/FRONTEND_FEATURES.md](docs/FRONTEND_FEATURES.md) | 前端增强：顶部多标签页 / 右侧抽屉表单 |
| [docs/PLUGIN_DEVELOPMENT.md](docs/PLUGIN_DEVELOPMENT.md) | 插件开发指南：Plugin 接口 / 菜单权限 / 前后端步骤 |
| [docs/OBSERVABILITY.md](docs/OBSERVABILITY.md) | OpenTelemetry 追踪 / Prometheus 指标 |
| [docs/OPENTELEMETRY_TRACEID.md](docs/OPENTELEMETRY_TRACEID.md) | TraceID 跨进程传递 |
| [docs/WITH_CONTEXT_USAGE.md](docs/WITH_CONTEXT_USAGE.md) | Context 传参与日志链路 |

## 🗂️ 需求文档对照

| 要求.md 交付物 | 实现位置 |
|---|---|
| 资源管理页面（搜索/筛选/分页/详情/新增编辑） | `web/src/views/plugin/resource/index.vue` |
| 后端资源接口（分页/详情/创建/更新） | `server/plugin/resource/` |
| 资源管理四接口（列表/详情/创建/更新 + 统一 Mock 切换） | `web/src/api/resource.js` + `api/mock/resource.js` |
| 环境切换配置说明 | [docs/ADAPTER_LAYER.md](docs/ADAPTER_LAYER.md) |
| 三层全链路联调示例 | [docs/ADAPTER_LAYER.md](docs/ADAPTER_LAYER.md) |
| Mermaid 架构图 + 职责边界 | 本文档架构节 |

## 🔑 常用命令

```bash
# 后端
cd server && go run ./cmd/server -c config/config.yaml

# Swagger 文档生成（依赖 swag，首次先 make swag）
make swagger

# 前端开发（真实模式）
cd web && npm run dev

# 前端开发（Mock 模式，免后端免登录）
cd web && VITE_USE_MOCK=true npm run dev

# Docker
docker compose up -d               # 正式模式（:80）
docker compose up -d web-mock      # Mock 模式前端（:8080）
```
