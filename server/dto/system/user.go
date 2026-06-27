// Package system system 模块 DTO
package system

import (
	"github.com/lihongsheng/go-admin/server/enum"
	"github.com/lihongsheng/go-admin/server/model/system"
)

// UserCreateReq 新增用户
type UserCreateReq struct {
	Username   string          `json:"username" binding:"required"`
	Password   string          `json:"password" binding:"required"`
	Nickname   string          `json:"nickname"`
	Email      string          `json:"email"`
	Phone      string          `json:"phone"`
	Status     int8            `json:"status"`
	RoleIDs    []uint          `json:"role_ids"`
	SystemType enum.SystemType `json:"system_type"`
	MchID      int64           `json:"mch_id"`
}

// UserUpdateReq 更新用户；Password 为空表示不修改密码
type UserUpdateReq struct {
	ID       uint   `json:"id" binding:"required"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
	// RoleIDs 为 nil 表示不调整角色；空切片表示清空角色
	RoleIDs    []uint          `json:"role_ids"`
	SystemType enum.SystemType `json:"system_type"`
	MchID      int64           `json:"mch_id"`
}

// UserListReq 用户列表查询
type UserListReq struct {
	Keyword    string `form:"keyword"`
	Page       int    `form:"page"`
	Size       int    `form:"size"`
	SystemType int    `form:"system_type"`
	MchID      int64  `form:"mch_id"`
}

// UserListResp 用户列表响应
type UserListResp struct {
	List  []system.SysUser `json:"list"`
	Total int64            `json:"total"`
}
