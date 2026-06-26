// Package install install 模块 DTO
package install

import (
	"github.com/lihongsheng/go-admin/server/config"
	"github.com/lihongsheng/go-admin/server/core/installer"
)

// CheckDBReq 数据库连接检测
type CheckDBReq struct {
	config.DB
	CreateIfMissing bool `json:"create_if_missing"`
}

// CheckDBResp 数据库检测响应
type CheckDBResp struct {
	Configured  bool   `json:"configured"`
	DBConnected bool   `json:"db_connected"`
	Installed   bool   `json:"installed"`
	Version     string `json:"version,omitempty"`
	Driver      string `json:"driver"`
	Reason      string `json:"reason,omitempty"`
	Created     bool   `json:"created"`
}

// InitReq 安装请求（同步 / SSE 公用）
type InitReq struct {
	DB    config.DB           `json:"db"    binding:"required"`
	Admin installer.AdminSeed `json:"admin" binding:"required"`
}

// InitResp 安装结果
type InitResp struct {
	OK        bool             `json:"ok"`
	CreatedDB bool             `json:"created_db"`
	Steps     []installer.Step `json:"steps"`
	Error     string           `json:"error,omitempty"`
}
