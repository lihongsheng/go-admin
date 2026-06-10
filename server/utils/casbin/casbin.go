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

// SuperAdminRole 超级管理员角色 code；中间件直通
const SuperAdminRole = "super_admin"

// RBAC 模型：sub/obj/act 三段 + g 角色继承 + keyMatch2 路径匹配（兼容 :id 占位符）
const rbacModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && r.act == p.act
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

// userSub 把数字 user_id 拼成 "u:1" 这种 sub，便于和角色 code 共存
func userSub(uid uint) string { return "u:" + strconv.FormatUint(uint64(uid), 10) }

// AssignRole 给用户挂角色：g(u:1, role_code)
func AssignRole(uid uint, roleCode string) error {
	if !Ready() {
		return nil
	}
	_, err := E().AddGroupingPolicy(userSub(uid), roleCode)
	return err
}

// ReplaceUserRoles 把用户的角色替换为给定 codes
func ReplaceUserRoles(uid uint, roleCodes []string) error {
	if !Ready() {
		return nil
	}
	sub := userSub(uid)
	if _, err := E().DeleteRolesForUser(sub); err != nil {
		return err
	}
	for _, c := range roleCodes {
		if _, err := E().AddRoleForUser(sub, c); err != nil {
			return err
		}
	}
	return nil
}

// AddPolicy 角色拥有某 API 权限：p(role_code, /api/x, GET)
// 重复添加返回 (false, nil)，不算错。
func AddPolicy(role, path, method string) (bool, error) {
	if !Ready() {
		return false, nil
	}
	return E().AddPolicy(role, path, method)
}

// ReplaceRolePolicies 把角色的策略全量替换为给定 (path, method) 列表
func ReplaceRolePolicies(role string, items [][2]string) error {
	if !Ready() {
		return nil
	}
	// 先移除当前 role 全部策略
	if _, err := E().RemoveFilteredPolicy(0, role); err != nil {
		return err
	}
	for _, it := range items {
		if _, err := E().AddPolicy(role, it[0], it[1]); err != nil {
			return err
		}
	}
	return nil
}

// GetRolePolicies 取角色全部策略，返回 [(path, method), ...]
func GetRolePolicies(role string) [][2]string {
	if !Ready() {
		return nil
	}
	rs, _ := E().GetFilteredPolicy(0, role)
	out := make([][2]string, 0, len(rs))
	for _, r := range rs {
		if len(r) >= 3 {
			out = append(out, [2]string{r[1], r[2]})
		}
	}
	return out
}

// RolesOfUser 取用户当前在 casbin 里的角色列表
func RolesOfUser(uid uint) []string {
	if !Ready() {
		return nil
	}
	rs, _ := E().GetRolesForUser(userSub(uid))
	return rs
}

// Enforce 真正的校验：任一角色匹配即通过；super_admin 直通
func Enforce(uid uint, path, method string) (bool, error) {
	if !Ready() {
		// enforcer 没就绪时拒绝（防止误放行），由调用方决定语义
		return false, fmt.Errorf("casbin not ready")
	}
	roles := RolesOfUser(uid)
	for _, r := range roles {
		if r == SuperAdminRole {
			return true, nil
		}
	}
	for _, r := range roles {
		ok, err := E().Enforce(r, path, method)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// RemoveRolePolicies 移除角色的全部 API 策略 p(role, *, *)
func RemoveRolePolicies(role string) error {
	if !Ready() {
		return nil
	}
	_, err := E().RemoveFilteredPolicy(0, role)
	return err
}

// RemoveRoleFromUsers 移除所有用户与该角色的绑定 g(*, role)
func RemoveRoleFromUsers(role string) error {
	if !Ready() {
		return nil
	}
	_, err := E().RemoveFilteredGroupingPolicy(1, role)
	return err
}

// RemoveUserRoles 移除某个用户的全部角色绑定 g(u:<uid>, *)
// 用户被删除时由 service 层调用，避免 casbin_rule 中留下脏数据。
func RemoveUserRoles(uid uint) error {
	if !Ready() {
		return nil
	}
	_, err := E().DeleteRolesForUser(userSub(uid))
	return err
}

// MigrateRoleCode 角色 code 变更后，将原有用户-角色绑定迁移到新 code
func MigrateRoleCode(oldCode, newCode string) error {
	if !Ready() {
		return nil
	}
	// 找到所有绑定了旧 code 的用户
	policies, err := E().GetFilteredGroupingPolicy(1, oldCode)
	if err != nil {
		return err
	}
	for _, p := range policies {
		if len(p) >= 2 {
			// 添加新绑定 g(user, newCode)
			if _, err := E().AddGroupingPolicy(p[0], newCode); err != nil {
				return err
			}
		}
	}
	// 移除旧绑定
	_, err = E().RemoveFilteredGroupingPolicy(1, oldCode)
	return err
}
