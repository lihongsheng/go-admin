// Package enum 系统枚举常量
package enum

// UserStatus 用户状态
type UserStatus int8

const (
  UserStatusDisabled UserStatus = 0 // 禁用
  UserStatusEnabled  UserStatus = 1 // 启用
)

type SystemType int8

const (
  SystemTypePlatform SystemType = 0 // 平台
  SystemTypeMch      SystemType = 1 // 商户
  SystemTypeProxy    SystemType = 2 // 代理
)

var SystemTypeMap = map[SystemType]string{
  SystemTypePlatform: "平台",
  SystemTypeMch:      "商户",
  SystemTypeProxy:    "代理",
}

type SystemTypeInfo struct {
  SystemType SystemType
  Name       string
}

func GetAllSystemTypeInfo() []SystemTypeInfo {
  return []SystemTypeInfo{
    {SystemTypePlatform, SystemTypeMap[SystemTypePlatform]},
    {SystemTypeMch, SystemTypeMap[SystemTypeMch]},
    {SystemTypeProxy, SystemTypeMap[SystemTypeProxy]},
  }
}
