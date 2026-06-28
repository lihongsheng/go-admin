# CLAUDE.md

该文件为 Claude Code 在本仓库中处理代码时提供操作指引。

## 项目概述

基于 Go + Vue 3 + Element Plus 的 Web 后台管理系统。核心功能：用户管理、角色管理、菜单管理、API 权限管理。支持插件化扩展（后端 API/菜单自动注册、前端页面热加载）。集成 OpenTelemetry 分布式追踪与 Prometheus 指标监控。

> main 分支为单租户版本；多商户功能由 `merchant` 分支维护。

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.22+ / Gin / GORM / Casbin / JWT / Viper |
| 前端 | Vue 3 (Composition API) / Element Plus / Pinia / Vue Router 4 |
| 可观测性 | OpenTelemetry (Trace) / Prometheus (Metrics) / Zap (Log) |
| 数据库 | MySQL / PostgreSQL / SQLite（通过 config.yaml 切换） |

## 目录结构

```
server/
├── cmd/server/main.go        # 入口
├── config/                    # 配置结构体 + Load/Save
├── global/                    # 进程级共享对象（Cfg/DB/Redis/Installed）—— 尽量少用，新代码走注入
├── initialize/                # 启动期装配：DB/Casbin/Redis/Router/Service/Plugin/Sync
├── api/v1/                    # Handler 层：参数校验 + 响应组装
├── service/                   # Service 层：业务逻辑，定义接口 + 实现
├── repo/                      # Repo 层：数据访问，定义接口 + 实现
├── dto/                       # 跨层传输对象（请求/响应）
├── model/                     # GORM Model（不直接暴露给 handler）
├── enum/                      # 枚举常量
├── middleware/                 # Gin 中间件（JWT/Casbin/CORS/Trace/Metrics/InstallGuard）
├── router/                    # 路由注册
├── plugin/                    # 插件接口 + 注册中心
│   └── example/               # 示例插件（Model/Repo/Service/API/DTO 完整结构）
├── core/installer/            # 在线安装向导（建库/建表/种子数据）
├── utils/
│   ├── casbin/                # Casbin Port 接口 + 实现
│   ├── jwt/                   # JWT 工具
│   ├── captcha/               # 验证码（Memory/Redis 存储）
│   ├── upload/                # 文件上传（Local/AliyunOSS/TencentCOS）
│   ├── genid/                 # ID 生成
│   └── response/              # 统一响应封装
└── log/                       # 日志接口（Logger + Valuer + Zap 实现 + Context 传递）

web/
├── src/
│   ├── api/                   # 后端 API 封装（base.js / system.js / install.js）
│   ├── views/                 # 页面组件
│   │   ├── system/            # 系统管理（user/role/menu）
│   │   ├── plugin/            # 插件页面
│   │   ├── install/           # 安装向导
│   │   ├── login/             # 登录
│   │   └── dashboard/         # 仪表盘
│   ├── store/modules/         # Pinia Store（user / permission）
│   ├── router/                # 路由（常驻路由 + 服务端菜单动态注册）
│   ├── layout/                # 布局框架
│   ├── directive/permission/  # v-permission 按钮权限指令
│   ├── plugin/                # 前端插件加载器（扫描 plugin/*/index.js）
│   ├── composables/           # 组合式函数
│   ├── utils/request.js       # Axios 封装（JWT 自动注入）
│   └── permission.js          # 路由守卫（安装检测 → 登录 → 动态路由 → 按钮权限）
└── config.yaml                # 前端无独立配置，走后端 API
```

## 编码规范

### 依赖注入（核心原则）

**后端代码尽量使用注入方式，不使用全局变量。** `global` 包仅用于启动期 wiring，业务代码不得直接引用。

- Service 层通过构造函数注入依赖：`NewXxxService(repo, casbinPort, ...)`，不允许在内部 `new` 或引用 `global.DB` 等全局变量
- Repo 层通过构造函数接收 `*gorm.DB`：`NewUserRepo(db)`，不直接读 `global.DB`
- 所有依赖在 `initialize/service.go` 中统一装配，赋值给包级单例（`DefaultXxx`）
- Handler 通过包级单例调用 Service（`serviceSys.DefaultUser.Create(req)`），单例仅做 wiring 入口
- Casbin 操作通过 `utils/casbin.Port` 接口暴露，Service 层持有 `Port` 接口，不直接 import casbin 包
- Logger 通过 `log.Logger` 接口注入，不直接使用 `log.Global()` 或 `log.Info()` 等包级便捷函数（中间件/初始化代码除外）

