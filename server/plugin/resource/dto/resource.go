// Package dto resource 插件 DTO
package dto

// ResourceCreateReq 新建云资源
type ResourceCreateReq struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status int    `json:"status"`
}

// ResourceUpdateReq 更新云资源（名称 / 类型 / 状态）
type ResourceUpdateReq struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status int    `json:"status"`
}

// ResourceListReq 云资源分页查询
type ResourceListReq struct {
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
	Keyword string `form:"keyword"` // 资源名称模糊搜索
}

// ResourceBatchDeleteReq 批量删除
type ResourceBatchDeleteReq struct {
	IDs []uint `json:"ids"`
}
