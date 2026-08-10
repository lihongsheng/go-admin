package system

import (
	"go-admin/server/model/system"
)

// RoleCreateReq 新增角色
type RoleCreateReq struct {
	Name          string `json:"name"           binding:"required"`
	Remark        string `json:"remark"`
	Status        int8   `json:"status"`
	DefaultRouter string `json:"default_router"`
}

// RoleUpdateReq 更新角色
type RoleUpdateReq struct {
	ID            uint   `json:"id"             binding:"required"`
	Name          string `json:"name"           binding:"required"`
	Remark        string `json:"remark"`
	Status        int8   `json:"status"`
	DefaultRouter string `json:"default_router"`
}

// RoleListReq 角色列表查询
type RoleListReq struct {
}

// RoleListResp 角色列表响应
type RoleListResp struct {
	List  []system.SysRole `json:"list"`
	Total int              `json:"total"`
}

// RoleAuthReq 角色授权（菜单 + API）
type RoleAuthReq struct {
	RoleID  uint   `json:"role_id"   binding:"required"`
	MenuIDs []uint `json:"menu_ids"`
	ApiIDs  []uint `json:"api_ids"`
}

// RoleAuthDetailResp 角色已分配的菜单 id 列表 + 默认首页路由
// (API 权限已合并到菜单 api_rules 中，不再单独返回)
type RoleAuthDetailResp struct {
	MenuIDs       []uint `json:"menu_ids"`
	DefaultRouter string `json:"default_router"`
}

// RoleSetDefaultRouterReq 设置角色默认首页路由
type RoleSetDefaultRouterReq struct {
	DefaultRouter string `json:"default_router"`
}
