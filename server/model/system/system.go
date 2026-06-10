// Package system 系统核心 Model
//
// 表结构设计说明：
//
//  1. RBAC 模型：用户(SysUser) ←多对多→ 角色(SysRole) ←多对多→ 菜单(SysMenu) / API(SysApi)
//     Casbin 负责 API 级别鉴权，策略存储在 casbin_rule 表（由 gorm-adapter 自动建表）。
//     casbin_rule 中 p 策略格式：p(role_code, /api/path, method)
//     casbin_rule 中 g 策略格式：g(u:user_id, role_code)
//
//  2. 菜单树设计：SysMenu 采用单表自引用树，通过 type 字段区分节点类型：
//     - catalog : 目录节点，仅作为容器分组，component=Layout，不渲染独立页面
//     - menu    : 菜单节点，对应一个可路由的前端页面
//     - button  : 按钮节点，作为菜单的子节点存在，仅承载权限码(permission)，
//                 不进入前端路由，由 v-permission 指令控制显隐
//     菜单树通过 parent_id 形成父子层级关系，Children 字段(gorm:"-")仅用于 JSON 序列化。
//
//  3. Casbin 表：casbin_rule 表由 gorm-adapter 在 casbin.Setup() 时自动 CREATE TABLE IF NOT EXISTS，
//     无需在 CoreModels 中注册，也不参与 AutoMigrate。
package system

import (
	"time"

	"gorm.io/gorm"
)

// Base 公共字段
type Base struct {
	ID        uint           `gorm:"primaryKey"        json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"             json:"-"`
}

// ---------- 用户 ----------

// SysUser 用户
type SysUser struct {
	Base
	Username string    `gorm:"size:64;uniqueIndex;not null" json:"username"`
	Password string    `gorm:"size:128;not null"            json:"-"`
	Nickname string    `gorm:"size:64"                       json:"nickname"`
	Avatar   string    `gorm:"size:255"                      json:"avatar"`
	Email    string    `gorm:"size:128"                      json:"email"`
	Phone    string    `gorm:"size:32"                       json:"phone"`
	Status   int8      `gorm:"default:1"                     json:"status"` // 1启用 0禁用
	Roles    []SysRole `gorm:"many2many:sys_user_roles"      json:"roles"`
}

// ---------- 角色 ----------

// SysRole 角色；Code 同时作为 Casbin 策略中的 sub（角色标识）
// 角色通过 sys_role_menus 关联菜单，通过 sys_role_apis 关联 API
type SysRole struct {
	Base
	Name   string    `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Code   string    `gorm:"size:64;uniqueIndex;not null" json:"code"` // Casbin 角色标识
	Remark string    `gorm:"size:255"                      json:"remark"`
	Status int8      `gorm:"default:1"                     json:"status"`
	Menus  []SysMenu `gorm:"many2many:sys_role_menus"      json:"menus"`
	Apis   []SysApi  `gorm:"many2many:sys_role_apis"       json:"apis"`
}

// ---------- 菜单树节点类型 ----------

const (
	MenuTypeCatalog = "catalog" // 目录：仅作为父节点分组，不渲染独立页面（component=Layout）
	MenuTypeMenu    = "menu"    // 菜单：可路由可渲染的实际页面
	MenuTypeButton  = "button"  // 按钮：不进入路由，只承载权限码（permission 字段）
)

// SysMenu 统一菜单树节点
//
// 三种类型通过 type 字段区分，通过 parent_id 形成层级关系：
//
//	type=catalog : 目录节点；component 通常为 "Layout"，path 为前端路由前缀（如 /system）
//	type=menu    : 菜单节点；path/name/component/icon 均生效，对应一个前端页面
//	type=button  : 按钮节点；只用 permission 字段做权限码（如 user:add），
//	               path/component 留空，不参与路由注册
//
// menu 的 Children 可包含 button 节点，用于前端 v-permission 指令控制按钮显隐。
// catalog 的 Children 可包含 menu 或其他 catalog 节点。
type SysMenu struct {
	Base
	Type       string     `gorm:"size:16;not null;default:menu;index" json:"type"`       // catalog / menu / button
	ParentID   uint       `gorm:"index"                                  json:"parent_id"` // 父节点 ID，0 表示根
	Path       string     `gorm:"size:255"                               json:"path"`      // 路由路径；catalog 如 /system，menu 如 user
	Name       string     `gorm:"size:64;index"                          json:"name"`      // 路由 name；button 为按钮名
	Component  string     `gorm:"size:255"                               json:"component"` // 前端组件路径；catalog 用 "Layout"
	Redirect   string     `gorm:"size:255"                               json:"redirect"`  // 重定向路径
	Permission string     `gorm:"size:128;index"                         json:"permission"` // 权限码；button 必有值，如 user:add
	Title      string     `gorm:"size:64"                                json:"title"`      // 显示标题
	Icon       string     `gorm:"size:64"                                json:"icon"`       // 图标名（element-ui icon）
	Sort       int        `gorm:"default:0"                              json:"sort"`       // 排序序号
	Hidden     bool       `gorm:"default:false"                          json:"hidden"`     // 是否隐藏
	KeepAlive  bool       `gorm:"default:false"                          json:"keep_alive"` // 是否缓存页面
	Children   []SysMenu  `gorm:"-"                                      json:"children,omitempty"` // 虚拟字段，仅用于 JSON 序列化时组装树
}

// ---------- API 元数据 ----------

// SysApi 后端 API 元数据
// 供"角色授权"界面列出可选 API，Casbin 真实策略（p 规则）写入 casbin_rule 表。
// 角色通过 sys_role_apis 中间表关联 API。
type SysApi struct {
	Base
	Path   string `gorm:"size:255;index"  json:"path"`   // 接口路径，如 /api/v1/system/user
	Method string `gorm:"size:16;index"   json:"method"` // HTTP 方法：GET/POST/PUT/DELETE
	Group  string `gorm:"size:64"          json:"group"`  // 分组标签，如 user / role / menu
	Desc   string `gorm:"size:255"         json:"desc"`   // 中文描述
}

// ---------- 安装状态 ----------

// SysInstall 安装状态表（仅一行，标记系统是否已完成初始化安装）
type SysInstall struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Version     string    `gorm:"size:32"    json:"version"`
	InstalledAt time.Time `json:"installed_at"`
	DBVersion   string    `gorm:"size:32"    json:"db_version"`
}

// ---------- 注册 ----------

// CoreModels 返回核心 Model 列表（供 installer 注册，参与 AutoMigrate）
// 注意：casbin_rule 表由 gorm-adapter 在 casbin.Setup() 时自动建表，不在此列表。
func CoreModels() []interface{} {
	return []interface{}{
		&SysUser{}, &SysRole{}, &SysMenu{},
		&SysApi{}, &SysInstall{},
	}
}
