package system

import (
	"github.com/lihongsheng/go-admin/server/enum"
	"github.com/lihongsheng/go-admin/server/model/system"

	"gorm.io/gorm"
)

// ApiRepo API 仓储接口
type ApiRepo interface {
	Create(a *system.SysApi) error
	Update(id uint, patch map[string]any) error
	Delete(id uint) error
	List(group string, systemType int) ([]system.SysApi, error)
	ListBySystemType(systemType enum.SystemType) ([]system.SysApi, error)
	FindByIDs(ids []uint) ([]system.SysApi, error)
	FindByPathMethod(path, method string) (*system.SysApi, error)
}

// NewApiRepo 构造 ApiRepo
func NewApiRepo(db *gorm.DB) ApiRepo { return &apiRepo{db: db} }

type apiRepo struct{ db *gorm.DB }

func (r *apiRepo) Create(a *system.SysApi) error {
	return r.db.Create(a).Error
}

func (r *apiRepo) Update(id uint, patch map[string]any) error {
	return r.db.Model(&system.SysApi{}).Where("id=?", id).Updates(patch).Error
}

func (r *apiRepo) Delete(id uint) error {
	return r.db.Delete(&system.SysApi{}, id).Error
}

func (r *apiRepo) List(group string, systemType int) ([]system.SysApi, error) {
	q := r.db.Model(&system.SysApi{})
	if group != "" {
		q = q.Where("`group` = ?", group)
	}
	if systemType > 0 {
		q = q.Where("system_type = ?", systemType)
	}
	var list []system.SysApi
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *apiRepo) ListBySystemType(systemType enum.SystemType) ([]system.SysApi, error) {
	var list []system.SysApi
	if err := r.db.Where("system_type = ?", systemType).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *apiRepo) FindByIDs(ids []uint) ([]system.SysApi, error) {
	var list []system.SysApi
	if len(ids) == 0 {
		return list, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *apiRepo) FindByPathMethod(path, method string) (*system.SysApi, error) {
	var a system.SysApi
	if err := r.db.Where("path = ? AND method = ?", path, method).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}
