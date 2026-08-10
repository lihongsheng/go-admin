package system

import (
	"context"
	"errors"
	"fmt"
	admin "go-admin/server/dto/system"
	"go-admin/server/enum"
	"go-admin/server/model/system"
	"go-admin/server/utils/genid"
	"gorm.io/gorm"
	"time"
)

type MchRepo interface {
	Get(ctx context.Context, mchId int64) (*system.Merchant, error)
	Create(ctx context.Context, mch admin.MchCreateRequest) error
	Save(ctx context.Context, mch admin.MchCreateRequest) error
	Search(ctx context.Context, mch admin.MchQueryRequest) ([]*system.Merchant, error)
	Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error)
	ChangeStatus(ctx context.Context, mchId int64, status enum.MchStatus) error
	GetByMchNo(ctx context.Context, mchNo string) (*system.Merchant, error)
}

type mchRepoImpl struct{ db *gorm.DB }

func NewMchRepo(db *gorm.DB) MchRepo {
	return &mchRepoImpl{db}
}

func (m *mchRepoImpl) Get(ctx context.Context, mchId int64) (*system.Merchant, error) {
	var mdl system.Merchant
	err := m.db.WithContext(ctx).Where("id = ?", mchId).First(&mdl).Error
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}

func (m *mchRepoImpl) Create(ctx context.Context, mch admin.MchCreateRequest) (err error) {
	mdl := &system.Merchant{
		MchNo:   "M" + genid.GenDeviceID.Generate0X(),
		MchName: mch.MchName,
		Linker:  mch.Linker,
		Phone:   mch.Phone,
		Email:   mch.Email,
		Address: mch.Address,
		Reason:  mch.Reason,
		Status:  int64(mch.Status),
	}
	return m.db.WithContext(ctx).Create(mdl).Error
}

func (m *mchRepoImpl) Save(ctx context.Context, mch admin.MchCreateRequest) (err error) {
	var mdl *system.Merchant
	if mch.Validate() != nil {
		return mch.Validate()
	}
	if mch.ID == 0 {
		return errors.New("商户ID不能为空")
	}
	mdl, err = m.Get(ctx, mch.ID)
	if err != nil {
		return fmt.Errorf("商户不存在: %w", err)
	}
	mdl.MchName = mch.MchName
	mdl.Linker = mch.Linker
	mdl.Phone = mch.Phone
	mdl.Email = mch.Email
	mdl.Address = mch.Address
	mdl.Reason = mch.Reason
	mdl.Status = int64(mch.Status)
	mdl.UpdatedAt = time.Now()
	return m.db.WithContext(ctx).Save(mdl).Error
}

func (m *mchRepoImpl) Search(ctx context.Context, mch admin.MchQueryRequest) ([]*system.Merchant, error) {
	// 实现搜索功能，根据请求参数进行查询
	var merchants []*system.Merchant
	query := m.buildQuery(ctx, mch)
	if mch.PageSize < 1 {
		mch.PageSize = 10
	}
	if mch.Page > 0 {
		query = query.Offset(mch.PageSize * (mch.Page - 1))
	}
	if mch.PageSize > 0 {
		query = query.Limit(mch.PageSize)
	}
	err := query.Find(&merchants).Error
	return merchants, err
}
func (m *mchRepoImpl) buildQuery(ctx context.Context, mch admin.MchQueryRequest) *gorm.DB {
	query := m.db.WithContext(ctx).Model(&system.Merchant{})
	if mch.MchNo != "" {
		query = query.Where("mch_no LIKE ?", "%"+mch.MchNo+"%")
	}
	if mch.MchName != "" {
		query = query.Where("mch_name LIKE ?", "%"+mch.MchName+"%")
	}
	if mch.ID != 0 {
		query = query.Where("id = ?", mch.ID)
	}
	if mch.Status > 0 {
		query = query.Where("status = ?", mch.Status)
	}
	if mch.IDList != nil {
		query = query.Where("id IN ?", mch.IDList)
	}
	if len(mch.MchNos) > 0 {
		query = query.Where("mch_no IN ?", mch.MchNos)
	}
	return query
}
func (m *mchRepoImpl) Count(ctx context.Context, mch admin.MchQueryRequest) (int64, error) {
	var count int64
	query := m.buildQuery(ctx, mch)
	err := query.Count(&count).Error
	return count, err
}

func (m *mchRepoImpl) ChangeStatus(ctx context.Context, mchId int64, status enum.MchStatus) error {
	mdl, err := m.Get(ctx, mchId)
	if err != nil {
		return err
	}
	if mdl.Status == int64(status) {
		return nil
	}
	result := m.db.WithContext(ctx).Model(&system.Merchant{}).
		Where("id = ?", mchId).
		Update("status", status)

	return result.Error
}

func (m *mchRepoImpl) GetByMchNo(ctx context.Context, mchNo string) (*system.Merchant, error) {
	var mdl system.Merchant
	err := m.db.WithContext(ctx).Where("mch_no = ?", mchNo).First(&mdl).Error
	if err != nil {
		return nil, err
	}
	return &mdl, nil
}
