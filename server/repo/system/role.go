package system

import (
	"go-admin/server/model/system"

	"gorm.io/gorm"
)

// RoleRepo 角色仓储接口
type RoleRepo interface {
	Create(r *system.SysRole) error
	Update(id uint, patch map[string]any) error
	Delete(id uint) error
	GetByID(id uint) (*system.SysRole, error)
	GetByIDWithMenus(id uint) (*system.SysRole, error)
	List() ([]system.SysRole, error)
	ReplaceMenus(roleID uint, menus []system.SysMenu) error
}

// NewRoleRepo 构造 RoleRepo
func NewRoleRepo(db *gorm.DB) RoleRepo { return &roleRepo{db: db} }

type roleRepo struct{ db *gorm.DB }

func (r *roleRepo) Create(role *system.SysRole) error {
	return r.db.Create(role).Error
}

func (r *roleRepo) Update(id uint, patch map[string]any) error {
	return r.db.Model(&system.SysRole{}).Where("id=?", id).Updates(patch).Error
}

func (r *roleRepo) Delete(id uint) error {
	return r.db.Delete(&system.SysRole{}, id).Error
}

func (r *roleRepo) GetByID(id uint) (*system.SysRole, error) {
	var role system.SysRole
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) GetByIDWithMenus(id uint) (*system.SysRole, error) {
	var role system.SysRole
	if err := r.db.Preload("Menus").First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) List() ([]system.SysRole, error) {
	var list []system.SysRole
	if err := r.db.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *roleRepo) ReplaceMenus(roleID uint, menus []system.SysMenu) error {
	return r.db.Model(&system.SysRole{Base: system.Base{ID: roleID}}).Association("Menus").Replace(menus)
}

// FindRoleIDsByUserID 通过用户 ID 找出全部角色 ID（供 menu 服务用）
func FindRoleIDsByUserID(db *gorm.DB, userID uint) ([]uint, error) {
	var user system.SysUser
	if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(user.Roles))
	for _, role := range user.Roles {
		ids = append(ids, role.ID)
	}
	return ids, nil
}
