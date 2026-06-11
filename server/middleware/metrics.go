package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// 包级 meter 和 instruments。
// 由 init() 一次性创建。如果 MeterProvider 未设置或为 noop，
// 所有 instruments 会自动降级为 noop — 中间件执行但无开销。
var (
	metricsMeter = otel.Meter("go-admin-server")

	httpRequestCount  metric.Int64Counter
	httpRequestDur    metric.Float64Histogram
	httpRequestErrors metric.Int64Counter
)

func init() {
	var err error

	httpRequestCount, err = metricsMeter.Int64Counter(
		"http.server.request_count",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		handleInstrumentErr(err)
	}

	httpRequestDur, err = metricsMeter.Float64Histogram(
		"http.server.request_duration_seconds",
		metric.WithDescription("Duration of HTTP requests in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		handleInstrumentErr(err)
	}

	httpRequestErrors, err = metricsMeter.Int64Counter(
		"http.server.request_errors",
		metric.WithDescription("Total number of HTTP requests that resulted in 5xx"),
		metric.WithUnit("1"),
	)
	if err != nil {
		handleInstrumentErr(err)
	}
}

func handleInstrumentErr(err error) {
	// OTel API 约定：即使创建失败也会返回 noop instrument，不会返回 error。
	// 此函数仅为防御性编程保留。
	_ = err
}

// Metrics 返回一个记录 HTTP 请求指标的 gin 中间件。
//
// 记录的指标：
//   - http.server.request_count            计数器
//   - http.server.request_duration_seconds  直方图
//   - http.server.request_errors           计数器（仅 5xx）
//
// 属性：http.method, http.route, http.scheme, http.status_code
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", c.FullPath()),
			attribute.String("http.scheme", c.Request.URL.Scheme),
			attribute.Int("http.status_code", c.Writer.Status()),
		}

		httpRequestCount.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))

		dur := time.Since(start).Seconds()
		httpRequestDur.Record(c.Request.Context(), dur, metric.WithAttributes(attrs...))

		if c.Writer.Status() >= 500 {
			httpRequestErrors.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
		}
	}
}
