// Package enum 系统枚举常量
package enum

// UserStatus 用户状态
type UserStatus int8

const (
	UserStatusDisabled UserStatus = 0 // 禁用
	UserStatusEnabled  UserStatus = 1 // 启用
)
