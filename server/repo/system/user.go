// Package system 系统模块仓储层
package system

import (
	"github.com/lihongsheng/go-admin/server/model/system"

	"gorm.io/gorm"
)

// UserRepo 用户仓储接口
type UserRepo interface {
	Create(u *system.SysUser) error
	Update(id uint, patch map[string]any) error
	Delete(id uint) error
	GetByID(id uint, preloadRoles bool) (*system.SysUser, error)
	GetByUsername(username string) (*system.SysUser, error)
	List(keyword string, page, size int, mchID int64, systemType int) ([]system.SysUser, int64, error)
	ReplaceRoles(uid uint, roles []system.SysRole) error
	FindRolesByIDs(ids []uint, mchID int64) ([]system.SysRole, error)
}

// NewUserRepo 构造 UserRepo
func NewUserRepo(db *gorm.DB) UserRepo { return &userRepo{db: db} }

type userRepo struct{ db *gorm.DB }

func (r *userRepo) Create(u *system.SysUser) error {
	return r.db.Create(u).Error
}

func (r *userRepo) Update(id uint, patch map[string]any) error {
	return r.db.Model(&system.SysUser{}).Where("id=?", id).Updates(patch).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&system.SysUser{}, id).Error
}

func (r *userRepo) GetByID(id uint, preloadRoles bool) (*system.SysUser, error) {
	var u system.SysUser
	q := r.db
	if preloadRoles {
		q = q.Preload("Roles")
	}
	if err := q.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetByUsername(username string) (*system.SysUser, error) {
	var u system.SysUser
	if err := r.db.Preload("Roles").Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) List(keyword string, page, size int, mchID int64, systemType int) ([]system.SysUser, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	q := r.db.Model(&system.SysUser{}).Preload("Roles")
	if keyword != "" {
		q = q.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if mchID > 0 {
		q = q.Where("mch_id = ?", mchID)
	}
	if systemType > 0 {
		q = q.Where("system_type = ?", systemType)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []system.SysUser
	if err := q.Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *userRepo) ReplaceRoles(uid uint, roles []system.SysRole) error {
	return r.db.Model(&system.SysUser{Base: system.Base{ID: uid}}).Association("Roles").Replace(roles)
}

func (r *userRepo) FindRolesByIDs(ids []uint, mchID int64) ([]system.SysRole, error) {
	var roles []system.SysRole
	if len(ids) == 0 {
		return roles, nil
	}
	q := r.db.Where("id IN ?", ids)
	if mchID > 0 {
		q = q.Where("mch_id = ?", mchID)
	}
	if err := q.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
