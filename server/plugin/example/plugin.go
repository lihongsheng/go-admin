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
			// button 权限节点作为 Children（type=button，仅承载 permission 字段）
			Children: []system.SysMenu{
				{Type: system.MenuTypeButton, Name: "新增笔记", Permission: "example:add"},
				{Type: system.MenuTypeButton, Name: "删除笔记", Permission: "example:del"},
			},
		},
	}
}

func (p) Apis() []system.SysApi {
	return []system.SysApi{
		{Path: "/api/v1/plugin/example/note", Method: "POST", Group: "example", Desc: "新增笔记"},
		{Path: "/api/v1/plugin/example/note/:id", Method: "DELETE", Group: "example", Desc: "删除笔记"},
		{Path: "/api/v1/plugin/example/note/list", Method: "GET", Group: "example", Desc: "笔记列表"},
	}
}

// InitServices 示例插件无需初始化服务层，返回 nil
func (p) InitServices(ctx plugin.InitContext) error {
	return nil
}

func (p) RegisterRoute(g *gin.Engine) {
	g.POST("/note", create)
	g.DELETE("/note/:id", del)
	g.GET("/note/list", list)
}

func (p) SeedTable(db *gorm.DB) error {
	return db.Create(&exampleModel.Note{Title: "Hello", Content: "示例插件初始化笔记"}).Error
}

func init() { plugin.Register(p{}) }
