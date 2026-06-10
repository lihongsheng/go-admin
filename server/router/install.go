package router

import (
	"go-admin/server/api/v1/install"

	"github.com/gin-gonic/gin"
)

// InstallRouter /install/* 安装向导（免鉴权，由 InstallGuard 控制可用性）
func InstallRouter(r *gin.RouterGroup) {
	g := r.Group("/install")
	{
		g.GET("/status", install.Status)
		g.POST("/check-db", install.CheckDB)
		g.POST("/init", install.Init)
		g.POST("/stream", install.Stream)
	}
}
