package initialize

import (
	"go-admin/server/global"
	"go-admin/server/plugin"
)

// SyncOnBoot 启动期同步：
//   - 系统已安装 + global.DB 可用时，自动 AutoMigrate 全部已注册 Model
//     （含核心 + 插件），新插件首次被加载即建表
//   - 幂等 upsert 各插件的菜单 / API，避免重复
//   - 若插件目标表为空，调用 SeedTable 写入初始数据
//
// 系统未安装时不做事；安装流程结束后由 /install/init 自行完成。
func SyncOnBoot() {
	if !global.Installed.Load() || global.DB == nil {
		return
	}
	if err := plugin.SyncOnBoot(global.DB); err != nil {
		global.Logger.Sugar().Errorf("plugin SyncOnBoot: %v", err)
		return
	}
	global.Logger.Info("plugin sync on boot completed")
}
