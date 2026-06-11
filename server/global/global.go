// Package global 持有进程级共享对象：配置 / DB / Logger / 安装状态
// 注意：尽量避免使用全局变量，新代码应该使用依赖注入
package global

import (
	"sync/atomic"

	"go-admin/server/config"

	"gorm.io/gorm"
)

var (
	Cfg    *config.Config
	DB     *gorm.DB

	// Installed 全局安装状态，安装中间件读取
	Installed atomic.Bool
)

// SetDB 安装完成 / 启动加载完成后写入
func SetDB(db *gorm.DB) { DB = db }
