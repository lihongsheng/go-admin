# CLAUDE.md

该文件可为 Claude Code（claude.ai/code）在本代码仓库中处理代码时提供操作指引。

## 项目概述

基于 Go (Gin) + Vue 3 + Element Plus 的 Web 后台管理系统，支持多商户、插件化扩展和可观测性监控。

核心功能：用户管理、角色管理、菜单管理、商户管理、API 权限管控。

### 多商户体系

系统通过 `SystemType` 枚举区分平台与商户：

| SystemType | 值 | 说明 |
|---|---|---|
| `SystemTypePlatform` | 0 | 平台（超级管理员），拥有所有权限 |
| `SystemTypeMch` | 1 | 商户，仅拥有商户管理员分配的权限 |
| `SystemTypeProxy` | 2 | 代理（预留） |

**数据隔离规则：**
- `SysUser.MchID` + `SysUser.SystemType` 标识用户归属，平台用户 `MchID=0`
- `SysRole.MchID` + `SysRole.SystemType` 标识角色归属
- `SysMenu.SystemType` 区分菜单归属平台或商户，商户管理员只能看到 `SystemTypeMch` 菜单
- 商户管理员拥有商户用户和商户角色管理权限，**不拥有**商户菜单和 API 维护权限
- JWT Token 携带 `mch_id` 和 `system_type`，前端/后端可据此做展示与鉴权区分

**安装流程初始化（`core/installer/installer.go`）：**
1. 创建平台超级管理员角色 + 平台管理员用户（`MchID=0, SystemType=Platform`）
2. 创建平台菜单树（`SystemType=Platform`）
3. 创建默认商户（`Merchant` 表）
4. 创建商户管理员角色 + 商户管理员用户（`MchID=merchant.ID, SystemType=Mch`）
5. 创建商户菜单树（`SystemType=Mch`，仅用户管理 + 角色管理 + Dashboard）
6. 写入 Casbin 策略：平台角色绑定平台菜单 API，商户角色绑定商户菜单 API

## 编码规范

本项目遵循 Go 高级开发实践，请注意以下原则：

### 依赖注入
- 所有 service 层通过构造函数注入依赖（如 `NewXxxService(repo, casbin, ...)`），不允许在内部直接 `new` 或使用全局变量拼凑
- repo / casbin / 其他外部依赖统一在 `initialize/service.go` 中装配好再注入
- 包级单例（`DefaultXxx`）只用于 wiring，业务代码不直接引用包级单例之外的全局变量

### 接口编程
- service 层必须定义接口（如 `type XxxService interface { ... }`），返回接口类型，让调用方面向接口编程
- repo 层同样定义接口，方便单测替换
- casbin 操作通过 `Port` 接口暴露，不直接导入 `casbin` 包

### 分层架构
- handler（`api/v1/`）→ service（`service/`）→ repo（`repo/`），不允许跨层调用
- handler 只做参数校验和响应组装，不写业务逻辑
- middleware 不直接操作 DB，通过 service 方法完成业务检查
- 跨层传递用 DTO（`dto/`），不直接暴露 model

### 不要硬编码
- 枚举值用 `enum/` 目录下的类型常量，不写 magic number
- 不要用 `"u:" + id` 这种前缀拼接 — casbin 策略直接使用数字 ID
- 不要用角色 code 做策略主体 — 一律用角色 ID
- 配置项走 `global.Cfg`，不硬编码在业务代码中
- 商户编号由 `genid` 自动生成（`"M" + GenDeviceID`），不手动拼接

## Commands

```bash
# 启动后端服务
cd server && go run cmd/server/main.go -c config/config.yaml

# 启动前端开发服务
cd web && npm run dev
```

## 架构

### Server 目录结构

