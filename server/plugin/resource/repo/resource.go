// Package repo resource 插件仓储层
package repo

import (
	resourceModel "go-admin/server/plugin/resource/model"

	"gorm.io/gorm"
)

// ResourceRepo 云资源仓储接口
type ResourceRepo interface {
	Create(r *resourceModel.Resource) error
	Update(r *resourceModel.Resource) error
	Delete(ids []uint) error
	Get(id uint) (*resourceModel.Resource, error)
	Page(page, limit int, keyword string) ([]resourceModel.Resource, int64, error)
}

// NewResourceRepo 构造 ResourceRepo
func NewResourceRepo(db *gorm.DB) ResourceRepo { return &resourceRepo{db: db} }

type resourceRepo struct{ db *gorm.DB }

func (r *resourceRepo) Create(res *resourceModel.Resource) error {
	return r.db.Create(res).Error
}

func (r *resourceRepo) Update(res *resourceModel.Resource) error {
	return r.db.Model(res).Updates(map[string]interface{}{
		"name":      res.Name,
		"type":      res.Type,
		"status":    res.Status,
		"status_at": res.StatusAt,
	}).Error
}

func (r *resourceRepo) Get(id uint) (*resourceModel.Resource, error) {
	var res resourceModel.Resource
	if err := r.db.First(&res, id).Error; err != nil {
		return nil, err
	}
	return &res, nil
}

// Delete 按 ID 批量删除（单个 ID 同样走此方法）
func (r *resourceRepo) Delete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Delete(&resourceModel.Resource{}, ids).Error
}

// Page 分页查询；keyword 匹配资源名称
func (r *resourceRepo) Page(page, limit int, keyword string) ([]resourceModel.Resource, int64, error) {
	var list []resourceModel.Resource
	var total int64

	q := r.db.Model(&resourceModel.Resource{})
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	err := q.Order("id DESC").Offset((page - 1) * limit).Limit(limit).Find(&list).Error
	return list, total, err
}
