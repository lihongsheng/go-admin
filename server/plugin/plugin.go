// Package plugin 插件接口与注册中心
//
// 插件以 Go 包形式存在于 server/plugin/<name>/，
// 在其 init() 内调用 plugin.Register(p) 完成自注册。
//
// 启动期 initialize.LoadPlugins() / SyncOnBoot() 行为：
//   - 把每个插件的 Models() 注册到 installer，参与在线初始化与启动期增量迁移
//   - 把每个插件的 Apis() / Menus() 在启动时按 (path,method) / (name) 做幂等 upsert
//   - SeedTable(db) 仅在「插件目标表当前为空」时执行，避免重复写入
//   - RegisterRoute(g) 挂到 /api/v1/plugin/<name>
//
// 菜单设计：
//   插件 Menus() 返回的是完整菜单树（含 catalog / menu / button 节点），
//   button 节点作为 menu 的 Children 存在，通过 type="button" 区分。
//   upsert 时递归处理整棵树，按 name 幂等。
package plugin

import (
	"sync"

	"go-admin/server/core/installer"
	"go-admin/server/model/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Plugin 插件契约
type Plugin interface {
	Name() string                                  // 唯一名（路由前缀 / 表前缀建议同名）
	Version() string                               // 版本号
	Models() []interface{}                         // 参与 AutoMigrate 的 Model
	Menus() []system.SysMenu                       // 注入菜单树（含 catalog/menu/button；按 Name 幂等）
	Apis() []system.SysApi                         // 注入 API（按 path+method 幂等）
	RegisterRoute(g *gin.RouterGroup)              // 注册自身路由（已在 /api/v1/plugin/<name> 下）
	SeedTable(db *gorm.DB) error                   // 插件自身业务表的初始数据；仅在目标表为空时调用
}

var (
	mu       sync.RWMutex
	registry []Plugin
)

// Register 由插件 init() 调用
func Register(p Plugin) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, p)

	// Model 进入 installer 注册中心
	installer.Register(p.Models()...)

	// 把插件菜单/API/SeedTable 推到 installer.Seed，用于「首次在线安装」
	// 启动期增量同步则由 initialize.SyncOnBoot 单独处理
	pl := p
	installer.RegisterSeed(func(d installer.SeedDeps) error {
		if err := upsertMenusAndApis(d.DB, pl, true); err != nil {
			return err
		}
		return pl.SeedTable(d.DB)
	})
}

// All 返回全部已注册插件
func All() []Plugin {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]Plugin, len(registry))
	copy(out, registry)
	return out
}

// ---------- 菜单 / API 幂等同步 ----------

// upsertMenusAndApis 递归插入插件菜单树 + 幂等写入 API
// attachSuper 决定是否自动挂到超级管理员角色
func upsertMenusAndApis(db *gorm.DB, p Plugin, attachSuper bool) error {
	// ----- Menus（递归 upsert 整棵树，按 name 幂等）-----
	for _, m := range p.Menus() {
		if err := upsertMenuTree(db, &m, 0, attachSuper); err != nil {
			return err
		}
	}

	// ----- APIs（按 path+method 幂等）-----
	for _, a := range p.Apis() {
		var exist system.SysApi
		err := db.Where("path = ? AND method = ?", a.Path, a.Method).First(&exist).Error
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&a).Error; err != nil {
				return err
			}
			exist = a
		} else if err != nil {
			return err
		} else {
			// 已存在则更新 desc / group
			db.Model(&exist).Updates(map[string]interface{}{
				"desc": a.Desc, "group": a.Group,
			})
		}
		if attachSuper {
			if err := attachApiToSuper(db, exist.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// upsertMenuTree 递归创建 / 更新单棵菜单子树
// parentID 为父节点 ID；创建成功后按 name 幂等，已存在则更新字段
func upsertMenuTree(db *gorm.DB, m *system.SysMenu, parentID uint, attachSuper bool) error {
	children := m.Children
	m.Children = nil
	m.ParentID = parentID

	var exist system.SysMenu
	err := db.Where("name = ?", m.Name).First(&exist).Error
	if err == gorm.ErrRecordNotFound {
		if err := db.Create(m).Error; err != nil {
			return err
		}
		exist = *m
	} else if err != nil {
		return err
	} else {
		// 已存在则更新关键字段
		db.Model(&exist).Updates(map[string]interface{}{
			"type": m.Type, "parent_id": parentID,
			"path": m.Path, "component": m.Component,
			"title": m.Title, "icon": m.Icon, "sort": m.Sort,
			"permission": m.Permission, "hidden": m.Hidden,
			"keep_alive": m.KeepAlive, "redirect": m.Redirect,
		})
	}

	// 超级管理员自动获得该菜单
	if attachSuper {
		if err := attachMenuToSuper(db, exist.ID); err != nil {
			return err
		}
	}

	// 递归处理子节点（catalog 的 menu 子节点，menu 的 button 子节点）
	for i := range children {
		if err := upsertMenuTree(db, &children[i], exist.ID, attachSuper); err != nil {
			return err
		}
	}
	return nil
}

// ---------- 超级管理员关联 ----------

func attachMenuToSuper(db *gorm.DB, menuID uint) error {
	var super system.SysRole
	if err := db.Where("code = ?", "super_admin").First(&super).Error; err != nil {
		return nil // 超级管理员还没建（极早期），跳过
	}
	return db.Model(&super).Association("Menus").Append(&system.SysMenu{Base: system.Base{ID: menuID}})
}

func attachApiToSuper(db *gorm.DB, apiID uint) error {
	var super system.SysRole
	if err := db.Where("code = ?", "super_admin").First(&super).Error; err != nil {
		return nil
	}
	return db.Model(&super).Association("Apis").Append(&system.SysApi{Base: system.Base{ID: apiID}})
}

// ---------- 启动期同步 ----------

// SyncOnBoot 启动期增量同步（所有插件）：
//   1) 对每个插件 Model AutoMigrate（installer.AllModels 内已包含）
//   2) 幂等 upsert 菜单 / API
//   3) 若插件自身业务表为空，执行 SeedTable
// 该函数对外暴露，由 initialize.SyncOnBoot 调用。
func SyncOnBoot(db *gorm.DB) error {
	// 1) 全量 AutoMigrate（系统 + 全部插件 Model）
	for _, m := range installer.AllModels() {
		if err := db.AutoMigrate(m); err != nil {
			return err
		}
	}
	// 2-3) 插件级 upsert + 条件 Seed
	for _, p := range All() {
		if err := upsertMenusAndApis(db, p, true); err != nil {
			return err
		}
		if err := seedIfEmpty(db, p); err != nil {
			return err
		}
	}
	return nil
}

// seedIfEmpty 仅当插件首张 Model 表为空时调用 SeedTable
func seedIfEmpty(db *gorm.DB, p Plugin) error {
	ms := p.Models()
	if len(ms) == 0 {
		return p.SeedTable(db) // 无 Model 的插件直接交给插件自行幂等
	}
	var cnt int64
	if err := db.Model(ms[0]).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt > 0 {
		return nil
	}
	return p.SeedTable(db)
}
