package initialize

import (
	"net/http"

	"go-admin/server/global"
	"go-admin/server/middleware"
	"go-admin/server/router"

	"github.com/gin-gonic/gin"
)

// Router 构造 gin 路由
func Router() *gin.Engine {
	gin.SetMode(global.Cfg.App.Mode)
	r := gin.New()
	r.Use(gin.Recovery(), middleware.RequestLog(), middleware.Cors())

	// 健康检查 / 安装状态前置（不走 InstallGuard）
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// 安装向导直接挂在根路径（无 /api 前缀，方便前端独立模块）
	router.InstallRouter(&r.RouterGroup)

	// 业务路由统一在 /api/v1 下，由 InstallGuard 拦截
	api := r.Group("/api/v1", middleware.InstallGuard())
	router.BaseRouter(api)   // 登录 / 当前用户 / 当前菜单
	router.SystemRouter(api) // user / role / menu / api
	router.PluginRouter(api) // 已装插件列表 + 插件自身路由
	return r
}
