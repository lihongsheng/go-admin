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

func NewMchService(mchRepo repoSys.MchRepo) MchService {
	return &mchService{
		mchRepo: mchRepo,
	}
}
