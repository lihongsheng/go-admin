// Package casbin Casbin enforcer 全局单例 + 策略管理
package casbin

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// Model：sub/obj/act 三段 + keyMatch2 路径匹配（兼容 :id 占位符）
// 仅基于角色 ID 做策略匹配，用户-角色关系由 sys_user_roles 表维护
const rbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act
`

var (
	mu       sync.RWMutex
	enforcer *casbin.SyncedEnforcer
)

// Setup 由 initialize 调用：基于 global.DB 初始化 enforcer + 自动建 casbin_rule
func Setup(db *gorm.DB) error {
	mu.Lock()
	defer mu.Unlock()

	// gorm-adapter 会自动 CREATE TABLE IF NOT EXISTS casbin_rule
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return fmt.Errorf("casbin adapter: %w", err)
	}
	m, err := model.NewModelFromString(rbacModel)
	if err != nil {
		return fmt.Errorf("casbin model: %w", err)
	}
	e, err := casbin.NewSyncedEnforcer(m, a)
	if err != nil {
		return fmt.Errorf("casbin enforcer: %w", err)
	}
	if err := e.LoadPolicy(); err != nil {
		return fmt.Errorf("casbin load policy: %w", err)
	}
	enforcer = e
	return nil
}

// E 取 enforcer
func E() *casbin.SyncedEnforcer {
	mu.RLock()
	defer mu.RUnlock()
	return enforcer
}

// Ready 是否就绪
func Ready() bool { return E() != nil }

// ridStr 把数字 role_id 转为字符串，用作 casbin sub
func ridStr(rid uint) string { return strconv.FormatUint(uint64(rid), 10) }

// AddPolicy 角色拥有某 API 权限：p(role_id, /api/x, GET)
// 重复添加返回 (false, nil)，不算错。
func AddPolicy(roleID uint, path, method string) (bool, error) {
	if !Ready() {
		return false, nil
	}
	return E().AddPolicy(ridStr(roleID), path, method)
}

// ReplaceRolePolicies 把角色的策略全量替换为给定 (path, method) 列表
func ReplaceRolePolicies(roleID uint, items [][2]string) error {
	if !Ready() {
		return nil
	}
	// 先移除当前 role 全部策略
	sub := ridStr(roleID)
	if _, err := E().RemoveFilteredPolicy(0, sub); err != nil {
		return err
	}
	for _, it := range items {
		if _, err := E().AddPolicy(sub, it[0], it[1]); err != nil {
			return err
		}
	}
	return nil
}

// GetRolePolicies 取角色全部策略，返回 [(path, method), ...]
func GetRolePolicies(roleID uint) [][2]string {
	if !Ready() {
		return nil
	}
	rs, _ := E().GetFilteredPolicy(0, ridStr(roleID))
	out := make([][2]string, 0, len(rs))
	for _, r := range rs {
		if len(r) >= 3 {
			out = append(out, [2]string{r[1], r[2]})
		}
	}
	return out
}

// Enforce 权限校验：检查指定角色是否有权限访问 path+method
// 匹配 p(role_id, path, method)
func Enforce(roleID uint, path, method string) (bool, error) {
	if !Ready() {
		return false, fmt.Errorf("casbin not ready")
	}
	ok, err := E().Enforce(ridStr(roleID), path, method)
	if err != nil {
		return false, err
	}
	return ok, nil
}

// RemoveRolePolicies 移除角色的全部 API 策略 p(role_id, *, *)
func RemoveRolePolicies(roleID uint) error {
	if !Ready() {
		return nil
	}
	_, err := E().RemoveFilteredPolicy(0, ridStr(roleID))
	return err
}
