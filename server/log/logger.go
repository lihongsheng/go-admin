// Package log 定义日志接口和实现，参考 go-kratos 设计
package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Valuer 是一个从 context 中提取值的函数类型
type Valuer func(ctx context.Context) interface{}

// Logger 定义日志接口
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})

	With(keysAndValues ...interface{}) Logger
	WithValuer(keys ...Valuer) Logger
	WithContext(ctx context.Context) Logger

	Sync() error
}

// 内置 Valuer 函数

const (
	TraceIDKey    = "trace_id"
	SpanIDKey     = "span_id"
	TraceFlagsKey = "trace_flags"
	RequestIDKey  = "request_id"
)

// TraceID 返回一个 Valuer，从 context 中提取 OTel trace_id
// 优先使用 OTel span context；没有时回退到手动设置的 trace_id
func TraceID() Valuer {
	return func(ctx context.Context) interface{} {
		if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
			return sc.TraceID().String()
		}
		// 回退：手动设置的 trace_id
		if id, ok := ctx.Value(traceIDKey{}).(string); ok {
			return id
		}
		return nil
	}
}

// SpanID 返回一个 Valuer，从 context 中提取 OTel span_id
func SpanID() Valuer {
	return func(ctx context.Context) interface{} {
		if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
			return sc.SpanID().String()
		}
		return nil
	}
}

// TraceFlags 返回一个 Valuer，从 context 中提取 OTel trace_flags
func TraceFlags() Valuer {
	return func(ctx context.Context) interface{} {
		if sc := trace.SpanContextFromContext(ctx); sc.IsValid() {
			return sc.TraceFlags().String()
		}
		return nil
	}
}

// RequestID 返回一个 Valuer，从 context 中提取 request_id
func RequestID() Valuer {
	return func(ctx context.Context) interface{} {
		if id, ok := ctx.Value(requestIDKey{}).(string); ok {
			return id
		}
		return nil
	}
}

// DefaultValuers 返回默认的 Valuer 列表
func DefaultValuers() []Valuer {
	return []Valuer{
		TraceID(),
		SpanID(),
		TraceFlags(),
		RequestID(),
	}
}

// requestIDKey 用于存储 request_id 的 context key
type requestIDKey struct{}

// traceIDKey 用于存储手动设置的 trace_id（回退方案）
type traceIDKey struct{}

// WithRequestID 将 request_id 设置到 context
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// WithTraceID 将自定义 trace_id 设置到 context（当没有 OTel span 时回退使用）
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// Extract 从 context 中提取键值对
func Extract(ctx context.Context, valuers []Valuer) []interface{} {
	if len(valuers) == 0 {
		return nil
	}
	var result []interface{}

	// 使用默认的 key 名（假设 valuers 使用 DefaultValuers 的顺序）
	for i, v := range valuers {
		value := v(ctx)
		if value == nil {
			continue
		}

		var key string
		switch i {
		case 0:
			key = TraceIDKey
		case 1:
			key = SpanIDKey
		case 2:
			key = TraceFlagsKey
		case 3:
			key = RequestIDKey
		default:
			key = "key"
		}
		result = append(result, key, value)
	}

	return result
}

// Nop Logger

type nopLogger struct{}

func (*nopLogger) Debug(msg string, keysAndValues ...interface{}) {}
func (*nopLogger) Info(msg string, keysAndValues ...interface{})  {}
func (*nopLogger) Warn(msg string, keysAndValues ...interface{})  {}
func (*nopLogger) Error(msg string, keysAndValues ...interface{}) {}
func (*nopLogger) Fatal(msg string, keysAndValues ...interface{}) {}
func (n *nopLogger) With(keysAndValues ...interface{}) Logger { return n }
func (n *nopLogger) WithValuer(keys ...Valuer) Logger          { return n }
func (n *nopLogger) WithContext(ctx context.Context) Logger    { return n }
func (*nopLogger) Sync() error                                 { return nil }

// Nop 返回一个不做任何操作的 Logger
func Nop() Logger {
	return &nopLogger{}
}

// ── 全局 logger ────

var (
	globalLogger Logger
	// globalZap 持有底层 *zap.Logger，跳过 zapLogger 和便捷函数两层包装，
	// 保证 log.Info() 等便捷函数也能输出正确的 caller 位置。
	globalZap *zap.Logger
)

// SetGlobal 设置全局 logger
func SetGlobal(logger Logger) {
	if logger != nil {
		globalLogger = logger
		// 提取底层 zap logger，为便捷函数增加一跳 skip（跳过便捷函数自身）
		if zl, ok := logger.(*zapLogger); ok {
			globalZap = zl.logger.WithOptions(zap.AddCallerSkip(1))
		} else {
			globalZap = nil
		}
	}
}

// Global 返回全局 logger
func Global() Logger {
	if globalLogger != nil {
		return globalLogger
	}
	return Nop()
}

// ── 包级便捷函数 ────
// 通过 globalZap 直接写日志，绕过 zapLogger 包装层，确保 caller 信息指向真实调用方。
// 若 globalZap 不可用（logger 不是 zap 实现），回退到 Global()。

// Debug 输出调试级别日志
func Debug(msg string, keysAndValues ...interface{}) {
	if globalZap != nil {
		globalZap.Debug(msg, ToZapFields(keysAndValues...)...)
	} else {
		Global().Debug(msg, keysAndValues...)
	}
}

// Info 输出信息级别日志
func Info(msg string, keysAndValues ...interface{}) {
	if globalZap != nil {
		globalZap.Info(msg, ToZapFields(keysAndValues...)...)
	} else {
		Global().Info(msg, keysAndValues...)
	}
}

// Warn 输出警告级别日志
func Warn(msg string, keysAndValues ...interface{}) {
	if globalZap != nil {
		globalZap.Warn(msg, ToZapFields(keysAndValues...)...)
	} else {
		Global().Warn(msg, keysAndValues...)
	}
}

// Error 输出错误级别日志
func Error(msg string, keysAndValues ...interface{}) {
	if globalZap != nil {
		globalZap.Error(msg, ToZapFields(keysAndValues...)...)
	} else {
		Global().Error(msg, keysAndValues...)
	}
}

// Fatal 输出致命级别日志
func Fatal(msg string, keysAndValues ...interface{}) {
	if globalZap != nil {
		globalZap.Fatal(msg, ToZapFields(keysAndValues...)...)
	} else {
		Global().Fatal(msg, keysAndValues...)
	}
}

// ── Context 相关 ──

type loggerContextKey struct{}

// NewContext 将 logger 存入 context
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

// FromContext 从 context 中获取 logger
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return logger
	}
	return Nop()
}
