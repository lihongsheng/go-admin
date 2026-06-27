package installer

import (
	"github.com/lihongsheng/go-admin/server/enum"
	"github.com/lihongsheng/go-admin/server/model/system"
)

// defaultMenus 默认菜单树（catalog -> menu -> button 三层）
// catalog 仅作为目录承载 children；menu 渲染页面；button 表示按钮权限
func defaultMenus() []system.SysMenu {
	return []system.SysMenu{
		// 单个独立菜单：Dashboard 不属于目录
		{
			Type: system.MenuTypeMenu,
			Path: "/dashboard", Name: "Dashboard", Component: "dashboard/index",
			Title: "仪表盘", Icon: "dashboard", Sort: 1,
		},
		// 系统管理（catalog）
		{
			Type: system.MenuTypeCatalog,
			Path: "/system", Name: "System", Component: "Layout",
			Title: "系统管理", Icon: "setting", Sort: 10,
			Children: []system.SysMenu{
				{
					Type: system.MenuTypeMenu,
					Path: "user", Name: "SysUser", Component: "system/user/index",
					Title: "用户管理", Icon: "user", Sort: 1,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增用户", Permission: "user:add"},
						{Type: system.MenuTypeButton, Name: "编辑用户", Permission: "user:edit"},
						{Type: system.MenuTypeButton, Name: "删除用户", Permission: "user:del"},
						{Type: system.MenuTypeButton, Name: "重置密码", Permission: "user:reset"},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "role", Name: "SysRole", Component: "system/role/index",
					Title: "角色管理", Icon: "peoples", Sort: 2,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增角色", Permission: "role:add"},
						{Type: system.MenuTypeButton, Name: "编辑角色", Permission: "role:edit"},
						{Type: system.MenuTypeButton, Name: "删除角色", Permission: "role:del"},
						{Type: system.MenuTypeButton, Name: "角色授权", Permission: "role:auth"},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "menu", Name: "SysMenu", Component: "system/menu/index",
					Title: "菜单管理", Icon: "tree-table", Sort: 3,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增菜单", Permission: "menu:add"},
						{Type: system.MenuTypeButton, Name: "编辑菜单", Permission: "menu:edit"},
						{Type: system.MenuTypeButton, Name: "删除菜单", Permission: "menu:del"},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "api", Name: "SysApi", Component: "system/api/index",
					Title: "API管理", Icon: "api", Sort: 4,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增API", Permission: "api:add"},
						{Type: system.MenuTypeButton, Name: "编辑API", Permission: "api:edit"},
						{Type: system.MenuTypeButton, Name: "删除API", Permission: "api:del"},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "mch", Name: "SysMch", Component: "plugin/mch/view/index",
					Title: "商户管理", Icon: "shop", Sort: 5,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增商户", Permission: "mch:add"},
						{Type: system.MenuTypeButton, Name: "编辑商户", Permission: "mch:edit"},
						{Type: system.MenuTypeButton, Name: "查看商户", Permission: "mch:view"},
						{Type: system.MenuTypeButton, Name: "商户状态", Permission: "mch:status"},
					},
				},
			},
		},

		// 插件中心
		{
			Type: system.MenuTypeCatalog,
			Path: "/plugin", Name: "Plugin", Component: "Layout",
			Title: "插件中心", Icon: "component", Sort: 90,
			Children: []system.SysMenu{
				{
					Type: system.MenuTypeMenu,
					Path: "list", Name: "PluginList", Component: "plugin/list/index",
					Title: "已装插件", Icon: "list", Sort: 1,
				},
			},
		},
	}
}

// defaultApis 默认 API 列表（与路由 / 控制器对应）
func defaultApis() []system.SysApi {
	return []system.SysApi{
		{Path: "/api/v1/system/user", Method: "POST", Group: "user", Desc: "新增用户"},
		{Path: "/api/v1/system/user", Method: "PUT", Group: "user", Desc: "编辑用户"},
		{Path: "/api/v1/system/user/:id", Method: "DELETE", Group: "user", Desc: "删除用户"},
		{Path: "/api/v1/system/user/list", Method: "GET", Group: "user", Desc: "用户列表"},

		{Path: "/api/v1/system/role", Method: "POST", Group: "role", Desc: "新增角色"},
		{Path: "/api/v1/system/role", Method: "PUT", Group: "role", Desc: "编辑角色"},
		{Path: "/api/v1/system/role/:id", Method: "DELETE", Group: "role", Desc: "删除角色"},
		{Path: "/api/v1/system/role/list", Method: "GET", Group: "role", Desc: "角色列表"},
		{Path: "/api/v1/system/role/auth", Method: "POST", Group: "role", Desc: "角色授权"},
		{Path: "/api/v1/system/role/auth/:id", Method: "GET", Group: "role", Desc: "查询角色授权详情"},

		{Path: "/api/v1/system/menu", Method: "POST", Group: "menu", Desc: "新增菜单"},
		{Path: "/api/v1/system/menu", Method: "PUT", Group: "menu", Desc: "编辑菜单"},
		{Path: "/api/v1/system/menu/:id", Method: "DELETE", Group: "menu", Desc: "删除菜单"},
		{Path: "/api/v1/system/menu/tree", Method: "GET", Group: "menu", Desc: "菜单树"},

		{Path: "/api/v1/system/api", Method: "POST", Group: "api", Desc: "新增API"},
		{Path: "/api/v1/system/api", Method: "PUT", Group: "api", Desc: "编辑API"},
		{Path: "/api/v1/system/api/:id", Method: "DELETE", Group: "api", Desc: "删除API"},
		{Path: "/api/v1/system/api/list", Method: "GET", Group: "api", Desc: "API列表"},

		{Path: "/api/v1/plugin/list", Method: "GET", Group: "plugin", Desc: "插件列表"},

		{Path: "/api/v1/system/mch", Method: "POST", Group: "mch", Desc: "新增商户"},
		{Path: "/api/v1/system/mch", Method: "PUT", Group: "mch", Desc: "编辑商户"},
		{Path: "/api/v1/system/mch/:id", Method: "GET", Group: "mch", Desc: "商户详情"},
		{Path: "/api/v1/system/mch/no/:mchNo", Method: "GET", Group: "mch", Desc: "按编号查商户"},
		{Path: "/api/v1/system/mch/list", Method: "GET", Group: "mch", Desc: "商户列表"},
		{Path: "/api/v1/system/mch/status", Method: "PUT", Group: "mch", Desc: "修改商户状态"},
	}
}

