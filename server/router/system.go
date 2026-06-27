package router

import (
	"github.com/lihongsheng/go-admin/server/api/v1/base"
	"github.com/lihongsheng/go-admin/server/api/v1/system"
	"github.com/lihongsheng/go-admin/server/middleware"

	"github.com/gin-gonic/gin"
)

// BaseRouter 登录 / 验证码 / 当前用户 / 当前菜单
func BaseRouter(g *gin.RouterGroup) {
	b := g.Group("/base")
	b.GET("/captcha", base.Captcha)
	b.POST("/login", base.Login)
	auth := b.Group("", middleware.JWTAuth())
	auth.POST("/logout", base.Logout)
	auth.GET("/info", base.Info)
	auth.GET("/menu", base.Menu)

	// 上传接口：登录用户均可上传，不走 CasbinAuth
	upload := b.Group("/upload", middleware.JWTAuth())
	upload.POST("", base.Upload)
	baseSystem := b.Group("/system", middleware.JWTAuth())
	baseSystem.GET("/types", base.GetSystemTypeInfo)
}

// SystemRouter 用户 / 角色 / 菜单 / API
func SystemRouter(g *gin.RouterGroup) {
	s := g.Group("/system", middleware.JWTAuth(), middleware.CasbinAuth())

	s.POST("/user", system.UserCreate)
	s.PUT("/user", system.UserUpdate)
	s.DELETE("/user/:id", system.UserDelete)
	s.GET("/user/list", system.UserList)

	s.POST("/role", system.RoleCreate)
	s.PUT("/role", system.RoleUpdate)
	s.DELETE("/role/:id", system.RoleDelete)
	s.GET("/role/list", system.RoleList)
	s.POST("/role/auth", system.RoleAuth)
	s.GET("/role/auth/:id", system.RoleAuthDetail)
	s.PUT("/role/:id/default-router", system.RoleSetDefaultRouter)

	s.POST("/menu", system.MenuCreate)
	s.PUT("/menu", system.MenuUpdate)
	s.DELETE("/menu/:id", system.MenuDelete)
	s.GET("/menu/tree", system.MenuTree)

	s.POST("/mch", system.MchCreate)
	s.PUT("/mch", system.MchUpdate)
	s.GET("/mch/:id", system.MchDetail)
	s.GET("/mch/no/:mchNo", system.MchDetailByNo)
	s.GET("/mch/list", system.MchList)
}
