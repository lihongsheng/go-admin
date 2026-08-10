# 资源管理模块

云 ECS 资源管理插件（`server/plugin/resource/` + `web/src/views/plugin/resource/`），以插件化方式完整实现：后端 Model / Repo / Service / API 分层、菜单与按钮权限自动注册、前端页面（搜索 / 表格 / 右侧抽屉 / 批量删除）。

## 功能清单

| 功能 | 说明 |
|---|---|
| 分页列表 | 关键字（资源名称）搜索 + 分页 |
| 详情查看 | 右侧抽屉展示资源详情（el-descriptions） |
| 新增资源 | 右侧抽屉表单（名称 / 类型 / 状态），工具栏右侧按钮 |
| 编辑资源 | 右侧抽屉表单；**状态变更时自动刷新状态更新时间** |
| 批量删除 | 表格多选 + 批量删除按钮，权限码 `resource:del` |

## 数据表

表名 `plugin_resource_ecs`（插件表统一 `plugin_<name>_` 前缀）：

| 字段 | 类型 | 说明 |
|---|---|---|
| id | uint PK | 主键 |
| name | varchar(128) | 资源名称 |
| type | varchar(64) | 资源类型（如 通用型 g6） |
| status | int | 1 运行中 / 2 已停止 / 3 已释放 |
| status_at | datetime | 状态更新时间（仅状态变更时刷新） |
| created_at / updated_at | datetime | 创建 / 更新时间 |

## 后端接口

路由挂在 `/api/plugin/resource/v1`（`/api/plugin` 组带 JWT + Casbin 鉴权）：

| 方法 | 路径 | 说明 | 权限 |
|---|---|---|---|
| GET | `/list` | 分页列表（page / limit / keyword） | 菜单级（resource 菜单） |
| GET | `/:id` | 资源详情 | 菜单级 |
| POST | `/` | 新增资源 | `resource:add` |
| PUT | `/:id` | 编辑资源 | `resource:edit` |
| DELETE | `/batch` | 批量删除（body `{ids: []}`） | `resource:del` |

统一响应结构：

```json
{ "code": 0, "msg": "ok", "data": { ... } }
```

### curl 示例

```bash
# 分页列表
curl 'http://localhost:8989/api/plugin/resource/v1/list?page=1&limit=10&keyword=web' \
  -H "Authorization: Bearer <token>"

# 详情
curl 'http://localhost:8989/api/plugin/resource/v1/1' -H "Authorization: Bearer <token>"

# 新增
curl -X POST 'http://localhost:8989/api/plugin/resource/v1' \
  -H "Authorization: Bearer <token>" -H 'Content-Type: application/json' \
  -d '{"name":"web-server-02","type":"通用型 g6","status":1}'

# 编辑（status 变更会刷新 status_at）
curl -X PUT 'http://localhost:8989/api/plugin/resource/v1/1' \
  -H "Authorization: Bearer <token>" -H 'Content-Type: application/json' \
  -d '{"name":"web-server-01","type":"通用型 g6","status":2}'

# 批量删除
curl -X DELETE 'http://localhost:8989/api/plugin/resource/v1/batch' \
  -H "Authorization: Bearer <token>" -H 'Content-Type: application/json' \
  -d '{"ids":[1,2]}'
```

## 菜单与权限

`Menus()` 返回菜单树，启动期自动 upsert 并写入 Casbin：

```
资源管理  /plugin/resource  （icon: monitor, sort: 92）
├── 新增资源     resource:add   → POST /api/plugin/resource/v1
├── 编辑资源     resource:edit  → PUT  /api/plugin/resource/v1/:id
└── 批量删除资源 resource:del   → DELETE /api/plugin/resource/v1/batch
```

> 插件按钮级 API 权限由 `plugin/plugin.go` 的 `collectApiRules` 递归提取并写入 Casbin（超级管理员启动时自动挂载；其他角色需在角色授权中勾选）。

## 前端页面

`web/src/views/plugin/resource/index.vue`：

- **搜索栏**：资源名称关键字 + 查询 / 重置
- **表格**：ID / 资源名称 / 资源类型（tag）/ 状态（tag）/ 状态更新时间 / 操作（查看 / 编辑）+ 多选列
- **新增 / 编辑**：右侧抽屉（480px），资源类型为可输入选择的 ECS 规格族
- **查看**：右侧抽屉，el-descriptions 详情
- **批量删除**：工具栏按钮（未勾选禁用），确认后调用 `batchDeleteResources`

页面所有数据请求走 `api/resource.js` 封装（不直接拼接 URL）：

```js
import { resourceList, resourceCreate, resourceUpdate, resourceBatchDelete } from '@/api/resource'

const { data } = await resourceList({ page, limit, keyword })
await resourceCreate({ name, type, status })
await resourceUpdate(id, { name, type, status })
await resourceBatchDelete(ids)
```

Mock 由统一方案覆盖（见 [MOCK_MODE.md](MOCK_MODE.md)）：mock 开启时上述请求被拦截器分发到 `api/mock/resource.js`，关闭时直达真实后端。
