package installer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/lihongsheng/go-admin/server/enum"
	"github.com/lihongsheng/go-admin/server/model/system"
	"github.com/lihongsheng/go-admin/server/utils/casbin"
	"github.com/lihongsheng/go-admin/server/utils/genid"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Version 系统版本号（写入 sys_install）
const Version = "0.1.0"

// AdminSeed 安装向导提交的管理员账号
type AdminSeed struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

// SeedDeps 传给插件 Seed 的依赖
type SeedDeps struct {
	DB    *gorm.DB
	Admin AdminSeed
}

// Step 单张表迁移进度
type Step struct {
	Table string `json:"table"`
	State string `json:"state"` // pending / migrating / done / failed / seeding
	Err   string `json:"err,omitempty"`
	At    int64  `json:"at"`
}

// Install 执行完整安装流程
// progress 可为 nil；非 nil 时实时推送进度
func Install(ctx context.Context, db *gorm.DB, admin AdminSeed, progress chan<- Step) error {
	if err := AcquireProcessLock(); err != nil {
		return err
	}
	defer ReleaseProcessLock()

	if err := EnsureNotInstalled(db); err != nil {
		return err
	}

	models := AllModels()
	if len(models) == 0 {
		return errors.New("no models registered")
	}

	emit := func(s Step) {
		s.At = time.Now().UnixMilli()
		if progress != nil {
			select {
			case progress <- s:
			case <-ctx.Done():
			}
		}
	}

	// 1. AutoMigrate 逐表执行，便于上报进度
	for _, m := range models {
		name := tableName(db, m)
		emit(Step{Table: name, State: "migrating"})
		if err := db.AutoMigrate(m); err != nil {
			emit(Step{Table: name, State: "failed", Err: err.Error()})
			return fmt.Errorf("migrate %s: %w", name, err)
		}
		emit(Step{Table: name, State: "done"})
	}

	// 2. Casbin 初始化（自动建 casbin_rule 表）
	emit(Step{Table: "casbin_rule", State: "migrating"})
	if err := casbin.Setup(db); err != nil {
		emit(Step{Table: "casbin_rule", State: "failed", Err: err.Error()})
		return err
	}
	emit(Step{Table: "casbin_rule", State: "done"})

	// 3. 核心种子数据 + casbin 策略
	emit(Step{Table: "_seed_core", State: "seeding"})
	if err := seedCore(db, admin); err != nil {
		emit(Step{Table: "_seed_core", State: "failed", Err: err.Error()})
		return err
	}
	emit(Step{Table: "_seed_core", State: "done"})

	// 4. 插件种子数据
	deps := SeedDeps{DB: db, Admin: admin}
	for i, fn := range AllSeeds() {
		label := fmt.Sprintf("_seed_plugin_%d", i)
		emit(Step{Table: label, State: "seeding"})
		if err := fn(deps); err != nil {
			emit(Step{Table: label, State: "failed", Err: err.Error()})
			return err
		}
		emit(Step{Table: label, State: "done"})
	}

	// 5. 写入 sys_install 标记完成
	if err := db.Create(&system.SysInstall{
		Version:     Version,
		InstalledAt: time.Now(),
		DBVersion:   Version,
	}).Error; err != nil {
		return err
	}
	emit(Step{Table: "_install_done", State: "done"})
	return nil
}

