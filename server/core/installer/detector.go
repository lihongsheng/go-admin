package installer

import (
	"errors"
	"fmt"

	"go-admin/server/config"
	"go-admin/server/model/system"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBStatus 安装向导前端展示用
type DBStatus struct {
	Configured bool   `json:"configured"`
	Connected  bool   `json:"db_connected"`
	Installed  bool   `json:"installed"`
	Version    string `json:"version,omitempty"`
	Driver     string `json:"driver,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// OpenWith 根据传入的 DB 配置打开连接（不挂到 global）
func OpenWith(c config.DB) (*gorm.DB, error) {
	if !c.Configured() {
		return nil, errors.New("db not configured")
	}
	dsn := c.DSN()
	switch c.Driver {
	case "mysql":
		return gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	case "sqlite":
		return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	}
	return nil, fmt.Errorf("unsupported driver: %s", c.Driver)
}

// Detect 探测 DB 连通性 + 是否已安装
func Detect(c config.DB) DBStatus {
	s := DBStatus{Configured: c.Configured(), Driver: c.Driver}
	if !s.Configured {
		s.Reason = "db not configured"
		return s
	}
	db, err := OpenWith(c)
	if err != nil {
		s.Reason = err.Error()
		return s
	}
	sqlDB, err := db.DB()
	if err != nil {
		s.Reason = err.Error()
		return s
	}
	if err := sqlDB.Ping(); err != nil {
		s.Reason = "ping failed: " + err.Error()
		return s
	}
	s.Connected = true
	if db.Migrator().HasTable(&system.SysInstall{}) {
		var rec system.SysInstall
		if err := db.First(&rec).Error; err == nil {
			s.Installed = true
			s.Version = rec.Version
		}
	}
	_ = sqlDB.Close()
	return s
}
