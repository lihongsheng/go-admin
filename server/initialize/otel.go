package initialize

import (
	"context"
	"fmt"
	"strings"

	"go-admin/server/config"
	"go-admin/server/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer 全局 tracer 实例（临时保留，后续可以移除）
var Tracer trace.Tracer

// otelState 保存 OpenTelemetry 状态
type otelState struct {
	tracerProvider *sdktrace.TracerProvider
}

var globalOtelState *otelState

// InitOpenTelemetryOption OpenTelemetry 初始化选项
type InitOpenTelemetryOption func(*otelInitConfig)

type otelInitConfig struct {
	logger      log.Logger
	serviceName string
	appConfig   config.App
}

// WithOTelLogger 为 OpenTelemetry 设置 logger
func WithOTelLogger(logger log.Logger) InitOpenTelemetryOption {
	return func(c *otelInitConfig) {
		c.logger = logger
	}
}

// WithOTelServiceName 设置服务名
func WithOTelServiceName(name string) InitOpenTelemetryOption {
	return func(c *otelInitConfig) {
		c.serviceName = name
	}
}

// WithOTelAppConfig 设置 app 配置
func WithOTelAppConfig(cfg config.App) InitOpenTelemetryOption {
	return func(c *otelInitConfig) {
		c.appConfig = cfg
	}
}

// InitOpenTelemetry 初始化 OpenTelemetry
func InitOpenTelemetry(cfg config.Observability, opts ...InitOpenTelemetryOption) error {
	initCfg := &otelInitConfig{
		logger: log.Nop(),
	}
	for _, opt := range opts {
		opt(initCfg)
	}

	if !cfg.Trace.Enable {
		initCfg.logger.Info("observability: trace disabled")
		return nil
	}

	// 设置服务名
	serviceName := cfg.Trace.ServiceName
	if serviceName == "" {
		serviceName = initCfg.serviceName
	}
	if serviceName == "" {
		serviceName = initCfg.appConfig.Name
	}
	if serviceName == "" {
		serviceName = "go-admin-server"
	}

	// 创建 exporter
	var exporter sdktrace.SpanExporter
	var err error

	switch cfg.Trace.Exporter {
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	case "otlp":
		endpoint := cfg.Trace.Endpoint
		if endpoint == "" {
			endpoint = "localhost:4317"
		}

		// 根据前缀判断协议：http(s):// → HTTP/protobuf，否则 gRPC
		if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
			cleanEndpoint := strings.TrimPrefix(endpoint, "http://")
			cleanEndpoint = strings.TrimPrefix(cleanEndpoint, "https://")
			exp, httpErr := otlptracehttp.New(context.Background(),
				otlptracehttp.WithEndpoint(cleanEndpoint),
				otlptracehttp.WithInsecure(),
			)
			exporter, err = exp, httpErr
		} else {
			exp, grpcErr := otlptracegrpc.New(context.Background(),
				otlptracegrpc.WithEndpoint(endpoint),
				otlptracegrpc.WithInsecure(),
			)
			exporter, err = exp, grpcErr
		}
	default:
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	}

	if err != nil {
		return fmt.Errorf("create trace exporter: %w", err)
	}

	// 设置采样率
	sampler := sdktrace.TraceIDRatioBased(cfg.Trace.SampleRate)
	if cfg.Trace.SampleRate <= 0 {
		sampler = sdktrace.NeverSample()
	} else if cfg.Trace.SampleRate >= 1 {
		sampler = sdktrace.AlwaysSample()
	}

	// 创建 resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("1.0.0"),
	)

	// 创建 tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	// 设置全局 tracer provider
	otel.SetTracerProvider(tp)

	// 创建全局 tracer
	Tracer = tp.Tracer(serviceName)

	// 保存状态用于 shutdown
	globalOtelState = &otelState{
		tracerProvider: tp,
	}

	initCfg.logger.Info("observability: trace initialized",
		"exporter", cfg.Trace.Exporter,
		"service_name", serviceName,
		"sample_rate", cfg.Trace.SampleRate,
	)

	return nil
}

// ShutdownOpenTelemetry 优雅关闭 OpenTelemetry
func ShutdownOpenTelemetry(ctx context.Context) error {
	if globalOtelState != nil && globalOtelState.tracerProvider != nil {
		return globalOtelState.tracerProvider.Shutdown(ctx)
	}
	if tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		return tp.Shutdown(ctx)
	}
	return nil
}
