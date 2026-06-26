package initialize

import (
	"os"

	"github.com/lihongsheng/go-admin/server/config"
	"github.com/lihongsheng/go-admin/server/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 初始化 zap logger，返回 log.Logger 接口
func Logger(cfg config.Log) (log.Logger, error) {
	if cfg.Dir == "" {
		cfg.Dir = "./logs"
	}
	if err := os.MkdirAll(cfg.Dir, 0o755); err != nil {
		return nil, err
	}

	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level))

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	fileWriter := log.NewDateWriter(cfg.Dir, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(fileWriter), level),
		zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(os.Stdout), level),
	)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return log.NewZapLoggerWithDefaults(zapLogger), nil
}