// seedCore 写入默认角色 / 管理员 / 商户 / 商户管理员 / 菜单 / API / Casbin 策略
func seedCore(db *gorm.DB, admin AdminSeed) error {
	var superRole system.SysRole
	var mchAdminRole system.SysRole
	var merchant system.Merchant

	err := db.Transaction(func(tx *gorm.DB) error {
		// 1) 平台超级管理员角色
		superRole = system.SysRole{
			Name: "超级管理员",
			Remark: "拥有全部权限", Status: 1,
			DefaultRouter: "/dashboard",
		}
		if err := tx.Create(&superRole).Error; err != nil {
			return err
		}

		// 2) 平台管理员用户
		hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		nickname := admin.Nickname
		if nickname == "" {
			nickname = admin.Username
		}
		user := system.SysUser{
			Username: admin.Username,
			Password: string(hash),
			Nickname: nickname,
			Email:    admin.Email,
			Status:   1,
			MchID:    0,
			SystemType: enum.SystemTypePlatform,
			Roles:    []system.SysRole{superRole},
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 3) 菜单树（平台）
		if err := createMenusTree(tx, defaultMenus(), 0); err != nil {
			return err
		}

		// 4) API 元数据（平台）
		if err := tx.Create(defaultApis()).Error; err != nil {
			return err
		}

		// 5) super_admin 拥有所有平台菜单
		var allMenus []system.SysMenu
		if err := tx.Find(&allMenus).Error; err != nil {
			return err
		}
		if err := tx.Model(&superRole).Association("Menus").Replace(allMenus); err != nil {
			return err
		}

		// 6) 默认商户
		merchant = defaultMerchant()
		merchant.MchNo = "M" + genid.GenDeviceID.Generate0X()
		if err := tx.Create(&merchant).Error; err != nil {
			return err
		}

		// 7) 商户管理员角色
		mchAdminRole = system.SysRole{
			Name:          "商户管理员",
			Remark:        "商户管理员，拥有用户管理和角色管理权限",
			Status:        1,
			DefaultRouter: "/dashboard",
			MchID:         merchant.ID,
			SystemType:    enum.SystemTypeMch,
		}
		if err := tx.Create(&mchAdminRole).Error; err != nil {
			return err
		}

		// 8) 商户管理员用户
		mchHash, err := bcrypt.GenerateFromPassword([]byte("mchadmin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		mchUser := system.SysUser{
			Username:   "mchadmin",
			Password:   string(mchHash),
			Nickname:   "商户管理员",
			Email:      "",
			Status:     1,
			MchID:      merchant.ID,
			SystemType: enum.SystemTypeMch,
			Roles:      []system.SysRole{mchAdminRole},
		}
		if err := tx.Create(&mchUser).Error; err != nil {
			return err
		}

		// 9) 商户管理员菜单树（systemType=Mch）
		if err := createMenusTree(tx, merchantAdminMenus(), 0); err != nil {
			return err
		}

		// 10) 商户管理员 API 元数据（systemType=Mch）
		if err := tx.Create(merchantAdminApis()).Error; err != nil {
			return err
		}

		// 11) 商户管理员角色拥有对应 systemType 的菜单
		var mchMenus []system.SysMenu
		if err := tx.Where("system_type = ?", enum.SystemTypeMch).Find(&mchMenus).Error; err != nil {
			return err
		}
		if err := tx.Model(&mchAdminRole).Association("Menus").Replace(mchMenus); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	// 12) Casbin 策略：超级管理员角色拥有所有平台 API
	var allApis []system.SysApi
	if err := db.Where("system_type = ?", enum.SystemTypePlatform).Find(&allApis).Error; err != nil {
		return err
	}
	items := make([][2]string, 0, len(allApis))
	for _, a := range allApis {
		items = append(items, [2]string{a.Path, a.Method})
	}
	if err := casbin.ReplaceRolePolicies(superRole.ID, items); err != nil {
		return err
	}

	// 13) Casbin 策略：商户管理员角色拥有商户 API
	var mchApis []system.SysApi
	if err := db.Where("system_type = ?", enum.SystemTypeMch).Find(&mchApis).Error; err != nil {
		return err
	}
	mchItems := make([][2]string, 0, len(mchApis))
	for _, a := range mchApis {
		mchItems = append(mchItems, [2]string{a.Path, a.Method})
	}
	if err := casbin.ReplaceRolePolicies(mchAdminRole.ID, mchItems); err != nil {
		return err
	}

	return nil
}

// createMenusTree 递归插入菜单树
func createMenusTree(tx *gorm.DB, nodes []system.SysMenu, parent uint) error {
	for i := range nodes {
		n := nodes[i]
		children := n.Children
		n.Children = nil
		n.ParentID = parent
		if err := tx.Create(&n).Error; err != nil {
			return err
		}
		if len(children) > 0 {
			if err := createMenusTree(tx, children, n.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// tableName 反射拿到 Model 对应的表名
func tableName(db *gorm.DB, m interface{}) string {
	// 通过正常 GORM 链式调用构建 Statement，避免内部字段未初始化导致 panic
	stmt := db.Session(&gorm.Session{NewDB: true}).Model(m).Statement
	if err := stmt.Parse(m); err == nil && stmt.Schema != nil {
		return stmt.Schema.Table
	}
	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return strings.ToLower(t.Name())
}
