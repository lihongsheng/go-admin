# 插件开发指南

在 `server/plugin/<name>/` 放一个目录即可注入 Model / 菜单 / API / 路由；前端在 `web/src/plugin/<name>` 提供视图。参考实现：`server/plugin/example`（最小示例）、`server/plugin/resource`（完整分层 + 权限示例）。

## 后端插件结构

```
server/plugin/<name>/
├── model/          GORM Model（表名统一 plugin_<name>_ 前缀）
├── dto/            请求/响应 DTO
├── repo/           Repo 接口 + 实现（构造注入 *gorm.DB）
├── service/        业务接口 + 实现（构造注入 Repo），包级单例 DefaultXxx
├── api.go          Handler 层（参数绑定 + 响应组装 + swagger 注解）
└── plugin.go       Plugin 接口实现 + init() 自注册
```

### Plugin 接口（server/plugin/plugin.go）

| 方法 | 职责 |
|---|---|
| `Name() string` | 唯一名（路由前缀 / 表前缀建议同名） |
| `Version() string` | 版本号 |
| `Models() []interface{}` | 参与 AutoMigrate 的 Model（自动注册到 installer） |
| `Menus() []system.SysMenu` | 注入菜单树（catalog/menu/button，按 Name 幂等 upsert） |
| `InitServices(ctx InitContext) error` | 装配服务层（DB/Redis/Config 通过 InitContext 注入，不依赖 global） |
| `RegisterRoute(g, privatePlugin)` | 注册自身路由（已挂在 `/api/plugin/<name>` 下，带 JWT+Casbin） |
| `SeedTable(db) error` | 初始数据（仅当插件表为空时执行） |

### 菜单树与按钮权限

```go
func (p) Menus() []system.SysMenu {
	return []system.SysMenu{
		{
			Type:      system.MenuTypeMenu,
			Path:      "/plugin/resource",
			Name:      "PluginResource",          // 幂等 upsert 的唯一键
			Component: "plugin/resource/index",   // 前端页面组件路径
			Title:     "资源管理",
			Icon:      "monitor",
			Sort:      92,
			ApiRules:  `[{"path":"/api/plugin/resource/v1/list","method":"GET"}]`,
			Children: []system.SysMenu{
				{Type: system.MenuTypeButton, Name: "新增资源", Permission: "resource:add",
					ApiRules: `[{"path":"/api/plugin/resource/v1","method":"POST"}]`},
			},
		},
	}
}
```

- button 节点作为 menu 的 Children（`type=button`），`Permission` 为按钮权限码（前端 `v-permission` 用）
- `ApiRules` 为 JSON 数组，启动期由 `collectApiRules` **递归收集整棵菜单树**写入 Casbin（含 button 子节点），超级管理员自动挂载

### 注册

```go
// plugin.go 末尾
func init() { plugin.Register(p{}) }

// server/initialize/plugin.go 空导入触发自注册
_ "go-admin/server/plugin/resource"
```

## 前端插件

```
web/src/plugin/<name>/index.js     前端插件入口（install(app, ctx)，可选）
web/src/views/plugin/<name>/index.vue   页面组件（菜单 Component 字段指向它）
```

- 页面由服务端菜单驱动：菜单 upsert 后，前端登录拉取菜单树 → 动态路由自动注册 → 侧边栏出现
- 组件路径 `plugin/resource/index` 对应 `web/src/views/plugin/resource/index.vue`
- 按钮权限码自动收集进 permission store，页面用 `v-permission="'resource:add'"` 控制

## 快速开始（新插件步骤）

1. 建目录：`mkdir server/plugin/<name>`，按上述结构建 model / dto / repo / service / api.go / plugin.go
2. 写 Model（`TableName()` 返回 `plugin_<name>_xxx`）
3. Repo / Service 面向接口，构造函数注入（`NewXxxRepo(db)` / `NewXxxService(repo)`），包级单例 `DefaultXxx` 在 `InitServices` 装配
4. `api.go` 写 handler，附 swagger 注解（@Router /api/plugin/<name>/v1/...）
5. `plugin.go` 实现 Plugin 接口：`Menus()` 定义菜单与按钮权限、`RegisterRoute` 挂路由、`InitServices` 装配、`SeedTable` 初始数据
6. `initialize/plugin.go` 加空导入
7. 前端：`web/src/views/plugin/<name>/index.vue` 写页面（数据请求走适配层或 api 封装）；如需前端插件入口建 `web/src/plugin/<name>/index.js`
8. 重启后端：AutoMigrate → 菜单/API 幂等 upsert → 路由注册，侧边栏自动出现

> 已安装系统新增插件：重启后端即自动迁移；超级管理员自动获得权限，其他角色在「角色授权」中勾选。
