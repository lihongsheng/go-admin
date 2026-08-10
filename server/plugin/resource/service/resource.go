// Package service resource 插件业务服务
package service

import (
	"errors"
	"time"

	dtoResource "go-admin/server/plugin/resource/dto"
	resourceModel "go-admin/server/plugin/resource/model"
	repoResource "go-admin/server/plugin/resource/repo"

	"gorm.io/gorm"
)

// ResourceService 云资源业务接口
type ResourceService interface {
	Create(req dtoResource.ResourceCreateReq) (*resourceModel.Resource, error)
	Update(req dtoResource.ResourceUpdateReq) (*resourceModel.Resource, error)
	Delete(ids []uint) error
	Get(id uint) (*resourceModel.Resource, error)
	List(page, limit int, keyword string) ([]resourceModel.Resource, int64, error)
}

// NewResourceService 构造 ResourceService
func NewResourceService(repo repoResource.ResourceRepo) ResourceService {
	return &resourceService{repo: repo}
}

type resourceService struct {
	repo repoResource.ResourceRepo
}

// DefaultResource 包级单例（由插件 InitServices 装配）
var DefaultResource ResourceService

func (s *resourceService) Create(req dtoResource.ResourceCreateReq) (*resourceModel.Resource, error) {
	if req.Name == "" {
		return nil, errors.New("资源名称不能为空")
	}
	if req.Type == "" {
		return nil, errors.New("资源类型不能为空")
	}
	now := time.Now()
	res := &resourceModel.Resource{
		Name:     req.Name,
		Type:     req.Type,
		Status:   req.Status,
		StatusAt: now,
	}
	if res.Status == 0 {
		res.Status = resourceModel.StatusRunning
	}
	if err := s.repo.Create(res); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *resourceService) Update(req dtoResource.ResourceUpdateReq) (*resourceModel.Resource, error) {
	cur, err := s.repo.Get(req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("资源不存在")
		}
		return nil, err
	}
	cur.Name = req.Name
	cur.Type = req.Type
	if req.Status > 0 {
		// 仅状态发生变更时刷新状态更新时间
		if cur.Status != req.Status {
			cur.Status = req.Status
			cur.StatusAt = time.Now()
		}
	}
	if err := s.repo.Update(cur); err != nil {
		return nil, err
	}
	return cur, nil
}

func (s *resourceService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的资源")
	}
	return s.repo.Delete(ids)
}

func (s *resourceService) Get(id uint) (*resourceModel.Resource, error) {
	return s.repo.Get(id)
}

func (s *resourceService) List(page, limit int, keyword string) ([]resourceModel.Resource, int64, error) {
	return s.repo.Page(page, limit, keyword)
}
