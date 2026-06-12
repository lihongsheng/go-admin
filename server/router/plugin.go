package router

import (
	"go-admin/server/log"
	"go-admin/server/middleware"
	"go-admin/server/plugin"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// PluginRouter ：/api/v1/plugin/list 与每个插件自身路由 /api/v1/plugin/<name>/*
func PluginRouter(g *gin.Engine) {
	p := g.Group("/private/v1/plugin", middleware.JWTAuth())
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
		// sub := p.Group("/" + pl.Name())
		pl.RegisterRoute(g)
		log.Global().Info("plugin route mounted: /api/v1/plugin/" + pl.Name())
	}
}
