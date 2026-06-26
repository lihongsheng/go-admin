package enum

// MchStatus 商户状态
type MchStatus int8

const (
	MchStatusDisabled MchStatus = 2 // 禁用
	MchStatusEnabled  MchStatus = 1 // 启用
)
