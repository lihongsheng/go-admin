package initialize

import (
	"go-admin/server/global"
	"go-admin/server/utils/casbin"
)

// SetupCasbin 在 global.DB 就绪后初始化 enforcer + 自动建 casbin_rule 表
func SetupCasbin() {
	if global.DB == nil {
		return
	}
	if err := casbin.Setup(global.DB); err != nil {
		getLogger().Error("casbin setup: " + err.Error())
		return
	}
	getLogger().Info("casbin enforcer ready")
}
