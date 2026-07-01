// Package example 示例插件：演示自有 Model + 菜单 + API + 路由
//
// 菜单树说明：
//
//	插件返回的 Menus() 是完整的菜单子树，button 权限节点作为 menu 的 Children 存在。
//	每个节点通过 type 字段标识：catalog / menu / button
//
// Note Model 已移到 plugin/example/model 子包，避免与 service/plugin/example 形成循环依赖。
package example

import (
	"github.com/lihongsheng/go-admin/server/model/system"
	"github.com/lihongsheng/go-admin/server/plugin"
	exampleModel "github.com/lihongsheng/go-admin/server/plugin/example/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type p struct{}

func (p) Name() string    { return "example" }
func (p) Version() string { return "0.1.0" }

func (p) Models() []interface{} { return []interface{}{&exampleModel.Note{}} }

func (p) Menus() []system.SysMenu {
	return []system.SysMenu{
		{
			Type:      system.MenuTypeMenu,
			Path:      "/plugin/example",
			Name:      "PluginExample",
			Component: "plugin/example/index",
			Title:     "示例插件",
			Icon:      "edit",
			Sort:      91,
			ApiRules:  `[{"path":"/api/plugin/example/v1/note/list","method":"GET"}]`,
			Children: []system.SysMenu{
				{Type: system.MenuTypeButton, Name: "新增笔记", Permission: "example:add", ApiRules: `[{"path":"/api/plugin/example/v1/note","method":"POST"}]`},
				{Type: system.MenuTypeButton, Name: "删除笔记", Permission: "example:del", ApiRules: `[{"path":"/api/plugin/example/v1/note/:id","method":"DELETE"}]`},
			},
		},
	}
}

func (p) RegisterRoute(g *gin.Engine, privatePlugin *gin.RouterGroup) {
	g.POST("/note", create)
	g.DELETE("/note/:id", del)
	g.GET("/note/list", list)
}

func (p) SeedTable(db *gorm.DB) error {
	return db.Create(&exampleModel.Note{Title: "Hello", Content: "示例插件初始化笔记"}).Error
}

func init() { plugin.Register(p{}) }