```
server/
├── api/v1/          # Handler 层：参数校验 + 响应组装
│   ├── base/        #   登录/验证码/当前用户/菜单/上传
│   ├── install/     #   安装向导
│   └── system/      #   用户/角色/菜单/商户 CRUD
├── service/         # Service 层：业务逻辑，定义接口
│   ├── base/        #   登录/验证码/用户信息
│   ├── install/     #   安装流程
│   └── system/      #   用户/角色/菜单/商户
├── repo/            # Repo 层：数据访问，定义接口
│   └── system/      #   用户/角色/菜单/商户
├── dto/             # DTO：跨层传递的请求/响应结构
│   ├── base/
│   ├── install/
│   └── system/
├── model/system/    # Model：数据库表映射（SysUser, SysRole, SysMenu, Merchant, SysInstall）
├── enum/            # 枚举常量（UserStatus, MchStatus, SystemType）
├── middleware/       # Gin 中间件
│   ├── jwt.go       #   JWT 鉴权，解析 token 写入 context
│   ├── casbin.go    #   Casbin API 鉴权，必须放在 JWTAuth 之后
│   ├── install_guard.go  # 安装守卫，未安装时拦截业务路由
│   ├── metrics.go   #   Prometheus HTTP 指标
│   └── trace.go     #   OpenTelemetry trace 注入日志
├── router/          # 路由注册
│   ├── system.go    #   /api/v1/system/*（需 JWTAuth + CasbinAuth）
│   ├── plugin.go    #   /private/v1/plugin/* + 插件自身路由
│   └── install.go   #   /install/*（无鉴权）
├── initialize/      # 启动期初始化
│   ├── service.go   #   依赖注入装配（核心：InitDBServices）
│   ├── router.go    #   路由构造（中间件顺序：Recovery → OTel → Trace → Metrics → Cors → RequestLog）
│   ├── plugin.go    #   插件加载（空导入触发 init()）
│   ├── sync.go      #   启动期插件增量同步
│   ├── gorm.go      #   DB 连接
│   ├── casbin.go    #   Casbin 初始化
│   ├── otel.go      #   OpenTelemetry Trace
│   └── metrics.go   #   Prometheus + OTLP Metrics
├── core/installer/  # 在线安装引擎
│   ├── installer.go #   安装主流程 + seedCore（角色/用户/商户/菜单/Casbin 策略）
│   ├── seed.go      #   默认菜单树（平台 + 商户）、默认商户
│   └── source.go    #   Model/Seed 注册中心
├── plugin/          # 插件框架
│   ├── plugin.go    #   Plugin 接口 + 注册中心 + 启动期同步
│   └── example/     #   示例插件（Model/Service/Repo/Handler/菜单 一条龙）
├── config/          # 配置结构 + Viper 加载
├── global/          # 全局变量（Cfg, DB, Installed）
├── utils/           # 工具包
│   ├── casbin/      #   Casbin 封装 + Port 接口
│   ├── jwt/         #   JWT 签发/解析（User 结构含 MchID + SystemType）
│   ├── genid/       #   ID 生成器（商户编号等）
│   ├── captcha/     #   验证码（memory/redis 存储）
│   ├── upload/      #   文件上传（本地/阿里云 OSS/腾讯 COS）
│   └── response/    #   统一响应封装
├── log/             # 日志（Zap + 按日期切割）
└── cmd/server/      # 入口 main.go
```

### Web 目录结构

```
web/src/
├── api/             # API 请求封装
│   ├── base.js      #   登录/验证码/用户信息/菜单/系统类型
│   ├── system.js    #   用户/角色/菜单/商户/插件 API
│   └── install.js   #   安装向导 API
├── views/           # 页面组件
│   ├── login/       #   登录页
│   ├── dashboard/   #   仪表盘
│   ├── system/      #   系统管理（用户/角色/菜单）
│   ├── plugin/      #   插件中心
│   │   ├── list/    #     已装插件列表
│   │   ├── mch/     #     商户管理（view + form）
│   │   └── example/ #     示例插件页面
│   ├── install/     #   安装向导
│   └── error/       #   404
├── store/modules/   # Pinia 状态管理
│   ├── user.js      #   用户信息/Token/登录登出
│   └── permission.js #  动态路由 + 按钮权限码（从服务端菜单树生成）
├── router/          # Vue Router 配置
├── plugin/          # 前端插件系统
│   ├── index.js     #   自动扫描 src/plugin/*/index.js，调用 install(app, ctx)
│   └── example/     #   示例前端插件
├── directive/       # 自定义指令
│   └── permission/  #   v-permission 按钮权限控制（对比 store 中的 btns 数组）
├── composables/     # 组合式函数
├── layout/          # 布局组件（含 SubTree 递归侧边栏）
├── utils/           # 工具函数（axios 封装等）
└── permission.js    # 路由守卫（安装检测 → 登录检测 → 动态路由 → 默认首页）
```

