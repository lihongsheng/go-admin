package initialize

import (
	"os"
	"path/filepath"

	"go-admin/server/config"
	"go-admin/server/global"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 初始化 zap + lumberjack
func Logger(c config.Log) {
	if c.Dir == "" {
		c.Dir = "./logs"
	}
	_ = os.MkdirAll(c.Dir, 0o755)

	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(c.Level))

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(c.Dir, "server.log"),
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   true,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(fileWriter), level),
		zapcore.NewCore(zapcore.NewConsoleEncoder(encCfg), zapcore.AddSync(os.Stdout), level),
	)
	global.Logger = zap.New(core, zap.AddCaller())
}
