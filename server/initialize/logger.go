package initialize

import (
	"os"
	"path/filepath"

	"go-admin/server/config"
	"go-admin/server/log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 初始化 zap logger，返回 log.Logger 接口
// 使用默认的 Valuers（trace_id, span_id, trace_flags, request_id）
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

	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Dir, "server.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   true,
	}

	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(fileWriter), level),
		zapcore.NewCore(zapcore.NewJSONEncoder(encCfg), zapcore.AddSync(os.Stdout), level),
	)
	zapLogger := zap.New(core, zap.AddCaller())

	// 创建带有默认 Valuers 的 logger
	return log.NewZapLoggerWithDefaults(zapLogger), nil
}

// ── 包级 logger（用于工具函数无法显式传 logger 的场景）────

var internalLogger log.Logger

// SetLogger 设置包级 logger
func SetLogger(logger log.Logger) {
	if logger != nil {
		internalLogger = logger
	}
}

// getLogger 返回包级 logger；未设置时返回 Nop
func getLogger() log.Logger {
	if internalLogger != nil {
		return internalLogger
	}
	return log.Nop()
}
