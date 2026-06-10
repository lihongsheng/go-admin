package main

import (
	"flag"
	"fmt"
	"log"

	"go-admin/server/config"
	"go-admin/server/global"
	"go-admin/server/initialize"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	global.Cfg = cfg

	initialize.Logger(cfg.Log)
	defer global.Logger.Sync() //nolint:errcheck

	// install service 早于 DB 装配——安装向导本身无需 DB 就绪
	initialize.InitInstallService()

	// 尝试连接 DB；连不上不致命，进入安装向导模式
	if cfg.DB.Configured() {
		if err := initialize.GormConnect(); err != nil {
			global.Logger.Warn("db connect failed, entering install mode: " + err.Error())
		} else {
			initialize.DetectInstalled()
			initialize.SetupCasbin()
		}
	} else {
		global.Logger.Warn("db not configured, entering install mode")
	}

	// 加载插件（注册 Model 到 installer 注册中心 / 注册路由）
	initialize.LoadPlugins()

	// 已安装且 DB 就绪：启动期增量同步（新插件 Model/菜单/API 自动迁入）
	initialize.SyncOnBoot()

	// 装配依赖 global.DB 的 service 单例（DB 未就绪时这步空跑，等安装回调）
	initialize.InitDBServices()

	r := initialize.Router()
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	global.Logger.Sugar().Infof("server start at %s (installed=%v)", addr, global.Installed.Load())
	if err := r.Run(addr); err != nil {
		global.Logger.Fatal("server run: " + err.Error())
	}
}
