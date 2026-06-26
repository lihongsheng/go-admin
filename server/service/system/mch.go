package system

import (
	"context"
	admin "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/enum"
	"github.com/lihongsheng/go-admin/server/model/system"
	repoSys "github.com/lihongsheng/go-admin/server/repo/system"
)

type MchService interface {
	Get(ctx context.Context, mchId int64) (*system.Merchant, error)
	Create(ctx context.Context, mch admin.MchCreateRequest) error
	Save(ctx context.Context, mch admin.MchCreateRequest) error
	Search(ctx context.Context, mch admin.MchQueryRequest) ([]*system.Merchant, error)
	Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error)
	ChangeStatus(ctx context.Context, mchId int64, status enum.MchStatus) error
	GetByMchNo(ctx context.Context, mchNo string) (*system.Merchant, error)
}
type mchService struct {
	mchRepo repoSys.MchRepo
}

// DefaultMch 包级单例
var DefaultMch MchService

func NewMchService(mchRepo repoSys.MchRepo) MchService {
	return &mchService{
		mchRepo: mchRepo,
	}
}

func (s *mchService) Get(ctx context.Context, mchId int64) (*system.Merchant, error) {
	return s.mchRepo.Get(ctx, mchId)
}

func (s *mchService) Create(ctx context.Context, mch admin.MchCreateRequest) error {
	return s.mchRepo.Create(ctx, mch)
}

func (s *mchService) Save(ctx context.Context, mch admin.MchCreateRequest) error {
	return s.mchRepo.Save(ctx, mch)
}

func (s *mchService) Search(ctx context.Context, mch admin.MchQueryRequest) ([]*system.Merchant, error) {
	return s.mchRepo.Search(ctx, mch)
}

func (s *mchService) Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error) {
	return s.mchRepo.Count(ctx, mch)
}

func (s *mchService) ChangeStatus(ctx context.Context, mchId int64, status enum.MchStatus) error {
	return s.mchRepo.ChangeStatus(ctx, mchId, status)
}

func (s *mchService) GetByMchNo(ctx context.Context, mchNo string) (*system.Merchant, error) {
	return s.mchRepo.GetByMchNo(ctx, mchNo)
}
