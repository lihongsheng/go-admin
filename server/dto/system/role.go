package system

import "go-admin/server/model/system"

// RoleCreateReq 新增角色
type RoleCreateReq struct {
	Name          string `json:"name"           binding:"required"`
	Code          string `json:"code"           binding:"required"`
	Remark        string `json:"remark"`
	Status        int8   `json:"status"`
	DefaultRouter string `json:"default_router"`
}

// RoleUpdateReq 更新角色；Code 变更会触发 Casbin 策略迁移
type RoleUpdateReq struct {
	ID            uint   `json:"id"             binding:"required"`
	Name          string `json:"name"           binding:"required"`
	Code          string `json:"code"           binding:"required"`
	Remark        string `json:"remark"`
	Status        int8   `json:"status"`
	DefaultRouter string `json:"default_router"`
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

// RoleAuthDetailResp 角色已分配的菜单/API id 列表 + 默认首页路由
type RoleAuthDetailResp struct {
	MenuIDs       []uint `json:"menu_ids"`
	ApiIDs        []uint `json:"api_ids"`
	DefaultRouter string `json:"default_router"`
}

// RoleSetDefaultRouterReq 设置角色默认首页路由
type RoleSetDefaultRouterReq struct {
	DefaultRouter string `json:"default_router"`
}
