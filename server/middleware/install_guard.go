package middleware

import (
	"net/http"
	"strings"

	"github.com/lihongsheng/go-admin/server/global"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// InstallGuard：未安装时只放行 /install/* 和 /health；已安装时关闭 /install/*
func InstallGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		p := c.Request.URL.Path
		isInstall := strings.HasPrefix(p, "/install") || strings.HasPrefix(p, "/api/v1/install")
		if !global.Installed.Load() {
			if isInstall || p == "/health" {
				c.Next()
				return
			}
			response.FailHTTP(c, http.StatusServiceUnavailable, response.CodeNotInstalled,
				"system not installed, please POST /install/init")
			c.Abort()
			return
		}
		if isInstall {
			response.FailHTTP(c, http.StatusForbidden, response.CodeAlreadyDone, "system already installed")
			c.Abort()
			return
		}
		c.Next()
	}
}
