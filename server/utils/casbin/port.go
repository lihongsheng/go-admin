// Package casbin port.go：暴露给 service 层用的端口接口
//
// 设计目的：避免 service 包直接 import 整个 casbin 包以及隐式依赖包级单例。
// service 层只持有 Port 接口指针，方便单测时替换为 mock。
package casbin

// Port 服务层使用的 Casbin 操作接口（按需逐步扩展）
type Port interface {
	ReplaceUserRoles(uid uint, roleCodes []string) error
	RemoveUserRoles(uid uint) error
	ReplaceRolePolicies(roleCode string, items [][2]string) error
	RemoveRolePolicies(roleCode string) error
	RemoveRoleFromUsers(roleCode string) error
	MigrateRoleCode(oldCode, newCode string) error
	GetRolePolicies(roleCode string) [][2]string
}

// NewPort 返回基于本包包级 enforcer 的默认 Port 实现
func NewPort() Port { return defaultPort{} }

type defaultPort struct{}

func (defaultPort) ReplaceUserRoles(uid uint, roleCodes []string) error {
	return ReplaceUserRoles(uid, roleCodes)
}

func (defaultPort) RemoveUserRoles(uid uint) error {
	return RemoveUserRoles(uid)
}

func (defaultPort) ReplaceRolePolicies(roleCode string, items [][2]string) error {
	return ReplaceRolePolicies(roleCode, items)
}

func (defaultPort) RemoveRolePolicies(roleCode string) error {
	return RemoveRolePolicies(roleCode)
}

func (defaultPort) RemoveRoleFromUsers(roleCode string) error {
	return RemoveRoleFromUsers(roleCode)
}

func (defaultPort) MigrateRoleCode(oldCode, newCode string) error {
	return MigrateRoleCode(oldCode, newCode)
}

func (defaultPort) GetRolePolicies(roleCode string) [][2]string {
	return GetRolePolicies(roleCode)
}
