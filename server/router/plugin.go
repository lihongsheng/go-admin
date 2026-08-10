package router

import (
	"go-admin/server/log"
	"go-admin/server/middleware"
	"go-admin/server/plugin"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// PluginRouter ：/api/plugin/list 与每个插件自身路由 /api/plugin/<name>/*
func PluginRouter(g *gin.Engine) {
	p := g.Group("/api/plugin", middleware.JWTAuth(), middleware.CasbinAuth())
	p.GET("/list", func(c *gin.Context) {
		type item struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		out := []item{}
		for _, pl := range plugin.All() {
			out = append(out, item{pl.Name(), pl.Version()})
		}
		response.OK(c, gin.H{"list": out, "total": len(out)})
	})
	for _, pl := range plugin.All() {
		pl.RegisterRoute(g, p)
		log.Global().Info("plugin route mounted: /api/plugin/" + pl.Name())
	}
}
