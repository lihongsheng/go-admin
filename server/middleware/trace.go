package middleware

import (
	"context"
	"time"

	"github.com/lihongsheng/go-admin/server/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// TraceOption Trace 中间件选项
type TraceOption func(*traceConfig)

type traceConfig struct {
	logger log.Logger
}

// WithLogger 为 Trace 中间件设置 logger
func WithLogger(logger log.Logger) TraceOption {
	return func(c *traceConfig) {
		c.logger = logger
	}
}

// Trace 追踪中间件
// 完全基于 OpenTelemetry 和 context.Context，不使用 gin.Context.Set
func Trace(opts ...TraceOption) gin.HandlerFunc {
	cfg := &traceConfig{
		logger: log.Nop(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// 1. 从 HTTP Header 提取 W3C Trace Context（标准方式）
		// 如果 header 中没有，otelgin 会自动创建新的 span
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		// 2. 生成 request_id（应用级的请求 ID）
		requestID := c.GetHeader("X-Request-Id")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		// 只设置到 context，不使用 gin.Set
		ctx = log.WithRequestID(ctx, requestID)

		// 3. 将 logger 存入 context（带有 valuers 的 logger）
		// WithContext 会自动从 ctx 中提取 OTel 的 trace_id
		ctxLogger := cfg.logger.WithContext(ctx)
		ctx = log.NewContext(ctx, ctxLogger)

		// 4. 更新 request context
		c.Request = c.Request.WithContext(ctx)

		// 5. 将 trace_id 和 request_id 写入响应 header（方便调试）
		if traceID := log.TraceID()(ctx); traceID != nil {
			if id, ok := traceID.(string); ok {
				c.Header("X-Trace-Id", id)
			}
		}
		c.Header("X-Request-Id", requestID)

		c.Next()
	}
}

// GetTraceID 从 context 提取 OTel trace_id
func GetTraceID(ctx context.Context) string {
	if id := log.TraceID()(ctx); id != nil {
		if s, ok := id.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// GetRequestID 从 context 提取 request_id
func GetRequestID(ctx context.Context) string {
	if id := log.RequestID()(ctx); id != nil {
		if s, ok := id.(string); ok && s != "" {
			return s
		}
	}
	return ""
}

// SetTraceIDToContext 将自定义 trace_id 存入 context（回退方案，优先使用 OTel）
func SetTraceIDToContext(ctx context.Context, traceID string) context.Context {
	return log.WithTraceID(ctx, traceID)
}

// SetRequestIDToContext 将自定义 request_id 存入 context
func SetRequestIDToContext(ctx context.Context, requestID string) context.Context {
	return log.WithRequestID(ctx, requestID)
}

// LoggerFromContext 从 context 获取 logger（没有则返回 Nop）
func LoggerFromContext(ctx context.Context) log.Logger {
	return log.FromContext(ctx)
}

// RequestLog 访问日志中间件
type RequestLogOption func(*requestLogConfig)

type requestLogConfig struct {
	logger log.Logger
}

func WithRequestLogger(logger log.Logger) RequestLogOption {
	return func(c *requestLogConfig) {
		c.logger = logger
	}
}

func RequestLog(opts ...RequestLogOption) gin.HandlerFunc {
	cfg := &requestLogConfig{
		logger: log.Nop(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		// 从 context 中获取 logger（已经带有 trace 信息）
		logger := log.FromContext(c.Request.Context())

		logger.Info("http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"cost", time.Since(start),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}
