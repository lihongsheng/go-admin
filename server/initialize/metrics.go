package initialize

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/lihongsheng/go-admin/server/config"
	"github.com/lihongsheng/go-admin/server/log"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// metricsState 保存 Metrics 状态，用于 shutdown
type metricsState struct {
	meterProvider *metric.MeterProvider
	hasPrometheus bool
}

var globalMetricsState *metricsState

// InitMetricsOption Metrics 初始化选项
type InitMetricsOption func(*metricsInitConfig)

type metricsInitConfig struct {
	logger      log.Logger
	serviceName string
	appConfig   config.App
}

// WithMetricsLogger 为 Metrics 设置 logger
func WithMetricsLogger(logger log.Logger) InitMetricsOption {
	return func(c *metricsInitConfig) { c.logger = logger }
}

// WithMetricsServiceName 设置服务名
func WithMetricsServiceName(name string) InitMetricsOption {
	return func(c *metricsInitConfig) { c.serviceName = name }
}

// WithMetricsAppConfig 设置 app 配置
func WithMetricsAppConfig(cfg config.App) InitMetricsOption {
	return func(c *metricsInitConfig) { c.appConfig = cfg }
}

// InitMetrics 初始化 Metrics。
//
// 支持两种导出模式：
//   - Prometheus（拉取）：始终在 /metrics 端点暴露，由 Prometheus Server 定期抓取
//   - OTLP（推送）：当 cfg.Endpoint 非空时，额外推送指标到 OTLP collector
//
// 两种模式可同时启用。当 cfg.Enable 为 false 时全部跳过。
func InitMetrics(cfg config.Metrics, opts ...InitMetricsOption) error {
	initCfg := &metricsInitConfig{
		logger: log.Nop(),
	}
	for _, opt := range opts {
		opt(initCfg)
	}

	if !cfg.Enable {
		initCfg.logger.Info("observability: metrics disabled")
		return nil
	}

	// 确定服务名（与 otel.go 保持一致的降级逻辑）
	serviceName := initCfg.serviceName
	if serviceName == "" {
		serviceName = initCfg.appConfig.Name
	}
	if serviceName == "" {
		serviceName = "go-admin-server"
	}

	// 创建 resource
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String("1.0.0"),
	)

	// 收集 readers
	var readers []metric.Reader
	state := &metricsState{}

	// ── Prometheus reader ──
	promExporter, err := prometheus.New()
	if err != nil {
		return fmt.Errorf("create prometheus exporter: %w", err)
	}
	readers = append(readers, promExporter)
	state.hasPrometheus = true

	// ── OTLP 推送通道（可选）──
	if cfg.Endpoint != "" {
		otlpReader, err := newOTLPMetricReader(cfg.Endpoint)
		if err != nil {
			return fmt.Errorf("create otlp metric exporter: %w", err)
		}
		readers = append(readers, otlpReader)
	}

	// 创建 MeterProvider（多 reader）
	mpOpts := []metric.Option{metric.WithResource(res)}
	for _, r := range readers {
		mpOpts = append(mpOpts, metric.WithReader(r))
	}
	mp := metric.NewMeterProvider(mpOpts...)

	// 设置全局 MeterProvider
	otel.SetMeterProvider(mp)

	state.meterProvider = mp
	globalMetricsState = state

	initCfg.logger.Info("observability: metrics initialized",
		"exporter", cfg.Exporter,
		"endpoint", cfg.Endpoint,
		"path", cfg.Path,
		"service_name", serviceName,
	)

	return nil
}

// newOTLPMetricReader 根据 endpoint 创建 OTLP metric exporter + PeriodicReader。
// http(s):// 前缀 → HTTP/protobuf，否则 → gRPC。
func newOTLPMetricReader(endpoint string) (metric.Reader, error) {
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		cleanEndpoint := strings.TrimPrefix(endpoint, "http://")
		cleanEndpoint = strings.TrimPrefix(cleanEndpoint, "https://")
		exp, err := otlpmetrichttp.New(context.Background(),
			otlpmetrichttp.WithEndpoint(cleanEndpoint),
			otlpmetrichttp.WithInsecure(),
		)
		if err != nil {
			return nil, err
		}
		return metric.NewPeriodicReader(exp), nil
	}

	exp, err := otlpmetricgrpc.New(context.Background(),
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return metric.NewPeriodicReader(exp), nil
}

// ShutdownMetrics 优雅关闭 Metrics
func ShutdownMetrics(ctx context.Context) error {
	if globalMetricsState != nil && globalMetricsState.meterProvider != nil {
		return globalMetricsState.meterProvider.Shutdown(ctx)
	}
	return nil
}

// MetricsHandler 返回 Prometheus HTTP handler（promhttp.Handler）。
// 当 Prometheus 未启用时返回 nil。
func MetricsHandler() http.Handler {
	if globalMetricsState != nil && globalMetricsState.hasPrometheus {
		return promhttp.Handler()
	}
	return nil
}
