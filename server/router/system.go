package router

import (
	"go-admin/server/api/v1/base"
	"go-admin/server/api/v1/system"
	"go-admin/server/middleware"

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
}

// SystemRouter 用户 / 角色 / 菜单 / API
func SystemRouter(g *gin.RouterGroup) {
	s := g.Group("/system", middleware.JWTAuth())

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

	s.POST("/api", system.ApiCreate)
	s.PUT("/api", system.ApiUpdate)
	s.DELETE("/api/:id", system.ApiDelete)
	s.GET("/api/list", system.ApiList)
}
