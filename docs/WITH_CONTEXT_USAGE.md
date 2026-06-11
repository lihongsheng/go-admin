# WithContext 使用指南

## 概述

`WithContext(ctx)` 方法让 logger 自动从 context 中提取 trace_id 和其他追踪信息，无需手动添加这些字段。

## 基本用法

### 1. 在中间件中自动设置

Trace 中间件会自动：
- 从 HTTP header 读取或生成 trace_id
- 将 trace_id 存入 context
- 创建带有 trace_id 的 logger 并也存入 context

```go
// main.go 中设置中间件
logger, _ := initialize.Logger(cfg.Log)

r := gin.New()
r.Use(
    gin.Recovery(),
    middleware.Trace(middleware.WithLogger(logger)),  // 使用 WithLogger 注入 logger
    middleware.Cors(),
    middleware.RequestLog(middleware.WithRequestLogger(logger)),
)
```

### 2. 在 Handler 中使用

```go
func SomeHandler(c *gin.Context) {
    // 方法1: 直接从 context 获取已经带有 trace_id 的 logger
    logger := middleware.LoggerFromContext(c.Request.Context())
    logger.Info("处理请求", "user_id", 123)  // 自动包含 trace_id 和 request_id

    // 方法2: 使用基础 logger 调用 WithContext 手动提取
    baseLogger := someInjectedLogger
    tracedLogger := baseLogger.WithContext(c.Request.Context())
    tracedLogger.Info("另一条日志")  // 也会包含 trace_id
}
```

### 3. 在 Service 层使用

```go
func SomeService(ctx context.Context, data Data) error {
    // 获取带有 trace 的 logger
    logger := log.FromContext(ctx)
    // 或者
    logger := middleware.LoggerFromContext(ctx)

    logger.Info("开始处理", "data_id", data.ID)

    // 使用 Sugar
    sugar := logger.Sugar()
    sugar.Infof("处理 %s 完成", data.Name)

    return nil
}
```

## WithContext 自动提取的字段

| 字段 | 来源 | 说明 |
|------|------|------|
| `trace_id` | `context.Value("trace_id")` | 请求追踪 ID |
| `request_id` | `context.Value("request_id")` | 单次请求 ID |
| `span_id` | OpenTelemetry span context | 分布式追踪 span ID |
| `trace_flags` | OpenTelemetry span context | 追踪标志 |

## 手动设置 trace_id

如果需要手动设置 trace_id（比如在非 HTTP 场景）：

```go
ctx := context.Background()
ctx = log.SetTraceIDToContext(ctx, "my-custom-trace-id-123")
ctx = log.SetRequestIDToContext(ctx, "my-request-id-456")

logger := baseLogger.WithContext(ctx)
logger.Info("手动设置了 trace_id")  // 会包含 trace_id
```

## 日志输出示例

### JSON 格式输出
```json
{
  "level": "info",
  "ts": "2024-06-11T10:30:00.000Z",
  "msg": "用户登录",
  "trace_id": "abc123-def456-ghi789",
  "request_id": "jkl012-mno345-pqr678",
  "span_id": "0123456789abcdef",
  "user_id": 123,
  "ip": "127.0.0.1"
}
```

### 控制台输出
```
INFO	2024-06-11T10:30:00.000Z	用户登录	{"trace_id":"abc123-def456-ghi789","request_id":"jkl012-mno345-pqr678","user_id":123,"ip":"127.0.0.1"}
```

## 链式调用

`WithContext` 可以和 `With` 结合使用：

```go
// 添加业务字段
userLogger := logger.WithContext(ctx).With(
    "user_id", 123,
    "role", "admin",
)

// 之后所有日志都会包含 trace_id、request_id、user_id 和 role
userLogger.Info("用户操作")
userLogger.Warn("敏感操作")
```

## SugaredLogger 也支持 WithContext

```go
// 使用 sugared logger
sugar := logger.Sugar().WithContext(ctx)
sugar.Infof("用户 %d 登录", 123)
```

## 最佳实践

1. **始终传递 context**: 在函数调用链中确保传递 ctx
2. **优先使用 FromContext**: 使用 `log.FromContext(ctx)` 获取 logger 而不是全局 logger
3. **中间件优先**: 让 Trace 中间件自动处理，不要手动设置 trace_id
4. **避免 NopLogger**: 确保正确注入 logger，避免使用 NopLogger 导致日志丢失
5. **结合 OpenTelemetry**: WithContext 会自动提取 OTel 的 span 信息
