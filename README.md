# go-admin

基于 **Go + Gin + GORM** 与 **Vue + Element-UI** 的插件化后台管理系统。

特性：
- 启动期自动检测 DB 连通性 → 未配置 / 未安装 时进入 **在线安装向导**（`/install/*`）
- 安装向导支持 **MySQL / PostgreSQL / SQLite**，一键 AutoMigrate + 种子数据
- 用户 / 角色 / 菜单 / API 四件套，菜单树驱动前端 **动态路由**，按钮级 `v-permission`
- 插件机制：放一个目录到 `server/plugin/<name>` 即可注入 Model、菜单、API、路由，
  插件 Model 会自动参与首次安装与日常迁移；前端同步在 `web/src/plugin/<name>` 提供视图

## 目录结构

```
server/   Go 后端
  cmd/server         启动入口
  config             config.yaml + 加载/写回
  core/installer     【关键】DB 检测 + 在线初始化
  api/v1/install     /install 路由实现
  api/v1/base        登录 / 当前用户 / 当前菜单
  api/v1/system      用户/角色/菜单/API 增删改查
  initialize         路由/日志/DB/插件 初始化
  middleware         InstallGuard / JWT / Cors / RequestLog
  model/system       核心 Model
  plugin/<name>      插件包（init 内调用 plugin.Register）
web/      Vue 前端
  src/views/install  三步安装向导
  src/views/system   系统管理页面
  src/views/plugin   插件视图
  src/plugin         前端插件目录（自动 import）
```

## 运行

```bash
# 后端
cd server
go mod tidy
go run ./cmd/server -c config/config.yaml

# 前端
cd web
npm i
npm run dev          # 默认代理 /install 与 /api 到 http://127.0.0.1:8080
```

首次访问 `http://127.0.0.1:5173` 会自动跳到 `/install` 完成数据库与管理员初始化。

## 在线初始化状态机

```
读 config.yaml ─► driver 为空 ──────────────► 安装向导
                │
                └► 尝试连接 ─► 失败 ────────► 安装向导
                            └► 成功
                                 │
                            sys_install 表无记录 ──► 安装向导
                            sys_install 表有记录 ──► 正常启动
```

`/install/*` 仅在「未安装」时开放；安装完成后由 `InstallGuard` 中间件拒绝。
