// Package initialize service / repo 装配
package initialize

import (
	"github.com/lihongsheng/go-admin/server/global"
	"github.com/lihongsheng/go-admin/server/log"
	repoExampleNote "github.com/lihongsheng/go-admin/server/plugin/example/repo"
	serviceExampleNote "github.com/lihongsheng/go-admin/server/plugin/example/service"
	repoSys "github.com/lihongsheng/go-admin/server/repo/system"
	serviceBase "github.com/lihongsheng/go-admin/server/service/base"
	serviceInstall "github.com/lihongsheng/go-admin/server/service/install"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	casbinUtil "github.com/lihongsheng/go-admin/server/utils/casbin"
)

// InitInstallService 仅装配 install service（与 DB 是否就绪无关）
// install service 自己不依赖 DB，由 handler 在请求时传入 DB 配置。
// 在 main.go 启动早期调用一次即可。
func InitInstallService() {
	if serviceInstall.Default != nil {
		return
	}
	svc := serviceInstall.NewService()
	// 设置 logger
	if setter, ok := svc.(interface{ SetLogger(log.Logger) }); ok {
		setter.SetLogger(log.Global())
	}
	// 注册"安装完成 + DB 就绪后"的回调：重新装配所有依赖 global.DB 的 service
	svc.OnReady(InitDBServices)
	serviceInstall.Default = svc
}

// InitDBServices 装配所有依赖 global.DB 的 service 单例。
//
// 触发时机：
//  1. main.go 启动期：当 DB 已配置且连接成功时直接调用一次；
//  2. 安装向导完成、global.DB 被替换后由 install service 回调触发。
//
// 这是幂等的——每次都用最新的 global.DB 重新装配。
func InitDBServices() {
	if global.DB == nil {
		// DB 尚未就绪（处于安装模式）；等安装完成后回调再装配
		return
	}
	casbinPort := casbinUtil.NewPort()

	// repo
	userRepo := repoSys.NewUserRepo(global.DB)
	roleRepo := repoSys.NewRoleRepo(global.DB)
	menuRepo := repoSys.NewMenuRepo(global.DB)
	mchRepo := repoSys.NewMchRepo(global.DB)
	noteRepo := repoExampleNote.NewNoteRepo(global.DB)

	// system service
	serviceSys.DefaultUser = serviceSys.NewUserService(userRepo)
	serviceSys.DefaultRole = serviceSys.NewRoleService(roleRepo, menuRepo, casbinPort)
	serviceSys.DefaultMenu = serviceSys.NewMenuService(menuRepo, userRepo)
	serviceSys.DefaultMch = serviceSys.NewMchService(mchRepo)

	// base service
	serviceBase.Default = serviceBase.NewService(userRepo)

	// example plugin service
	serviceExampleNote.DefaultNote = serviceExampleNote.NewNoteService(noteRepo)

	log.Info("db services initialized")
}
