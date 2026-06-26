package system

import (
	"time"
)

const TableNameMerchant = "merchant"

// Merchant 公司信息
type Merchant struct {
	ID        int64     `gorm:"column:id;type:int;primaryKey;autoIncrement:true" json:"id"`
	MchName   string    `gorm:"column:mch_name;type:varchar(100);not null;comment:公司名字" json:"mch_name"`       // 公司名字
	Linker    string    `gorm:"column:linker;type:varchar(20);not null;comment:联系人" json:"linker"`              // 联系人
	Phone     string    `gorm:"column:phone;type:varchar(20);not null;comment:联系电话" json:"phone"`              // 联系电话
	Email     string    `gorm:"column:email;type:varchar(100);not null;comment:联系电话" json:"email"`             // 联系电话
	MchNo     string    `gorm:"column:mch_no;type:varchar(100);not null;comment:编号" json:"mch_no"`               // 编号
	Status    int64     `gorm:"column:status;type:tinyint;not null;default:1;comment:2 停用 1 正常" json:"status"` // 2 停用 1 正常
	Address   string    `gorm:"column:address;type:varchar(200);not null;comment:地址" json:"address"`             // 地址
	Reason    string    `gorm:"column:reason;type:varchar(200);not null;comment:备注" json:"reason"`               // 备注
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName Merchant's table name
func (*Merchant) TableName() string {
	return TableNameMerchant
}
