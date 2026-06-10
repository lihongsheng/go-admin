package installer

import (
	"errors"
	"sync"

	"go-admin/server/model/system"

	"gorm.io/gorm"
)

// 进程级安装锁（防止 /install/init 被并发触发）
var installMu sync.Mutex

// AcquireProcessLock 非阻塞获取进程锁
func AcquireProcessLock() error {
	if !installMu.TryLock() {
		return errors.New("install already in progress")
	}
	return nil
}

// ReleaseProcessLock 释放进程锁
func ReleaseProcessLock() { installMu.Unlock() }

// EnsureNotInstalled DB 层防重复安装：sys_install 表已有记录则拒绝
func EnsureNotInstalled(db *gorm.DB) error {
	if !db.Migrator().HasTable(&system.SysInstall{}) {
		return nil
	}
	var n int64
	if err := db.Model(&system.SysInstall{}).Count(&n).Error; err != nil {
		return err
	}
	if n > 0 {
		return errors.New("system already installed")
	}
	return nil
}