### 接口编程

- Service 层必须定义接口（`type XxxService interface`），返回接口类型
- Repo 层必须定义接口（`type XxxRepo interface`），方便单测 mock 替换
- 跨包调用面向接口，不依赖具体实现

### 分层架构

```
Handler (api/v1/) → Service (service/) → Repo (repo/)
```

- **Handler**：只做参数绑定（`ShouldBindJSON/ShouldBindQuery`）+ 调用 Service + 组装响应，不写业务逻辑
- **Service**：业务逻辑核心，通过注入的 Repo/CasbinPort 等完成操作
- **Repo**：纯数据访问，只操作 GORM，不含业务判断
- **Middleware**：不直接操作 DB，通过 Service 方法完成业务检查
- **跨层传递用 DTO**（`dto/`），不直接暴露 Model 给 Handler

### 不要硬编码

- 枚举值用 `enum/` 下的类型常量，不写 magic number
- Casbin 策略直接使用数字 ID，不用 `"u:" + id` 前缀拼接
- Casbin 策略主体一律用角色 ID，不用角色 code
- 配置项走 `global.Cfg`（启动期注入），不硬编码在业务代码中
- 菜单类型使用 `model/system` 中定义的常量（`MenuTypeCatalog` / `MenuTypeMenu` / `MenuTypeButton`）

### 前端规范

- 使用 Vue 3 Composition API（`<script setup>`），不用 Options API
- 状态管理用 Pinia，不用 Vuex
- API 调用统一走 `src/api/` 封装，不直接在组件中写 axios
- 按钮权限用 `v-permission="'xxx:yyy'"` 指令，权限码由后端菜单 `type=button` 节点提供
- 动态路由由 `permission store` 根据服务端菜单自动注册，不在 `router/index.js` 中硬编码业务路由
- 前端插件放 `src/plugin/<name>/`，导出 `install(app, ctx)` 函数，由 `plugin/index.js` 自动扫描加载

## 插件开发

### 后端插件

1. 在 `server/plugin/<name>/` 下创建包，实现 `plugin.Plugin` 接口：
   - `Name() / Version() / Models() / Menus() / RegisterRoute(g) / SeedTable(db)`
2. 在包的 `init()` 中调用 `plugin.Register(p{})` 自注册
3. 菜单树通过 `Menus()` 返回，button 权限节点作为 menu 的 Children（`type=button`）
4. API 规则通过菜单 `ApiRules` 字段注入（JSON 数组），自动写入 Casbin 策略
5. 启动期自动执行：AutoMigrate → 幂等 upsert 菜单/API → 条件 SeedTable
6. 插件内部遵循分层：`model/ → repo/ → service/ → api/` + `dto/`

### 前端插件

1. 在 `web/src/plugin/<name>/` 下创建目录
2. 导出 `index.js`，默认导出 `install(app, ctx)` 函数
3. 页面组件放 `web/src/views/plugin/<name>/`，由菜单管理配置组件路径

## 启动流程

```
main.go
 ├─ config.Load()              → global.Cfg
 ├─ initialize.Logger()        → log.SetGlobal()
 ├─ initialize.InitOpenTelemetry()
 ├─ initialize.InitMetrics()
 ├─ initialize.InitInstallService()   ← 不依赖 DB
 ├─ initialize.GormConnect()          → global.DB（失败则进入安装模式）
 │   └─ initialize.DetectInstalled()  → global.Installed
 │   └─ initialize.SetupCasbin()
 ├─ initialize.InitRedis()
 ├─ initialize.InitCaptcha()
 ├─ initialize.LoadPlugins()          → plugin.Register() → installer 注册
 ├─ initialize.SyncOnBoot()           → AutoMigrate + upsert 菜单/API + SeedTable
 ├─ initialize.InitDBServices()       → 装配所有 Service 单例（依赖 global.DB）
 └─ initialize.Router()               → Gin 路由 + 中间件
```

安装向导完成后，install service 回调 `InitDBServices()` 重新装配所有依赖 DB 的 Service。

## 常用命令

```bash
# 后端
cd server && go run cmd/server/main.go -c config/config.yaml

# 前端
cd web && npm run dev

# 代码生成（utils/gen.go）
cd server && go run utils/gen.go
```
