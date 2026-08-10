package system

import (
	"go-admin/server/model/system"
)

// MenuCreateReq 新增菜单
type MenuCreateReq struct {
	Type       string `json:"type"` // catalog / menu / button，空字符串表示 menu
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Component  string `json:"component"`
	Redirect   string `json:"redirect"`
	Permission string `json:"permission"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Hidden     bool   `json:"hidden"`
	KeepAlive  bool   `json:"keep_alive"`
	ApiRules   string `json:"api_rules"`
}

// MenuUpdateReq 更新菜单
type MenuUpdateReq struct {
	ID         uint   `json:"id" binding:"required"`
	Type       string `json:"type"`
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Component  string `json:"component"`
	Redirect   string `json:"redirect"`
	Permission string `json:"permission"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Sort       int    `json:"sort"`
	Hidden     bool   `json:"hidden"`
	KeepAlive  bool   `json:"keep_alive"`
	ApiRules   string `json:"api_rules"`
}

// MenuTreeResp 菜单树响应
type MenuTreeResp struct {
	List []system.SysMenu `json:"list"`
}

// MenuTreeReq 菜单树查询
type MenuTreeReq struct {
}
