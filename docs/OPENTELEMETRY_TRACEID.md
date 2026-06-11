# 使用 OpenTelemetry TraceID

## 概述

现在我们优先使用 OpenTelemetry 的 traceID，符合云原生可观测性标准。

## 为什么优先使用 OTel TraceID？

1. **标准化格式**：符合 W3C Trace Context 标准
2. **分布式追踪**：可以跨服务追踪整个请求链路
3. **生态系统集成**：与 Jaeger、Zipkin、Grafana Tempo 等兼容
4. **多语言支持**：OTel 在各种语言中都有实现
5. **自动传播**：通过 HTTP header 自动传播 trace 信息

## 工作流程

### 1. Trace 传播链

```
上游服务 → HTTP Header → otelgin 中间件 → Trace 中间件 → 业务代码 → 日志
   (W3C Trace Context)  (创建 span)     (读取 span)    (使用 WithContext)
```

### 2. 中间件顺序

```
Recovery() 
  → otelgin.Middleware()  [创建 span，设置 span context]
  → middleware.Trace()    [读取 span，设置 logger]
  → Cors()
  → RequestLog()
  → 业务处理
```

## Context 中的字段

当使用 OpenTelemetry 时，`WithContext(ctx)` 会自动提取以下字段：

| 字段 | 来源 | 说明 |
|------|------|------|
| `trace_id` | `trace.SpanContextFromContext(ctx)` | 标准 16 字节 traceID |
| `span_id` | `trace.SpanContextFromContext(ctx)` | 标准 8 字节 spanID |
| `trace_flags` | `trace.SpanContextFromContext(ctx)` | 采样等标志 |
| `trace_state` | `trace.SpanContextFromContext(ctx)` | 供应商特定的 trace state（可选） |
| `request_id` | `context.Value("request_id")` | 应用级请求 ID（总是存在） |

## 回退机制

如果没有 OpenTelemetry span，系统会优雅回退：

1. **优先使用 OTel span context**（标准做法）
2. **如果没有 OTel，使用自定义 X-Trace-Id header**
3. **如果都没有，生成新的 UUID**

## 使用示例

### 基本使用（自动）

```go
// 在 handler 中 - 什么都不用改！
func Handler(c *gin.Context) {
    // 自动获取带有 OTel trace_id 的 logger
    logger := middleware.LoggerFromContext(c.Request.Context())
    logger.Info("处理请求", "user_id", 123)
    // 日志自动包含 trace_id, span_id, request_id
}
```

### 手动创建 Span

```go
import "go.opentelemetry.io/otel"

func SomeService(ctx context.Context) {
    // 创建一个 span
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "operation-name")
    defer span.End()

    // WithContext 自动提取 span 的 trace_id
    logger := baseLogger.WithContext(ctx)
    logger.Info("在 span 中记录日志")
}
```

## HTTP Header 格式

### 传入请求支持的 Header

| Header | 说明 | 示例 |
|--------|------|------|
| `traceparent` | W3C Trace Context 标准 | `00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01` |
| `tracestate` | 供应商特定的 state（可选） | `vendor1=value1,vendor2=value2` |
| `X-Trace-Id` | 自定义 trace_id（回退使用） | `abc123-def456` |

### 响应包含的 Header

| Header | 说明 |
|--------|------|
| `X-Trace-Id` | 本次请求使用的 trace_id（可用于关联日志） |
| `X-Request-Id` | 本次请求的 request_id |

## 日志输出示例

### 有 OTel Trace 时

```json
{
  "level": "info",
  "ts": "2024-06-11T10:30:00.000Z",
  "msg": "用户登录",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "trace_flags": "01",
  "request_id": "789abc-def012",
  "user_id": 123
}
```

### 没有 OTel 时（回退）

```json
{
  "level": "info",
  "ts": "2024-06-11T10:30:00.000Z",
  "msg": "用户登录",
  "trace_id": "abc123-def456-ghi789",
  "request_id": "789abc-def012",
  "user_id": 123
}
```

## 配置

在 `config/config.yaml` 中启用 OpenTelemetry：

```yaml
observability:
  trace:
    enable: true
    exporter: stdout  # 或 otlp（未来支持）
    service_name: go-admin-server
    sample_rate: 1.0
```

## 最佳实践

1. **让 otelgin 先执行**：中间件顺序很重要
2. **在整个调用链中传递 ctx**：确保 trace 传播
3. **使用 `log.FromContext(ctx)`**：而不是全局 logger
4. **创建有意义的 span**：在重要操作处创建 span
5. **利用 trace_id 关联日志**：在排查问题时通过 trace_id 过滤所有相关日志