// defaultMerchant 默认商户
func defaultMerchant() system.Merchant {
	return system.Merchant{
		MchName: "默认商户",
		Linker:  "管理员",
		Phone:   "13800000000",
		Email:   "admin@merchant.local",
		Address: "默认地址",
		Status:  1,
	}
}

// merchantAdminMenus 商户管理员菜单（仅用户管理 + 角色管理 + Dashboard）
func merchantAdminMenus() []system.SysMenu {
	return []system.SysMenu{
		{
			Type: system.MenuTypeMenu,
			Path: "/dashboard", Name: "Dashboard", Component: "dashboard/index",
			Title: "仪表盘", Icon: "dashboard", Sort: 1,
			SystemType: enum.SystemTypeMch,
		},
		{
			Type: system.MenuTypeCatalog,
			Path: "/system", Name: "System", Component: "Layout",
			Title: "系统管理", Icon: "setting", Sort: 10,
			SystemType: enum.SystemTypeMch,
			Children: []system.SysMenu{
				{
					Type: system.MenuTypeMenu,
					Path: "user", Name: "SysUser", Component: "system/user/index",
					Title: "用户管理", Icon: "user", Sort: 1,
					SystemType: enum.SystemTypeMch,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增用户", Permission: "user:add", SystemType: enum.SystemTypeMch},
						{Type: system.MenuTypeButton, Name: "编辑用户", Permission: "user:edit", SystemType: enum.SystemTypeMch},
						{Type: system.MenuTypeButton, Name: "删除用户", Permission: "user:del", SystemType: enum.SystemTypeMch},
						{Type: system.MenuTypeButton, Name: "重置密码", Permission: "user:reset", SystemType: enum.SystemTypeMch},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "role", Name: "SysRole", Component: "system/role/index",
					Title: "角色管理", Icon: "peoples", Sort: 2,
					SystemType: enum.SystemTypeMch,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增角色", Permission: "role:add", SystemType: enum.SystemTypeMch},
						{Type: system.MenuTypeButton, Name: "编辑角色", Permission: "role:edit", SystemType: enum.SystemTypeMch},
						{Type: system.MenuTypeButton, Name: "删除角色", Permission: "role:del", SystemType: enum.SystemTypeMch},
						{Type: system.MenuTypeButton, Name: "角色授权", Permission: "role:auth", SystemType: enum.SystemTypeMch},
					},
				},
			},
		},
	}
}

// merchantAdminApis 商户管理员 API（基础 + 用户 + 角色 + 菜单树 + API列表）
func merchantAdminApis() []system.SysApi {
	return []system.SysApi{

		{Path: "/api/v1/system/user", Method: "POST", Group: "user", Desc: "新增用户", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/user", Method: "PUT", Group: "user", Desc: "编辑用户", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/user/:id", Method: "DELETE", Group: "user", Desc: "删除用户", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/user/list", Method: "GET", Group: "user", Desc: "用户列表", SystemType: enum.SystemTypeMch},

		{Path: "/api/v1/system/role", Method: "POST", Group: "role", Desc: "新增角色", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/role", Method: "PUT", Group: "role", Desc: "编辑角色", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/role/:id", Method: "DELETE", Group: "role", Desc: "删除角色", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/role/list", Method: "GET", Group: "role", Desc: "角色列表", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/role/auth", Method: "POST", Group: "role", Desc: "角色授权", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/role/auth/:id", Method: "GET", Group: "role", Desc: "查询角色授权详情", SystemType: enum.SystemTypeMch},

		{Path: "/api/v1/system/menu/tree", Method: "GET", Group: "menu", Desc: "菜单树", SystemType: enum.SystemTypeMch},
		{Path: "/api/v1/system/api/list", Method: "GET", Group: "api", Desc: "API列表", SystemType: enum.SystemTypeMch},
	}
}

func init() {
	// 自动注册系统核心 Model
	Register(system.CoreModels()...)
}
