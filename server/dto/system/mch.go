package system

import (
	"errors"
	"go-admin/server/enum"
)

type MchCreateRequest struct {
	ID      int64          `json:"id"`
	MchName string         `json:"mch_name"` // 公司名字
	Linker  string         `json:"linker"`   // 联系人
	Phone   string         `json:"phone"`    // 联系电话
	Email   string         `json:"email"`    // 联系电话
	Status  enum.MchStatus `json:"status"`   // 0 停用 1 正常
	Address string         `json:"address"`  // 地址
	Reason  string         `json:"reason"`   // 备注
}

func (m MchCreateRequest) Validate() error {
	if m.MchName == "" {
		return errors.New("商户名称不能为空")
	}
	if m.Linker == "" {
		return errors.New("联系人不能为空")
	}
	if m.Phone == "" {
		return errors.New("联系电话不能为空")
	}
	if m.Email == "" {
		return errors.New("邮箱不能为空")
	}
	if m.Status < 1 {
		return errors.New("商户状态不能为空")
	}
	return nil
}

type MchQueryRequest struct {
	MchName  string         `json:"mch_name" form:"mch_name"` // 公司名字
	MchNo    string         `json:"mch_no" form:"mch_no"`     // 编号
	Status   enum.MchStatus `json:"status" form:"status"`     // 0 停用 1 正常
	ID       int64          `json:"id" form:"id"`             // 商户ID
	PageSize int            `json:"limit" form:"limit"`       // 限制数量
	Page     int            `json:"page" form:"page"`         // 页大小
	IDList   []int64        `json:"id_list" form:"id_list"`   // 商户ID
	MchNos   []string       `json:"mch_nos" form:"mch_nos"`   // 编号
}

func (m MchQueryRequest) Validate() error {
	if m.PageSize <= 0 {
		return errors.New("页大小不能小于0")
	}
	if m.PageSize > 50 {
		return errors.New("页大小不能大于50")
	}
	if m.Page <= 0 {
		return errors.New("页数不能小于0")
	}
	return nil
}

type MchStatusRequest struct {
	ID     int64          `json:"id" form:"id" binding:"required,gt=0"`
	MchNo  string         `json:"mch_no" form:"mch_no"`                         // 编号
	Status enum.MchStatus `json:"status" form:"status" binding:"required,gt=0"` // 0 停用 1 正常
}
