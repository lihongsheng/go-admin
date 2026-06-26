package enum

// MchStatus 商户状态
type MchStatus int8

const (
	MchStatusDisabled UserStatus = 2 // 禁用
	MchStatusEnabled  UserStatus = 1 // 启用
)
