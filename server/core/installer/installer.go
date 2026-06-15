package installer

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"go-admin/server/model/system"
	"go-admin/server/utils/casbin"

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
// progress 可为 nil；非 nil 时实时推送步骤
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

// seedCore 写入默认角色 / 管理员 / 菜单 / API / Casbin 策略
func seedCore(db *gorm.DB, admin AdminSeed) error {
	var superRole system.SysRole
	err := db.Transaction(func(tx *gorm.DB) error {
		// 1) 角色
		superRole = system.SysRole{
			Name: "超级管理员", Code: casbin.SuperAdminRole,
			Remark: "拥有全部权限", Status: 1,
			DefaultRouter: "/dashboard",
		}
		if err := tx.Create(&superRole).Error; err != nil {
			return err
		}

		// 2) 管理员用户
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
			Roles:    []system.SysRole{superRole},
		}
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 3) 菜单树
		if err := createMenusTree(tx, defaultMenus(), 0); err != nil {
			return err
		}

		// 4) API 元数据
		if err := tx.Create(defaultApis()).Error; err != nil {
			return err
		}

		// 5) super_admin 拥有所有菜单
		var allMenus []system.SysMenu
		if err := tx.Find(&allMenus).Error; err != nil {
			return err
		}
		if err := tx.Model(&superRole).Association("Menus").Replace(allMenus); err != nil {
			return err
		}

		// 把 admin user_id 暂存到外层，供事务后写 casbin g 关系
		_ = user
		return nil
	})
	if err != nil {
		return err
	}

	// 6) Casbin 策略：super_admin 拥有所有 API；admin 用户绑 super_admin 角色
	var allApis []system.SysApi
	if err := db.Find(&allApis).Error; err != nil {
		return err
	}
	items := make([][2]string, 0, len(allApis))
	for _, a := range allApis {
		items = append(items, [2]string{a.Path, a.Method})
	}
	if err := casbin.ReplaceRolePolicies(superRole.Code, items); err != nil {
		return err
	}
	var adminUser system.SysUser
	if err := db.Where("username = ?", admin.Username).First(&adminUser).Error; err != nil {
		return err
	}
	if err := casbin.ReplaceUserRoles(adminUser.ID, []string{superRole.Code}); err != nil {
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
