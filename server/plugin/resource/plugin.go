// Package resource 资源管理插件：云资源（ECS）实例管理
//
// 菜单树说明：
//
//	资源管理 menu 节点下挂 2 个 button 权限节点（新增/编辑），API 规则通过菜单 ApiRules
//	注入，自动写入 Casbin 策略。
package resource

import (
	"go-admin/server/model/system"
	"go-admin/server/plugin"
	resourceModel "go-admin/server/plugin/resource/model"
	repoResource "go-admin/server/plugin/resource/repo"
	serviceResource "go-admin/server/plugin/resource/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type p struct{}

func (p) Name() string    { return "resource" }
func (p) Version() string { return "0.1.0" }

func (p) Models() []interface{} { return []interface{}{&resourceModel.Resource{}} }

func (p) Menus() []system.SysMenu {
	return []system.SysMenu{
		{
			Type:      system.MenuTypeMenu,
			Path:      "/plugin/resource",
			Name:      "PluginResource",
			Component: "plugin/resource/index",
			Title:     "资源管理",
			Icon:      "monitor",
			Sort:      92,
			ApiRules:  `[{"path":"/api/plugin/resource/v1/list","method":"GET"},{"path":"/api/plugin/resource/v1/:id","method":"GET"}]`,
			Children: []system.SysMenu{
				{Type: system.MenuTypeButton, Name: "新增资源", Permission: "resource:add", ApiRules: `[{"path":"/api/plugin/resource/v1","method":"POST"}]`},
				{Type: system.MenuTypeButton, Name: "编辑资源", Permission: "resource:edit", ApiRules: `[{"path":"/api/plugin/resource/v1/:id","method":"PUT"}]`},
				{Type: system.MenuTypeButton, Name: "批量删除资源", Permission: "resource:del", ApiRules: `[{"path":"/api/plugin/resource/v1/batch","method":"DELETE"}]`},
			},
		},
	}
}

// InitServices 装配资源服务：Repo(DB)，不依赖 global
func (p) InitServices(ctx plugin.InitContext) error {
	repo := repoResource.NewResourceRepo(ctx.DB)
	serviceResource.DefaultResource = serviceResource.NewResourceService(repo)
	return nil
}

func (p) RegisterRoute(g *gin.Engine, privatePlugin *gin.RouterGroup) {
	privatePlugin.POST("resource/v1", create)
	privatePlugin.PUT("resource/v1/:id", update)
	privatePlugin.DELETE("resource/v1/batch", batchDel)
	privatePlugin.GET("resource/v1/:id", detail)
	privatePlugin.GET("resource/v1/list", list)
}

// SeedTable 资源表无需初始数据
func (p) SeedTable(db *gorm.DB) error { return nil }

func init() { plugin.Register(p{}) }
