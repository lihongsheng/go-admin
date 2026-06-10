package installer

import (
	"errors"
	"fmt"
	"strings"

	"go-admin/server/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// EnsureDatabase 若目标库不存在则创建之。
// 仅对 mysql / postgres 生效；sqlite 文件由驱动自动创建。
// 返回是否实际新建了库。
func EnsureDatabase(c config.DB) (created bool, err error) {
	switch c.Driver {
	case "mysql":
		return ensureMySQL(c)
	case "postgres":
		return ensurePostgres(c)
	case "sqlite":
		return false, nil
	}
	return false, fmt.Errorf("unsupported driver: %s", c.Driver)
}

// 检查 mysql 错误是否表示「库不存在」
func isMySQLUnknownDB(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	// MySQL: Error 1049 (42000): Unknown database 'xxx'
	return strings.Contains(msg, "1049") || strings.Contains(strings.ToLower(msg), "unknown database")
}

// 检查 postgres 错误是否表示「库不存在」
func isPgUnknownDB(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	// pq/pgx: SQLSTATE 3D000 "database \"xxx\" does not exist"
	return strings.Contains(msg, "3d000") || strings.Contains(msg, "does not exist")
}

func ensureMySQL(c config.DB) (bool, error) {
	if c.Database == "" {
		return false, errors.New("mysql: database name is empty")
	}
	// 1) 先尝试连目标库；能通就什么都不做
	if db, err := gorm.Open(mysql.Open(c.DSN()), silent()); err == nil {
		closeGorm(db)
		return false, nil
	} else if !isMySQLUnknownDB(err) {
		return false, fmt.Errorf("connect target db: %w", err)
	}

	// 2) 连默认实例 → CREATE DATABASE
	srv, err := gorm.Open(mysql.Open(c.ServerDSN()), silent())
	if err != nil {
		return false, fmt.Errorf("connect server: %w", err)
	}
	defer closeGorm(srv)

	charset := c.Charset
	if charset == "" {
		charset = "utf8mb4"
	}
	// 用反引号避免库名特殊字符；charset 是白名单写死，注入风险可控
	stmt := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET %s COLLATE %s_unicode_ci",
		strings.ReplaceAll(c.Database, "`", ""),
		charset, charset)
	if err := srv.Exec(stmt).Error; err != nil {
		return false, fmt.Errorf("create database: %w", err)
	}
	return true, nil
}

func ensurePostgres(c config.DB) (bool, error) {
	if c.Database == "" {
		return false, errors.New("postgres: database name is empty")
	}
	if db, err := gorm.Open(postgres.Open(c.DSN()), silent()); err == nil {
		closeGorm(db)
		return false, nil
	} else if !isPgUnknownDB(err) {
		return false, fmt.Errorf("connect target db: %w", err)
	}

	srv, err := gorm.Open(postgres.Open(c.ServerDSN()), silent())
	if err != nil {
		return false, fmt.Errorf("connect postgres system db: %w", err)
	}
	defer closeGorm(srv)

	// Postgres 不支持 IF NOT EXISTS for CREATE DATABASE 在所有版本通用，
	// 这里先查询 pg_database 再决定是否创建
	var exists bool
	if err := srv.Raw(`SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)`, c.Database).Scan(&exists).Error; err != nil {
		return false, fmt.Errorf("query pg_database: %w", err)
	}
	if exists {
		return false, nil
	}
	stmt := fmt.Sprintf(`CREATE DATABASE "%s" ENCODING 'UTF8'`,
		strings.ReplaceAll(c.Database, `"`, ""))
	if err := srv.Exec(stmt).Error; err != nil {
		return false, fmt.Errorf("create database: %w", err)
	}
	return true, nil
}

func silent() *gorm.Config { return &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)} }

func closeGorm(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