### RBAC 权限模型

```
用户(SysUser) ←多对多→ 角色(SysRole) ←多对多→ 菜单(SysMenu)
```

- **Casbin 负责 API 级鉴权**：策略存储在 `casbin_rule` 表
  - `p` 策略格式：`p(role_id, /api/path, method)` — 基于角色 ID，不用角色名
  - `g` 策略格式：`g(user_id, role_id)` — 用户-角色关联
- **菜单三层结构**：
  - `catalog`：目录节点，`component=Layout`，仅作路由容器
  - `menu`：菜单节点，对应一个前端页面，`ApiRules` 字段关联 API 路径
  - `button`：按钮节点，作为 menu 的子节点，`permission` 字段承载权限码（如 `user:add`），由前端 `v-permission` 指令控制显隐
- **API 规则注入**：菜单 `ApiRules` 字段为 JSON 数组 `[{"path":"/api/...","method":"GET"},...]`，角色授权时自动写入 Casbin 策略

### 插件系统

**后端插件契约（`plugin.Plugin` 接口）：**
- `Name()` / `Version()` — 唯一标识
- `Models()` — 参与 AutoMigrate 的 Model 列表
- `Menus()` — 注入菜单树（含 catalog/menu/button；按 `Name` 幂等 upsert），API 规则通过 `ApiRules` 字段注入
- `RegisterRoute(g *gin.Engine)` — 注册自身路由
- `SeedTable(db)` — 首次安装时的种子数据（仅目标表为空时执行）

**插件生命周期：**
1. 插件包 `init()` 调用 `plugin.Register(p)` 自注册
2. `initialize.LoadPlugins()` 通过空导入触发
3. 启动期 `plugin.SyncOnBoot()` 执行：AutoMigrate → 幂等 upsert 菜单/API → 条件 SeedTable
4. 安装向导完成时，插件 Seed 注册到 `installer` 一并执行

**前端插件：**
- 在 `src/plugin/<name>/index.js` 导出 `install(app, ctx)` 函数
- `src/plugin/index.js` 通过 `import.meta.glob` 自动扫描加载

### 路由与中间件链

```
Request → Recovery → OTel(可选) → Trace → Metrics(可选) → Cors → RequestLog
        → InstallGuard → [业务路由]
```

- `/health` — 健康检查，无中间件
- `/metrics` — Prometheus 指标端点，InstallGuard 之前
- `/uploads/*` — 静态文件，长缓存
- `/install/*` — 安装向导，无鉴权
- `/api/v1/base/*` — 登录/验证码/上传，仅需 JWTAuth
- `/api/v1/system/*` — 系统管理，需 JWTAuth + CasbinAuth
- `/private/v1/plugin/*` — 插件列表，需 JWTAuth + CasbinAuth

### 可观测性

**OpenTelemetry Trace（`initialize/otel.go`）：**
- 支持导出器：`stdout`（默认）、`otlp`（gRPC/HTTP）、`none`（仅写日志）
- 通过 `otelgin` 中间件自动注入 span context
- 配置项：`observability.trace.enable / exporter / endpoint / service_name / sample_rate`

**Prometheus Metrics（`initialize/metrics.go`）：**
- 始终暴露 `/metrics` 端点供 Prometheus 拉取
- 可选 OTLP 推送到 collector
- HTTP 指标中间件记录请求计数、延迟、错误率
- 配置项：`observability.metrics.enable / endpoint / path`

### 安装向导

系统首次启动时 DB 未配置，进入安装模式：
1. 前端 `/install` 页面提交 DB 配置 + 管理员账号
2. 后端 `installer.Install()` 执行：建库 → AutoMigrate → Casbin 初始化 → seedCore → 插件 Seed → 写入 `sys_install`
3. 安装完成后回调 `InitDBServices()` 装配所有 service
4. `InstallGuard` 中间件此后放行所有业务路由
