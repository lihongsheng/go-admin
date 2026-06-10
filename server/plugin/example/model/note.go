// Package model example 插件 Model（独立子包以避免与 service/plugin/example 形成循环依赖）
package model

// Note 插件自有表（带 plugin_example_ 前缀）
type Note struct {
	ID      uint   `gorm:"primaryKey" json:"id"`
	Title   string `gorm:"size:128"   json:"title"`
	Content string `gorm:"type:text"  json:"content"`
}

// TableName 自定义表名，避免与系统表冲突
func (Note) TableName() string { return "plugin_example_notes" }
