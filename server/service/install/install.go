// Package install 安装向导业务服务
//
// 把原本散落在 api/v1/install handler 里的"ensure db → open → install → save config →
// 重建 global.DB → 重新初始化 casbin"这套主流程下沉到 service。
// handler 只剩参数校验、SSE 适配，以及在成功后触发 service 重新装配。
package install

import (
	"context"
	"errors"
	"time"

	"go-admin/server/config"
	"go-admin/server/core/installer"
	dtoInstall "go-admin/server/dto/install"
	"go-admin/server/global"
	"go-admin/server/log"
	casbinUtil "go-admin/server/utils/casbin"

	"gorm.io/gorm"
)

// Service 安装向导服务接口
type Service interface {
	Status() installer.DBStatus
	CheckDB(req dtoInstall.CheckDBReq) dtoInstall.CheckDBResp
	// Install 同步执行完整安装；progress 可为 nil
	Install(ctx context.Context, req dtoInstall.InitReq, progress chan<- installer.Step) (createdDB bool, err error)
	// AfterInstalled 安装成功后善后：保存 config、切换 global.DB、重置 casbin、触发 OnReady 回调
	AfterInstalled(newDB config.DB) error
	// OnReady 注册"安装完成 / DB 切换后"的回调；由 initialize 注册 InitService
	OnReady(fn func())
}

// NewService 构造 install.Service
func NewService() Service { return &service{} }

type service struct {
	readyHooks []func()
}

// Default 包级单例
var Default Service

// 预定义错误
var (
	ErrAlreadyInstalled = errors.New("already installed")
)

func (s *service) Status() installer.DBStatus {
	return installer.Detect(global.Cfg.DB)
}

func (s *service) CheckDB(req dtoInstall.CheckDBReq) dtoInstall.CheckDBResp {
	resp := dtoInstall.CheckDBResp{
		Configured: req.Configured(),
		Driver:     req.Driver,
	}
	if req.CreateIfMissing {
		created, err := installer.EnsureDatabase(req.DB)
		if err != nil {
			resp.Reason = "ensure db: " + err.Error()
			return resp
		}
		resp.Created = created
	}
	st := installer.Detect(req.DB)
	resp.Configured = st.Configured
	resp.DBConnected = st.Connected
	resp.Installed = st.Installed
	resp.Version = st.Version
	resp.Driver = st.Driver
	if st.Reason != "" {
		resp.Reason = st.Reason
	}
	return resp
}

// Install 真正跑安装流程；不负责善后（由 AfterInstalled 完成）
func (s *service) Install(ctx context.Context, req dtoInstall.InitReq, progress chan<- installer.Step) (bool, error) {
	if global.Installed.Load() {
		return false, ErrAlreadyInstalled
	}
	createdDB, err := installer.EnsureDatabase(req.DB)
	if err != nil {
		return false, err
	}
	db, err := installer.OpenWith(req.DB)
	if err != nil {
		return createdDB, err
	}
	defer closeGormDB(db)

	// 调用方传入的 ctx 可能已有 timeout；这里再叠加一个上限兜底
	subCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	if err := installer.Install(subCtx, db, req.Admin, progress); err != nil {
		return createdDB, err
	}
	return createdDB, nil
}

// AfterInstalled 完成 config 落盘 + 切换 global.DB + 重置 casbin + 触发 onReady
func (s *service) AfterInstalled(newDB config.DB) error {
	global.Cfg.DB = newDB
	if err := config.Save(global.Cfg); err != nil {
		return errors.New("config save: " + err.Error())
	}
	if err := refreshGlobalDB(); err != nil {
		return errors.New("refresh db: " + err.Error())
	}
	if err := casbinUtil.Setup(global.DB); err != nil {
		return errors.New("casbin setup: " + err.Error())
	}
	global.Installed.Store(true)
	log.Info("install completed")
	// 通知所有等待 DB 就绪的回调（重新装配 service / repo 等）
	for _, fn := range s.readyHooks {
		fn()
	}
	return nil
}

func (s *service) OnReady(fn func()) {
	if fn == nil {
		return
	}
	s.readyHooks = append(s.readyHooks, fn)
}

// refreshGlobalDB 用最新 cfg 重新打开连接并赋值给 global.DB
func refreshGlobalDB() error {
	if global.DB != nil {
		if sqlDB, err := global.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	db, err := installer.OpenWith(global.Cfg.DB)
	if err != nil {
		return err
	}
	global.SetDB(db)
	return nil
}

// closeGormDB 安全关闭 GORM 背后的 sql.DB 连接
func closeGormDB(db *gorm.DB) {
	if db == nil {
		return
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
