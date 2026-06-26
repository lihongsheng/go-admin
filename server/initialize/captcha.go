package initialize

import (
	"github.com/lihongsheng/go-admin/server/global"
	"github.com/lihongsheng/go-admin/server/utils/captcha"
)

// InitCaptcha 初始化验证码存储
func InitCaptcha() {
	drive := global.Cfg.Captcha.Drive
	if drive == "" {
		drive = "memory"
	}
	if drive == "redis" {
		if global.Redis == nil {
			return // redis 未启用，退回默认 memory
		}
		captcha.SetupRedisStore(global.Redis)
	}
	// memory 模式下不做任何操作（captcha init() 已使用默认 memory store）
}
