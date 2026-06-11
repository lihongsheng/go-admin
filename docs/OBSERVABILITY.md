# 可观测性文档

本文档说明如何在 go-admin 项目中使用可观测性功能（Trace ID、分布式追踪等）。

## 功能概述

- **Trace ID**: 每个请求自动生成或传播 trace_id，符合 W3C Trace Context 规范
- **日志增强**: 所有日志自动包含 trace_id 和 request_id，便于问题追踪
- **OpenTelemetry 集成**: 支持分布式追踪，可导出到 Jaeger、Zipkin 等
- **标准 Header**: 支持 `X-Trace-Id` 和 `X-Request-Id` 自定义 header

## 配置说明

在 `config/config.yaml` 中配置可观测性：

```yaml
observability:
    trace:
        enable: true              # 是否启用分布式追踪
        exporter: stdout          # 导出器类型: stdout, otlp (暂未实现)
        service_name: go-admin-server  # 服务名称
        sample_rate: 1.0          # 采样率 (0.0-1.0)
    metrics:
        enable: false             # 指标暂未实现
        path: /metrics
```

## 使用方式

### 1. 在 HTTP 请求中传递 Trace ID

客户端可以通过 header 传递 trace_id：

```bash
# 自定义 trace_id
curl -H "X-Trace-Id: your-custom-trace-id" http://localhost:8989/api/v1/...

# 或使用 W3C Trace Context (与 OpenTelemetry 兼容)
curl -H "traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01" http://localhost:8989/api/v1/...
```

### 2. 在 Handler 中获取 Trace ID

```go
import (
    "go-admin/server/middleware"
    // 或使用工具包
    "go-admin/server/utils"
)

func SomeHandler(c *gin.Context) {
    // 从 gin context 获取
    traceID := middleware.GetTraceID(c.Request.Context())
    
    // 或使用 utils
    traceID = utils.GetTraceID(c.Request.Context())
    
    // 获取带有 trace 的 logger
    logger := utils.Logger(c.Request.Context())
    logger.Info("handling request", zap.String("data", "..."))
}
```

### 3. 在 Service 层使用

```go
func SomeService(ctx context.Context) error {
    // 获取 trace_id
    traceID := utils.GetTraceID(ctx)
    
    // 使用带 trace 的 logger
    log := utils.SugarLogger(ctx)
    log.Infof("processing with trace_id: %s", traceID)
    
    // 创建 span
    ctx, span := utils.StartSpan(ctx, "SomeService")
    defer span.End()
    
    // 添加属性
    utils.SetSpanAttributes(ctx, map[string]string{
        "user_id": "123",
        "action": "update",
    })
    
    return nil
}
```

### 4. 响应 Header

每个响应会自动包含以下 header：

```
X-Trace-Id: abc123...
X-Request-Id: def456...
```

## 日志示例

启用 trace 后，日志格式如下：

```json
{
  "level": "info",
  "ts": "2024-06-11T10:00:00.000Z",
  "caller": "middleware/middleware.go:23",
  "msg": "http",
  "trace_id": "abc123-def456-ghi789",
  "request_id": "jkl012-mno345-pqr678",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "cost": 0.001234,
  "ip": "127.0.0.1"
}
```

## 架构

### 中间件顺序

```
Recovery -> Trace -> (OpenTelemetry) -> CORS -> RequestLog -> 业务
```

### Trace 传播

- 优先读取 `X-Trace-Id` header
- 然后尝试读取 W3C `traceparent` header
- 都没有则生成新的 UUID
- 同时支持 `X-Request-Id`

## OpenTelemetry

目前支持 stdout 导出器（用于调试），计划支持：

- OTLP 导出器 (gRPC/HTTP)
- Jaeger 直接导出
- Zipkin 导出

## 最佳实践

1. **始终传递 context**: 确保所有函数调用都传递 ctx，以便 trace 传播
2. **使用 utils.Logger**: 代替 global.Logger，自动包含 trace 上下文
3. **创建关键 span**: 在重要的业务逻辑处创建 span
4. **添加有用属性**: 在 span 中添加业务属性（如 user_id、order_id 等）

## 示例

完整示例见 `server/examples/trace` (待添加)
