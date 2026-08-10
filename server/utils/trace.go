package utils

import (
	"context"

	"go-admin/server/log"
	"go-admin/server/middleware"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan 开始一个 span（简化版）
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer("").Start(ctx, name, opts...)
}

// StartSpanWithAttrs 开始一个带属性的 span
func StartSpanWithAttrs(ctx context.Context, name string, attrs map[string]string) (context.Context, trace.Span) {
	attributes := make([]attribute.KeyValue, 0, len(attrs))
	for k, v := range attrs {
		attributes = append(attributes, attribute.String(k, v))
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attributes...),
	}
	return otel.Tracer("").Start(ctx, name, opts...)
}

// AddEventToSpan 给当前 span 添加事件
func AddEventToSpan(ctx context.Context, name string, attrs map[string]string) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		attributes := make([]attribute.KeyValue, 0, len(attrs))
		for k, v := range attrs {
			attributes = append(attributes, attribute.String(k, v))
		}
		span.AddEvent(name, trace.WithAttributes(attributes...))
	}
}

// SetSpanAttributes 设置当前 span 的属性
func SetSpanAttributes(ctx context.Context, attrs map[string]string) {
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		attributes := make([]attribute.KeyValue, 0, len(attrs))
		for k, v := range attrs {
			attributes = append(attributes, attribute.String(k, v))
		}
		span.SetAttributes(attributes...)
	}
}

// GetTraceID 获取 trace_id
func GetTraceID(ctx context.Context) string {
	return middleware.GetTraceID(ctx)
}

// GetRequestID 获取 request_id
func GetRequestID(ctx context.Context) string {
	return middleware.GetRequestID(ctx)
}

// Logger 从 context 获取带 trace 的 logger
func Logger(ctx context.Context) log.Logger {
	return middleware.LoggerFromContext(ctx)
}
