package initialize

import (
	"fmt"

	"go-admin/server/core/installer"
	"go-admin/server/global"
	applog "go-admin/server/log"
	"go-admin/server/model/system"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormConnect 用 global.Cfg.DB 打开连接，挂到 global.DB
func GormConnect() error {
	db, err := installer.OpenWith(global.Cfg.DB)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	// 调整日志级别
	switch global.Cfg.DB.LogMode {
	case "silent":
		db.Logger = db.Logger.LogMode(logger.Silent)
	case "error":
		db.Logger = db.Logger.LogMode(logger.Error)
	case "warn":
		db.Logger = db.Logger.LogMode(logger.Warn)
	case "info":
		db.Logger = db.Logger.LogMode(logger.Info)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	if global.Cfg.DB.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(global.Cfg.DB.MaxIdle)
	}
	if global.Cfg.DB.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(global.Cfg.DB.MaxOpen)
	}
	global.SetDB(db)
	return nil
}

// DetectInstalled 在已连接 DB 上探测是否安装过
func DetectInstalled() {
	if global.DB == nil {
		return
	}
	if !global.DB.Migrator().HasTable(&system.SysInstall{}) {
		global.Installed.Store(false)
		return
	}
	var n int64
	if err := global.DB.Model(&system.SysInstall{}).Count(&n).Error; err != nil {
		applog.Warn("count sys_install: " + err.Error())
		return
	}
	global.Installed.Store(n > 0)
}

// MustRefreshDB 安装完成后用最新 cfg 重新打开连接
func MustRefreshDB() error {
	if global.DB != nil {
		if sqlDB, err := global.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
	return GormConnect()
}

// Tx 简写：在 global.DB 上开事务
func Tx(fn func(*gorm.DB) error) error { return global.DB.Transaction(fn) }
