package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

// zapLogger 基于 zap 的 Logger 实现
type zapLogger struct {
	logger  *zap.Logger
	valuers []Valuer
}

// NewZapLogger 创建一个基于 zap 的 Logger
func NewZapLogger(logger *zap.Logger) Logger {
	if logger == nil {
		return Nop()
	}
	return &zapLogger{
		logger: logger,
	}
}

func (z *zapLogger) Debug(msg string, keysAndValues ...interface{}) {
	z.logger.Debug(msg, z.toZapFields(keysAndValues...)...)
}

func (z *zapLogger) Info(msg string, keysAndValues ...interface{}) {
	z.logger.Info(msg, z.toZapFields(keysAndValues...)...)
}

func (z *zapLogger) Warn(msg string, keysAndValues ...interface{}) {
	z.logger.Warn(msg, z.toZapFields(keysAndValues...)...)
}

func (z *zapLogger) Error(msg string, keysAndValues ...interface{}) {
	z.logger.Error(msg, z.toZapFields(keysAndValues...)...)
}

func (z *zapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	z.logger.Fatal(msg, z.toZapFields(keysAndValues...)...)
}

func (z *zapLogger) With(keysAndValues ...interface{}) Logger {
	newLogger := *z
	newLogger.logger = z.logger.With(z.toZapFields(keysAndValues...)...)
	return &newLogger
}

func (z *zapLogger) WithValuer(valuers ...Valuer) Logger {
	newLogger := *z
	newLogger.valuers = make([]Valuer, 0, len(z.valuers)+len(valuers))
	newLogger.valuers = append(newLogger.valuers, z.valuers...)
	newLogger.valuers = append(newLogger.valuers, valuers...)
	return &newLogger
}

func (z *zapLogger) WithContext(ctx context.Context) Logger {
	if len(z.valuers) == 0 {
		return z
	}

	fields := Extract(ctx, z.valuers)
	if len(fields) == 0 {
		return z
	}
	return z.With(fields...)
}

func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}

// toZapFields 将键值对转换为 zap.Field
func (z *zapLogger) toZapFields(keysAndValues ...interface{}) []zap.Field {
	if len(keysAndValues) == 0 {
		return nil
	}
	if len(keysAndValues)%2 != 0 {
		keysAndValues = append(keysAndValues, "(MISSING)")
	}

	fields := make([]zap.Field, 0, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		key, ok := keysAndValues[i].(string)
		if !ok {
			key = fmt.Sprintf("%v", keysAndValues[i])
		}
		fields = append(fields, zap.Any(key, keysAndValues[i+1]))
	}
	return fields
}

// Helper 函数

// NewZapLoggerWithValuers 创建一个预配置 valuers 的 zap logger
func NewZapLoggerWithValuers(logger *zap.Logger, valuers ...Valuer) Logger {
	l := NewZapLogger(logger)
	if len(valuers) > 0 {
		l = l.WithValuer(valuers...)
	}
	return l
}

// NewZapLoggerWithDefaults 创建一个预配置默认 valuers 的 zap logger
func NewZapLoggerWithDefaults(logger *zap.Logger) Logger {
	return NewZapLoggerWithValuers(logger, DefaultValuers()...)
}

// New 创建一个默认的 zap logger（方便迁移）
func New(logger *zap.Logger) Logger {
	return NewZapLoggerWithDefaults(logger)
}
