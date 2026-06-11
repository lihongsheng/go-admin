// Package log 定义日志接口和实现，参考 go-kratos 设计
package log

import (
	"context"

	"go.opentelemetry.io/otel/trace"
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

// ── 全局 logger（临时兼容旧代码）────

var globalLogger Logger

// SetGlobal 设置全局 logger
func SetGlobal(logger Logger) {
	if logger != nil {
		globalLogger = logger
	}
}

// Global 返回全局 logger
func Global() Logger {
	if globalLogger != nil {
		return globalLogger
	}
	return Nop()
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
