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
					ApiRules: `[{"path":"/api/v1/system/user/list","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增用户", Permission: "user:add", ApiRules: `[{"path":"/api/v1/system/user","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑用户", Permission: "user:edit", ApiRules: `[{"path":"/api/v1/system/user","method":"PUT"}]`},
						{Type: system.MenuTypeButton, Name: "删除用户", Permission: "user:del", ApiRules: `[{"path":"/api/v1/system/user/:id","method":"DELETE"}]`},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "role", Name: "SysRole", Component: "system/role/index",
					Title: "角色管理", Icon: "peoples", Sort: 2,
					ApiRules: `[{"path":"/api/v1/system/role/list","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增角色", Permission: "role:add", ApiRules: `[{"path":"/api/v1/system/role","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑角色", Permission: "role:edit", ApiRules: `[{"path":"/api/v1/system/role","method":"PUT"}]`},
						{Type: system.MenuTypeButton, Name: "删除角色", Permission: "role:del", ApiRules: `[{"path":"/api/v1/system/role/:id","method":"DELETE"}]`},
						{Type: system.MenuTypeButton, Name: "角色授权", Permission: "role:auth", ApiRules: `[{"path":"/api/v1/system/role/auth","method":"POST"},{"path":"/api/v1/system/role/auth/:id","method":"GET"}]`},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "menu", Name: "SysMenu", Component: "system/menu/index",
					Title: "菜单管理", Icon: "tree-table", Sort: 3,
					ApiRules: `[{"path":"/api/v1/system/menu/tree","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增菜单", Permission: "menu:add", ApiRules: `[{"path":"/api/v1/system/menu","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑菜单", Permission: "menu:edit", ApiRules: `[{"path":"/api/v1/system/menu","method":"PUT"}]`},
						{Type: system.MenuTypeButton, Name: "删除菜单", Permission: "menu:del", ApiRules: `[{"path":"/api/v1/system/menu/:id","method":"DELETE"}]`},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "mch", Name: "SysMch", Component: "plugin/mch/view/index",
					Title: "商户管理", Icon: "shop", Sort: 5,
					ApiRules: `[{"path":"/api/v1/system/mch/list","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增商户", Permission: "mch:add", ApiRules: `[{"path":"/api/v1/system/mch","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑商户", Permission: "mch:edit", ApiRules: `[{"path":"/api/v1/system/mch","method":"PUT"}]`},
						{Type: system.MenuTypeButton, Name: "查看商户", Permission: "mch:view", ApiRules: `[{"path":"/api/v1/system/mch/:id","method":"GET"},{"path":"/api/v1/system/mch/no/:mchNo","method":"GET"}]`},
						{Type: system.MenuTypeButton, Name: "商户状态", Permission: "mch:status", ApiRules: `[{"path":"/api/v1/system/mch/status","method":"PUT"}]`},
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
					ApiRules: `[{"path":"/api/plugin/list","method":"GET"}]`,
				},
			},
		},
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
					ApiRules:   `[{"path":"/api/v1/system/user/list","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增用户", Permission: "user:add", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/user","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑用户", Permission: "user:edit", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/user","method":"PUT"}]`},
						{Type: system.MenuTypeButton, Name: "删除用户", Permission: "user:del", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/user/:id","method":"DELETE"}]`},
					},
				},
				{
					Type: system.MenuTypeMenu,
					Path: "role", Name: "SysRole", Component: "system/role/index",
					Title: "角色管理", Icon: "peoples", Sort: 2,
					SystemType: enum.SystemTypeMch,
					ApiRules:   `[{"path":"/api/v1/system/role/list","method":"GET"}]`,
					Children: []system.SysMenu{
						{Type: system.MenuTypeButton, Name: "新增角色", Permission: "role:add", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/role","method":"POST"}]`},
						{Type: system.MenuTypeButton, Name: "编辑角色", Permission: "role:edit", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/role","method":"PUT"}]`},
						{Type: system.MenuTypeButton, Name: "删除角色", Permission: "role:del", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/role/:id","method":"DELETE"}]`},
						{Type: system.MenuTypeButton, Name: "角色授权", Permission: "role:auth", SystemType: enum.SystemTypeMch, ApiRules: `[{"path":"/api/v1/system/role/auth","method":"POST"},{"path":"/api/v1/system/menu/tree","method":"GET"},{"path":"/api/v1/system/role/auth/:id","method":"GET"}]`},
					},
				},
			},
		},
	}
}

func init() {
	// 自动注册系统核心 Model
	Register(system.CoreModels()...)
}
