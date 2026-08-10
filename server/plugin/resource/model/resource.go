// Package model resource 插件 Model（独立子包以避免与服务层形成循环依赖）
package model

import "time"

// 资源运行状态
const (
	StatusRunning  = 1 // 运行中
	StatusStopped  = 2 // 已停止
	StatusReleased = 3 // 已释放
)

// Resource 云资源（ECS 等）实例
type Resource struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"` // 资源名称
	Type      string    `gorm:"size:64;not null" json:"type"`  // 资源类型（如 通用型 g6）
	Status    int       `gorm:"default:1" json:"status"`       // 运行状态：1 运行中 / 2 已停止 / 3 已释放
	StatusAt  time.Time `json:"status_at"`                     // 状态更新时间
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 自定义表名，避免与系统表冲突
func (Resource) TableName() string { return "plugin_resource_ecs" }
