package log_test

import (
	"context"
	"os"

	"github.com/lihongsheng/go-admin/server/log"
	"github.com/lihongsheng/go-admin/server/middleware"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Example_logger_WithContext 展示如何使用 WithContext 自动添加 trace_id
func Example_logger_WithContext() {
	// 1. 创建一个基于 zap 的 logger（只输出到控制台，方便演示）
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	})
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	zapLogger := zap.New(core)
	logger := log.NewZapLogger(zapLogger)

	// 2. 创建一个带有 trace_id 的 context
	ctx := context.Background()
	ctx = middleware.SetTraceIDToContext(ctx, "test-trace-id-12345")
	ctx = middleware.SetRequestIDToContext(ctx, "test-request-id-67890")

	// 3. 使用 WithContext 获取带有 trace 信息的 logger
	tracedLogger := logger.WithContext(ctx)

	// 4. 打印日志，会自动包含 trace_id 和 request_id
	tracedLogger.Info("这是一条带有 trace_id 的日志",
		"key1", "value1",
		"key2", "value2")

	// Output:
	// 这是一条带有 trace_id 的日志，会自动包含 trace_id 和 request_id 字段
}

// Example_usage_in_service 展示在 service 层如何使用
func Example_usage_in_service() {
	// 假设 logger 已经通过依赖注入传入
	var logger log.Logger

	// 处理请求时
	ctx := context.Background()
	// trace_id 已经由中间件设置好了

	// 在 service 中获取带有 trace 的 logger
	tracedLogger := logger.WithContext(ctx)
	tracedLogger.Info("处理业务逻辑", "user_id", 123)

	// 结构化日志
	tracedLogger.Info("用户操作完成", "user_id", 123, "action", "login")
}

// Example_with_structured_fields 展示结构化日志使用
func Example_with_structured_fields() {
	var logger log.Logger
	ctx := context.Background()

	// 添加多个字段
	userLogger := logger.WithContext(ctx).With(
		"user_id", 123,
		"role", "admin",
		"tenant_id", "tenant-abc",
	)

	// 之后使用时会自动包含所有这些字段
	userLogger.Info("用户登录")
	userLogger.Warn("用户尝试敏感操作")
}
